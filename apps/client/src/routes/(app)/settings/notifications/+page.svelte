<script lang="ts">
	import { SettingsRow, SettingsSection } from '@facile/muse';
	import NotificationControl from '$lib/components/NotificationControl.svelte';
	import { isNative } from '$lib/platform';

	const device = isNative() ? 'this phone' : 'this browser';
</script>

<div class="flex flex-col gap-10">
	<SettingsSection
		title="Notifications"
		description="Courrier checks your inboxes every few minutes and tells you when a person writes to you. Newsletters, lists and robots never notify."
	>
		<SettingsRow
			label="New mail"
			description="Stored per device. Turn it on for {device}; your other devices keep their own choice."
			stacked
		>
			<NotificationControl showTest />
		</SettingsRow>
		<SettingsRow
			label="What a notification shows"
			description="The sender's name and the address it reached. Never the subject or the message."
		/>
	</SettingsSection>
</div>
