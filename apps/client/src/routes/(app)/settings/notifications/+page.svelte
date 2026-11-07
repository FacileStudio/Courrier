<script lang="ts">
	import { SettingsRow, SettingsSection } from '@facile/muse';
	import NotificationControl from '$lib/components/NotificationControl.svelte';
	import { isDesktop, isNative } from '$lib/platform';
	import { t } from '$lib/i18n/index.svelte';

	const device = $derived(
		isNative()
			? t('notifySettings.devicePhone')
			: isDesktop()
				? t('notifySettings.deviceComputer')
				: t('notifySettings.deviceBrowser')
	);
</script>

<div class="flex flex-col gap-10">
	<SettingsSection
		title={t('notifySettings.title')}
		description={t('notifySettings.body')}
	>
		<SettingsRow
			label={t('notifySettings.newMail')}
			description={t('notifySettings.newMailBody', { device })}
			stacked
		>
			<NotificationControl showTest />
		</SettingsRow>
		<SettingsRow
			label={t('notifySettings.shows')}
			description={t('notifySettings.showsBody')}
		/>
	</SettingsSection>
</div>
