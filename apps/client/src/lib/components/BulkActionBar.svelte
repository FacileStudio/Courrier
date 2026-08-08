<script lang="ts">
	import { Badge, IconButton, icons } from '@facile/muse';

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
	<div class="bulk-action-bar flex items-center gap-2 border-b border-fc-border bg-fc-surface px-3 py-1.5">
		<Badge tone="accent" class="shrink-0 tabular-nums">{count}</Badge>
		<div class="flex min-w-0 items-center gap-0.5">
			<IconButton
				variant="ghost"
				aria-label="Archive"
				title="Archive"
				disabled={loading}
				onclick={onarchive}
			>
				<iconify-icon icon={mailIcons.archive} width="18" height="18" class="block size-4.5"
				></iconify-icon>
			</IconButton>
			<IconButton
				variant="ghost"
				aria-label="Mark as read"
				title="Mark as read"
				disabled={loading}
				onclick={onmarkread}
			>
				<iconify-icon icon={mailIcons.mailOpen} width="18" height="18" class="block size-4.5"
				></iconify-icon>
			</IconButton>
			<IconButton
				variant="ghost"
				aria-label="Mark as unread"
				title="Mark as unread"
				disabled={loading}
				onclick={onmarkunread}
			>
				<iconify-icon icon={mailIcons.mailUnread} width="18" height="18" class="block size-4.5"
				></iconify-icon>
			</IconButton>
			<IconButton
				variant="danger"
				aria-label="Delete"
				title="Delete"
				disabled={loading}
				onclick={ondelete}
			>
				<iconify-icon icon={icons.remove} width="18" height="18" class="block size-4.5"
				></iconify-icon>
			</IconButton>
		</div>
		<IconButton variant="ghost" aria-label="Clear selection" class="ml-auto" onclick={onclear}>
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
