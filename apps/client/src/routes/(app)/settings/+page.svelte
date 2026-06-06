<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import { backend, type MailAccount } from '$lib/backend';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { toast } from 'svelte-sonner';
	import { Plus, Trash2 } from 'lucide-svelte';

	const app = getContext<{
		token: string;
		refreshAccounts: () => Promise<void>;
	}>('app');

	let accounts = $state<MailAccount[]>([]);
	let showForm = $state(false);
	let saving = $state(false);
	let testing = $state(false);
	let deleting = $state<number | null>(null);

	let name = $state('');
	let email = $state('');
	let imapHost = $state('');
	let imapPort = $state(993);
	let imapUser = $state('');
	let imapPassword = $state('');
	let smtpHost = $state('');
	let smtpPort = $state(587);
	let smtpUser = $state('');
	let smtpPassword = $state('');

	function resetForm() {
		name = '';
		email = '';
		imapHost = '';
		imapPort = 993;
		imapUser = '';
		imapPassword = '';
		smtpHost = '';
		smtpPort = 587;
		smtpUser = '';
		smtpPassword = '';
		showForm = false;
	}

	async function loadAccounts() {
		const result = await backend.listAccounts(app.token);
		accounts = result.accounts;
	}

	async function testConnection() {
		testing = true;
		try {
			await backend.testConnection({
				imap_host: imapHost,
				imap_port: imapPort,
				imap_user: imapUser,
				imap_password: imapPassword,
				smtp_host: smtpHost,
				smtp_port: smtpPort,
				smtp_user: smtpUser,
				smtp_password: smtpPassword
			});
			toast.success('Connection test passed');
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Connection test failed');
		}
		testing = false;
	}

	async function addAccount() {
		if (!name || !email || !imapHost || !smtpHost) return;
		saving = true;
		try {
			await backend.createAccount(app.token, {
				name,
				email,
				imap_host: imapHost,
				imap_port: imapPort,
				imap_user: imapUser,
				imap_password: imapPassword,
				smtp_host: smtpHost,
				smtp_port: smtpPort,
				smtp_user: smtpUser,
				smtp_password: smtpPassword,
				signature: '',
				is_default: accounts.length === 0
			});
			toast.success('Account added');
			resetForm();
			await loadAccounts();
			await app.refreshAccounts();
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to add account');
		}
		saving = false;
	}

	async function deleteAccount(id: number) {
		deleting = id;
		try {
			await backend.deleteAccount(app.token, id);
			toast.success('Account deleted');
			await loadAccounts();
			await app.refreshAccounts();
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to delete account');
		}
		deleting = null;
	}

	onMount(loadAccounts);
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
			{#if accounts.length === 0 && !showForm}
				<p class="text-sm text-muted-foreground mb-4">No mail accounts configured yet.</p>
			{:else}
				<div class="space-y-3 mb-4">
					{#each accounts as account}
						<div class="flex items-center justify-between rounded-lg border p-3">
							<div>
								<p class="font-medium text-sm">{account.name}</p>
								<p class="text-xs text-muted-foreground">{account.email}</p>
							</div>
							<div class="flex items-center gap-2">
								<p class="text-xs text-muted-foreground">{account.imap_host}</p>
								{#if account.is_default}
									<span class="text-xs bg-foreground text-background px-2 py-0.5 rounded-full">Default</span>
								{/if}
								<Button
									variant="ghost"
									size="icon"
									class="h-8 w-8 text-muted-foreground hover:text-destructive"
									onclick={() => deleteAccount(account.id)}
									disabled={deleting === account.id}
								>
									<Trash2 class="h-4 w-4" />
								</Button>
							</div>
						</div>
					{/each}
				</div>
			{/if}

			{#if showForm}
				<div class="space-y-4 border rounded-lg p-4">
					<div class="grid grid-cols-2 gap-4">
						<div class="space-y-2">
							<Label for="name">Display name</Label>
							<Input id="name" bind:value={name} placeholder="Work Email" />
						</div>
						<div class="space-y-2">
							<Label for="email">Email address</Label>
							<Input id="email" type="email" bind:value={email} placeholder="you@example.com" />
						</div>
					</div>

					<div class="pt-2">
						<p class="text-sm font-medium mb-3">IMAP (Incoming)</p>
						<div class="grid grid-cols-2 gap-4">
							<div class="space-y-2">
								<Label for="imap-host">Host</Label>
								<Input id="imap-host" bind:value={imapHost} placeholder="imap.example.com" />
							</div>
							<div class="space-y-2">
								<Label for="imap-port">Port</Label>
								<Input id="imap-port" type="number" bind:value={imapPort} />
							</div>
							<div class="space-y-2">
								<Label for="imap-user">Username</Label>
								<Input id="imap-user" bind:value={imapUser} placeholder="you@example.com" />
							</div>
							<div class="space-y-2">
								<Label for="imap-password">Password</Label>
								<Input id="imap-password" type="password" bind:value={imapPassword} />
							</div>
						</div>
					</div>

					<div class="pt-2">
						<p class="text-sm font-medium mb-3">SMTP (Outgoing)</p>
						<div class="grid grid-cols-2 gap-4">
							<div class="space-y-2">
								<Label for="smtp-host">Host</Label>
								<Input id="smtp-host" bind:value={smtpHost} placeholder="smtp.example.com" />
							</div>
							<div class="space-y-2">
								<Label for="smtp-port">Port</Label>
								<Input id="smtp-port" type="number" bind:value={smtpPort} />
							</div>
							<div class="space-y-2">
								<Label for="smtp-user">Username</Label>
								<Input id="smtp-user" bind:value={smtpUser} placeholder="you@example.com" />
							</div>
							<div class="space-y-2">
								<Label for="smtp-password">Password</Label>
								<Input id="smtp-password" type="password" bind:value={smtpPassword} />
							</div>
						</div>
					</div>

					<div class="flex items-center gap-2 pt-2">
						<Button variant="outline" size="sm" onclick={testConnection} disabled={testing || !imapHost}>
							{testing ? 'Testing...' : 'Test Connection'}
						</Button>
						<div class="flex-1"></div>
						<Button variant="ghost" size="sm" onclick={resetForm}>Cancel</Button>
						<Button size="sm" class="gap-1.5" onclick={addAccount} disabled={saving || !name || !email || !imapHost || !smtpHost}>
							<Plus class="h-4 w-4" />
							{saving ? 'Saving...' : 'Add Account'}
						</Button>
					</div>
				</div>
			{:else}
				<Button variant="outline" size="sm" class="gap-1.5" onclick={() => (showForm = true)}>
					<Plus class="h-4 w-4" />
					Add Account
				</Button>
			{/if}
		</CardContent>
	</Card>
</div>
