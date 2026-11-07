<script lang="ts">
	import {
		Alert,
		Button,
		Divider,
		Drawer,
		Field,
		Input,
		REDACTED,
		SecretField,
		Spinner,
		StatusDot,
		Switch,
		Textarea,
		isRedacted,
		icons
	} from '@facile/muse';
	import { ApiError, backend, type MailAccount } from '$lib/backend';
	import type { PlanLimit } from '$lib/plan-limit';
	import PlanLimitAlert from './PlanLimitAlert.svelte';
	import { t } from '$lib/i18n/index.svelte';
	import { spaceStore } from '$lib/stores/space.svelte';

	let {
		open = $bindable(false),
		account = null,
		isFirst = false,
		onSaved,
		onRemove
	}: {
		open?: boolean;
		account?: MailAccount | null;
		isFirst?: boolean;
		onSaved?: () => void;
		onRemove?: (account: MailAccount) => void;
	} = $props();

	let name = $state('');
	let email = $state('');
	let imapHost = $state('');
	let imapPort = $state('993');
	let imapUser = $state('');
	let imapPassword = $state('');
	let smtpHost = $state('');
	let smtpPort = $state('587');
	let smtpUser = $state('');
	let smtpPassword = $state('');
	let signature = $state('');
	let isDefault = $state(false);

	let saving = $state(false);
	let error = $state('');
	let planLimit = $state<PlanLimit | undefined>();
	let testState = $state<'idle' | 'testing' | 'passed' | 'failed'>('idle');
	let testError = $state('');

	let loadedKey = $state<string>('');
	const key = $derived(String(account?.id ?? 'new'));

	$effect(() => {
		if (!open) {
			loadedKey = '';
			return;
		}
		if (key === loadedKey) return;
		loadedKey = key;
		error = '';
		testState = 'idle';
		testError = '';
		name = account?.name ?? '';
		email = account?.email ?? '';
		imapHost = account?.imap_host ?? '';
		imapPort = String(account?.imap_port ?? 993);
		imapUser = account?.imap_user ?? '';
		smtpHost = account?.smtp_host ?? '';
		smtpPort = String(account?.smtp_port ?? 587);
		smtpUser = account?.smtp_user ?? '';
		signature = account?.signature ?? '';
		isDefault = account?.is_default ?? isFirst;
		imapPassword = account ? REDACTED : '';
		smtpPassword = account ? REDACTED : '';
	});

	const editing = $derived(account !== null);
	const typedBoth = $derived(
		imapPassword.length > 0 &&
			!isRedacted(imapPassword) &&
			smtpPassword.length > 0 &&
			!isRedacted(smtpPassword)
	);
	const complete = $derived(
		name.trim().length > 0 &&
			email.trim().length > 0 &&
			imapHost.trim().length > 0 &&
			smtpHost.trim().length > 0
	);

	function port(value: string, fallback: number) {
		const parsed = Number.parseInt(value, 10);
		return Number.isFinite(parsed) && parsed > 0 ? parsed : fallback;
	}

	async function test() {
		testState = 'testing';
		testError = '';
		try {
			await backend.testConnection({
				imap_host: imapHost.trim(),
				imap_port: port(imapPort, 993),
				imap_user: imapUser.trim() || email.trim(),
				imap_password: imapPassword,
				smtp_host: smtpHost.trim(),
				smtp_port: port(smtpPort, 587),
				smtp_user: smtpUser.trim() || email.trim(),
				smtp_password: smtpPassword
			});
			testState = 'passed';
		} catch (err) {
			testState = 'failed';
			testError = err instanceof Error ? err.message : t('drawer.testError');
		}
	}

	async function save(event: Event) {
		event.preventDefault();
		saving = true;
		error = '';
		planLimit = undefined;
		try {
			const common = {
				name: name.trim(),
				email: email.trim(),
				imap_host: imapHost.trim(),
				imap_port: port(imapPort, 993),
				imap_user: imapUser.trim() || email.trim(),
				smtp_host: smtpHost.trim(),
				smtp_port: port(smtpPort, 587),
				smtp_user: smtpUser.trim() || email.trim(),
				signature,
				is_default: isDefault
			};

			if (account) {
				await backend.updateAccount(account.id, {
					...common,
					...(isRedacted(imapPassword) || imapPassword.length === 0
						? {}
						: { imap_password: imapPassword }),
					...(isRedacted(smtpPassword) || smtpPassword.length === 0
						? {}
						: { smtp_password: smtpPassword })
				});
			} else {
				await backend.createAccount({
					...common,
					imap_password: imapPassword,
					smtp_password: smtpPassword,
					space_id: spaceStore.active?.id
				});
			}
			open = false;
			onSaved?.();
		} catch (err) {
			planLimit = err instanceof ApiError ? err.planLimit : undefined;
			error = err instanceof Error ? err.message : t('drawer.saveFailed');
		}
		saving = false;
	}
</script>

