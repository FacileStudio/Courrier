<script lang="ts">
	import { Tabs } from '@facile/muse';
	import type { InboxTab } from '$lib/backend';
	import { t } from '$lib/i18n/index.svelte';

	let { tab = $bindable('primary') }: { tab?: InboxTab } = $props();

	const items = $derived([
		{ id: 'primary', label: t('tabs.primary') },
		{ id: 'newsletters', label: t('tabs.newsletters') },
		{ id: 'notifications', label: t('tabs.notifications') }
	]);
</script>

<div class="px-3 sm:px-4">
	<Tabs {items} value={tab} label={t('tabs.label')} onChange={(id) => (tab = id as InboxTab)} />
</div>
