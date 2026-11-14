<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { getContext, onMount, onDestroy } from 'svelte';
	import { backend, type MailAccount } from '$lib/backend';
	import { spaceStore } from '$lib/stores/space.svelte';
	import {
		Badge,
		Button,
		EmptyState,
		Field,
		IconButton,
		Input,
		Modal,
		Select,
		Textarea,
		Tooltip,
		icons,
		toast
	} from '@facile/muse';
	import type { EmailTemplate } from '$lib/backend';
	import { plainOf } from '$lib/assist/text';
	import TiptapEditor from '$lib/components/TiptapEditor.svelte';
	import AddressInput from '$lib/components/AddressInput.svelte';
	import { formatBytes, formatDate, formatTime } from '$lib/i18n/format';
	import { t } from '$lib/i18n/index.svelte';

	/* Glyphs an email client needs that muse's `icons` map has no key for yet. */
	const mailIcons = {
		paperclip: 'solar:paperclip-linear',
		send: 'solar:plain-linear',
		template: 'solar:document-text-linear',
		richText: 'solar:text-square-linear',
		plainText: 'solar:text-linear'
	};

	const app = getContext<{
		activeAccountId: number | null;
		accounts: MailAccount[];
	}>('app');

	const urlAccountId = Number(page.url.searchParams.get('accountId'));
	let pickedFromId = $state<number | null>(null);
	const fromId = $derived(
		pickedFromId ??
			(app.accounts.some((a) => a.id === urlAccountId) ? urlAccountId : app.activeAccountId)
	);

	let to = $state('');
	let cc = $state('');
	let subject = $state('');
	let bodyHtml = $state('');
	let bodyPlainText = $state('');
	let assistContext = $state('');
	let initialContent = $state('');
	let sending = $state(false);
	let inReplyTo = $state('');
	let references = $state<string[]>([]);
	let ready = $state(false);

	let attachedFiles = $state<File[]>([]);
	let fileInput: HTMLInputElement;

	let plainTextMode = $state(false);
	let templates = $state<EmailTemplate[]>([]);
	let templatesOpen = $state(false);

	let draftId = $state<number | null>(null);
	let draftAccountId: number | null = null;
	let lastSaved = $state<string | null>(null);
	let savingDraft = $state(false);
	let saveTimer: ReturnType<typeof setTimeout>;

	function escapeHtml(s: string): string {
		return s
			.replace(/&/g, '&amp;')
			.replace(/</g, '&lt;')
			.replace(/>/g, '&gt;')
			.replace(/"/g, '&quot;');
	}

	function buildSignatureHtml(signature: string): string {
		const sigLines = signature
			.split('\n')
			.map((l) => `<p>${escapeHtml(l) || '<br>'}</p>`)
			.join('');
		return `<p><br></p><p>--</p>${sigLines}`;
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
		if (!fromId || savingDraft) return;
		if (!to.trim() && !subject.trim() && !bodyHtml && !bodyPlainText) return;

		savingDraft = true;
		try {
			const toAddrs = to.split(',').map((s) => s.trim()).filter(Boolean);
			const ccAddrs = cc ? cc.split(',').map((s) => s.trim()).filter(Boolean) : [];
			const body = plainTextMode ? bodyPlainText : bodyHtml.replace(/<[^>]*>/g, '');
			const bodyHtmlValue = plainTextMode ? undefined : bodyHtml;

			if (draftId && draftAccountId) {
				await backend.deleteDraft(draftAccountId, draftId);
			}

			const result = await backend.saveDraft(fromId, {
				to: toAddrs,
				cc: ccAddrs,
				subject,
				body,
				body_html: bodyHtmlValue,
				in_reply_to: inReplyTo || undefined,
				references: references.length > 0 ? references : undefined
			});
			draftId = result.id;
			draftAccountId = fromId;
			lastSaved = formatTime(new Date());
		} catch {
			// silent fail
		}
		savingDraft = false;
	}

	function applyTemplate(tmpl: EmailTemplate) {
		subject = tmpl.subject;
		if (tmpl.body_html) {
			bodyHtml = tmpl.body_html;
			initialContent = tmpl.body_html;
		} else if (tmpl.body_text) {
			const paragraphs = tmpl.body_text.split('\n').map((line: string) => `<p>${line || '<br>'}</p>`).join('');
			bodyHtml = paragraphs;
			initialContent = paragraphs;
		}
		bodyPlainText = tmpl.body_text || tmpl.body_html?.replace(/<[^>]*>/g, '') || '';
		templatesOpen = false;
	}

	$effect(() => {
		// Reload templates when the active space changes.
		spaceStore.active;
		backend.listTemplates(spaceStore.active?.id).then((r) => { templates = r.templates; }).catch(() => {});
	});

	onMount(async () => {
		const params = page.url.searchParams;
		const replyId = params.get('reply');
		const replyAllId = params.get('replyall');
		const forwardId = params.get('forward');
		const accountId = params.get('accountId');
		const emailId = replyId || replyAllId || forwardId;
		const loadDraftId = params.get('draft');

		const account = app.accounts.find((a) => a.id === fromId);
		let signatureHtml = '';
		if (account?.signature) {
			signatureHtml = buildSignatureHtml(account.signature);
		}

		if (loadDraftId && fromId) {
			try {
				const draft = await backend.getEmail(fromId, parseInt(loadDraftId));
				to = draft.to_addresses.map((a) => a.email).join(', ');
				cc = draft.cc_addresses.map((a) => a.email).join(', ');
				subject = draft.subject;
				bodyHtml = draft.body_html || '';
				bodyPlainText = draft.body_text || '';
				initialContent = draft.body_html || `<p>${escapeHtml(draft.body_text || '')}</p>`;
				inReplyTo = draft.in_reply_to || '';
				references = draft.references ? draft.references.split(/\s+/).filter(Boolean) : [];
				draftId = parseInt(loadDraftId);
				draftAccountId = fromId;
			} catch {
				toast.danger(t('compose.draftFailed'));
			}
		} else if (emailId && accountId) {
			try {
				const original = await backend.getEmail(parseInt(accountId), parseInt(emailId));
				const acct = app.accounts.find((a) => a.id === parseInt(accountId));
				const currentEmail = acct?.email ?? '';
				const senderDate = escapeHtml(formatDate(original.date, { dateStyle: 'medium', timeStyle: 'short' }));
				const senderName = escapeHtml(original.from_name || original.from_address);

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

					assistContext = `From: ${original.from_name} <${original.from_address}>\nSubject: ${original.subject}\n\n${plainOf(original.body_html || original.body_text)}`;
					inReplyTo = original.message_id;
					const existingRefs = original.references ? original.references.split(/\s+/).filter(Boolean) : [];
					references = [...existingRefs, original.message_id];

					const originalBody = original.body_html || `<p>${escapeHtml(original.body_text)}</p>`;
					const quotedBlock = `<blockquote style="border-left: 2px solid #ccc; padding-left: 12px; margin-left: 0; color: #666;"><p>${t('compose.wrote', { date: senderDate, name: senderName })}</p>${originalBody}</blockquote>`;
					initialContent = `<p><br></p>${signatureHtml}<p><br></p>${quotedBlock}`;
				}

				if (forwardId) {
					subject = original.subject.replace(/^(Fwd:\s*)+/i, '');
					subject = `Fwd: ${subject}`;

					const toList = escapeHtml(original.to_addresses.map((a) => a.email).join(', '));
					const originalBody = original.body_html || `<p>${escapeHtml(original.body_text)}</p>`;
					const forwardBlock = `<p>${t('compose.forwardedHeader')}</p><p>${t('compose.fwdFrom')}: ${senderName} &lt;${escapeHtml(original.from_address)}&gt;</p><p>${t('compose.fwdDate')}: ${senderDate}</p><p>${t('compose.fwdSubject')}: ${escapeHtml(original.subject)}</p><p>${t('compose.fwdTo')}: ${toList}</p><br>${originalBody}`;
					initialContent = `<p><br></p>${signatureHtml}<p><br></p>${forwardBlock}`;
				}
			} catch {
				toast.danger(t('compose.originalFailed'));
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
		if (!to.trim() || !fromId) return;
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
				await backend.sendEmailWithAttachments(fromId, formData);
			} else {
				const body = plainTextMode ? bodyPlainText : bodyHtml.replace(/<[^>]*>/g, '');
				await backend.sendEmail(fromId, {
					to: toAddrs,
					cc: ccAddrs,
					subject,
					body,
					body_html: plainTextMode ? undefined : bodyHtml,
					in_reply_to: inReplyTo || undefined,
					references: references.length > 0 ? references : undefined
				});
			}

			if (draftId && draftAccountId) {
				try {
					await backend.deleteDraft(draftAccountId, draftId);
				} catch {
					// best effort
				}
			}

			toast.success(t('compose.sent'));
			goto('/mail');
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : t('compose.sendFailed'));
		}
		sending = false;
	}
