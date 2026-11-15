<script lang="ts">
	import { Avatar, Badge, Checkbox, ContextMenu, icons, type MenuEntry } from '@facile/muse';
	import type { EmailMessage } from '$lib/backend';
	import { formatListDate } from '$lib/i18n/format';
	import { t } from '$lib/i18n/index.svelte';

	/* Glyphs an email client needs that muse's `icons` map has no key for yet. */
	const mailIcons = {
		reply: 'solar:reply-linear',
		forward: 'solar:forward-linear',
		archive: 'solar:archive-linear',
		mailOpen: 'solar:letter-opened-linear',
		star: 'solar:star-linear'
	};

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

	const sender = $derived(email.from_name || email.from_address);

	function handleClick() {
		if (selectionActive) {
			ontogglecheck();
			return;
		}
		onopen();
	}

	function handleCheckboxClick(e: Event) {
		e.stopPropagation();
		ontogglecheck();
	}

	const menuItems = $derived.by(() => {
		const groups: MenuEntry[][] = [
			[
				...(onreply ? [{ label: t('mail.reply'), icon: mailIcons.reply, onselect: onreply }] : []),
				...(onforward ? [{ label: t('mail.forward'), icon: mailIcons.forward, onselect: onforward }] : [])
			],
			[
				...(ontoggleread
					? [{ label: email.is_read ? t('mail.markUnread') : t('mail.markRead'), icon: email.is_read ? icons.mail : mailIcons.mailOpen, onselect: ontoggleread }]
					: []),
				...(ontogglestar ? [{ label: email.is_starred ? t('mail.unstar') : t('mail.star'), icon: mailIcons.star, onselect: ontogglestar }] : [])
			],
			[
				...(onarchive ? [{ label: t('mail.archive'), icon: mailIcons.archive, onselect: onarchive }] : []),
				...(ondelete ? [{ label: t('mail.delete'), icon: icons.remove, tone: 'danger' as const, onselect: ondelete }] : [])
			]
		];
		return groups.filter((g) => g.length > 0).flatMap((g, i) => (i === 0 ? g : [{ separator: true } as const, ...g]));
	});
</script>

<ContextMenu label={t('mail.rowActions')} items={menuItems} class="relative">
	<button
		type="button"
		class="group flex w-full items-center gap-3 px-4 py-2.5 text-left transition-colors focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-fc-ring
			{selected ? 'bg-fc-accent text-fc-accent-fg' : checked ? 'bg-fc-surface text-fc-fg' : 'text-fc-fg hover:bg-fc-surface'}"
		onclick={handleClick}
	>
		<span class="flex size-8 shrink-0 items-center justify-center">
			{#if selectionActive || checked}
				<span
					class="flex size-8 cursor-pointer items-center justify-center"
					onclick={handleCheckboxClick}
					onkeydown={(e) => {
						if (e.key === 'Enter' || e.key === ' ') handleCheckboxClick(e);
					}}
					role="checkbox"
					aria-checked={checked}
					aria-label={t('mail.selectConversation')}
					tabindex={0}
				>
					<Checkbox {checked} class="pointer-events-none" />
				</span>
			{:else}
				<span class="block group-hover:hidden">
					<Avatar name={sender} size="sm" />
				</span>
				<span
					class="hidden size-8 cursor-pointer items-center justify-center group-hover:flex"
					onclick={handleCheckboxClick}
					onkeydown={(e) => {
						if (e.key === 'Enter' || e.key === ' ') handleCheckboxClick(e);
					}}
					role="checkbox"
					aria-checked={checked}
					aria-label={t('mail.selectConversation')}
					tabindex={0}
				>
					<Checkbox {checked} class="pointer-events-none" />
				</span>
			{/if}
		</span>

		<span class="min-w-0 flex-1">
			<span class="flex items-center justify-between gap-2">
				<span class="flex min-w-0 items-center gap-1.5">
					<span class="truncate text-fc-sm {email.is_read ? '' : 'font-semibold'}">{sender}</span>
					{#if (email.message_count ?? 1) > 1}
						<Badge tone="neutral" class="tabular-nums">{email.message_count}</Badge>
					{/if}
				</span>
				<span class="shrink-0 text-fc-xs {selected ? 'text-fc-accent-fg/70' : 'text-fc-fg-muted'}">
					{formatListDate(email.date)}
				</span>
			</span>
			<span
				class="mt-0.5 block truncate text-fc-sm {selected ? 'text-fc-accent-fg/70' : 'text-fc-fg-muted'}"
			>
				{email.subject || t('common.noSubject')}
			</span>
		</span>

		{#if !email.is_read}
			<span
				class="size-2 shrink-0 rounded-fc-pill {selected ? 'bg-fc-accent-fg' : 'bg-fc-accent'}"
				aria-label={t('mail.unreadMark')}
			></span>
		{/if}
	</button>
</ContextMenu>
