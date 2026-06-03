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
	}
};
