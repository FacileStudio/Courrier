<script lang="ts">
	import { getContext } from 'svelte';
	import {
		Alert,
		Avatar,
		Badge,
		Card,
		Button,
		ConfirmModal,
		EmptyState,
		SettingsSection,
		StatusDot,
		Tooltip,
		icons
	} from '@facile/muse';
	import { ApiError, backend, type CheckResult, type MailAccount } from '$lib/backend';
	import { spaceStore } from '$lib/stores/space.svelte';
	import AccountDrawer from '$lib/components/settings/AccountDrawer.svelte';
	import { t } from '$lib/i18n/index.svelte';

	type Probe =
		| { state: 'checking' }
		| { state: 'checked'; result: CheckResult }
		| { state: 'throttled'; retryAt: number }
		| { state: 'unreachable'; error: string };
	type Tone = 'neutral' | 'success' | 'warning' | 'danger';

	const app = getContext<{ refreshAccounts: () => Promise<void> }>('app');

	let accounts = $state<MailAccount[]>([]);
	let loadError = $state('');

	let probes = $state<Record<number, Probe>>({});

	let editorOpen = $state(false);
	let editing = $state<MailAccount | null>(null);

	let removeTarget = $state<MailAccount | null>(null);
	let removeOpen = $state(false);
	let removeError = $state('');

	let promoting = $state<number | null>(null);

	$effect(() => {
		spaceStore.active;
		void load();
	});

	async function load() {
		try {
			accounts = (await backend.listAccounts(spaceStore.active?.id)).accounts;
			loadError = '';
		} catch (err) {
			accounts = [];
			loadError = err instanceof Error ? err.message : t('accounts.loadFailed');
		}
	}

	function legs(account: MailAccount): { tone: Tone; label: string; pulse: boolean }[] {
		const probe = probes[account.id];

		if (probe?.state === 'checking') {
			return [{ tone: 'warning', label: t('accounts.checkRunning'), pulse: true }];
		}
		if (probe?.state === 'throttled') {
			return [{ tone: 'warning', label: t('accounts.notChecked', { reason: retryIn(probe) }), pulse: false }];
		}
		if (probe?.state === 'unreachable') {
			return [{ tone: 'danger', label: t('accounts.checkFailedToRun'), pulse: false }];
		}

		return (['imap', 'smtp'] as const).map((protocol) => {
			const name = protocol.toUpperCase();
			const host = protocol === 'imap' ? account.imap_host : account.smtp_host;
			if (!host) return { tone: 'neutral' as Tone, label: t('accounts.noServer', { protocol: name }), pulse: false };

			const leg = probe?.state === 'checked' ? probe.result[protocol] : undefined;
			if (!leg) return { tone: 'neutral' as Tone, label: t('accounts.legNotChecked', { protocol: name }), pulse: false };
			if (leg.ok) return { tone: 'success' as Tone, label: t('accounts.legConnected', { protocol: name }), pulse: false };
			return { tone: 'danger' as Tone, label: failureLabel(name, leg.error), pulse: false };
		});
	}

	function failureLabel(protocol: string, error?: string) {
		return /auth|password|credential|login|invalid/i.test(error ?? '')
			? t('accounts.legRejected', { protocol })
			: t('accounts.legUnreachable', { protocol });
	}

	let now = $state(Date.now());

	$effect(() => {
		if (!Object.values(probes).some((probe) => probe.state === 'throttled')) return;
		const timer = setInterval(() => {
			now = Date.now();
			for (const [id, probe] of Object.entries(probes)) {
				if (probe.state === 'throttled' && probe.retryAt <= now) delete probes[Number(id)];
			}
		}, 1000);
		return () => clearInterval(timer);
	});

	function retryIn(probe: Extract<Probe, { state: 'throttled' }>) {
		const seconds = Math.max(0, Math.ceil((probe.retryAt - now) / 1000));
		return seconds > 0 ? t('accounts.retryIn', { seconds }) : t('accounts.retryNow');
	}

	async function check(account: MailAccount) {
		probes[account.id] = { state: 'checking' };
		try {
			probes[account.id] = { state: 'checked', result: await backend.checkAccount(account.id) };
		} catch (err) {
			if (err instanceof ApiError && err.status === 429) {
				probes[account.id] = {
					state: 'throttled',
					retryAt: Date.now() + (err.retryAfter ?? 60) * 1000
				};
				return;
			}
			probes[account.id] = {
				state: 'unreachable',
				error: err instanceof Error ? err.message : t('accounts.checkRunFailed')
			};
		}
	}

	async function makeDefault(account: MailAccount) {
		promoting = account.id;
		try {
			await backend.updateAccount(account.id, { is_default: true });
			await load();
			await app.refreshAccounts();
		} catch (err) {
			loadError = err instanceof Error ? err.message : t('accounts.defaultFailed');
		}
		promoting = null;
	}

	function openNew() {
		editing = null;
		editorOpen = true;
	}

	function openEdit(account: MailAccount) {
		editing = account;
		editorOpen = true;
	}

	async function remove() {
		const target = removeTarget;
		if (!target) return;
		removeError = '';
		try {
			await backend.deleteAccount(target.id);
			delete probes[target.id];
			removeTarget = null;
			await load();
			await app.refreshAccounts();
		} catch (err) {
			removeError = err instanceof Error ? err.message : t('accounts.removeFailed');
		}
	}

	const failures = $derived(
		accounts.flatMap((account) => {
			const probe = probes[account.id];
			if (probe?.state === 'unreachable') {
				return [{ key: `${account.id}-run`, account, protocol: '', error: probe.error }];
			}
			if (probe?.state !== 'checked') return [];
			return (['imap', 'smtp'] as const)
				.filter((protocol) => probe.result[protocol].error)
				.map((protocol) => ({
					key: `${account.id}-${protocol}`,
					account,
					protocol: protocol.toUpperCase(),
					error: probe.result[protocol].error as string
				}));
		})
	);
