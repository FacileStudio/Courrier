<script lang="ts">
	import { goto } from '$app/navigation';
	import { getContext, onMount } from 'svelte';
	import { MediaQuery } from 'svelte/reactivity';
	import DOMPurify from 'dompurify';
	import { backend, type EmailMessage, type EmailAttachment, type MailAccount } from '$lib/backend';
	import { mailCache } from '$lib/stores/mail-cache';
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import * as Resizable from '$lib/components/ui/resizable';
	import { toast } from 'svelte-sonner';
	import { RefreshCw, Paperclip, Download, Reply, ReplyAll, Forward, Loader2, Trash2, Archive, ArrowLeft } from 'lucide-svelte';
	import EmailItem from '$lib/components/EmailItem.svelte';
	import BulkActionBar from '$lib/components/BulkActionBar.svelte';
	import DeleteConfirmDialog from '$lib/components/DeleteConfirmDialog.svelte';

	const app = getContext<{
		defaultAccountId: number | null;
		accounts: MailAccount[];
		folders: { id: number; type: string }[];
		refreshAccounts: () => Promise<void>;
	}>('app');

	let emails = $state<EmailMessage[]>([]);
	let selectedId = $state<number | null>(null);
	let loading = $state(false);
	let syncing = $state(false);
	let currentPage = $state(1);
	let totalEmails = $state(0);
	let loadingMore = $state(false);
	let listContainer = $state<HTMLDivElement | null>(null);
	let sentinel = $state<HTMLDivElement | null>(null);
	let checkedIds = $state<Set<number>>(new Set());
	let deleteDialogOpen = $state(false);
	let deleteTarget = $state<number[]>([]);
	let bulkLoading = $state(false);

	const selected = $derived(emails.find((e) => e.id === selectedId) ?? null);
	const hasMore = $derived(emails.length < totalEmails);
	const selectionActive = $derived(checkedIds.size > 0);
	const allChecked = $derived(emails.length > 0 && emails.every((e) => checkedIds.has(e.id)));
	const LIMIT = 30;

	const isDesktop = new MediaQuery('(min-width: 768px)');

	function toggleCheck(id: number) {
		const next = new Set(checkedIds);
		if (next.has(id)) next.delete(id);
		else next.add(id);
		checkedIds = next;
	}

	function toggleAll() {
		if (allChecked) {
			checkedIds = new Set();
		} else {
			checkedIds = new Set(emails.map((e) => e.id));
		}
	}

	function clearSelection() {
		checkedIds = new Set();
	}

	function resolveCIDImages(html: string, accountId: number, emailId: number): string {
		return html.replace(/src=["']cid:([^"']+)["']/gi, (_match, cid) => {
			return `src="${backend.getCIDImageUrl(accountId, emailId, cid)}"`;
		});
	}

	function sanitizeHTML(html: string): string {
		if (!html || !app.defaultAccountId || !selectedId) return html;
		const resolved = resolveCIDImages(html, app.defaultAccountId, selectedId);
		return DOMPurify.sanitize(resolved);
	}

	async function loadEmails() {
		if (!app.defaultAccountId) return;

		const cached = mailCache.get(app.defaultAccountId, 'inbox', 1);
		if (cached) {
			emails = cached.emails;
			totalEmails = cached.total;
		}

		loading = !cached;
		currentPage = 1;
		try {
			const result = await backend.getEmailsByFolder(app.defaultAccountId, 'inbox', 1, LIMIT);
			emails = result.emails;
			totalEmails = result.total;
			mailCache.set(app.defaultAccountId, 'inbox', 1, result.emails, result.total);
		} catch {
			if (!cached) {
				emails = [];
				totalEmails = 0;
			}
		}
		loading = false;
	}

	async function loadMoreEmails() {
		if (!app.defaultAccountId || loadingMore || !hasMore) return;
		loadingMore = true;
		const nextPage = currentPage + 1;

		const cached = mailCache.get(app.defaultAccountId, 'inbox', nextPage);
		if (cached) {
			emails = [...emails, ...cached.emails];
			totalEmails = cached.total;
			currentPage = nextPage;
			loadingMore = false;
			return;
		}

		try {
			const result = await backend.getEmailsByFolder(app.defaultAccountId, 'inbox', nextPage, LIMIT);
			emails = [...emails, ...result.emails];
			totalEmails = result.total;
			currentPage = nextPage;
			mailCache.set(app.defaultAccountId, 'inbox', nextPage, result.emails, result.total);
		} catch {}
		loadingMore = false;
	}

	async function syncAndLoad() {
		if (!app.defaultAccountId) return;
		syncing = true;
		try {
			await backend.syncAccount(app.defaultAccountId);
			await app.refreshAccounts();
			const inboxFolder = app.folders.find((f) => f.type === 'inbox');
			if (inboxFolder) {
				await backend.syncFolder(app.defaultAccountId, inboxFolder.id);
			}
			mailCache.invalidateFolder(app.defaultAccountId, 'inbox');
			await loadEmails();
		} catch {}
		syncing = false;
	}

	async function openEmail(email: EmailMessage) {
		selectedId = email.id;
		if (!app.defaultAccountId) return;

		if (!email.body_text && !email.body_html) {
			try {
				const full = await backend.getEmail(app.defaultAccountId, email.id);
				emails = emails.map((e) => (e.id === email.id ? full : e));
			} catch {}
		}

		if (!email.is_read) {
			try {
				await backend.updateEmail(app.defaultAccountId, email.id, { is_read: true });
				emails = emails.map((e) => (e.id === email.id ? { ...e, is_read: true } : e));
			} catch {}
		}
	}

	function formatFileSize(bytes: number): string {
		if (bytes < 1024) return `${bytes} B`;
		if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
		return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
	}

	function formatDate(dateStr: string) {
		const date = new Date(dateStr);
		const now = new Date();
		const diffMs = now.getTime() - date.getTime();
		const diffDays = diffMs / (1000 * 60 * 60 * 24);
		if (diffDays < 1) return date.toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' });
		if (diffDays < 7) return date.toLocaleDateString('fr-FR', { weekday: 'short' });
		return date.toLocaleDateString('fr-FR', { day: 'numeric', month: 'short' });
	}

	async function downloadAttachment(attachment: EmailAttachment) {
		if (!app.defaultAccountId || !selected) return;
		await backend.downloadAttachment(app.defaultAccountId, selected.id, attachment.id, attachment.filename);
	}

	async function handleBulkDelete() {
		deleteTarget = [...checkedIds];
		deleteDialogOpen = true;
	}

	async function confirmDelete() {
		if (!app.defaultAccountId || deleteTarget.length === 0) return;
		bulkLoading = true;
		const count = deleteTarget.length;
		try {
			await backend.bulkAction(app.defaultAccountId, deleteTarget, 'delete');
			emails = emails.filter((e) => !deleteTarget.includes(e.id));
			totalEmails = Math.max(0, totalEmails - count);
			mailCache.removeEmails(deleteTarget);
			if (selectedId && deleteTarget.includes(selectedId)) selectedId = null;
			checkedIds = new Set();
			deleteTarget = [];
			toast.success(count === 1 ? 'Moved to trash' : `${count} moved to trash`);
		} catch {
			toast.error('Failed to delete');
		}
		bulkLoading = false;
	}

	async function handleBulkArchive() {
		if (!app.defaultAccountId) return;
		bulkLoading = true;
		const ids = [...checkedIds];
		try {
			await backend.bulkAction(app.defaultAccountId, ids, 'archive');
			emails = emails.filter((e) => !ids.includes(e.id));
			totalEmails = Math.max(0, totalEmails - ids.length);
			mailCache.removeEmails(ids);
			if (selectedId && ids.includes(selectedId)) selectedId = null;
			checkedIds = new Set();
			toast.success(`${ids.length} archived`);
		} catch {
			toast.error('Failed to archive');
		}
		bulkLoading = false;
	}

	async function handleBulkMarkRead() {
		if (!app.defaultAccountId) return;
		bulkLoading = true;
		const ids = [...checkedIds];
		try {
			await backend.bulkAction(app.defaultAccountId, ids, 'mark_read');
			emails = emails.map((e) => (ids.includes(e.id) ? { ...e, is_read: true } : e));
			checkedIds = new Set();
		} catch {
			toast.error('Failed to mark as read');
		}
		bulkLoading = false;
	}

	async function handleBulkMarkUnread() {
		if (!app.defaultAccountId) return;
		bulkLoading = true;
		const ids = [...checkedIds];
		try {
			await backend.bulkAction(app.defaultAccountId, ids, 'mark_unread');
			emails = emails.map((e) => (ids.includes(e.id) ? { ...e, is_read: false } : e));
			checkedIds = new Set();
		} catch {
			toast.error('Failed to mark as unread');
		}
		bulkLoading = false;
	}

	async function handleSingleDelete(emailId: number) {
		deleteTarget = [emailId];
		deleteDialogOpen = true;
	}

	async function handleSingleArchive(emailId: number) {
		if (!app.defaultAccountId) return;
		try {
			await backend.bulkAction(app.defaultAccountId, [emailId], 'archive');
			emails = emails.filter((e) => e.id !== emailId);
			totalEmails = Math.max(0, totalEmails - 1);
			mailCache.removeEmails([emailId]);
			if (selectedId === emailId) selectedId = null;
			toast.success('Archived');
		} catch {
			toast.error('Failed to archive');
		}
	}

	async function handleToggleRead(email: EmailMessage) {
		if (!app.defaultAccountId) return;
		try {
			await backend.updateEmail(app.defaultAccountId, email.id, { is_read: !email.is_read });
			emails = emails.map((e) => (e.id === email.id ? { ...e, is_read: !email.is_read } : e));
		} catch {}
	}

	async function handleToggleStar(email: EmailMessage) {
		if (!app.defaultAccountId) return;
		try {
			await backend.updateEmail(app.defaultAccountId, email.id, { is_starred: !email.is_starred });
			emails = emails.map((e) => (e.id === email.id ? { ...e, is_starred: !email.is_starred } : e));
		} catch {}
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

<DeleteConfirmDialog bind:open={deleteDialogOpen} count={deleteTarget.length} onconfirm={confirmDelete} />

{#snippet listPane()}
	<div class="flex h-full flex-col">
		<div class="flex items-center justify-between border-b px-4 py-3">
			<div class="flex items-center gap-2">
				{#if emails.length > 0}
					<button
						class="flex h-6 w-6 items-center justify-center rounded hover:bg-muted"
						onclick={toggleAll}
					>
						<Checkbox checked={allChecked} class="h-4 w-4" />
					</button>
				{/if}
				<h2 class="text-lg font-semibold">Inbox</h2>
			</div>
			<Button variant="ghost" size="icon" class="h-8 w-8" onclick={syncAndLoad} disabled={syncing}>
				<RefreshCw class="h-4 w-4 {syncing ? 'animate-spin' : ''}" />
			</Button>
		</div>

		<BulkActionBar
			count={checkedIds.size}
			loading={bulkLoading}
			ondelete={handleBulkDelete}
			onarchive={handleBulkArchive}
			onmarkread={handleBulkMarkRead}
			onmarkunread={handleBulkMarkUnread}
			onclear={clearSelection}
		/>

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
					{#if app.accounts.length === 0}
						<p class="text-sm">No mail accounts configured</p>
						<p class="text-xs mt-1">Add one in Settings</p>
					{:else}
						<p class="text-sm">No emails yet</p>
						<Button variant="outline" size="sm" class="mt-3 gap-1.5" onclick={syncAndLoad} disabled={syncing}>
							<RefreshCw class="h-4 w-4 {syncing ? 'animate-spin' : ''}" />
							{syncing ? 'Syncing...' : 'Sync now'}
						</Button>
					{/if}
				</div>
			{:else}
				{#each emails as email, i (email.id)}
					<div style="--delay: {Math.min(i, 15) * 30}ms" class="mail-list-item">
						<EmailItem
							{email}
							selected={selectedId === email.id}
							checked={checkedIds.has(email.id)}
							{selectionActive}
							onopen={() => openEmail(email)}
							ontogglecheck={() => toggleCheck(email.id)}
							onreply={() => goto(`/mail/compose?reply=${email.id}&accountId=${app.defaultAccountId}`)}
							onforward={() => goto(`/mail/compose?forward=${email.id}&accountId=${app.defaultAccountId}`)}
							onarchive={() => handleSingleArchive(email.id)}
							ondelete={() => handleSingleDelete(email.id)}
							ontoggleread={() => handleToggleRead(email)}
							ontogglestar={() => handleToggleStar(email)}
						/>
					</div>
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
{/snippet}

{#snippet readingPane()}
	<div class="flex h-full flex-col">
		{#if selected}
			<div class="border-b px-4 py-4 sm:px-6 mail-detail-header">
				<h1 class="text-xl font-semibold">{selected.subject || '(no subject)'}</h1>
				<div class="mt-2 flex flex-wrap items-center gap-x-3 gap-y-1 text-sm text-muted-foreground">
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
						<div class="mx-1 h-4 w-px bg-border"></div>
						<Button variant="ghost" size="icon" class="h-7 w-7" onclick={() => handleSingleArchive(selected!.id)}>
							<Archive class="h-4 w-4" />
						</Button>
						<Button variant="ghost" size="icon" class="h-7 w-7 text-destructive hover:text-destructive hover:bg-destructive/10" onclick={() => handleSingleDelete(selected!.id)}>
							<Trash2 class="h-4 w-4" />
						</Button>
						<span class="ml-2 text-xs">{formatDate(selected.date)}</span>
					</div>
				</div>
			</div>
			{#if selected.attachments && selected.attachments.length > 0}
				<div class="border-b px-4 py-3 sm:px-6 mail-attachments">
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
			<div class="flex-1 overflow-auto px-4 py-4 sm:px-6 mail-body-content">
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
{/snippet}

{#if isDesktop.current}
	<Resizable.PaneGroup direction="horizontal" class="h-full">
		<Resizable.Pane defaultSize={30} minSize={20} maxSize={50}>
			{@render listPane()}
		</Resizable.Pane>
		<Resizable.Handle />
		<Resizable.Pane defaultSize={70}>
			{@render readingPane()}
		</Resizable.Pane>
	</Resizable.PaneGroup>
{:else if selected}
	<div class="flex h-full flex-col">
		<div class="flex flex-shrink-0 items-center border-b px-2 py-1.5">
			<Button variant="ghost" size="sm" class="gap-1.5" onclick={() => (selectedId = null)}>
				<ArrowLeft class="h-4 w-4" />
				Back
			</Button>
		</div>
		<div class="min-h-0 flex-1">
			{@render readingPane()}
		</div>
	</div>
{:else}
	{@render listPane()}
{/if}

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
