<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { backend, type Space } from '$lib/backend';
	import { spaceStore } from '$lib/stores/space.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Plus, Users, Settings, LogOut, Crown, Shield } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';

	let spaces = $state<Space[]>([]);
	let loading = $state(true);
	let leaving = $state<string | null>(null);

	function roleLabel(role: string) {
		return role.charAt(0).toUpperCase() + role.slice(1);
	}

	async function loadSpaces() {
		loading = true;
		try {
			const result = await backend.listSpaces();
			spaces = result.spaces;
		} catch {
			spaces = [];
		}
		loading = false;
	}

	async function leaveSpace(space: Space) {
		leaving = space.id;
		try {
			await backend.leaveSpace(space.id);
			toast.success(`Left ${space.name}`);
			if (spaceStore.active?.id === space.id) {
				spaceStore.clear();
			}
			await loadSpaces();
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to leave space');
		}
		leaving = null;
	}

	onMount(() => {
		loadSpaces();
	});
</script>

<svelte:head>
	<title>Spaces — Courrier</title>
</svelte:head>

<div class="flex flex-col gap-0 h-full">
	<div class="px-6 pt-6 pb-0">
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-2xl font-semibold">Spaces</h1>
				<p class="text-sm text-muted-foreground mt-1">Collaborate with your team in shared spaces.</p>
			</div>
			<Button class="gap-1.5" onclick={() => goto('/spaces/new')}>
				<Plus class="h-4 w-4" />
				New Space
			</Button>
		</div>
	</div>

	<div class="flex-1 overflow-auto p-6">
		<div class="max-w-2xl">
			{#if loading}
				<p class="text-sm text-muted-foreground">Loading...</p>
			{:else if spaces.length === 0}
				<div class="flex flex-col items-center justify-center gap-4 py-16">
					<Users class="h-12 w-12 text-muted-foreground/40" />
					<div class="text-center">
						<p class="text-sm font-medium">No spaces yet</p>
						<p class="text-sm text-muted-foreground mt-1">Create a space to start collaborating.</p>
					</div>
					<Button variant="outline" class="gap-1.5" onclick={() => goto('/spaces/new')}>
						<Plus class="h-4 w-4" />
						Create your first space
					</Button>
				</div>
			{:else}
				<div class="space-y-0">
					{#each spaces as space, i}
						<div class="flex items-center justify-between py-3 {i < spaces.length - 1 ? 'border-b' : ''}">
							<a href="/spaces/{space.id}" class="min-w-0 flex-1 group">
								<p class="font-medium text-sm group-hover:underline">{space.name}</p>
								{#if space.description}
									<p class="text-xs text-muted-foreground mt-0.5">{space.description}</p>
								{/if}
							</a>
							<div class="flex items-center gap-2 shrink-0 ml-4">
								<span class="inline-flex items-center gap-1 text-xs text-muted-foreground">
									{#if space.role === 'owner'}
										<Crown class="h-3 w-3" />
									{:else if space.role === 'admin'}
										<Shield class="h-3 w-3" />
									{:else}
										<Users class="h-3 w-3" />
									{/if}
									{roleLabel(space.role)}
								</span>
								{#if space.role === 'owner' || space.role === 'admin'}
									<Button
										variant="ghost"
										size="icon"
										class="h-8 w-8"
										onclick={() => goto(`/spaces/${space.id}/settings`)}
									>
										<Settings class="h-4 w-4" />
									</Button>
								{/if}
								{#if space.role !== 'owner'}
									<Button
										variant="ghost"
										size="icon"
										class="h-8 w-8 text-muted-foreground hover:text-destructive"
										onclick={() => leaveSpace(space)}
										disabled={leaving === space.id}
									>
										<LogOut class="h-4 w-4" />
									</Button>
								{/if}
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</div>
	</div>
</div>
