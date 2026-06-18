<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { backend, type Space } from '$lib/backend';
	import { spaceStore } from '$lib/stores/space.svelte';
	import { toast } from 'svelte-sonner';

	const spaceId = $derived(page.params.id as string);

	let space = $state<Space | null>(null);
	let loading = $state(true);
	let name = $state('');
	let description = $state('');
	let saving = $state(false);
	let message = $state('');
	let confirmDelete = $state(false);
	let deleting = $state(false);

	onMount(async () => {
		await loadSpace();
	});

	async function loadSpace() {
		loading = true;
		try {
			space = await backend.getSpace(spaceId);
			name = space.name;
			description = space.description ?? '';
		} catch {
			space = null;
		}
		loading = false;
	}

	async function saveSettings() {
		saving = true;
		message = '';
		try {
			space = await backend.updateSpace(spaceId, {
				name: name.trim(),
				description: description.trim()
			});
			message = 'Paramètres enregistrés';
			if (spaceStore.active?.id === spaceId) {
				spaceStore.set({ id: space.id, name: space.name, role: space.role });
			}
		} catch (e: any) {
			message = e.message || 'Impossible d\'enregistrer';
		}
		saving = false;
		setTimeout(() => { message = ''; }, 3000);
	}

	async function deleteSpace() {
		deleting = true;
		try {
			await backend.deleteSpace(spaceId);
			if (spaceStore.active?.id === spaceId) {
				spaceStore.clear();
			}
			toast.success('Espace supprimé');
			goto('/spaces');
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Impossible de supprimer l\'espace');
		}
		deleting = false;
		confirmDelete = false;
	}

	let isOwner = $derived(space?.role === 'owner');
</script>

<svelte:head>
	<title>Paramètres — {space?.name ?? 'Espace'} — Courrier</title>
</svelte:head>

<div class="flex h-full flex-col">
	<div class="border-b border-border px-4 py-4 md:px-8 md:py-5">
		<div class="flex items-center gap-3">
			<a href="/spaces/{spaceId}" class="text-muted-foreground transition-colors hover:text-foreground" aria-label="Retour à l'espace">
				<iconify-icon icon="solar:arrow-left-linear" width="20"></iconify-icon>
			</a>
			<div>
				<h1 class="text-lg font-semibold">Paramètres de l'espace</h1>
				{#if space}
					<p class="mt-0.5 text-sm text-muted-foreground">{space.name}</p>
				{/if}
			</div>
		</div>
	</div>

	<div class="flex-1 overflow-auto px-4 py-6 md:px-8">
		{#if loading}
			<div class="flex items-center justify-center py-20">
				<div class="h-6 w-6 animate-spin rounded-full border-2 border-foreground border-t-transparent"></div>
			</div>
		{:else if !space}
			<div class="flex flex-col items-center justify-center py-20 text-center">
				<p class="text-sm text-muted-foreground">Espace introuvable ou accès refusé.</p>
			</div>
		{:else}
			<div class="max-w-xl space-y-8">
				<div class="space-y-4">
					<h2 class="text-sm font-semibold uppercase tracking-wider text-muted-foreground">Général</h2>
					<div>
						<label for="space-name" class="mb-1.5 block text-sm font-medium">Nom</label>
						<input
							id="space-name"
							type="text"
							bind:value={name}
							class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
						/>
					</div>
					<div>
						<label for="space-desc" class="mb-1.5 block text-sm font-medium">Description</label>
						<textarea
							id="space-desc"
							bind:value={description}
							rows="3"
							class="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring resize-none"
						></textarea>
					</div>
					<div class="flex items-center gap-3">
						<button
							class="inline-flex h-9 items-center rounded-md bg-primary px-4 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90 disabled:opacity-50"
							onclick={saveSettings}
							disabled={saving}
						>
							{saving ? 'Enregistrement...' : 'Enregistrer'}
						</button>
						{#if message}
							<span class="text-sm text-muted-foreground">{message}</span>
						{/if}
					</div>
				</div>

				{#if isOwner}
					<div class="space-y-4 border-t border-border pt-8">
						<h2 class="text-sm font-semibold uppercase tracking-wider text-destructive">Zone de danger</h2>
						<p class="text-sm text-muted-foreground">
							Supprimer un espace retire tous les membres. Les données associées deviendront non-assignées.
						</p>
						{#if confirmDelete}
							<div class="flex items-center gap-3">
								<button
									class="inline-flex h-9 items-center gap-2 rounded-md border border-destructive/30 bg-destructive/10 px-4 text-sm font-medium text-destructive transition-colors hover:bg-destructive/20 disabled:opacity-50"
									onclick={deleteSpace}
									disabled={deleting}
								>
									<iconify-icon icon="solar:trash-bin-2-linear" width="16"></iconify-icon>
									{deleting ? 'Suppression...' : 'Confirmer la suppression'}
								</button>
								<button
									class="inline-flex h-9 items-center rounded-md border border-border px-4 text-sm font-medium transition-colors hover:bg-accent"
									onclick={() => confirmDelete = false}
								>
									Annuler
								</button>
							</div>
						{:else}
							<button
								class="inline-flex h-9 items-center gap-2 rounded-md border border-destructive/30 bg-destructive/10 px-4 text-sm font-medium text-destructive transition-colors hover:bg-destructive/20 disabled:opacity-50"
								onclick={() => confirmDelete = true}
								disabled={deleting}
							>
								<iconify-icon icon="solar:trash-bin-2-linear" width="16"></iconify-icon>
								Supprimer l'espace
							</button>
						{/if}
					</div>
				{/if}
			</div>
		{/if}
	</div>
</div>
