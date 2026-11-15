<script lang="ts">
	import { Button, Field, Modal, Select, Textarea, toast } from '@facile/muse';
	import { CAPSULE_EXPIRIES, sealLink, type CapsuleExpiry } from '$lib/capsule';
	import { t } from '$lib/i18n/index.svelte';

	let {
		open = $bindable(false),
		icon,
		oninsert
	}: { open?: boolean; icon: string; oninsert: (link: string) => void } = $props();

	let secret = $state('');
	let expiry = $state<CapsuleExpiry>('24h');
	let sealing = $state(false);

	$effect(() => {
		if (open) {
			secret = '';
			expiry = '24h';
		}
	});

	async function seal(e: SubmitEvent) {
		e.preventDefault();
		if (!secret.trim() || sealing) return;
		sealing = true;
		try {
			oninsert(await sealLink(secret, expiry));
			secret = '';
			open = false;
			toast.success(t('compose.capsuleAdded'));
		} catch {
			toast.danger(t('compose.capsuleFailed'));
		}
		sealing = false;
	}
</script>

<Modal bind:open title={t('compose.capsuleTitle')} showClose>
	<form class="flex flex-col gap-4" onsubmit={seal}>
		<p class="text-fc-sm text-fc-fg-muted">{t('compose.capsuleBody')}</p>
		<Field label={t('compose.capsuleSecret')}>
			<Textarea bind:value={secret} rows={3} placeholder={t('compose.capsulePlaceholder')} autocomplete="off" spellcheck={false} />
		</Field>
		<Field label={t('compose.capsuleExpiry')}>
			<Select
				value={expiry}
				options={CAPSULE_EXPIRIES.map((value) => ({ value, label: t(`compose.capsule${value}`) }))}
				onchange={(e) => (expiry = e.currentTarget.value as CapsuleExpiry)}
			/>
		</Field>
		<div class="flex justify-end gap-2">
			<Button type="button" variant="ghost" onclick={() => (open = false)}>{t('common.cancel')}</Button>
			<Button type="submit" {icon} disabled={sealing || !secret.trim()}>
				{sealing ? t('compose.capsuleSealing') : t('compose.capsuleInsert')}
			</Button>
		</div>
	</form>
</Modal>
