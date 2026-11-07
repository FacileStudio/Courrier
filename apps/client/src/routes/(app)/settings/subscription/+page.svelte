<script lang="ts">
	import { Alert, SubscriptionSection } from '@facile/muse';
	import { backend } from '$lib/backend';
	import { t } from '$lib/i18n/index.svelte';
	import { SPACES_URL } from '$lib/plan-limit';
	import { spaceStore } from '$lib/stores/space.svelte';

	let usage = $state<{ used: number; limit: number } | null>(null);
	let loadError = $state('');

	$effect(() => {
		const spaceId = spaceStore.active?.id;
		usage = null;
		backend
			.accountUsage(spaceId)
			.then((next) => {
				usage = next;
				loadError = '';
			})
			.catch((err: unknown) => {
				loadError = err instanceof Error ? err.message : t('api.readFailed');
			});
	});

	const unlimited = $derived(usage !== null && usage.limit < 0);
	const plan = $derived(
		usage === null ? '' : unlimited ? t('subscription.unlimited') : t('subscription.allowance', { n: usage.limit })
	);
	const rows = $derived(
		usage !== null && usage.limit > 0
			? [{ label: t('subscription.mailboxes'), used: usage.used, limit: usage.limit }]
			: []
	);
</script>

{#if loadError}
	<Alert tone="danger">{loadError}</Alert>
{:else if usage}
	<SubscriptionSection {plan} usage={rows} manageHref={SPACES_URL} upgradeHref={unlimited ? undefined : SPACES_URL} />
{/if}
