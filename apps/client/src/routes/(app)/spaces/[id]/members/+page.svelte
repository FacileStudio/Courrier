<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import {
		Avatar,
		Badge,
		Button,
		ConfirmModal,
		EmptyState,
		Field,
		Input,
		Select,
		Skeleton,
		Table,
		icons,
		toast
	} from '@facile/muse';
	import { backend, type SpaceMember, type Space } from '$lib/backend';
	import { roleLabel, roleTone } from '../../roles';

	const spaceId = $derived(page.params.id as string);

	let space = $state<Space | null>(null);
	let members = $state<SpaceMember[]>([]);
	let loading = $state(true);
	let addUserId = $state('');
	let addRole = $state('member');
	let addError = $state('');
	let adding = $state(false);

	let pendingRemoval = $state<SpaceMember | null>(null);
	let confirmOpen = $state(false);

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

	async function refreshMembers() {
		const res = await backend.listSpaceMembers(spaceId);
		members = res.members ?? [];
	}

	async function addMember() {
		const uid = parseInt(addUserId, 10);
		if (isNaN(uid)) {
			addError = 'Entrez un ID utilisateur valide';
			return;
		}

		adding = true;
		addError = '';
		try {
			await backend.addSpaceMember(spaceId, { user_id: uid, role: addRole });
			addUserId = '';
			addRole = 'member';
			await refreshMembers();
			toast.success('Membre ajouté');
		} catch (err) {
			addError = err instanceof Error ? err.message : "Impossible d'ajouter le membre";
		}
		adding = false;
	}

	async function updateRole(memberId: string, role: string) {
		try {
			await backend.updateSpaceMember(spaceId, memberId, { role });
			await refreshMembers();
			toast.success('Rôle mis à jour');
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : 'Impossible de mettre à jour le rôle');
			await refreshMembers();
		}
	}

	function askRemove(member: SpaceMember) {
		pendingRemoval = member;
		confirmOpen = true;
	}

	async function removeMember() {
		const member = pendingRemoval;
		if (!member) return;
		try {
			await backend.removeSpaceMember(spaceId, member.id);
			members = members.filter((m) => m.id !== member.id);
			toast.success('Membre retiré');
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : 'Impossible de retirer le membre');
		}
		pendingRemoval = null;
	}

	const isOwnerOrAdmin = $derived(space?.role === 'owner' || space?.role === 'admin');
	const isOwner = $derived(space?.role === 'owner');
</script>

<svelte:head>
	<title>Membres — {space?.name ?? 'Espace'} — Courrier</title>
</svelte:head>

<div class="w-full px-4 py-6 md:px-8">
	<div class="mx-auto flex max-w-4xl flex-col gap-10">
		{#if loading}
			<div class="flex flex-col gap-4">
				<Skeleton class="h-8 w-56" />
				<Skeleton class="h-40 w-full" />
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
					<h1 class="text-fc-2xl font-semibold text-fc-fg">Membres</h1>
					<p class="text-fc-sm text-fc-fg-muted">{space.name}</p>
				</div>
				<Button variant="ghost" href="/spaces/{spaceId}" icon={icons.chevronLeft}>Espace</Button>
			</div>

			{#if isOwnerOrAdmin}
				<section class="flex flex-col gap-4">
					<div class="flex flex-col gap-1">
						<h2 class="text-fc-lg font-semibold text-fc-fg">Ajouter un membre</h2>
						<p class="text-fc-sm text-fc-fg-muted">
							Un membre voit les comptes mail rattachés à cet espace.
						</p>
					</div>
					<div class="flex flex-col gap-3 sm:flex-row sm:items-end">
						<div class="min-w-0 flex-1">
							<Field label="ID utilisateur" error={addError}>
								<Input bind:value={addUserId} placeholder="ex. 42" />
							</Field>
						</div>
						<Field label="Rôle">
							<Select bind:value={addRole} class="min-w-36">
								<option value="member">Membre</option>
								<option value="admin">Admin</option>
							</Select>
						</Field>
						<Button icon={icons.plus} size="lg" disabled={adding} onclick={addMember}>
							{adding ? 'Ajout…' : 'Ajouter'}
						</Button>
					</div>
				</section>
			{/if}

			<section class="flex flex-col gap-4">
				<h2 class="text-fc-lg font-semibold text-fc-fg">
					{members.length} membre{members.length !== 1 ? 's' : ''}
				</h2>

				{#if members.length === 0}
					<EmptyState
						icon={icons.usersGroup}
						title="Personne ici pour l'instant"
						description="Ajoutez un coéquipier au-dessus et il verra cet espace immédiatement."
					/>
				{:else}
					<Table>
						<thead>
							<tr>
								<th scope="col">Membre</th>
								<th scope="col">Rôle</th>
								<th scope="col"><span class="sr-only">Actions</span></th>
							</tr>
						</thead>
						<tbody>
							{#each members as member (member.id)}
								<tr>
									<td>
										<div class="flex min-w-0 items-center gap-3">
											<Avatar name={member.name || member.email} size="sm" />
											<div class="min-w-0">
												<p class="truncate font-medium text-fc-fg">{member.name || member.email}</p>
												{#if member.name}
													<p class="truncate text-fc-xs text-fc-fg-muted">{member.email}</p>
												{/if}
											</div>
										</div>
									</td>
									<td>
										{#if member.role === 'owner' || !isOwnerOrAdmin}
											<Badge tone={roleTone(member.role)}>{roleLabel(member.role)}</Badge>
										{:else}
											<Select
												value={member.role}
												class="min-w-32"
												aria-label="Rôle de {member.name || member.email}"
												onchange={(e) => updateRole(member.id, (e.currentTarget as HTMLSelectElement).value)}
											>
												<option value="member">Membre</option>
												<option value="admin">Admin</option>
												{#if isOwner}
													<option value="owner">Propriétaire</option>
												{/if}
											</Select>
										{/if}
									</td>
									<td>
										{#if member.role !== 'owner' && isOwnerOrAdmin}
											<Button
												variant="ghost-danger"
												size="sm"
												icon={icons.remove}
												onclick={() => askRemove(member)}
											>
												Retirer
											</Button>
										{/if}
									</td>
								</tr>
							{/each}
						</tbody>
					</Table>
				{/if}
			</section>
		{/if}
	</div>
</div>

<ConfirmModal
	bind:open={confirmOpen}
	tone="danger"
	title="Retirer {pendingRemoval?.name || pendingRemoval?.email || 'ce membre'} ?"
	description="Cette personne perd l'accès à l'espace et aux comptes mail qui y sont rattachés. Son propre compte Courrier et ses boîtes personnelles ne sont pas touchés."
	confirmLabel="Retirer le membre"
	onConfirm={removeMember}
	onCancel={() => (pendingRemoval = null)}
/>
