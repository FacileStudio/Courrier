<script lang="ts">
	import { Tabs } from '@facile/muse';
	import type { InboxTab } from '$lib/backend';

	let { tab = $bindable('primary') }: { tab?: InboxTab } = $props();

	const items = [
		{ id: 'primary', label: 'Primary' },
		{ id: 'newsletters', label: 'Newsletters' },
		{ id: 'notifications', label: 'Notifications' }
	];
</script>

<div class="px-3 sm:px-4">
	<Tabs {items} value={tab} label="Inbox tabs" onChange={(id) => (tab = id as InboxTab)} />
</div>
