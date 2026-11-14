<script lang="ts">
	import { Bar, Roll, Stage, Tile, Window, type Act } from '@facile/mise-en-scene';
	import { Avatar, Icon } from '@facile/muse';
	import { t } from '$lib/i18n/index.svelte';

	type Row = { id: string; person?: string; icon?: string; w: string };

	const ROWS: Row[] = [
		{ id: 'lea', person: 'Léa Martin', w: '70%' },
		{ id: 'news', icon: 'solar:letter-linear', w: '58%' },
		{ id: 'marc', person: 'Marc Dubois', w: '62%' },
		{ id: 'bill', icon: 'solar:bell-linear', w: '46%' },
		{ id: 'nora', person: 'Nora Petit', w: '0' }
	];

	let sorted = $state(false);

	function slot(row: Row, index: number): number {
		if (!sorted) return index;
		return row.person ? ROWS.filter((r) => r.person).indexOf(row) : index;
	}

	function reset() {
		sorted = false;
	}

	async function script(act: Act) {
		await act.wait(1500);
		sorted = true;
		await act.wait(3200);
		reset();
		await act.wait(1300);
	}
</script>

<Tile title={t('landing.cardTabsTitle')} body={t('landing.cardTabsBody')}>
	<Stage width={320} height={240} {script} {reset} poster={() => (sorted = true)}>
		{#snippet children()}
			<Window style="top:20px;left:20px;width:280px;height:204px">
				<div class="tt-body">
					<div class="tt-tabs">
						<span class="tt-tab tt-on">{t('landing.demoPrimary')}</span>
						<span class="tt-tab">{t('landing.demoNewsletters')}<Roll value={sorted ? '4' : '3'} class="tt-count" /></span>
					</div>
					<div class="tt-rows">
						{#each ROWS as row, i (row.id)}
							<div
								class="tt-row"
								class:tt-out={sorted && !row.person}
								style="transform: translate({sorted && !row.person ? 48 : 0}px, {slot(row, i) * 36}px)"
							>
								{#if row.person}
									<Avatar name={row.person} size="sm" class="size-6" />
									<b>{row.person}</b>
								{:else}
									<span class="tt-icon"><Icon icon={row.icon ?? ''} size={14} /></span>
									<Bar w={row.w} />
								{/if}
							</div>
						{/each}
					</div>
				</div>
			</Window>
		{/snippet}
	</Stage>
</Tile>

<style>
	.tt-body {
		display: flex;
		flex-direction: column;
		gap: 8px;
		padding: 10px 12px;
		color: var(--color-fc-fg);
		font-family: var(--font-fc-body);
		font-size: 12px;
	}
	.tt-tabs {
		display: flex;
		gap: 4px;
	}
	.tt-tab {
		display: inline-flex;
		align-items: center;
		gap: 6px;
		height: 26px;
		padding: 0 10px;
		border-radius: 999px;
		color: var(--color-fc-fg-muted);
		white-space: nowrap;
	}
	.tt-on {
		background: var(--color-fc-surface);
		color: var(--color-fc-fg);
		font-weight: 700;
	}
	.tt-tab :global(.tt-count) {
		width: 10px;
		font-variant-numeric: tabular-nums;
	}
	.tt-rows {
		position: relative;
		height: 108px;
		overflow: hidden;
	}
	.tt-row {
		position: absolute;
		top: 0;
		left: 0;
		display: flex;
		align-items: center;
		gap: 8px;
		width: 100%;
		height: 36px;
		padding: 0 6px;
		transition: transform 0.55s cubic-bezier(0.45, 0, 0.1, 1), opacity 0.4s ease;
	}
	.tt-out {
		opacity: 0;
	}
	.tt-icon {
		display: inline-flex;
		flex: none;
		align-items: center;
		justify-content: center;
		width: 24px;
		height: 24px;
		border-radius: 999px;
		background: var(--color-fc-surface);
	}
</style>
