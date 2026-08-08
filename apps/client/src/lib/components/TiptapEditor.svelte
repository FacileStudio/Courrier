<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { Editor } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import Placeholder from '@tiptap/extension-placeholder';
	import Link from '@tiptap/extension-link';
	import Underline from '@tiptap/extension-underline';
	import { IconButton, icons } from '@facile/muse';

	let {
		content = '',
		onchange,
		placeholder = 'Write your message...'
	}: {
		content?: string;
		onchange?: (html: string) => void;
		placeholder?: string;
	} = $props();

	let element: HTMLDivElement;
	let editor = $state<Editor | null>(null);

	onMount(() => {
		editor = new Editor({
			element,
			extensions: [
				StarterKit,
				Placeholder.configure({ placeholder }),
				Link.configure({ openOnClick: false }),
				Underline
			],
			content,
			onUpdate: ({ editor: e }) => {
				onchange?.(e.getHTML());
			},
			onTransaction: ({ editor: e }) => {
				editor = e;
			}
		});
	});

	onDestroy(() => {
		editor?.destroy();
	});

	function setLink() {
		if (!editor) return;
		const previousUrl = editor.getAttributes('link').href;
		const url = window.prompt('URL', previousUrl);
		if (url === null) return;
		if (url === '') {
			editor.chain().focus().extendMarkRange('link').unsetLink().run();
			return;
		}
		editor.chain().focus().extendMarkRange('link').setLink({ href: url }).run();
	}

	/*
	 * Formatting glyphs muse's `icons` map has no key for. `mdi:` for the two list marks and
	 * the quote mark because Solar ships no `-linear` equivalent of any of the three — the
	 * same reason CHARTE §8 already sends plus/close/chevrons to MDI.
	 */
	const marks = [
		{ name: 'bold', label: 'Bold', icon: 'solar:text-bold-linear', run: () => editor?.chain().focus().toggleBold().run() },
		{ name: 'italic', label: 'Italic', icon: 'solar:text-italic-linear', run: () => editor?.chain().focus().toggleItalic().run() },
		{ name: 'underline', label: 'Underline', icon: 'solar:text-underline-linear', run: () => editor?.chain().focus().toggleUnderline().run() },
		{ name: 'strike', label: 'Strikethrough', icon: 'solar:text-cross-linear', run: () => editor?.chain().focus().toggleStrike().run() },
		{ separator: true },
		{ name: 'link', label: 'Link', icon: 'solar:link-linear', run: setLink },
		{ separator: true },
		{ name: 'bulletList', label: 'Bullet list', icon: 'solar:list-linear', run: () => editor?.chain().focus().toggleBulletList().run() },
		{ name: 'orderedList', label: 'Numbered list', icon: 'mdi:format-list-numbered', run: () => editor?.chain().focus().toggleOrderedList().run() },
		{ separator: true },
		{ name: 'blockquote', label: 'Quote', icon: 'mdi:format-quote-close', run: () => editor?.chain().focus().toggleBlockquote().run() },
		{ name: 'code', label: 'Code', icon: icons.code, run: () => editor?.chain().focus().toggleCode().run() }
	] as const;
</script>

<div class="flex min-h-0 flex-1 flex-col">
	{#if editor}
		<!--
			muse's `IconButton` has no pressed/active variant, so an active mark switches from
			`ghost` to `default` (which carries the outline) and says so through `aria-pressed`.
		-->
		<div class="flex flex-wrap items-center gap-0.5 border-b border-fc-border px-3 py-1.5">
			{#each marks as mark, i (i)}
				{#if 'separator' in mark}
					<span class="mx-1 h-5 w-px shrink-0 bg-fc-border"></span>
				{:else}
					{@const active = editor?.isActive(mark.name) ?? false}
					<IconButton
						variant={active ? 'default' : 'ghost'}
						aria-label={mark.label}
						aria-pressed={active}
						title={mark.label}
						onclick={mark.run}
					>
						<iconify-icon icon={mark.icon} width="18" height="18" class="block size-4.5"
						></iconify-icon>
					</IconButton>
				{/if}
			{/each}
		</div>
	{/if}

	<div class="flex-1 overflow-auto px-6 py-4">
		<div bind:this={element} class="tiptap-editor h-full"></div>
	</div>
</div>

<style>
	/* The editor's content is prose, not chrome: styled straight from the tokens, never with
	   a muse component's classes. */
	:global(.tiptap-editor .tiptap) {
		outline: none;
		min-height: 100%;
		font-family: var(--font-fc-body);
		font-size: var(--text-fc-sm);
		line-height: 1.625;
		color: var(--color-fc-fg);
	}

	:global(.tiptap-editor .tiptap p.is-editor-empty:first-child::before) {
		content: attr(data-placeholder);
		float: left;
		height: 0;
		pointer-events: none;
		color: var(--color-fc-fg-muted);
	}

	:global(.tiptap-editor .tiptap p) {
		margin: 0.25em 0;
	}

	:global(.tiptap-editor .tiptap ul),
	:global(.tiptap-editor .tiptap ol) {
		padding-left: 1.5em;
		margin: 0.5em 0;
	}

	:global(.tiptap-editor .tiptap ul) {
		list-style: disc;
	}

	:global(.tiptap-editor .tiptap ol) {
		list-style: decimal;
	}

	:global(.tiptap-editor .tiptap blockquote) {
		border-left: 2px solid var(--color-fc-border);
		padding-left: 1em;
		margin: 0.5em 0;
		color: var(--color-fc-fg-muted);
	}

	:global(.tiptap-editor .tiptap code) {
		background: var(--color-fc-surface);
		border-radius: var(--radius-fc-xs);
		padding: 0.15em 0.3em;
		font-family: var(--font-fc-mono);
		font-size: 0.85em;
	}

	:global(.tiptap-editor .tiptap a) {
		color: var(--color-fc-accent);
		text-decoration: underline;
	}

	:global(.tiptap-editor .tiptap strong) {
		font-weight: 600;
	}
</style>
