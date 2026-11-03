<script lang="ts">
	import { Badge, Button, SettingsRow, SettingsSection, icons } from '@facile/muse';
	import { updater } from '$lib/updater.svelte';
</script>

<div class="flex flex-col gap-10">
	<SettingsSection
		title="Desktop updates"
		description="Keep Courrier desktop up to date with new features and fixes."
	>
		<SettingsRow
			label="Current version"
			description="The release build installed on this computer."
		>
			<span class="font-mono text-fc-sm text-fc-fg">v{updater.currentVersion}</span>
		</SettingsRow>

		<SettingsRow
			label="Update status"
			description="The latest state from the Courrier release channel."
		>
			{#if updater.state === 'checking'}
				<Badge tone="info">Checking for updates…</Badge>
			{:else if updater.state === 'available'}
				<Badge tone="accent">Version v{updater.availableVersion} is available</Badge>
			{:else if updater.state === 'downloading'}
				<Badge tone="info">Downloading {updater.downloadProgress}%</Badge>
			{:else if updater.state === 'ready'}
				<Badge tone="success">Ready to install</Badge>
			{:else if updater.state === 'error'}
				<Badge tone="danger">{updater.errorMessage || 'Update failed'}</Badge>
			{:else if updater.state === 'upToDate'}
				<Badge tone="success">Courrier is up to date</Badge>
			{:else}
				<Badge tone="neutral">Ready to check</Badge>
			{/if}
		</SettingsRow>

		{#if updater.state === 'downloading'}
			<SettingsRow
				label="Download progress"
				description="Fetching update package from GitHub releases."
				stacked
			>
				<div class="flex w-full flex-col gap-2">
					<div class="h-2 w-full overflow-hidden rounded-fc-pill bg-fc-surface">
						<div
							class="h-full bg-fc-accent transition-all duration-200"
							style="width: {updater.downloadProgress}%"
						></div>
					</div>
					<span class="font-mono text-fc-xs text-fc-fg-muted">{updater.downloadProgress}%</span>
				</div>
			</SettingsRow>
		{/if}

		<SettingsRow
			label="Actions"
			description="Manually check or install available updates."
		>
			{#if updater.state === 'available'}
				<Button
					variant="primary"
					icon={icons.download}
					onclick={() => void updater.installNow()}
				>
					Download & Install
				</Button>
			{:else if updater.state === 'ready'}
				<Button
					variant="primary"
					icon={icons.refresh}
					onclick={() => void updater.restartNow()}
				>
					Restart & Install
				</Button>
			{:else}
				<Button
					variant="outline"
					icon={icons.refresh}
					disabled={updater.state === 'checking' || updater.state === 'downloading'}
					onclick={() => void updater.checkNow(false)}
				>
					{updater.state === 'checking' ? 'Checking…' : 'Check for updates'}
				</Button>
			{/if}
		</SettingsRow>

		<SettingsRow
			label="Release history"
			description="View changelog and previous release notes on GitHub."
		>
			<Button
				variant="ghost"
				size="sm"
				icon={icons.history}
				href="https://github.com/FacileStudio/Courrier/releases"
				target="_blank"
				rel="noopener noreferrer"
			>
				View releases
			</Button>
		</SettingsRow>
	</SettingsSection>
</div>
