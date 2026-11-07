<script lang="ts">
	import AppMark from './AppMark.svelte';
	import Ambient from './Ambient.svelte';
	import { t } from '$lib/i18n/index.svelte';
	import SignInButton from './SignInButton.svelte';

	let { error = $bindable('') }: { error?: string } = $props();
</script>

<div class="fixed inset-0 flex flex-col overflow-hidden bg-fc-page text-fc-fg">
	<Ambient />

	<header class="relative z-10 flex h-14 items-center gap-2 px-5">
		<AppMark size={22} />
		<span class="font-fc-title text-fc-md font-semibold tracking-tight">Courrier</span>
	</header>

	<main
		class="relative z-10 flex min-h-0 flex-1 flex-col items-center justify-center gap-8 overflow-y-auto px-8 text-center"
	>
		<div
			class="flex size-32 items-center justify-center rounded-[2rem] bg-fc-component shadow-sm"
		>
			<AppMark size={64} />
		</div>
		<div class="flex max-w-xs flex-col gap-3">
			<h1 class="font-fc-title text-fc-3xl leading-tight font-semibold tracking-tight">
				{t('login.nativeTitle')}
			</h1>
			<p class="text-fc-md leading-relaxed text-fc-fg-muted">
				{t('login.nativeBody')}
			</p>
		</div>
	</main>

	<footer class="relative z-10 px-6 pt-4 pb-6">
		<SignInButton bind:error />
	</footer>
</div>
