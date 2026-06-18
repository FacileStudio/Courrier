<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { backend, type Space } from '$lib/backend';
	import { spaceStore } from '$lib/stores/space.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { toast } from 'svelte-sonner';
	import { ArrowLeft, Save, Trash2 } from 'lucide-svelte';

	const spaceId = $derived(page.params.id as string);

	let space = $state<Space | null>(null);
	let loading = $state(true);
	let name = $state('');
	let description = $state('');
	let saving = $state(false);
	let deleting = $state(false);
	let confirmDelete = $state(false);

	async function loadSpace() {
		loading = true;
		try {
			space = await backend.getSpace(spaceId);
			name = space.name;
			description = space.description;
		} catch {
			space = null;
		}
		loading = false;
	}

	async function saveSpace() {
		saving = true;
		try {
			const updated = await backend.updateSpace(spaceId, {
				name: name.trim(),
				description: description.trim()
			});
			space = updated;
			if (spaceStore.active?.id === spaceId) {
				spaceStore.set({ id: updated.id, name: updated.name, role: updated.role });
			}
			toast.success('Space updated');
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to update space');
		}
		saving = false;
	}

	async function deleteSpace() {
		deleting = true;
		try {
			await backend.deleteSpace(spaceId);
			if (spaceStore.active?.id === spaceId) {
				spaceStore.clear();
			}
			toast.success('Space deleted');
			goto('/spaces');
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to delete space');
		}
		deleting = false;
	}

	onMount(() => {
		loadSpace();
	});
</script>

<svelte:head>
	<title>Settings — {space?.name ?? 'Space'} — Courrier</title>
</svelte:head>

<div class="flex flex-col gap-0 h-full">
	<div class="px-6 pt-6 pb-0">
		<button class="flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors mb-4" onclick={() => goto(`/spaces/${spaceId}`)}>
			<ArrowLeft class="h-4 w-4" />
			Back to Space
		</button>
		<h1 class="text-2xl font-semibold">Space Settings</h1>
		<p class="text-sm text-muted-foreground mt-1">{space?.name ?? ''}</p>
	</div>

	<div class="flex-1 overflow-auto p-6">
		<div class="max-w-lg">
			{#if loading}
				<p class="text-sm text-muted-foreground">Loading...</p>
			{:else if !space}
				<p class="text-sm text-muted-foreground">Space not found.</p>
			{:else}
				<div class="space-y-8">
					<div class="space-y-6">
						<div class="space-y-2">
							<Label for="space-name">Name</Label>
							<Input id="space-name" bind:value={name} />
						</div>

						<div class="space-y-2">
							<Label for="space-desc">Description</Label>
							<textarea
								id="space-desc"
								bind:value={description}
								class="h-24 w-full resize-none rounded-md border border-input bg-transparent px-3 py-2 text-sm leading-relaxed outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 placeholder:text-muted-foreground"
							></textarea>
						</div>

						<Button class="gap-1.5" onclick={saveSpace} disabled={saving || !name.trim()}>
							<Save class="h-4 w-4" />
							{saving ? 'Saving...' : 'Save Changes'}
						</Button>
					</div>

					{#if space.role === 'owner'}
						<div class="border-t pt-6">
							<h2 class="text-sm font-medium text-destructive">Danger Zone</h2>
							<p class="text-xs text-muted-foreground mt-1 mb-4">Permanently delete this space and remove all members.</p>

							{#if confirmDelete}
								<div class="flex items-center gap-2">
									<Button
										variant="destructive"
										size="sm"
										class="gap-1.5"
										onclick={deleteSpace}
										disabled={deleting}
									>
										<Trash2 class="h-4 w-4" />
										{deleting ? 'Deleting...' : 'Yes, delete this space'}
									</Button>
									<Button variant="ghost" size="sm" onclick={() => (confirmDelete = false)}>
										Cancel
									</Button>
								</div>
							{:else}
								<Button
									variant="outline"
									size="sm"
									class="gap-1.5 text-destructive hover:text-destructive"
									onclick={() => (confirmDelete = true)}
								>
									<Trash2 class="h-4 w-4" />
									Delete Space
								</Button>
							{/if}
						</div>
					{/if}
				</div>
			{/if}
		</div>
	</div>
</div>
