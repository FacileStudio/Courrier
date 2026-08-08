<script lang="ts">
	import {
		Alert,
		Button,
		Divider,
		Drawer,
		Field,
		Input,
		REDACTED,
		SecretField,
		Spinner,
		StatusDot,
		Switch,
		Textarea,
		isRedacted,
		icons
	} from '@facile/muse';
	import { backend, type MailAccount } from '$lib/backend';
	import { spaceStore } from '$lib/stores/space.svelte';

	let {
		open = $bindable(false),
		account = null,
		isFirst = false,
		onSaved
	}: {
		open?: boolean;
		account?: MailAccount | null;
		isFirst?: boolean;
		onSaved?: () => void;
	} = $props();

	let name = $state('');
	let email = $state('');
	let imapHost = $state('');
	let imapPort = $state('993');
	let imapUser = $state('');
	let imapPassword = $state('');
	let smtpHost = $state('');
	let smtpPort = $state('587');
	let smtpUser = $state('');
	let smtpPassword = $state('');
	let signature = $state('');
	let isDefault = $state(false);

	let saving = $state(false);
	let error = $state('');
	let testState = $state<'idle' | 'testing' | 'passed' | 'failed'>('idle');
	let testError = $state('');

	/*
	 * The drawer is mounted once and reused, so it repopulates whenever it is opened rather
	 * than on mount — otherwise editing the second account shows the first one's server.
	 */
	let loadedKey = $state<string>('');
	const key = $derived(String(account?.id ?? 'new'));

	$effect(() => {
		if (!open) {
			loadedKey = '';
			return;
		}
		if (key === loadedKey) return;
		loadedKey = key;
		error = '';
		testState = 'idle';
		testError = '';
		name = account?.name ?? '';
		email = account?.email ?? '';
		imapHost = account?.imap_host ?? '';
		imapPort = String(account?.imap_port ?? 993);
		imapUser = account?.imap_user ?? '';
		smtpHost = account?.smtp_host ?? '';
		smtpPort = String(account?.smtp_port ?? 587);
		smtpUser = account?.smtp_user ?? '';
		signature = account?.signature ?? '';
		isDefault = account?.is_default ?? isFirst;
		/*
		 * The API never sends a stored password back, so an existing account gets the eight-dot
		 * placeholder: SecretField goes inert on it, and a save that leaves it untouched omits
		 * the field entirely rather than overwriting a live credential with punctuation.
		 */
		imapPassword = account ? REDACTED : '';
		smtpPassword = account ? REDACTED : '';
	});

	const editing = $derived(account !== null);
	const typedBoth = $derived(
		imapPassword.length > 0 &&
			!isRedacted(imapPassword) &&
			smtpPassword.length > 0 &&
			!isRedacted(smtpPassword)
	);
	const complete = $derived(
		name.trim().length > 0 &&
			email.trim().length > 0 &&
			imapHost.trim().length > 0 &&
			smtpHost.trim().length > 0
	);

	function port(value: string, fallback: number) {
		const parsed = Number.parseInt(value, 10);
		return Number.isFinite(parsed) && parsed > 0 ? parsed : fallback;
	}

	async function test() {
		testState = 'testing';
		testError = '';
		try {
			await backend.testConnection({
				imap_host: imapHost.trim(),
				imap_port: port(imapPort, 993),
				imap_user: imapUser.trim() || email.trim(),
				imap_password: imapPassword,
				smtp_host: smtpHost.trim(),
				smtp_port: port(smtpPort, 587),
				smtp_user: smtpUser.trim() || email.trim(),
				smtp_password: smtpPassword
			});
			testState = 'passed';
		} catch (err) {
			testState = 'failed';
			testError = err instanceof Error ? err.message : 'Connection test failed';
		}
	}

	async function save(event: Event) {
		event.preventDefault();
		saving = true;
		error = '';
		try {
			const common = {
				name: name.trim(),
				email: email.trim(),
				imap_host: imapHost.trim(),
				imap_port: port(imapPort, 993),
				imap_user: imapUser.trim() || email.trim(),
				smtp_host: smtpHost.trim(),
				smtp_port: port(smtpPort, 587),
				smtp_user: smtpUser.trim() || email.trim(),
				signature,
				is_default: isDefault
			};

			if (account) {
				await backend.updateAccount(account.id, {
					...common,
					...(isRedacted(imapPassword) || imapPassword.length === 0
						? {}
						: { imap_password: imapPassword }),
					...(isRedacted(smtpPassword) || smtpPassword.length === 0
						? {}
						: { smtp_password: smtpPassword })
				});
			} else {
				await backend.createAccount({
					...common,
					imap_password: imapPassword,
					smtp_password: smtpPassword,
					space_id: spaceStore.active?.id
				});
			}
			open = false;
			onSaved?.();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Could not save the account';
		}
		saving = false;
	}
