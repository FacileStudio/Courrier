<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { backend, type Space, type SpaceMember } from '$lib/backend';
	import { spaceStore } from '$lib/stores/space.svelte';

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
		if (space) {
			spaceStore.set({ id: space.id, name: space.name, role: space.role });
		}
	}

	function roleBadgeClass(role: string): string {
		switch (role) {
			case 'owner': return 'bg-amber-500/10 text-amber-600';
			case 'admin': return 'bg-blue-500/10 text-blue-600';
			default: return 'bg-muted text-muted-foreground';
		}
	}

	let isOwnerOrAdmin = $derived(space?.role === 'owner' || space?.role === 'admin');
</script>

<svelte:head>
	<title>{space?.name ?? 'Espace'} — Courrier</title>
</svelte:head>

<div class="flex h-full flex-col">
	<div class="border-b border-border px-4 py-4 md:px-8 md:py-5">
		<div class="flex items-center gap-3">
			<a href="/spaces" class="text-muted-foreground transition-colors hover:text-foreground" aria-label="Retour aux espaces">
				<iconify-icon icon="solar:arrow-left-linear" width="20"></iconify-icon>
			</a>
			<div class="flex-1">
				{#if space}
					<h1 class="text-lg font-semibold">{space.name}</h1>
					{#if space.description}
						<p class="mt-0.5 text-sm text-muted-foreground">{space.description}</p>
					{/if}
				{:else}
					<h1 class="text-lg font-semibold">Espace</h1>
				{/if}
			</div>
			{#if space}
				<div class="flex items-center gap-2">
					<button
						class="inline-flex h-9 items-center gap-2 rounded-md bg-primary px-4 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
						onclick={switchToSpace}
					>
						<iconify-icon icon="solar:login-3-linear" width="16"></iconify-icon>
						Basculer vers cet espace
					</button>
					{#if isOwnerOrAdmin}
						<a
							href="/spaces/{spaceId}/settings"
							class="inline-flex h-9 items-center gap-2 rounded-md border border-border px-4 text-sm font-medium transition-colors hover:bg-accent"
						>
							<iconify-icon icon="solar:settings-linear" width="16"></iconify-icon>
							Paramètres
						</a>
					{/if}
				</div>
			{/if}
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
				<a href="/spaces" class="mt-3 text-sm text-primary hover:underline">Retour aux espaces</a>
			</div>
		{:else}
			<div class="max-w-2xl space-y-8">
				<div>
					<div class="flex items-center justify-between">
						<h2 class="text-sm font-semibold uppercase tracking-wider text-muted-foreground">Membres</h2>
						{#if isOwnerOrAdmin}
							<a
								href="/spaces/{spaceId}/members"
								class="inline-flex items-center gap-1.5 text-sm text-primary transition-colors hover:underline"
							>
								<iconify-icon icon="solar:user-plus-linear" width="16"></iconify-icon>
								Gérer les membres
							</a>
						{/if}
					</div>

					<div class="mt-4 space-y-2">
						{#each members as member}
							<div class="flex items-center justify-between rounded-lg border border-border p-3">
								<div class="flex items-center gap-3">
									<div
										class="flex h-8 w-8 items-center justify-center rounded-full border border-border bg-muted text-xs font-semibold"
									>
										{(member.name || member.email || '?').slice(0, 2).toUpperCase()}
									</div>
									<div>
										<p class="text-sm font-medium">{member.name || member.email}</p>
										{#if member.name}
											<p class="text-xs text-muted-foreground">{member.email}</p>
										{/if}
									</div>
								</div>
								<span class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium {roleBadgeClass(member.role)}">
									{member.role}
								</span>
							</div>
						{/each}
					</div>
				</div>

				<div class="rounded-lg border border-border/60 bg-muted/20 p-4">
					<div class="flex items-center gap-2 text-sm text-muted-foreground">
						<iconify-icon icon="solar:info-circle-linear" width="16"></iconify-icon>
						<span>Créé le {new Date(space.created_at).toLocaleDateString('fr-FR')}</span>
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>
