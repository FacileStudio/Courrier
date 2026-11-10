<script lang="ts">
	import '../app.css';
	import { browser } from '$app/environment';
	import { Toaster, provideLocale } from '@facile/muse';
	import { bindLocale, localeSource } from '$lib/i18n/index.svelte';
	import TitleBar from '$lib/components/native/TitleBar.svelte';
	import { registerDeepLinks } from '$lib/native';
	import { forgetLegacyWebToken, isApp, isDesktop } from '$lib/platform';
	import { theme } from '$lib/theme.svelte';
	import { updater } from '$lib/updater.svelte';

	let { children } = $props();

	const i18n = provideLocale(localeSource);
	bindLocale(() => i18n.locale);

	if (browser) {
		theme.restore();
		forgetLegacyWebToken(localStorage);
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
