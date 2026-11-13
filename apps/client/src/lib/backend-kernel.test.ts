import { afterEach, beforeAll, beforeEach, describe, expect, it, vi } from 'vitest';

const platform = vi.hoisted(() => ({ app: false, token: 'native-token' as string | null, cleared: 0 }));

vi.mock('$lib/platform', () => ({
	isApp: () => platform.app,
	apiBase: () => (platform.app ? 'https://courrier.facile.studio' : ''),
	getToken: async () => platform.token,
	clearToken: async () => {
		platform.cleared += 1;
	}
}));

type Call = { url: string; init: RequestInit };

function stubFetch(status: number, body: string, headers: Record<string, string> = {}): Call[] {
	const calls: Call[] = [];
	vi.stubGlobal('fetch', async (url: string, init: RequestInit) => {
		calls.push({ url, init });
		return new Response(status === 204 ? null : body, { status, headers });
	});
	return calls;
}

type Backend = typeof import('./backend');

const LOAD_TIMEOUT = 60_000;

async function load(app: boolean): Promise<Backend> {
	platform.app = app;
	vi.resetModules();
	return import('./backend');
}

beforeEach(() => {
	platform.cleared = 0;
});

afterEach(() => {
	vi.unstubAllGlobals();
});

describe('web kernel', () => {
	let mod: Backend;
	beforeAll(async () => {
		mod = await load(false);
	}, LOAD_TIMEOUT);

	it('keeps the same-origin path, sends the cookie and the CSRF header on a write', async () => {
		const calls = stubFetch(200, '{"id":7}');
		const { backend } = mod;
		await backend.createTemplate({ name: 'n', subject: 's', body_html: '', body_text: '' });
		const { url, init } = calls[0];
		const headers = new Headers(init.headers);
		expect(url).toBe('/api/templates');
		expect(init.method).toBe('POST');
		expect(init.credentials).toBe('same-origin');
		expect(headers.get('X-Facile-CSRF')).toBe('1');
		expect(headers.get('Content-Type')).toBe('application/json');
		expect(headers.has('Authorization')).toBe(false);
	});

	it('sends no CSRF header on a read', async () => {
		const calls = stubFetch(200, '{"spaces":[]}');
		const { backend } = mod;
		await backend.listSpaces();
		expect(calls[0].url).toBe('/api/spaces');
		expect(new Headers(calls[0].init.headers).has('X-Facile-CSRF')).toBe(false);
	});

	it('reads the envelope into Courrier ApiError with the wait and the plan', async () => {
		stubFetch(429, '{"error":{"code":"rate_limited","message":"slow down"}}', { 'Retry-After': '30' });
		const { backend, ApiError } = mod;
		const error = await backend.listSpaces().catch((e: unknown) => e);
		expect(error).toBeInstanceOf(ApiError);
		expect(error).toMatchObject({ status: 429, message: 'slow down', retryAfter: 30 });

		stubFetch(402, '{"error":{"code":"plan_limit","message":"full","key":"courrier.mailboxes","limit":1}}');
		const refused = await backend.listSpaces().catch((e: unknown) => e);
		expect(refused).toMatchObject({ status: 402, planLimit: { message: 'full', key: 'courrier.mailboxes', limit: 1 } });
	});

	it('answers undefined for a 204', async () => {
		stubFetch(204, '');
		const { backend } = mod;
		expect(await backend.registerDevice({ kind: 'web', token: 't' })).toBeUndefined();
	});
});

describe('native kernel', () => {
	let mod: Backend;
	beforeAll(async () => {
		mod = await load(true);
	}, LOAD_TIMEOUT);

	it('sends the bearer token to production, with no cookie and no CSRF header', async () => {
		const calls = stubFetch(200, '{"id":7}');
		const { backend } = mod;
		await backend.createTemplate({ name: 'n', subject: 's', body_html: '', body_text: '' });
		const { url, init } = calls[0];
		const headers = new Headers(init.headers);
		expect(url).toBe('https://courrier.facile.studio/api/templates');
		expect(init.credentials).toBe('omit');
		expect(headers.get('Authorization')).toBe('Bearer native-token');
		expect(headers.has('X-Facile-CSRF')).toBe(false);
		expect(headers.get('Content-Type')).toBe('application/json');
	});

	it('forgets the token on a 401', async () => {
		stubFetch(401, '{"error":{"code":"unauthorized","message":"signed out"}}');
		const { backend, ApiError } = mod;
		const error = await backend.listSpaces().catch((e: unknown) => e);
		expect(error).toBeInstanceOf(ApiError);
		expect(platform.cleared).toBe(1);
	});
});
