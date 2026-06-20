<script lang="ts">
	import { goto } from '$app/navigation';
	import { backend, type MailAccount } from '$lib/backend';
	import { spaceStore } from '$lib/stores/space.svelte';

	let accounts = $state<MailAccount[]>([]);
	let loading = $state(true);

	$effect(() => {
		const spaceId = spaceStore.active?.id;
		loading = true;
		backend
			.listAccounts(spaceId)
			.then((res) => {
				accounts = res.accounts ?? [];
			})
			.catch(() => {
				accounts = [];
			})
			.finally(() => {
				loading = false;
			});
	});
</script>

<svelte:head>
	<title>Accounts — Courrier</title>
</svelte:head>

<div class="flex h-full flex-col">
	<div class="border-b border-border px-4 py-4 md:px-8 md:py-5">
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-lg font-semibold">Accounts</h1>
				<p class="mt-0.5 text-sm text-muted-foreground">Manage your email accounts</p>
			</div>
			<a
				href="/accounts/new"
				class="inline-flex h-9 items-center gap-2 rounded-md bg-primary px-4 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
			>
				<iconify-icon icon="solar:add-circle-linear" width="16"></iconify-icon>
				Add account
			</a>
		</div>
	</div>

	<div class="flex-1 overflow-auto p-4 md:p-8">
		{#if loading}
			<div class="flex items-center justify-center py-12 text-muted-foreground">
				<iconify-icon icon="solar:refresh-linear" width="20" class="animate-spin"></iconify-icon>
			</div>
		{:else if accounts.length === 0}
			<div class="flex flex-col items-center justify-center py-16 text-center">
				<iconify-icon icon="solar:mailbox-bold-duotone" width="48" class="text-muted-foreground/50"></iconify-icon>
				<p class="mt-4 text-sm text-muted-foreground">No accounts yet</p>
				<a
					href="/accounts/new"
					class="mt-4 inline-flex h-9 items-center gap-2 rounded-md border border-border px-4 text-sm font-medium transition-colors hover:bg-accent"
				>
					<iconify-icon icon="solar:add-circle-linear" width="16"></iconify-icon>
					Add your first account
				</a>
			</div>
		{:else}
			<div class="grid gap-3">
				{#each accounts as account (account.id)}
					<div class="flex items-center gap-4 rounded-lg border border-border p-4 transition-colors hover:bg-muted/50">
						<div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-primary/10">
							<iconify-icon icon="solar:letter-bold-duotone" width="20" class="text-primary"></iconify-icon>
						</div>
						<div class="min-w-0 flex-1">
							<p class="truncate text-sm font-medium">{account.name}</p>
							<p class="truncate text-xs text-muted-foreground">{account.email}</p>
						</div>
						{#if account.is_default}
							<span class="shrink-0 inline-flex rounded-full bg-amber-500/10 px-2.5 py-0.5 text-xs font-medium text-amber-600">
								Default
							</span>
						{/if}
						<iconify-icon icon="solar:alt-arrow-right-linear" width="16" class="shrink-0 text-muted-foreground"></iconify-icon>
					</div>
				{/each}
			</div>
		{/if}
	</div>
</div>
