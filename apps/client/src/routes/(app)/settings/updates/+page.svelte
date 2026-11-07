<script lang="ts">
	import { Badge, Button, SettingsRow, SettingsSection, icons } from '@facile/muse';
	import { updater } from '$lib/updater.svelte';
	import { t } from '$lib/i18n/index.svelte';
</script>

<div class="flex flex-col gap-10">
	<SettingsSection
		title={t('updates.title')}
		description={t('updates.body')}
	>
		<SettingsRow
			label={t('updates.version')}
			description={t('updates.versionBody')}
		>
			<span class="font-mono text-fc-sm text-fc-fg">v{updater.currentVersion}</span>
		</SettingsRow>

		<SettingsRow
			label={t('updates.status')}
			description={t('updates.statusBody')}
		>
			{#if updater.state === 'checking'}
				<Badge tone="info">{t('updates.checking')}</Badge>
			{:else if updater.state === 'available'}
				<Badge tone="accent">{t('updates.available', { version: updater.availableVersion ?? '' })}</Badge>
			{:else if updater.state === 'downloading'}
				<Badge tone="info">{t('updates.downloading', { progress: updater.downloadProgress })}</Badge>
			{:else if updater.state === 'ready'}
				<Badge tone="success">{t('updates.ready')}</Badge>
			{:else if updater.state === 'error'}
				<Badge tone="danger">{updater.errorMessage || t('updates.failed')}</Badge>
			{:else if updater.state === 'upToDate'}
				<Badge tone="success">{t('updates.upToDate')}</Badge>
			{:else}
				<Badge tone="neutral">{t('updates.idle')}</Badge>
			{/if}
		</SettingsRow>

		{#if updater.state === 'downloading'}
			<SettingsRow
				label={t('updates.progress')}
				description={t('updates.progressBody')}
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
			label={t('updates.actions')}
			description={t('updates.actionsBody')}
		>
			{#if updater.state === 'available'}
				<Button
					variant="primary"
					icon={icons.download}
					onclick={() => void updater.installNow()}
				>
					{t('updates.install')}
				</Button>
			{:else if updater.state === 'ready'}
				<Button
					variant="primary"
					icon={icons.refresh}
					onclick={() => void updater.restartNow()}
				>
					{t('updates.restart')}
				</Button>
			{:else}
				<Button
					variant="outline"
					icon={icons.refresh}
					disabled={updater.state === 'checking' || updater.state === 'downloading'}
					onclick={() => void updater.checkNow(false)}
				>
					{updater.state === 'checking' ? t('updates.checkingShort') : t('updates.check')}
				</Button>
			{/if}
		</SettingsRow>

		<SettingsRow
			label={t('updates.history')}
			description={t('updates.historyBody')}
		>
			<Button
				variant="ghost"
				size="sm"
				icon={icons.history}
				href="https://github.com/FacileStudio/Courrier/releases"
				target="_blank"
				rel="noopener noreferrer"
			>
				{t('updates.viewReleases')}
			</Button>
		</SettingsRow>
	</SettingsSection>
</div>
