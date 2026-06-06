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

type ApiErrorPayload = {
	error?: { message?: string };
};

async function apiFetch<T>(path: string, options: RequestInit = {}, token?: string) {
	const headers = new Headers(options.headers);
	if (!headers.has('Content-Type') && options.body) {
		headers.set('Content-Type', 'application/json');
	}
	if (token) {
		headers.set('Authorization', `Bearer ${token}`);
	}
	const response = await fetch(`${backendBaseUrl}${path}`, { ...options, headers });
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
	me(token: string) {
		return apiFetch<MeResponse>('/users/me', {}, token).then((r) => ({
			user: normalizeUser(r.user)
		}));
	},

	listAccounts(token: string) {
		return apiFetch<{ accounts: MailAccount[] }>('/accounts', {}, token);
	},
	getAccount(token: string, id: number) {
		return apiFetch<MailAccount>(`/accounts/${id}`, {}, token);
	},
	createAccount(token: string, data: Omit<MailAccount, 'id' | 'created_at' | 'updated_at'> & { imap_password: string; smtp_password: string }) {
		return apiFetch<MailAccount>('/accounts', {
			method: 'POST',
			body: JSON.stringify(data)
		}, token);
	},
	updateAccount(token: string, id: number, data: Partial<MailAccount & { imap_password: string; smtp_password: string }>) {
		return apiFetch<MailAccount>(`/accounts/${id}`, {
			method: 'PUT',
			body: JSON.stringify(data)
		}, token);
	},
	deleteAccount(token: string, id: number) {
		return apiFetch<{ deleted: boolean }>(`/accounts/${id}`, { method: 'DELETE' }, token);
	},

	syncProfile(token: string) {
		return apiFetch<{ synced: boolean }>('/auth/sync-profile', { method: 'POST' }, token);
	},

	syncAccount(token: string, accountId: number) {
		return apiFetch<{ synced: boolean }>(`/accounts/${accountId}/mail/sync`, { method: 'POST' }, token);
	},
	syncFolder(token: string, accountId: number, folderId: number) {
		return apiFetch<{ synced: boolean }>(`/accounts/${accountId}/mail/folders/${folderId}/sync`, { method: 'POST' }, token);
	},
	getFolders(token: string, accountId: number) {
		return apiFetch<{ folders: Folder[] }>(`/accounts/${accountId}/mail/folders`, {}, token);
	},
	getEmailsByFolder(token: string, accountId: number, folderType: string, page = 1, limit = 50) {
		return apiFetch<{ emails: EmailMessage[]; total: number; page: number; limit: number }>(
			`/accounts/${accountId}/mail/folders/${folderType}/emails?page=${page}&limit=${limit}`,
			{},
			token
		);
	},
	getEmail(token: string, accountId: number, emailId: number) {
		return apiFetch<EmailMessage>(`/accounts/${accountId}/mail/emails/${emailId}`, {}, token);
	},
	updateEmail(token: string, accountId: number, emailId: number, data: { is_read?: boolean; is_starred?: boolean }) {
		return apiFetch<EmailMessage>(`/accounts/${accountId}/mail/emails/${emailId}`, {
			method: 'PATCH',
			body: JSON.stringify(data)
		}, token);
	},
	sendEmail(token: string, accountId: number, data: { to: string[]; cc?: string[]; subject: string; body: string; body_html?: string; in_reply_to?: string; references?: string[] }) {
		return apiFetch<{ sent: boolean }>(`/accounts/${accountId}/mail/send`, {
			method: 'POST',
			body: JSON.stringify(data)
		}, token);
	},
	testConnection(data: { imap_host: string; imap_port: number; imap_user: string; imap_password: string; smtp_host: string; smtp_port: number; smtp_user: string; smtp_password: string }) {
		return apiFetch<{ ok: boolean }>('/mail/test-connection', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	},

	getAttachmentUrl(token: string, accountId: number, emailId: number, attachmentId: number): string {
		return `${backendBaseUrl}/accounts/${accountId}/mail/emails/${emailId}/attachments/${attachmentId}/download?token=${encodeURIComponent(token)}`;
	},

	getCIDImageUrl(token: string, accountId: number, emailId: number, cid: string): string {
		return `${backendBaseUrl}/accounts/${accountId}/mail/emails/${emailId}/cid/${encodeURIComponent(cid)}?token=${encodeURIComponent(token)}`;
	},

	searchContacts(token: string, accountId: number, query: string) {
		return apiFetch<{ contacts: Array<{ name: string; email: string; count: number }> }>(
			`/accounts/${accountId}/mail/contacts?q=${encodeURIComponent(query)}`,
			{},
			token
		);
	},

	async sendEmailWithAttachments(token: string, accountId: number, data: FormData): Promise<{ sent: boolean }> {
		const response = await fetch(`${backendBaseUrl}/accounts/${accountId}/mail/send`, {
			method: 'POST',
			headers: { Authorization: `Bearer ${token}` },
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

	saveDraft(token: string, accountId: number, data: { to: string[]; cc?: string[]; subject: string; body: string; body_html?: string; in_reply_to?: string; references?: string[] }) {
		return apiFetch<{ id: number }>(`/accounts/${accountId}/mail/drafts`, {
			method: 'POST',
			body: JSON.stringify(data)
		}, token);
	},

	deleteDraft(token: string, accountId: number, draftId: number) {
		return apiFetch<{ deleted: boolean }>(`/accounts/${accountId}/mail/drafts/${draftId}`, {
			method: 'DELETE'
		}, token);
	},

	async downloadAttachment(token: string, accountId: number, emailId: number, attachmentId: number, filename: string): Promise<void> {
		const response = await fetch(
			`${backendBaseUrl}/accounts/${accountId}/mail/emails/${emailId}/attachments/${attachmentId}/download`,
			{ headers: { Authorization: `Bearer ${token}` } }
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
	}
};
