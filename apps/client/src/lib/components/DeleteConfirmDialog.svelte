<script lang="ts">
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import { Trash2 } from 'lucide-svelte';

	let {
		open = $bindable(false),
		count = 1,
		permanent = false,
		onconfirm
	}: {
		open: boolean;
		count?: number;
		permanent?: boolean;
		onconfirm: () => void;
	} = $props();

	function handleConfirm() {
		onconfirm();
		open = false;
	}
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Portal>
		<AlertDialog.Overlay class="fixed inset-0 z-50 bg-black/60 backdrop-blur-sm data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0" />
		<AlertDialog.Content class="fixed left-[50%] top-[50%] z-50 grid w-full max-w-md translate-x-[-50%] translate-y-[-50%] gap-6 rounded-2xl border bg-background p-8 shadow-2xl duration-200 data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[state=closed]:slide-out-to-left-1/2 data-[state=closed]:slide-out-to-top-[48%] data-[state=open]:slide-in-from-left-1/2 data-[state=open]:slide-in-from-top-[48%]">
			<AlertDialog.Header class="flex flex-col items-center gap-4 text-center">
				<div class="flex h-16 w-16 items-center justify-center rounded-full bg-destructive/10">
					<Trash2 class="h-8 w-8 text-destructive" />
				</div>
				<AlertDialog.Title class="text-xl font-semibold">
					{#if permanent}
						Permanently delete {count === 1 ? 'this email' : `${count} emails`}?
					{:else}
						Delete {count === 1 ? 'this email' : `${count} emails`}?
					{/if}
				</AlertDialog.Title>
				<AlertDialog.Description class="text-sm text-muted-foreground">
					{#if permanent}
						This action cannot be undone. {count === 1 ? 'This email' : 'These emails'} will be permanently removed.
					{:else}
						{count === 1 ? 'This email' : 'These emails'} will be moved to Trash. You can restore {count === 1 ? 'it' : 'them'} from there.
					{/if}
				</AlertDialog.Description>
			</AlertDialog.Header>
			<AlertDialog.Footer class="flex justify-center gap-3 sm:justify-center">
				<AlertDialog.Cancel class="min-w-[100px]">Cancel</AlertDialog.Cancel>
				<AlertDialog.Action
					class="min-w-[100px] gap-2 bg-destructive text-destructive-foreground hover:bg-destructive/90"
					onclick={handleConfirm}
				>
					<Trash2 class="h-4 w-4" />
					Delete
				</AlertDialog.Action>
			</AlertDialog.Footer>
		</AlertDialog.Content>
	</AlertDialog.Portal>
</AlertDialog.Root>
