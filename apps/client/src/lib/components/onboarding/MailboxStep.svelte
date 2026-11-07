<script lang="ts">
	import { Alert, Button, OptionCards, Spinner, StatusDot, icons } from '@facile/muse';
	import type { MailAccount } from '$lib/backend';
	import AccountDrawer from '$lib/components/settings/AccountDrawer.svelte';
	import { t } from '$lib/i18n/index.svelte';

	let {
		accounts,
		loading,
		choice = $bindable(''),
		drawerOpen = $bindable(false),
		onSaved
	}: {
		accounts: MailAccount[];
		loading: boolean;
		choice?: string;
		drawerOpen?: boolean;
		onSaved: () => void;
	} = $props();

	const options = $derived(
		accounts.map((account) => ({ value: String(account.id), label: account.email, icon: icons.mail }))
	);
</script>

<div class="flex w-full max-w-fc-sm flex-col gap-4 text-left" data-reveal>
	{#if loading}
		<div class="flex justify-center"><Spinner /></div>
	{:else if accounts.length === 0}
		<Alert tone="info" title={t('onboarding.noMailboxTitle')}>
			{t('onboarding.noMailboxBody')}
		</Alert>
	{:else if accounts.length === 1}
		<div class="flex items-center justify-between gap-3 rounded-fc-md bg-fc-component p-4">
			<span class="min-w-0 truncate text-fc-sm text-fc-fg">{accounts[0].email}</span>
			<StatusDot tone="success" label={t('onboarding.connected')} />
		</div>
	{:else}
		<OptionCards {options} bind:value={choice} name="default-mailbox" label={t('onboarding.firstMailbox')} />
	{/if}
	{#if !loading && accounts.length > 0}
		<Button variant="ghost" class="min-h-11 self-start" icon={icons.plus} onclick={() => (drawerOpen = true)}>
			{t('onboarding.addAnother')}
		</Button>
	{/if}
</div>

<AccountDrawer bind:open={drawerOpen} isFirst={accounts.length === 0} {onSaved} />
