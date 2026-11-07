<script lang="ts">
	import { Alert, Button, icons } from '@facile/muse';
	import type { PlanLimit } from '$lib/plan-limit';
	import { t } from '$lib/i18n/index.svelte';

	let { limit }: { limit: PlanLimit } = $props();
</script>

<Alert tone="warning" title={t('plan.reached')}>
	<div class="flex flex-col gap-3">
		<p>{limit.message}</p>
		{#if limit.upgradeUrl}
			<div>
				<Button size="sm" icon={icons.card} href={limit.upgradeUrl} target="_blank" rel="noopener noreferrer">
					{t('plan.upgrade')}
				</Button>
			</div>
		{:else}
			<p class="text-fc-xs">{t('plan.askOwner')}</p>
		{/if}
	</div>
</Alert>
