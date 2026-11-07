<script lang="ts">
	import MailView from '$lib/components/MailView.svelte';
	import { t } from '$lib/i18n/index.svelte';
</script>

<svelte:head>
	<title>{t('mail.inboxTitle')}</title>
</svelte:head>

<MailView folder="inbox" />
