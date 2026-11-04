<script lang="ts">
	import { Select, toast } from '@facile/muse';
	import { backend, type InboxTab } from '$lib/backend';

	let {
		accountId,
		address,
		onchange
	}: { accountId: number; address: string; onchange?: () => void } = $props();

	let value = $state<InboxTab | 'auto' | ''>('');
	let saving = $state(false);

	const labels: Record<InboxTab | 'auto', string> = {
		primary: 'Primary',
		newsletters: 'Newsletters',
		notifications: 'Notifications',
		auto: 'Let Courrier sort it'
	};

	async function pin(next: InboxTab | 'auto') {
		saving = true;
		try {
			await backend.setSenderRule(accountId, address, next);
			toast.success(next === 'auto' ? `${address} is sorted automatically again` : `Mail from ${address} goes to ${labels[next]}`);
			onchange?.();
		} catch {
			toast.danger('Could not move this sender');
		}
		saving = false;
		value = '';
	}
</script>

<Select
	bind:value
	class="h-8 max-w-44 text-fc-xs"
	aria-label="Always file mail from this sender in"
	disabled={saving}
	onchange={() => value && pin(value)}
>
	<option value="">Move sender to…</option>
	<option value="primary">{labels.primary}</option>
	<option value="newsletters">{labels.newsletters}</option>
	<option value="notifications">{labels.notifications}</option>
	<option value="auto">{labels.auto}</option>
</Select>
