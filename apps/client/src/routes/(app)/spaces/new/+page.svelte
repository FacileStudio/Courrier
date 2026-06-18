<script lang="ts">
	import { goto } from '$app/navigation';
	import { backend } from '$lib/backend';
	import { spaceStore } from '$lib/stores/space.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { toast } from 'svelte-sonner';
	import { ArrowLeft } from 'lucide-svelte';

	let name = $state('');
	let description = $state('');
	let saving = $state(false);

	async function createSpace() {
		if (!name.trim()) return;
		saving = true;
		try {
			const space = await backend.createSpace({ name: name.trim(), description: description.trim() });
			spaceStore.set({ id: space.id, name: space.name, role: space.role });
			toast.success('Space created');
			goto(`/spaces/${space.id}`);
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to create space');
		}
		saving = false;
	}
</script>

<svelte:head>
	<title>New Space — Courrier</title>
</svelte:head>

<div class="flex flex-col gap-0 h-full">
	<div class="px-6 pt-6 pb-0">
		<button class="flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors mb-4" onclick={() => goto('/spaces')}>
			<ArrowLeft class="h-4 w-4" />
			Back to Spaces
		</button>
		<h1 class="text-2xl font-semibold">Create Space</h1>
		<p class="text-sm text-muted-foreground mt-1">Set up a new shared workspace for your team.</p>
	</div>

	<div class="flex-1 overflow-auto p-6">
		<div class="max-w-lg space-y-6">
			<div class="space-y-2">
				<Label for="space-name">Name</Label>
				<Input id="space-name" bind:value={name} placeholder="e.g. Marketing Team" />
			</div>

			<div class="space-y-2">
				<Label for="space-desc">Description</Label>
				<textarea
					id="space-desc"
					bind:value={description}
					placeholder="What is this space for?"
					class="h-24 w-full resize-none rounded-md border border-input bg-transparent px-3 py-2 text-sm leading-relaxed outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 placeholder:text-muted-foreground"
				></textarea>
			</div>

			<div class="flex items-center gap-2">
				<Button onclick={createSpace} disabled={saving || !name.trim()}>
					{saving ? 'Creating...' : 'Create Space'}
				</Button>
				<Button variant="ghost" onclick={() => goto('/spaces')}>
					Cancel
				</Button>
			</div>
		</div>
	</div>
</div>
