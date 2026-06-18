<script lang="ts">
	import { onMount } from 'svelte';
	import { backend, type EmailTemplate } from '$lib/backend';

	let templates = $state<EmailTemplate[]>([]);
	let loading = $state(true);

	onMount(async () => {
		try {
			const res = await backend.listTemplates();
			templates = res.templates ?? [];
		} catch {
			templates = [];
		}
		loading = false;
	});
</script>

<svelte:head>
	<title>Templates — Courrier</title>
</svelte:head>

<div class="flex h-full flex-col">
	<div class="border-b border-border px-4 py-4 md:px-8 md:py-5">
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-lg font-semibold">Templates</h1>
				<p class="mt-0.5 text-sm text-muted-foreground">Reusable email templates</p>
			</div>
			<a
				href="/templates/new"
				class="inline-flex h-9 items-center gap-2 rounded-md bg-primary px-4 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
			>
				<iconify-icon icon="solar:add-circle-linear" width="16"></iconify-icon>
				New template
			</a>
		</div>
	</div>

	<div class="flex-1 overflow-auto p-4 md:p-8">
		{#if loading}
			<div class="flex items-center justify-center py-12 text-muted-foreground">
				<iconify-icon icon="solar:refresh-linear" width="20" class="animate-spin"></iconify-icon>
			</div>
		{:else if templates.length === 0}
			<div class="flex flex-col items-center justify-center py-16 text-center">
				<iconify-icon icon="solar:document-bold-duotone" width="48" class="text-muted-foreground/50"></iconify-icon>
				<p class="mt-4 text-sm text-muted-foreground">No templates yet</p>
				<a
					href="/templates/new"
					class="mt-4 inline-flex h-9 items-center gap-2 rounded-md border border-border px-4 text-sm font-medium transition-colors hover:bg-accent"
				>
					<iconify-icon icon="solar:add-circle-linear" width="16"></iconify-icon>
					Create your first template
				</a>
			</div>
		{:else}
			<div class="grid gap-3">
				{#each templates as template (template.id)}
					<div class="flex items-center gap-4 rounded-lg border border-border p-4 transition-colors hover:bg-muted/50">
						<div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-primary/10">
							<iconify-icon icon="solar:document-bold-duotone" width="20" class="text-primary"></iconify-icon>
						</div>
						<div class="min-w-0 flex-1">
							<p class="truncate text-sm font-medium">{template.name}</p>
							<p class="truncate text-xs text-muted-foreground">{template.subject}</p>
						</div>
						<iconify-icon icon="solar:alt-arrow-right-linear" width="16" class="shrink-0 text-muted-foreground"></iconify-icon>
					</div>
				{/each}
			</div>
		{/if}
	</div>
</div>
