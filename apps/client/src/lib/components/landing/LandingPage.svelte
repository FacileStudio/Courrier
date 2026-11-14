<script lang="ts">
	import {
		DownloadSection,
		ForDevelopers,
		LANDING_HEADING,
		LandingCta,
		LandingFooter,
		LandingHeader,
		PricingPlans,
		SignInWithFacile,
		WorksWith,
		type DeveloperRow,
		type WorksWithApp
	} from '@facile/muse';
	import AccountsTile from './AccountsTile.svelte';
	import InboxDemo from './InboxDemo.svelte';
	import TabsTile from './TabsTile.svelte';
	import TeamTile from './TeamTile.svelte';
	import { CLI_COMMANDS, DESKTOP_CHECKSUMS_URL, DOWNLOADS, RELEASES_URL } from '$lib/downloads';
	import { rememberLocale, t } from '$lib/i18n/index.svelte';
	import { registerLandingIcons } from '$lib/landing-icons';

	let { signInHref }: { signInHref: string } = $props();

	registerLandingIcons();

	const ICON = 'solar:letter-opened-bold-duotone';
	const WIDTH = 'mx-auto w-full max-w-5xl px-4 sm:px-6';

	const links = $derived([
		{ href: '#features', label: t('landing.navFeatures') },
		{ href: '#download', label: t('landing.navApps') },
		{ href: '#pricing', label: t('landing.navPricing') }
	]);
	const columns = $derived([
		{
			title: 'Courrier',
			links: [
				{ href: '#features', label: t('landing.navFeatures') },
				{ href: '/downloads', label: t('landing.navApps') },
				{ href: '#pricing', label: t('landing.navPricing') }
			]
		},
		{ title: 'Facile', links: [{ href: 'https://facile.studio', label: 'facile.studio', external: true }] }
	]);
	const apps: WorksWithApp[] = $derived([
		{ name: 'Bulle', icon: 'solar:smile-circle-bold-duotone', line: t('landing.withBulle'), href: 'https://bulle.facile.studio' },
		{ name: 'Glouton', icon: 'solar:magnet-wave-bold-duotone', line: t('landing.withGlouton'), href: 'https://glouton.facile.studio' }
	]);
	const developer: DeveloperRow[] = $derived([
		{ kind: 'cli', text: t('landing.devCli'), command: CLI_COMMANDS[0].command },
		{ kind: 'mcp', text: t('landing.devMcp') },
		{ kind: 'source', text: t('landing.devSource'), href: 'https://github.com/FacileStudio/Courrier' }
	]);
</script>

<svelte:head>
	<title>{t('landing.title')}</title>
	<meta name="description" content={t('landing.metaDescription')} />
</svelte:head>

<div class="min-h-dvh overflow-x-clip bg-fc-page text-fc-fg">
	<LandingHeader appName="Courrier" icon={ICON} {links} skipHref="#main">
		{#snippet actions()}
			<SignInWithFacile href={signInHref} size="sm" data-sveltekit-reload />
		{/snippet}
	</LandingHeader>

	<main id="main">
		<section class="mx-auto flex max-w-6xl flex-col items-center px-4 pt-16 pb-12 text-center sm:px-6 sm:pt-24">
			<h1 class="max-w-4xl font-fc-title text-fc-3xl leading-tight font-bold tracking-tight sm:text-5xl">
				{t('landing.heroTitle')}
			</h1>
			<p class="mt-5 max-w-xl text-fc-lg text-fc-fg-muted">{t('landing.heroBody')}</p>
			<LandingCta class="mt-8 w-full sm:w-auto" {signInHref} align="center" data-sveltekit-reload />
			<p class="mt-4 text-fc-sm text-fc-fg-muted">{t('landing.heroNote')}</p>
			<div class="mt-12 w-full sm:mt-16"><InboxDemo /></div>
		</section>

		<section id="features" class="{WIDTH} scroll-mt-20 py-16 sm:py-20">
			<h2 class="text-center {LANDING_HEADING}">{t('landing.featuresTitle')}</h2>
			<div class="mt-10 grid gap-4 md:grid-cols-3">
				<AccountsTile />
				<TabsTile />
				<TeamTile />
			</div>
		</section>

		<WorksWith class="{WIDTH} py-16 sm:py-20" {apps} description={t('landing.withLine')} />

		<section id="download" class="mx-auto w-full max-w-4xl scroll-mt-20 px-4 py-16 sm:px-6 sm:py-20">
			<div class="mb-8 text-center">
				<h2 class={LANDING_HEADING}>{t('landing.downloadTitle')}</h2>
				<p class="mt-3 text-fc-fg-muted">{t('landing.downloadBody')}</p>
			</div>
			<DownloadSection
				downloads={[...DOWNLOADS]}
				releasesHref={RELEASES_URL}
				checksumsHref={DESKTOP_CHECKSUMS_URL}
				headingLevel={3}
				mode="device"
				allHref="/downloads"
			/>
		</section>

		<PricingPlans
			id="pricing"
			class="{WIDTH} scroll-mt-20 py-16 sm:py-20"
			title={t('landing.pricing')}
			headingClass={LANDING_HEADING}
			highlight={['courrier.mailboxes', 'ai.access']}
			highlightOnly
			notes={{ pro: t('landing.planAi'), studio: t('landing.planAi') }}
			recommended="pro"
			suiteHref="https://suite.facile.studio/#tarifs"
		/>

		<div class="{WIDTH} py-16 sm:py-20">
			<ForDevelopers rows={developer} description={t('landing.devLine')} />
		</div>

		<section class="{WIDTH} py-16 text-center sm:py-24">
			<h2 class={LANDING_HEADING}>{t('landing.ctaTitle')}</h2>
			<p class="mt-3 text-fc-fg-muted">{t('landing.ctaBody')}</p>
			<LandingCta class="mt-8" {signInHref} align="center" data-sveltekit-reload />
		</section>

	</main>

	<LandingFooter appName="Courrier" icon={ICON} tagline={t('landing.tagline')} {columns} legal={{ onLocaleChange: rememberLocale }} />
</div>
