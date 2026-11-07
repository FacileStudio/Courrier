<script lang="ts">
	import { goto } from '$app/navigation';
	import { Avatar, icons, springPress, twMerge } from '@facile/muse';
	import { backend } from '$lib/backend';
	import { isApp } from '$lib/platform';
	import { t } from '$lib/i18n/index.svelte';

	type Entry = { href: string; label: string; icon: string };

	let {
		user,
		collapsed = false
	}: {
		user?: { name: string; email?: string; avatar?: string };
		collapsed?: boolean;
	} = $props();

	const entries: Entry[] = $derived([
		{ href: '/settings', label: t('shell.settings'), icon: icons.settings },
		...(isApp() ? [] : [{ href: '/downloads', label: t('shell.desktopApp'), icon: icons.download }])
	]);

	let open = $state(false);
	let root = $state<HTMLElement | null>(null);
	let trigger = $state<HTMLButtonElement | null>(null);

	async function signOut() {
		open = false;
		await backend.logout().catch(() => undefined);
		await goto('/login', { replaceState: true });
	}

	$effect(() => {
		if (!open) return;
		const outside = (event: MouseEvent) => {
			if (root && !root.contains(event.target as Node)) open = false;
		};
		const escape = (event: KeyboardEvent) => {
			if (event.key !== 'Escape') return;
			event.stopPropagation();
			open = false;
			trigger?.focus();
		};
		document.addEventListener('click', outside);
		document.addEventListener('keydown', escape);
		return () => {
			document.removeEventListener('click', outside);
			document.removeEventListener('keydown', escape);
		};
	});
</script>

<div bind:this={root} class="relative flex w-0 min-w-full shrink-0 flex-col px-3 pb-3">
	{#if user}
		{#if open}
			<div
				role="menu"
				aria-label={t('shell.account')}
				class="absolute bottom-full left-0 z-30 mb-2 flex flex-col gap-0.5 rounded-fc-lg bg-fc-page p-1.5 shadow-lg {collapsed
					? 'w-60'
					: 'right-0'}"
			>
				<div class="flex flex-col px-2.5 pt-1.5 pb-2">
					<span class="truncate text-fc-sm font-medium text-fc-fg">{user.name}</span>
					{#if user.email && user.email !== user.name}
						<span class="truncate text-fc-xs text-fc-fg-muted">{user.email}</span>
					{/if}
				</div>
				{#each entries as entry (entry.label)}
					<a
						href={entry.href}
						role="menuitem"
						onclick={() => (open = false)}
						class="flex min-h-10 items-center gap-2.5 whitespace-nowrap rounded-fc-md px-2.5 text-fc-sm text-fc-fg transition-colors hover:bg-fc-surface focus-visible:bg-fc-surface focus-visible:outline-none"
					>
						<iconify-icon icon={entry.icon} width="16" height="16" class="block shrink-0 text-fc-fg-muted"
						></iconify-icon>
						<span class="truncate">{entry.label}</span>
					</a>
				{/each}
				<div class="my-1 h-px bg-fc-surface"></div>
				<button
					type="button"
					role="menuitem"
					onclick={signOut}
					class="flex min-h-10 items-center gap-2.5 rounded-fc-md px-2.5 text-left text-fc-sm text-fc-danger transition-colors hover:bg-fc-surface focus-visible:bg-fc-surface focus-visible:outline-none"
				>
					<iconify-icon icon={icons.logout} width="16" height="16" class="block shrink-0"></iconify-icon>
					{t('shell.signOut')}
				</button>
			</div>
		{/if}
		<button
			bind:this={trigger}
			type="button"
			aria-haspopup="menu"
			aria-expanded={open}
			aria-label={t('shell.accountMenu', { name: user.name })}
			onclick={() => (open = !open)}
			use:springPress
			class={twMerge(
				'flex items-center gap-2.5 overflow-hidden rounded-fc-md text-left text-fc-sm text-fc-fg transition-colors hover:bg-fc-surface focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-fc-ring',
				collapsed ? 'size-fc-nav-item shrink-0 self-center justify-center gap-0 p-0' : 'min-h-11 w-full px-2.5 py-2',
				open && 'bg-fc-surface'
			)}
		>
			<Avatar src={user.avatar} name={user.name} alt="" size="sm" class="size-7 text-fc-xs" />
			{#if !collapsed}
				<span class="min-w-0 flex-1 truncate">{user.name}</span>
			{/if}
		</button>
	{/if}
</div>
