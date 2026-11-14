import { ApiError as KernelError, CSRF_HEADER, createApi, redirectToLogin, type CallOptions } from '@facile/client';
import { createNativeApi } from '@facile/client/native';
import { apiBase, clearToken, getToken, isApp } from '$lib/platform';
import { planLimitFrom, type PlanLimit } from '$lib/plan-limit';
import { setSessionLocale, t } from '$lib/i18n/index.svelte';

const backendBaseUrl = apiBase();

const kernel = isApp()
	? createNativeApi({ origin: backendBaseUrl, base: '/', token: getToken, onUnauthorized: clearToken })
	: createApi({ base: '/', onUnauthorized: redirectToLogin });

/** An inbox tab: where the classifier, Bulle's triage or a sender pin filed a conversation. */
/** A Bulle helper action, as POST /api/assist takes it. */
export type AssistAction =
	| 'rewrite'
	| 'shorten'
	| 'expand'
	| 'fix'
	| 'formal'
	| 'friendly'
	| 'translate'
	| 'reply'
	| 'summarize';

export type AssistRequest = {
	action: AssistAction;
	text?: string;
	context?: string;
	language?: string;
	instructions?: string;
};

export type AssistStatus = {
	available: boolean;
	reason?: 'off' | 'plan_required' | 'not_linked' | 'unreachable';
	actions: AssistAction[];
	bulle_url?: string;
	upgrade_url?: string;
};

export type InboxTab = 'primary' | 'newsletters' | 'notifications';

export type PushConfig = { enabled: boolean; native: boolean; web_push_key?: string };

export type PushDevice = { kind: 'fcm' | 'web'; token: string; p256dh?: string; auth?: string };

export type UserProfile = {
	id: string;
	email: string;
	name: string;
	avatar_url: string;
	avatar_source: string;
	locale?: string;
	created_at: string;
};

export type User = UserProfile;

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
	// Conversation aggregates, present only on folder-list (collapsed) rows.
	message_count?: number;
	unread_count?: number;
	email_ids?: number[];
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
	manage_url?: string;
	created_at: string;
	updated_at: string;
	members?: SpaceMember[];
};

export type SpaceMember = {
	id: string;
	user_id: string;
	email: string;
	name: string;
	avatar_url: string;
	role: string;
	joined_at: string;
};

/** One half of a connection check. `configured` false means no host is set for it at all. */
export type CheckLeg = {
	configured: boolean;
	ok: boolean;
	error?: string;
};

export type CheckResult = {
	imap: CheckLeg;
	smtp: CheckLeg;
};

export type AuthConfig = {
	sso_only: boolean;
	oidc_enabled: boolean;
};

export type InstanceSettings = {
	encryption_key_set: boolean;
};

export type ApiTokenStatus = {
	has_token: boolean;
	name?: string;
	created_at?: string;
};

export type ApiTokenCreated = {
	token: string;
	name: string;
	created_at: string;
};

/**
 * Carries the status alongside the message, because the two say different things and a
 * caller usually needs both: a 429 from the rate limiter and a 502 from a refused IMAP
 * handshake are the same `Error` otherwise, and a UI that cannot tell them apart ends up
 * blaming the mail server for the rate limiter. `retryAfter` is the header, in seconds.
 * `planLimit` is set on a 402: the scope's plan has no room, and where to upgrade it.
 */
export class ApiError extends Error {
	status: number;
	retryAfter?: number;
	planLimit?: PlanLimit;

	constructor(message: string, status: number, retryAfter?: number, planLimit?: PlanLimit) {
		super(message);
		this.name = 'ApiError';
		this.status = status;
		this.retryAfter = retryAfter;
		this.planLimit = planLimit;
	}
}

const SAFE_METHODS: ReadonlySet<string> = new Set(['GET', 'HEAD', 'OPTIONS']);

/**
 * Adds the credential to a request. The installed apps hold no cookie for the API's origin, so
 * it sends the bearer token from the mobile sign-in, which porte never asks a CSRF header of.
 * The web app sends its session cookie and the CSRF header on writes.
 */
export async function authed(options: RequestInit = {}): Promise<RequestInit> {
	const headers = new Headers(options.headers);
	if (isApp()) {
		const token = await getToken();
		if (token) headers.set('Authorization', `Bearer ${token}`);
		return { ...options, headers };
	}
	if (!SAFE_METHODS.has((options.method ?? 'GET').toUpperCase())) headers.set(CSRF_HEADER, '1');
	return { ...options, headers, credentials: 'include' };
}

export type IdeaSource = { label: string; ref: string };

export type BulleIdea = {
	id: string;
	kind: string;
	title: string;
	why: string;
	draft: string;
	sources: IdeaSource[] | null;
	status: string;
	expires_at: string;
};

export type IdeasResponse = {
	ai_access: boolean;
	upgrade_url?: string;
	ideas: BulleIdea[];
	generating: boolean;
	last_batch_at: string | null;
};

