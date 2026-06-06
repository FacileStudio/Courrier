<script lang="ts">
	import { goto } from '$app/navigation';
	import { getContext, onMount } from 'svelte';
	import { backend, type MailAccount } from '$lib/backend';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { toast } from 'svelte-sonner';
	import { X, SendHorizonal } from 'lucide-svelte';

	const app = getContext<{
		token: string;
		defaultAccountId: number | null;
		accounts: MailAccount[];
	}>('app');

	let to = $state('');
	let cc = $state('');
	let subject = $state('');
	let body = $state('');
	let sending = $state(false);

	onMount(() => {
		const account = app.accounts.find((a) => a.id === app.defaultAccountId);
		if (account?.signature) {
			body = `\n\n--\n${account.signature}`;
		}
	});

	async function send() {
		if (!to.trim() || !app.defaultAccountId) return;
		sending = true;
		try {
			const toAddrs = to.split(',').map((s) => s.trim()).filter(Boolean);
			const ccAddrs = cc ? cc.split(',').map((s) => s.trim()).filter(Boolean) : [];
			await backend.sendEmail(app.token, app.defaultAccountId, {
				to: toAddrs,
				cc: ccAddrs,
				subject,
				body
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

		<div class="flex-1 overflow-auto p-6">
			<textarea
				bind:value={body}
				placeholder="Write your message..."
				class="h-full w-full resize-none bg-transparent text-sm leading-relaxed outline-none placeholder:text-muted-foreground"
			></textarea>
		</div>
	{/if}
</div>
