<script lang="ts">
	import { page } from '$app/state';
	import type { Folder, UserProfile } from '$lib/backend';

	let { folders = [], user = null }: { folders?: Folder[]; user?: UserProfile | null } = $props();

	let avatarFailed = $state(false);

	$effect(() => {
		void user?.avatar_url;
		avatarFailed = false;
	});

	function getInitials(value: string) {
		const parts = value.trim().split(/\s+/).filter(Boolean);
		if (parts.length === 0) return '?';
		if (parts.length === 1) return parts[0].slice(0, 2).toUpperCase();
		return `${parts[0][0] ?? ''}${parts[1][0] ?? ''}`.toUpperCase();
	}

	const userLabel = $derived(user?.name?.trim() || user?.email || '');
	const profileActive = $derived(isActive('/settings') || isActive('/profile'));

	const items = [
		{ href: '/mail', label: 'Mail', icon: 'solar:letter-linear' },
		{ href: '/mail/compose', label: 'Compose', icon: 'solar:pen-new-square-linear' }
	];

	function isActive(href: string) {
		return page.url.pathname === href || page.url.pathname.startsWith(href + '/');
	}

	function isItemActive(href: string) {
		if (href === '/mail/compose') return page.url.pathname === '/mail/compose';
		if (href === '/mail') {
			return (
				page.url.pathname === '/mail' ||
				(page.url.pathname.startsWith('/mail/') && page.url.pathname !== '/mail/compose')
			);
		}
		return isActive(href);
	}

	const inboxUnread = $derived(folders.find((f) => f.type === 'inbox')?.unread_count ?? 0);
</script>

<nav
	class="fixed inset-x-0 z-50 flex justify-center px-4 md:hidden"
	style="bottom: max(0.75rem, env(safe-area-inset-bottom))"
>
	<div
		class="flex items-center gap-1 rounded-full border border-border/40 bg-background/55 p-1.5 shadow-lg shadow-black/10 ring-1 ring-white/10 backdrop-blur-2xl backdrop-saturate-150"
	>
		{#each items as item (item.href)}
			{@const active = isItemActive(item.href)}
			<a
				href={item.href}
				aria-label={item.label}
				title={item.label}
				class="relative flex items-center justify-center rounded-full px-3.5 py-2 transition-all duration-200 {active
					? 'bg-foreground text-background shadow-sm'
					: 'text-muted-foreground hover:bg-muted/60 hover:text-foreground'}"
			>
				<iconify-icon icon={item.icon} width="22"></iconify-icon>
				{#if item.href === '/mail' && inboxUnread > 0 && !active}
					<span
						class="absolute right-1 top-0 flex h-4 min-w-4 items-center justify-center rounded-full bg-foreground px-1 text-[9px] font-semibold text-background"
					>
						{inboxUnread > 99 ? '99+' : inboxUnread}
					</span>
				{/if}
			</a>
		{/each}

		<a
			href="/settings"
			aria-label="Profile"
			title="Profile"
			class="flex items-center justify-center rounded-full px-2.5 py-1.5 transition-all duration-200 {profileActive
				? 'bg-foreground shadow-sm'
				: 'hover:bg-muted/60'}"
		>
			{#if user?.avatar_url && !avatarFailed}
				<img
					src={user.avatar_url}
					alt={userLabel}
					class="h-7 w-7 rounded-full border border-border object-cover"
					onerror={() => (avatarFailed = true)}
				/>
			{:else}
				<span
					class="flex h-7 w-7 items-center justify-center rounded-full border border-border bg-foreground text-[10px] font-semibold text-background"
				>
					{getInitials(userLabel)}
				</span>
			{/if}
		</a>
	</div>
</nav>
