const backendBaseUrl = (import.meta.env.VITE_API_BASE_URL as string | undefined)?.replace(/\/$/, '') || '';

export type AuthResponse = {
	user_id: string;
	token: string;
};

export type UserProfile = {
	id: string;
	email: string;
	name: string;
	avatar_url: string;
	avatar_source: string;
	created_at: string;
};

export type MeResponse = {
	user: UserProfile;
};

export type MailAccount = {
	id: number;
	name: string;
	email: string;
	imap_host: string;
	imap_port: number;
	imap_user: string;
	smtp_host: string;
	smtp_port: number;
	smtp_user: string;
	signature: string;
	is_default: boolean;
	created_at: string;
	updated_at: string;
};

export type Folder = {
	id: number;
	account_id: number;
	path: string;
	name: string;
	type: string;
	unread_count: number;
	total_count: number;
};

export type EmailAddress = {
	name: string;
	email: string;
};

export type EmailMessage = {
	id: number;
	account_id: number;
	folder_id: number;
	message_id: string;
	thread_id: string;
	subject: string;
	from_address: string;
	from_name: string;
	to_addresses: EmailAddress[];
	cc_addresses: EmailAddress[];
	date: string;
	body_text: string;
	body_html: string;
	is_read: boolean;
	is_starred: boolean;
	has_attachments: boolean;
	attachments?: EmailAttachment[];
	in_reply_to?: string;
	references?: string;
};

export type EmailAttachment = {
	id: number;
	filename: string;
	mime_type: string;
	size: number;
};

export type EmailTemplate = {
	id: number;
	name: string;
	subject: string;
	body_html: string;
	body_text: string;
	created_at: string;
	updated_at: string;
};

type ApiErrorPayload = {
	error?: { message?: string };
};

async function apiFetch<T>(path: string, options: RequestInit = {}) {
	const headers = new Headers(options.headers);
	if (!headers.has('Content-Type') && options.body) {
		headers.set('Content-Type', 'application/json');
	}
	const response = await fetch(`${backendBaseUrl}${path}`, { ...options, headers, credentials: 'include' });
	if (!response.ok) {
		let payload: ApiErrorPayload | undefined;
		try {
			payload = (await response.json()) as ApiErrorPayload;
		} catch {
			payload = undefined;
		}
		throw new Error(payload?.error?.message || `Request failed with status ${response.status}`);
	}
	return (await response.json()) as T;
}

