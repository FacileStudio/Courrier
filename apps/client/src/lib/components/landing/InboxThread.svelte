<script lang="ts">
	import { Bar } from '@facile/mise-en-scene';
	import { Blob } from '@facile/blob';
	import { Avatar, Icon } from '@facile/muse';
	import { t } from '$lib/i18n/index.svelte';

	let {
		sender,
		opened,
		summarizing,
		summary,
		replying,
		draft,
		sent,
		summarizeButton = $bindable(),
		replyButton = $bindable(),
		sendButton = $bindable()
	}: {
		sender: string;
		opened: boolean;
		summarizing: boolean;
		summary: boolean;
		replying: boolean;
		draft: string;
		sent: boolean;
		summarizeButton?: HTMLElement;
		replyButton?: HTMLElement;
		sendButton?: HTMLElement;
	} = $props();

	const SUMMARY = ['landing.demoSummary1', 'landing.demoSummary2', 'landing.demoSummary3'] as const;
</script>

{#if opened}
	<div class="it-thread">
		<div class="it-title">
			<h3>{t('landing.demoSubject')}</h3>
			<span bind:this={summarizeButton} class="it-chip" class:it-pressed={summarizing}>
				<Blob name="Bulle" size={16} class="size-4" />
				{t('landing.demoSummarize')}
			</span>
			<span bind:this={replyButton} class="it-reply" class:it-pressed={replying}>
				<Icon icon="solar:reply-linear" size={18} />
			</span>
		</div>
		<div class="it-slot it-slot-a">
			<div class="it-layer it-earlier" class:it-hide={summary}>
				<span class="it-label">{t('landing.demoEarlier')}</span>
				<div class="it-fold"><span class="it-dot"></span><Bar w="64%" /></div>
				<div class="it-fold"><span class="it-dot"></span><Bar w="48%" /></div>
			</div>
			<div class="it-layer it-summary" class:it-hide={!summary}>
				<span class="it-label it-bulle"><Blob name="Bulle" size={16} class="size-4" />{t('assist.summary')}</span>
				{#each SUMMARY as key (key)}<p>{t(key)}</p>{/each}
			</div>
		</div>
		<div class="it-message">
			<div class="it-head"><Avatar name={sender} size="sm" /><b>{sender}</b><span class="it-time">09:41</span></div>
			<p>{t('landing.demoBody')}</p>
		</div>
		<div class="it-slot it-slot-b">
			<div class="it-layer it-composer" class:it-hide={!replying}>
				<span class="it-draft">{draft}<span class="it-caret"></span></span>
				<span bind:this={sendButton} class="it-send">{t('landing.demoSend')}</span>
			</div>
			<div class="it-layer it-message it-mine" class:it-hide={!sent}>
				<div class="it-head"><b>{t('landing.demoYou')}</b><span class="it-time">09:44</span></div>
				<p>{t('landing.demoReply')}</p>
			</div>
		</div>
	</div>
{:else}
	<div class="it-empty"><Bar w="52%" h={12} /><Bar w="80%" /><Bar w="66%" /><Bar w="72%" /></div>
{/if}

<style>
	.it-thread {
		display: flex;
		flex-direction: column;
		gap: 8px;
		padding: 16px;
		font-size: 13px;
		line-height: 1.45;
	}
	.it-title {
		display: flex;
		align-items: center;
		gap: 8px;
		height: 36px;
	}
	.it-title h3 {
		flex: 1;
		overflow: hidden;
		font-family: var(--font-fc-title);
		font-size: 17px;
		font-weight: 700;
		white-space: nowrap;
		text-overflow: ellipsis;
	}
	.it-chip,
	.it-reply {
		display: inline-flex;
		flex: none;
		align-items: center;
		justify-content: center;
		height: 32px;
		border-radius: 999px;
		background: var(--color-fc-surface);
		transition: background 0.2s ease, color 0.2s ease;
	}
	.it-chip {
		gap: 6px;
		padding: 0 12px;
		font-size: 12px;
		font-weight: 700;
	}
	.it-reply {
		width: 32px;
	}
	.it-pressed,
	.it-send {
		background: var(--color-fc-accent);
		color: var(--color-fc-accent-fg);
	}
	.it-slot {
		position: relative;
		flex: none;
	}
	.it-slot-a {
		height: 116px;
	}
	.it-slot-b {
		height: 86px;
	}
	.it-layer {
		position: absolute;
		inset: 0;
		transition: opacity 0.35s ease, transform 0.45s cubic-bezier(0.25, 0.8, 0.25, 1);
	}
	.it-hide {
		opacity: 0;
		transform: translateY(6px);
	}
	.it-earlier,
	.it-summary {
		display: flex;
		flex-direction: column;
		gap: 4px;
		padding: 10px 14px;
		border-radius: 12px;
		background: var(--color-fc-page);
	}
	.it-summary {
		background: var(--color-fc-surface);
	}
	.it-summary p {
		font-size: 12px;
		line-height: 17px;
	}
	.it-label {
		display: flex;
		align-items: center;
		gap: 6px;
		margin-bottom: 4px;
		color: var(--color-fc-fg-muted);
		font-size: 11px;
		font-weight: 700;
	}
	.it-bulle {
		color: var(--color-fc-fg);
	}
	.it-fold {
		display: flex;
		align-items: center;
		gap: 10px;
		height: 30px;
	}
	.it-dot {
		flex: none;
		width: 24px;
		height: 24px;
		border-radius: 999px;
		background: var(--color-fc-surface);
	}
	.it-message {
		display: flex;
		flex-direction: column;
		gap: 8px;
		padding: 12px 14px;
		border-radius: 12px;
		background: var(--color-fc-page);
	}
	.it-mine {
		gap: 4px;
	}
	.it-head {
		display: flex;
		align-items: center;
		gap: 10px;
	}
	.it-head b {
		flex: 1;
		font-weight: 700;
	}
	.it-time {
		color: var(--color-fc-fg-muted);
		font-size: 11px;
		font-variant-numeric: tabular-nums;
	}
	.it-composer {
		display: flex;
		align-items: center;
		gap: 10px;
		height: 52px;
		padding: 0 10px 0 16px;
		border-radius: 26px;
		background: var(--color-fc-surface);
	}
	.it-draft {
		flex: 1;
		line-height: 18px;
	}
	.it-caret {
		display: inline-block;
		width: 1.5px;
		height: 14px;
		margin-left: 1px;
		vertical-align: -2px;
		background: var(--color-fc-fg);
	}
	.it-send {
		flex: none;
		height: 32px;
		padding: 0 16px;
		border-radius: 999px;
		font-size: 12px;
		font-weight: 700;
		line-height: 32px;
	}
	.it-empty {
		display: grid;
		gap: 14px;
		padding: 28px 24px;
	}
</style>