export type RunState = 'running' | 'needs_approval' | 'done' | 'failed';

export type IdeaRunResponse = {
	id: string;
	state: RunState;
	summary: string;
	outputs: { call_id: string; tool: string; url: string; preview: string }[] | null;
	pending_approvals: { call_id: string; tool: string; args: unknown }[] | null;
	continue_url: string;
};

const spaceQuery = (spaceId: string | null): string =>
	spaceId === null ? '' : `?space_id=${encodeURIComponent(spaceId)}`;

/** A draft as a pending action would send or delete it. */
export interface PendingDraft {
	id: number;
	to: string[];
	cc: string[];
	bcc: string[];
	subject: string;
	body: string;
}

/** An action an AI client asked for, waiting for the person to confirm it. */
export interface PendingAction {
	id: string;
	kind: 'send' | 'delete_draft';
	client_id: string;
	created_at: string;
	expires_at: string;
	account_id: number;
	account_email: string;
	state: 'ready' | 'changed' | 'missing';
	draft?: PendingDraft;
}

async function send(path: string, options: RequestInit = {}): Promise<Response> {
	const response = await fetch(`${backendBaseUrl}${path}`, await authed(options));
	if (response.status === 401 && isApp()) await clearToken();
	return response;
}

/** Turns the kernel's error into Courrier's, which call sites read for the 429 wait and the 402 plan. */
function toApiError(error: KernelError): ApiError {
	const message = typeof error.details.message === 'string' ? error.details.message : '';
	const seconds = error.retryAfterSeconds;
	return new ApiError(
		message || t('error.requestFailed', { status: error.status }),
		error.status,
		seconds !== null && seconds > 0 ? seconds : undefined,
		planLimitFrom(error.status, { error: error.details })
	);
}

/**
 * Calls the API through @facile/client's kernel. The web app sends its cookie and the CSRF header
 * on writes, and a 401 not marked quiet goes to the login page. The installed apps send the bearer
 * token to production and forget it on a 401. A JSON body gets its Content-Type either way.
 */
