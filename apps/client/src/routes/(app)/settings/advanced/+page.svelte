<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import {
		Alert,
		Button,
		ConfirmModal,
		SecretField,
		SettingsRow,
		SettingsSection,
		StatusDot,
		icons
	} from '@facile/muse';
	import { backend, type AuthConfig, type MailAccount, type UserProfile } from '$lib/backend';
	import { spaceStore } from '$lib/stores/space.svelte';

	const app = getContext<{
		user: UserProfile | null;
		refreshAccounts: () => Promise<void>;
	}>('app');

	let config = $state<AuthConfig | null>(null);
	let encryptionKeySet = $state<boolean | null>(null);
	let accounts = $state<MailAccount[]>([]);

	let wipeOpen = $state(false);
	let wiping = $state(false);
	let wipeError = $state('');
	let wiped = $state(0);

	const endpoint = `${backend.baseUrl || (typeof location === 'undefined' ? '' : location.origin)}/api`;

	onMount(async () => {
		try {
			config = await backend.authConfig();
		} catch {
			config = null;
		}
		try {
			encryptionKeySet = (await backend.instanceSettings()).encryption_key_set;
		} catch {
			encryptionKeySet = null;
		}
		await loadAccounts();
	});

	async function loadAccounts() {
		try {
			accounts = (await backend.listAccounts(spaceStore.active?.id)).accounts;
		} catch {
			accounts = [];
		}
	}

	async function wipeAccounts() {
		wiping = true;
		wipeError = '';
		wiped = 0;
		for (const account of accounts) {
			try {
				await backend.deleteAccount(account.id);
				wiped += 1;
			} catch (err) {
				wipeError = err instanceof Error ? err.message : 'Could not remove every account';
				break;
			}
		}
		await loadAccounts();
		await app.refreshAccounts();
		wiping = false;
	}
</script>

<div class="flex flex-col gap-10">
	<SettingsSection
		title="Instance"
		description="The facts worth quoting when you file a bug against a self-hosted install."
	>
		<SettingsRow label="API base URL" description="Same origin as this page in production." stacked>
			<SecretField value={endpoint} sensitive={false} class="w-full" />
		</SettingsRow>

		<SettingsRow label="Your user ID" description="What the API knows you as." stacked>
			<SecretField value={app.user?.id ?? '—'} sensitive={false} class="w-full" />
		</SettingsRow>

		<SettingsRow
			label="Single sign-on"
			description="OIDC is optional: when it is off, this instance only knows local passwords."
		>
			{#if config === null}
				<StatusDot tone="neutral" label="Unknown" />
			{:else if config.oidc_enabled}
				<StatusDot tone="success" label="Federated over OIDC" />
			{:else}
				<StatusDot tone="neutral" label="Not configured" />
			{/if}
		</SettingsRow>

		<SettingsRow
			label="Password sign-in"
			description="SSO_ONLY hides the password form and closes registration."
		>
			{#if config === null}
				<StatusDot tone="neutral" label="Unknown" />
			{:else if config.sso_only}
				<StatusDot tone="neutral" label="Disabled — SSO only" />
			{:else}
				<StatusDot tone="success" label="Enabled" />
			{/if}
		</SettingsRow>

		<SettingsRow
			label="Credential encryption"
			description="ENCRYPTION_KEY is what keeps your IMAP and SMTP passwords unreadable in the database."
		>
			{#if encryptionKeySet === null}
				<StatusDot tone="neutral" label="Unknown" />
			{:else if encryptionKeySet}
				<StatusDot tone="success" label="Key configured" />
			{:else}
				<StatusDot tone="danger" label="No key set" />
			{/if}
		</SettingsRow>

		{#if encryptionKeySet === false}
			<Alert tone="danger" title="Mail credentials are not encrypted at rest">
				This instance is running without an ENCRYPTION_KEY. Set one in the server environment and
				restart before you add another mail account.
			</Alert>
		{/if}
	</SettingsSection>

	<!--
		Danger zone lives at the bottom of the last tab, never as a tab of its own: a destructive
		action should cost a scroll instead of sitting one mis-click from every other section.
	-->
	<SettingsSection title="Danger zone" description="Irreversible, and nobody can undo it for you.">
		<SettingsRow
			label="Disconnect every mail account"
			description={accounts.length === 0
				? 'Nothing to disconnect — there are no mail accounts in this space.'
				: `Removes all ${accounts.length} account${accounts.length === 1 ? '' : 's'} in this space, along with everything Courrier has cached from them.`}
		>
			<Button
				variant="danger"
				icon={icons.remove}
				disabled={wiping || accounts.length === 0}
				onclick={() => (wipeOpen = true)}
			>
				{wiping ? 'Disconnecting…' : 'Disconnect all'}
			</Button>
		</SettingsRow>

		{#if wipeError}
			<Alert tone="danger" title={`Stopped after ${wiped} account${wiped === 1 ? '' : 's'}`}>
				{wipeError}
			</Alert>
		{:else if wiped > 0}
			<Alert tone="success">
				Disconnected {wiped} account{wiped === 1 ? '' : 's'}. Your mail is untouched on the server.
			</Alert>
		{/if}
	</SettingsSection>
</div>

<ConfirmModal
	bind:open={wipeOpen}
	tone="danger"
	title="Disconnect every mail account?"
	description={`Courrier deletes all ${accounts.length} account${accounts.length === 1 ? '' : 's'} in this space, their stored IMAP and SMTP credentials, and every folder, message and draft cached from them — the mail views go empty. Your mail itself stays on the servers, but adding the accounts back means re-entering every password and syncing from scratch.`}
	confirmLabel="Disconnect all"
	cancelLabel="Leave them alone"
	onConfirm={wipeAccounts}
/>
