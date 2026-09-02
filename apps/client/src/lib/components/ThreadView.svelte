<script lang="ts">
	import DOMPurify from 'dompurify';
	import { backend, type EmailMessage, type EmailAttachment } from '$lib/backend';
	import { Avatar, Button, Spinner, icons } from '@facile/muse';

	let {
		accountId,
		email,
		onthread
	}: {
		accountId: number;
		email: EmailMessage;
		onthread?: (messages: EmailMessage[]) => void;
	} = $props();

	/* Glyphs an email client needs that muse's `icons` map has no key for yet. */
	const mailIcons = {
		paperclip: 'solar:paperclip-linear'
	};

	let messages = $state<EmailMessage[]>([]);
	let expanded = $state<Set<number>>(new Set());
	let loadingBodies = $state<Set<number>>(new Set());
	let loadingThread = $state(false);
	let loadSeq = 0;

	function resolveCIDImages(html: string, emailId: number): string {
		return html.replace(/src=["']cid:([^"']+)["']/gi, (_match, cid) => {
			return `src="${backend.getCIDImageUrl(accountId, emailId, cid)}"`;
		});
	}

	function sanitize(msg: EmailMessage): string {
		if (!msg.body_html) return '';
		const resolved = resolveCIDImages(msg.body_html, msg.id);
		return DOMPurify.sanitize(resolved).replace(/<img /gi, '<img loading="lazy" decoding="async" ');
	}

	function snippet(msg: EmailMessage): string {
		const text = (msg.body_text || msg.body_html?.replace(/<[^>]+>/g, ' ') || '').replace(/\s+/g, ' ').trim();
		return text.slice(0, 120);
	}

	function formatDate(dateStr: string) {
		const date = new Date(dateStr);
		const now = new Date();
		const diffDays = (now.getTime() - date.getTime()) / (1000 * 60 * 60 * 24);
		if (diffDays < 1) return date.toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' });
		if (diffDays < 7) return date.toLocaleDateString('fr-FR', { weekday: 'short', hour: '2-digit', minute: '2-digit' });
		return date.toLocaleDateString('fr-FR', { day: 'numeric', month: 'short', year: 'numeric' });
	}

	function formatFileSize(bytes: number): string {
		if (bytes < 1024) return `${bytes} B`;
		if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
		return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
	}

	async function ensureBody(msg: EmailMessage) {
		if (msg.body_text || msg.body_html || loadingBodies.has(msg.id)) return;
		loadingBodies = new Set(loadingBodies).add(msg.id);
		try {
			const full = await backend.getEmail(accountId, msg.id);
			messages = messages.map((m) => (m.id === msg.id ? full : m));
		} catch {
			/* leave body empty; UI shows a fallback */
		}
		const next = new Set(loadingBodies);
		next.delete(msg.id);
		loadingBodies = next;
	}

	async function markRead(msg: EmailMessage) {
		if (msg.is_read) return;
		try {
			await backend.updateEmail(accountId, msg.id, { is_read: true });
			messages = messages.map((m) => (m.id === msg.id ? { ...m, is_read: true } : m));
		} catch {
			/* non-fatal */
		}
	}

	function toggle(msg: EmailMessage) {
		const next = new Set(expanded);
		if (next.has(msg.id)) {
			next.delete(msg.id);
		} else {
			next.add(msg.id);
			void ensureBody(msg);
			void markRead(msg);
		}
		expanded = next;
	}

	async function downloadAttachment(msg: EmailMessage, attachment: EmailAttachment) {
		await backend.downloadAttachment(accountId, msg.id, attachment.id, attachment.filename);
	}

	let loadedId = -1;

	$effect(() => {
		const seedId = email.id;
		const threadId = email.thread_id;
		// Only (re)load when the opened message actually changes. The parent passes
		// a fresh object for the open row when it marks it read, etc.; re-running
		// here would refetch the thread and discard expanded/loaded messages.
		if (seedId === loadedId) return;
		loadedId = seedId;
		const reqId = ++loadSeq;

		// Always render the opened message immediately; enrich with the thread next.
		messages = [email];
		expanded = new Set([seedId]);
		void markRead(email);

		if (!threadId) {
			onthread?.([email]);
			return;
		}

		loadingThread = true;
		(async () => {
			try {
				const result = await backend.getThread(accountId, threadId);
				if (reqId !== loadSeq) return;
				// Use the thread rows as-is: they carry real per-message read state and
				// bodies; the list representative omits bodies and holds aggregate state.
				messages = result.emails.length > 0 ? result.emails : [email];
				// Expand the opened message; if it is not the newest, expand the newest too.
				const newest = messages[messages.length - 1];
				const next = new Set<number>([seedId]);
				if (newest) {
					next.add(newest.id);
					void ensureBody(newest);
					void markRead(newest);
				}
				expanded = next;
				onthread?.(messages);
			} catch {
				if (reqId !== loadSeq) return;
				messages = [email];
				onthread?.([email]);
			}
			if (reqId === loadSeq) loadingThread = false;
		})();
	});
