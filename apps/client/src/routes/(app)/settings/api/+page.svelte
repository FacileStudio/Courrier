<script lang="ts">
	import { onMount } from 'svelte';
	import {
		Alert,
		Button,
		ConfirmModal,
		Drawer,
		Field,
		Input,
		SecretField,
		SettingsSection,
		Spinner,
		Table,
		icons
	} from '@facile/muse';
	import { backend, type ApiTokenStatus } from '$lib/backend';
	import { formatDate } from '$lib/i18n/format';
	import { t } from '$lib/i18n/index.svelte';

	let status = $state<ApiTokenStatus>({ has_token: false });
	let loadError = $state('');

	let createOpen = $state(false);
	let creating = $state(false);
	let createError = $state('');
	let newName = $state('');
	let issued = $state('');

	let revokeOpen = $state(false);
	let revokeError = $state('');

	const endpoint = `${backend.baseUrl || (typeof location === 'undefined' ? '' : location.origin)}/api`;

	onMount(load);

	async function load() {
		try {
			status = await backend.apiTokenStatus();
			loadError = '';
		} catch (err) {
			loadError = err instanceof Error ? err.message : t('api.readFailed');
		}
	}

	/* Reopening the drawer must never re-show the token from a previous run. */
	function openCreate() {
		issued = '';
		newName = '';
		createError = '';
		createOpen = true;
	}

	async function create(event: Event) {
		event.preventDefault();
		creating = true;
		createError = '';
		try {
			issued = (await backend.createApiToken(newName.trim())).token;
			await load();
		} catch (err) {
			createError = err instanceof Error ? err.message : t('api.createFailed');
		}
		creating = false;
	}

	async function revoke() {
		revokeError = '';
		try {
			await backend.deleteApiToken();
			await load();
		} catch (err) {
			revokeError = err instanceof Error ? err.message : t('api.revokeFailed');
		}
	}

	function when(raw?: string) {
		if (!raw) return '—';
		const date = new Date(raw);
		return Number.isNaN(date.getTime()) ? '—' : formatDate(date);
	}
</script>

<div class="flex flex-col gap-10">
	<SettingsSection
		title={t('api.endpoint')}
		description={t('api.endpointBody')}
	>
		<SecretField value={endpoint} sensitive={false} label={t('api.baseUrl')} />
	</SettingsSection>

	<SettingsSection
		title={t('api.tokenTitle')}
		description={t('api.tokenBody')}
		bare
	>
		{#snippet actions()}
			<Button icon={icons.plus} onclick={openCreate}>
				{status.has_token ? t('api.replace') : t('api.new')}
			</Button>
		{/snippet}

		{#if loadError}
			<Alert tone="danger">{loadError}</Alert>
		{/if}

		{#if !status.has_token}
			<Alert tone="info">
				{t('api.none')}
			</Alert>
		{:else}
			<Table>
				<thead>
					<tr>
						<th>{t('api.colName')}</th>
						<th>{t('api.colToken')}</th>
						<th>{t('api.colCreated')}</th>
						<th class="text-right">{t('api.colActions')}</th>
					</tr>
				</thead>
				<tbody>
					<tr>
						<td class="font-medium text-fc-fg">{status.name || 'CLI'}</td>
						<td class="font-fc-mono text-fc-xs text-fc-fg-muted">{t('api.storedHashed')}</td>
						<td class="whitespace-nowrap text-fc-fg-muted">{when(status.created_at)}</td>
						<td class="text-right">
							<Button
								variant="ghost-danger"
								size="sm"
								icon={icons.revoke}
								onclick={() => {
									revokeError = '';
									revokeOpen = true;
								}}
							>
								{t('common.revoke')}
							</Button>
						</td>
					</tr>
				</tbody>
			</Table>
		{/if}

		{#if revokeError}
			<Alert tone="danger">{revokeError}</Alert>
		{/if}
	</SettingsSection>
</div>

<Drawer bind:open={createOpen} title={status.has_token && !issued ? t('api.drawerReplace') : t('api.drawerNew')}>
	{#if issued}
		<div class="flex flex-col gap-4">
			<Alert tone="warning" title={t('api.copyNow')}>
				{t('api.copyNowBody')}
			</Alert>

			<!--
				The one-time token is the exception to the auto-hide rule: it starts revealed and stays
				that way, because hiding a value nobody has copied yet is theatre.
			-->
			<SecretField
				value={issued}
				visible
				autoHideMs={0}
				label={t('api.tokenLabel')}
				helper={t('api.tokenHelp')}
			/>

			<div class="flex justify-end">
				<Button onclick={() => (createOpen = false)}>{t('api.done')}</Button>
			</div>
		</div>
	{:else}
		<form class="flex flex-col gap-4" onsubmit={create}>
			{#if status.has_token}
				<Alert tone="warning" title={t('api.replacesTitle')}>
					{t('api.replacesBody', { name: status.name || 'CLI' })}
				</Alert>
			{/if}

			<Field label={t('api.nameLabel')} helper={t('api.nameHelp')}>
				<Input bind:value={newName} placeholder={t('api.namePlaceholder')} required disabled={creating} />
			</Field>

			{#if createError}
				<Alert tone="danger">{createError}</Alert>
			{/if}

			<div class="flex justify-end gap-2 pt-1">
				<Button type="button" variant="ghost" disabled={creating} onclick={() => (createOpen = false)}>
					{t('common.cancel')}
				</Button>
				<Button type="submit" disabled={creating || newName.trim().length === 0}>
					{#if creating}<Spinner size="sm" />{/if}
					{creating ? t('api.creating') : t('api.create')}
				</Button>
			</div>
		</form>
	{/if}
</Drawer>

<ConfirmModal
	bind:open={revokeOpen}
	tone="danger"
	title={t('api.revokeTitle')}
	description={t('api.revokeBody')}
	confirmLabel={t('api.revokeConfirm')}
	cancelLabel={t('accounts.keepIt')}
	onConfirm={revoke}
/>
