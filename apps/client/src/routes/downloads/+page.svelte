<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Button, Card, Divider, icons } from '@facile/muse';
	import AppMark from '$lib/components/native/AppMark.svelte';
	import { APK_SHA256_URL, APK_URL, RELEASES_URL } from '$lib/downloads';
	import { isNative } from '$lib/platform';

	const STEPS = [
		{ title: 'Download the APK', body: 'Tap the button above on your phone. The file is about 4 MB.' },
		{
			title: 'Allow this source',
			body: 'Android asks once. Allow your browser to install unknown apps, then come back.'
		},
		{ title: 'Open and sign in', body: 'Open the downloaded file, install it, and sign in with your Facile account.' }
	];

	onMount(() => {
		if (isNative()) void goto('/mail', { replaceState: true });
	});
</script>

<svelte:head>
	<title>Get the app — Courrier</title>
	<meta name="description" content="Download Courrier for Android, or use it in any browser." />
</svelte:head>

<div class="min-h-dvh bg-fc-page text-fc-fg">
	<header class="border-b border-fc-border">
		<div class="mx-auto flex max-w-fc-lg items-center justify-between px-4 py-3 sm:px-6">
			<a href="/" class="flex min-h-11 items-center gap-2" aria-label="Courrier home">
				<AppMark size={24} />
				<span class="font-fc-title text-fc-lg font-semibold tracking-tight">Courrier</span>
			</a>
			<Button href="/login" variant="ghost" class="min-h-11">Log in</Button>
		</div>
	</header>

	<main class="mx-auto flex max-w-fc-md flex-col gap-10 px-4 py-12 sm:px-6 sm:py-16">
		<section class="flex flex-col items-center gap-4 text-center">
			<h1 class="font-fc-title text-fc-3xl font-semibold tracking-tight sm:text-5xl">Courrier on your phone</h1>
			<p class="max-w-md text-fc-md text-fc-fg-muted">
				The Android app reads the same mailboxes as the web, and tells you when a person writes to you.
			</p>
		</section>

		<Card class="flex flex-col gap-6">
			<div class="flex items-start gap-3">
				<iconify-icon icon="solar:smartphone-linear" width="24" height="24" aria-hidden="true"></iconify-icon>
				<div class="flex flex-col gap-1">
					<h2 class="font-fc-title text-fc-lg font-semibold">Android</h2>
					<p class="text-fc-sm text-fc-fg-muted">Signed APK, Android 7 or later.</p>
				</div>
			</div>
			<Button href={APK_URL} size="lg" icon={icons.download} class="h-12 w-full rounded-fc-md sm:w-auto sm:self-start">
				Download the APK
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
				To check the file, compare <code class="font-mono">sha256sum courrier-android.apk</code> with the
				<a href={APK_SHA256_URL} class="underline hover:text-fc-fg">published checksum</a>. Older versions are on the
				<a href={RELEASES_URL} class="underline hover:text-fc-fg">releases page</a>.
			</p>
		</Card>

		<Card class="flex items-start gap-3">
			<iconify-icon icon={icons.globe} width="24" height="24" aria-hidden="true"></iconify-icon>
			<div class="flex flex-col gap-1">
				<h2 class="font-fc-title text-fc-lg font-semibold">Any browser</h2>
				<p class="text-fc-sm text-fc-fg-muted">
					No install needed: Courrier works in any modern browser, on any computer or phone.
				</p>
				<div class="mt-2"><Button href="/login" variant="outline" class="min-h-11">Open the web app</Button></div>
			</div>
		</Card>
	</main>
</div>
