<script lang="ts">
	import { getContext } from 'svelte';
	import {
		Alert,
		Badge,
		Button,
		ConfirmModal,
		EmptyState,
		SettingsSection,
		StatusDot,
		Table,
		icons
	} from '@facile/muse';
	import { backend, type MailAccount } from '$lib/backend';
	import { spaceStore } from '$lib/stores/space.svelte';
	import AccountDrawer from '$lib/components/settings/AccountDrawer.svelte';

	type ProbeState = 'unchecked' | 'checking' | 'connected' | 'failed';
	type Probe = { state: ProbeState; error?: string };
	type Tone = 'neutral' | 'success' | 'warning' | 'danger';

	const app = getContext<{ refreshAccounts: () => Promise<void> }>('app');

	let accounts = $state<MailAccount[]>([]);
	let loadError = $state('');

	let probes = $state<Record<number, Probe>>({});

	let editorOpen = $state(false);
	let editing = $state<MailAccount | null>(null);

	let removeTarget = $state<MailAccount | null>(null);
	let removeOpen = $state(false);
	let removeError = $state('');

	let promoting = $state<number | null>(null);

	$effect(() => {
		spaceStore.active; // reload when the active space changes
		void load();
	});

	async function load() {
		try {
			accounts = (await backend.listAccounts(spaceStore.active?.id)).accounts;
			loadError = '';
		} catch (err) {
			accounts = [];
			loadError = err instanceof Error ? err.message : 'Could not load your mail accounts';
		}
	}

	/*
	 * Connection state is not a boolean, and it is not stored either: nothing in the schema
	 * remembers whether the last IMAP handshake worked. So "unchecked" is a real, distinct
	 * state — it says nobody has asked this server anything yet — and an account missing a
	 * host is a fifth one, because the fix for it is a form field rather than a network.
	 */
	function status(account: MailAccount): { tone: Tone; label: string; pulse: boolean } {
		if (!account.imap_host || !account.smtp_host) {
			const missing = !account.imap_host ? 'IMAP' : 'SMTP';
			return { tone: 'neutral', label: `No ${missing} server set`, pulse: false };
		}
		const probe = probes[account.id];
		switch (probe?.state) {
			case 'checking':
				return { tone: 'warning', label: 'Connecting to IMAP…', pulse: true };
			case 'connected':
				return { tone: 'success', label: 'IMAP connected', pulse: false };
			case 'failed':
				return { tone: 'danger', label: 'IMAP refused the connection', pulse: false };
			default:
				return { tone: 'neutral', label: 'Not checked yet', pulse: false };
		}
	}

	/*
	 * The stored password never comes back over the wire, so /mail/test-connection cannot be
	 * used on a saved account. A folder sync is the probe that can: it opens IMAP with the
	 * decrypted credentials and lists mailboxes without fetching a single message.
	 */
	async function check(account: MailAccount) {
		probes[account.id] = { state: 'checking' };
		try {
			await backend.syncAccount(account.id);
			probes[account.id] = { state: 'connected' };
			await app.refreshAccounts();
		} catch (err) {
			probes[account.id] = {
				state: 'failed',
				error: err instanceof Error ? err.message : 'The connection failed'
			};
		}
	}

	async function makeDefault(account: MailAccount) {
		promoting = account.id;
		try {
			await backend.updateAccount(account.id, { is_default: true });
			await load();
			await app.refreshAccounts();
		} catch (err) {
			loadError = err instanceof Error ? err.message : 'Could not change the default account';
		}
		promoting = null;
	}

	function openNew() {
		editing = null;
		editorOpen = true;
	}

	function openEdit(account: MailAccount) {
		editing = account;
		editorOpen = true;
	}

	async function remove() {
		const target = removeTarget;
		if (!target) return;
		removeError = '';
		try {
			await backend.deleteAccount(target.id);
			delete probes[target.id];
			removeTarget = null;
			await load();
			await app.refreshAccounts();
		} catch (err) {
			removeError = err instanceof Error ? err.message : 'Could not remove the account';
		}
	}

	const failures = $derived(
		accounts
			.map((account) => ({ account, probe: probes[account.id] }))
			.filter((entry) => entry.probe?.state === 'failed' && entry.probe.error)
	);
