<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Trash2, Archive, MailOpen, MailX, X } from 'lucide-svelte';

	let {
		count = 0,
		ondelete,
		onarchive,
		onmarkread,
		onmarkunread,
		onclear
	}: {
		count: number;
		ondelete: () => void;
		onarchive: () => void;
		onmarkread: () => void;
		onmarkunread: () => void;
		onclear: () => void;
	} = $props();
</script>

{#if count > 0}
	<div class="bulk-action-bar flex items-center gap-2 border-b bg-muted/50 px-4 py-2">
		<span class="text-sm font-medium">{count} selected</span>
		<div class="ml-auto flex items-center gap-1">
			<Button variant="ghost" size="sm" class="gap-1.5 text-xs" onclick={onarchive}>
				<Archive class="h-3.5 w-3.5" />
				Archive
			</Button>
			<Button variant="ghost" size="sm" class="gap-1.5 text-xs" onclick={onmarkread}>
				<MailOpen class="h-3.5 w-3.5" />
				Read
			</Button>
			<Button variant="ghost" size="sm" class="gap-1.5 text-xs" onclick={onmarkunread}>
				<MailX class="h-3.5 w-3.5" />
				Unread
			</Button>
			<Button variant="ghost" size="sm" class="gap-1.5 text-xs text-destructive hover:text-destructive hover:bg-destructive/10" onclick={ondelete}>
				<Trash2 class="h-3.5 w-3.5" />
				Delete
			</Button>
			<div class="ml-1 h-4 w-px bg-border"></div>
			<Button variant="ghost" size="icon" class="h-7 w-7" onclick={onclear}>
				<X class="h-3.5 w-3.5" />
			</Button>
		</div>
	</div>
{/if}

<style>
	@media (prefers-reduced-motion: no-preference) {
		.bulk-action-bar {
			animation: bar-slide-in 150ms ease-out both;
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
