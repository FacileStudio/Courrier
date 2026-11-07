<script lang="ts">
	import { onMount, tick } from 'svelte';
	import { goto } from '$app/navigation';
	import { Alert, Button, PageTransition, Spinner, icons } from '@facile/muse';
	import { backend, type MailAccount } from '$lib/backend';
	import AppMark from '$lib/components/native/AppMark.svelte';
	import Ambient from '$lib/components/native/Ambient.svelte';
	import NotificationControl from '$lib/components/NotificationControl.svelte';
	import { STEPS, clampStep, isLastStep, markOnboardingSeen, primaryLabel } from '$lib/onboarding/steps';
	import { t } from '$lib/i18n/index.svelte';
	import MailboxStep from './MailboxStep.svelte';
	import StepProgress from './StepProgress.svelte';

	let { name }: { name: string } = $props();

	const COPY = $derived({
		welcome: { title: t('onboarding.welcomeTitle', { name }), body: t('onboarding.welcomeBody') },
		mailbox: { title: t('onboarding.mailboxTitle'), body: t('onboarding.mailboxBody') },
		notifications: { title: t('onboarding.notificationsTitle'), body: t('onboarding.notificationsBody') },
		done: { title: t('onboarding.doneTitle'), body: t('onboarding.doneBody') }
	});

	let step = $state(0);
	let accounts = $state<MailAccount[]>([]);
	let loading = $state(true);
	let choice = $state('');
	let drawerOpen = $state(false);
	let saving = $state(false);
	let error = $state('');
	let content = $state<HTMLElement | null>(null);

	const current = $derived(STEPS[step]);
	const copy = $derived(COPY[current.id]);

	async function loadAccounts() {
		try {
			accounts = (await backend.listAccounts()).accounts;
			choice = String((accounts.find((a) => a.is_default) ?? accounts[0])?.id ?? '');
		} catch {
			accounts = [];
		} finally {
			loading = false;
		}
	}

	onMount(loadAccounts);

	$effect(() => {
		void step;
		void tick().then(() => content?.querySelector<HTMLElement>('[data-step-heading]')?.focus({ preventScroll: true }));
	});

	async function saveDefault() {
		const picked = accounts.find((a) => String(a.id) === choice);
		if (!picked || picked.is_default || accounts.length < 2) return;
		await backend.updateAccount(picked.id, { is_default: true });
	}

	async function finish() {
		markOnboardingSeen();
		await goto('/mail', { replaceState: true });
	}

	async function onContinue() {
		error = '';
		if (current.id === 'mailbox' && accounts.length === 0) {
			drawerOpen = true;
			return;
		}
		if (isLastStep(step)) return finish();
		saving = true;
		try {
			if (current.id === 'mailbox') await saveDefault();
			step = clampStep(step + 1);
		} catch (err) {
			error = err instanceof Error ? err.message : t('onboarding.saveFailed');
		} finally {
			saving = false;
		}
	}
</script>

<div class="fixed inset-0 z-50 flex h-dvh flex-col overflow-hidden bg-fc-page text-fc-fg">
	<Ambient />

	<header class="relative z-10 flex flex-col items-center gap-4 px-4 pt-4 sm:px-8 sm:pt-6">
		<div class="flex w-full max-w-fc-sm items-center justify-between">
			<div class="flex items-center gap-2 text-fc-fg-muted">
				<AppMark size={20} />
				<span class="text-fc-sm font-medium">Courrier</span>
			</div>
			{#if !isLastStep(step)}
				<Button variant="ghost" size="sm" class="min-h-11 px-3" onclick={finish}>{t('onboarding.skipSetup')}</Button>
			{/if}
		</div>
		<StepProgress index={step} />
	</header>

	<main class="relative z-10 min-h-0 flex-1 overflow-y-auto overscroll-contain">
		<PageTransition key={step} distance={16} duration={0.3} class="min-h-full">
			<div
				bind:this={content}
				class="mx-auto flex min-h-full w-full max-w-fc-sm flex-col items-center justify-center gap-6 px-4 py-8 text-center"
			>
				<div
					class="flex size-24 items-center justify-center rounded-[1.75rem] bg-fc-component shadow-sm"
				>
					{#if current.id === 'notifications'}
						<iconify-icon icon={icons.notification} width="40" height="40" aria-hidden="true"></iconify-icon>
					{:else if current.id === 'done'}
						<iconify-icon icon={icons.check} width="40" height="40" aria-hidden="true"></iconify-icon>
					{:else}
						<AppMark size={44} />
					{/if}
				</div>
				<div class="flex flex-col items-center gap-2">
					<h1
						data-step-heading
						tabindex="-1"
						class="font-fc-title text-fc-2xl font-semibold tracking-tight outline-none sm:text-fc-3xl"
					>
						{copy.title}
					</h1>
					<p class="text-fc-md text-fc-fg-muted">{copy.body}</p>
				</div>

				{#if current.id === 'mailbox'}
					<MailboxStep {accounts} {loading} bind:choice bind:drawerOpen onSaved={loadAccounts} />
				{:else if current.id === 'notifications'}
					<div class="flex w-full justify-center"><NotificationControl /></div>
				{/if}

				{#if error}
					<div class="w-full" role="alert"><Alert tone="danger">{error}</Alert></div>
				{/if}
			</div>
		</PageTransition>
	</main>

	<footer
		class="relative z-10 bg-fc-page px-4 pt-3 pb-[max(0.75rem,env(safe-area-inset-bottom))] sm:px-8"
	>
		<div class="mx-auto flex w-full max-w-fc-sm flex-col gap-1">
			<Button
				size="lg"
				class="h-12 w-full rounded-fc-md"
				iconRight={isLastStep(step) ? undefined : icons.arrow}
				disabled={saving || (current.id === 'mailbox' && loading)}
				onclick={onContinue}
			>
				{#if saving}<Spinner size="sm" />{/if}
				{t(primaryLabel(step, accounts.length))}
			</Button>
			<div class="flex min-h-11 items-center justify-between">
				{#if step > 0 && !isLastStep(step)}
					<Button variant="ghost" class="min-h-11 px-3" disabled={saving} onclick={() => (step = clampStep(step - 1))}>
						{t('common.back')}
					</Button>
				{:else}
					<span></span>
				{/if}
				{#if current.skippable}
					<Button variant="ghost" class="min-h-11 px-3" onclick={() => (step = clampStep(step + 1))}>
						{t('onboarding.skipForNow')}
					</Button>
				{/if}
			</div>
		</div>
	</footer>
</div>
