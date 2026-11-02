<script lang="ts">
	import { icons, springPress, twMerge } from '@facile/muse';

	let {
		user,
		href,
		active = false,
		collapsed = false
	}: {
		user?: { name: string; avatar?: string };
		href: string;
		active?: boolean;
		collapsed?: boolean;
	} = $props();
</script>

<div class="flex w-0 min-w-full shrink-0 flex-col gap-1 rounded-fc-lg bg-fc-component p-1.5">
	{#if user}
		<a
			{href}
			aria-current={active ? 'page' : undefined}
			aria-label="{user.name}, settings"
			use:springPress
			class={twMerge(
				'flex items-center gap-2.5 overflow-hidden rounded-fc-md text-fc-sm text-fc-fg transition-colors hover:bg-fc-surface focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-fc-ring',
				collapsed
					? 'size-fc-nav-item shrink-0 self-center justify-center gap-0 p-0'
					: 'min-h-11 w-full px-2.5 py-2',
				active && 'bg-fc-surface'
			)}
		>
			{#if user.avatar}
				<img src={user.avatar} alt="" class="size-7 shrink-0 rounded-fc-pill object-cover" />
			{:else}
				<span
					class="flex size-7 shrink-0 items-center justify-center rounded-fc-pill bg-fc-accent text-fc-xs font-semibold text-fc-accent-fg"
				>
					{user.name.charAt(0).toUpperCase()}
				</span>
			{/if}
			{#if !collapsed}
				<span class="flex min-w-0 flex-1 items-center justify-between gap-2">
					<span class="truncate">{user.name}</span>
					<iconify-icon icon={icons.settings} width="18" height="18" class="block shrink-0"></iconify-icon>
				</span>
			{/if}
		</a>
	{/if}
</div>
