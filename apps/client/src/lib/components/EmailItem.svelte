<script lang="ts">
	import * as ContextMenu from '$lib/components/ui/context-menu';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Reply, Forward, Archive, Trash2, MailOpen, Mail, Star, StarOff } from 'lucide-svelte';
	import type { EmailMessage } from '$lib/backend';

	let {
		email,
		selected = false,
		checked = false,
		selectionActive = false,
		onopen,
		ontogglecheck,
		onreply,
		onforward,
		onarchive,
		ondelete,
		ontoggleread,
		ontogglestar
	}: {
		email: EmailMessage;
		selected?: boolean;
		checked?: boolean;
		selectionActive?: boolean;
		onopen: () => void;
		ontogglecheck: () => void;
		onreply?: () => void;
		onforward?: () => void;
		onarchive?: () => void;
		ondelete?: () => void;
		ontoggleread?: () => void;
		ontogglestar?: () => void;
	} = $props();

	function getInitials(name: string) {
		const parts = name.trim().split(/\s+/).filter(Boolean);
		if (parts.length === 0) return '?';
		if (parts.length === 1) return parts[0][0].toUpperCase();
		return `${parts[0][0]}${parts[1][0]}`.toUpperCase();
	}

	function formatDate(dateStr: string) {
		const date = new Date(dateStr);
		const now = new Date();
		const diffMs = now.getTime() - date.getTime();
		const diffDays = diffMs / (1000 * 60 * 60 * 24);

		if (diffDays < 1) {
			return date.toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' });
		}
		if (diffDays < 7) {
			return date.toLocaleDateString('fr-FR', { weekday: 'short' });
		}
		return date.toLocaleDateString('fr-FR', { day: 'numeric', month: 'short' });
	}

	function handleClick(e: MouseEvent) {
		if (selectionActive) {
			ontogglecheck();
			return;
		}
		onopen();
	}

	function handleCheckboxClick(e: MouseEvent) {
		e.stopPropagation();
		ontogglecheck();
	}
</script>

<ContextMenu.Root>
	<ContextMenu.Trigger>
		<!-- svelte-ignore node_invalid_placement_ssr -->
		<button
			class="mail-list-item group w-full text-left px-4 py-3 border-b transition-colors duration-150 hover:bg-muted/50 cursor-pointer
				{selected ? 'bg-muted' : ''}
				{checked ? 'bg-primary/5' : ''}
				{!email.is_read ? 'font-medium' : ''}"
			onclick={handleClick}
		>
			<div class="flex items-center gap-3">
				<div class="relative flex h-8 w-8 shrink-0 items-center justify-center">
					{#if selectionActive || checked}
						<span
							class="flex h-8 w-8 items-center justify-center cursor-pointer"
							onclick={handleCheckboxClick}
							onkeydown={(e) => { if (e.key === 'Enter' || e.key === ' ') handleCheckboxClick(e as any); }}
							role="checkbox"
							aria-checked={checked}
							tabindex={0}
						>
							<Checkbox
								checked={checked}
								class="h-4 w-4 pointer-events-none"
							/>
						</span>
					{:else}
						<div
							class="relative flex h-8 w-8 items-center justify-center"
						>
							<div class="flex h-8 w-8 items-center justify-center rounded-full bg-muted text-xs font-medium group-hover:hidden">
								{getInitials(email.from_name || email.from_address)}
							</div>
							<span
								class="hidden h-8 w-8 items-center justify-center group-hover:flex cursor-pointer"
								onclick={handleCheckboxClick}
								onkeydown={(e) => { if (e.key === 'Enter' || e.key === ' ') handleCheckboxClick(e as any); }}
								role="checkbox"
								aria-checked={checked}
								tabindex={0}
							>
								<Checkbox
									checked={checked}
									class="h-4 w-4 pointer-events-none"
								/>
							</span>
						</div>
					{/if}
				</div>
				<div class="min-w-0 flex-1">
					<div class="flex items-center justify-between gap-2">
						<span class="truncate text-sm">{email.from_name || email.from_address}</span>
						<span class="shrink-0 text-xs text-muted-foreground">{formatDate(email.date)}</span>
					</div>
					<p class="truncate text-sm text-muted-foreground">{email.subject || '(no subject)'}</p>
				</div>
				{#if !email.is_read}
					<div class="h-2 w-2 shrink-0 rounded-full bg-blue-600"></div>
				{/if}
			</div>
		</button>
	</ContextMenu.Trigger>
	<ContextMenu.Content class="w-56">
		{#if onreply}
			<ContextMenu.Item class="gap-2" onclick={onreply}>
				<Reply class="h-4 w-4" />
				Reply
			</ContextMenu.Item>
		{/if}
		{#if onforward}
			<ContextMenu.Item class="gap-2" onclick={onforward}>
				<Forward class="h-4 w-4" />
				Forward
			</ContextMenu.Item>
		{/if}
		<ContextMenu.Separator />
		{#if ontoggleread}
			<ContextMenu.Item class="gap-2" onclick={ontoggleread}>
				{#if email.is_read}
					<Mail class="h-4 w-4" />
					Mark as unread
				{:else}
					<MailOpen class="h-4 w-4" />
					Mark as read
				{/if}
			</ContextMenu.Item>
		{/if}
		{#if ontogglestar}
			<ContextMenu.Item class="gap-2" onclick={ontogglestar}>
				{#if email.is_starred}
					<StarOff class="h-4 w-4" />
					Unstar
				{:else}
					<Star class="h-4 w-4" />
					Star
				{/if}
			</ContextMenu.Item>
		{/if}
		<ContextMenu.Separator />
		{#if onarchive}
			<ContextMenu.Item class="gap-2" onclick={onarchive}>
				<Archive class="h-4 w-4" />
				Archive
			</ContextMenu.Item>
		{/if}
		{#if ondelete}
			<ContextMenu.Item class="gap-2 text-destructive focus:text-destructive" onclick={ondelete}>
				<Trash2 class="h-4 w-4" />
				Delete
			</ContextMenu.Item>
		{/if}
	</ContextMenu.Content>
</ContextMenu.Root>
