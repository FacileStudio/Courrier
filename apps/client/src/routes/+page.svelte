<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { Spinner } from '@facile/muse';
	import { backend } from '$lib/backend';
	import LandingPage from '$lib/components/landing/LandingPage.svelte';
	import NativeIntro from '$lib/components/native/NativeIntro.svelte';
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
</script>

{#if !ready}
	{#if native || desktop}
		<div class="flex min-h-dvh items-center justify-center bg-fc-page"><Spinner /></div>
	{/if}
{:else if native}
	<NativeIntro onfinish={finishIntro} />
{:else}
	<LandingPage {signInHref} />
{/if}
