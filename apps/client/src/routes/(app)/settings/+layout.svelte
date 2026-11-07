<script lang="ts">
	import { page } from '$app/state';
	import { Divider, PageTransition, Tabs, icons } from '@facile/muse';
	import { isDesktop } from '$lib/platform';
	import { t } from '$lib/i18n/index.svelte';

	let { children } = $props();

	const sections = $derived([
		{ id: 'profile', label: t('settings.profile'), icon: icons.userCircle },
		{ id: 'appearance', label: t('settings.appearance'), icon: icons.palette },
		{ id: 'accounts', label: t('settings.accounts'), icon: icons.mail },
		{ id: 'notifications', label: t('settings.notifications'), icon: icons.notification },
		{ id: 'templates', label: t('settings.templates'), icon: icons.edit },
		{ id: 'api', label: t('settings.api'), icon: icons.key },
		{ id: 'advanced', label: t('settings.advanced'), icon: icons.shield },
		...(isDesktop() ? [{ id: 'updates', label: t('settings.updates'), icon: icons.refresh }] : [])
	]);

	const items = $derived(sections.map((section) => ({ ...section, href: `/settings/${section.id}` })));
	const active = $derived(
		sections.find((section) => page.url.pathname.startsWith(`/settings/${section.id}`))?.id ??
			sections[0].id
	);
</script>

<svelte:head>
	<title>{t('settings.title')}</title>
</svelte:head>

<div class="mx-auto flex w-full max-w-4xl flex-col gap-8 p-6">
	<div class="flex flex-col gap-2">
		<h1 class="text-fc-2xl font-semibold text-fc-fg">{t('settings.heading')}</h1>
		<p class="text-fc-sm text-fc-fg-muted">
			{t('settings.intro')}
		</p>
	</div>

	<!-- gap-4, not tighter: the rule separates a page header from its body, so it needs air. -->
	<div class="flex flex-col gap-4">
		<Tabs {items} value={active} label={t('settings.sections')} />
		<Divider class="my-0" />
	</div>

	<PageTransition key={active} distance={8} duration={0.25}>
		{@render children()}
	</PageTransition>
</div>
