<script lang="ts">
	import { Roll, Stage, Tile, Window, type Act } from '@facile/mise-en-scene';
	import { Avatar } from '@facile/muse';
	import { t } from '$lib/i18n/index.svelte';

	const ACCOUNTS = [
		{ email: 'camille@gmail.com', name: 'Camille', unread: '3' },
		{ email: 'contact@atelier.fr', name: 'Atelier', unread: '5' },
		{ email: 'compta@atelier.fr', name: 'Compta', unread: '1' }
	];

	let active = $state(0);

	function reset() {
		active = 0;
	}

	async function script(act: Act) {
		await act.wait(1400);
		active = 1;
		await act.wait(1800);
		active = 2;
		await act.wait(1800);
		reset();
		await act.wait(1200);
	}
</script>

<Tile title={t('landing.cardAccountsTitle')} body={t('landing.cardAccountsBody')}>
	<Stage width={320} height={240} {script} {reset} poster={() => (active = 1)}>
		{#snippet children()}
			<Window style="top:20px;left:20px;width:280px;height:204px">
				<div class="at-body">
					<div class="at-current">
						<Avatar name={ACCOUNTS[active].name} size="sm" />
						<Roll value={ACCOUNTS[active].email} class="at-email" />
					</div>
					<div class="at-list">
						<span class="at-mark" style="transform: translateY({active * 36}px)"></span>
						{#each ACCOUNTS as account, i (account.email)}
							<div class="at-row" class:at-on={i === active}>
								<span class="at-dot"></span>
								<span class="at-text">{account.email}</span>
								<span class="at-count">{account.unread}</span>
							</div>
						{/each}
					</div>
				</div>
			</Window>
		{/snippet}
	</Stage>
</Tile>

<style>
	.at-body {
		display: flex;
		flex-direction: column;
		gap: 8px;
		padding: 10px 12px;
		color: var(--color-fc-fg);
		font-family: var(--font-fc-body);
		font-size: 12px;
	}
	.at-current {
		display: flex;
		align-items: center;
		gap: 8px;
		height: 32px;
		font-weight: 700;
	}
	.at-current :global(.at-email) {
		flex: 1;
	}
	.at-list {
		position: relative;
	}
	.at-mark {
		position: absolute;
		top: 0;
		left: 0;
		width: 100%;
		height: 36px;
		border-radius: 10px;
		background: var(--color-fc-accent);
		transition: transform 0.5s cubic-bezier(0.45, 0, 0.1, 1);
	}
	.at-row {
		position: relative;
		display: flex;
		align-items: center;
		gap: 8px;
		height: 36px;
		padding: 0 10px;
		transition: color 0.3s ease;
	}
	.at-on {
		color: var(--color-fc-accent-fg);
	}
	.at-dot {
		flex: none;
		width: 8px;
		height: 8px;
		border-radius: 999px;
		background: currentColor;
		opacity: 0.4;
	}
	.at-text {
		flex: 1;
		overflow: hidden;
		white-space: nowrap;
		text-overflow: ellipsis;
	}
	.at-count {
		font-variant-numeric: tabular-nums;
		font-weight: 700;
	}
</style>
