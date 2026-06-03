<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { page } from '$app/stores';

	const TOKEN_KEY = 'courrier.token';

	onMount(() => {
		const token = $page.url.searchParams.get('token');
		if (token) {
			localStorage.setItem(TOKEN_KEY, token);
		}
		goto('/mail');
	});
</script>
