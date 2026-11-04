<script lang="ts">
	import type { InboxTab } from '$lib/backend';

	let { tab = $bindable('primary') }: { tab?: InboxTab } = $props();

	const tabs: { value: InboxTab; label: string }[] = [
		{ value: 'primary', label: 'Primary' },
		{ value: 'newsletters', label: 'Newsletters' },
		{ value: 'notifications', label: 'Notifications' }
	];
</script>

<div role="tablist" aria-label="Inbox tabs" class="flex gap-1 overflow-x-auto border-b border-fc-border px-3 sm:px-4">
	{#each tabs as item (item.value)}
		<button
			type="button"
			role="tab"
			aria-selected={tab === item.value}
			class="-mb-px min-h-10 shrink-0 border-b-2 px-3 text-fc-sm font-medium transition-colors focus-visible:outline-2 focus-visible:outline-fc-ring {tab ===
			item.value
				? 'border-fc-accent text-fc-fg'
				: 'border-transparent text-fc-fg-muted hover:text-fc-fg'}"
			onclick={() => (tab = item.value)}
		>
			{item.label}
		</button>
	{/each}
</div>
