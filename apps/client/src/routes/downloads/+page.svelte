<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Button, Card, Divider, icons } from '@facile/muse';
	import AppMark from '$lib/components/native/AppMark.svelte';
	import {
		APK_SHA256_URL,
		APK_URL,
		DEB_URL,
		DESKTOP_BUILDS,
		DESKTOP_CHECKSUMS_URL,
		RELEASES_URL,
		detectOS,
		type DesktopOS
	} from '$lib/downloads';
	import { slots } from '$lib/i18n/format';
	import { t } from '$lib/i18n/index.svelte';
	import { isNative } from '$lib/platform';

	const STEPS = $derived([
		{ title: t('downloads.step1Title'), body: t('downloads.step1Body') },
		{ title: t('downloads.step2Title'), body: t('downloads.step2Body') },
		{ title: t('downloads.step3Title'), body: t('downloads.step3Body') }
	]);

	let os = $state<DesktopOS | null>(null);

	onMount(() => {
		if (isNative()) void goto('/mail', { replaceState: true });
		os = detectOS(navigator.userAgent);
	});
</script>

<svelte:head>
	<title>{t('downloads.title')}</title>
	<meta name="description" content={t('downloads.description')} />
</svelte:head>

<div class="min-h-dvh bg-fc-page text-fc-fg">
	<header class="">
		<div class="mx-auto flex max-w-fc-lg items-center justify-between px-4 py-3 sm:px-6">
			<a href="/" class="flex min-h-11 items-center gap-2" aria-label={t('downloads.home')}>
				<AppMark size={24} />
				<span class="font-fc-title text-fc-lg font-semibold tracking-tight">Courrier</span>
			</a>
			<Button href="/mail" variant="ghost" class="min-h-11">{t('downloads.open')}</Button>
		</div>
	</header>

	<main class="mx-auto flex max-w-fc-md flex-col gap-10 px-4 py-12 sm:px-6 sm:py-16">
		<section class="flex flex-col items-center gap-4 text-center">
			<h1 class="font-fc-title text-fc-3xl font-semibold tracking-tight sm:text-5xl">{t('downloads.heading')}</h1>
			<p class="max-w-md text-fc-md text-fc-fg-muted">
				{t('downloads.intro')}
			</p>
		</section>

		<Card class="flex flex-col gap-6">
			<div class="flex items-start gap-3">
				<iconify-icon icon="solar:monitor-linear" width="24" height="24" aria-hidden="true"></iconify-icon>
				<div class="flex flex-col gap-1">
					<h2 class="font-fc-title text-fc-lg font-semibold">{t('downloads.desktop')}</h2>
					<p class="text-fc-sm text-fc-fg-muted">{t('downloads.desktopBody')}</p>
				</div>
			</div>
			<ul class="flex flex-col gap-2">
				{#each DESKTOP_BUILDS as build (build.os)}
					<li class="flex flex-wrap items-center justify-between gap-3 rounded-fc-md bg-fc-surface px-3 py-2">
						<div class="flex flex-col">
							<span class="text-fc-sm font-medium">{build.label}</span>
							<span class="text-fc-xs text-fc-fg-muted">{t(build.detailKey)}</span>
						</div>
						<Button
							href={build.url}
							variant={os === build.os ? 'primary' : 'secondary'}
							icon={icons.download}
							class="min-h-11"
						>
							{t('downloads.download')}
						</Button>
					</li>
				{/each}
			</ul>
			<p class="text-fc-sm text-fc-fg-muted">
				{#each slots(t('downloads.debNote')) as part, i (i)}
					{#if typeof part === 'string'}
						{part}
					{:else if part.slot === 'link'}
						<a href={DEB_URL} class="underline hover:text-fc-fg">{t('downloads.debLink')}</a>
					{:else}
						<a href={DESKTOP_CHECKSUMS_URL} class="underline hover:text-fc-fg">{t('downloads.checksumsLink')}</a>
					{/if}
				{/each}
			</p>
		</Card>

		<Card class="flex flex-col gap-6">
			<div class="flex items-start gap-3">
				<iconify-icon icon="solar:smartphone-linear" width="24" height="24" aria-hidden="true"></iconify-icon>
				<div class="flex flex-col gap-1">
					<h2 class="font-fc-title text-fc-lg font-semibold">{t('downloads.android')}</h2>
					<p class="text-fc-sm text-fc-fg-muted">{t('downloads.androidBody')}</p>
				</div>
			</div>
			<Button href={APK_URL} size="lg" icon={icons.download} class="h-12 w-full rounded-fc-md sm:w-auto sm:self-start">
				{t('downloads.apk')}
			</Button>
			<ol class="flex flex-col gap-4">
				{#each STEPS as step, i (step.title)}
					<li class="flex gap-3">
						<span
							class="flex size-7 shrink-0 items-center justify-center rounded-fc-pill bg-fc-surface text-fc-sm font-semibold"
							aria-hidden="true">{i + 1}</span
						>
						<div class="flex flex-col gap-0.5">
							<span class="text-fc-sm font-medium">{step.title}</span>
							<span class="text-fc-sm text-fc-fg-muted">{step.body}</span>
						</div>
					</li>
				{/each}
			</ol>
			<Divider class="my-0" />
			<p class="text-fc-sm text-fc-fg-muted">
				{#each slots(t('downloads.apkNote')) as part, i (i)}
					{#if typeof part === 'string'}
						{part}
					{:else if part.slot === 'command'}
						<code class="font-mono">sha256sum courrier-android.apk</code>
					{:else if part.slot === 'checksum'}
						<a href={APK_SHA256_URL} class="underline hover:text-fc-fg">{t('downloads.checksumLink')}</a>
					{:else}
						<a href={RELEASES_URL} class="underline hover:text-fc-fg">{t('downloads.releasesLink')}</a>
					{/if}
				{/each}
			</p>
		</Card>

		<Card class="flex items-start gap-3">
			<iconify-icon icon={icons.globe} width="24" height="24" aria-hidden="true"></iconify-icon>
			<div class="flex flex-col gap-1">
				<h2 class="font-fc-title text-fc-lg font-semibold">{t('downloads.browser')}</h2>
				<p class="text-fc-sm text-fc-fg-muted">
					{t('downloads.browserBody')}
				</p>
				<div class="mt-2"><Button href="/login" variant="secondary" class="min-h-11">{t('downloads.openWeb')}</Button></div>
			</div>
		</Card>
	</main>
</div>