</script>

<div class="flex flex-col gap-10">
	<SettingsSection
		title="Mail accounts"
		description={spaceStore.active
			? `The IMAP and SMTP servers Courrier reads and sends through in ${spaceStore.active.name}.`
			: 'The IMAP and SMTP servers Courrier reads and sends through.'}
		bare
	>
		{#snippet actions()}
			<Button icon={icons.plus} onclick={openNew}>Add account</Button>
		{/snippet}

		{#if loadError}
			<Alert tone="danger">{loadError}</Alert>
		{/if}

		{#if accounts.length === 0}
			<EmptyState
				icon={icons.mail}
				title="No mail accounts yet"
				description="Courrier has nothing to read until you point it at an IMAP server. Credentials are encrypted at rest."
			>
				<Button icon={icons.plus} onclick={openNew}>Add your first account</Button>
			</EmptyState>
		{:else}
			<Table>
				<thead>
					<tr>
						<th>Account</th>
						<th>Servers</th>
						<th>Connection</th>
						<th class="text-right">Actions</th>
					</tr>
				</thead>
				<tbody>
					{#each accounts as account (account.id)}
						{@const state = status(account)}
						<tr>
							<td>
								<div class="flex min-w-0 flex-col gap-1">
									<span class="flex items-center gap-2 font-medium text-fc-fg">
										{account.name}
										{#if account.is_default}<Badge tone="accent">Default</Badge>{/if}
									</span>
									<span class="text-fc-xs text-fc-fg-muted">{account.email}</span>
								</div>
							</td>
							<td class="font-fc-mono text-fc-xs whitespace-nowrap text-fc-fg-muted">
								<div class="flex flex-col gap-1">
									<span>{account.imap_host || '—'}:{account.imap_port}</span>
									<span>{account.smtp_host || '—'}:{account.smtp_port}</span>
								</div>
							</td>
							<td>
								<StatusDot tone={state.tone} label={state.label} pulse={state.pulse} />
							</td>
							<td>
								<div class="flex flex-wrap items-center justify-end gap-1">
									<Button
										variant="ghost"
										size="sm"
										icon={icons.plug}
										disabled={probes[account.id]?.state === 'checking' || !account.imap_host}
										onclick={() => check(account)}
									>
										Check
									</Button>
									{#if !account.is_default}
										<Button
											variant="ghost"
											size="sm"
											icon={icons.check}
											disabled={promoting === account.id}
											onclick={() => makeDefault(account)}
										>
											Make default
										</Button>
									{/if}
									<Button
										variant="ghost"
										size="sm"
										icon={icons.edit}
										onclick={() => openEdit(account)}
									>
										Edit
									</Button>
									<Button
										variant="ghost-danger"
										size="sm"
										icon={icons.remove}
										onclick={() => {
											removeTarget = account;
											removeError = '';
											removeOpen = true;
										}}
									>
										Remove
									</Button>
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</Table>

			<p class="text-fc-xs text-fc-fg-muted">
				Check opens IMAP with the stored password and refreshes the folder list — it never
				downloads messages. SMTP is only reachable from Test connection inside the account, which
				needs the passwords typed in.
			</p>
		{/if}

		{#each failures as failure (failure.account.id)}
			<Alert tone="danger" title={`${failure.account.name} could not connect`}>
				{failure.probe?.error}
			</Alert>
		{/each}

		{#if removeError}
			<Alert tone="danger">{removeError}</Alert>
		{/if}
	</SettingsSection>
</div>

<AccountDrawer
	bind:open={editorOpen}
	account={editing}
	isFirst={accounts.length === 0}
	onSaved={async () => {
		await load();
		await app.refreshAccounts();
	}}
/>

<ConfirmModal
	bind:open={removeOpen}
	tone="danger"
	title="Remove this mail account?"
	description={`Courrier deletes ${removeTarget?.name ?? 'this account'}, its stored credentials, and every folder, message and draft it has cached locally. Your mail itself stays on ${removeTarget?.imap_host || 'the server'} — but re-adding the account means syncing it all again from scratch.`}
	confirmLabel="Remove account"
	cancelLabel="Keep it"
	onConfirm={remove}
/>
