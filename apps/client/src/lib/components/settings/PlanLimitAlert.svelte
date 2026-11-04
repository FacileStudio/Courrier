<script lang="ts">
	import { Alert, Button, icons } from '@facile/muse';
	import type { PlanLimit } from '$lib/plan-limit';

	let { limit }: { limit: PlanLimit } = $props();
</script>

<Alert tone="warning" title="Plan limit reached">
	<div class="flex flex-col gap-3">
		<p>{limit.message}</p>
		{#if limit.upgradeUrl}
			<div>
				<Button size="sm" icon={icons.card} href={limit.upgradeUrl} target="_blank" rel="noopener noreferrer">
					Upgrade the plan
				</Button>
			</div>
		{:else}
			<p class="text-fc-xs">Ask a space owner to upgrade its plan in Registre.</p>
		{/if}
	</div>
</Alert>