</script>

<svelte:head>
	<title>{t('compose.title')}</title>
</svelte:head>

<input type="file" multiple class="hidden" bind:this={fileInput} onchange={handleFileSelect} />

<div class="flex h-full flex-col">
	<div class="flex items-center justify-between gap-3 px-4 py-3 sm:px-6">
		<div class="flex min-w-0 items-center gap-3">
			<h2 class="truncate text-fc-lg font-semibold text-fc-fg">{t('compose.heading')}</h2>
			{#if lastSaved}
				<span class="shrink-0 text-fc-xs text-fc-fg-muted">{t('compose.draftSaved', { time: lastSaved })}</span>
			{/if}
		</div>
		<div class="flex shrink-0 items-center gap-1">
			{#if templates.length > 0}
				<IconButton
					variant="ghost"
					aria-label={t('compose.useTemplate')}
					onclick={() => (templatesOpen = true)}
				>
					<iconify-icon icon={mailIcons.template} width="18" height="18" class="block size-4.5"
					></iconify-icon>
				</IconButton>
			{/if}
			<IconButton
				variant="ghost"
				aria-label={plainTextMode ? t('compose.toRich') : t('compose.toPlain')}
				onclick={togglePlainText}
			>
				<iconify-icon
					icon={plainTextMode ? mailIcons.richText : mailIcons.plainText}
					width="18"
					height="18"
					class="block size-4.5"
				></iconify-icon>
			</IconButton>
			<IconButton
				variant="ghost"
				aria-label={t('compose.attach')}
				onclick={() => fileInput.click()}
			>
				<iconify-icon icon={mailIcons.paperclip} width="18" height="18" class="block size-4.5"
				></iconify-icon>
			</IconButton>
			<span class="mx-1 h-5 w-px shrink-0 bg-fc-surface"></span>
			<Button variant="ghost" size="sm" icon={icons.close} onclick={() => goto('/mail')}>{t('common.cancel')}</Button>
			<Button
				size="sm"
				icon={mailIcons.send}
				disabled={sending || !to.trim() || !fromId}
				onclick={send}
			>
				{sending ? t('compose.sending') : t('compose.send')}
			</Button>
		</div>
	</div>

	{#if app.accounts.length === 0}
		<div class="flex flex-1 items-center justify-center px-4">
			<EmptyState
				bare
				icon={icons.mail}
				title={t('compose.noAccount')}
				description={t('compose.noAccountBody')}
			>
				<Button href="/settings" iconRight={icons.arrow}>{t('compose.openSettings')}</Button>
			</EmptyState>
		</div>
	{:else}
		<div class="flex flex-col gap-4 px-4 py-4 sm:px-6">
			{#if app.accounts.length > 1}
				<Field label={t('compose.from')}>
					<Select
						value={String(fromId ?? '')}
						options={app.accounts.map((account) => ({
							value: String(account.id),
							label: `${account.name} <${account.email}>`
						}))}
						onchange={(e) => {
							pickedFromId = Number(e.currentTarget.value);
							scheduleDraftSave();
						}}
					/>
				</Field>
			{/if}
			<Field label={t('compose.to')}>
				<AddressInput
					value={to}
					onchange={(v) => {
						to = v;
						scheduleDraftSave();
					}}
					accountId={fromId ?? 0}
					placeholder={t('compose.toPlaceholder')}
				/>
			</Field>
			<Field label={t('compose.cc')}>
				<AddressInput
					value={cc}
					onchange={(v) => {
						cc = v;
						scheduleDraftSave();
					}}
					accountId={fromId ?? 0}
					placeholder={t('compose.ccPlaceholder')}
				/>
			</Field>
			<Field label={t('compose.subject')}>
				<Input bind:value={subject} oninput={scheduleDraftSave} placeholder={t('compose.subject')} />
			</Field>
		</div>

		{#if attachedFiles.length > 0}
			<div class="flex flex-wrap gap-2 px-4 py-3 sm:px-6">
				{#each attachedFiles as file, i (`${file.name}-${i}`)}
					<Badge tone="neutral" class="max-w-full pr-1">
						<span class="max-w-48 truncate">{file.name}</span>
						<span>({formatBytes(file.size)})</span>
						<Tooltip text={t('compose.removeFile', { name: file.name })} describe={false} class="contents">
							<button
								type="button"
								aria-label={t('compose.removeFile', { name: file.name })}
								class="flex size-4 shrink-0 items-center justify-center rounded-fc-pill text-fc-fg-muted transition-colors hover:bg-fc-component hover:text-fc-fg focus-visible:outline-2 focus-visible:outline-offset-1 focus-visible:outline-fc-ring"
								onclick={() => removeFile(i)}
							>
								<iconify-icon icon={icons.close} width="12" height="12" class="block size-3"
								></iconify-icon>
							</button>
						</Tooltip>
					</Badge>
				{/each}
			</div>
		{/if}

		{#if ready}
			{#if plainTextMode}
				<div class="flex-1 overflow-auto px-4 py-4 pb-24 sm:px-6 md:pb-4">
					<Textarea
						bind:value={bodyPlainText}
						oninput={scheduleDraftSave}
						placeholder={t('editor.placeholder')}
						aria-label={t('compose.bodyLabel')}
						class="h-full resize-none text-fc-sm"
					/>
				</div>
			{:else}
				<TiptapEditor
					assist
					{assistContext}
					accountId={fromId ?? undefined}
					content={initialContent}
					onchange={(html) => {
						bodyHtml = html;
						scheduleDraftSave();
					}}
				/>
			{/if}
		{/if}
	{/if}
</div>

<Modal bind:open={templatesOpen} title={t('compose.templates')} showClose>
	<div class="flex flex-col gap-1">
		{#each templates as tmpl (tmpl.id)}
			<button
				type="button"
				class="flex w-full flex-col items-start gap-0.5 rounded-fc-sm px-3 py-2.5 text-left transition-colors hover:bg-fc-surface focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-fc-ring"
				onclick={() => applyTemplate(tmpl)}
			>
				<span class="text-fc-sm font-medium text-fc-fg">{tmpl.name}</span>
				{#if tmpl.subject}
					<span class="truncate text-fc-xs text-fc-fg-muted">{tmpl.subject}</span>
				{/if}
			</button>
		{/each}
	</div>
</Modal>
