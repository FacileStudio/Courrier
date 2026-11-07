<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import {
		Avatar,
		Badge,
		Button,
		Card,
		EmptyState,
		Skeleton,
		StatCard,
		icons,
		toast
	} from '@facile/muse';
	import { backend, type Space, type SpaceMember } from '$lib/backend';
	import { spaceStore } from '$lib/stores/space.svelte';
	import { formatDate } from '$lib/i18n/format';
	import { t } from '$lib/i18n/index.svelte';
	import { roleLabelKey, roleTone } from '../roles';

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
		toast.success(t('space.activeToast', { name: space.name }));
	}

	const isOwnerOrAdmin = $derived(
		!space?.manage_url && (space?.role === 'owner' || space?.role === 'admin')
	);
	const isActive = $derived(spaceStore.active?.id === spaceId);

	function roleName(role: string): string {
		const key = roleLabelKey(role);
		return key ? t(key) : role;
	}
</script>

<svelte:head>
	<title>{space?.name ?? t('space.fallbackName')} — Courrier</title>
</svelte:head>

<div class="w-full px-4 py-6 md:px-8">
	<div class="mx-auto flex max-w-4xl flex-col gap-10">
		{#if loading}
			<div class="flex flex-col gap-4">
				<Skeleton class="h-8 w-56" />
				<Skeleton class="h-28 w-full" />
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
					<h1 class="text-fc-2xl font-semibold text-fc-fg">{space.name}</h1>
					<p class="text-fc-sm text-fc-fg-muted">
						{space.description || t('spaces.noDescription')}
					</p>
				</div>
				<div class="flex flex-wrap items-center gap-2">
					<Button variant="ghost" href="/spaces" icon={icons.chevronLeft}>{t('space.spaces')}</Button>
					<Button icon={icons.check} disabled={isActive} onclick={switchToSpace}>
						{isActive ? t('space.active') : t('space.switch')}
					</Button>
					{#if space.manage_url}
						<Button variant="secondary" href={space.manage_url} icon={icons.space}>
							{t('spaces.manageInRegistre')}
						</Button>
					{:else if isOwnerOrAdmin}
						<Button variant="secondary" href="/spaces/{spaceId}/settings" icon={icons.settings}>
							{t('space.settings')}
						</Button>
					{/if}
				</div>
			</div>

			<section class="grid gap-4 sm:grid-cols-2">
				<StatCard label={t('space.members')} value={members.length} />
				<StatCard
					label={t('space.createdOn')}
					value={formatDate(space.created_at)}
					delta={t('space.yourRole', { role: roleName(space.role) })}
				/>
			</section>

			<section class="flex flex-col gap-4">
				<div class="flex flex-wrap items-start justify-between gap-4">
					<div class="flex min-w-0 flex-col gap-1">
						<h2 class="text-fc-lg font-semibold text-fc-fg">{t('space.members')}</h2>
						<p class="text-fc-sm text-fc-fg-muted">
							{t('space.membersIntro')}
						</p>
					</div>
					{#if isOwnerOrAdmin}
						<Button variant="secondary" href="/spaces/{spaceId}/members" icon={icons.usersGroup}>
							{t('space.manageMembers')}
						</Button>
					{/if}
				</div>

				{#if members.length === 0}
					<EmptyState
						icon={icons.usersGroup}
						title={t('space.nobody')}
						description={space.manage_url ? t('space.nobodyRegistre') : t('space.nobodyLocal')}
					/>
				{:else}
					<Card class="flex flex-col">
						{#each members as member (member.id)}
							<div
								class="flex items-center justify-between gap-3 py-4 first:pt-0 last:pb-0"
							>
								<div class="flex min-w-0 items-center gap-3">
									<Avatar src={member.avatar_url || undefined} name={member.name || member.email} size="sm" />
									<div class="min-w-0">
										<p class="truncate text-fc-sm font-medium text-fc-fg">{member.name || member.email}</p>
										{#if member.name}
											<p class="truncate text-fc-xs text-fc-fg-muted">{member.email}</p>
										{/if}
									</div>
								</div>
								<Badge tone={roleTone(member.role)}>{roleName(member.role)}</Badge>
							</div>
						{/each}
					</Card>
				{/if}
			</section>
		{/if}
	</div>
</div>