<Drawer bind:open title={editing ? t('drawer.editTitle', { name: account?.name ?? '' }) : t('drawer.addTitle')}>
	<form class="flex flex-col gap-4" onsubmit={save}>
		<Field label={t('drawer.displayName')} helper={t('drawer.displayNameHelp')}>
			<Input bind:value={name} placeholder={t('drawer.namePlaceholder')} required disabled={saving} />
		</Field>

		<Field label={t('drawer.email')} helper={t('drawer.emailHelp')}>
			<Input bind:value={email} type="email" placeholder="you@example.com" required disabled={saving} />
		</Field>

		<Divider class="my-1" />

		<p class="text-fc-sm font-medium text-fc-fg">{t('drawer.imap')}</p>

		<div class="grid gap-4 sm:grid-cols-[1fr_7rem]">
			<Field label={t('drawer.host')}>
				<Input bind:value={imapHost} placeholder="imap.example.com" required disabled={saving} />
			</Field>
			<Field label={t('drawer.port')}>
				<Input bind:value={imapPort} inputmode="numeric" disabled={saving} />
			</Field>
		</div>

		<Field label={t('drawer.username')} helper={t('drawer.usernameHelp')}>
			<Input bind:value={imapUser} placeholder={email || 'you@example.com'} disabled={saving} />
		</Field>

		<SecretField
			bind:value={imapPassword}
			editable
			label={t('drawer.imapPassword')}
			helper={isRedacted(imapPassword) ? t('drawer.storedHelp') : t('drawer.imapPasswordHelp')}
			disabled={saving}
		/>

		<Divider class="my-1" />

		<p class="text-fc-sm font-medium text-fc-fg">{t('drawer.smtp')}</p>

		<div class="grid gap-4 sm:grid-cols-[1fr_7rem]">
			<Field label={t('drawer.host')}>
				<Input bind:value={smtpHost} placeholder="smtp.example.com" required disabled={saving} />
			</Field>
			<Field label={t('drawer.port')}>
				<Input bind:value={smtpPort} inputmode="numeric" disabled={saving} />
			</Field>
		</div>

		<Field label={t('drawer.username')} helper={t('drawer.usernameHelp')}>
			<Input bind:value={smtpUser} placeholder={email || 'you@example.com'} disabled={saving} />
		</Field>

		<SecretField
			bind:value={smtpPassword}
			editable
			label={t('drawer.smtpPassword')}
			helper={isRedacted(smtpPassword) ? t('drawer.storedHelp') : t('drawer.smtpPasswordHelp')}
			disabled={saving}
		/>

		<Divider class="my-1" />

		<Field label={t('drawer.signature')} helper={t('drawer.signatureHelp')}>
			<Textarea bind:value={signature} rows={4} placeholder={t('drawer.signaturePlaceholder')} disabled={saving} />
		</Field>

		<div class="flex flex-col gap-1">
			<Switch bind:checked={isDefault} label={t('drawer.useDefault')} disabled={saving} />
			<p class="pl-14 text-fc-xs text-fc-fg-muted">
				{t('drawer.useDefaultHelp')}
			</p>
		</div>

		<Divider class="my-1" />

		<div class="flex flex-col gap-3">
			<div class="flex flex-wrap items-center gap-2">
				<Button
					type="button"
					variant="outline"
					icon={icons.plug}
					disabled={saving || testState === 'testing' || !imapHost.trim() || !typedBoth}
					onclick={test}
				>
					{testState === 'testing' ? t('drawer.testing') : t('drawer.test')}
				</Button>

				{#if testState === 'testing'}
					<StatusDot tone="warning" label={t('drawer.testConnecting')} pulse />
				{:else if testState === 'passed'}
					<StatusDot tone="success" label={t('drawer.testPassed')} />
				{:else if testState === 'failed'}
					<StatusDot tone="danger" label={t('drawer.testFailedShort')} />
				{/if}
			</div>

			{#if !typedBoth}
				<p class="text-fc-xs text-fc-fg-muted">
					{t('drawer.testNeedsPasswords')}
				</p>
			{/if}

			{#if testState === 'failed' && testError}
				<Alert tone="danger" title={t('drawer.testFailedTitle')}>{testError}</Alert>
			{/if}
		</div>

		{#if planLimit}
			<PlanLimitAlert limit={planLimit} />
		{:else if error}
			<Alert tone="danger">{error}</Alert>
		{/if}

		<div class="flex justify-end gap-2 pt-1">
			<Button type="button" variant="ghost" disabled={saving} onclick={() => (open = false)}>
				{t('common.cancel')}
			</Button>
			<Button type="submit" disabled={saving || !complete}>
				{#if saving}<Spinner size="sm" />{/if}
				{saving ? t('common.saving') : editing ? t('drawer.saveAccount') : t('drawer.addAccount')}
			</Button>
		</div>

		{#if account && onRemove}
			<Divider class="my-1" />
			<Button
				type="button"
				variant="ghost-danger"
				icon={icons.remove}
				disabled={saving}
				onclick={() => {
					const target = account;
					open = false;
					onRemove(target);
				}}
			>
				{t('drawer.removeAccount')}
			</Button>
		{/if}
	</form>
</Drawer>
