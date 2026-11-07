<script lang="ts">
	import { getContext } from 'svelte';
	import { Avatar, icons, twMerge } from '@facile/muse';
	import type { MailAccount } from '$lib/backend';
	import { t } from '$lib/i18n/index.svelte';

	let {
		align = 'start',
		rail = false,
		collapsed = false
	}: { align?: 'start' | 'end'; rail?: boolean; collapsed?: boolean } = $props();

	const menuPosition = $derived(
		rail ? 'bottom-full left-0 mb-1.5' : align === 'end' ? 'top-full right-0 mt-1.5' : 'top-full left-0 mt-1.5'
	);
	const triggerShape = $derived(
		rail
			? collapsed
				? 'size-11 justify-center rounded-fc-md'
				: 'min-h-11 w-full rounded-fc-md px-2.5'
			: 'min-h-10 max-w-64 rounded-fc-pill px-2'
	);

	const app = getContext<{
		activeAccountId: number | null;
		accounts: MailAccount[];
		setActiveAccount: (id: number) => Promise<void>;
	}>('app');

	let open = $state(false);
	let rootEl: HTMLElement | null = $state(null);
	let triggerEl: HTMLButtonElement | null = $state(null);
	let menuEl: HTMLElement | null = $state(null);

	const current = $derived(
		app.accounts.find((a) => a.id === app.activeAccountId) ?? app.accounts[0] ?? null
	);

	function items(): HTMLElement[] {
		return menuEl ? Array.from(menuEl.querySelectorAll<HTMLElement>('[role="menuitemradio"], a')) : [];
	}

	function choose(id: number) {
		open = false;
		triggerEl?.focus();
		void app.setActiveAccount(id);
	}

	function handleClickOutside(e: MouseEvent) {
		if (rootEl && !rootEl.contains(e.target as Node)) open = false;
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') {
			e.stopPropagation();
			open = false;
			triggerEl?.focus();
			return;
		}
		if (e.key !== 'ArrowDown' && e.key !== 'ArrowUp') return;
		const list = items();
		if (list.length === 0) return;
		e.preventDefault();
		const index = list.indexOf(document.activeElement as HTMLElement);
		const step = e.key === 'ArrowDown' ? 1 : -1;
		list[(index + step + list.length) % list.length].focus();
	}

	$effect(() => {
		if (!open) return;
		items()[0]?.focus();
		document.addEventListener('click', handleClickOutside);
		document.addEventListener('keydown', handleKeydown);
		return () => {
			document.removeEventListener('click', handleClickOutside);
			document.removeEventListener('keydown', handleKeydown);
		};
	});
</script>

{#if app.accounts.length > 1 && current}
	<div bind:this={rootEl} class={twMerge('relative min-w-0', rail && 'w-full shrink-0', rail && collapsed && 'w-auto self-center')}>
		<button
			bind:this={triggerEl}
			type="button"
			class={twMerge(
				'flex items-center gap-2 text-fc-sm font-medium text-fc-fg transition-colors hover:bg-fc-surface focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-fc-ring',
				triggerShape
			)}
			aria-haspopup="menu"
			aria-expanded={open}
			aria-label={t('shell.switchAccount', { email: current.email })}
			onclick={() => (open = !open)}
		>
			<Avatar name={current.name} size="sm" class="size-6" />
			{#if !collapsed}
				<span class={twMerge('min-w-0 truncate', rail ? 'flex-1 text-left' : 'hidden sm:inline')}>{current.email}</span>
			{/if}
			{#if !collapsed}<iconify-icon
				icon={icons.chevronDown}
				width="16"
				height="16"
				class={twMerge('block size-4 shrink-0 text-fc-fg-muted transition-transform', open && 'rotate-180')}
			></iconify-icon>{/if}
		</button>

		{#if open}
			<div
				bind:this={menuEl}
				role="menu"
				class="absolute z-50 {menuPosition} flex w-64 flex-col overflow-hidden rounded-fc-md bg-fc-component p-1 shadow-lg"
			>
				{#each app.accounts as account (account.id)}
					{@const active = account.id === current.id}
					<button
						type="button"
						role="menuitemradio"
						aria-checked={active}
						class={twMerge(
							'flex min-h-11 items-center gap-2.5 rounded-fc-sm px-2.5 py-2 text-left text-fc-sm transition-colors focus-visible:outline-2 focus-visible:outline-fc-ring',
							active ? 'bg-fc-accent text-fc-accent-fg' : 'text-fc-fg hover:bg-fc-surface'
						)}
						onclick={() => choose(account.id)}
					>
						<span class="flex min-w-0 flex-1 flex-col">
							<span class="truncate font-medium">{account.name}</span>
							<span class="truncate text-fc-xs opacity-70">{account.email}</span>
						</span>
						{#if active}
							<iconify-icon icon={icons.check} width="16" height="16" class="block size-4 shrink-0"
							></iconify-icon>
						{/if}
					</button>
				{/each}
				<a
					href="/settings/accounts"
					class="mt-1 flex min-h-11 items-center rounded-fc-sm px-2.5 py-2 text-fc-sm text-fc-fg-muted transition-colors hover:bg-fc-surface hover:text-fc-fg focus-visible:outline-2 focus-visible:outline-fc-ring"
					onclick={() => (open = false)}
				>
					{t('shell.manageAccounts')}
				</a>
			</div>
		{/if}
	</div>
{/if}
