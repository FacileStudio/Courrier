<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { getContext, onMount, onDestroy } from 'svelte';
	import { backend, type MailAccount } from '$lib/backend';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Badge } from '$lib/components/ui/badge';
	import { toast } from 'svelte-sonner';
	import { X, SendHorizonal, Paperclip, FileText, Type } from 'lucide-svelte';
	import TiptapEditor from '$lib/components/TiptapEditor.svelte';
	import AddressInput from '$lib/components/AddressInput.svelte';

	const app = getContext<{
		token: string;
		defaultAccountId: number | null;
		accounts: MailAccount[];
	}>('app');

	let to = $state('');
	let cc = $state('');
	let subject = $state('');
	let bodyHtml = $state('');
	let bodyPlainText = $state('');
	let initialContent = $state('');
	let sending = $state(false);
	let inReplyTo = $state('');
	let references = $state<string[]>([]);
	let ready = $state(false);

	let attachedFiles = $state<File[]>([]);
	let fileInput: HTMLInputElement;

	let plainTextMode = $state(false);

	let draftId = $state<number | null>(null);
	let lastSaved = $state<string | null>(null);
	let savingDraft = $state(false);
	let saveTimer: ReturnType<typeof setTimeout>;

	function buildSignatureHtml(signature: string): string {
		const sigLines = signature
			.split('\n')
			.map((l) => `<p>${l || '<br>'}</p>`)
			.join('');
		return `<p><br></p><p>--</p>${sigLines}`;
	}

	function formatFileSize(bytes: number): string {
		if (bytes < 1024) return `${bytes} B`;
		if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
		return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
	}

	function handleFileSelect(e: Event) {
		const input = e.target as HTMLInputElement;
		if (!input.files) return;
		attachedFiles = [...attachedFiles, ...Array.from(input.files)];
		input.value = '';
	}

	function removeFile(index: number) {
		attachedFiles = attachedFiles.filter((_, i) => i !== index);
	}

	function togglePlainText() {
		if (plainTextMode) {
			const paragraphs = bodyPlainText.split('\n').map((line) => `<p>${line || '<br>'}</p>`).join('');
			bodyHtml = paragraphs;
			initialContent = paragraphs;
		} else {
			bodyPlainText = bodyHtml.replace(/<[^>]*>/g, '');
		}
		plainTextMode = !plainTextMode;
	}

	function scheduleDraftSave() {
		clearTimeout(saveTimer);
		saveTimer = setTimeout(saveDraft, 3000);
	}

	async function saveDraft() {
		if (!app.defaultAccountId || savingDraft) return;
		if (!to.trim() && !subject.trim() && !bodyHtml && !bodyPlainText) return;

		savingDraft = true;
		try {
			const toAddrs = to.split(',').map((s) => s.trim()).filter(Boolean);
			const ccAddrs = cc ? cc.split(',').map((s) => s.trim()).filter(Boolean) : [];
			const body = plainTextMode ? bodyPlainText : bodyHtml.replace(/<[^>]*>/g, '');
			const bodyHtmlValue = plainTextMode ? undefined : bodyHtml;

			if (draftId) {
				await backend.deleteDraft(app.token, app.defaultAccountId, draftId);
			}

			const result = await backend.saveDraft(app.token, app.defaultAccountId, {
				to: toAddrs,
				cc: ccAddrs,
				subject,
				body,
				body_html: bodyHtmlValue,
				in_reply_to: inReplyTo || undefined,
				references: references.length > 0 ? references : undefined
			});
			draftId = result.id;
			lastSaved = new Date().toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' });
		} catch {
			// silent fail
		}
		savingDraft = false;
	}

	onMount(async () => {
		const params = page.url.searchParams;
		const replyId = params.get('reply');
		const replyAllId = params.get('replyall');
		const forwardId = params.get('forward');
		const accountId = params.get('accountId');
		const emailId = replyId || replyAllId || forwardId;
		const loadDraftId = params.get('draft');

		const account = app.accounts.find((a) => a.id === app.defaultAccountId);
		let signatureHtml = '';
		if (account?.signature) {
			signatureHtml = buildSignatureHtml(account.signature);
		}

		if (loadDraftId && app.defaultAccountId) {
			try {
				const draft = await backend.getEmail(app.token, app.defaultAccountId, parseInt(loadDraftId));
				to = draft.to_addresses.map((a) => a.email).join(', ');
				cc = draft.cc_addresses.map((a) => a.email).join(', ');
				subject = draft.subject;
				bodyHtml = draft.body_html || '';
				bodyPlainText = draft.body_text || '';
				initialContent = draft.body_html || `<p>${draft.body_text || ''}</p>`;
				inReplyTo = draft.in_reply_to || '';
				references = draft.references ? draft.references.split(/\s+/).filter(Boolean) : [];
				draftId = parseInt(loadDraftId);
			} catch {
				toast.error('Failed to load draft');
			}
		} else if (emailId && accountId) {
			try {
				const original = await backend.getEmail(app.token, parseInt(accountId), parseInt(emailId));
				const acct = app.accounts.find((a) => a.id === parseInt(accountId));
				const currentEmail = acct?.email ?? '';
				const senderDate = new Date(original.date).toLocaleString('fr-FR');
				const senderName = original.from_name || original.from_address;

				if (replyId || replyAllId) {
					to = original.from_address;

					if (replyAllId) {
						const allTo = original.to_addresses
							.map((a) => a.email)
							.filter((e) => e !== currentEmail);
						const allCc = original.cc_addresses
							.map((a) => a.email)
							.filter((e) => e !== currentEmail && e !== original.from_address);
						if (allTo.length > 0) {
							to = [original.from_address, ...allTo].join(', ');
						}
						if (allCc.length > 0) {
							cc = allCc.join(', ');
						}
					}

					subject = original.subject.replace(/^(Re:\s*)+/i, '');
					subject = `Re: ${subject}`;

					inReplyTo = original.message_id;
					const existingRefs = original.references ? original.references.split(/\s+/).filter(Boolean) : [];
					references = [...existingRefs, original.message_id];

					const originalBody = original.body_html || `<p>${original.body_text}</p>`;
					const quotedBlock = `<blockquote style="border-left: 2px solid #ccc; padding-left: 12px; margin-left: 0; color: #666;"><p>On ${senderDate}, ${senderName} wrote:</p>${originalBody}</blockquote>`;
					initialContent = `<p><br></p>${signatureHtml}<p><br></p>${quotedBlock}`;
				}

				if (forwardId) {
					subject = original.subject.replace(/^(Fwd:\s*)+/i, '');
					subject = `Fwd: ${subject}`;

					const toList = original.to_addresses.map((a) => a.email).join(', ');
					const originalBody = original.body_html || `<p>${original.body_text}</p>`;
					const forwardBlock = `<p>---------- Forwarded message ----------</p><p>From: ${senderName} &lt;${original.from_address}&gt;</p><p>Date: ${senderDate}</p><p>Subject: ${original.subject}</p><p>To: ${toList}</p><br>${originalBody}`;
					initialContent = `<p><br></p>${signatureHtml}<p><br></p>${forwardBlock}`;
				}
			} catch {
				toast.error('Failed to load original email');
			}
		}

		if (!initialContent) {
			initialContent = `<p><br></p>${signatureHtml}`;
		}

		ready = true;
	});

	onDestroy(() => {
		clearTimeout(saveTimer);
	});

	async function send() {
		if (!to.trim() || !app.defaultAccountId) return;
		sending = true;
		try {
			const toAddrs = to.split(',').map((s) => s.trim()).filter(Boolean);
			const ccAddrs = cc ? cc.split(',').map((s) => s.trim()).filter(Boolean) : [];

			if (attachedFiles.length > 0) {
				const formData = new FormData();
				formData.set('to', toAddrs.join(','));
				if (ccAddrs.length > 0) formData.set('cc', ccAddrs.join(','));
				formData.set('subject', subject);
				if (plainTextMode) {
					formData.set('body', bodyPlainText);
				} else {
					formData.set('body', bodyHtml.replace(/<[^>]*>/g, ''));
					formData.set('body_html', bodyHtml);
				}
				if (inReplyTo) formData.set('in_reply_to', inReplyTo);
				if (references.length > 0) formData.set('references', references.join(','));
				for (const file of attachedFiles) {
					formData.append('attachments', file, file.name);
				}
				await backend.sendEmailWithAttachments(app.token, app.defaultAccountId, formData);
			} else {
				const body = plainTextMode ? bodyPlainText : bodyHtml.replace(/<[^>]*>/g, '');
				await backend.sendEmail(app.token, app.defaultAccountId, {
					to: toAddrs,
					cc: ccAddrs,
					subject,
					body,
					body_html: plainTextMode ? undefined : bodyHtml,
					in_reply_to: inReplyTo || undefined,
					references: references.length > 0 ? references : undefined
				});
			}

			if (draftId && app.defaultAccountId) {
				try {
					await backend.deleteDraft(app.token, app.defaultAccountId, draftId);
				} catch {
					// best effort
				}
			}

			toast.success('Email sent');
			goto('/mail');
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to send');
		}
		sending = false;
	}
