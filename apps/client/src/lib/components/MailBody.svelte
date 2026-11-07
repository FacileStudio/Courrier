<script lang="ts">
	import { t } from '$lib/i18n/index.svelte';

	let { html }: { html: string } = $props();

	let frame = $state<HTMLIFrameElement>();
	let height = $state(0);

	const reset =
		'html{color-scheme:light}body{margin:0;display:flow-root;background:#fff;color:#1a1a1a;font:14px/1.5 system-ui,sans-serif;overflow-wrap:anywhere}img{max-width:100%;height:auto}';

	const srcdoc = $derived(
		`<!doctype html><html><head><meta charset="utf-8"><base target="_blank"><style>${reset}</style></head><body>${html}</body></html>`
	);

	$effect(() => {
		srcdoc;
		const el = frame;
		if (!el) return;
		let observer: ResizeObserver | undefined;

		const measure = () => {
			const doc = el.contentDocument;
			if (doc?.body) height = Math.ceil(doc.body.getBoundingClientRect().height);
		};

		const attach = () => {
			observer?.disconnect();
			const body = el.contentDocument?.body;
			if (!body) return;
			measure();
			observer = new ResizeObserver(measure);
			observer.observe(body);
		};

		el.addEventListener('load', attach);
		if (el.contentDocument?.readyState === 'complete') attach();

		return () => {
			el.removeEventListener('load', attach);
			observer?.disconnect();
		};
	});
</script>

<div class="overflow-hidden rounded-fc-md">
	<iframe
		bind:this={frame}
		title={t('mail.messageBody')}
		{srcdoc}
		sandbox="allow-same-origin allow-popups allow-popups-to-escape-sandbox"
		scrolling="no"
		class="block w-full border-0 bg-white"
		style:height={height ? `${height}px` : undefined}
	></iframe>
</div>
