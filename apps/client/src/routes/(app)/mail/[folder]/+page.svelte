<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { getContext, onMount } from 'svelte';
	import DOMPurify from 'dompurify';
	import { backend, type EmailMessage, type EmailAttachment, type MailAccount } from '$lib/backend';
	import { Button } from '$lib/components/ui/button';
	import { RefreshCw, Paperclip, Download, Reply, ReplyAll, Forward, Loader2 } from 'lucide-svelte';

	const app = getContext<{
		token: string;
		defaultAccountId: number | null;
		accounts: MailAccount[];
		folders: { id: number; type: string }[];
		refreshAccounts: () => Promise<void>;
	}>('app');

	const folderLabels: Record<string, string> = {
		sent: 'Sent',
		drafts: 'Drafts',
		archive: 'Archive',
		junk: 'Junk',
		trash: 'Trash'
	};

	const folderSlug = $derived(page.params.folder ?? '');
	const folderLabel = $derived(folderLabels[folderSlug] ?? folderSlug);
	const validFolder = $derived(folderSlug !== '' && folderSlug in folderLabels);

	let emails = $state<EmailMessage[]>([]);
	let selectedId = $state<number | null>(null);
	let loading = $state(false);
	let syncing = $state(false);
	let currentPage = $state(1);
	let totalEmails = $state(0);
	let loadingMore = $state(false);
	let listContainer = $state<HTMLDivElement | null>(null);
	let sentinel = $state<HTMLDivElement | null>(null);
	let resourceToken = $state<string | null>(null);

	const selected = $derived(emails.find((e) => e.id === selectedId) ?? null);
	const hasMore = $derived(emails.length < totalEmails);
	const LIMIT = 30;

	$effect(() => {
		if (!validFolder) {
			goto('/mail');
		}
	});

	async function ensureResourceToken(): Promise<string> {
		if (!resourceToken) {
			const res = await backend.getResourceToken(app.token);
			resourceToken = res.token;
		}
		return resourceToken;
	}

	function resolveCIDImages(html: string, accountId: number, emailId: number, token: string): string {
		return html.replace(/src=["']cid:([^"']+)["']/gi, (_match, cid) => {
			return `src="${backend.getCIDImageUrl(token, accountId, emailId, cid)}"`;
		});
	}

	function sanitizeHTML(html: string): string {
		if (!html || !app.defaultAccountId || !selectedId || !resourceToken) return html;
		const resolved = resolveCIDImages(html, app.defaultAccountId, selectedId, resourceToken);
		return DOMPurify.sanitize(resolved);
	}

	async function loadEmails() {
		if (!app.defaultAccountId || !validFolder) return;
		loading = true;
		currentPage = 1;
		try {
			const result = await backend.getEmailsByFolder(app.token, app.defaultAccountId, folderSlug, 1, LIMIT);
			emails = result.emails;
			totalEmails = result.total;
		} catch {
			emails = [];
			totalEmails = 0;
		}
		loading = false;
	}

	async function loadMoreEmails() {
		if (!app.defaultAccountId || loadingMore || !hasMore || !validFolder) return;
		loadingMore = true;
		const nextPage = currentPage + 1;
		try {
			const result = await backend.getEmailsByFolder(app.token, app.defaultAccountId, folderSlug, nextPage, LIMIT);
			emails = [...emails, ...result.emails];
			totalEmails = result.total;
			currentPage = nextPage;
		} catch {}
		loadingMore = false;
	}

	async function syncAndLoad() {
		if (!app.defaultAccountId) return;
		syncing = true;
		try {
			await backend.syncAccount(app.token, app.defaultAccountId);
			await app.refreshAccounts();
			const folder = app.folders.find((f) => f.type === folderSlug);
			if (folder) {
				await backend.syncFolder(app.token, app.defaultAccountId, folder.id);
			}
			await loadEmails();
		} catch {
		}
		syncing = false;
	}

	async function openEmail(email: EmailMessage) {
		selectedId = email.id;
		if (!app.defaultAccountId) return;

		try {
			await ensureResourceToken();
		} catch {}

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

	function formatFileSize(bytes: number): string {
		if (bytes < 1024) return `${bytes} B`;
		if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
		return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
	}

	async function downloadAttachment(attachment: EmailAttachment) {
		if (!app.defaultAccountId || !selected) return;
		await backend.downloadAttachment(app.token, app.defaultAccountId, selected.id, attachment.id, attachment.filename);
	}

	onMount(async () => {
		await loadEmails();
		if (emails.length === 0 && app.defaultAccountId) {
			await syncAndLoad();
		}
	});

	$effect(() => {
		if (!sentinel) return;
		const observer = new IntersectionObserver(
			(entries) => {
				if (entries[0]?.isIntersecting && hasMore && !loadingMore) {
					loadMoreEmails();
				}
			},
			{ root: listContainer, threshold: 0.1 }
		);
		observer.observe(sentinel);
		return () => observer.disconnect();
	});
</script>

<svelte:head>
	<title>{folderLabel} — Courrier</title>
</svelte:head>

<div class="flex h-full">
	<div class="w-80 flex-shrink-0 border-r flex flex-col">
		<div class="flex items-center justify-between border-b px-4 py-3">
			<h2 class="text-lg font-semibold">{folderLabel}</h2>
			<Button variant="ghost" size="icon" class="h-8 w-8" onclick={syncAndLoad} disabled={syncing}>
				<RefreshCw class="h-4 w-4 {syncing ? 'animate-spin' : ''}" />
			</Button>
		</div>

		<div class="flex-1 overflow-auto" bind:this={listContainer}>
			{#if loading}
				<div class="flex flex-col gap-0">
					{#each Array(5) as _, i}
						<div class="px-4 py-3 border-b mail-skeleton" style="--delay: {i * 80}ms">
							<div class="flex items-center gap-3">
								<div class="h-8 w-8 shrink-0 rounded-full bg-muted skeleton-pulse"></div>
								<div class="min-w-0 flex-1 space-y-2">
									<div class="flex items-center justify-between gap-2">
										<div class="h-3.5 w-28 rounded bg-muted skeleton-pulse"></div>
										<div class="h-3 w-10 rounded bg-muted skeleton-pulse"></div>
									</div>
									<div class="h-3 w-40 rounded bg-muted skeleton-pulse"></div>
								</div>
							</div>
						</div>
					{/each}
				</div>
			{:else if emails.length === 0}
				<div class="flex flex-col items-center justify-center h-full text-muted-foreground mail-fade-in">
					<p class="text-sm">No emails in {folderLabel}</p>
				</div>
			{:else}
				{#each emails as email, i}
					<button
						class="mail-list-item w-full text-left px-4 py-3 border-b transition-colors duration-150 hover:bg-muted/50
							{selectedId === email.id ? 'bg-muted' : ''}
							{!email.is_read ? 'font-medium' : ''}"
						style="--delay: {Math.min(i, 15) * 30}ms"
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
				{#if loadingMore}
					<div class="flex items-center justify-center py-4">
						<Loader2 class="h-4 w-4 animate-spin text-muted-foreground" />
					</div>
				{/if}
				<div bind:this={sentinel} class="h-1"></div>
			{/if}
		</div>
	</div>

	<div class="flex-1 flex flex-col">
		{#if selected}
			<div class="border-b px-6 py-4 mail-detail-header">
				<h1 class="text-xl font-semibold">{selected.subject || '(no subject)'}</h1>
				<div class="mt-2 flex items-center gap-3 text-sm text-muted-foreground">
					<span class="font-medium text-foreground">{selected.from_name || selected.from_address}</span>
					<span>&lt;{selected.from_address}&gt;</span>
					<div class="ml-auto flex items-center gap-1">
						<Button variant="ghost" size="icon" class="h-7 w-7" onclick={() => goto(`/mail/compose?reply=${selected!.id}&accountId=${app.defaultAccountId}`)}>
							<Reply class="h-4 w-4" />
						</Button>
						<Button variant="ghost" size="icon" class="h-7 w-7" onclick={() => goto(`/mail/compose?replyall=${selected!.id}&accountId=${app.defaultAccountId}`)}>
							<ReplyAll class="h-4 w-4" />
						</Button>
						<Button variant="ghost" size="icon" class="h-7 w-7" onclick={() => goto(`/mail/compose?forward=${selected!.id}&accountId=${app.defaultAccountId}`)}>
							<Forward class="h-4 w-4" />
						</Button>
						<span class="ml-2 text-xs">{formatDate(selected.date)}</span>
					</div>
				</div>
			</div>
			{#if selected.attachments && selected.attachments.length > 0}
				<div class="border-b px-6 py-3 mail-attachments">
					<div class="flex items-center gap-2 text-sm text-muted-foreground mb-2">
						<Paperclip class="h-4 w-4" />
						<span>{selected.attachments.length} attachment{selected.attachments.length > 1 ? 's' : ''}</span>
					</div>
					<div class="flex flex-wrap gap-2">
						{#each selected.attachments as attachment}
							<Button
								variant="outline"
								size="sm"
								class="gap-2 text-xs"
								onclick={() => downloadAttachment(attachment)}
							>
								<Download class="h-3.5 w-3.5" />
								<span class="max-w-48 truncate">{attachment.filename}</span>
								<span class="text-muted-foreground">({formatFileSize(attachment.size)})</span>
							</Button>
						{/each}
					</div>
				</div>
			{/if}
			<div class="flex-1 overflow-auto px-6 py-4 mail-body-content">
				{#if selected.body_html}
					{@html sanitizeHTML(selected.body_html)}
				{:else if selected.body_text}
					<pre class="whitespace-pre-wrap text-sm">{selected.body_text}</pre>
				{:else}
					<div class="flex items-center gap-2 text-sm text-muted-foreground">
						<Loader2 class="h-4 w-4 animate-spin" />
						<span>Loading message body...</span>
					</div>
				{/if}
			</div>
		{:else}
			<div class="flex flex-1 items-center justify-center text-muted-foreground">
				<p class="text-sm">Select an email to read</p>
			</div>
		{/if}
	</div>
</div>

<style>
	@media (prefers-reduced-motion: no-preference) {
		.mail-list-item {
			animation: mail-fade-in 200ms ease-out both;
			animation-delay: var(--delay, 0ms);
		}

		.mail-detail-header {
			animation: mail-slide-in 180ms ease-out both;
		}

		.mail-body-content {
			animation: mail-fade-in 200ms ease-out 60ms both;
		}

		.mail-attachments {
			animation: mail-slide-down 180ms ease-out both;
		}

		.mail-fade-in {
			animation: mail-fade-in 200ms ease-out both;
		}

		.skeleton-pulse {
			animation: skeleton-pulse 1.5s ease-in-out infinite;
		}

		.mail-skeleton {
			animation: mail-fade-in 200ms ease-out both;
			animation-delay: var(--delay, 0ms);
		}
	}

	@keyframes mail-fade-in {
		from {
			opacity: 0;
			transform: translateY(4px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	@keyframes mail-slide-in {
		from {
			opacity: 0;
			transform: translateX(8px);
		}
		to {
			opacity: 1;
			transform: translateX(0);
		}
	}

	@keyframes mail-slide-down {
		from {
			opacity: 0;
			transform: translateY(-4px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	@keyframes skeleton-pulse {
		0%, 100% {
			opacity: 0.4;
		}
		50% {
			opacity: 0.8;
		}
	}
</style>
