<script lang="ts">
	import { tick } from 'svelte';
	import { Bar, Cursor, Roll, Stage, Window, cursorAt, type Act } from '@facile/mise-en-scene';
	import { Icon } from '@facile/muse';
	import { t } from '$lib/i18n/index.svelte';
	import InboxList from './InboxList.svelte';
	import InboxThread from './InboxThread.svelte';

	const SENDER = 'Studio Brume';

	let arrived = $state(false);
	let opened = $state(false);
	let summarizing = $state(false);
	let summary = $state(false);
	let replying = $state(false);
	let draft = $state('');
	let sent = $state(false);
	let cursor = $state(cursorAt());
	let row = $state<HTMLElement>();
	let summarizeButton = $state<HTMLElement>();
	let replyButton = $state<HTMLElement>();
	let sendButton = $state<HTMLElement>();

	const unread = $derived(String(arrived && !opened ? 3 : 2));

	function reset() {
		arrived = false;
		opened = false;
		summarizing = false;
		summary = false;
		replying = false;
		draft = '';
		sent = false;
		cursor.visible = false;
	}

	function poster() {
		arrived = true;
		opened = true;
		summary = true;
		sent = true;
	}

	async function press(act: Act, target: HTMLElement | undefined, onPress: () => void) {
		await act.move(cursor, act.point(target, 0.6, 0.65));
		await act.click(cursor, onPress);
		await tick();
	}

	async function script(act: Act) {
		await act.wait(900);
		arrived = true;
		await act.wait(900);
		cursor.visible = true;
		await press(act, row, () => (opened = true));
		await act.wait(800);
		await press(act, summarizeButton, () => (summarizing = true));
		await act.wait(450);
		summary = true;
		summarizing = false;
		await act.wait(2000);
		await press(act, replyButton, () => (replying = true));
		await act.wait(250);
		await act.type(t('landing.demoReply'), (v) => (draft = v), 38);
		await act.wait(350);
		await press(act, sendButton, () => {
			replying = false;
			sent = true;
		});
		await act.wait(2400);
		reset();
		await act.wait(600);
	}
</script>

{#snippet list()}
	<InboxList sender={SENDER} {arrived} {opened} bind:row />
{/snippet}

{#snippet thread()}
	<InboxThread
		sender={SENDER}
		{opened}
		{summarizing}
		{summary}
		{replying}
		{draft}
		{sent}
		bind:summarizeButton
		bind:replyButton
		bind:sendButton
	/>
{/snippet}

<Stage
	width={960}
	height={490}
	compact={{ below: 640, width: 358, height: 488 }}
	{script}
	{reset}
	{poster}
	label={t('landing.demoLabel')}
>
	{#snippet children({ compact })}
		{#if compact}
			<Window style="top:12px;left:12px;width:334px;height:446px">
				<div class="id-compact">
					{#if opened}{@render thread()}{:else}{@render list()}{/if}
				</div>
			</Window>
		{:else}
			<Window style="top:16px;left:40px;width:880px;height:452px">
				{#snippet sidebar()}
					<div class="id-folders">
						<div class="id-folder id-active">
							<Icon icon="solar:inbox-linear" size={16} />
							<span class="id-name">{t('landing.demoInbox')}</span>
							<Roll value={unread} class="id-count" />
						</div>
						<div class="id-folder"><Bar w="62%" /></div>
						<div class="id-folder"><Bar w="48%" /></div>
						<div class="id-folder"><Bar w="56%" /></div>
					</div>
				{/snippet}
				<div class="id-body">
					<div class="id-panel">{@render list()}</div>
					<div class="id-panel">{@render thread()}</div>
				</div>
			</Window>
		{/if}
		<Cursor {cursor} />
	{/snippet}
</Stage>

<style>
	.id-body,
	.id-compact {
		display: grid;
		flex: 1;
		grid-template-columns: 300px minmax(0, 1fr);
		gap: 8px;
		min-height: 0;
		padding: 8px;
		color: var(--color-fc-fg);
		font-family: var(--font-fc-body);
	}
	.id-compact {
		display: block;
		padding: 4px;
	}
	.id-panel {
		min-height: 0;
		overflow: hidden;
		border-radius: 12px;
		background: var(--color-fc-component);
	}
	.id-folders {
		display: flex;
		flex-direction: column;
		gap: 4px;
		padding: 12px 10px;
		font-family: var(--font-fc-body);
	}
	.id-folder {
		display: flex;
		align-items: center;
		gap: 10px;
		height: 36px;
		padding: 0 12px;
		border-radius: 10px;
		font-size: 13px;
	}
	.id-active {
		background: var(--color-fc-surface);
		color: var(--color-fc-fg);
		font-weight: 700;
	}
	.id-name {
		flex: 1;
	}
	.id-folder :global(.id-count) {
		width: 14px;
		text-align: right;
		font-variant-numeric: tabular-nums;
	}
</style>