function resolveFileUrl(path: string) {
	if (!path) return '';
	if (/^https?:\/\//.test(path)) return path;
	return `${backendBaseUrl}${path.startsWith('/') ? path : `/${path}`}`;
}

function normalizeUser(user: UserProfile): UserProfile {
	return { ...user, avatar_url: resolveFileUrl(user.avatar_url) };
}

export const backend = {
	baseUrl: backendBaseUrl,

	register(email: string, password: string) {
		return apiFetch<AuthResponse>('/auth/register', {
			method: 'POST',
			body: JSON.stringify({ email, password })
		});
	},
	login(email: string, password: string) {
		return apiFetch<AuthResponse>('/auth/login', {
			method: 'POST',
			body: JSON.stringify({ email, password })
		});
	},
	logout() {
		return apiFetch<{ ok: boolean }>('/auth/logout', { method: 'POST' });
	},
	me() {
		return apiFetch<MeResponse>('/users/me').then((r) => ({
			user: normalizeUser(r.user)
		}));
	},

	listAccounts() {
		return apiFetch<{ accounts: MailAccount[] }>('/accounts');
	},
	getAccount(id: number) {
		return apiFetch<MailAccount>(`/accounts/${id}`);
	},
	createAccount(data: Omit<MailAccount, 'id' | 'created_at' | 'updated_at'> & { imap_password: string; smtp_password: string }) {
		return apiFetch<MailAccount>('/accounts', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	},
	updateAccount(id: number, data: Partial<MailAccount & { imap_password: string; smtp_password: string }>) {
		return apiFetch<MailAccount>(`/accounts/${id}`, {
			method: 'PUT',
			body: JSON.stringify(data)
		});
	},
	deleteAccount(id: number) {
		return apiFetch<{ deleted: boolean }>(`/accounts/${id}`, { method: 'DELETE' });
	},

	syncProfile() {
		return apiFetch<{ synced: boolean }>('/auth/sync-profile', { method: 'POST' });
	},

	syncAccount(accountId: number) {
		return apiFetch<{ synced: boolean }>(`/accounts/${accountId}/mail/sync`, { method: 'POST' });
	},
	syncFolder(accountId: number, folderId: number) {
		return apiFetch<{ synced: boolean }>(`/accounts/${accountId}/mail/folders/${folderId}/sync`, { method: 'POST' });
	},
	getFolders(accountId: number) {
		return apiFetch<{ folders: Folder[] }>(`/accounts/${accountId}/mail/folders`);
	},
	getEmailsByFolder(accountId: number, folderType: string, page = 1, limit = 50, unreadOnly = false) {
		const params = new URLSearchParams({ page: String(page), limit: String(limit) });
		if (unreadOnly) params.set('unread', 'true');
		return apiFetch<{ emails: EmailMessage[]; total: number; page: number; limit: number }>(
			`/accounts/${accountId}/mail/folders/${folderType}/emails?${params}`
		);
	},
	getEmail(accountId: number, emailId: number) {
		return apiFetch<EmailMessage>(`/accounts/${accountId}/mail/emails/${emailId}`);
	},
	updateEmail(accountId: number, emailId: number, data: { is_read?: boolean; is_starred?: boolean }) {
		return apiFetch<EmailMessage>(`/accounts/${accountId}/mail/emails/${emailId}`, {
			method: 'PATCH',
			body: JSON.stringify(data)
		});
	},
	sendEmail(accountId: number, data: { to: string[]; cc?: string[]; subject: string; body: string; body_html?: string; in_reply_to?: string; references?: string[] }) {
		return apiFetch<{ sent: boolean }>(`/accounts/${accountId}/mail/send`, {
			method: 'POST',
			body: JSON.stringify(data)
		});
	},
	testConnection(data: { imap_host: string; imap_port: number; imap_user: string; imap_password: string; smtp_host: string; smtp_port: number; smtp_user: string; smtp_password: string }) {
		return apiFetch<{ ok: boolean }>('/mail/test-connection', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	},

	getAttachmentUrl(accountId: number, emailId: number, attachmentId: number): string {
		return `${backendBaseUrl}/accounts/${accountId}/mail/emails/${emailId}/attachments/${attachmentId}/download`;
	},

	getCIDImageUrl(accountId: number, emailId: number, cid: string): string {
		return `${backendBaseUrl}/accounts/${accountId}/mail/emails/${emailId}/cid/${encodeURIComponent(cid)}`;
	},

	searchContacts(accountId: number, query: string) {
		return apiFetch<{ contacts: Array<{ name: string; email: string; count: number }> }>(
			`/accounts/${accountId}/mail/contacts?q=${encodeURIComponent(query)}`
		);
	},

	async sendEmailWithAttachments(accountId: number, data: FormData): Promise<{ sent: boolean }> {
		const response = await fetch(`${backendBaseUrl}/accounts/${accountId}/mail/send`, {
			method: 'POST',
			credentials: 'include',
			body: data
		});
		if (!response.ok) {
			let payload: ApiErrorPayload | undefined;
			try {
				payload = (await response.json()) as ApiErrorPayload;
			} catch {
				payload = undefined;
			}
			throw new Error(payload?.error?.message || `Request failed with status ${response.status}`);
		}
		return (await response.json()) as { sent: boolean };
	},

	saveDraft(accountId: number, data: { to: string[]; cc?: string[]; subject: string; body: string; body_html?: string; in_reply_to?: string; references?: string[] }) {
		return apiFetch<{ id: number }>(`/accounts/${accountId}/mail/drafts`, {
			method: 'POST',
			body: JSON.stringify(data)
		});
	},

	deleteDraft(accountId: number, draftId: number) {
		return apiFetch<{ deleted: boolean }>(`/accounts/${accountId}/mail/drafts/${draftId}`, {
			method: 'DELETE'
		});
	},

	async downloadAttachment(accountId: number, emailId: number, attachmentId: number, filename: string): Promise<void> {
		const response = await fetch(
			`${backendBaseUrl}/accounts/${accountId}/mail/emails/${emailId}/attachments/${attachmentId}/download`,
			{ credentials: 'include' }
		);
		if (!response.ok) {
			throw new Error(`Download failed with status ${response.status}`);
		}
		const blob = await response.blob();
		const url = URL.createObjectURL(blob);
		const a = document.createElement('a');
		a.href = url;
		a.download = filename;
		document.body.appendChild(a);
		a.click();
		document.body.removeChild(a);
		URL.revokeObjectURL(url);
	},

	bulkAction(accountId: number, emailIds: number[], action: 'delete' | 'archive' | 'mark_read' | 'mark_unread') {
		return apiFetch<{ ok: boolean }>(`/accounts/${accountId}/mail/emails/bulk-action`, {
			method: 'POST',
			body: JSON.stringify({ email_ids: emailIds, action })
		});
	},

	listTemplates() {
		return apiFetch<{ templates: EmailTemplate[] }>('/templates');
	},

	createTemplate(data: { name: string; subject: string; body_html: string; body_text: string }) {
		return apiFetch<EmailTemplate>('/templates', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	},

	updateTemplate(templateId: number, data: { name: string; subject: string; body_html: string; body_text: string }) {
		return apiFetch<EmailTemplate>(`/templates/${templateId}`, {
			method: 'PUT',
			body: JSON.stringify(data)
		});
	},

	deleteTemplate(templateId: number) {
		return apiFetch<{ deleted: boolean }>(`/templates/${templateId}`, {
			method: 'DELETE'
		});
	}
};
