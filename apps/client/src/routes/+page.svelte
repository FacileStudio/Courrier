<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import {
		DownloadCta,
		LandingFooter,
		LandingHeader,
		PricingPlans,
		SignInWithFacile,
		Spinner,
		icons
	} from '@facile/muse';
	import { backend } from '$lib/backend';
	import NativeIntro from '$lib/components/native/NativeIntro.svelte';
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
		{ href: '#pricing', label: t('landing.navPricing') },
		{ href: '/downloads', label: t('landing.navApps') }
	]);
	const columns = $derived([
		{
			title: 'Courrier',
			links: [
				{ href: '#features', label: t('landing.navFeatures') },
				{ href: '#pricing', label: t('landing.navPricing') },
				{ href: '/downloads', label: t('landing.navApps') }
			]
		},
		{ title: 'Facile', links: [{ href: 'https://facile.studio', label: 'facile.studio', external: true }] }
	]);
	const features = $derived([
		{ icon: 'solar:inbox-in-linear', title: t('landing.imapTitle'), body: t('landing.imapBody') },
		{ icon: 'solar:users-group-rounded-linear', title: t('landing.spacesTitle'), body: t('landing.spacesBody') },
		{ icon: 'solar:lock-linear', title: t('landing.selfHostedTitle'), body: t('landing.selfHostedBody') }
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
		<LandingHeader appName="Courrier" icon="solar:letter-bold-duotone" {links} skipHref="#main">
			{#snippet actions()}
				<SignInWithFacile href={signInHref} data-sveltekit-reload />
			{/snippet}
		</LandingHeader>

		<main id="main">
			<section class="mx-auto max-w-5xl px-6 py-24 text-center">
				<h1 class="font-fc-title text-fc-3xl font-semibold sm:text-5xl tracking-tight leading-tight">
					{t('landing.heroTitle')}
				</h1>
				<p class="mx-auto mt-6 max-w-xl text-fc-lg text-fc-fg-muted leading-relaxed">
					{t('landing.heroBody')}
				</p>
				<div class="mt-10 flex justify-center">
					<SignInWithFacile href={signInHref} size="lg" data-sveltekit-reload />
				</div>
				<div class="mt-4 flex justify-center">
					<DownloadCta
						href="/downloads"
						platforms={['macos', 'windows', 'linux', 'android']}
						label={t('landing.downloadCta')}
					/>
				</div>
			</section>

			<section id="features" class="mx-auto max-w-5xl px-6 py-20">
				<div class="grid gap-6 md:grid-cols-3">
					{#each features as feature (feature.icon)}
						<div class="rounded-fc-lg p-6">
							<div class="mb-3 flex size-10 items-center justify-center rounded-fc-md">
								<iconify-icon icon={feature.icon} width="20" height="20" class="size-5"></iconify-icon>
							</div>
							<h3 class="text-fc-md font-semibold">{feature.title}</h3>
							<p class="mt-1.5 text-fc-sm text-fc-fg-muted">{feature.body}</p>
						</div>
					{/each}
				</div>
			</section>

			<PricingPlans
				id="pricing"
				highlight={['courrier.mailboxes']}
				title={t('landing.pricing')}
				class="mx-auto max-w-5xl px-6 py-20"
			/>

			<section class="mx-auto max-w-5xl px-6 py-20 text-center">
				<h2 class="font-fc-title text-fc-3xl font-semibold tracking-tight">{t('landing.ctaTitle')}</h2>
				<p class="mt-4 text-fc-fg-muted">{t('landing.ctaBody')}</p>
				<div class="mt-8 flex justify-center">
					<SignInWithFacile href={signInHref} size="lg" data-sveltekit-reload />
				</div>
			</section>
		</main>

		<LandingFooter
			appName="Courrier"
			icon="solar:letter-bold-duotone"
			tagline={t('landing.tagline')}
			{columns}
			legal={{ onLocaleChange: rememberLocale }}
		/>
	</div>
{/if}
