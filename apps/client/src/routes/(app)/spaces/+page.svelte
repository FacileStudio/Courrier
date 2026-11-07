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
	import { roleLabelKey, roleTone } from './roles';
	import { t } from '$lib/i18n/index.svelte';

	let spaces = $state<Space[]>([]);
	let manageUrl = $state('');
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
			manageUrl = res.manage_url ?? '';
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
			createError = t('spaces.nameRequired');
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
			toast.success(t('spaces.created', { name: space.name }));
			await goto(`/spaces/${space.id}`);
		} catch (err) {
			createError = err instanceof Error ? err.message : t('spaces.createFailed');
		}
		creating = false;
	}
</script>

<svelte:head>
	<title>{t('spaces.title')}</title>
</svelte:head>

<div class="w-full px-4 py-6 md:px-8">
	<div class="mx-auto flex max-w-4xl flex-col gap-10">
		<div class="flex flex-wrap items-start justify-between gap-4">
			<div class="flex min-w-0 flex-col gap-1">
				<h1 class="text-fc-2xl font-semibold text-fc-fg">{t('spaces.heading')}</h1>
				<p class="text-fc-sm text-fc-fg-muted">
					{t('spaces.intro')}
				</p>
			</div>
			<div class="flex items-center gap-2">
				<Button variant="ghost" href="/settings" icon={icons.chevronLeft}>{t('spaces.settings')}</Button>
				{#if manageUrl}
					<Button href={manageUrl} icon={icons.space}>{t('spaces.manageInRegistre')}</Button>
				{:else}
					<Button icon={icons.plus} onclick={openCreate}>{t('spaces.new')}</Button>
				{/if}
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
				icon={icons.space}
				title={t('spaces.none')}
				description={manageUrl ? t('spaces.noneRegistre') : t('spaces.noneLocal')}
			>
				{#if manageUrl}
					<Button href={manageUrl} icon={icons.space}>{t('spaces.manageInRegistre')}</Button>
				{:else}
					<Button icon={icons.plus} onclick={openCreate}>{t('spaces.createFirst')}</Button>
				{/if}
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
										icon={icons.space}
										width="20"
										height="20"
										class="block size-5"
									></iconify-icon>
								</span>
								<p class="truncate text-fc-md font-medium text-fc-fg">{space.name}</p>
							</div>
							<Badge tone={roleTone(space.role)}>{roleLabelKey(space.role) ? t(roleLabelKey(space.role)!) : space.role}</Badge>
						</div>
						<p class="line-clamp-2 text-fc-sm text-fc-fg-muted">
							{space.description || t('spaces.noDescription')}
						</p>
					</Card>
				{/each}
			</section>
		{/if}
	</div>
</div>

<Modal bind:open={createOpen} title={t('spaces.createTitle')} showClose dismissible={!creating}>
	<div class="flex flex-col gap-4">
		<Field label={t('spaces.nameLabel')} error={createError}>
			<Input bind:value={draftName} placeholder={t('spaces.namePlaceholder')} />
		</Field>
		<Field label={t('spaces.descriptionLabel')} helper={t('spaces.descriptionHelp')}>
			<Textarea bind:value={draftDescription} rows={3} placeholder={t('spaces.descriptionPlaceholder')} />
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
				{t('spaces.cancel')}
			</Button>
			<Button icon={icons.plus} class="w-full sm:w-auto" disabled={creating} onclick={createSpace}>
				{creating ? t('spaces.creating') : t('spaces.createSubmit')}
			</Button>
		</div>
	{/snippet}
</Modal>
