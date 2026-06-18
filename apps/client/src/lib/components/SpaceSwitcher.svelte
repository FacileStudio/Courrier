<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { backend, type Space } from '$lib/backend';
	import { spaceStore } from '$lib/stores/space.svelte';
	import { Button } from '$lib/components/ui/button';
	import { ChevronDown, Plus, Users, Check } from 'lucide-svelte';

	let spaces = $state<Space[]>([]);
	let open = $state(false);

	async function loadSpaces() {
		try {
			const result = await backend.listSpaces();
			spaces = result.spaces;
		} catch {
			spaces = [];
		}
	}

	function selectSpace(space: Space) {
		spaceStore.set({ id: space.id, name: space.name, role: space.role });
		open = false;
	}

	function clearSpace() {
		spaceStore.clear();
		open = false;
	}

	function handleClickOutside(e: MouseEvent) {
		const target = e.target as HTMLElement;
		if (!target.closest('.space-switcher')) {
			open = false;
		}
	}

	onMount(() => {
		loadSpaces();
		document.addEventListener('click', handleClickOutside);
		return () => document.removeEventListener('click', handleClickOutside);
	});
</script>

<div class="space-switcher relative">
	<button
		class="flex w-full items-center gap-2 rounded-md border border-border/60 bg-muted/30 px-3 py-2 text-sm transition-colors hover:bg-muted"
		onclick={() => (open = !open)}
	>
		<Users class="h-4 w-4 shrink-0 text-muted-foreground" />
		<span class="flex-1 truncate text-left">
			{spaceStore.active?.name ?? 'Personnel'}
		</span>
		<ChevronDown class="h-3.5 w-3.5 shrink-0 text-muted-foreground transition-transform {open ? 'rotate-180' : ''}" />
	</button>

	{#if open}
		<div class="absolute left-0 right-0 top-full z-50 mt-1 rounded-md border bg-popover p-1 shadow-md">
			<button
				class="flex w-full items-center gap-2 rounded-sm px-2 py-1.5 text-sm transition-colors hover:bg-accent"
				onclick={clearSpace}
			>
				<span class="flex-1 text-left">Personnel</span>
				{#if !spaceStore.active}
					<Check class="h-3.5 w-3.5" />
				{/if}
			</button>

			{#each spaces as space}
				<button
					class="flex w-full items-center gap-2 rounded-sm px-2 py-1.5 text-sm transition-colors hover:bg-accent"
					onclick={() => selectSpace(space)}
				>
					<span class="flex-1 truncate text-left">{space.name}</span>
					{#if spaceStore.active?.id === space.id}
						<Check class="h-3.5 w-3.5" />
					{/if}
				</button>
			{/each}

			<div class="my-1 border-t"></div>

			<button
				class="flex w-full items-center gap-2 rounded-sm px-2 py-1.5 text-sm text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
				onclick={() => { open = false; goto('/spaces/new'); }}
			>
				<Plus class="h-3.5 w-3.5" />
				Nouvel espace
			</button>

			<button
				class="flex w-full items-center gap-2 rounded-sm px-2 py-1.5 text-sm text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
				onclick={() => { open = false; goto('/spaces'); }}
			>
				<Users class="h-3.5 w-3.5" />
				Gérer les espaces
			</button>
		</div>
	{/if}
</div>
