<script lang="ts">
	import '../app.css';
	import { browser } from '$app/environment';
	import { Toaster } from '@facile/muse';
	import { theme } from '$lib/theme.svelte';

	let { children } = $props();

	if (browser) {
		theme.restore();
		void import('iconify-icon');
	}
</script>

{@render children()}

<!-- One Toaster for the whole app, outside every page so a route change cannot unmount a
     toast mid-flight. The bottom padding clears MobileNav on a phone. -->
<Toaster class="pb-28 md:pb-6" />
