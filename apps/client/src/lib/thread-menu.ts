import { goto } from '$app/navigation';
import { icons, toast, type MenuEntry } from '@facile/muse';
import type { EmailAttachment, EmailMessage } from '$lib/backend';
import { t } from '$lib/i18n/index.svelte';

const glyphs = {
	reply: 'solar:reply-linear',
	replyAll: 'solar:multiple-forward-left-linear',
	forward: 'solar:forward-linear'
};

async function copy(text: string) {
	try {
		await navigator.clipboard.writeText(text);
		toast.success(t('mail.copied'));
	} catch {
		toast.danger(t('mail.copyFailed'));
	}
}

/** The right-click menu of one message in a thread: answer it, or copy who sent it. */
export function messageMenu(msg: EmailMessage, accountId: number): MenuEntry[] {
	const open = (kind: string) => () => goto(`/mail/compose?${kind}=${msg.id}&accountId=${accountId}`);
	return [
		{ label: t('mail.reply'), icon: glyphs.reply, onselect: open('reply') },
		{ label: t('mail.replyAll'), icon: glyphs.replyAll, onselect: open('replyall') },
		{ label: t('mail.forward'), icon: glyphs.forward, onselect: open('forward') },
		{ separator: true },
		{ label: t('mail.copyAddress'), icon: icons.copy, onselect: () => copy(msg.from_address) }
	];
}

/** The right-click menu of an attachment chip. */
export function attachmentMenu(attachment: EmailAttachment, download: () => void): MenuEntry[] {
	return [
		{ label: t('mail.download'), icon: icons.download, onselect: download },
		{ label: t('mail.copyName'), icon: icons.copy, onselect: () => copy(attachment.filename) }
	];
}
