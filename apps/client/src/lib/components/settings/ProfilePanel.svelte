<script lang="ts">
	import { getContext } from 'svelte';
	import { ProfileSection } from '@facile/muse';
	import type { UserProfile } from '$lib/backend';
	import { t } from '$lib/i18n/index.svelte';
	import { spaceStore } from '$lib/stores/space.svelte';

	const app = getContext<{ user: UserProfile | null }>('app');
</script>

{#if app.user}
	<ProfileSection
		name={app.user.name?.trim() || app.user.email || t('common.you')}
		email={app.user.email}
		avatar={app.user.avatar_url || undefined}
		role={spaceStore.active?.role}
	/>
{/if}
