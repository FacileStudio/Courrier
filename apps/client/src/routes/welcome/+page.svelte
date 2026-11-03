<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Spinner } from '@facile/muse';
	import { ApiError, backend } from '$lib/backend';
	import Onboarding from '$lib/components/onboarding/Onboarding.svelte';
	import { needsOnboarding, onboardingSeen } from '$lib/onboarding/steps';
	import { isNative } from '$lib/platform';

	let name = $state<string | null>(null);
	let failed = $state(false);

	onMount(async () => {
		try {
			const { user } = await backend.me();
			const { accounts } = await backend.listAccounts();
			if (!needsOnboarding({ native: isNative(), seen: onboardingSeen(), accounts: accounts.length })) {
				await goto('/mail', { replaceState: true });
				return;
			}
			name = user.name?.trim().split(/\s+/)[0] || user.email;
		} catch (err) {
			if (err instanceof ApiError && (err.status === 401 || err.status === 403)) {
				await goto('/login', { replaceState: true });
				return;
			}
			failed = true;
		}
	});
</script>

<svelte:head>
	<title>Welcome — Courrier</title>
</svelte:head>

{#if name}
	<Onboarding {name} />
{:else if failed}
	<div class="flex h-dvh items-center justify-center bg-fc-page px-4 text-center text-fc-sm text-fc-fg-muted">
		Courrier is not answering. Reload the page to try again.
	</div>
{:else}
	<div class="flex h-dvh items-center justify-center bg-fc-page">
		<Spinner />
	</div>
{/if}
