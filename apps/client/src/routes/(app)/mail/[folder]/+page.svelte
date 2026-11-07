<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import MailView from '$lib/components/MailView.svelte';
	import { MAIL_FOLDERS } from '$lib/mail-folders';
	import { t } from '$lib/i18n/index.svelte';

	const folderSlug = $derived(page.params.folder ?? '');
	const entry = $derived(MAIL_FOLDERS.find((folder) => folder.type === folderSlug && folder.type !== 'inbox'));
	const folderLabel = $derived(entry ? t(entry.labelKey) : folderSlug);
	const validFolder = $derived(entry !== undefined);

	$effect(() => {
		if (!validFolder) {
			goto('/mail');
		}
	});
</script>

<svelte:head>
	<title>{t('mail.folderTitle', { folder: folderLabel })}</title>
</svelte:head>

{#if validFolder}
	{#key folderSlug}
		<MailView folder={folderSlug} />
	{/key}
{/if}
