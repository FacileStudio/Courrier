<script lang="ts">
	import { onMount } from 'svelte';
	import type { Editor } from '@tiptap/core';
	import { Blob } from '@facile/blob';
	import { Button, Input, Spinner, Tooltip, toast } from '@facile/muse';
	import { backend, type AssistAction, type AssistRequest } from '$lib/backend';
	import { assistStatus } from '$lib/assist/status.svelte';
	import { reasonText, textToHtml, upgradeLink } from '$lib/assist/text';
	import { t } from '$lib/i18n/index.svelte';

	let { editor, context = '' }: { editor: Editor; context?: string } = $props();

	type Range = { from: number; to: number };
	type Result = { text: string; range: Range; request: AssistRequest };

	const languages = $derived([
		{ value: 'English', label: t('assist.langEnglish') },
		{ value: 'French', label: t('assist.langFrench') },
		{ value: 'Spanish', label: t('assist.langSpanish') },
		{ value: 'German', label: t('assist.langGerman') },
		{ value: 'Italian', label: t('assist.langItalian') }
	]);
	const rewrites = $derived<{ action: AssistAction; label: string }[]>([
		{ action: 'rewrite', label: t('assist.rewrite') },
		{ action: 'shorten', label: t('assist.shorten') },
		{ action: 'expand', label: t('assist.expand') },
		{ action: 'fix', label: t('assist.fix') },
		{ action: 'formal', label: t('assist.formal') },
		{ action: 'friendly', label: t('assist.friendly') }
	]);

	let open = $state(false);
	let busy = $state(false);
	let result = $state<Result | null>(null);
	let instructions = $state('');
	let bubble = $state<{ top: number; left: number } | null>(null);
	let root = $state<HTMLDivElement | null>(null);

	const status = $derived(assistStatus.current);
	const available = $derived(status?.available ?? false);
	const upgrade = $derived(upgradeLink(status));
	let hasSelection = $state(false);
	const scope = $derived(hasSelection ? t('assist.scopeSelection') : t('assist.scopeAll'));

	onMount(() => {
		void assistStatus.load();
		const place = () => {
			const { from, to, empty } = editor.state.selection;
			hasSelection = !empty;
			if (empty || !editor.isFocused || to - from < 3) {
				bubble = null;
				return;
			}
			const start = editor.view.coordsAtPos(from);
			bubble = { top: Math.max(start.top - 40, 8), left: start.left };
		};
		const hide = () => (bubble = null);
		editor.on('selectionUpdate', place);
		editor.on('blur', hide);
		window.addEventListener('scroll', hide, true);
		return () => {
			editor.off('selectionUpdate', place);
			editor.off('blur', hide);
			window.removeEventListener('scroll', hide, true);
		};
	});

	function target(): { text: string; range: Range } {
		const { from, to, empty } = editor.state.selection;
		if (!empty) return { text: editor.state.doc.textBetween(from, to, '\n\n'), range: { from, to } };
		return { text: editor.getText({ blockSeparator: '\n\n' }), range: { from: 0, to: editor.state.doc.content.size } };
	}

	async function run(request: Omit<AssistRequest, 'text' | 'context'>, range?: Range) {
		const picked = range ? { text: editor.state.doc.textBetween(range.from, range.to, '\n\n'), range } : target();
		const full: AssistRequest = { ...request, text: picked.text, context: context || undefined };
		open = false;
		bubble = null;
		busy = true;
		try {
			const out = await backend.assist(full);
			result = { text: out.text, range: picked.range, request: full };
		} catch (err) {
			assistStatus.refused(err);
			toast.danger(err instanceof Error ? err.message : t('assist.failed'));
		}
		busy = false;
	}

	function replace() {
		if (!result) return;
		const { from, to } = result.range;
		const html = textToHtml(result.text);
		if (from === 0 && to >= editor.state.doc.content.size) editor.chain().focus().setContent(html, { emitUpdate: true }).run();
		else editor.chain().focus().insertContentAt({ from, to }, html).run();
		result = null;
	}

	function insertBelow() {
		if (!result) return;
		editor.chain().focus().insertContentAt(result.range.to, textToHtml(result.text)).run();
		result = null;
	}

	function custom(event: SubmitEvent) {
		event.preventDefault();
		if (!instructions.trim()) return;
		void run({ action: 'rewrite', instructions: instructions.trim() });
		instructions = '';
	}

	function outside(event: MouseEvent) {
		if (open && root && !root.contains(event.target as Node)) open = false;
	}
</script>

<svelte:window onclick={outside} />

