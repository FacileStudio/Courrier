<script lang="ts">
	import { page } from '$app/state';
	import { Inbox, Send, PenLine, Trash2, User } from 'lucide-svelte';
	import type { Folder } from '$lib/backend';

	let { folders = [] }: { folders?: Folder[] } = $props();

	const inboxUnread = $derived(folders.find((f) => f.type === 'inbox')?.unread_count ?? 0);

	const items = [
		{ href: '/mail', label: 'Inbox', icon: Inbox, exact: true },
		{ href: '/mail/sent', label: 'Sent', icon: Send, exact: false },
		{ href: '/mail/compose', label: 'Compose', icon: PenLine, exact: false },
		{ href: '/mail/trash', label: 'Trash', icon: Trash2, exact: false },
		{ href: '/profile', label: 'Profile', icon: User, exact: false }
	];

	function isActive(item: (typeof items)[number]) {
		return item.exact
			? page.url.pathname === item.href
			: page.url.pathname.startsWith(item.href);
	}
</script>

<nav
	class="fixed inset-x-0 z-50 flex justify-center px-4 md:hidden"
	style="bottom: max(0.75rem, env(safe-area-inset-bottom))"
>
	<div
		class="flex items-center gap-1 rounded-full border border-border/40 bg-background/55 p-1.5 shadow-lg shadow-black/10 ring-1 ring-white/10 backdrop-blur-2xl backdrop-saturate-150"
	>
		{#each items as item (item.href)}
			{@const active = isActive(item)}
			<a
				href={item.href}
				aria-label={item.label}
				title={item.label}
				class="relative flex items-center justify-center rounded-full px-3.5 py-2 transition-all duration-200 {active
					? 'bg-foreground text-background shadow-sm'
					: 'text-muted-foreground hover:bg-muted/60 hover:text-foreground'}"
			>
				<item.icon class="h-[22px] w-[22px]" />
				{#if item.href === '/mail' && inboxUnread > 0 && !active}
					<span
						class="absolute right-1 top-0 flex h-4 min-w-4 items-center justify-center rounded-full bg-foreground px-1 text-[9px] font-semibold text-background"
					>
						{inboxUnread > 99 ? '99+' : inboxUnread}
					</span>
				{/if}
			</a>
		{/each}
	</div>
</nav>
