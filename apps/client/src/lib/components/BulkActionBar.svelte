<script lang="ts">
	import { Badge, IconButton, icons } from '@facile/muse';
	import { t } from '$lib/i18n/index.svelte';

	let {
		count = 0,
		loading = false,
		ondelete,
		onarchive,
		onmarkread,
		onmarkunread,
		onclear
	}: {
		count: number;
		loading?: boolean;
		ondelete: () => void;
		onarchive: () => void;
		onmarkread: () => void;
		onmarkunread: () => void;
		onclear: () => void;
	} = $props();

	/* Glyphs an email client needs that muse's `icons` map has no key for yet. */
	const mailIcons = {
		archive: 'solar:archive-linear',
		mailOpen: 'solar:letter-opened-linear',
		mailUnread: 'solar:letter-unread-linear'
	};
</script>

{#if count > 0}
	<div class="bulk-action-bar flex items-center gap-2 bg-fc-surface px-3 py-1.5">
		<Badge tone="accent" class="shrink-0 tabular-nums">{count}</Badge>
		<div class="flex min-w-0 items-center gap-0.5">
			<IconButton
				variant="ghost"
				aria-label={t('mail.archive')}
				title={t('mail.archive')}
				disabled={loading}
				onclick={onarchive}
			>
				<iconify-icon icon={mailIcons.archive} width="18" height="18" class="block size-4.5"
				></iconify-icon>
			</IconButton>
			<IconButton
				variant="ghost"
				aria-label={t('mail.markRead')}
				title={t('mail.markRead')}
				disabled={loading}
				onclick={onmarkread}
			>
				<iconify-icon icon={mailIcons.mailOpen} width="18" height="18" class="block size-4.5"
				></iconify-icon>
			</IconButton>
			<IconButton
				variant="ghost"
				aria-label={t('mail.markUnread')}
				title={t('mail.markUnread')}
				disabled={loading}
				onclick={onmarkunread}
			>
				<iconify-icon icon={mailIcons.mailUnread} width="18" height="18" class="block size-4.5"
				></iconify-icon>
			</IconButton>
			<IconButton
				variant="ghost-danger"
				aria-label={t('mail.delete')}
				title={t('mail.delete')}
				disabled={loading}
				onclick={ondelete}
			>
				<iconify-icon icon={icons.remove} width="18" height="18" class="block size-4.5"
				></iconify-icon>
			</IconButton>
		</div>
		<IconButton variant="ghost" aria-label={t('mail.clearSelection')} class="ml-auto" onclick={onclear}>
			<iconify-icon icon={icons.close} width="18" height="18" class="block size-4.5"></iconify-icon>
		</IconButton>
	</div>
{/if}

<style>
	@media (prefers-reduced-motion: no-preference) {
		.bulk-action-bar {
			animation: bar-slide-in 150ms var(--ease-fc, ease-out) both;
		}
	}

	@keyframes bar-slide-in {
		from {
			opacity: 0;
			transform: translateY(-100%);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}
</style>
