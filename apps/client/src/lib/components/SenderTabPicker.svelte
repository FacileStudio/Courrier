<script lang="ts">
	import { Select, toast } from '@facile/muse';
	import { backend, type InboxTab } from '$lib/backend';
	import { t } from '$lib/i18n/index.svelte';

	let {
		accountId,
		address,
		onchange
	}: { accountId: number; address: string; onchange?: () => void } = $props();

	let value = $state<InboxTab | 'auto' | ''>('');
	let saving = $state(false);

	const labels = $derived<Record<InboxTab | 'auto', string>>({
		primary: t('tabs.primary'),
		newsletters: t('tabs.newsletters'),
		notifications: t('tabs.notifications'),
		auto: t('sender.auto')
	});

	async function pin(next: InboxTab | 'auto') {
		saving = true;
		try {
			await backend.setSenderRule(accountId, address, next);
			toast.success(
				next === 'auto'
					? t('sender.autoDone', { address })
					: t('sender.pinned', { address, tab: labels[next] })
			);
			onchange?.();
		} catch {
			toast.danger(t('sender.failed'));
		}
		saving = false;
		value = '';
	}
</script>

<Select
	bind:value
	class="h-8 max-w-44 text-fc-xs"
	aria-label={t('sender.label')}
	disabled={saving}
	placeholder={t('sender.moveTo')}
	options={[
		{ value: 'primary', label: labels.primary },
		{ value: 'newsletters', label: labels.newsletters },
		{ value: 'notifications', label: labels.notifications },
		{ value: 'auto', label: labels.auto }
	]}
	onchange={() => value && pin(value)}
/>