</script>

<div class="flex flex-col gap-10">
	<SettingsSection
		title={t('accounts.title')}
		description={spaceStore.active
			? t('accounts.bodySpace', { space: spaceStore.active.name })
			: t('accounts.body')}
		bare
	>
		{#snippet actions()}
			<Button icon={icons.plus} onclick={openNew}>{t('accounts.add')}</Button>
		{/snippet}

		{#if loadError}
			<Alert tone="danger">{loadError}</Alert>
		{/if}

		{#if accounts.length === 0}
			<EmptyState
				icon={icons.mail}
				title={t('accounts.empty')}
				description={t('accounts.emptyBody')}
			>
				<Button icon={icons.plus} onclick={openNew}>{t('accounts.addFirst')}</Button>
			</EmptyState>
		{:else}
			<div class="flex flex-col gap-3">
				{#each accounts as account (account.id)}
					{@const state = legs(account)}
					{@const accountFailures = failures.filter((failure) => failure.account.id === account.id)}
					{@const checkDisabled =
						probes[account.id]?.state === 'checking' ||
						probes[account.id]?.state === 'throttled' ||
						(!account.imap_host && !account.smtp_host)}
					<Card class="flex flex-col gap-4">
						<div class="flex min-w-0 items-center gap-3">
							<Avatar name={account.name} />
							<div class="flex min-w-0 flex-col">
								<span class="flex min-w-0 items-center gap-2 font-medium text-fc-fg">
									<span class="truncate">{account.name}</span>
									{#if account.is_default}<Badge tone="accent">{t('accounts.default')}</Badge>{/if}
								</span>
								<span class="truncate text-fc-xs text-fc-fg-muted">{account.email}</span>
							</div>
						</div>

						<div class="flex flex-wrap items-center gap-x-4 gap-y-1">
							{#each state as leg (leg.label)}
								<StatusDot tone={leg.tone} label={leg.label} pulse={leg.pulse} />
							{/each}
						</div>

						<p class="font-fc-mono text-fc-xs truncate text-fc-fg-muted">
							IMAP {account.imap_host || '-'}:{account.imap_port} · SMTP {account.smtp_host || '-'}:{account.smtp_port}
						</p>

						{#each accountFailures as failure (failure.key)}
							<Alert
								tone="danger"
								title={failure.protocol ? t('accounts.protocolFailed', { protocol: failure.protocol }) : t('accounts.couldNotCheck')}
							>
								{failure.error}
							</Alert>
						{/each}

						<div class="flex flex-col gap-2 sm:flex-row sm:flex-wrap sm:items-center">
							<Tooltip text={t('accounts.checkHint')} class="contents">
								<Button
									variant="ghost"
									size="sm"
									icon={icons.plug}
									disabled={checkDisabled}
									onclick={() => check(account)}
								>
									{t('accounts.checkConnection')}
								</Button>
							</Tooltip>
							{#if !account.is_default}
								<Tooltip text={t('accounts.makeDefaultHint')} class="contents">
									<Button
										variant="ghost"
										size="sm"
										icon={icons.check}
										disabled={promoting === account.id}
										onclick={() => makeDefault(account)}
									>
										{t('accounts.makeDefault')}
									</Button>
								</Tooltip>
							{/if}
							<Button
								variant="ghost"
								size="sm"
								icon={icons.edit}
								onclick={() => openEdit(account)}
							>
								{t('common.edit')}
							</Button>
						</div>
					</Card>
				{/each}
			</div>

			<p class="text-fc-xs text-fc-fg-muted">
				{t('accounts.defaultNote')}
			</p>
		{/if}

		{#if removeError}
			<Alert tone="danger">{removeError}</Alert>
		{/if}
	</SettingsSection>
</div>

<AccountDrawer
	bind:open={editorOpen}
	account={editing}
	isFirst={accounts.length === 0}
	onRemove={(target) => {
		removeTarget = target;
		removeError = '';
		removeOpen = true;
	}}
	onSaved={async () => {
		await load();
		await app.refreshAccounts();
	}}
/>

<ConfirmModal
	bind:open={removeOpen}
	tone="danger"
	title={t('accounts.removeTitle')}
	description={t('accounts.removeBody', {
		name: removeTarget?.name ?? t('accounts.thisAccount'),
		host: removeTarget?.imap_host || t('accounts.theServer')
	})}
	confirmLabel={t('accounts.removeConfirm')}
	cancelLabel={t('accounts.keepIt')}
	onConfirm={remove}
/>
