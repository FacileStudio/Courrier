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
	import { t } from '$lib/i18n/index.svelte';

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
			nameError = t('spaces.nameRequired');
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
			toast.success(t('spaceSettings.saved'));
		} catch (err) {
			nameError = err instanceof Error ? err.message : t('spaceSettings.saveFailed');
		}
		saving = false;
	}

	async function deleteSpace() {
		deleteError = '';
		try {
			await backend.deleteSpace(spaceId);
		} catch (err) {
			deleteError = err instanceof Error ? err.message : t('spaceSettings.deleteFailed');
			throw err;
		}
		if (spaceStore.active?.id === spaceId) {
			spaceStore.clear();
		}
		await goto('/spaces');
		toast.success(t('spaceSettings.deleted'));
	}

	const isOwner = $derived(space?.role === 'owner');
</script>

<svelte:head>
	<title>{t('spaceSettings.title', { name: space?.name ?? t('space.fallbackName') })}</title>
</svelte:head>

<div class="w-full px-4 py-6 md:px-8">
	<div class="mx-auto flex max-w-2xl flex-col gap-10">
		{#if loading}
			<div class="flex flex-col gap-4">
				<Skeleton class="h-8 w-56" />
				<Skeleton class="h-56 w-full" />
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
					<h1 class="text-fc-2xl font-semibold text-fc-fg">{t('spaceSettings.heading')}</h1>
					<p class="text-fc-sm text-fc-fg-muted">{space.name}</p>
				</div>
				<Button variant="ghost" href="/spaces/{spaceId}" icon={icons.chevronLeft}>{t('spaceSettings.space')}</Button>
			</div>

			{#if space.manage_url}
				<SettingsSection
					title={t('spaceSettings.managedTitle')}
					description={t('spaceSettings.managedBody')}
				>
					<div>
						<Button href={space.manage_url} icon={icons.space}>{t('spaces.manageInRegistre')}</Button>
					</div>
				</SettingsSection>
			{:else}
				<SettingsSection title={t('spaceSettings.general')} description={t('spaceSettings.generalBody')}>
					<Field label={t('spaces.nameLabel')} error={nameError}>
						<Input bind:value={name} placeholder={t('spaces.namePlaceholder')} />
					</Field>
					<Field label={t('spaces.descriptionLabel')} helper={t('spaces.descriptionHelp')}>
						<Textarea bind:value={description} rows={3} />
					</Field>
					<div>
						<Button icon={icons.check} disabled={saving} onclick={saveSettings}>
							{saving ? t('spaceSettings.saving') : t('spaceSettings.save')}
						</Button>
					</div>
				</SettingsSection>

				{#if isOwner}
					<SettingsSection
						title={t('spaceSettings.danger')}
						description={t('spaceSettings.dangerBody')}
					>
						<p class="text-fc-sm text-fc-fg-muted">
							{t('spaceSettings.dangerDetail')}
						</p>
						<div>
							<Button variant="danger" icon={icons.remove} onclick={() => (confirmDelete = true)}>
								{t('spaceSettings.delete')}
							</Button>
						</div>
					</SettingsSection>
				{/if}
			{/if}
		{/if}
	</div>
</div>

<ConfirmModal
	bind:open={confirmDelete}
	tone="danger"
	title={t('spaceSettings.deleteTitle', { name: space?.name ?? t('spaceSettings.deleteThis') })}
	description={t('spaceSettings.deleteBody')}
	confirmLabel={t('spaceSettings.delete')}
	onConfirm={deleteSpace}
	onCancel={() => (deleteError = '')}
>
	{#if deleteError}
		<Alert tone="danger">{deleteError}</Alert>
	{/if}
</ConfirmModal>
