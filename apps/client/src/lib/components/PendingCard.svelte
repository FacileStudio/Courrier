<script lang="ts">
	import { Badge, Button, Card, icons } from '@facile/muse';
	import type { PendingAction } from '$lib/backend';
	import { formatDate } from '$lib/i18n/format';
	import { t } from '$lib/i18n/index.svelte';

	let {
		action,
		busy,
		onconfirm,
		ondismiss
	}: {
		action: PendingAction;
		busy: boolean;
		onconfirm: (action: PendingAction) => void;
		ondismiss: (action: PendingAction) => void;
	} = $props();

	const sending = $derived(action.kind === 'send');
	const ready = $derived(action.state === 'ready');
	const draftHref = $derived(
		action.draft ? `/mail/compose?accountId=${action.account_id}&draft=${action.draft.id}` : ''
	);

	/** Names the client that asked: the host of a client URL, else its id. */
	function clientName(clientId: string): string {
		if (clientId === 'bulle') return 'Bulle';
		try {
			return new URL(clientId).host;
		} catch {
			return clientId || '?';
		}
	}
</script>

<Card class="flex flex-col gap-4">
	<div class="flex flex-wrap items-center justify-between gap-2">
		<Badge tone={sending ? 'accent' : 'danger'}>{sending ? t('pending.send') : t('pending.deleteDraft')}</Badge>
		<span class="text-fc-xs text-fc-fg-muted">
			{t('pending.askedBy', { client: clientName(action.client_id) })} · {formatDate(action.created_at, {
				day: 'numeric',
				month: 'short',
				hour: '2-digit',
				minute: '2-digit'
			})}
		</span>
	</div>

	{#if action.draft}
		<dl class="grid grid-cols-[auto_1fr] gap-x-3 gap-y-1 text-fc-sm">
			<dt class="text-fc-fg-muted">{t('pending.from')}</dt>
			<dd class="min-w-0 break-words text-fc-fg">{action.account_email}</dd>
			<dt class="text-fc-fg-muted">{t('pending.to')}</dt>
			<dd class="min-w-0 break-words text-fc-fg">{action.draft.to.join(', ')}</dd>
			{#if action.draft.cc.length}
				<dt class="text-fc-fg-muted">{t('pending.cc')}</dt>
				<dd class="min-w-0 break-words text-fc-fg">{action.draft.cc.join(', ')}</dd>
			{/if}
			{#if action.draft.bcc.length}
				<dt class="text-fc-fg-muted">{t('pending.bcc')}</dt>
				<dd class="min-w-0 break-words text-fc-fg">{action.draft.bcc.join(', ')}</dd>
			{/if}
		</dl>
		<div class="flex flex-col gap-2 rounded-fc-md bg-fc-surface p-4">
			<p class="text-fc-md font-medium text-fc-fg">{action.draft.subject || t('common.noSubject')}</p>
			<p class="max-h-80 overflow-y-auto whitespace-pre-wrap break-words text-fc-sm text-fc-fg-muted">
				{action.draft.body}
			</p>
		</div>
	{/if}

	{#if action.state === 'changed'}
		<p class="text-fc-sm text-fc-warning">{t('pending.changed')}</p>
	{:else if action.state === 'missing'}
		<p class="text-fc-sm text-fc-warning">{t('pending.missing')}</p>
	{/if}

	<div class="flex flex-col-reverse gap-2 sm:flex-row sm:justify-end">
		<Button variant="ghost" class="w-full sm:w-auto" disabled={busy} onclick={() => ondismiss(action)}>
			{t('pending.dismiss')}
		</Button>
		{#if draftHref}
			<Button variant="secondary" class="w-full sm:w-auto" href={draftHref} icon={icons.edit}>
				{t('pending.openDraft')}
			</Button>
		{/if}
		{#if ready}
			<Button
				variant={sending ? 'primary' : 'danger'}
				class="w-full sm:w-auto"
				disabled={busy}
				icon={sending ? icons.mail : icons.remove}
				onclick={() => onconfirm(action)}
			>
				{sending ? t('pending.confirmSend') : t('pending.confirmDelete')}
			</Button>
		{/if}
	</div>
</Card>
