<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import MailView from '$lib/components/MailView.svelte';

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

	$effect(() => {
		if (!validFolder) {
			goto('/mail');
		}
	});
</script>

<svelte:head>
	<title>{folderLabel} — Courrier</title>
</svelte:head>

{#if validFolder}
	{#key folderSlug}
		<MailView folder={folderSlug} />
	{/key}
{/if}
