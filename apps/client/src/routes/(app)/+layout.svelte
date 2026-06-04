<script lang="ts">
	import { onMount, setContext } from 'svelte';
	import { goto } from '$app/navigation';
	import { backend, type UserProfile, type MailAccount, type Folder } from '$lib/backend';
	import Sidebar from '$lib/components/Sidebar.svelte';

	let { children } = $props();

	let token = $state('');
	let user = $state<UserProfile | null>(null);
	let loaded = $state(false);
	let accounts = $state<MailAccount[]>([]);
	let defaultAccountId = $state<number | null>(null);
	let folders = $state<Folder[]>([]);

	function setUser(nextUser: UserProfile) {
		user = nextUser;
	}

	setContext('app', {
		get token() { return token; },
		get user() { return user; },
		get accounts() { return accounts; },
		get defaultAccountId() { return defaultAccountId; },
		get folders() { return folders; },
		setUser,
		refreshAccounts
	});

	async function refreshAccounts() {
		if (!token) return;
		try {
			const result = await backend.listAccounts(token);
			accounts = result.accounts;
			const def = accounts.find((a) => a.is_default) ?? accounts[0] ?? null;
			defaultAccountId = def?.id ?? null;

			if (defaultAccountId) {
				try {
					const folderResult = await backend.getFolders(token, defaultAccountId);
					folders = folderResult.folders;
				} catch {
					folders = [];
				}
			}
		} catch {
			accounts = [];
		}
	}

	onMount(async () => {
		const stored = localStorage.getItem('courrier.token') ?? '';
		if (!stored) {
			goto('/login');
			return;
		}
		try {
			const result = await backend.me(stored);
			token = stored;
			user = result.user;
			loaded = true;
			backend.syncProfile(stored).then(async (r) => {
				if (r.synced) {
					const fresh = await backend.me(stored);
					user = fresh.user;
				}
			}).catch(() => {});
			await refreshAccounts();
		} catch {
			goto('/login');
		}
	});
</script>

{#if loaded}
	<div class="flex h-screen w-full overflow-hidden">
		<Sidebar {user} {folders} />
		<main class="flex-1 overflow-auto">
			{@render children()}
		</main>
	</div>
{/if}
