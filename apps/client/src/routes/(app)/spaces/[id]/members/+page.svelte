<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { backend, type SpaceMember, type Space } from '$lib/backend';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { toast } from 'svelte-sonner';
	import { ArrowLeft, Plus, Trash2, Crown, Shield, Users, X } from 'lucide-svelte';

	const spaceId = $derived(page.params.id as string);

	let space = $state<Space | null>(null);
	let members = $state<SpaceMember[]>([]);
	let loading = $state(true);
	let showAddForm = $state(false);
	let addUserId = $state('');
	let addRole = $state('member');
	let adding = $state(false);
	let removing = $state<string | null>(null);
	let updatingRole = $state<string | null>(null);

	async function loadData() {
		loading = true;
		try {
			const [spaceResult, membersResult] = await Promise.all([
				backend.getSpace(spaceId),
				backend.listSpaceMembers(spaceId)
			]);
			space = spaceResult;
			members = membersResult.members;
		} catch {
			space = null;
			members = [];
		}
		loading = false;
	}

	async function addMember() {
		const uid = parseInt(addUserId, 10);
		if (isNaN(uid)) {
			toast.error('Enter a valid user ID');
			return;
		}
		adding = true;
		try {
			await backend.addSpaceMember(spaceId, { user_id: uid, role: addRole });
			toast.success('Member added');
			addUserId = '';
			addRole = 'member';
			showAddForm = false;
			await loadData();
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to add member');
		}
		adding = false;
	}

	async function removeMember(member: SpaceMember) {
		removing = member.id;
		try {
			await backend.removeSpaceMember(spaceId, member.id);
			toast.success('Member removed');
			await loadData();
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to remove member');
		}
		removing = null;
	}

	async function changeRole(member: SpaceMember, newRole: string) {
		updatingRole = member.id;
		try {
			await backend.updateSpaceMember(spaceId, member.id, { role: newRole });
			toast.success('Role updated');
			await loadData();
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to update role');
		}
		updatingRole = null;
	}

	const canManage = $derived(space?.role === 'owner' || space?.role === 'admin');

	onMount(() => {
		loadData();
	});
</script>

<svelte:head>
	<title>Members — {space?.name ?? 'Space'} — Courrier</title>
</svelte:head>

<div class="flex flex-col gap-0 h-full">
	<div class="px-6 pt-6 pb-0">
		<button class="flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors mb-4" onclick={() => goto(`/spaces/${spaceId}`)}>
			<ArrowLeft class="h-4 w-4" />
			Back to Space
		</button>
		<h1 class="text-2xl font-semibold">Members</h1>
		<p class="text-sm text-muted-foreground mt-1">{space?.name ?? ''}</p>
	</div>

	<div class="flex-1 overflow-auto p-6">
		<div class="max-w-2xl">
			{#if loading}
				<p class="text-sm text-muted-foreground">Loading...</p>
			{:else}
				<div class="space-y-4">
					<div class="space-y-0">
						{#each members as member, i}
							<div class="flex items-center justify-between py-3 {i < members.length - 1 ? 'border-b' : ''}">
								<div class="min-w-0">
									<p class="font-medium text-sm">{member.name || member.email}</p>
									{#if member.name}
										<p class="text-xs text-muted-foreground">{member.email}</p>
									{/if}
								</div>
								<div class="flex items-center gap-2 shrink-0">
									{#if member.role === 'owner'}
										<span class="inline-flex items-center gap-1 text-xs text-muted-foreground">
											<Crown class="h-3 w-3" />
											Owner
										</span>
									{:else if canManage}
										<select
											class="h-8 rounded-md border border-input bg-transparent px-2 text-xs outline-none"
											value={member.role}
											onchange={(e) => changeRole(member, (e.target as HTMLSelectElement).value)}
											disabled={updatingRole === member.id}
										>
											<option value="admin">Admin</option>
											<option value="member">Member</option>
										</select>
										<Button
											variant="ghost"
											size="icon"
											class="h-8 w-8 text-muted-foreground hover:text-destructive"
											onclick={() => removeMember(member)}
											disabled={removing === member.id}
										>
											<Trash2 class="h-4 w-4" />
										</Button>
									{:else}
										<span class="inline-flex items-center gap-1 text-xs text-muted-foreground">
											{#if member.role === 'owner'}
												<Crown class="h-3 w-3" />
											{:else if member.role === 'admin'}
												<Shield class="h-3 w-3" />
											{:else}
												<Users class="h-3 w-3" />
											{/if}
											{member.role.charAt(0).toUpperCase() + member.role.slice(1)}
										</span>
									{/if}
								</div>
							</div>
						{/each}
					</div>

					{#if canManage}
						{#if showAddForm}
							<div class="space-y-4 pt-2 border-t">
								<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
									<div class="space-y-2">
										<Label for="add-user-id">User ID</Label>
										<Input id="add-user-id" bind:value={addUserId} placeholder="Enter user ID" />
									</div>
									<div class="space-y-2">
										<Label for="add-role">Role</Label>
										<select
											id="add-role"
											bind:value={addRole}
											class="h-9 w-full rounded-md border border-input bg-transparent px-3 text-sm outline-none"
										>
											<option value="member">Member</option>
											<option value="admin">Admin</option>
										</select>
									</div>
								</div>
								<div class="flex items-center gap-2">
									<Button size="sm" class="gap-1.5" onclick={addMember} disabled={adding || !addUserId}>
										<Plus class="h-4 w-4" />
										{adding ? 'Adding...' : 'Add Member'}
									</Button>
									<Button variant="ghost" size="sm" class="gap-1.5" onclick={() => (showAddForm = false)}>
										<X class="h-4 w-4" />
										Cancel
									</Button>
								</div>
							</div>
						{:else}
							<Button variant="outline" size="sm" class="gap-1.5" onclick={() => (showAddForm = true)}>
								<Plus class="h-4 w-4" />
								Add Member
							</Button>
						{/if}
					{/if}
				</div>
			{/if}
		</div>
	</div>
</div>
