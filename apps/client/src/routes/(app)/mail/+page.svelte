<script lang="ts">
	import { getContext } from 'svelte';
	import type { EmailMessage } from '$lib/backend';

	const app = getContext<{ token: string }>('app');

	let emails = $state<EmailMessage[]>([]);
	let selectedId = $state<number | null>(null);

	const selected = $derived(emails.find((e) => e.id === selectedId) ?? null);

	function formatDate(dateStr: string) {
		const date = new Date(dateStr);
		const now = new Date();
		const diffMs = now.getTime() - date.getTime();
		const diffDays = diffMs / (1000 * 60 * 60 * 24);

		if (diffDays < 1) {
			return date.toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' });
		}
		if (diffDays < 7) {
			return date.toLocaleDateString('fr-FR', { weekday: 'short' });
		}
		return date.toLocaleDateString('fr-FR', { day: 'numeric', month: 'short' });
	}

	function getInitials(name: string) {
		const parts = name.trim().split(/\s+/).filter(Boolean);
		if (parts.length === 0) return '?';
		if (parts.length === 1) return parts[0][0].toUpperCase();
		return `${parts[0][0]}${parts[1][0]}`.toUpperCase();
	}
</script>

<div class="flex h-full">
	<div class="w-80 flex-shrink-0 border-r flex flex-col">
		<div class="flex items-center justify-between border-b px-4 py-3">
			<h2 class="text-lg font-semibold">Inbox</h2>
		</div>

		<div class="flex-1 overflow-auto">
			{#if emails.length === 0}
				<div class="flex flex-col items-center justify-center h-full text-muted-foreground">
					<p class="text-sm">No emails yet</p>
					<p class="text-xs mt-1">Add a mail account in Settings</p>
				</div>
			{:else}
				{#each emails as email}
					<button
						class="w-full text-left px-4 py-3 border-b transition-colors hover:bg-muted/50
							{selectedId === email.id ? 'bg-muted' : ''}
							{!email.is_read ? 'font-medium' : ''}"
						onclick={() => (selectedId = email.id)}
					>
						<div class="flex items-center gap-3">
							<div class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-muted text-xs font-medium">
								{getInitials(email.from_name || email.from_address)}
							</div>
							<div class="min-w-0 flex-1">
								<div class="flex items-center justify-between gap-2">
									<span class="truncate text-sm">{email.from_name || email.from_address}</span>
									<span class="shrink-0 text-xs text-muted-foreground">{formatDate(email.date)}</span>
								</div>
								<p class="truncate text-sm text-muted-foreground">{email.subject || '(no subject)'}</p>
							</div>
							{#if !email.is_read}
								<div class="h-2 w-2 shrink-0 rounded-full bg-blue-600"></div>
							{/if}
						</div>
					</button>
				{/each}
			{/if}
		</div>
	</div>

	<div class="flex-1 flex flex-col">
		{#if selected}
			<div class="border-b px-6 py-4">
				<h1 class="text-xl font-semibold">{selected.subject || '(no subject)'}</h1>
				<div class="mt-2 flex items-center gap-3 text-sm text-muted-foreground">
					<span class="font-medium text-foreground">{selected.from_name || selected.from_address}</span>
					<span>&lt;{selected.from_address}&gt;</span>
					<span class="ml-auto">{formatDate(selected.date)}</span>
				</div>
			</div>
			<div class="flex-1 overflow-auto px-6 py-4">
				{#if selected.body_html}
					{@html selected.body_html}
				{:else}
					<pre class="whitespace-pre-wrap text-sm">{selected.body_text}</pre>
				{/if}
			</div>
		{:else}
			<div class="flex flex-1 items-center justify-center text-muted-foreground">
				<p class="text-sm">Select an email to read</p>
			</div>
		{/if}
	</div>
</div>
