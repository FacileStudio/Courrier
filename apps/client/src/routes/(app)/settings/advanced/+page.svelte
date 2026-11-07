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
	import { t } from '$lib/i18n/index.svelte';

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
				wipeError = err instanceof Error ? err.message : t('advanced.wipeFailed');
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
		title={t('advanced.instance')}
		description={t('advanced.instanceBody')}
	>
		<SettingsRow label={t('advanced.apiBase')} description={t('advanced.apiBaseBody')} stacked>
			<SecretField value={endpoint} sensitive={false} class="w-full" />
		</SettingsRow>

		<SettingsRow label={t('advanced.userId')} description={t('advanced.userIdBody')} stacked>
			<SecretField value={app.user?.id ?? '—'} sensitive={false} class="w-full" />
		</SettingsRow>

		<SettingsRow label={t('advanced.signIn')} description={t('advanced.signInBody')}>
			{#if config === null}
				<StatusDot tone="neutral" label={t('advanced.unknown')} />
			{:else if config.oidc_enabled}
				<StatusDot tone="success" label={t('advanced.federated')} />
			{:else}
				<StatusDot tone="neutral" label={t('advanced.notConfigured')} />
			{/if}
		</SettingsRow>

		<SettingsRow label={t('advanced.encryption')} description={t('advanced.encryptionBody')}>
			{#if encryptionKeySet === null}
				<StatusDot tone="neutral" label={t('advanced.unknown')} />
			{:else if encryptionKeySet}
				<StatusDot tone="success" label={t('advanced.keyConfigured')} />
			{:else}
				<StatusDot tone="danger" label={t('advanced.noKey')} />
			{/if}
		</SettingsRow>

		{#if encryptionKeySet === false}
			<Alert tone="danger" title={t('advanced.noKeyTitle')}>
				{t('advanced.noKeyBody')}
			</Alert>
		{/if}
	</SettingsSection>

	<!--
		Danger zone lives at the bottom of the last tab, never as a tab of its own: a destructive
		action should cost a scroll instead of sitting one mis-click from every other section.
	-->
	<SettingsSection title={t('advanced.danger')} description={t('advanced.dangerBody')}>
		<SettingsRow
			label={t('advanced.disconnectAll')}
			description={accounts.length === 0
				? t('advanced.disconnectNone')
				: t('advanced.disconnectBody', { count: accounts.length })}
		>
			<Button
				variant="danger"
				icon={icons.remove}
				disabled={wiping || accounts.length === 0}
				onclick={() => (wipeOpen = true)}
			>
				{wiping ? t('advanced.disconnecting') : t('advanced.disconnectAllButton')}
			</Button>
		</SettingsRow>

		{#if wipeError}
			<Alert tone="danger" title={t('advanced.stoppedAfter', { count: wiped })}>
				{wipeError}
			</Alert>
		{:else if wiped > 0}
			<Alert tone="success">{t('advanced.disconnected', { count: wiped })}</Alert>
		{/if}
	</SettingsSection>
</div>

<ConfirmModal
	bind:open={wipeOpen}
	tone="danger"
	title={t('advanced.wipeTitle')}
	description={t('advanced.wipeBody', { count: accounts.length })}
	confirmLabel={t('advanced.disconnectAllButton')}
	cancelLabel={t('advanced.leaveThem')}
	onConfirm={wipeAccounts}
/>
