<script lang="ts">
	import { onMount } from 'svelte';
	import { Alert, Button, Spinner, StatusDot, icons, toast } from '@facile/muse';
	import { backend, type PushConfig } from '$lib/backend';
	import { enable, permission, supported, type Permission } from '$lib/notifications';
	import { isDesktop, isNative } from '$lib/platform';

	let { showTest = false }: { showTest?: boolean } = $props();

	let push = $state<PushConfig | null>(null);
	let access = $state<Permission | 'loading'>('loading');
	let busy = $state(false);
	let testing = $state(false);
	const native = isNative();
	const desktop = isDesktop();

	onMount(async () => {
		try {
			push = await backend.getPush();
		} catch {
			push = null;
		}
		access = await permission().catch(() => 'unsupported' as const);
	});

	async function turnOn() {
		busy = true;
		try {
			access = await enable(push);
			if (access === 'granted') toast.success('Notifications are on');
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : 'Could not turn on notifications');
		} finally {
			busy = false;
		}
	}

	async function sendTest() {
		testing = true;
		try {
			await backend.testPush();
			toast.success('Test notification sent');
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : 'Could not send the test');
		} finally {
			testing = false;
		}
	}
</script>

<div class="flex flex-col gap-4" aria-live="polite">
	{#if access === 'loading'}
		<Spinner size="sm" />
	{:else if desktop}
		<p class="text-fc-sm text-fc-fg-muted">
			The desktop app does not show notifications yet. The Android app and a browser tab do.
		</p>
	{:else if !supported(push)}
		<p class="text-fc-sm text-fc-fg-muted">
			{native ? 'Phone notifications' : 'Browser notifications'} are not set up on this server yet.
		</p>
	{:else if access === 'unsupported'}
		<p class="text-fc-sm text-fc-fg-muted">This browser cannot show notifications.</p>
	{:else if access === 'denied'}
		<Alert tone="warning">
			{native
				? 'Notifications are blocked. Turn them on in Settings > Apps > Courrier > Notifications.'
				: 'Notifications are blocked for this site. Allow them in your browser settings.'}
		</Alert>
	{:else if access === 'prompt'}
		<div>
			<Button variant="primary" class="min-h-11" icon={busy ? undefined : icons.notification} disabled={busy} onclick={turnOn}>
				{#if busy}<Spinner size="sm" />{/if}
				Turn on notifications
			</Button>
		</div>
	{:else}
		<StatusDot tone="success" label="On for this device" />
	{/if}

	{#if showTest && supported(push) && access === 'granted'}
		<div>
			<Button variant="outline" class="min-h-11" disabled={testing} onclick={sendTest}>
				{#if testing}<Spinner size="sm" />{/if}
				Send a test
			</Button>
		</div>
	{/if}
</div>
