<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Button, DownloadSection, Icon, LandingFooter, LandingHeader, icons } from '@facile/muse';
	import { CLI_COMMANDS, DESKTOP_CHECKSUMS_URL, DOWNLOADS, RELEASES_URL } from '$lib/downloads';
	import { rememberLocale, t } from '$lib/i18n/index.svelte';
	import { isNative } from '$lib/platform';

	const links = $derived([
		{ href: '/#features', label: t('landing.navFeatures') },
		{ href: '/downloads', label: t('landing.navApps') },
		{ href: '/#pricing', label: t('landing.navPricing') }
	]);
	const columns = $derived([
		{ title: 'Courrier', links },
		{ title: 'Facile', links: [{ href: 'https://facile.studio', label: 'facile.studio', external: true }] }
	]);

	onMount(() => {
		if (isNative()) void goto('/mail', { replaceState: true });
	});
</script>

<svelte:head>
	<title>{t('downloads.title')}</title>
	<meta name="description" content={t('downloads.description')} />
</svelte:head>

<div class="min-h-dvh bg-fc-page text-fc-fg">
	<LandingHeader appName="Courrier" icon="solar:letter-opened-bold-duotone" {links} skipHref="#main" />

	<main id="main" class="mx-auto flex max-w-4xl flex-col gap-10 px-4 py-12 sm:px-6 sm:py-16">
		<section class="flex flex-col items-center gap-4 text-center">
			<h1 class="font-fc-title text-fc-3xl font-bold tracking-tight sm:text-5xl">{t('downloads.heading')}</h1>
			<p class="max-w-md text-fc-md text-fc-fg-muted">{t('downloads.intro')}</p>
		</section>

		<DownloadSection
			downloads={[...DOWNLOADS]}
			commands={[...CLI_COMMANDS]}
			releasesHref={RELEASES_URL}
			checksumsHref={DESKTOP_CHECKSUMS_URL}
		/>

		<section class="flex items-start gap-3 rounded-fc-lg bg-fc-component p-5">
			<span class="flex size-10 shrink-0 items-center justify-center rounded-fc-md bg-fc-surface">
				<Icon icon={icons.globe} size={20} />
			</span>
			<div class="flex flex-col gap-1">
				<h2 class="text-fc-md font-bold">{t('downloads.browser')}</h2>
				<p class="text-fc-sm text-fc-fg-muted">{t('downloads.browserBody')}</p>
				<div class="mt-2"><Button href="/login" variant="secondary" class="min-h-11">{t('downloads.openWeb')}</Button></div>
			</div>
		</section>
	</main>

	<LandingFooter
		appName="Courrier"
		icon="solar:letter-opened-bold-duotone"
		tagline={t('landing.tagline')}
		{columns}
		legal={{ onLocaleChange: rememberLocale }}
	/>
</div>
