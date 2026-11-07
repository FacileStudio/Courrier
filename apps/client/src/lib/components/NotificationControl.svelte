<script lang="ts">
	import { onMount } from 'svelte';
	import { Alert, Button, Spinner, StatusDot, icons, toast } from '@facile/muse';
	import { backend, type PushConfig } from '$lib/backend';
	import { enable, permission, supported, type Permission } from '$lib/notifications';
	import { isDesktop, isNative } from '$lib/platform';
	import { t } from '$lib/i18n/index.svelte';

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
			if (access === 'granted') toast.success(t('notify.on'));
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : t('notify.enableFailed'));
		} finally {
			busy = false;
		}
	}

	async function sendTest() {
		testing = true;
		try {
			await backend.testPush();
			toast.success(t('notify.testSent'));
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : t('notify.testFailed'));
		} finally {
			testing = false;
		}
	}
</script>

<div class="flex flex-col gap-4" aria-live="polite">
	{#if access === 'loading'}
		<Spinner size="sm" />
	{:else if !supported(push)}
		<p class="text-fc-sm text-fc-fg-muted">
			{native ? t('notify.notSetUpPhone') : t('notify.notSetUpBrowser')}
		</p>
	{:else if access === 'unsupported'}
		<p class="text-fc-sm text-fc-fg-muted">{t('notify.unsupported')}</p>
	{:else if access === 'denied'}
		<Alert tone="warning">
			{native
				? t('notify.blockedNative')
				: desktop
					? t('notify.blockedDesktop')
					: t('notify.blockedBrowser')}
		</Alert>
	{:else if access === 'prompt'}
		<div>
			<Button variant="primary" class="min-h-11" icon={busy ? undefined : icons.notification} disabled={busy} onclick={turnOn}>
				{#if busy}<Spinner size="sm" />{/if}
				{t('notify.turnOn')}
			</Button>
		</div>
	{:else}
		<StatusDot tone="success" label={t('notify.onThisDevice')} />
	{/if}

	{#if showTest && supported(push) && access === 'granted'}
		<div>
			<Button variant="outline" class="min-h-11" disabled={testing} onclick={sendTest}>
				{#if testing}<Spinner size="sm" />{/if}
				{t('notify.sendTest')}
			</Button>
		</div>
	{/if}
</div>
