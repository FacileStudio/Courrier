<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import {
		Badge,
		Button,
		Card,
		EmptyState,
		Field,
		Input,
		Modal,
		Skeleton,
		Textarea,
		icons,
		toast
	} from '@facile/muse';
	import { backend, type Space } from '$lib/backend';
	import { spaceStore } from '$lib/stores/space.svelte';
	import { roleLabel, roleTone } from './roles';

	let spaces = $state<Space[]>([]);
	let loading = $state(true);

	let createOpen = $state(false);
	let draftName = $state('');
	let draftDescription = $state('');
	let creating = $state(false);
	let createError = $state('');

	onMount(async () => {
		await loadSpaces();
	});

	async function loadSpaces() {
		loading = true;
		try {
			const res = await backend.listSpaces();
			spaces = res.spaces ?? [];
		} catch {
			spaces = [];
		}
		loading = false;
	}

	function openCreate() {
		draftName = '';
		draftDescription = '';
		createError = '';
		createOpen = true;
	}

	async function createSpace() {
		if (!draftName.trim()) {
			createError = 'Le nom est requis';
			return;
		}

		creating = true;
		createError = '';
		try {
			const space = await backend.createSpace({
				name: draftName.trim(),
				description: draftDescription.trim()
			});
			spaceStore.set({ id: space.id, name: space.name, role: space.role });
			createOpen = false;
			toast.success(`Espace « ${space.name} » créé`);
			await goto(`/spaces/${space.id}`);
		} catch (err) {
			createError = err instanceof Error ? err.message : "Impossible de créer l'espace";
		}
		creating = false;
	}
</script>

<svelte:head>
	<title>Espaces — Courrier</title>
</svelte:head>

<div class="w-full px-4 py-6 md:px-8">
	<div class="mx-auto flex max-w-4xl flex-col gap-10">
		<div class="flex flex-wrap items-start justify-between gap-4">
			<div class="flex min-w-0 flex-col gap-1">
				<h1 class="text-fc-2xl font-semibold text-fc-fg">Espaces</h1>
				<p class="text-fc-sm text-fc-fg-muted">
					Collaborez avec votre équipe dans des espaces partagés : chaque espace a ses membres et
					ses comptes mail.
				</p>
			</div>
			<div class="flex items-center gap-2">
				<Button variant="ghost" href="/settings" icon={icons.chevronLeft}>Réglages</Button>
				<Button icon={icons.plus} onclick={openCreate}>Nouvel espace</Button>
			</div>
		</div>

		{#if loading}
			<div class="grid gap-4 sm:grid-cols-2">
				{#each [0, 1, 2, 3] as row (row)}
					<Skeleton class="h-28 w-full" />
				{/each}
			</div>
		{:else if spaces.length === 0}
			<EmptyState
				icon={icons.usersGroup}
				title="Aucun espace"
				description="Créez un espace pour partager des comptes mail et des dossiers avec votre équipe."
			>
				<Button icon={icons.plus} onclick={openCreate}>Créer votre premier espace</Button>
			</EmptyState>
		{:else}
			<section class="grid gap-4 sm:grid-cols-2">
				{#each spaces as space (space.id)}
					<Card href="/spaces/{space.id}" class="flex flex-col gap-4">
						<div class="flex items-start justify-between gap-3">
							<div class="flex min-w-0 items-center gap-3">
								<span
									class="flex size-10 shrink-0 items-center justify-center rounded-fc-md bg-fc-surface text-fc-fg-muted transition-colors group-hover:bg-fc-accent group-hover:text-fc-accent-fg"
								>
									<iconify-icon
										icon={icons.usersGroup}
										width="20"
										height="20"
										class="block size-5"
									></iconify-icon>
								</span>
								<p class="truncate text-fc-md font-medium text-fc-fg">{space.name}</p>
							</div>
							<Badge tone={roleTone(space.role)}>{roleLabel(space.role)}</Badge>
						</div>
						<p class="line-clamp-2 text-fc-sm text-fc-fg-muted">
							{space.description || 'Pas de description.'}
						</p>
					</Card>
				{/each}
			</section>
		{/if}
	</div>
</div>

<Modal bind:open={createOpen} title="Nouvel espace" showClose dismissible={!creating}>
	<div class="flex flex-col gap-4">
		<Field label="Nom" error={createError}>
			<Input bind:value={draftName} placeholder="ex. Équipe Marketing" />
		</Field>
		<Field label="Description" helper="Optionnel — à quoi sert cet espace ?">
			<Textarea bind:value={draftDescription} rows={3} placeholder="Boîtes partagées de l'équipe." />
		</Field>
	</div>
	{#snippet footer()}
		<div class="flex flex-col-reverse gap-2 sm:flex-row sm:justify-end">
			<Button
				variant="ghost"
				class="w-full sm:w-auto"
				disabled={creating}
				onclick={() => (createOpen = false)}
			>
				Annuler
			</Button>
			<Button icon={icons.plus} class="w-full sm:w-auto" disabled={creating} onclick={createSpace}>
				{creating ? 'Création…' : "Créer l'espace"}
			</Button>
		</div>
	{/snippet}
</Modal>
