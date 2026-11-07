<script lang="ts">
	import { getCurrentWindow } from '@tauri-apps/api/window';
	import type { AppKey } from '$lib/i18n/messages';
	import { hasNativeTrafficLights } from '$lib/platform';
	import { t } from '$lib/i18n/index.svelte';

	type Control = { labelKey: AppKey; color: string; glyph: string; run: () => Promise<void> };

	const drawn = !hasNativeTrafficLights();
	const win = getCurrentWindow();

	const controls: Control[] = [
		{ labelKey: 'window.close', color: 'bg-[#ff5f57]', glyph: 'M3.5 3.5l5 5m0-5l-5 5', run: () => win.close() },
		{ labelKey: 'window.minimize', color: 'bg-[#febc2e]', glyph: 'M3 6h6', run: () => win.minimize() },
		{ labelKey: 'window.maximize', color: 'bg-[#28c840]', glyph: 'M3.5 6h5M6 3.5v5', run: () => win.toggleMaximize() }
	];
</script>

<div
	data-tauri-drag-region
	class="fixed inset-x-0 top-0 z-[60] flex h-[var(--titlebar-h)] items-center bg-fc-page px-3.5 select-none"
>
	{#if drawn}
		<div class="group flex items-center gap-2" role="group" aria-label={t('window.label')}>
			{#each controls as control (control.labelKey)}
				<button
					type="button"
					aria-label={t(control.labelKey)}
					title={t(control.labelKey)}
					onclick={() => void control.run()}
					class="flex size-3 items-center justify-center rounded-full {control.color} focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-fc-ring"
				>
					<svg viewBox="0 0 12 12" class="size-2 opacity-0 transition-opacity group-hover:opacity-100" aria-hidden="true">
						<path d={control.glyph} stroke="rgb(0 0 0 / 0.6)" stroke-width="1.3" stroke-linecap="round" fill="none" />
					</svg>
				</button>
			{/each}
		</div>
	{/if}
</div>
