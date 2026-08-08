<script lang="ts">
	import { Avatar, Badge, Checkbox, Divider, icons } from '@facile/muse';
	import type { EmailMessage } from '$lib/backend';

	/* Glyphs an email client needs that muse's `icons` map has no key for yet. */
	const mailIcons = {
		reply: 'solar:reply-linear',
		forward: 'solar:forward-linear',
		archive: 'solar:archive-linear',
		mailOpen: 'solar:letter-opened-linear',
		star: 'solar:star-linear'
	};

	let {
		email,
		selected = false,
		checked = false,
		selectionActive = false,
		onopen,
		ontogglecheck,
		onreply,
		onforward,
		onarchive,
		ondelete,
		ontoggleread,
		ontogglestar
	}: {
		email: EmailMessage;
		selected?: boolean;
		checked?: boolean;
		selectionActive?: boolean;
		onopen: () => void;
		ontogglecheck: () => void;
		onreply?: () => void;
		onforward?: () => void;
		onarchive?: () => void;
		ondelete?: () => void;
		ontoggleread?: () => void;
		ontogglestar?: () => void;
	} = $props();

	const MENU_WIDTH = 224;
	const MENU_HEIGHT = 260;

	let menuOpen = $state(false);
	let menuX = $state(0);
	let menuY = $state(0);

	const sender = $derived(email.from_name || email.from_address);

	function formatDate(dateStr: string) {
		const date = new Date(dateStr);
		const now = new Date();
		const diffMs = now.getTime() - date.getTime();
		const diffDays = diffMs / (1000 * 60 * 60 * 24);

		if (diffDays < 1) {
			return date.toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' });
		}
		if (diffDays < 7) {
			return date.toLocaleDateString('fr-FR', { weekday: 'short' });
		}
		return date.toLocaleDateString('fr-FR', { day: 'numeric', month: 'short' });
	}

	function handleClick() {
		if (selectionActive) {
			ontogglecheck();
			return;
		}
		onopen();
	}

	function handleCheckboxClick(e: Event) {
		e.stopPropagation();
		ontogglecheck();
	}

	function openMenu(e: MouseEvent) {
		e.preventDefault();
		menuX = Math.min(e.clientX, Math.max(8, window.innerWidth - MENU_WIDTH - 8));
		menuY = Math.min(e.clientY, Math.max(8, window.innerHeight - MENU_HEIGHT - 8));
		menuOpen = true;
	}

	function run(action?: () => void) {
		menuOpen = false;
		action?.();
	}
</script>

<svelte:window
	onkeydown={(e) => {
		if (menuOpen && e.key === 'Escape') menuOpen = false;
	}}
/>

<div class="relative border-b border-fc-border last:border-b-0">
	<button
		type="button"
		class="group flex w-full items-center gap-3 px-4 py-2.5 text-left transition-colors focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-fc-ring
			{selected ? 'bg-fc-accent text-fc-accent-fg' : checked ? 'bg-fc-surface text-fc-fg' : 'text-fc-fg hover:bg-fc-surface'}"
		onclick={handleClick}
		oncontextmenu={openMenu}
	>
		<span class="flex size-8 shrink-0 items-center justify-center">
			{#if selectionActive || checked}
				<span
					class="flex size-8 cursor-pointer items-center justify-center"
					onclick={handleCheckboxClick}
					onkeydown={(e) => {
						if (e.key === 'Enter' || e.key === ' ') handleCheckboxClick(e);
					}}
					role="checkbox"
					aria-checked={checked}
					aria-label="Select conversation"
					tabindex={0}
				>
					<Checkbox {checked} class="pointer-events-none" />
				</span>
			{:else}
				<span class="block group-hover:hidden">
					<Avatar name={sender} size="sm" />
				</span>
				<span
					class="hidden size-8 cursor-pointer items-center justify-center group-hover:flex"
					onclick={handleCheckboxClick}
					onkeydown={(e) => {
						if (e.key === 'Enter' || e.key === ' ') handleCheckboxClick(e);
					}}
					role="checkbox"
					aria-checked={checked}
					aria-label="Select conversation"
					tabindex={0}
				>
					<Checkbox {checked} class="pointer-events-none" />
				</span>
			{/if}
		</span>

		<span class="min-w-0 flex-1">
			<span class="flex items-center justify-between gap-2">
				<span class="flex min-w-0 items-center gap-1.5">
					<span class="truncate text-fc-sm {email.is_read ? '' : 'font-semibold'}">{sender}</span>
					{#if (email.message_count ?? 1) > 1}
						<Badge tone="neutral" class="tabular-nums">{email.message_count}</Badge>
					{/if}
				</span>
				<span class="shrink-0 text-fc-xs {selected ? 'text-fc-accent-fg/70' : 'text-fc-fg-muted'}">
					{formatDate(email.date)}
				</span>
			</span>
			<span
				class="mt-0.5 block truncate text-fc-sm {selected ? 'text-fc-accent-fg/70' : 'text-fc-fg-muted'}"
			>
				{email.subject || '(no subject)'}
			</span>
		</span>

		{#if !email.is_read}
			<span
				class="size-2 shrink-0 rounded-fc-pill {selected ? 'bg-fc-accent-fg' : 'bg-fc-accent'}"
				aria-label="Unread"
			></span>
		{/if}
	</button>
</div>

<!--
	muse ships no menu primitive (no DropdownMenu, ContextMenu or Popover), so the row menu is
	local: a floating surface on muse tokens, dismissed by Escape, by a click anywhere else and
	by a second right-click. Everything inside it is a muse component.
-->
{#if menuOpen}
	<div
		class="fixed inset-0 z-40"
		role="presentation"
		onclick={() => (menuOpen = false)}
		oncontextmenu={(e) => {
			e.preventDefault();
			menuOpen = false;
		}}
	></div>
	<div
		class="fixed z-50 w-56 rounded-fc-md border border-fc-border bg-fc-component p-1 shadow-lg"
		style="left: {menuX}px; top: {menuY}px"
		role="menu"
		tabindex="-1"
	>
		{#snippet item(label: string, icon: string, action: () => void, danger = false)}
			<button
				type="button"
				role="menuitem"
				class="flex w-full items-center gap-2.5 rounded-fc-sm px-2.5 py-2 text-left text-fc-sm transition-colors focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-fc-ring
					{danger
					? 'text-fc-danger hover:bg-fc-danger/10'
					: 'text-fc-fg-muted hover:bg-fc-surface hover:text-fc-fg'}"
				onclick={() => run(action)}
			>
				<iconify-icon {icon} width="16" height="16" class="block size-4 shrink-0"></iconify-icon>
				{label}
			</button>
		{/snippet}

		{#if onreply}
			{@render item('Reply', mailIcons.reply, onreply)}
		{/if}
		{#if onforward}
			{@render item('Forward', mailIcons.forward, onforward)}
		{/if}
		{#if onreply || onforward}
			<Divider class="my-1" />
		{/if}
		{#if ontoggleread}
			{@render item(
				email.is_read ? 'Mark as unread' : 'Mark as read',
				email.is_read ? icons.mail : mailIcons.mailOpen,
				ontoggleread
			)}
		{/if}
		{#if ontogglestar}
			{@render item(email.is_starred ? 'Unstar' : 'Star', mailIcons.star, ontogglestar)}
		{/if}
		{#if onarchive || ondelete}
			<Divider class="my-1" />
		{/if}
		{#if onarchive}
			{@render item('Archive', mailIcons.archive, onarchive)}
		{/if}
		{#if ondelete}
			{@render item('Delete', icons.remove, ondelete, true)}
		{/if}
	</div>
{/if}
