<script lang="ts">
	import { OptionCards, SettingsRow, SettingsSection, StatusDot, icons } from '@facile/muse';
	import { theme, type ThemePreference } from '$lib/theme.svelte';

	const modes = [
		{ value: 'system', label: 'System', icon: icons.monitor },
		{ value: 'light', label: 'Light', icon: icons.sun },
		{ value: 'dark', label: 'Dark', icon: icons.moon }
	];

	let mode = $state<string>(theme.preference);

	$effect(() => {
		theme.set(mode as ThemePreference);
	});
</script>

<div class="flex flex-col gap-10">
	<SettingsSection
		title="Theme"
		description="Stored in this browser. Every device you sign in from keeps its own choice."
	>
		<SettingsRow
			label="Colour scheme"
			description="System follows your operating system, and keeps following it while the app is open."
			stacked
		>
			<OptionCards options={modes} bind:value={mode} name="theme-mode" label="Colour scheme" />
		</SettingsRow>

		<SettingsRow
			label="Currently applied"
			description="What System resolved to, or the scheme you picked."
		>
			<StatusDot tone="neutral" label={theme.resolved === 'dark' ? 'Dark' : 'Light'} />
		</SettingsRow>
	</SettingsSection>
</div>
