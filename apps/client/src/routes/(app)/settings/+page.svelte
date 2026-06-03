<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import { backend, type MailAccount } from '$lib/backend';
	import { Button } from '$lib/components/ui/button';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';

	const app = getContext<{ token: string }>('app');

	let accounts = $state<MailAccount[]>([]);

	onMount(async () => {
		const result = await backend.listAccounts(app.token);
		accounts = result.accounts;
	});
</script>

<div class="mx-auto max-w-2xl p-8">
	<div class="flex items-center justify-between mb-6">
		<h1 class="text-2xl font-bold">Settings</h1>
	</div>

	<Card>
		<CardHeader>
			<CardTitle>Mail Accounts</CardTitle>
		</CardHeader>
		<CardContent>
			{#if accounts.length === 0}
				<p class="text-sm text-muted-foreground mb-4">No mail accounts configured yet.</p>
			{:else}
				<div class="space-y-3 mb-4">
					{#each accounts as account}
						<div class="flex items-center justify-between rounded-lg border p-3">
							<div>
								<p class="font-medium text-sm">{account.name}</p>
								<p class="text-xs text-muted-foreground">{account.email}</p>
							</div>
							<p class="text-xs text-muted-foreground">{account.imap_host}</p>
						</div>
					{/each}
				</div>
			{/if}
			<Button variant="outline" size="sm">Add Account</Button>
		</CardContent>
	</Card>
</div>
