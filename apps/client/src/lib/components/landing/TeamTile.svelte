<script lang="ts">
	import { Stage, Tile, Window, type Act } from '@facile/mise-en-scene';
	import { Avatar, Icon } from '@facile/muse';
	import { t } from '$lib/i18n/index.svelte';

	const ADDRESS = 'contact@atelier.fr';
	const MEMBERS = ['Camille', 'Léa', 'Hugo'];

	let read = $state(false);
	let replied = $state(false);

	function reset() {
		read = false;
		replied = false;
	}

	function poster() {
		read = true;
		replied = true;
	}

	async function script(act: Act) {
		await act.wait(1000);
		read = true;
		await act.wait(1400);
		replied = true;
		await act.wait(2800);
		reset();
		await act.wait(1100);
	}
</script>

<Tile title={t('landing.cardTeamTitle')} body={t('landing.cardTeamBody')}>
	<Stage width={320} height={240} {script} {reset} {poster}>
		{#snippet children()}
			<Window style="top:20px;left:20px;width:280px;height:204px">
				<div class="tm-body">
					<div class="tm-head">
						<span class="tm-who">
							<b>{t('landing.demoTeam')}</b>
							<span class="tm-muted">{ADDRESS}</span>
						</span>
						<span class="tm-members">
							{#each MEMBERS as name (name)}<Avatar {name} size="sm" class="size-6" />{/each}
						</span>
					</div>
					<div class="tm-card">
						<Avatar name="Nora Petit" size="sm" class="size-6" />
						<span class="tm-text"><b>Nora Petit</b><span class="tm-muted">{t('landing.demoQuote')}</span></span>
						<span class="tm-unread" class:tm-gone={read}></span>
					</div>
					<div class="tm-card tm-reply" class:tm-hide={!replied}>
						<span class="tm-text">
							<span class="tm-sent"><Icon icon="solar:check-circle-linear" size={14} />{t('landing.demoSent')} · {ADDRESS}</span>
							<span>{t('landing.demoQuoteReply')}</span>
						</span>
					</div>
				</div>
			</Window>
		{/snippet}
	</Stage>
</Tile>

<style>
	.tm-body {
		display: flex;
		flex-direction: column;
		gap: 8px;
		padding: 10px 12px;
		color: var(--color-fc-fg);
		font-family: var(--font-fc-body);
		font-size: 12px;
	}
	.tm-head {
		display: flex;
		align-items: center;
		gap: 8px;
		height: 32px;
	}
	.tm-who {
		display: flex;
		flex: 1;
		flex-direction: column;
		min-width: 0;
	}
	.tm-muted {
		overflow: hidden;
		color: var(--color-fc-fg-muted);
		white-space: nowrap;
		text-overflow: ellipsis;
	}
	.tm-members {
		display: flex;
		gap: 2px;
	}
	.tm-card {
		display: flex;
		align-items: center;
		gap: 8px;
		height: 48px;
		padding: 0 10px;
		border-radius: 10px;
		background: var(--color-fc-surface);
		transition: opacity 0.4s ease, transform 0.5s cubic-bezier(0.25, 0.8, 0.25, 1);
	}
	.tm-reply {
		height: 48px;
		margin-left: 24px;
		background: var(--color-fc-accent);
		color: var(--color-fc-accent-fg);
	}
	.tm-hide {
		opacity: 0;
		transform: translateY(8px);
	}
	.tm-unread {
		flex: none;
		width: 8px;
		height: 8px;
		margin-left: auto;
		border-radius: 999px;
		background: var(--color-fc-accent);
		transition: opacity 0.3s ease;
	}
	.tm-gone {
		opacity: 0;
	}
	.tm-text {
		display: flex;
		flex-direction: column;
		gap: 2px;
		min-width: 0;
	}
	.tm-reply .tm-text > span {
		overflow: hidden;
		white-space: nowrap;
		text-overflow: ellipsis;
	}
	.tm-sent {
		display: flex;
		align-items: center;
		gap: 4px;
		font-size: 11px;
		opacity: 0.8;
	}
</style>
