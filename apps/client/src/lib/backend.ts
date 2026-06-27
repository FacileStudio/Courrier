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
	thread_id?: string;
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

export type Space = {
	id: string;
	name: string;
	description: string;
	role: string;
	created_at: string;
	updated_at: string;
	members?: SpaceMember[];
};

export type SpaceMember = {
	id: string;
	user_id: string;
	email: string;
	name: string;
	role: string;
	joined_at: string;
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
		return apiFetch<AuthResponse>('/api/auth/register', {
			method: 'POST',
			body: JSON.stringify({ email, password })
		});
	},
	login(email: string, password: string) {
		return apiFetch<AuthResponse>('/api/auth/login', {
			method: 'POST',
			body: JSON.stringify({ email, password })
		});
	},
	logout() {
		return apiFetch<{ ok: boolean }>('/api/auth/logout', { method: 'POST' });
	},
	me() {
		return apiFetch<MeResponse>('/api/users/me').then((r) => ({
			user: normalizeUser(r.user)
		}));
	},

	listAccounts(spaceId?: string) {
		const params = spaceId ? `?space_id=${encodeURIComponent(spaceId)}` : '';
		return apiFetch<{ accounts: MailAccount[] }>(`/api/accounts${params}`);
	},
	getAccount(id: number) {
		return apiFetch<MailAccount>(`/api/accounts/${id}`);
	},
	createAccount(data: Omit<MailAccount, 'id' | 'created_at' | 'updated_at'> & { imap_password: string; smtp_password: string; space_id?: string }) {
		return apiFetch<MailAccount>('/api/accounts', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	},
	updateAccount(id: number, data: Partial<MailAccount & { imap_password: string; smtp_password: string }>) {
		return apiFetch<MailAccount>(`/api/accounts/${id}`, {
			method: 'PUT',
			body: JSON.stringify(data)
		});
	},
	deleteAccount(id: number) {
		return apiFetch<{ deleted: boolean }>(`/api/accounts/${id}`, { method: 'DELETE' });
	},

	syncProfile() {
		return apiFetch<{ synced: boolean }>('/api/auth/sync-profile', { method: 'POST' });
	},

	syncAccount(accountId: number) {
		return apiFetch<{ synced: boolean }>(`/api/accounts/${accountId}/mail/sync`, { method: 'POST' });
	},
	syncFolder(accountId: number, folderId: number) {
		return apiFetch<{ synced: boolean }>(`/api/accounts/${accountId}/mail/folders/${folderId}/sync`, { method: 'POST' });
	},
	getFolders(accountId: number) {
		return apiFetch<{ folders: Folder[] }>(`/api/accounts/${accountId}/mail/folders`);
	},
	getEmailsByFolder(accountId: number, folderType: string, page = 1, limit = 50, unreadOnly = false) {
		const params = new URLSearchParams({ page: String(page), limit: String(limit) });
		if (unreadOnly) params.set('unread', 'true');
		return apiFetch<{ emails: EmailMessage[]; total: number; page: number; limit: number }>(
			`/api/accounts/${accountId}/mail/folders/${folderType}/emails?${params}`
		);
	},
	searchEmails(accountId: number, query: string, page = 1, limit = 30) {
		const params = new URLSearchParams({ q: query, page: String(page), limit: String(limit) });
		return apiFetch<{ emails: EmailMessage[]; total: number; page: number; limit: number }>(
			`/api/accounts/${accountId}/mail/search?${params}`
		);
	},
	getEmail(accountId: number, emailId: number) {
		return apiFetch<EmailMessage>(`/api/accounts/${accountId}/mail/emails/${emailId}`);
	},
	getThread(accountId: number, threadId: string) {
		return apiFetch<{ emails: EmailMessage[] }>(
			`/api/accounts/${accountId}/mail/threads/${encodeURIComponent(threadId)}`
		);
	},
	updateEmail(accountId: number, emailId: number, data: { is_read?: boolean; is_starred?: boolean }) {
		return apiFetch<EmailMessage>(`/api/accounts/${accountId}/mail/emails/${emailId}`, {
			method: 'PATCH',
			body: JSON.stringify(data)
		});
	},
	sendEmail(accountId: number, data: { to: string[]; cc?: string[]; subject: string; body: string; body_html?: string; in_reply_to?: string; references?: string[] }) {
		return apiFetch<{ sent: boolean }>(`/api/accounts/${accountId}/mail/send`, {
			method: 'POST',
			body: JSON.stringify(data)
		});
	},
	testConnection(data: { imap_host: string; imap_port: number; imap_user: string; imap_password: string; smtp_host: string; smtp_port: number; smtp_user: string; smtp_password: string }) {
		return apiFetch<{ ok: boolean }>('/api/mail/test-connection', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	},

	getAttachmentUrl(accountId: number, emailId: number, attachmentId: number): string {
		return `${backendBaseUrl}/api/accounts/${accountId}/mail/emails/${emailId}/attachments/${attachmentId}/download`;
	},

	getCIDImageUrl(accountId: number, emailId: number, cid: string): string {
		return `${backendBaseUrl}/api/accounts/${accountId}/mail/emails/${emailId}/cid/${encodeURIComponent(cid)}`;
	},

	searchContacts(accountId: number, query: string) {
		return apiFetch<{ contacts: Array<{ name: string; email: string; count: number }> }>(
			`/api/accounts/${accountId}/mail/contacts?q=${encodeURIComponent(query)}`
		);
	},

	async sendEmailWithAttachments(accountId: number, data: FormData): Promise<{ sent: boolean }> {
		const response = await fetch(`${backendBaseUrl}/api/accounts/${accountId}/mail/send`, {
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
		return apiFetch<{ id: number }>(`/api/accounts/${accountId}/mail/drafts`, {
			method: 'POST',
			body: JSON.stringify(data)
		});
	},

	deleteDraft(accountId: number, draftId: number) {
		return apiFetch<{ deleted: boolean }>(`/api/accounts/${accountId}/mail/drafts/${draftId}`, {
			method: 'DELETE'
		});
	},

	async downloadAttachment(accountId: number, emailId: number, attachmentId: number, filename: string): Promise<void> {
		const response = await fetch(
			`${backendBaseUrl}/api/accounts/${accountId}/mail/emails/${emailId}/attachments/${attachmentId}/download`,
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
		return apiFetch<{ ok: boolean }>(`/api/accounts/${accountId}/mail/emails/bulk-action`, {
			method: 'POST',
			body: JSON.stringify({ email_ids: emailIds, action })
		});
	},

	listTemplates(spaceId?: string) {
		const params = spaceId ? `?space_id=${encodeURIComponent(spaceId)}` : '';
		return apiFetch<{ templates: EmailTemplate[] }>(`/api/templates${params}`);
	},

	createTemplate(data: { name: string; subject: string; body_html: string; body_text: string; space_id?: string }) {
		return apiFetch<EmailTemplate>('/api/templates', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	},

	updateTemplate(templateId: number, data: { name: string; subject: string; body_html: string; body_text: string }) {
		return apiFetch<EmailTemplate>(`/api/templates/${templateId}`, {
			method: 'PUT',
			body: JSON.stringify(data)
		});
	},

	deleteTemplate(templateId: number) {
		return apiFetch<{ deleted: boolean }>(`/api/templates/${templateId}`, {
			method: 'DELETE'
		});
	},

	listSpaces() {
		return apiFetch<{ spaces: Space[] }>('/api/spaces');
	},

	getSpace(spaceId: string) {
		return apiFetch<Space>(`/api/spaces/${spaceId}`);
	},

	createSpace(data: { name: string; description?: string }) {
		return apiFetch<Space>('/api/spaces', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	},

	updateSpace(spaceId: string, data: { name?: string; description?: string }) {
		return apiFetch<Space>(`/api/spaces/${spaceId}`, {
			method: 'PUT',
			body: JSON.stringify(data)
		});
	},

	deleteSpace(spaceId: string) {
		return apiFetch<{ deleted: boolean }>(`/api/spaces/${spaceId}`, {
			method: 'DELETE'
		});
	},

	leaveSpace(spaceId: string) {
		return apiFetch<{ left: boolean }>(`/api/spaces/${spaceId}/leave`, {
			method: 'POST'
		});
	},

	listSpaceMembers(spaceId: string) {
		return apiFetch<{ members: SpaceMember[] }>(`/api/spaces/${spaceId}/members`);
	},

	addSpaceMember(spaceId: string, data: { user_id: number; role?: string }) {
		return apiFetch<SpaceMember>(`/api/spaces/${spaceId}/members`, {
			method: 'POST',
			body: JSON.stringify(data)
		});
	},

	updateSpaceMember(spaceId: string, memberId: string, data: { role: string }) {
		return apiFetch<{ id: string; role: string }>(`/api/spaces/${spaceId}/members/${memberId}`, {
			method: 'PUT',
			body: JSON.stringify(data)
		});
	},

	removeSpaceMember(spaceId: string, memberId: string) {
		return apiFetch<{ removed: boolean }>(`/api/spaces/${spaceId}/members/${memberId}`, {
			method: 'DELETE'
		});
	}
};
