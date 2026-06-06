<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { getContext, onMount } from 'svelte';
	import { backend, type MailAccount } from '$lib/backend';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { toast } from 'svelte-sonner';
	import { X, SendHorizonal } from 'lucide-svelte';
	import TiptapEditor from '$lib/components/TiptapEditor.svelte';

	const app = getContext<{
		token: string;
		defaultAccountId: number | null;
		accounts: MailAccount[];
	}>('app');

	let to = $state('');
	let cc = $state('');
	let subject = $state('');
	let bodyHtml = $state('');
	let initialContent = $state('');
	let sending = $state(false);
	let inReplyTo = $state('');
	let references = $state<string[]>([]);
	let ready = $state(false);

	onMount(async () => {
		const params = page.url.searchParams;
		const replyId = params.get('reply');
		const replyAllId = params.get('replyall');
		const forwardId = params.get('forward');
		const accountId = params.get('accountId');
		const emailId = replyId || replyAllId || forwardId;

		if (emailId && accountId) {
			try {
				const original = await backend.getEmail(app.token, parseInt(accountId), parseInt(emailId));
				const account = app.accounts.find((a) => a.id === parseInt(accountId));
				const currentEmail = account?.email ?? '';
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
					initialContent = `<p><br></p><p><br></p><blockquote style="border-left: 2px solid #ccc; padding-left: 12px; margin-left: 0; color: #666;"><p>On ${senderDate}, ${senderName} wrote:</p>${originalBody}</blockquote>`;
				}

				if (forwardId) {
					subject = original.subject.replace(/^(Fwd:\s*)+/i, '');
					subject = `Fwd: ${subject}`;

					const toList = original.to_addresses.map((a) => a.email).join(', ');
					const originalBody = original.body_html || `<p>${original.body_text}</p>`;
					initialContent = `<p><br></p><p><br></p><p>---------- Forwarded message ----------</p><p>From: ${senderName} &lt;${original.from_address}&gt;</p><p>Date: ${senderDate}</p><p>Subject: ${original.subject}</p><p>To: ${toList}</p><br>${originalBody}`;
				}
			} catch {
				toast.error('Failed to load original email');
			}
		}

		const account = app.accounts.find((a) => a.id === app.defaultAccountId);
		if (account?.signature && !initialContent) {
			initialContent = `<p><br></p><p><br></p><p>--</p><p>${account.signature}</p>`;
		}

		ready = true;
	});

	async function send() {
		if (!to.trim() || !app.defaultAccountId) return;
		sending = true;
		try {
			const toAddrs = to.split(',').map((s) => s.trim()).filter(Boolean);
			const ccAddrs = cc ? cc.split(',').map((s) => s.trim()).filter(Boolean) : [];

			const plainText = bodyHtml.replace(/<[^>]*>/g, '');

			await backend.sendEmail(app.token, app.defaultAccountId, {
				to: toAddrs,
				cc: ccAddrs,
				subject,
				body: plainText,
				body_html: bodyHtml,
				in_reply_to: inReplyTo || undefined,
				references: references.length > 0 ? references : undefined
			});
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

<div class="flex h-full flex-col">
	<div class="flex items-center justify-between border-b px-6 py-3">
		<h2 class="text-lg font-semibold">New message</h2>
		<div class="flex items-center gap-2">
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
				<Label for="to" class="w-16 shrink-0 text-sm text-muted-foreground">To</Label>
				<Input
					id="to"
					type="email"
					bind:value={to}
					placeholder="recipient@example.com"
					class="border-0 shadow-none focus-visible:ring-0 rounded-none"
				/>
			</div>
			<div class="flex items-center border-b px-6">
				<Label for="cc" class="w-16 shrink-0 text-sm text-muted-foreground">Cc</Label>
				<Input
					id="cc"
					bind:value={cc}
					placeholder="cc@example.com"
					class="border-0 shadow-none focus-visible:ring-0 rounded-none"
				/>
			</div>
			<div class="flex items-center px-6">
				<Label for="subject" class="w-16 shrink-0 text-sm text-muted-foreground">Subject</Label>
				<Input
					id="subject"
					bind:value={subject}
					placeholder="Subject"
					class="border-0 shadow-none focus-visible:ring-0 rounded-none"
				/>
			</div>
		</div>

		{#if ready}
			<TiptapEditor
				content={initialContent}
				onchange={(html) => { bodyHtml = html; }}
			/>
		{/if}
	{/if}
</div>
