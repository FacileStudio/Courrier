<script lang="ts">
	import { Button, icons } from '@facile/muse';
	import AccountSwitcher from '$lib/components/AccountSwitcher.svelte';
</script>

<div class="fixed top-3 left-1/2 z-50 hidden -translate-x-1/2 md:flex">
	<div
		class="flex items-center gap-1.5 rounded-fc-pill bg-fc-bg/70 p-1.5 shadow-lg backdrop-blur-2xl backdrop-saturate-150"
	>
		<AccountSwitcher />
		<Button href="/mail/compose" icon={icons.edit} class="rounded-fc-pill">Compose</Button>
	</div>
</div>
