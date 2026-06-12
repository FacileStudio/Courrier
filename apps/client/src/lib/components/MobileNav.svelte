<script lang="ts">
	import { page } from '$app/state';
	import { Inbox, Send, PenLine, Trash2, User } from 'lucide-svelte';
	import type { Folder } from '$lib/backend';

	let { folders = [] }: { folders?: Folder[] } = $props();

	const inboxUnread = $derived(folders.find((f) => f.type === 'inbox')?.unread_count ?? 0);

	const items = [
		{ href: '/mail', label: 'Inbox', icon: Inbox, exact: true },
		{ href: '/mail/sent', label: 'Sent', icon: Send, exact: false },
		{ href: '/mail/compose', label: 'Compose', icon: PenLine, exact: false, primary: true },
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
		class="flex items-center gap-0.5 rounded-2xl border border-border/70 bg-background/90 p-1.5 shadow-xl backdrop-blur-md"
	>
		{#each items as item (item.href)}
			{@const active = isActive(item)}
			<a
				href={item.href}
				class="relative flex flex-col items-center gap-0.5 rounded-xl px-3 py-1.5 text-[10px] font-medium transition-colors {active
					? 'bg-foreground text-background'
					: 'text-muted-foreground hover:text-foreground'}"
			>
				<item.icon class="h-5 w-5" />
				<span>{item.label}</span>
				{#if item.href === '/mail' && inboxUnread > 0 && !active}
					<span
						class="absolute right-1.5 top-0.5 flex h-4 min-w-4 items-center justify-center rounded-full bg-foreground px-1 text-[9px] font-semibold text-background"
					>
						{inboxUnread > 99 ? '99+' : inboxUnread}
					</span>
				{/if}
			</a>
		{/each}
	</div>
</nav>
