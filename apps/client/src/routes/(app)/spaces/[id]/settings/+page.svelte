<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import {
		Alert,
		Button,
		ConfirmModal,
		EmptyState,
		Field,
		Input,
		SettingsSection,
		Skeleton,
		Textarea,
		icons,
		toast
	} from '@facile/muse';
	import { backend, type Space } from '$lib/backend';
	import { spaceStore } from '$lib/stores/space.svelte';

	const spaceId = $derived(page.params.id as string);

	let space = $state<Space | null>(null);
	let loading = $state(true);
	let name = $state('');
	let description = $state('');
	let saving = $state(false);
	let nameError = $state('');
	let confirmDelete = $state(false);
	let deleteError = $state('');

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
		if (!name.trim()) {
			nameError = 'Le nom est requis';
			return;
		}

		saving = true;
		nameError = '';
		try {
			space = await backend.updateSpace(spaceId, {
				name: name.trim(),
				description: description.trim()
			});
			if (spaceStore.active?.id === spaceId) {
				spaceStore.set({ id: space.id, name: space.name, role: space.role });
			}
			toast.success('Paramètres enregistrés');
		} catch (err) {
			nameError = err instanceof Error ? err.message : "Impossible d'enregistrer";
		}
		saving = false;
	}

	async function deleteSpace() {
		deleteError = '';
		try {
			await backend.deleteSpace(spaceId);
		} catch (err) {
			deleteError = err instanceof Error ? err.message : "Impossible de supprimer l'espace";
			throw err;
		}
		if (spaceStore.active?.id === spaceId) {
			spaceStore.clear();
		}
		await goto('/spaces');
		toast.success('Espace supprimé');
	}

	const isOwner = $derived(space?.role === 'owner');
</script>

<svelte:head>
	<title>Paramètres — {space?.name ?? 'Espace'} — Courrier</title>
</svelte:head>

<div class="h-full overflow-auto px-4 py-6 md:px-8">
	<div class="mx-auto flex max-w-2xl flex-col gap-10">
		{#if loading}
			<div class="flex flex-col gap-4">
				<Skeleton class="h-8 w-56" />
				<Skeleton class="h-56 w-full" />
			</div>
		{:else if !space}
			<EmptyState
				icon={icons.warning}
				title="Espace introuvable"
				description="Il a peut-être été supprimé, ou votre accès a été retiré."
			>
				<Button variant="outline" href="/spaces" icon={icons.chevronLeft}>Retour aux espaces</Button>
			</EmptyState>
		{:else}
			<div class="flex flex-wrap items-start justify-between gap-4">
				<div class="flex min-w-0 flex-col gap-1">
					<h1 class="text-fc-2xl font-semibold text-fc-fg">Paramètres de l'espace</h1>
					<p class="text-fc-sm text-fc-fg-muted">{space.name}</p>
				</div>
				<Button variant="ghost" href="/spaces/{spaceId}" icon={icons.chevronLeft}>Espace</Button>
			</div>

			<SettingsSection title="Général" description="Ce que les membres voient de cet espace.">
				<Field label="Nom" error={nameError}>
					<Input bind:value={name} placeholder="ex. Équipe Marketing" />
				</Field>
				<Field label="Description" helper="Optionnel — à quoi sert cet espace ?">
					<Textarea bind:value={description} rows={3} />
				</Field>
				<div>
					<Button icon={icons.check} disabled={saving} onclick={saveSettings}>
						{saving ? 'Enregistrement…' : 'Enregistrer'}
					</Button>
				</div>
			</SettingsSection>

			{#if isOwner}
				<SettingsSection
					title="Zone de danger"
					description="Supprimer cet espace est immédiat et définitif."
				>
					<p class="text-fc-sm text-fc-fg-muted">
						Tous les membres perdent leur accès, et les comptes mail rattachés à cet espace
						disparaissent de la liste : ils ne sont ni supprimés ni rendus personnels, ils ne sont
						simplement plus rattachés à rien. Les identifiants IMAP/SMTP restent en base.
					</p>
					<div>
						<Button variant="danger" icon={icons.remove} onclick={() => (confirmDelete = true)}>
							Supprimer l'espace
						</Button>
					</div>
				</SettingsSection>
			{/if}
		{/if}
	</div>
</div>

<ConfirmModal
	bind:open={confirmDelete}
	tone="danger"
	title="Supprimer « {space?.name ?? 'cet espace'} » ?"
	description="Tous les membres perdent leur accès, et les comptes mail rattachés à cet espace quittent la liste — ils ne réapparaissent pas dans vos comptes personnels. Cette action est définitive."
	confirmLabel="Supprimer l'espace"
	onConfirm={deleteSpace}
	onCancel={() => (deleteError = '')}
>
	{#if deleteError}
		<Alert tone="danger">{deleteError}</Alert>
	{/if}
</ConfirmModal>
