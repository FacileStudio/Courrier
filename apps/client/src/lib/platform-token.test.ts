import { beforeEach, describe, expect, it, vi } from 'vitest';

const KEY = 'courrier.token';

const env = vi.hoisted(() => ({
	native: false,
	desktop: false,
	refuseWrites: false,
	preferences: new Map<string, string>(),
	secure: new Map<string, string>(),
	keychain: null as string | null
}));

vi.mock('@capacitor/core', () => ({ Capacitor: { isNativePlatform: () => env.native } }));

vi.mock('@capacitor/preferences', () => ({
	Preferences: {
		get: async ({ key }: { key: string }) => ({ value: env.preferences.get(key) ?? null }),
		set: async ({ key, value }: { key: string; value: string }) => void env.preferences.set(key, value),
		remove: async ({ key }: { key: string }) => void env.preferences.delete(key)
	}
}));

vi.mock('capacitor-secure-storage-plugin', () => ({
	SecureStoragePlugin: {
		get: async ({ key }: { key: string }) => {
			const value = env.secure.get(key);
			if (value === undefined) throw new Error('Item with given key does not exist');
			return { value };
		},
		set: async ({ key, value }: { key: string; value: string }) => {
			if (env.refuseWrites) throw new Error('keystore unavailable');
			env.secure.set(key, value);
			return { value: true };
		},
		remove: async ({ key }: { key: string }) => {
			if (!env.secure.delete(key)) throw new Error('Item with given key does not exist');
			return { value: true };
		}
	}
}));

vi.mock('@tauri-apps/api/core', () => ({
	isTauri: () => env.desktop,
	invoke: async (command: string, args?: { token: string }) => {
		if (command === 'token_get') return env.keychain;
		if (command === 'token_set') {
			if (env.refuseWrites) throw new Error('no secret service');
			env.keychain = args?.token ?? null;
			return null;
		}
		if (command === 'token_clear') {
			env.keychain = null;
			return null;
		}
		throw new Error(`unexpected command ${command}`);
	}
}));

async function platform() {
	vi.resetModules();
	return import('./platform');
}

beforeEach(() => {
	env.native = false;
	env.desktop = false;
	env.refuseWrites = false;
	env.preferences.clear();
	env.secure.clear();
	env.keychain = null;
	vi.spyOn(console, 'warn').mockImplementation(() => undefined);
});

describe('getToken on Android', () => {
	beforeEach(() => {
		env.native = true;
	});

	it('moves a Preferences token into the Keystore and deletes the old copy', async () => {
		env.preferences.set(KEY, 'old-token');
		const { getToken } = await platform();
		expect(await getToken()).toBe('old-token');
		expect(env.secure.get(KEY)).toBe('old-token');
		expect(env.preferences.has(KEY)).toBe(false);
	});

	it('keeps the Preferences token when the Keystore refuses it, so nobody is signed out', async () => {
		env.preferences.set(KEY, 'old-token');
		env.refuseWrites = true;
		const { getToken } = await platform();
		expect(await getToken()).toBe('old-token');
		expect(env.preferences.get(KEY)).toBe('old-token');
	});

	it('prefers the Keystore token and drops a stale Preferences copy', async () => {
		env.secure.set(KEY, 'new-token');
		env.preferences.set(KEY, 'old-token');
		const { getToken } = await platform();
		expect(await getToken()).toBe('new-token');
		expect(env.preferences.has(KEY)).toBe(false);
	});

	it('answers null when no token exists anywhere', async () => {
		const { getToken } = await platform();
		expect(await getToken()).toBeNull();
	});

	it('stores a new token in the Keystore only, and clears every copy on sign-out', async () => {
		const { setToken, clearToken, getToken } = await platform();
		await setToken('fresh');
		expect(env.secure.get(KEY)).toBe('fresh');
		expect(env.preferences.has(KEY)).toBe(false);
		env.preferences.set(KEY, 'leftover');
		await clearToken();
		expect(env.secure.has(KEY)).toBe(false);
		expect(env.preferences.has(KEY)).toBe(false);
		expect(await getToken()).toBeNull();
	});
});

describe('getToken on the desktop', () => {
	beforeEach(() => {
		env.desktop = true;
	});

	it('moves the localStorage Preferences token into the OS keychain', async () => {
		env.preferences.set(KEY, 'old-token');
		const { getToken } = await platform();
		expect(await getToken()).toBe('old-token');
		expect(env.keychain).toBe('old-token');
		expect(env.preferences.has(KEY)).toBe(false);
	});

	it('keeps the old copy when there is no keychain to write to', async () => {
		env.preferences.set(KEY, 'old-token');
		env.refuseWrites = true;
		const { getToken } = await platform();
		expect(await getToken()).toBe('old-token');
		expect(env.preferences.get(KEY)).toBe('old-token');
	});
});

describe('getToken on the web', () => {
	it('never reads or moves anything', async () => {
		env.preferences.set(KEY, 'app-token');
		const { getToken, setToken } = await platform();
		expect(await getToken()).toBeNull();
		await setToken('ignored');
		expect(env.secure.size).toBe(0);
		expect(env.preferences.get(KEY)).toBe('app-token');
	});
});
