<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import DOMPurify from 'dompurify';
	import { backend, type EmailMessage, type MailAccount } from '$lib/backend';
	import { Button } from '$lib/components/ui/button';
	import { RefreshCw } from 'lucide-svelte';

	const app = getContext<{
		token: string;
		defaultAccountId: number | null;
		accounts: MailAccount[];
		refreshAccounts: () => Promise<void>;
	}>('app');

	let emails = $state<EmailMessage[]>([]);
	let selectedId = $state<number | null>(null);
	let loading = $state(false);
	let syncing = $state(false);

	const selected = $derived(emails.find((e) => e.id === selectedId) ?? null);

	async function loadEmails() {
		if (!app.defaultAccountId) return;
		loading = true;
		try {
			const result = await backend.getEmailsByFolder(app.token, app.defaultAccountId, 'inbox');
			emails = result.emails;
		} catch {
			emails = [];
		}
		loading = false;
	}

	async function syncAndLoad() {
		if (!app.defaultAccountId) return;
		syncing = true;
		try {
			await backend.syncAccount(app.token, app.defaultAccountId);
			await app.refreshAccounts();
			await loadEmails();
		} catch {
		}
		syncing = false;
	}

	async function openEmail(email: EmailMessage) {
		selectedId = email.id;
		if (!app.defaultAccountId) return;

		if (!email.body_text && !email.body_html) {
			try {
				const full = await backend.getEmail(app.token, app.defaultAccountId, email.id);
				emails = emails.map((e) => (e.id === email.id ? full : e));
			} catch {}
		}

		if (!email.is_read) {
			try {
				await backend.updateEmail(app.token, app.defaultAccountId, email.id, { is_read: true });
				emails = emails.map((e) => (e.id === email.id ? { ...e, is_read: true } : e));
			} catch {}
		}
	}

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

	onMount(() => {
		loadEmails();
	});
</script>

<div class="flex h-full">
	<div class="w-80 flex-shrink-0 border-r flex flex-col">
		<div class="flex items-center justify-between border-b px-4 py-3">
			<h2 class="text-lg font-semibold">Inbox</h2>
			<Button variant="ghost" size="icon" class="h-8 w-8" onclick={syncAndLoad} disabled={syncing}>
				<RefreshCw class="h-4 w-4 {syncing ? 'animate-spin' : ''}" />
			</Button>
		</div>

		<div class="flex-1 overflow-auto">
			{#if loading}
				<div class="flex items-center justify-center h-full text-muted-foreground">
					<p class="text-sm">Loading...</p>
				</div>
			{:else if emails.length === 0}
				<div class="flex flex-col items-center justify-center h-full text-muted-foreground">
					{#if app.accounts.length === 0}
						<p class="text-sm">No mail accounts configured</p>
						<p class="text-xs mt-1">Add one in Settings</p>
					{:else}
						<p class="text-sm">No emails yet</p>
						<Button variant="outline" size="sm" class="mt-3" onclick={syncAndLoad} disabled={syncing}>
							{syncing ? 'Syncing...' : 'Sync now'}
						</Button>
					{/if}
				</div>
			{:else}
				{#each emails as email}
					<button
						class="w-full text-left px-4 py-3 border-b transition-colors hover:bg-muted/50
							{selectedId === email.id ? 'bg-muted' : ''}
							{!email.is_read ? 'font-medium' : ''}"
						onclick={() => openEmail(email)}
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
					{@html DOMPurify.sanitize(selected.body_html)}
				{:else if selected.body_text}
					<pre class="whitespace-pre-wrap text-sm">{selected.body_text}</pre>
				{:else}
					<p class="text-sm text-muted-foreground">Loading message body...</p>
				{/if}
			</div>
		{:else}
			<div class="flex flex-1 items-center justify-center text-muted-foreground">
				<p class="text-sm">Select an email to read</p>
			</div>
		{/if}
	</div>
</div>
