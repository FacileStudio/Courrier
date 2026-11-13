<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import {
		DownloadCta,
		DownloadSection,
		Icon,
		LandingFooter,
		LandingHeader,
		PricingPlans,
		SignInWithFacile,
		Spinner
	} from '@facile/muse';
	import { backend } from '$lib/backend';
	import InboxDemo from '$lib/components/landing/InboxDemo.svelte';
	import NativeIntro from '$lib/components/native/NativeIntro.svelte';
	import { CLI_COMMANDS, DESKTOP_CHECKSUMS_URL, DOWNLOADS, RELEASES_URL } from '$lib/downloads';
	import { rememberLocale, t } from '$lib/i18n/index.svelte';
	import { isDesktop, isNative } from '$lib/platform';

	const INTRO_SEEN_KEY = 'courrier.intro-seen';
	const native = isNative();
	const desktop = isDesktop();
	const signInHref = `${backend.baseUrl}/api/auth/oidc`;

	let ready = $state(false);

	function introSeen(): boolean {
		try {
			return localStorage.getItem(INTRO_SEEN_KEY) === '1';
		} catch {
			return false;
		}
	}

	function markIntroSeen() {
		try {
			localStorage.setItem(INTRO_SEEN_KEY, '1');
		} catch {
			return;
		}
	}

	async function finishIntro() {
		markIntroSeen();
		await goto('/login', { replaceState: true });
	}

	onMount(async () => {
		try {
			await backend.me();
			goto('/mail', { replaceState: true });
			return;
		} catch {}

		if (desktop) {
			await goto('/login', { replaceState: true });
			return;
		}

		if (native) {
			if (introSeen()) await goto('/login', { replaceState: true });
			else ready = true;
			return;
		}

		ready = true;
	});

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
	const features = $derived([
		{ icon: 'solar:inbox-in-linear', title: t('landing.imapTitle'), body: t('landing.imapBody') },
		{ icon: 'solar:filter-linear', title: t('landing.tabsTitle'), body: t('landing.tabsBody') },
		{ icon: 'solar:users-group-rounded-linear', title: t('landing.spacesTitle'), body: t('landing.spacesBody') },
		{ icon: 'solar:server-linear', title: t('landing.selfHostedTitle'), body: t('landing.selfHostedBody') }
	]);
</script>

{#if !ready}
	{#if native || desktop}
		<div class="flex min-h-dvh items-center justify-center bg-fc-page"><Spinner /></div>
	{/if}
{:else if native}
	<NativeIntro onfinish={finishIntro} />
{:else}
	<div class="min-h-dvh bg-fc-page text-fc-fg">
		<LandingHeader appName="Courrier" icon="solar:letter-opened-bold-duotone" {links} skipHref="#main">
			{#snippet actions()}
				<SignInWithFacile href={signInHref} data-sveltekit-reload />
			{/snippet}
		</LandingHeader>

		<main id="main">
			<section class="mx-auto flex max-w-6xl flex-col items-center px-4 pt-16 pb-12 text-center sm:px-6 sm:pt-24">
				<h1 class="font-fc-title text-fc-3xl font-bold tracking-tight leading-tight sm:text-5xl">
					{t('landing.heroTitle')}
				</h1>
				<p class="mt-5 max-w-lg text-fc-lg text-fc-fg-muted">{t('landing.heroBody')}</p>
				<div class="mt-8 flex flex-col items-center gap-3">
					<SignInWithFacile href={signInHref} size="lg" data-sveltekit-reload />
					<DownloadCta
						href="/downloads"
						platforms={['macos', 'windows', 'linux', 'android']}
						label={t('landing.downloadCta')}
					/>
				</div>
				<div class="mt-12 w-full sm:mt-16"><InboxDemo /></div>
			</section>

			<section id="features" class="mx-auto max-w-6xl scroll-mt-20 px-4 py-16 sm:px-6 sm:py-20">
				<h2 class="text-center font-fc-title text-fc-2xl font-bold tracking-tight sm:text-fc-3xl">
					{t('landing.featuresTitle')}
				</h2>
				<div class="mt-10 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
					{#each features as feature (feature.icon)}
						<div class="rounded-fc-lg bg-fc-component p-6">
							<span class="flex size-10 items-center justify-center rounded-fc-md bg-fc-surface">
								<Icon icon={feature.icon} size={20} />
							</span>
							<h3 class="mt-4 text-fc-md font-bold">{feature.title}</h3>
							<p class="mt-1.5 text-fc-sm text-fc-fg-muted">{feature.body}</p>
						</div>
					{/each}
				</div>
			</section>

			<section id="download" class="mx-auto max-w-4xl scroll-mt-20 px-4 py-16 sm:px-6 sm:py-20">
				<div class="mb-8 text-center">
					<h2 class="font-fc-title text-fc-2xl font-bold tracking-tight sm:text-fc-3xl">{t('landing.downloadTitle')}</h2>
					<p class="mt-3 text-fc-fg-muted">{t('landing.downloadBody')}</p>
				</div>
				<DownloadSection
					downloads={[...DOWNLOADS]}
					commands={[...CLI_COMMANDS]}
					releasesHref={RELEASES_URL}
					checksumsHref={DESKTOP_CHECKSUMS_URL}
					headingLevel={3}
				/>
			</section>

			<PricingPlans
				id="pricing"
				highlight={['courrier.mailboxes']}
				title={t('landing.pricing')}
				class="mx-auto max-w-5xl scroll-mt-20 px-4 py-16 sm:px-6 sm:py-20"
			/>

			<section class="mx-auto max-w-5xl px-4 py-16 text-center sm:px-6 sm:py-24">
				<h2 class="font-fc-title text-fc-2xl font-bold tracking-tight sm:text-fc-3xl">{t('landing.ctaTitle')}</h2>
				<p class="mt-3 text-fc-fg-muted">{t('landing.ctaBody')}</p>
				<div class="mt-8 flex justify-center">
					<SignInWithFacile href={signInHref} size="lg" data-sveltekit-reload />
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
{/if}
