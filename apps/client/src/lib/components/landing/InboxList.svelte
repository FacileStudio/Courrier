<script lang="ts">
	import { Bar, grow } from '@facile/mise-en-scene';
	import { Avatar } from '@facile/muse';
	import { t } from '$lib/i18n/index.svelte';

	let {
		sender,
		arrived,
		opened,
		row = $bindable()
	}: { sender: string; arrived: boolean; opened: boolean; row?: HTMLElement } = $props();

	const FILLER = [
		{ id: 1, w: '58%', s: '72%' },
		{ id: 2, w: '44%', s: '60%' },
		{ id: 3, w: '52%', s: '78%' },
		{ id: 4, w: '40%', s: '55%' }
	];
</script>

<div class="il-tabs">
	<span class="il-tab il-on">{t('landing.demoPrimary')}</span>
	<span class="il-tab">{t('landing.demoNewsletters')}</span>
</div>
<div class="il-rows">
	{#if arrived}
		<div transition:grow={{ duration: 340 }}>
			<div bind:this={row} class="il-row" class:il-selected={opened}>
				<Avatar name={sender} size="sm" />
				<span class="il-text">
					<span class="il-top">
						<b class:il-read={opened}>{sender}</b>
						<span class="il-muted il-time">09:41</span>
					</span>
					<span class="il-muted il-subject">{t('landing.demoSubject')}</span>
				</span>
				{#if !opened}<span class="il-dot"></span>{/if}
			</div>
		</div>
	{/if}
	{#each FILLER as item (item.id)}
		<div class="il-row">
			<span class="il-avatar"></span>
			<span class="il-text il-gap"><Bar w={item.w} /><Bar w={item.s} h={6} /></span>
		</div>
	{/each}
</div>

<style>
	.il-tabs {
		display: flex;
		gap: 4px;
		padding: 12px 12px 8px;
	}
	.il-tab {
		height: 28px;
		padding: 0 12px;
		border-radius: 999px;
		color: var(--color-fc-fg-muted);
		font-size: 12px;
		line-height: 28px;
		white-space: nowrap;
	}
	.il-on {
		background: var(--color-fc-surface);
		color: var(--color-fc-fg);
		font-weight: 700;
	}
	.il-rows {
		display: flex;
		flex-direction: column;
	}
	.il-row {
		display: flex;
		align-items: center;
		gap: 12px;
		height: 60px;
		padding: 0 16px;
		transition: background 0.2s ease, color 0.2s ease;
	}
	.il-selected {
		background: var(--color-fc-accent);
		color: var(--color-fc-accent-fg);
	}
	.il-text {
		display: flex;
		flex: 1;
		flex-direction: column;
		gap: 2px;
		min-width: 0;
		font-size: 13px;
	}
	.il-gap {
		gap: 8px;
	}
	.il-top {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 8px;
		font-weight: 700;
	}
	.il-read {
		font-weight: 400;
	}
	.il-muted {
		color: var(--color-fc-fg-muted);
		white-space: nowrap;
	}
	.il-selected .il-muted {
		color: inherit;
		opacity: 0.75;
	}
	.il-time {
		font-size: 11px;
		font-weight: 400;
		font-variant-numeric: tabular-nums;
	}
	.il-subject {
		overflow: hidden;
		text-overflow: ellipsis;
	}
	.il-dot {
		flex: none;
		width: 8px;
		height: 8px;
		border-radius: 999px;
		background: var(--color-fc-accent);
	}
	.il-avatar {
		flex: none;
		width: 32px;
		height: 32px;
		border-radius: 999px;
		background: var(--color-fc-surface);
	}
</style>
