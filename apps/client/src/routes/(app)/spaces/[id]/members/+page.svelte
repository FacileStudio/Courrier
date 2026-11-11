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
	import { backend, type SpaceMember, type Space, type User } from '$lib/backend';
	import { t } from '$lib/i18n/index.svelte';
	import { roleLabelKey, roleTone } from '../../roles';

	const spaceId = $derived(page.params.id as string);

	let space = $state<Space | null>(null);
	let members = $state<SpaceMember[]>([]);
	let allUsers = $state<User[]>([]);
	let loading = $state(true);
	let selectedUserId = $state('');
	let addEmail = $state('');
	let addRole = $state('member');
	let addError = $state('');
	let adding = $state(false);

	let pendingRemoval = $state<SpaceMember | null>(null);
	let confirmOpen = $state(false);

	const availableUsers = $derived(
		allUsers.filter((u) => !members.some((m) => String(m.user_id) === String(u.id)))
	);

	onMount(async () => {
		await loadData();
	});

	async function loadData() {
		loading = true;
		try {
			const [spaceRes, membersRes, usersRes] = await Promise.all([
				backend.getSpace(spaceId),
				backend.listSpaceMembers(spaceId),
				backend.listUsers()
			]);
			space = spaceRes;
			members = membersRes.members ?? [];
			allUsers = usersRes.users ?? [];
		} catch {
			space = null;
			members = [];
			allUsers = [];
		}
		loading = false;
	}

	async function refreshMembers() {
		const res = await backend.listSpaceMembers(spaceId);
		members = res.members ?? [];
	}

	async function addMember() {
		const email = addEmail.trim();
		const uid = parseInt(selectedUserId, 10);
		if (!email && isNaN(uid)) {
			addError = t('members.pickUserError');
			return;
		}

		adding = true;
		addError = '';
		try {
			await backend.addSpaceMember(spaceId, email ? { email, role: addRole } : { user_id: uid, role: addRole });
			selectedUserId = '';
			addEmail = '';
			addRole = 'member';
			await refreshMembers();
			toast.success(t('members.added'));
		} catch (err) {
			addError = err instanceof Error ? err.message : t('members.addFailed');
		}
		adding = false;
	}

	async function updateRole(memberId: string, role: string) {
		try {
			await backend.updateSpaceMember(spaceId, memberId, { role });
			await refreshMembers();
			toast.success(t('members.roleUpdated'));
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : t('members.roleFailed'));
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
			toast.success(t('members.removed'));
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : t('members.removeFailed'));
		}
		pendingRemoval = null;
	}

	const isOwnerOrAdmin = $derived(
		!space?.manage_url && (space?.role === 'owner' || space?.role === 'admin')
	);
	const isOwner = $derived(space?.role === 'owner');

	function roleName(role: string): string {
		const key = roleLabelKey(role);
		return key ? t(key) : role;
	}

	const userOptions = $derived(
		availableUsers.map((user) => ({ value: String(user.id), label: user.name ? `${user.name} (${user.email})` : user.email }))
	);
	const addRoleOptions = $derived([
		{ value: 'member', label: t('spaces.roleMember') },
		{ value: 'admin', label: t('spaces.roleAdmin') }
	]);
</script>

<svelte:head>
	<title>{t('members.title', { name: space?.name ?? t('space.fallbackName') })}</title>
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
				icon={icons.space}
				title={t('space.notFound')}
				description={t('space.notFoundBody')}
			>
				<Button variant="secondary" href="/spaces" icon={icons.chevronLeft}>{t('space.backToSpaces')}</Button>
			</EmptyState>
		{:else}
			<div class="flex flex-wrap items-start justify-between gap-4">
				<div class="flex min-w-0 flex-col gap-1">
					<h1 class="text-fc-2xl font-semibold text-fc-fg">{t('members.heading')}</h1>
					<p class="text-fc-sm text-fc-fg-muted">{space.name}</p>
				</div>
				<div class="flex flex-wrap items-center gap-2">
					<Button variant="ghost" href="/spaces/{spaceId}" icon={icons.chevronLeft}>{t('spaceSettings.space')}</Button>
					{#if space.manage_url}
						<Button href={space.manage_url} icon={icons.space}>{t('spaces.manageInRegistre')}</Button>
					{/if}
				</div>
			</div>

			{#if isOwnerOrAdmin}
				<section class="flex flex-col gap-4">
					<div class="flex flex-col gap-1">
						<h2 class="text-fc-lg font-semibold text-fc-fg">{t('members.addTitle')}</h2>
						<p class="text-fc-sm text-fc-fg-muted">
							{t('members.addBody')}
						</p>
					</div>
					<div class="flex flex-col gap-3 sm:flex-row sm:items-end">
						<div class="min-w-0 flex-1">
							<Field label={t('members.email')} error={addError}>
								<Input
									type="email"
									bind:value={addEmail}
									class="w-full"
									autocomplete="off"
									placeholder={t('members.emailPlaceholder')}
								/>
							</Field>
						</div>
						{#if userOptions.length > 0}
							<div class="min-w-0 flex-1">
								<Field label={t('members.user')}>
									<Select
										bind:value={selectedUserId}
										class="w-full"
										placeholder={t('members.pickUser')}
										options={userOptions}
										disabled={addEmail.trim() !== ''}
									/>
								</Field>
							</div>
						{/if}
						<Field label={t('members.role')}>
							<Select bind:value={addRole} class="min-w-36" options={addRoleOptions} />
						</Field>
						<Button icon={icons.plus} size="lg" disabled={adding || (!selectedUserId && !addEmail.trim())} onclick={addMember}>
							{adding ? t('members.adding') : t('members.add')}
						</Button>
					</div>
				</section>
			{/if}

			<section class="flex flex-col gap-4">
				<h2 class="text-fc-lg font-semibold text-fc-fg">
					{t('members.count', { count: members.length })}
				</h2>

				{#if members.length === 0}
					<EmptyState
						icon={icons.usersGroup}
						title={t('space.nobody')}
						description={t('members.nobodyBody')}
					/>
				{:else}
					<Table>
						<thead>
							<tr>
								<th scope="col">{t('members.colMember')}</th>
								<th scope="col">{t('members.colRole')}</th>
								<th scope="col"><span class="sr-only">{t('members.colActions')}</span></th>
							</tr>
						</thead>
						<tbody>
							{#each members as member (member.id)}
								<tr>
									<td>
										<div class="flex min-w-0 items-center gap-3">
											<Avatar src={member.avatar_url || undefined} name={member.name || member.email} size="sm" />
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
											<Badge tone={roleTone(member.role)}>{roleName(member.role)}</Badge>
										{:else}
											<Select
												value={member.role}
												class="min-w-32"
												aria-label={t('members.roleOf', { name: member.name || member.email })}
												options={[
													...addRoleOptions,
													...(isOwner ? [{ value: 'owner', label: t('spaces.roleOwner') }] : [])
												]}
												onchange={(e) => updateRole(member.id, (e.currentTarget as HTMLSelectElement).value)}
											/>
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
												{t('members.remove')}
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
	title={t('members.removeTitle', { name: pendingRemoval?.name || pendingRemoval?.email || t('members.removeThis') })}
	description={t('members.removeBody')}
	confirmLabel={t('members.removeConfirm')}
	onConfirm={removeMember}
	onCancel={() => (pendingRemoval = null)}
/>
