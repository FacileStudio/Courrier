export type MailFolderLink = { href: string; label: string; icon: string; type: string };

export const MAIL_FOLDERS: MailFolderLink[] = [
	{ href: '/mail', label: 'Inbox', icon: 'solar:inbox-linear', type: 'inbox' },
	{ href: '/mail/sent', label: 'Sent', icon: 'solar:plain-2-linear', type: 'sent' },
	{ href: '/mail/drafts', label: 'Drafts', icon: 'solar:file-text-linear', type: 'drafts' },
	{ href: '/mail/archive', label: 'Archive', icon: 'solar:archive-linear', type: 'archive' },
	{ href: '/mail/junk', label: 'Junk', icon: 'solar:danger-triangle-linear', type: 'junk' },
	{ href: '/mail/trash', label: 'Trash', icon: 'solar:trash-bin-trash-linear', type: 'trash' }
];