</script>

<div class="flex flex-col">
	{#each messages as msg (msg.id)}
		{@const isOpen = expanded.has(msg.id)}
		<div class="border-b border-fc-border last:border-b-0">
			<button
				type="button"
				class="flex w-full items-center gap-3 px-4 py-3 text-left text-fc-fg transition-colors hover:bg-fc-surface focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-fc-ring sm:px-6"
				onclick={() => toggle(msg)}
				aria-expanded={isOpen}
			>
				<span class="shrink-0">
					<Avatar name={msg.from_name || msg.from_address || '?'} size="sm" />
				</span>
				<span class="min-w-0 flex-1">
					<span class="flex items-center justify-between gap-2">
						<span class="truncate text-fc-sm {msg.is_read ? '' : 'font-semibold'}">
							{msg.from_name || msg.from_address}
						</span>
						<span class="flex shrink-0 items-center gap-2 text-fc-xs text-fc-fg-muted">
							{#if !msg.is_read}
								<span class="size-2 rounded-fc-pill bg-fc-accent"></span>
							{/if}
							{#if msg.has_attachments}
								<iconify-icon
									icon={mailIcons.paperclip}
									width="14"
									height="14"
									class="block size-3.5"
									aria-label="Has attachments"
								></iconify-icon>
							{/if}
							<span>{formatDate(msg.date)}</span>
							<iconify-icon
								icon={icons.chevronDown}
								width="16"
								height="16"
								class="block size-4 transition-transform {isOpen ? 'rotate-180' : ''}"
							></iconify-icon>
						</span>
					</span>
					{#if !isOpen}
						<span class="mt-0.5 block truncate text-fc-xs text-fc-fg-muted">
							{snippet(msg) || msg.from_address}
						</span>
					{:else}
						<span class="mt-0.5 block truncate text-fc-xs text-fc-fg-muted">
							to {msg.to_addresses.map((a) => a.name || a.email).join(', ')}
						</span>
					{/if}
				</span>
			</button>

			{#if isOpen}
				<div class="px-4 pb-5 sm:px-6">
					{#if msg.attachments && msg.attachments.length > 0}
						<div class="mb-4 flex flex-wrap items-center gap-2">
							<span class="flex items-center gap-1.5 text-fc-xs text-fc-fg-muted">
								<iconify-icon
									icon={mailIcons.paperclip}
									width="14"
									height="14"
									class="block size-3.5"
								></iconify-icon>
								{msg.attachments.length} attachment{msg.attachments.length > 1 ? 's' : ''}
							</span>
							{#each msg.attachments as attachment (attachment.id)}
								<Button
									variant="outline"
									size="sm"
									icon={icons.download}
									onclick={() => downloadAttachment(msg, attachment)}
								>
									<span class="max-w-48 truncate">{attachment.filename}</span>
									<span class="text-fc-fg-muted">({formatFileSize(attachment.size)})</span>
								</Button>
							{/each}
						</div>
					{/if}

					{#if msg.body_html}
						<div class="mail-body-content" style="color: var(--color-fc-fg);">{@html sanitize(msg)}</div>
					{:else if msg.body_text}
						<pre class="whitespace-pre-wrap font-fc-body text-fc-sm text-fc-fg">{msg.body_text}</pre>
					{:else if loadingBodies.has(msg.id)}
						<div class="flex items-center gap-3 text-fc-sm text-fc-fg-muted">
							<Spinner size="sm" />
							<span>Loading message…</span>
						</div>
					{:else}
						<p class="text-fc-sm text-fc-fg-muted">(empty message)</p>
					{/if}
				</div>
			{/if}
		</div>
	{/each}

	{#if loadingThread && messages.length <= 1}
		<div class="flex items-center gap-3 px-4 py-3 text-fc-xs text-fc-fg-muted sm:px-6">
			<Spinner size="sm" label="Loading conversation" />
			<span>Loading conversation…</span>
		</div>
	{/if}
</div>
