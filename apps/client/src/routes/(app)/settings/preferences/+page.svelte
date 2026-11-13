<script lang="ts">
	import { PreferencesSection } from '@facile/muse';
	import { theme } from '$lib/theme.svelte';
	import { rememberLocale } from '$lib/i18n/index.svelte';
</script>

<PreferencesSection
	theme={theme.mode}
	onThemeChange={(next) => theme.set(next)}
	onLocaleChange={rememberLocale}
/>
