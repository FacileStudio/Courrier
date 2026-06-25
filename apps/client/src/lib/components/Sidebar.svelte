<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { backend, type UserProfile, type Folder } from '$lib/backend';
	import { MAIL_FOLDERS } from '$lib/mail-folders';
	import SpaceSwitcher from '$lib/components/SpaceSwitcher.svelte';

	let { user, folders = [] }: { user: UserProfile | null; folders?: Folder[] } = $props();

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

	function userLabel(currentUser: UserProfile | null) {
		return currentUser?.name?.trim() || currentUser?.email || '';
	}

	async function logout() {
		try { await backend.logout(); } catch {}
		goto('/login');
	}

	function folderUnread(type: string): number {
		const f = folders.find((f) => f.type === type);
		return f?.unread_count ?? 0;
	}

	const folderLinks = MAIL_FOLDERS;

	const manageLinks: { href: string; label: string; icon: string }[] = [
		{ href: '/accounts', label: 'Accounts', icon: 'solar:mailbox-linear' },
		{ href: '/templates', label: 'Templates', icon: 'solar:documents-linear' },
		{ href: '/spaces', label: 'Spaces', icon: 'solar:users-group-rounded-linear' }
	];
</script>

<aside class="sticky top-0 hidden h-[100dvh] w-60 flex-col border-r bg-background md:flex">
	<div class="flex items-center gap-3 px-5 pt-8 pb-4">
		<iconify-icon icon="solar:letter-bold-duotone" width="28" class="text-foreground"></iconify-icon>
		<span class="text-2xl font-bold font-heading tracking-tight">Courrier</span>
	</div>

	<SpaceSwitcher />

	<div class="px-3 py-3">
		<a
			href="/mail/compose"
			class="flex items-center justify-center gap-2 rounded-md bg-foreground px-3 py-2.5 text-sm font-medium text-background transition-opacity hover:opacity-90"
		>
			<iconify-icon icon="solar:pen-new-square-linear" width="16"></iconify-icon>
			Compose
		</a>
	</div>

	<nav class="flex flex-1 flex-col gap-1 overflow-y-auto px-3">
		{#each folderLinks as link}
			{@const active = page.url.pathname === link.href}
			{@const unread = folderUnread(link.type)}
			<a
				href={link.href}
				class="flex items-center gap-3 rounded-md px-3 py-2.5 text-sm transition-colors {active
					? 'bg-foreground text-background font-medium'
					: 'text-muted-foreground hover:bg-muted hover:text-foreground'}"
			>
				<iconify-icon icon={link.icon} width="16"></iconify-icon>
				<span class="flex-1">{link.label}</span>
				{#if unread > 0}
					<span class="text-xs font-medium">{unread}</span>
				{/if}
			</a>
		{/each}

		<div class="my-2 h-px bg-border"></div>

		{#each manageLinks as link}
			{@const active = page.url.pathname === link.href || page.url.pathname.startsWith(link.href + '/')}
			<a
				href={link.href}
				class="flex items-center gap-3 rounded-md px-3 py-2.5 text-sm transition-colors {active
					? 'bg-foreground text-background font-medium'
					: 'text-muted-foreground hover:bg-muted hover:text-foreground'}"
			>
				<iconify-icon icon={link.icon} width="16"></iconify-icon>
				<span class="flex-1">{link.label}</span>
			</a>
		{/each}
	</nav>

	<div class="h-px bg-border"></div>

	<div class="flex flex-col gap-2 p-4">
		<a
			href="/settings"
			class="flex items-center gap-3 rounded-xl border border-border/70 bg-muted/40 p-2.5 transition-colors hover:bg-muted"
		>
			{#if user?.avatar_url && !avatarFailed}
				<img
					src={user.avatar_url}
					alt={userLabel(user)}
					class="h-9 w-9 rounded-full border border-border object-cover shrink-0"
					onerror={() => (avatarFailed = true)}
				/>
			{:else}
				<div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full border border-border bg-foreground text-xs font-semibold text-background">
					{getInitials(userLabel(user))}
				</div>
			{/if}
			<div class="min-w-0 flex-1">
				<p class="truncate text-sm font-medium">{user?.name || 'Set your profile'}</p>
				<p class="truncate text-xs text-muted-foreground">{user?.email ?? ''}</p>
			</div>
		</a>
		<button
			onclick={logout}
			class="flex w-full items-center gap-2 rounded-md px-3 py-2 text-sm text-muted-foreground transition-colors hover:text-destructive hover:bg-destructive/10"
		>
			<iconify-icon icon="solar:logout-2-linear" width="16"></iconify-icon>
			Logout
		</button>
	</div>
</aside>
