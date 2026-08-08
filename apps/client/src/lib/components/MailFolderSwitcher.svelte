<script lang="ts">
	import { page } from '$app/state';
	import { icons, twMerge } from '@facile/muse';
	import { MAIL_FOLDERS } from '$lib/mail-folders';

	let { folders = [] }: { folders?: { type: string; unread_count?: number }[] } = $props();

	/*
	 * Local on purpose: mail folders are Courrier's own vocabulary and this sits in the mail
	 * header as a compact heading, where muse's `Tabs` strip does not fit. Everything below is
	 * muse tokens and muse icons — no shadcn, no lucide.
	 */
	let open = $state(false);
	let rootEl: HTMLElement | null = $state(null);
	let triggerEl: HTMLButtonElement | null = $state(null);

	const current = $derived(MAIL_FOLDERS.find((f) => f.href === page.url.pathname) ?? MAIL_FOLDERS[0]);

	function unread(type: string): number {
		return folders.find((f) => f.type === type)?.unread_count ?? 0;
	}

	function handleClickOutside(e: MouseEvent) {
		if (rootEl && !rootEl.contains(e.target as Node)) open = false;
	}

	/* Escape has to hand focus back to the trigger: closing while focus sits on a row that is
	   about to leave the document drops a keyboard user at the top of the page. */
	function handleKeydown(e: KeyboardEvent) {
		if (e.key !== 'Escape') return;
		e.stopPropagation();
		open = false;
		triggerEl?.focus();
	}

	$effect(() => {
		if (!open) return;
		document.addEventListener('click', handleClickOutside);
		document.addEventListener('keydown', handleKeydown);
		return () => {
			document.removeEventListener('click', handleClickOutside);
			document.removeEventListener('keydown', handleKeydown);
		};
	});
</script>

<div bind:this={rootEl} class="relative">
	<button
		bind:this={triggerEl}
		type="button"
		class="flex min-h-11 items-center gap-1.5 rounded-fc-md px-1 text-fc-lg font-semibold text-fc-fg transition-colors hover:text-fc-fg focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-fc-ring"
		aria-expanded={open}
		onclick={() => (open = !open)}
	>
		<iconify-icon icon={current.icon} width="18" height="18" class="block size-[18px] shrink-0"
		></iconify-icon>
		<span>{current.label}</span>
		<iconify-icon
			icon={icons.chevronDown}
			width="16"
			height="16"
			class={twMerge('block size-4 shrink-0 text-fc-fg-muted transition-transform', open && 'rotate-180')}
		></iconify-icon>
	</button>

	{#if open}
		<div
			class="absolute left-0 top-full z-40 mt-1.5 flex w-56 flex-col overflow-hidden rounded-fc-md border border-fc-border bg-fc-component p-1 shadow-lg"
		>
			{#each MAIL_FOLDERS as folder (folder.href)}
				{@const count = unread(folder.type)}
				{@const active = folder.href === current.href}
				<a
					href={folder.href}
					aria-current={active ? 'page' : undefined}
					class={twMerge(
						'flex items-center gap-2.5 rounded-fc-sm px-2.5 py-2 text-fc-sm transition-colors',
						active
							? 'bg-fc-accent font-medium text-fc-accent-fg'
							: 'text-fc-fg hover:bg-fc-surface'
					)}
					onclick={() => (open = false)}
				>
					<iconify-icon icon={folder.icon} width="16" height="16" class="block size-4 shrink-0"
					></iconify-icon>
					<span class="min-w-0 flex-1 truncate">{folder.label}</span>
					{#if count > 0}
						<span class="shrink-0 text-fc-xs opacity-70">{count}</span>
					{/if}
				</a>
			{/each}
		</div>
	{/if}
</div>
