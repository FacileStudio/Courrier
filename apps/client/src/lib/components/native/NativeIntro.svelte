<script lang="ts">
	import { Button, icons } from '@facile/muse';
	import AppMark from './AppMark.svelte';
	import Ambient from './Ambient.svelte';

	type Slide = { title: string; body: string; icon?: string };

	const SLIDES: Slide[] = [
		{
			title: 'Your mail, in your pocket.',
			body: 'Every IMAP mailbox you connected to Courrier, threaded, searchable and ready to answer.'
		},
		{
			title: 'Only the mail that matters.',
			body: 'Newsletters and robots stay quiet. A notification means a person wrote to you.',
			icon: icons.notification
		},
		{
			title: 'Your server, your rules.',
			body: 'Courrier is self-hosted. Your messages stay on your server, never on ours.',
			icon: icons.shield
		}
	];

	let { onfinish }: { onfinish: () => void } = $props();

	let track = $state<HTMLElement | null>(null);
	let index = $state(0);

	const last = $derived(index === SLIDES.length - 1);

	function onscroll() {
		if (!track) return;
		index = Math.round(track.scrollLeft / track.clientWidth);
	}

	function show(target: number) {
		const smooth = !window.matchMedia('(prefers-reduced-motion: reduce)').matches;
		track?.scrollTo({ left: target * track.clientWidth, behavior: smooth ? 'smooth' : 'auto' });
	}

	function next() {
		if (last) onfinish();
		else show(index + 1);
	}
</script>

<div class="fixed inset-0 flex flex-col overflow-hidden bg-fc-page text-fc-fg">
	<Ambient />

	<header class="relative z-10 flex h-14 items-center justify-between px-5">
		<div class="flex items-center gap-2">
			<AppMark size={22} />
			<span class="font-fc-title text-fc-md font-semibold tracking-tight">Courrier</span>
		</div>
		{#if !last}
			<Button variant="ghost" size="sm" class="min-h-11 px-3" onclick={onfinish}>Skip</Button>
		{/if}
	</header>

	<div
		bind:this={track}
		{onscroll}
		class="relative z-10 flex min-h-0 flex-1 snap-x snap-mandatory overflow-x-auto overscroll-x-contain [scrollbar-width:none] [&::-webkit-scrollbar]:hidden"
		aria-roledescription="carousel"
	>
		{#each SLIDES as slide, i (slide.title)}
			<section
				class="flex w-full shrink-0 snap-center flex-col items-center justify-center gap-8 px-8 text-center"
				aria-roledescription="slide"
				aria-label="{i + 1} of {SLIDES.length}"
				aria-hidden={i !== index}
			>
				<div
					class="flex size-36 items-center justify-center rounded-[2rem] border border-fc-border bg-fc-component shadow-sm"
				>
					{#if slide.icon}
						<iconify-icon icon={slide.icon} width="56" height="56" aria-hidden="true"></iconify-icon>
					{:else}
						<AppMark size={72} />
					{/if}
				</div>
				<div class="flex max-w-xs flex-col gap-3">
					<h1 class="font-fc-title text-fc-3xl leading-tight font-semibold tracking-tight">
						{slide.title}
					</h1>
					<p class="text-fc-md leading-relaxed text-fc-fg-muted">{slide.body}</p>
				</div>
			</section>
		{/each}
	</div>

	<footer class="relative z-10 flex flex-col items-center gap-6 px-6 pt-4 pb-6">
		<div class="flex items-center gap-1" role="tablist" aria-label="Slides">
			{#each SLIDES as slide, i (slide.title)}
				<button
					type="button"
					role="tab"
					aria-selected={i === index}
					aria-label="Go to slide {i + 1}"
					class="flex size-11 items-center justify-center"
					onclick={() => show(i)}
				>
					<span
						class="h-1.5 rounded-fc-pill transition-all duration-300 motion-reduce:transition-none {i === index
							? 'w-6 bg-fc-fg'
							: 'w-1.5 bg-fc-fg/25'}"
					></span>
				</button>
			{/each}
		</div>
		<Button size="lg" class="h-12 w-full max-w-sm rounded-fc-md" iconRight={icons.arrow} onclick={next}>
			{last ? 'Get started' : 'Next'}
		</Button>
	</footer>
</div>
