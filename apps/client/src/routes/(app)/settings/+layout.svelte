<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { PageTransition, SettingsLayout, icons, orderSections, settingsRoute } from '@facile/muse';
	import { isDesktop } from '$lib/platform';
	import { t } from '$lib/i18n/index.svelte';

	let { children: panel } = $props();

	const sections = $derived(
		orderSections([
			{ id: 'profile', label: t('settings.profile'), icon: icons.userCircle },
			{ id: 'preferences', label: t('settings.preferences'), icon: icons.palette },
			{ id: 'notifications', label: t('settings.notifications'), icon: icons.notification },
			{ id: 'accounts', label: t('settings.accounts'), icon: icons.mail },
			{ id: 'templates', label: t('settings.templates'), icon: icons.edit },
			{ id: 'tokens', label: t('settings.tokens'), icon: icons.key },
			{ id: 'subscription', label: t('settings.subscription'), icon: icons.card },
			{ id: 'advanced', label: t('settings.advanced'), icon: icons.shield },
			...(isDesktop() ? [{ id: 'updates', label: t('settings.updates'), icon: icons.refresh }] : [])
		]).map((section) => ({ ...section, href: `/settings/${section.id}` }))
	);

	const requested = $derived(page.url.pathname.split('/')[2] || undefined);
	const route = $derived(settingsRoute(requested, sections));

	$effect(() => {
		if (!route.known) void goto('/settings', { replaceState: true });
	});
</script>

<svelte:head>
	<title>{t('settings.title')}</title>
</svelte:head>

<div class="mx-auto w-full max-w-5xl p-6">
	<SettingsLayout
		{sections}
		active={requested}
		backHref="/settings"
		ordered={false}
		title={t('settings.heading')}
		description={t('settings.intro')}
		navLabel={t('settings.sections')}
	>
		{#snippet children(current)}
			<PageTransition key={current} distance={8} duration={0.25}>
				{@render panel()}
			</PageTransition>
		{/snippet}
	</SettingsLayout>
</div>