</script>

<Drawer bind:open title={editing ? `Edit ${account?.name}` : 'Add a mail account'}>
	<form class="flex flex-col gap-4" onsubmit={save}>
		<Field label="Display name" helper="What this inbox is called inside Courrier.">
			<Input bind:value={name} placeholder="Work" required disabled={saving} />
		</Field>

		<Field label="Email address" helper="The address messages are sent from.">
			<Input bind:value={email} type="email" placeholder="you@example.com" required disabled={saving} />
		</Field>

		<Divider class="my-1" />

		<p class="text-fc-sm font-medium text-fc-fg">IMAP — reading mail</p>

		<div class="grid gap-4 sm:grid-cols-[1fr_7rem]">
			<Field label="Host">
				<Input bind:value={imapHost} placeholder="imap.example.com" required disabled={saving} />
			</Field>
			<Field label="Port">
				<Input bind:value={imapPort} inputmode="numeric" disabled={saving} />
			</Field>
		</div>

		<Field label="Username" helper="Left empty, the email address is used.">
			<Input bind:value={imapUser} placeholder={email || 'you@example.com'} disabled={saving} />
		</Field>

		<SecretField
			bind:value={imapPassword}
			editable
			label="IMAP password"
			helper={isRedacted(imapPassword)
				? 'Already stored, encrypted at rest. Type over it to replace it; leave it to keep it.'
				: 'Encrypted at rest with this instance’s ENCRYPTION_KEY.'}
			disabled={saving}
		/>

		<Divider class="my-1" />

		<p class="text-fc-sm font-medium text-fc-fg">SMTP — sending mail</p>

		<div class="grid gap-4 sm:grid-cols-[1fr_7rem]">
			<Field label="Host">
				<Input bind:value={smtpHost} placeholder="smtp.example.com" required disabled={saving} />
			</Field>
			<Field label="Port">
				<Input bind:value={smtpPort} inputmode="numeric" disabled={saving} />
			</Field>
		</div>

		<Field label="Username" helper="Left empty, the email address is used.">
			<Input bind:value={smtpUser} placeholder={email || 'you@example.com'} disabled={saving} />
		</Field>

		<SecretField
			bind:value={smtpPassword}
			editable
			label="SMTP password"
			helper={isRedacted(smtpPassword)
				? 'Already stored, encrypted at rest. Type over it to replace it; leave it to keep it.'
				: 'Often the same as the IMAP password, but not always.'}
			disabled={saving}
		/>

		<Divider class="my-1" />

		<Field label="Signature" helper="Appended to messages you send from this account.">
			<Textarea bind:value={signature} rows={4} placeholder="— Your name" disabled={saving} />
		</Field>

		<div class="flex flex-col gap-1">
			<Switch bind:checked={isDefault} label="Use as the default account" disabled={saving} />
			<p class="pl-14 text-fc-xs text-fc-fg-muted">
				The account Courrier opens on and composes from unless you pick another.
			</p>
		</div>

		<Divider class="my-1" />

		<div class="flex flex-col gap-3">
			<div class="flex flex-wrap items-center gap-2">
				<Button
					type="button"
					variant="outline"
					icon={icons.plug}
					disabled={saving || testState === 'testing' || !imapHost.trim() || !typedBoth}
					onclick={test}
				>
					{testState === 'testing' ? 'Testing…' : 'Test connection'}
				</Button>

				{#if testState === 'testing'}
					<StatusDot tone="warning" label="Connecting to both servers…" pulse />
				{:else if testState === 'passed'}
					<StatusDot tone="success" label="IMAP and SMTP both answered" />
				{:else if testState === 'failed'}
					<StatusDot tone="danger" label="Test failed" />
				{/if}
			</div>

			{#if !typedBoth}
				<p class="text-fc-xs text-fc-fg-muted">
					Testing sends both passwords to the servers, so it needs them typed in. For an account
					already saved, use Check on its row instead — that one uses the stored credentials.
				</p>
			{/if}

			{#if testState === 'failed' && testError}
				<Alert tone="danger" title="Connection test failed">{testError}</Alert>
			{/if}
		</div>

		{#if error}
			<Alert tone="danger">{error}</Alert>
		{/if}

		<div class="flex justify-end gap-2 pt-1">
			<Button type="button" variant="ghost" disabled={saving} onclick={() => (open = false)}>
				Cancel
			</Button>
			<Button type="submit" disabled={saving || !complete}>
				{#if saving}<Spinner size="sm" />{/if}
				{saving ? 'Saving…' : editing ? 'Save account' : 'Add account'}
			</Button>
		</div>
	</form>
</Drawer>
