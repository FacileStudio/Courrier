<script lang="ts">
	import '../app.css';
	import { browser } from '$app/environment';
	import { Toaster, provideLocale } from '@facile/muse';
	import TitleBar from '$lib/components/native/TitleBar.svelte';
	import { registerDeepLinks } from '$lib/native';
	import { isApp, isDesktop } from '$lib/platform';
	import { theme } from '$lib/theme.svelte';
	import { updater } from '$lib/updater.svelte';

	let { children } = $props();

	provideLocale(() => 'en');

	if (browser) {
		theme.restore();
		void import('iconify-icon');
		if (isApp()) registerDeepLinks();
		if (isDesktop()) {
			document.documentElement.dataset.desktop = '';
			void updater.checkNow(true);
		}
	}
</script>

{#if browser && isDesktop()}
	<TitleBar />
{/if}

{@render children()}

<!-- One Toaster for the whole app, outside every page so a route change cannot unmount a
     toast mid-flight. The bottom padding clears MobileNav on a phone. -->
<Toaster class="pb-28 md:pb-6" />
