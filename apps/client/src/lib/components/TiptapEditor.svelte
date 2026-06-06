<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { Editor } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import Placeholder from '@tiptap/extension-placeholder';
	import Link from '@tiptap/extension-link';
	import Underline from '@tiptap/extension-underline';
	import { Button } from '$lib/components/ui/button';
	import {
		Bold,
		Italic,
		Underline as UnderlineIcon,
		Strikethrough,
		Link as LinkIcon,
		List,
		ListOrdered,
		Quote,
		Code
	} from 'lucide-svelte';

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
</script>

<div class="flex flex-col flex-1 min-h-0">
	{#if editor}
		<div class="flex items-center gap-0.5 border-b px-3 py-1.5">
			<Button
				variant="ghost"
				size="icon"
				class="h-7 w-7 {editor.isActive('bold') ? 'bg-muted' : ''}"
				onclick={() => editor?.chain().focus().toggleBold().run()}
			>
				<Bold class="h-3.5 w-3.5" />
			</Button>
			<Button
				variant="ghost"
				size="icon"
				class="h-7 w-7 {editor.isActive('italic') ? 'bg-muted' : ''}"
				onclick={() => editor?.chain().focus().toggleItalic().run()}
			>
				<Italic class="h-3.5 w-3.5" />
			</Button>
			<Button
				variant="ghost"
				size="icon"
				class="h-7 w-7 {editor.isActive('underline') ? 'bg-muted' : ''}"
				onclick={() => editor?.chain().focus().toggleUnderline().run()}
			>
				<UnderlineIcon class="h-3.5 w-3.5" />
			</Button>
			<Button
				variant="ghost"
				size="icon"
				class="h-7 w-7 {editor.isActive('strike') ? 'bg-muted' : ''}"
				onclick={() => editor?.chain().focus().toggleStrike().run()}
			>
				<Strikethrough class="h-3.5 w-3.5" />
			</Button>

			<div class="mx-1 h-4 w-px bg-border"></div>

			<Button
				variant="ghost"
				size="icon"
				class="h-7 w-7 {editor.isActive('link') ? 'bg-muted' : ''}"
				onclick={setLink}
			>
				<LinkIcon class="h-3.5 w-3.5" />
			</Button>

			<div class="mx-1 h-4 w-px bg-border"></div>

			<Button
				variant="ghost"
				size="icon"
				class="h-7 w-7 {editor.isActive('bulletList') ? 'bg-muted' : ''}"
				onclick={() => editor?.chain().focus().toggleBulletList().run()}
			>
				<List class="h-3.5 w-3.5" />
			</Button>
			<Button
				variant="ghost"
				size="icon"
				class="h-7 w-7 {editor.isActive('orderedList') ? 'bg-muted' : ''}"
				onclick={() => editor?.chain().focus().toggleOrderedList().run()}
			>
				<ListOrdered class="h-3.5 w-3.5" />
			</Button>

			<div class="mx-1 h-4 w-px bg-border"></div>

			<Button
				variant="ghost"
				size="icon"
				class="h-7 w-7 {editor.isActive('blockquote') ? 'bg-muted' : ''}"
				onclick={() => editor?.chain().focus().toggleBlockquote().run()}
			>
				<Quote class="h-3.5 w-3.5" />
			</Button>
			<Button
				variant="ghost"
				size="icon"
				class="h-7 w-7 {editor.isActive('code') ? 'bg-muted' : ''}"
				onclick={() => editor?.chain().focus().toggleCode().run()}
			>
				<Code class="h-3.5 w-3.5" />
			</Button>
		</div>
	{/if}

	<div class="flex-1 overflow-auto px-6 py-4">
		<div bind:this={element} class="tiptap-editor h-full"></div>
	</div>
</div>

<style>
	:global(.tiptap-editor .tiptap) {
		outline: none;
		min-height: 100%;
		font-size: 0.875rem;
		line-height: 1.625;
	}

	:global(.tiptap-editor .tiptap p.is-editor-empty:first-child::before) {
		content: attr(data-placeholder);
		float: left;
		height: 0;
		pointer-events: none;
		color: var(--color-muted-foreground);
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
		border-left: 2px solid var(--color-border);
		padding-left: 1em;
		margin: 0.5em 0;
		color: var(--color-muted-foreground);
	}

	:global(.tiptap-editor .tiptap code) {
		background: var(--color-muted);
		border-radius: 0.25em;
		padding: 0.15em 0.3em;
		font-size: 0.85em;
	}

	:global(.tiptap-editor .tiptap a) {
		color: var(--color-primary);
		text-decoration: underline;
	}

	:global(.tiptap-editor .tiptap strong) {
		font-weight: 600;
	}
</style>
