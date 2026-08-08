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
	import { ApiError, backend, type CheckResult, type MailAccount } from '$lib/backend';
	import { spaceStore } from '$lib/stores/space.svelte';
	import AccountDrawer from '$lib/components/settings/AccountDrawer.svelte';

	type Probe =
		| { state: 'checking' }
		| { state: 'checked'; result: CheckResult }
		| { state: 'throttled'; retryAt: number }
		| { state: 'unreachable'; error: string };
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
	 * Nothing in the schema remembers a handshake, so "not checked yet" is a real state rather
	 * than a euphemism. The two protocols are reported separately because they fail
	 * separately — an account whose IMAP is fine and whose SMTP is refusing can read mail and
	 * not send it, and a single dot for both can only ever describe one of them.
	 */
	function legs(account: MailAccount): { tone: Tone; label: string; pulse: boolean }[] {
		const probe = probes[account.id];

		if (probe?.state === 'checking') {
			return [{ tone: 'warning', label: 'Checking…', pulse: true }];
		}
		if (probe?.state === 'throttled') {
			return [{ tone: 'warning', label: `Not checked — ${retryIn(probe)}`, pulse: false }];
		}
		if (probe?.state === 'unreachable') {
			return [{ tone: 'danger', label: 'Courrier could not run the check', pulse: false }];
		}

		return (['imap', 'smtp'] as const).map((protocol) => {
			const name = protocol.toUpperCase();
			const host = protocol === 'imap' ? account.imap_host : account.smtp_host;
			if (!host) return { tone: 'neutral' as Tone, label: `No ${name} server set`, pulse: false };

			const leg = probe?.state === 'checked' ? probe.result[protocol] : undefined;
			if (!leg) return { tone: 'neutral' as Tone, label: `${name} not checked yet`, pulse: false };
			if (leg.ok) return { tone: 'success' as Tone, label: `${name} connected`, pulse: false };
			return { tone: 'danger' as Tone, label: `${name} ${failureKind(leg.error)}`, pulse: false };
		});
	}

	/*
	 * "Refused the connection" was wrong for the commonest failure there is: with a bad
	 * password the connection is established and it is the *login* that fails. Saying so is
	 * the difference between checking the port and checking the password.
	 */
	function failureKind(error?: string) {
		return /auth|password|credential|login|invalid/i.test(error ?? '')
			? 'rejected the credentials'
			: 'unreachable';
	}

	/*
	 * A 429 never reached the mail server, so it says nothing about the account: the legs stay
	 * unknown and the wait is the whole message. The limiter sends `Retry-After`, and counting
	 * it down from a ticking clock is what stops the label going stale as it sits on screen.
	 */
	let now = $state(Date.now());

	$effect(() => {
		if (!Object.values(probes).some((probe) => probe.state === 'throttled')) return;
		const timer = setInterval(() => {
			now = Date.now();
			for (const [id, probe] of Object.entries(probes)) {
				if (probe.state === 'throttled' && probe.retryAt <= now) delete probes[Number(id)];
			}
		}, 1000);
		return () => clearInterval(timer);
	});

	function retryIn(probe: Extract<Probe, { state: 'throttled' }>) {
		const seconds = Math.max(0, Math.ceil((probe.retryAt - now) / 1000));
		return seconds > 0 ? `retry in ${seconds}s` : 'you can retry now';
	}

	/*
	 * A failed handshake comes back 200 with the reason in the body — the diagnostic worked,
	 * the server did not — so a thrown error here means something else entirely: the rate
	 * limiter, an expired session, a network drop. Those are not the account's fault and must
	 * not be reported as though they were.
	 */
	async function check(account: MailAccount) {
		probes[account.id] = { state: 'checking' };
		try {
			probes[account.id] = { state: 'checked', result: await backend.checkAccount(account.id) };
		} catch (err) {
			if (err instanceof ApiError && err.status === 429) {
				probes[account.id] = {
					state: 'throttled',
					retryAt: Date.now() + (err.retryAfter ?? 60) * 1000
				};
				return;
			}
			probes[account.id] = {
				state: 'unreachable',
				error: err instanceof Error ? err.message : 'The check could not be run'
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

	/*
	 * The dot names the kind of failure; the server's own words are what someone pastes into a
	 * search or sends to their host, so they belong here in full, per protocol.
	 */
	const failures = $derived(
		accounts.flatMap((account) => {
			const probe = probes[account.id];
			if (probe?.state === 'unreachable') {
				return [{ key: `${account.id}-run`, account, protocol: '', error: probe.error }];
			}
			if (probe?.state !== 'checked') return [];
			return (['imap', 'smtp'] as const)
				.filter((protocol) => probe.result[protocol].error)
				.map((protocol) => ({
					key: `${account.id}-${protocol}`,
					account,
					protocol: protocol.toUpperCase(),
					error: probe.result[protocol].error as string
				}));
		})
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
						{@const state = legs(account)}
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
								<div class="flex flex-col gap-1">
									{#each state as leg (leg.label)}
										<StatusDot tone={leg.tone} label={leg.label} pulse={leg.pulse} />
									{/each}
								</div>
							</td>
							<td>
								<div class="flex flex-wrap items-center justify-end gap-1">
									<Button
										variant="ghost"
										size="sm"
										icon={icons.plug}
										disabled={probes[account.id]?.state === 'checking' ||
										probes[account.id]?.state === 'throttled' ||
										(!account.imap_host && !account.smtp_host)}
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
				Check signs in to IMAP and SMTP with the stored passwords and hangs up again — it reads
				no mailboxes and downloads no messages. Each protocol is reported on its own, because
				an account can read mail perfectly well and still be unable to send it.
			</p>
		{/if}

		{#each failures as failure (failure.key)}
			<Alert
				tone="danger"
				title={failure.protocol
					? `${failure.account.name} — ${failure.protocol} failed`
					: `${failure.account.name} could not be checked`}
			>
				{failure.error}
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