<div class="relative ml-auto" bind:this={root}>
	<Tooltip text={available ? t('assist.askAbout', { scope }) : reasonText(status)} class="contents">
		<button
			type="button"
			class="flex h-8 items-center gap-1.5 rounded-fc-pill px-2 text-fc-xs font-medium text-fc-fg-muted hover:bg-fc-surface hover:text-fc-fg focus-visible:outline-2 focus-visible:outline-fc-ring"
			aria-haspopup="menu"
			aria-expanded={open}
			onclick={() => (open = !open)}
		>
			{#if busy}
				<Spinner size="sm" />
			{:else}
				<Blob name="Bulle" size={20} class="size-5" />
			{/if}
			<span>Bulle</span>
		</button>
	</Tooltip>

	{#if open}
		<div
			role="menu"
			class="absolute right-0 top-9 z-30 flex w-64 flex-col gap-0.5 rounded-fc-lg bg-fc-component p-1.5 text-fc-sm shadow-lg"
		>
			{#if !available}
				<p class="px-2 py-2 text-fc-fg-muted">{reasonText(status)}</p>
				{#if upgrade}
					<a role="menuitem" class="menu-item font-medium" href={upgrade} target="_blank" rel="noopener">{t('plan.upgrade')}</a>
				{/if}
			{:else}
				<p class="px-2 pb-1 pt-0.5 text-fc-xs text-fc-fg-muted">{t('assist.worksOn', { scope })}</p>
				{#if context}
					<button role="menuitem" type="button" class="menu-item" onclick={() => run({ action: 'reply' }, { from: editor.state.selection.from, to: editor.state.selection.to })}>
						{t('assist.draftReply')}
					</button>
				{/if}
				{#each rewrites as item (item.action)}
					<button role="menuitem" type="button" class="menu-item" onclick={() => run({ action: item.action })}>{item.label}</button>
				{/each}
				<details class="group">
					<summary class="menu-item cursor-pointer list-none">{t('assist.translate')}</summary>
					<div class="flex flex-col pl-3">
						{#each languages as language (language.value)}
							<button role="menuitem" type="button" class="menu-item" onclick={() => run({ action: 'translate', language: language.value })}>{language.label}</button>
						{/each}
					</div>
				</details>
				<form class="mt-1 flex flex-col gap-1 px-1 pt-2" onsubmit={custom}>
					<Input bind:value={instructions} placeholder={t('assist.instructions')} aria-label={t('assist.instructions')} maxlength={500} />
				</form>
			{/if}
		</div>
	{/if}
</div>

{#if bubble && available && !open && !busy}
	<button
		type="button"
		class="fixed z-40 flex h-8 items-center gap-1 rounded-fc-pill bg-fc-component px-2 text-fc-xs font-medium shadow-md hover:bg-fc-surface"
		style="top: {bubble.top}px; left: {bubble.left}px"
		onmousedown={(event) => event.preventDefault()}
		onclick={(event) => {
			event.stopPropagation();
			open = true;
		}}
	>
		<Blob name="Bulle" size={18} class="size-4.5" />
		{t('assist.askBulle')}
	</button>
{/if}

{#if result}
	<div class="mt-1 flex basis-full flex-col gap-2 rounded-fc-lg bg-fc-surface p-3" role="region" aria-label={t('assist.suggestion')}>
		<div class="flex items-center gap-2 text-fc-xs text-fc-fg-muted">
			<Blob name="Bulle" size={18} class="size-4.5" />
			{t('assist.suggestionFor', { scope: result.range.from === 0 ? t('assist.scopeAll') : t('assist.scopeSelection') })}
		</div>
		<p class="max-h-60 overflow-auto whitespace-pre-wrap text-fc-sm text-fc-fg">{result.text}</p>
		<div class="flex flex-wrap gap-2">
			<Button size="sm" onclick={replace}>{result.request.action === 'reply' ? t('assist.insert') : t('assist.replace')}</Button>
			{#if result.request.action !== 'reply'}
				<Button size="sm" variant="secondary" onclick={insertBelow}>{t('assist.insertBelow')}</Button>
			{/if}
			<Button size="sm" variant="ghost" disabled={busy} onclick={() => result && run(result.request, result.range)}>{t('assist.tryAgain')}</Button>
			<Button size="sm" variant="ghost" onclick={() => (result = null)}>{t('assist.discard')}</Button>
		</div>
	</div>
{/if}

<style>
	.menu-item {
		display: block;
		width: 100%;
		border-radius: var(--radius-fc-md);
		padding: 0.4rem 0.5rem;
		text-align: left;
		color: var(--color-fc-fg);
	}

	.menu-item:hover,
	.menu-item:focus-visible {
		background: var(--color-fc-surface);
		outline: none;
	}
</style>
