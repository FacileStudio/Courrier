<script lang="ts">
	import { onMount } from 'svelte';
	import { Alert, SecretField, SettingsSection, TokensSection, type TokenDraft } from '@facile/muse';
	import { backend, type ApiTokenStatus } from '$lib/backend';
	import { t } from '$lib/i18n/index.svelte';

	let status = $state<ApiTokenStatus>({ has_token: false });
	let loadError = $state('');

	const endpoint = `${backend.baseUrl || (typeof location === 'undefined' ? '' : location.origin)}/api`;

	const tokens = $derived(
		status.has_token
			? [
					{
						id: 'cli',
						name: status.name || 'CLI',
						prefix: t('api.storedHashed')
					}
				]
			: []
	);

	onMount(load);

	async function load() {
		try {
			status = await backend.apiTokenStatus();
			loadError = '';
		} catch (err) {
			loadError = err instanceof Error ? err.message : t('api.readFailed');
		}
	}

	async function create(draft: Omit<TokenDraft, 'days'>) {
		const created = await backend.createApiToken(draft.name);
		await load();
		return created.token;
	}

	async function revoke() {
		await backend.deleteApiToken();
		await load();
	}
</script>

<div class="flex flex-col gap-10">
	<SettingsSection title={t('api.endpoint')} description={t('api.endpointBody')}>
		<SecretField value={endpoint} sensitive={false} label={t('api.baseUrl')} />
	</SettingsSection>

	{#if loadError}
		<Alert tone="danger">{loadError}</Alert>
	{/if}

	<TokensSection
		{tokens}
		expiry={false}
		lastUsed={false}
		title={t('api.tokenTitle')}
		description={t('api.tokenBody')}
		onCreate={create}
		onRevoke={revoke}
	/>
</div>
