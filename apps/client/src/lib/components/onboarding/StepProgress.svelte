<script lang="ts">
	import { STEPS } from '$lib/onboarding/steps';
	import { t } from '$lib/i18n/index.svelte';

	let { index }: { index: number } = $props();
</script>

<div class="flex w-full max-w-fc-sm flex-col gap-2">
	<div class="flex items-center justify-between text-fc-xs text-fc-fg-muted">
		<span>{t('onboarding.stepOf', { index: index + 1, total: STEPS.length })}</span>
		<span>{t(STEPS[index].labelKey)}</span>
	</div>
	<div
		role="progressbar"
		aria-label={t('onboarding.progress')}
		aria-valuemin={1}
		aria-valuemax={STEPS.length}
		aria-valuenow={index + 1}
		aria-valuetext={t('onboarding.stepText', { index: index + 1, total: STEPS.length, label: t(STEPS[index].labelKey) })}
		class="flex gap-1.5"
	>
		{#each STEPS as item, position (item.id)}
			<span
				class="h-1.5 flex-1 rounded-fc-pill transition-colors duration-500 motion-reduce:transition-none {position <=
				index
					? 'bg-fc-accent'
					: 'bg-fc-surface'}"
			></span>
		{/each}
	</div>
</div>
