<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { MAIL_FOLDERS } from '$lib/mail-folders';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { ChevronDown } from 'lucide-svelte';

	let { folders = [] }: { folders?: { type: string; unread_count?: number }[] } = $props();

	const current = $derived(MAIL_FOLDERS.find((f) => f.href === page.url.pathname) ?? MAIL_FOLDERS[0]);

	function unread(type: string): number {
		return folders.find((f) => f.type === type)?.unread_count ?? 0;
	}
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger
		class="flex items-center gap-1.5 text-lg font-semibold outline-none"
	>
		<iconify-icon icon={current.icon} width="18"></iconify-icon>
		<span>{current.label}</span>
		<ChevronDown class="h-4 w-4 text-muted-foreground" />
	</DropdownMenu.Trigger>
	<DropdownMenu.Content align="start" class="w-48">
		{#each MAIL_FOLDERS as folder}
			{@const count = unread(folder.type)}
			<DropdownMenu.Item onclick={() => goto(folder.href)}>
				<iconify-icon icon={folder.icon} width="16"></iconify-icon>
				<span class="flex-1">{folder.label}</span>
				{#if count > 0}
					<span class="text-xs text-muted-foreground">{count}</span>
				{/if}
			</DropdownMenu.Item>
		{/each}
	</DropdownMenu.Content>
</DropdownMenu.Root>
