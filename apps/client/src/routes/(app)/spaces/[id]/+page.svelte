<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { backend, type Space } from '$lib/backend';
	import { spaceStore } from '$lib/stores/space.svelte';
	import { Button } from '$lib/components/ui/button';
	import { ArrowLeft, Users, Settings, Crown, Shield } from 'lucide-svelte';

	let space = $state<Space | null>(null);
	let loading = $state(true);

	const spaceId = $derived(page.params.id as string);

	function roleLabel(role: string) {
		return role.charAt(0).toUpperCase() + role.slice(1);
	}

	async function loadSpace() {
		loading = true;
		try {
			space = await backend.getSpace(spaceId);
		} catch {
			space = null;
		}
		loading = false;
	}

	function switchToSpace() {
		if (!space) return;
		spaceStore.set({ id: space.id, name: space.name, role: space.role });
	}

	onMount(() => {
		loadSpace();
	});
</script>

<svelte:head>
	<title>{space?.name ?? 'Space'} — Courrier</title>
</svelte:head>

<div class="flex flex-col gap-0 h-full">
	<div class="px-6 pt-6 pb-0">
		<button class="flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors mb-4" onclick={() => goto('/spaces')}>
			<ArrowLeft class="h-4 w-4" />
			Back to Spaces
		</button>

		{#if loading}
			<p class="text-sm text-muted-foreground">Loading...</p>
		{:else if !space}
			<p class="text-sm text-muted-foreground">Space not found.</p>
		{:else}
			<div class="flex items-center justify-between">
				<div>
					<h1 class="text-2xl font-semibold">{space.name}</h1>
					{#if space.description}
						<p class="text-sm text-muted-foreground mt-1">{space.description}</p>
					{/if}
				</div>
				<div class="flex items-center gap-2">
					{#if spaceStore.active?.id !== space.id}
						<Button variant="outline" size="sm" onclick={switchToSpace}>
							Switch to this space
						</Button>
					{:else}
						<span class="text-xs text-muted-foreground bg-muted px-2 py-1 rounded">Active</span>
					{/if}
					{#if space.role === 'owner' || space.role === 'admin'}
						<Button variant="ghost" size="icon" class="h-8 w-8" onclick={() => goto(`/spaces/${spaceId}/settings`)}>
							<Settings class="h-4 w-4" />
						</Button>
					{/if}
				</div>
			</div>
		{/if}
	</div>

	{#if space}
		<div class="flex-1 overflow-auto p-6">
			<div class="max-w-2xl">
				<div class="flex items-center justify-between mb-4">
					<h2 class="text-lg font-medium">Members</h2>
					{#if space.role === 'owner' || space.role === 'admin'}
						<Button variant="outline" size="sm" onclick={() => goto(`/spaces/${spaceId}/members`)}>
							Manage Members
						</Button>
					{/if}
				</div>

				{#if space.members && space.members.length > 0}
					<div class="space-y-0">
						{#each space.members as member, i}
							<div class="flex items-center justify-between py-3 {i < space.members.length - 1 ? 'border-b' : ''}">
								<div class="min-w-0">
									<p class="font-medium text-sm">{member.name || member.email}</p>
									{#if member.name}
										<p class="text-xs text-muted-foreground">{member.email}</p>
									{/if}
								</div>
								<span class="inline-flex items-center gap-1 text-xs text-muted-foreground shrink-0">
									{#if member.role === 'owner'}
										<Crown class="h-3 w-3" />
									{:else if member.role === 'admin'}
										<Shield class="h-3 w-3" />
									{:else}
										<Users class="h-3 w-3" />
									{/if}
									{roleLabel(member.role)}
								</span>
							</div>
						{/each}
					</div>
				{:else}
					<p class="text-sm text-muted-foreground">No members yet.</p>
				{/if}
			</div>
		</div>
	{/if}
</div>
