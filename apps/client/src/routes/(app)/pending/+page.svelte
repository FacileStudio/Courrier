<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { Alert, Button, EmptyState, Skeleton, icons, toast } from '@facile/muse';
	import { rememberReturn } from '$lib/after-login';
	import { ApiError, backend, type PendingAction } from '$lib/backend';
	import PendingCard from '$lib/components/PendingCard.svelte';
	import { t } from '$lib/i18n/index.svelte';
	import { isApp } from '$lib/platform';

	let actions = $state<PendingAction[]>([]);
	let loading = $state(true);
	let busy = $state('');
	let notOwnSession = $state(false);
	const native = isApp();

	const focused = $derived(page.url.searchParams.get('action') ?? '');
	const ordered = $derived([
		...actions.filter((action) => action.id === focused),
		...actions.filter((action) => action.id !== focused)
	]);

	onMount(load);

	/** Loads the open actions. */
	async function load() {
		loading = true;
		try {
			actions = (await backend.listPending()).actions;
		} catch (err) {
			if (!refusedSession(err)) toast.danger(t('pending.loadFailed'));
		}
		loading = false;
	}

	/** Runs one action once, then drops it from the list. */
	async function confirm(action: PendingAction) {
		busy = action.id;
		try {
			await backend.confirmPending(action.id);
			toast.success(action.kind === 'send' ? t('pending.sent') : t('pending.deleted'));
			actions = actions.filter((item) => item.id !== action.id);
		} catch (err) {
			if (!refusedSession(err)) {
				toast.danger(err instanceof Error ? err.message : t('pending.failed'));
				await load();
			}
		}
		busy = '';
	}

	/** Closes one action without running it. */
	async function dismiss(action: PendingAction) {
		busy = action.id;
		try {
			await backend.dismissPending(action.id);
			toast.success(t('pending.dismissed'));
			actions = actions.filter((item) => item.id !== action.id);
		} catch (err) {
			if (!refusedSession(err)) toast.danger(err instanceof Error ? err.message : t('pending.failed'));
		}
		busy = '';
	}

	/** A 403 means this is not the person's own browser session, which alone may confirm. */
	function refusedSession(err: unknown): boolean {
		if (!(err instanceof ApiError) || err.status !== 403) return false;
		notOwnSession = true;
		return true;
	}

	/** Ends the old session and comes back here after a fresh browser sign-in. */
	async function signInAgain() {
		rememberReturn(page.url);
		await backend.logout().catch(() => undefined);
		await goto('/login');
	}
</script>

<svelte:head>
	<title>{t('pending.title')}</title>
</svelte:head>

<div class="w-full px-4 py-6 md:px-8">
	<div class="mx-auto flex max-w-2xl flex-col gap-8">
		<div class="flex flex-col gap-1">
			<h1 class="text-fc-2xl font-semibold text-fc-fg">{t('pending.title')}</h1>
			<p class="text-fc-sm text-fc-fg-muted">{t('pending.intro')}</p>
		</div>

		{#if notOwnSession}
			<Alert tone="warning" title={t('pending.signInTitle')}>
				<div class="flex flex-col gap-3">
					<p>{native ? t('pending.appCannotConfirm') : t('pending.signInBody')}</p>
					{#if !native}
						<div><Button variant="secondary" class="min-h-11" onclick={signInAgain}>{t('pending.signInAgain')}</Button></div>
					{/if}
				</div>
			</Alert>
		{:else if loading}
			<Skeleton class="h-48 w-full" />
		{:else if actions.length === 0}
			<EmptyState icon={icons.check} title={t('pending.none')} description={t('pending.noneBody')} />
		{:else}
			<section class="flex flex-col gap-4">
				{#each ordered as action (action.id)}
					<PendingCard {action} busy={busy === action.id} onconfirm={confirm} ondismiss={dismiss} />
				{/each}
			</section>
		{/if}
	</div>
</div>
