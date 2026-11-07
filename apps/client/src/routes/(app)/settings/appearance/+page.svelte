<script lang="ts">
	import { OptionCards, SettingsRow, SettingsSection, StatusDot, icons } from '@facile/muse';
	import { theme, type ThemePreference } from '$lib/theme.svelte';
	import { t } from '$lib/i18n/index.svelte';

	const modes = $derived([
		{ value: 'system', label: t('appearance.system'), icon: icons.monitor },
		{ value: 'light', label: t('appearance.light'), icon: icons.sun },
		{ value: 'dark', label: t('appearance.dark'), icon: icons.moon }
	]);

	let mode = $state<string>(theme.preference);

	$effect(() => {
		theme.set(mode as ThemePreference);
	});
</script>

<div class="flex flex-col gap-10">
	<SettingsSection
		title={t('appearance.theme')}
		description={t('appearance.themeBody')}
	>
		<SettingsRow
			label={t('appearance.scheme')}
			description={t('appearance.schemeBody')}
			stacked
		>
			<OptionCards options={modes} bind:value={mode} name="theme-mode" label={t('appearance.scheme')} />
		</SettingsRow>

		<SettingsRow
			label={t('appearance.applied')}
			description={t('appearance.appliedBody')}
		>
			<StatusDot tone="neutral" label={theme.resolved === 'dark' ? t('appearance.dark') : t('appearance.light')} />
		</SettingsRow>
	</SettingsSection>
</div>
