<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { Alert, AuthScreen, icons } from '@facile/muse';
	import { backend } from '$lib/backend';
	import NativeLogin from '$lib/components/native/NativeLogin.svelte';
	import { rememberLocale, t } from '$lib/i18n/index.svelte';
	import { isApp } from '$lib/platform';

	const native = isApp();
	let failed = $state(false);
	let nativeError = $state('');

	$effect(() => {
		nativeError = page.url.searchParams.get('error') ?? '';
	});

	onMount(async () => {
		failed = page.url.searchParams.has('error');
		try {
			await backend.me();
			goto('/mail');
		} catch {}
	});
</script>

<svelte:head>
	<title>{t('login.title')}</title>
</svelte:head>

{#if native}
	<NativeLogin bind:error={nativeError} />
{:else}
	<AuthScreen
		appName="Courrier"
		href="{backend.baseUrl}/api/auth/oidc"
		data-sveltekit-reload
		onLocaleChange={rememberLocale}
	>
		{#snippet logo()}
			<iconify-icon icon={icons.mail} width="40" height="40" class="block text-fc-fg"></iconify-icon>
		{/snippet}
		{#snippet tagline()}
			{t('login.tagline')}
		{/snippet}
		{#snippet notice()}
			{#if failed}
				<Alert tone="danger">{t('login.failed')}</Alert>
			{/if}
		{/snippet}
	</AuthScreen>
{/if}
