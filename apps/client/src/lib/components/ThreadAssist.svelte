<script lang="ts">
	import { onMount } from 'svelte';
	import { Blob } from '@facile/blob';
	import { Button, Spinner, toast } from '@facile/muse';
	import { backend, type EmailMessage } from '$lib/backend';
	import { assistStatus } from '$lib/assist/status.svelte';
	import { plainOf, reasonText } from '$lib/assist/text';

	let { accountId, messages }: { accountId: number; messages: EmailMessage[] } = $props();

	let busy = $state(false);
	let summary = $state('');

	const status = $derived(assistStatus.current);
	const available = $derived(status?.available ?? false);

	onMount(() => {
		void assistStatus.load();
	});

	async function withBody(msg: EmailMessage): Promise<EmailMessage> {
		if (msg.body_text || msg.body_html) return msg;
		try {
			return await backend.getEmail(accountId, msg.id);
		} catch {
			return msg;
		}
	}

	async function threadText(): Promise<string> {
		const full = await Promise.all(messages.map(withBody));
		return full
			.map((m) => `From: ${m.from_name || ''} <${m.from_address}>\nDate: ${m.date}\nSubject: ${m.subject}\n\n${plainOf(m.body_html || m.body_text || '', 6000)}`)
			.join('\n\n---\n\n')
			.slice(-38000);
	}

	async function summarize() {
		busy = true;
		try {
			summary = (await backend.assist({ action: 'summarize', context: await threadText() })).text;
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : 'Bulle could not answer');
		}
		busy = false;
	}
</script>

{#if status && (available || status.reason !== 'off')}
	<div class="flex flex-col gap-2">
		<div class="flex flex-wrap items-center gap-2">
			<Button size="sm" variant="outline" disabled={!available || busy || messages.length === 0} title={available ? '' : reasonText(status)} onclick={summarize}>
				{#if busy}<Spinner size="sm" />{:else}<Blob name="Bulle" size={16} class="size-4" />{/if}
				Summarize
			</Button>
			{#if status.bulle_url}
				<Button size="sm" variant="ghost" href={status.bulle_url} target="_blank" rel="noopener">Ask Bulle</Button>
			{/if}
		</div>
		{#if summary}
			<div class="flex flex-col gap-2 rounded-fc-lg border border-fc-border bg-fc-surface p-3" role="region" aria-label="Thread summary">
				<p class="whitespace-pre-wrap text-fc-sm text-fc-fg">{summary}</p>
				<Button size="sm" variant="ghost" class="self-start" onclick={() => (summary = '')}>Hide</Button>
			</div>
		{/if}
	</div>
{/if}
