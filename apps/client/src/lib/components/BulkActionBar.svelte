<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Trash2, Archive, MailOpen, MailX, X } from 'lucide-svelte';

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
</script>

{#if count > 0}
	<div class="bulk-action-bar flex items-center gap-1.5 border-b bg-muted/50 px-3 py-1.5">
		<span class="shrink-0 text-xs font-medium">{count}</span>
		<div class="flex items-center gap-0.5 overflow-hidden">
			<Button variant="ghost" size="icon" class="h-7 w-7 shrink-0" onclick={onarchive} disabled={loading} title="Archive">
				<Archive class="h-3.5 w-3.5" />
			</Button>
			<Button variant="ghost" size="icon" class="h-7 w-7 shrink-0" onclick={onmarkread} disabled={loading} title="Mark read">
				<MailOpen class="h-3.5 w-3.5" />
			</Button>
			<Button variant="ghost" size="icon" class="h-7 w-7 shrink-0" onclick={onmarkunread} disabled={loading} title="Mark unread">
				<MailX class="h-3.5 w-3.5" />
			</Button>
			<Button variant="ghost" size="icon" class="h-7 w-7 shrink-0 text-destructive hover:text-destructive hover:bg-destructive/10" onclick={ondelete} disabled={loading} title="Delete">
				<Trash2 class="h-3.5 w-3.5" />
			</Button>
		</div>
		<div class="ml-auto">
			<Button variant="ghost" size="icon" class="h-6 w-6 shrink-0" onclick={onclear}>
				<X class="h-3 w-3" />
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
