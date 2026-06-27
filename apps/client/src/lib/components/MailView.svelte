<script lang="ts">
	import { goto } from '$app/navigation';
	import { getContext } from 'svelte';
	import { MediaQuery } from 'svelte/reactivity';
	import DOMPurify from 'dompurify';
	import { backend, type EmailMessage, type EmailAttachment, type MailAccount } from '$lib/backend';
	import { mailCache } from '$lib/stores/mail-cache';
	import { searchStore } from '$lib/stores/search.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Input } from '$lib/components/ui/input';
	import * as Resizable from '$lib/components/ui/resizable';
	import { toast } from 'svelte-sonner';
	import { RefreshCw, Paperclip, Download, Reply, ReplyAll, Forward, Loader2, Trash2, Archive, Mail, ArrowLeft, Search, X } from 'lucide-svelte';
	import EmailItem from '$lib/components/EmailItem.svelte';
	import BulkActionBar from '$lib/components/BulkActionBar.svelte';
	import DeleteConfirmDialog from '$lib/components/DeleteConfirmDialog.svelte';
	import MailFolderSwitcher from '$lib/components/MailFolderSwitcher.svelte';

	let { folder }: { folder: string } = $props();

	const app = getContext<{
		defaultAccountId: number | null;
		accounts: MailAccount[];
		folders: { id: number; type: string }[];
		refreshAccounts: () => Promise<void>;
	}>('app');

	const folderLabels: Record<string, string> = {
		inbox: 'Inbox',
		sent: 'Sent',
		drafts: 'Drafts',
		archive: 'Archive',
		junk: 'Junk',
		trash: 'Trash'
	};

	const folderLabel = $derived(folderLabels[folder] ?? folder);
	const isTrash = $derived(folder === 'trash');

	const isDesktop = new MediaQuery('(min-width: 768px)');

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
	let showUnreadOnly = $state(false);
	let searchQuery = $state('');
	let debouncedQuery = $state('');
	let inputEl = $state<HTMLInputElement | null>(null);
	let loadSeq = 0;
	let seenFocusSeq = searchStore.focusSeq;

	const selected = $derived(emails.find((e) => e.id === selectedId) ?? null);
	const hasMore = $derived(emails.length < totalEmails);
	const selectionActive = $derived(checkedIds.size > 0);
	const allChecked = $derived(emails.length > 0 && emails.every((e) => checkedIds.has(e.id)));
	const searchMode = $derived(debouncedQuery.trim() !== '');
	const LIMIT = 30;

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

	function clearSearch() {
		searchQuery = '';
		debouncedQuery = '';
	}

	function onSearchKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			clearSearch();
			inputEl?.blur();
		}
	}

	$effect(() => {
		const q = searchQuery;
		const timer = setTimeout(() => {
			debouncedQuery = q;
		}, 300);
		return () => clearTimeout(timer);
	});

	$effect(() => {
		if (searchStore.focusSeq > seenFocusSeq) {
			seenFocusSeq = searchStore.focusSeq;
			inputEl?.focus();
		}
	});

	function resolveCIDImages(html: string, accountId: number, emailId: number): string {
		return html.replace(/src=["']cid:([^"']+)["']/gi, (_match, cid) => {
			return `src="${backend.getCIDImageUrl(accountId, emailId, cid)}"`;
		});
	}

	const sanitizedBody = $derived.by(() => {
		const html = selected?.body_html;
		if (!html || !app.defaultAccountId || !selected) return '';
		const resolved = resolveCIDImages(html, app.defaultAccountId, selected.id);
		return DOMPurify.sanitize(resolved).replace(/<img /gi, '<img loading="lazy" decoding="async" ');
	});

	async function loadEmails() {
		if (!app.defaultAccountId) return;

		const cached = showUnreadOnly ? null : mailCache.get(app.defaultAccountId, folder, 1);
		if (cached) {
			emails = cached.emails;
			totalEmails = cached.total;
		}

		loading = !cached;
		currentPage = 1;
		try {
			const result = await backend.getEmailsByFolder(app.defaultAccountId, folder, 1, LIMIT, showUnreadOnly);
			emails = result.emails;
			totalEmails = result.total;
			if (!showUnreadOnly) {
				mailCache.set(app.defaultAccountId, folder, 1, result.emails, result.total);
			}
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
		const reqId = loadSeq;
		const nextPage = currentPage + 1;

		if (searchMode) {
			try {
				const result = await backend.searchEmails(app.defaultAccountId, debouncedQuery.trim(), nextPage, LIMIT);
				if (reqId === loadSeq) {
					emails = [...emails, ...result.emails];
					totalEmails = result.total;
					currentPage = nextPage;
				}
			} catch {}
			loadingMore = false;
			return;
		}

		const cached = showUnreadOnly ? null : mailCache.get(app.defaultAccountId, folder, nextPage);
		if (cached) {
			emails = [...emails, ...cached.emails];
			totalEmails = cached.total;
			currentPage = nextPage;
			loadingMore = false;
			return;
		}

		try {
			const result = await backend.getEmailsByFolder(app.defaultAccountId, folder, nextPage, LIMIT, showUnreadOnly);
			if (reqId === loadSeq) {
				emails = [...emails, ...result.emails];
				totalEmails = result.total;
				currentPage = nextPage;
				if (!showUnreadOnly) {
					mailCache.set(app.defaultAccountId, folder, nextPage, result.emails, result.total);
				}
			}
		} catch {}
		loadingMore = false;
	}

	async function syncAndLoad() {
		if (!app.defaultAccountId) return;
		syncing = true;
		try {
			await backend.syncAccount(app.defaultAccountId);
			await app.refreshAccounts();
			const target = app.folders.find((f) => f.type === folder);
			if (target) {
				await backend.syncFolder(app.defaultAccountId, target.id);
			}
			mailCache.invalidateFolder(app.defaultAccountId, folder);
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

	function formatDate(dateStr: string) {
		const date = new Date(dateStr);
		const now = new Date();
		const diffMs = now.getTime() - date.getTime();
		const diffDays = diffMs / (1000 * 60 * 60 * 24);
		if (diffDays < 1) return date.toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' });
		if (diffDays < 7) return date.toLocaleDateString('fr-FR', { weekday: 'short' });
		return date.toLocaleDateString('fr-FR', { day: 'numeric', month: 'short' });
	}

	function formatFileSize(bytes: number): string {
		if (bytes < 1024) return `${bytes} B`;
		if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
		return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
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
			toast.success(isTrash ? `${count} permanently deleted` : `${count} moved to trash`);
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
			if (showUnreadOnly) {
				emails = emails.filter((e) => !ids.includes(e.id));
				totalEmails = Math.max(0, totalEmails - ids.length);
				if (selectedId && ids.includes(selectedId)) selectedId = null;
			} else {
				emails = emails.map((e) => (ids.includes(e.id) ? { ...e, is_read: true } : e));
			}
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
			if (showUnreadOnly && !email.is_read) {
				emails = emails.filter((e) => e.id !== email.id);
				totalEmails = Math.max(0, totalEmails - 1);
				if (selectedId === email.id) selectedId = null;
			} else {
				emails = emails.map((e) => (e.id === email.id ? { ...e, is_read: !email.is_read } : e));
			}
		} catch {}
	}

	async function handleToggleStar(email: EmailMessage) {
		if (!app.defaultAccountId) return;
		try {
			await backend.updateEmail(app.defaultAccountId, email.id, { is_starred: !email.is_starred });
			emails = emails.map((e) => (e.id === email.id ? { ...e, is_starred: !email.is_starred } : e));
		} catch {}
	}

	let syncedAccount = -1;
	$effect(() => {
		const accountId = app.defaultAccountId;
		if (!accountId || accountId === syncedAccount) return;
		syncedAccount = accountId;
		void syncAndLoad();
	});

	$effect(() => {
		const _folder = folder;
		const _unreadOnly = showUnreadOnly;
		const _query = debouncedQuery.trim();
		if (!app.defaultAccountId) return;
		const accountId = app.defaultAccountId;
		const reqId = ++loadSeq;

		checkedIds = new Set();
		selectedId = null;
		currentPage = 1;
		totalEmails = 0;

		if (_query !== '') {
			emails = [];
			loading = true;
			(async () => {
				try {
					const result = await backend.searchEmails(accountId, _query, 1, LIMIT);
					if (reqId !== loadSeq) return;
					emails = result.emails;
					totalEmails = result.total;
				} catch {
					if (reqId !== loadSeq) return;
					emails = [];
					totalEmails = 0;
				}
				if (reqId !== loadSeq) return;
				loading = false;
			})();
			return;
		}

		const cached = _unreadOnly ? null : mailCache.get(accountId, _folder, 1);
		if (cached) {
			emails = cached.emails;
			totalEmails = cached.total;
			loading = false;
		} else {
			emails = [];
			loading = true;
		}

		(async () => {
			try {
				const result = await backend.getEmailsByFolder(accountId, _folder, 1, LIMIT, _unreadOnly);
				if (reqId !== loadSeq) return;
				emails = result.emails;
				totalEmails = result.total;
				if (!_unreadOnly) {
					mailCache.set(accountId, _folder, 1, result.emails, result.total);
				}
			} catch {
				if (reqId !== loadSeq) return;
				if (!cached) {
					emails = [];
					totalEmails = 0;
				}
			}
			if (reqId !== loadSeq) return;
			loading = false;
		})();
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

<DeleteConfirmDialog bind:open={deleteDialogOpen} count={deleteTarget.length} permanent={isTrash} onconfirm={confirmDelete} />

{#snippet listPane()}
	<div class="flex h-full flex-col">
		<div class="border-b">
			<div class="px-4 pt-3 pb-2">
				<div class="relative">
					<Search class="pointer-events-none absolute left-2.5 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
					<Input
						bind:ref={inputEl}
						bind:value={searchQuery}
						onkeydown={onSearchKeydown}
						placeholder="Search all mail…"
						aria-label="Search all mail"
						aria-keyshortcuts="Meta+K Control+K"
						class="h-9 pl-8 pr-12"
					/>
					{#if searchQuery}
						<button
							type="button"
							aria-label="Clear search"
							class="absolute right-2 top-1/2 flex h-6 w-6 -translate-y-1/2 items-center justify-center rounded text-muted-foreground hover:bg-muted hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
							onclick={clearSearch}
						>
							<X class="h-3.5 w-3.5" />
						</button>
					{:else}
						<kbd
							class="pointer-events-none absolute right-2 top-1/2 hidden -translate-y-1/2 items-center rounded border border-border bg-muted px-1.5 py-0.5 text-[10px] font-medium text-muted-foreground md:flex"
						>⌘K</kbd>
					{/if}
				</div>
			</div>
			<div class="flex items-center justify-between px-4 pb-3">
				<div class="flex items-center gap-2">
					{#if emails.length > 0}
						<button
							aria-label={allChecked ? 'Deselect all' : 'Select all'}
							class="flex h-6 w-6 items-center justify-center rounded hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
							onclick={toggleAll}
						>
							<Checkbox checked={allChecked} class="h-4 w-4" />
						</button>
					{/if}
					{#if searchMode}
						<h2 class="text-lg font-semibold">Search</h2>
						<span class="text-sm text-muted-foreground">{totalEmails} result{totalEmails === 1 ? '' : 's'}</span>
					{:else}
						<h2 class="hidden text-lg font-semibold md:block">{folderLabel}</h2>
						<div class="md:hidden">
							<MailFolderSwitcher folders={app.folders} />
						</div>
					{/if}
				</div>
				<div class="flex items-center gap-1">
					<Button
						variant="ghost"
						size="sm"
						class="h-8 gap-1.5 px-2 text-xs {showUnreadOnly ? 'bg-primary/10 text-primary font-medium' : 'text-muted-foreground'}"
						onclick={() => { showUnreadOnly = !showUnreadOnly; }}
					>
						<Mail class="h-3.5 w-3.5" />
						Unread
					</Button>
					<Button variant="ghost" size="icon" aria-label="Refresh" class="h-8 w-8" onclick={syncAndLoad} disabled={syncing}>
						<RefreshCw class="h-4 w-4 {syncing ? 'animate-spin' : ''}" />
					</Button>
				</div>
			</div>
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
					{#if searchMode}
						<p class="text-sm">No messages match “{debouncedQuery.trim()}”</p>
					{:else}
						<p class="text-sm">{showUnreadOnly ? `No unread emails in ${folderLabel}` : `No emails in ${folderLabel}`}</p>
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
							onarchive={folder !== 'archive' ? () => handleSingleArchive(email.id) : undefined}
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
						<Button variant="ghost" size="icon" aria-label="Reply" class="h-7 w-7" onclick={() => goto(`/mail/compose?reply=${selected!.id}&accountId=${app.defaultAccountId}`)}>
							<Reply class="h-4 w-4" />
						</Button>
						<Button variant="ghost" size="icon" aria-label="Reply all" class="h-7 w-7" onclick={() => goto(`/mail/compose?replyall=${selected!.id}&accountId=${app.defaultAccountId}`)}>
							<ReplyAll class="h-4 w-4" />
						</Button>
						<Button variant="ghost" size="icon" aria-label="Forward" class="h-7 w-7" onclick={() => goto(`/mail/compose?forward=${selected!.id}&accountId=${app.defaultAccountId}`)}>
							<Forward class="h-4 w-4" />
						</Button>
						<div class="mx-1 h-4 w-px bg-border"></div>
						{#if folder !== 'archive'}
							<Button variant="ghost" size="icon" aria-label="Archive" class="h-7 w-7" onclick={() => handleSingleArchive(selected!.id)}>
								<Archive class="h-4 w-4" />
							</Button>
						{/if}
						<Button variant="ghost" size="icon" aria-label="Delete" class="h-7 w-7 text-destructive hover:text-destructive hover:bg-destructive/10" onclick={() => handleSingleDelete(selected!.id)}>
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
					{@html sanitizedBody}
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
