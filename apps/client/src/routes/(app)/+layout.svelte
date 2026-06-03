<script lang="ts">
	import { onMount, setContext } from 'svelte';
	import { goto } from '$app/navigation';
	import { backend, type UserProfile } from '$lib/backend';
	import Sidebar from '$lib/components/Sidebar.svelte';

	let { children } = $props();

	let token = $state('');
	let user = $state<UserProfile | null>(null);
	let loaded = $state(false);

	function setUser(nextUser: UserProfile) {
		user = nextUser;
	}

	setContext('app', {
		get token() { return token; },
		get user() { return user; },
		setUser
	});

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
		} catch {
			goto('/login');
		}
	});
</script>

{#if loaded}
	<div class="flex h-screen w-full overflow-hidden">
		<Sidebar {user} />
		<main class="flex-1 overflow-auto">
			{@render children()}
		</main>
	</div>
{/if}
