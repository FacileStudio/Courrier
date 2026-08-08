<script lang="ts">
	import { page } from '$app/state';
	import { Divider, PageTransition, Tabs, icons } from '@facile/muse';

	let { children } = $props();

	/*
	 * The section is a route, not a `$state`: /settings/accounts survives a reload, can be
	 * linked to, and browser-back walks the sections instead of leaving the page.
	 */
	const sections = [
		{ id: 'profile', label: 'Profile', icon: icons.userCircle },
		{ id: 'appearance', label: 'Appearance', icon: icons.palette },
		{ id: 'accounts', label: 'Accounts', icon: icons.mail },
		{ id: 'templates', label: 'Templates', icon: icons.edit },
		{ id: 'api', label: 'API', icon: icons.key },
		{ id: 'advanced', label: 'Advanced', icon: icons.shield }
	];

	const items = sections.map((section) => ({ ...section, href: `/settings/${section.id}` }));
	const active = $derived(
		sections.find((section) => page.url.pathname.startsWith(`/settings/${section.id}`))?.id ??
			sections[0].id
	);
</script>

<svelte:head>
	<title>Settings — Courrier</title>
</svelte:head>

<div class="mx-auto flex w-full max-w-4xl flex-col gap-8 p-6">
	<div class="flex flex-col gap-2">
		<h1 class="text-fc-2xl font-semibold text-fc-fg">Settings</h1>
		<p class="text-fc-sm text-fc-fg-muted">
			You, this browser, and the mail servers Courrier talks to on your behalf.
		</p>
	</div>

	<!-- gap-4, not tighter: the rule separates a page header from its body, so it needs air. -->
	<div class="flex flex-col gap-4">
		<Tabs {items} value={active} label="Settings sections" />
		<Divider class="my-0" />
	</div>

	<PageTransition key={active} distance={8} duration={0.25}>
		{@render children()}
	</PageTransition>
</div>
