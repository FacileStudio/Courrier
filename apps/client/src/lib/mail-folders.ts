import type { AppKey } from '$lib/i18n/messages';

export type MailFolderLink = { href: string; labelKey: AppKey; icon: string; type: string; primary?: boolean };

export const MAIL_FOLDERS: MailFolderLink[] = [
	{ href: '/mail', labelKey: 'folder.inbox', icon: 'solar:inbox-linear', type: 'inbox', primary: true },
	{ href: '/mail/sent', labelKey: 'folder.sent', icon: 'solar:plain-2-linear', type: 'sent', primary: true },
	{ href: '/mail/drafts', labelKey: 'folder.drafts', icon: 'solar:file-text-linear', type: 'drafts', primary: true },
	{ href: '/mail/archive', labelKey: 'folder.archive', icon: 'solar:archive-linear', type: 'archive', primary: true },
	{ href: '/mail/junk', labelKey: 'folder.junk', icon: 'solar:danger-triangle-linear', type: 'junk' },
	{ href: '/mail/trash', labelKey: 'folder.trash', icon: 'solar:trash-bin-trash-linear', type: 'trash' }
];