</script>

<svelte:head>
	<title>Compose — Courrier</title>
</svelte:head>

<input
	type="file"
	multiple
	class="hidden"
	bind:this={fileInput}
	onchange={handleFileSelect}
/>

<div class="flex h-full flex-col">
	<div class="flex items-center justify-between border-b px-6 py-3">
		<div class="flex items-center gap-3">
			<h2 class="text-lg font-semibold">New message</h2>
			{#if lastSaved}
				<span class="text-xs text-muted-foreground">Draft saved {lastSaved}</span>
			{/if}
		</div>
		<div class="flex items-center gap-1">
			<Button
				variant="ghost"
				size="icon"
				class="h-8 w-8"
				title={plainTextMode ? 'Switch to rich text' : 'Switch to plain text'}
				onclick={togglePlainText}
			>
				{#if plainTextMode}
					<Type class="h-4 w-4" />
				{:else}
					<FileText class="h-4 w-4" />
				{/if}
			</Button>
			<Button
				variant="ghost"
				size="icon"
				class="h-8 w-8"
				title="Attach files"
				onclick={() => fileInput.click()}
			>
				<Paperclip class="h-4 w-4" />
			</Button>
			<div class="mx-1 h-4 w-px bg-border"></div>
			<Button variant="ghost" size="sm" class="gap-1.5" onclick={() => goto('/mail')}>
				<X class="h-4 w-4" />
				Cancel
			</Button>
			<Button size="sm" class="gap-1.5" disabled={sending || !to.trim() || !app.defaultAccountId} onclick={send}>
				<SendHorizonal class="h-4 w-4" />
				{sending ? 'Sending...' : 'Send'}
			</Button>
		</div>
	</div>

	{#if app.accounts.length === 0}
		<div class="flex flex-1 items-center justify-center text-muted-foreground">
			<div class="text-center">
				<p class="text-sm">No mail accounts configured</p>
				<p class="text-xs mt-1">Add one in Settings first</p>
			</div>
		</div>
	{:else}
		<div class="flex flex-col gap-0 border-b">
			<div class="flex items-center border-b px-6">
				<AddressInput
					id="to"
					label="To"
					value={to}
					onchange={(v) => { to = v; scheduleDraftSave(); }}
					token={app.token}
					accountId={app.defaultAccountId ?? 0}
					placeholder="recipient@example.com"
				/>
			</div>
			<div class="flex items-center border-b px-6">
				<AddressInput
					id="cc"
					label="Cc"
					value={cc}
					onchange={(v) => { cc = v; scheduleDraftSave(); }}
					token={app.token}
					accountId={app.defaultAccountId ?? 0}
					placeholder="cc@example.com"
				/>
			</div>
			<div class="flex items-center px-6">
				<Label for="subject" class="w-16 shrink-0 text-sm text-muted-foreground">Subject</Label>
				<Input
					id="subject"
					bind:value={subject}
					oninput={scheduleDraftSave}
					placeholder="Subject"
					class="border-0 shadow-none focus-visible:ring-0 rounded-none"
				/>
			</div>
		</div>

		{#if attachedFiles.length > 0}
			<div class="flex flex-wrap gap-2 border-b px-6 py-2">
				{#each attachedFiles as file, i}
					<Badge variant="secondary" class="gap-1.5 pr-1">
						<span class="max-w-[200px] truncate">{file.name}</span>
						<span class="text-muted-foreground">({formatFileSize(file.size)})</span>
						<button
							type="button"
							class="ml-0.5 rounded-full p-0.5 hover:bg-muted"
							onclick={() => removeFile(i)}
						>
							<X class="h-3 w-3" />
						</button>
					</Badge>
				{/each}
			</div>
		{/if}

		{#if ready}
			{#if plainTextMode}
				<div class="flex-1 overflow-auto px-6 py-4">
					<textarea
						bind:value={bodyPlainText}
						oninput={scheduleDraftSave}
						placeholder="Write your message..."
						class="h-full w-full resize-none bg-transparent text-sm leading-relaxed outline-none"
					></textarea>
				</div>
			{:else}
				<TiptapEditor
					content={initialContent}
					onchange={(html) => { bodyHtml = html; scheduleDraftSave(); }}
				/>
			{/if}
		{/if}
	{/if}
</div>
