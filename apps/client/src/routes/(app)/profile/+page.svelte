<script lang="ts">
	import { getContext } from 'svelte';
	import type { UserProfile } from '$lib/backend';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';

	const ctx = getContext<{
		token: string;
		user: UserProfile | null;
		setUser: (user: UserProfile) => void;
	}>('app');

	let name = $state(ctx.user?.name ?? '');

	$effect(() => {
		name = ctx.user?.name ?? '';
	});

	function getInitials(value: string) {
		const parts = value.trim().split(/\s+/).filter(Boolean);
		if (parts.length === 0) return '?';
		if (parts.length === 1) return parts[0].slice(0, 2).toUpperCase();
		return `${parts[0][0] ?? ''}${parts[1][0] ?? ''}`.toUpperCase();
	}

	function displayName(user: UserProfile | null) {
		return user?.name?.trim() || user?.email || '';
	}
</script>

<svelte:head>
	<title>Profile — Courrier</title>
</svelte:head>

<div class="flex flex-col gap-6 p-6">
	<div class="space-y-2">
		<h1 class="text-2xl font-semibold">Profile</h1>
		<p class="text-sm text-muted-foreground">Your account details.</p>
	</div>

	<Card.Root class="max-w-2xl">
		<Card.Header>
			<Card.Title>Identity</Card.Title>
		</Card.Header>
		<Card.Content class="space-y-6">
			<div class="flex items-center gap-4">
				{#if ctx.user?.avatar_url}
					<img
						src={ctx.user.avatar_url}
						alt={displayName(ctx.user)}
						class="h-24 w-24 rounded-full border border-border object-cover"
					/>
				{:else}
					<div class="flex h-24 w-24 items-center justify-center rounded-full border border-border bg-foreground text-2xl font-semibold text-background">
						{getInitials(displayName(ctx.user))}
					</div>
				{/if}
				{#if ctx.user?.avatar_source === 'oidc'}
					<p class="text-xs text-muted-foreground">Avatar synced from SSO</p>
				{/if}
			</div>

			<div class="space-y-4">
				<div class="space-y-2">
					<Label for="name">Name</Label>
					<Input id="name" value={name} disabled />
				</div>

				<div class="space-y-2">
					<Label for="email">Email</Label>
					<Input id="email" value={ctx.user?.email ?? ''} disabled />
				</div>
			</div>
		</Card.Content>
	</Card.Root>
</div>
