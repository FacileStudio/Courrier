<script lang="ts">
	import { Bar, grow } from '@facile/mise-en-scene';
	import { Avatar, Icon } from '@facile/muse';
	import { t } from '$lib/i18n/index.svelte';

	let {
		sender,
		opened,
		replying,
		draft,
		sent,
		replyButton = $bindable(),
		sendButton = $bindable()
	}: {
		sender: string;
		opened: boolean;
		replying: boolean;
		draft: string;
		sent: boolean;
		replyButton?: HTMLElement;
		sendButton?: HTMLElement;
	} = $props();
</script>

{#snippet message(name: string, time: string, body: string)}
	<div class="it-message">
		<div class="it-head">
			<Avatar {name} size="sm" />
			<b>{name}</b>
			<span class="it-time">{time}</span>
		</div>
		<p>{body}</p>
	</div>
{/snippet}

{#if opened}
	<div class="it-thread">
		<div class="it-title">
			<h3>{t('landing.demoSubject')}</h3>
			<span bind:this={replyButton} class="it-reply" class:it-pressed={replying}>
				<Icon icon="solar:reply-linear" size={18} />
			</span>
		</div>
		{@render message(sender, '09:41', t('landing.demoBody'))}
		{#if sent}
			<div transition:grow={{ duration: 320 }}>
				<div class="it-pad">{@render message(t('landing.demoYou'), '09:44', t('landing.demoReply'))}</div>
			</div>
		{/if}
		{#if replying}
			<div transition:grow={{ duration: 320 }}>
				<div class="it-pad">
					<div class="it-composer">
						<span class="it-draft">{draft}<span class="it-caret"></span></span>
						<span bind:this={sendButton} class="it-send">{t('landing.demoSend')}</span>
					</div>
				</div>
			</div>
		{/if}
	</div>
{:else}
	<div class="it-empty"><Bar w="52%" h={12} /><Bar w="80%" /><Bar w="66%" /><Bar w="72%" /></div>
{/if}

<style>
	.it-thread {
		display: flex;
		flex-direction: column;
		padding: 20px;
	}
	.it-title {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 12px;
		margin-bottom: 12px;
	}
	.it-title h3 {
		overflow: hidden;
		font-family: var(--font-fc-title);
		font-size: 18px;
		font-weight: 700;
		white-space: nowrap;
		text-overflow: ellipsis;
	}
	.it-reply {
		display: inline-flex;
		flex: none;
		align-items: center;
		justify-content: center;
		width: 36px;
		height: 36px;
		border-radius: 999px;
		background: var(--color-fc-surface);
		transition: background 0.2s ease, color 0.2s ease;
	}
	.it-pressed,
	.it-send {
		background: var(--color-fc-accent);
		color: var(--color-fc-accent-fg);
	}
	.it-message {
		display: flex;
		flex-direction: column;
		gap: 10px;
		padding: 14px 16px;
		border-radius: 12px;
		background: var(--color-fc-page);
		font-size: 13px;
		line-height: 1.5;
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
	.it-pad {
		padding-top: 12px;
	}
	.it-composer {
		display: flex;
		align-items: flex-end;
		gap: 10px;
		padding: 10px 10px 10px 16px;
		border-radius: 20px;
		background: var(--color-fc-surface);
		font-size: 13px;
		line-height: 20px;
	}
	.it-draft {
		flex: 1;
		min-height: 20px;
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
