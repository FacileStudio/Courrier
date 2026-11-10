import { beforeEach, describe, expect, it } from 'vitest';
import { rememberReturn, takeReturn } from './after-login';

beforeEach(() => {
	const store = new Map<string, string>();
	globalThis.sessionStorage = {
		getItem: (key: string) => store.get(key) ?? null,
		setItem: (key: string, value: string) => void store.set(key, value),
		removeItem: (key: string) => void store.delete(key)
	} as Storage;
});

describe('after-login return', () => {
	it('brings a confirm link back once', () => {
		rememberReturn(new URL('https://courrier.example/pending?action=abc'));
		expect(takeReturn()).toBe('/pending?action=abc');
		expect(takeReturn()).toBeNull();
	});

	it('ignores other pages and foreign paths', () => {
		rememberReturn(new URL('https://courrier.example/mail'));
		expect(takeReturn()).toBeNull();
		sessionStorage.setItem('courrier.after-login', '//evil.example/pending');
		expect(takeReturn()).toBeNull();
	});
});
