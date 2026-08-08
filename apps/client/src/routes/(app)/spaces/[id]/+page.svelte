<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import {
		Badge,
		Button,
		Card,
		EmptyState,
		SettingsRow,
		Skeleton,
		StatCard,
		icons,
		toast
	} from '@facile/muse';
	import { backend, type Space, type SpaceMember } from '$lib/backend';
	import { spaceStore } from '$lib/stores/space.svelte';
	import { roleLabel, roleTone } from '../roles';

	let space = $state<Space | null>(null);
	let members = $state<SpaceMember[]>([]);
	let loading = $state(true);

	const spaceId = $derived(page.params.id as string);

	onMount(async () => {
		await loadData();
	});

	async function loadData() {
		loading = true;
		try {
			const [spaceRes, membersRes] = await Promise.all([
				backend.getSpace(spaceId),
				backend.listSpaceMembers(spaceId)
			]);
			space = spaceRes;
			members = membersRes.members ?? [];
		} catch {
			space = null;
			members = [];
		}
		loading = false;
	}

	function switchToSpace() {
		if (!space) return;
		spaceStore.set({ id: space.id, name: space.name, role: space.role });
		toast.success(`Espace actif : ${space.name}`);
	}

	const isOwnerOrAdmin = $derived(space?.role === 'owner' || space?.role === 'admin');
	const isActive = $derived(spaceStore.active?.id === spaceId);
</script>

<svelte:head>
	<title>{space?.name ?? 'Espace'} — Courrier</title>
</svelte:head>

<div class="h-full overflow-auto px-4 py-6 md:px-8">
	<div class="mx-auto flex max-w-4xl flex-col gap-10">
		{#if loading}
			<div class="flex flex-col gap-4">
				<Skeleton class="h-8 w-56" />
				<Skeleton class="h-28 w-full" />
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
					<h1 class="text-fc-2xl font-semibold text-fc-fg">{space.name}</h1>
					<p class="text-fc-sm text-fc-fg-muted">
						{space.description || 'Pas de description.'}
					</p>
				</div>
				<div class="flex flex-wrap items-center gap-2">
					<Button variant="ghost" href="/spaces" icon={icons.chevronLeft}>Espaces</Button>
					<Button icon={icons.check} disabled={isActive} onclick={switchToSpace}>
						{isActive ? 'Espace actif' : 'Basculer vers cet espace'}
					</Button>
					{#if isOwnerOrAdmin}
						<Button variant="outline" href="/spaces/{spaceId}/settings" icon={icons.settings}>
							Paramètres
						</Button>
					{/if}
				</div>
			</div>

			<section class="grid gap-4 sm:grid-cols-2">
				<StatCard label="Membres" value={members.length} />
				<StatCard
					label="Créé le"
					value={new Date(space.created_at).toLocaleDateString('fr-FR')}
					delta="Votre rôle : {roleLabel(space.role)}"
				/>
			</section>

			<section class="flex flex-col gap-4">
				<div class="flex flex-wrap items-start justify-between gap-4">
					<div class="flex min-w-0 flex-col gap-1">
						<h2 class="text-fc-lg font-semibold text-fc-fg">Membres</h2>
						<p class="text-fc-sm text-fc-fg-muted">
							Tout le monde ici voit les comptes mail rattachés à cet espace.
						</p>
					</div>
					{#if isOwnerOrAdmin}
						<Button variant="outline" href="/spaces/{spaceId}/members" icon={icons.usersGroup}>
							Gérer les membres
						</Button>
					{/if}
				</div>

				{#if members.length === 0}
					<EmptyState
						icon={icons.usersGroup}
						title="Personne ici pour l'instant"
						description="Ajoutez un coéquipier et il verra cet espace immédiatement."
					/>
				{:else}
					<Card class="flex flex-col">
						{#each members as member (member.id)}
							<SettingsRow
								label={member.name || member.email}
								description={member.name ? member.email : undefined}
							>
								<Badge tone={roleTone(member.role)}>{roleLabel(member.role)}</Badge>
							</SettingsRow>
						{/each}
					</Card>
				{/if}
			</section>
		{/if}
	</div>
</div>
