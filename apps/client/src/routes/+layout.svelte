<script lang="ts">
	import '../app.css';
	import { browser } from '$app/environment';
	import { Toaster } from '@facile/muse';
	import { registerDeepLinks } from '$lib/native';
	import { isApp, isDesktop } from '$lib/platform';
	import { theme } from '$lib/theme.svelte';
	import { offerUpdate } from '$lib/updater';

	let { children } = $props();

	if (browser) {
		theme.restore();
		void import('iconify-icon');
		if (isApp()) registerDeepLinks();
		if (isDesktop()) void offerUpdate();
	}
</script>

{@render children()}

<!-- One Toaster for the whole app, outside every page so a route change cannot unmount a
     toast mid-flight. The bottom padding clears MobileNav on a phone. -->
<Toaster class="pb-28 md:pb-6" />
