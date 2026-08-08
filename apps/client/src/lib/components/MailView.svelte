<script lang="ts">
	import { goto } from '$app/navigation';
	import { getContext } from 'svelte';
	import type { Attachment } from 'svelte/attachments';
	import { MediaQuery } from 'svelte/reactivity';
	import { backend, type EmailMessage, type MailAccount } from '$lib/backend';
	import { mailCache } from '$lib/stores/mail-cache';
	import { searchStore } from '$lib/stores/search.svelte';
	import {
		Badge,
		Button,
		Checkbox,
		ConfirmModal,
		EmptyState,
		IconButton,
		Input,
		Skeleton,
		Spinner,
		icons,
		toast
	} from '@facile/muse';
	import * as Resizable from '$lib/components/ui/resizable';
	import EmailItem from '$lib/components/EmailItem.svelte';
	import ThreadView from '$lib/components/ThreadView.svelte';
	import BulkActionBar from '$lib/components/BulkActionBar.svelte';
	import MailFolderSwitcher from '$lib/components/MailFolderSwitcher.svelte';

	let { folder }: { folder: string } = $props();

	/* Glyphs an email client needs that muse's `icons` map has no key for yet. */
	const mailIcons = {
		reply: 'solar:reply-linear',
		replyAll: 'solar:multiple-forward-left-linear',
		forward: 'solar:forward-linear',
		archive: 'solar:archive-linear'
	};

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
	/* muse's `Input` renders the <input> itself and exposes no element binding, so the node
	   comes back through an attachment — the same trick muse's own ConfirmModal uses. */
	const captureInput: Attachment<HTMLInputElement> = (node) => {
		inputEl = node;
		return () => {
			inputEl = null;
		};
	};
	let loadSeq = 0;
	let seenFocusSeq = searchStore.focusSeq;

	let threadMessages = $state<EmailMessage[]>([]);

	const selected = $derived(emails.find((e) => e.id === selectedId) ?? null);
	const replyTarget = $derived(threadMessages.length > 0 ? threadMessages[threadMessages.length - 1] : selected);
	const hasMore = $derived(emails.length < totalEmails);
	const selectionActive = $derived(checkedIds.size > 0);
	const allChecked = $derived(emails.length > 0 && emails.every((e) => checkedIds.has(e.id)));
	const searchMode = $derived(debouncedQuery.trim() !== '');
	const deleteCount = $derived(deleteTarget.length);
	const LIMIT = 30;

	let animatedKeys = new Set<string>();

	function animKey(e: EmailMessage): string {
		return e.message_id || `id:${e.id}`;
	}

	function isNew(e: EmailMessage): boolean {
		return !animatedKeys.has(animKey(e));
	}

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

	async function loadEmails() {
		if (!app.defaultAccountId) return;

		const cached = showUnreadOnly ? null : mailCache.get(app.defaultAccountId, folder, 1);
		if (cached) {
			emails = cached.emails;
			totalEmails = cached.total;
		}

		loading = !cached && emails.length === 0;
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

	// A list row stands for a whole conversation; map the checked/targeted
	// representative ids back to every underlying message id for server actions.
	function expandIds(repIds: number[]): number[] {
		const byId = new Map(emails.map((e) => [e.id, e]));
		const out: number[] = [];
		for (const id of repIds) {
			out.push(...(byId.get(id)?.email_ids ?? [id]));
		}
		return out;
	}

	async function openEmail(email: EmailMessage) {
		selectedId = email.id;
		threadMessages = [];
		if (!app.defaultAccountId) return;

		const unread = (email.unread_count ?? (email.is_read ? 0 : 1)) > 0;
		if (unread) {
			try {
				await backend.bulkAction(app.defaultAccountId, email.email_ids ?? [email.id], 'mark_read');
				emails = emails.map((e) => (e.id === email.id ? { ...e, is_read: true, unread_count: 0 } : e));
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

	async function handleBulkDelete() {
		deleteTarget = [...checkedIds];
		deleteDialogOpen = true;
	}

	async function confirmDelete() {
		if (!app.defaultAccountId || deleteTarget.length === 0) return;
		bulkLoading = true;
		const count = deleteTarget.length;
		try {
			await backend.bulkAction(app.defaultAccountId, expandIds(deleteTarget), 'delete');
			emails = emails.filter((e) => !deleteTarget.includes(e.id));
			totalEmails = Math.max(0, totalEmails - count);
			mailCache.removeEmails(deleteTarget);
			if (selectedId && deleteTarget.includes(selectedId)) selectedId = null;
			checkedIds = new Set();
			deleteTarget = [];
			toast.success(isTrash ? `${count} permanently deleted` : `${count} moved to trash`);
		} catch {
			toast.danger('Failed to delete');
		}
		bulkLoading = false;
	}

	async function handleBulkArchive() {
		if (!app.defaultAccountId) return;
		bulkLoading = true;
		const ids = [...checkedIds];
		try {
			await backend.bulkAction(app.defaultAccountId, expandIds(ids), 'archive');
			emails = emails.filter((e) => !ids.includes(e.id));
			totalEmails = Math.max(0, totalEmails - ids.length);
			mailCache.removeEmails(ids);
			if (selectedId && ids.includes(selectedId)) selectedId = null;
			checkedIds = new Set();
			toast.success(`${ids.length} archived`);
		} catch {
			toast.danger('Failed to archive');
		}
		bulkLoading = false;
	}

	async function handleBulkMarkRead() {
		if (!app.defaultAccountId) return;
		bulkLoading = true;
		const ids = [...checkedIds];
		try {
			await backend.bulkAction(app.defaultAccountId, expandIds(ids), 'mark_read');
			if (showUnreadOnly) {
				emails = emails.filter((e) => !ids.includes(e.id));
				totalEmails = Math.max(0, totalEmails - ids.length);
				if (selectedId && ids.includes(selectedId)) selectedId = null;
			} else {
				emails = emails.map((e) => (ids.includes(e.id) ? { ...e, is_read: true, unread_count: 0 } : e));
			}
			checkedIds = new Set();
		} catch {
			toast.danger('Failed to mark as read');
		}
		bulkLoading = false;
	}

	async function handleBulkMarkUnread() {
		if (!app.defaultAccountId) return;
		bulkLoading = true;
		const ids = [...checkedIds];
		try {
			await backend.bulkAction(app.defaultAccountId, expandIds(ids), 'mark_unread');
			emails = emails.map((e) => (ids.includes(e.id) ? { ...e, is_read: false, unread_count: e.message_count ?? 1 } : e));
			checkedIds = new Set();
		} catch {
			toast.danger('Failed to mark as unread');
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
			await backend.bulkAction(app.defaultAccountId, expandIds([emailId]), 'archive');
			emails = emails.filter((e) => e.id !== emailId);
			totalEmails = Math.max(0, totalEmails - 1);
			mailCache.removeEmails([emailId]);
			if (selectedId === emailId) selectedId = null;
			toast.success('Archived');
		} catch {
			toast.danger('Failed to archive');
		}
	}

	async function handleToggleRead(email: EmailMessage) {
		if (!app.defaultAccountId) return;
		const markRead = !email.is_read;
		try {
			await backend.bulkAction(app.defaultAccountId, email.email_ids ?? [email.id], markRead ? 'mark_read' : 'mark_unread');
			if (showUnreadOnly && markRead) {
				emails = emails.filter((e) => e.id !== email.id);
				totalEmails = Math.max(0, totalEmails - 1);
				if (selectedId === email.id) selectedId = null;
			} else {
				emails = emails.map((e) => (e.id === email.id ? { ...e, is_read: markRead, unread_count: markRead ? 0 : (e.message_count ?? 1) } : e));
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
		const accountId = app.defaultAccountId;
		const reqId = ++loadSeq;

		checkedIds = new Set();
		selectedId = null;
		currentPage = 1;
		totalEmails = 0;
		animatedKeys = new Set();

		if (!accountId) {
			emails = [];
			loading = false;
			return;
		}

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
		for (const e of emails) animatedKeys.add(animKey(e));
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

<ConfirmModal
	bind:open={deleteDialogOpen}
	tone="danger"
	icon={icons.remove}
	title={isTrash
		? `Permanently delete ${deleteCount === 1 ? 'this email' : `${deleteCount} emails`}?`
		: `Delete ${deleteCount === 1 ? 'this email' : `${deleteCount} emails`}?`}
	description={isTrash
		? `This cannot be undone. ${deleteCount === 1 ? 'This email' : 'These emails'} will be permanently removed.`
		: `${deleteCount === 1 ? 'This email' : 'These emails'} will be moved to Trash. You can restore ${deleteCount === 1 ? 'it' : 'them'} from there.`}
	confirmLabel="Delete"
	onConfirm={confirmDelete}
/>

{#snippet listPane()}
	<div class="flex h-full flex-col">
		<div class="flex flex-col gap-3 border-b border-fc-border px-4 py-3">
			<div class="relative">
				<iconify-icon
					icon={icons.search}
					width="16"
					height="16"
					class="pointer-events-none absolute left-3 top-1/2 block size-4 -translate-y-1/2 text-fc-fg-muted"
				></iconify-icon>
				<Input
					{@attach captureInput}
					bind:value={searchQuery}
					onkeydown={onSearchKeydown}
					placeholder="Search all mail…"
					aria-label="Search all mail"
					aria-keyshortcuts="Meta+K Control+K"
					class="pl-9 pr-14 text-fc-sm"
				/>
				{#if searchQuery}
					<IconButton
						variant="ghost"
						aria-label="Clear search"
						class="absolute right-0 top-0"
						onclick={clearSearch}
					>
						<iconify-icon icon={icons.close} width="16" height="16" class="block size-4"
						></iconify-icon>
					</IconButton>
				{:else}
					<kbd
						class="pointer-events-none absolute right-3 top-1/2 hidden -translate-y-1/2 items-center rounded-fc-xs bg-fc-surface px-1.5 py-0.5 text-fc-xs font-medium text-fc-fg-muted md:flex"
						>⌘K</kbd
					>
				{/if}
			</div>
			<div class="flex items-center justify-between gap-2">
				<div class="flex min-w-0 items-center gap-2">
					{#if emails.length > 0}
						<Checkbox
							checked={allChecked}
							aria-label={allChecked ? 'Deselect all' : 'Select all'}
							onchange={toggleAll}
							class="size-8 shrink-0 justify-center"
						/>
					{/if}
					{#if searchMode}
						<h2 class="truncate text-fc-lg font-semibold text-fc-fg">Search</h2>
						<span class="shrink-0 text-fc-sm text-fc-fg-muted">
							{totalEmails} result{totalEmails === 1 ? '' : 's'}
						</span>
					{:else}
						<h2 class="hidden truncate text-fc-lg font-semibold text-fc-fg md:block">{folderLabel}</h2>
						<div class="min-w-0 md:hidden">
							<MailFolderSwitcher folders={app.folders} />
						</div>
					{/if}
				</div>
				<div class="flex shrink-0 items-center gap-1">
					<Button
						size="sm"
						variant={showUnreadOnly ? 'primary' : 'ghost'}
						icon={icons.mail}
						aria-pressed={showUnreadOnly}
						onclick={() => {
							showUnreadOnly = !showUnreadOnly;
						}}
					>
						Unread
					</Button>
					<IconButton
						variant="ghost"
						aria-label="Refresh"
						title="Refresh"
						disabled={syncing}
						onclick={syncAndLoad}
					>
						<iconify-icon
							icon={icons.refresh}
							width="18"
							height="18"
							class="block size-4.5 {syncing ? 'animate-spin motion-reduce:animate-none' : ''}"
						></iconify-icon>
					</IconButton>
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
				<!-- Skeleton rows, not a centred spinner: same padding, same avatar size and the
				     same two text lines as a real row, so nothing moves when the data lands. -->
				{#each Array(8) as _, i (i)}
					<div
						class="mail-skeleton flex items-center gap-3 border-b border-fc-border px-4 py-2.5"
						style="--delay: {i * 60}ms"
					>
						<Skeleton class="size-8 shrink-0 rounded-fc-pill" />
						<div class="flex min-w-0 flex-1 flex-col gap-1.5">
							<div class="flex items-center justify-between gap-2">
								<Skeleton class="h-4 w-28" />
								<Skeleton class="h-3 w-10" />
							</div>
							<Skeleton class="h-3.5 w-48" />
						</div>
					</div>
				{/each}
			{:else if emails.length === 0}
				<div class="mail-fade-in flex h-full items-center justify-center px-4">
					{#if searchMode}
						<EmptyState
							bare
							icon={icons.search}
							title="No messages match “{debouncedQuery.trim()}”"
							description="Try a different sender, subject or word from the body."
						/>
					{:else}
						<EmptyState
							bare
							icon={icons.mail}
							title={showUnreadOnly ? `No unread email in ${folderLabel}` : `No email in ${folderLabel}`}
							description={showUnreadOnly
								? 'Everything here has been read.'
								: 'New messages appear here after a sync.'}
						/>
					{/if}
				</div>
			{:else}
				{#each emails as email, i (email.id)}
					<div style="--delay: {Math.min(i, 15) * 30}ms" class="mail-list-item" class:mail-animate={isNew(email)}>
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
						<Spinner size="sm" label="Loading more messages" />
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
			<div class="mail-detail-header flex flex-col gap-2 border-b border-fc-border px-4 py-4 sm:px-6">
				<div class="flex items-center gap-2">
					<h1 class="text-fc-xl font-semibold text-fc-fg">{selected.subject || '(no subject)'}</h1>
					{#if threadMessages.length > 1}
						<Badge tone="neutral" class="tabular-nums">{threadMessages.length}</Badge>
					{/if}
				</div>
				<div class="flex flex-wrap items-center gap-x-3 gap-y-1 text-fc-sm text-fc-fg-muted">
					<span class="font-medium text-fc-fg">{selected.from_name || selected.from_address}</span>
					<span>&lt;{selected.from_address}&gt;</span>
					<div class="ml-auto flex items-center gap-1">
						<IconButton
							variant="ghost"
							aria-label="Reply"
							title="Reply"
							onclick={() => goto(`/mail/compose?reply=${replyTarget!.id}&accountId=${app.defaultAccountId}`)}
						>
							<iconify-icon icon={mailIcons.reply} width="18" height="18" class="block size-4.5"
							></iconify-icon>
						</IconButton>
						<IconButton
							variant="ghost"
							aria-label="Reply all"
							title="Reply all"
							onclick={() => goto(`/mail/compose?replyall=${replyTarget!.id}&accountId=${app.defaultAccountId}`)}
						>
							<iconify-icon icon={mailIcons.replyAll} width="18" height="18" class="block size-4.5"
							></iconify-icon>
						</IconButton>
						<IconButton
							variant="ghost"
							aria-label="Forward"
							title="Forward"
							onclick={() => goto(`/mail/compose?forward=${replyTarget!.id}&accountId=${app.defaultAccountId}`)}
						>
							<iconify-icon icon={mailIcons.forward} width="18" height="18" class="block size-4.5"
							></iconify-icon>
						</IconButton>
						<span class="mx-1 h-5 w-px shrink-0 bg-fc-border"></span>
						{#if folder !== 'archive'}
							<IconButton
								variant="ghost"
								aria-label="Archive"
								title="Archive"
								onclick={() => handleSingleArchive(selected!.id)}
							>
								<iconify-icon icon={mailIcons.archive} width="18" height="18" class="block size-4.5"
								></iconify-icon>
							</IconButton>
						{/if}
						<IconButton
							variant="danger"
							aria-label="Delete"
							title="Delete"
							onclick={() => handleSingleDelete(selected!.id)}
						>
							<iconify-icon icon={icons.remove} width="18" height="18" class="block size-4.5"
							></iconify-icon>
						</IconButton>
						<span class="ml-2 shrink-0 text-fc-xs">{formatDate(selected.date)}</span>
					</div>
				</div>
			</div>
			<div class="flex-1 overflow-auto">
				{#key selected.id}
					<ThreadView accountId={app.defaultAccountId!} email={selected} onthread={(m) => (threadMessages = m)} />
				{/key}
			</div>
		{:else}
			<div class="flex flex-1 items-center justify-center px-4">
				<EmptyState
					bare
					icon={icons.mail}
					title="Nothing open"
					description="Pick a message on the left to read it here."
				/>
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
		<div class="flex flex-shrink-0 items-center border-b border-fc-border px-2 py-1.5">
			<Button variant="ghost" size="sm" icon={icons.chevronLeft} onclick={() => (selectedId = null)}>
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
		.mail-list-item.mail-animate {
			animation: mail-fade-in 200ms ease-out both;
			animation-delay: var(--delay, 0ms);
		}

		.mail-detail-header {
			animation: mail-slide-in 180ms ease-out both;
		}

		.mail-fade-in {
			animation: mail-fade-in 200ms ease-out both;
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
</style>