async function apiFetch<T>(path: string, options: RequestInit = {}, call?: CallOptions): Promise<T> {
	try {
		return await kernel<T>(path, options, call);
	} catch (error) {
		throw error instanceof KernelError ? toApiError(error) : error;
	}
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
	listIdeas(spaceId: string | null) {
		return apiFetch<IdeasResponse>(`/api/ideas${spaceQuery(spaceId)}`);
	},

	generateIdeas(spaceId: string | null) {
		return apiFetch<{ batch_id: string }>(`/api/ideas/generate${spaceQuery(spaceId)}`, { method: 'POST' });
	},

	runIdea(id: string, spaceId: string | null) {
		return apiFetch<{ run_id: string }>(`/api/ideas/${encodeURIComponent(id)}/run${spaceQuery(spaceId)}`, {
			method: 'POST'
		});
	},

	dismissIdea(id: string) {
		return apiFetch<object>(`/api/ideas/${encodeURIComponent(id)}/dismiss`, { method: 'POST' });
	},

	ideaRun(runId: string) {
		return apiFetch<IdeaRunResponse>(`/api/ideas/runs/${encodeURIComponent(runId)}`);
	},

	decideIdea(runId: string, callId: string, approve: boolean) {
		return apiFetch<object>(`/api/ideas/runs/${encodeURIComponent(runId)}/approvals/${encodeURIComponent(callId)}`, {
			method: 'POST',
			body: JSON.stringify({ approve })
		});
	},

	baseUrl: backendBaseUrl,

	async logout() {
		try {
			return await apiFetch<{ ok: boolean }>('/api/auth/logout', { method: 'POST' });
		} finally {
			await clearToken();
		}
	},
	/** The pending Registre space invitations for the user menu, or null when there is none to show. */
	invitations() {
		return apiFetch<{ count: number; href?: string }>('/api/me/invitations')
			.then((r) => (r.href ? { count: r.count, href: r.href } : null))
			.catch(() => null);
	},
	me() {
		return apiFetch<MeResponse>('/api/users/me', {}, { quiet: true }).then((r) => {
			setSessionLocale(r.user.locale ?? null);
			return { user: normalizeUser(r.user) };
		});
	},

	listAccounts(spaceId?: string) {
		const params = spaceId ? `?space_id=${encodeURIComponent(spaceId)}` : '';
		return apiFetch<{ accounts: MailAccount[] }>(`/api/accounts${params}`);
	},
	accountUsage(spaceId?: string) {
		const params = spaceId ? `?space_id=${encodeURIComponent(spaceId)}` : '';
		return apiFetch<{ used: number; limit: number }>(`/api/accounts/usage${params}`);
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

	checkAccount(accountId: number) {
		return apiFetch<CheckResult>(`/api/accounts/${accountId}/mail/check`, { method: 'POST' });
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
	getEmailsByFolder(accountId: number, folderType: string, page = 1, limit = 50, unreadOnly = false, tab?: InboxTab) {
		const params = new URLSearchParams({ page: String(page), limit: String(limit) });
		if (unreadOnly) params.set('unread', 'true');
		if (tab) params.set('tab', tab);
		return apiFetch<{ emails: EmailMessage[]; total: number; page: number; limit: number }>(
			`/api/accounts/${accountId}/mail/folders/${folderType}/emails?${params}`
		);
	},
	assistStatus() {
		return apiFetch<AssistStatus>('/api/assist');
	},
	assist(request: AssistRequest) {
		return apiFetch<{ text: string }>('/api/assist', { method: 'POST', body: JSON.stringify(request) });
	},
	setSenderRule(accountId: number, address: string, tab: InboxTab | 'auto') {
		return apiFetch<void>(`/api/accounts/${accountId}/mail/sender-rules`, {
			method: 'PUT',
			body: JSON.stringify({ address, tab })
		});
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

	getCIDImageUrl(accountId: number, emailId: number, cid: string, token = ''): string {
		const query = token ? `?token=${encodeURIComponent(token)}` : '';
		return `${backendBaseUrl}/api/accounts/${accountId}/mail/emails/${emailId}/cid/${encodeURIComponent(cid)}${query}`;
	},

	resourceToken() {
		return apiFetch<{ token: string }>('/api/auth/resource-token');
	},

	searchContacts(accountId: number, query: string) {
		return apiFetch<{ contacts: Array<{ name: string; email: string; count: number }> }>(
			`/api/accounts/${accountId}/mail/contacts?q=${encodeURIComponent(query)}`
		);
	},

	sendEmailWithAttachments(accountId: number, data: FormData): Promise<{ sent: boolean }> {
		return apiFetch<{ sent: boolean }>(`/api/accounts/${accountId}/mail/send`, { method: 'POST', body: data });
	},

	saveDraft(accountId: number, data: { to: string[]; cc?: string[]; subject: string; body: string; body_html?: string; in_reply_to?: string; references?: string[] }) {
		return apiFetch<{ id: number }>(`/api/accounts/${accountId}/mail/drafts`, {
			method: 'POST',
			body: JSON.stringify(data)
		});
	},

	listPending() {
		return apiFetch<{ actions: PendingAction[] }>('/api/mcp/pending');
	},

	confirmPending(id: string) {
		return apiFetch<{ status: string }>(`/api/mcp/pending/${encodeURIComponent(id)}/confirm`, { method: 'POST' });
	},

	dismissPending(id: string) {
		return apiFetch<{ status: string }>(`/api/mcp/pending/${encodeURIComponent(id)}/dismiss`, { method: 'POST' });
	},

	deleteDraft(accountId: number, draftId: number) {
		return apiFetch<{ deleted: boolean }>(`/api/accounts/${accountId}/mail/drafts/${draftId}`, {
			method: 'DELETE'
		});
	},

	async downloadAttachment(accountId: number, emailId: number, attachmentId: number, filename: string): Promise<void> {
		const response = await send(
			`/api/accounts/${accountId}/mail/emails/${emailId}/attachments/${attachmentId}/download`
		);
		if (!response.ok) {
			throw new Error(t('error.downloadFailed', { status: response.status }));
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
		return apiFetch<{ spaces: Space[]; manage_url?: string }>('/api/spaces');
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
		return apiFetch<{ members: SpaceMember[] }>(`/api/spaces/${spaceId}/members`).then((r) => ({
			members: (r.members ?? []).map((m) => ({ ...m, avatar_url: resolveFileUrl(m.avatar_url) }))
		}));
	},

	addSpaceMember(spaceId: string, data: { user_id?: number; email?: string; role?: string }) {
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
	},

	listUsers(): Promise<{ users: User[] }> {
		return apiFetch<{ users: User[] }>('/api/users');
	},

	authConfig() {
		return apiFetch<AuthConfig>('/api/auth/config');
	},
	instanceSettings() {
		return apiFetch<{ settings: InstanceSettings }>('/api/settings/').then((r) => r.settings);
	},

	apiTokenStatus() {
		return apiFetch<ApiTokenStatus>('/api/users/me/api-token');
	},
	createApiToken(name: string) {
		return apiFetch<ApiTokenCreated>('/api/users/me/api-token', {
			method: 'POST',
			body: JSON.stringify({ name })
		});
	},
	deleteApiToken() {
		return apiFetch<{ deleted: boolean }>('/api/users/me/api-token', { method: 'DELETE' });
	},

	getPush() {
		return apiFetch<PushConfig>('/api/push');
	},
	testPush() {
		return apiFetch<void>('/api/push/test', { method: 'POST' });
	},
	registerDevice(device: PushDevice) {
		return apiFetch<void>('/api/push/devices', { method: 'POST', body: JSON.stringify(device) });
	},
	forgetDevice(token: string) {
		return apiFetch<void>('/api/push/devices', { method: 'DELETE', body: JSON.stringify({ token }) });
	}
};
