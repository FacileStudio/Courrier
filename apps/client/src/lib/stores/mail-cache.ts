import { writable, get } from 'svelte/store';
import type { EmailMessage } from '$lib/backend';

const CACHE_KEY = 'courrier.mail-cache';
const CACHE_TTL = 10 * 60 * 1000;

type CacheEntry = {
	emails: EmailMessage[];
	total: number;
	timestamp: number;
};

type CacheStore = Record<string, CacheEntry>;

function loadCache(): CacheStore {
	try {
		const raw = localStorage.getItem(CACHE_KEY);
		if (!raw) return {};
		const data = JSON.parse(raw) as CacheStore;
		const now = Date.now();
		for (const key of Object.keys(data)) {
			if (now - data[key].timestamp > CACHE_TTL) {
				delete data[key];
			}
		}
		return data;
	} catch {
		return {};
	}
}

function saveCache(cache: CacheStore) {
	try {
		localStorage.setItem(CACHE_KEY, JSON.stringify(cache));
	} catch {
		// quota exceeded — clear old entries
		localStorage.removeItem(CACHE_KEY);
	}
}

function cacheKey(accountId: number, folderType: string, page: number): string {
	return `${accountId}:${folderType}:${page}`;
}

const store = writable<CacheStore>(loadCache());

store.subscribe((value) => {
	saveCache(value);
});

export const mailCache = {
	get(accountId: number, folderType: string, page: number): CacheEntry | null {
		const cache = get(store);
		const key = cacheKey(accountId, folderType, page);
		const entry = cache[key];
		if (!entry) return null;
		if (Date.now() - entry.timestamp > CACHE_TTL) {
			store.update((c) => {
				delete c[key];
				return { ...c };
			});
			return null;
		}
		return entry;
	},

	set(accountId: number, folderType: string, page: number, emails: EmailMessage[], total: number) {
		const stripped = emails.map((e) => ({
			...e,
			body_text: '',
			body_html: ''
		}));
		store.update((c) => ({
			...c,
			[cacheKey(accountId, folderType, page)]: {
				emails: stripped,
				total,
				timestamp: Date.now()
			}
		}));
	},

	invalidateFolder(accountId: number, folderType: string) {
		store.update((c) => {
			const prefix = `${accountId}:${folderType}:`;
			const next = { ...c };
			for (const key of Object.keys(next)) {
				if (key.startsWith(prefix)) delete next[key];
			}
			return next;
		});
	},

	invalidateAll() {
		store.set({});
	},

	removeEmails(emailIds: number[]) {
		const idSet = new Set(emailIds);
		store.update((c) => {
			const next = { ...c };
			for (const key of Object.keys(next)) {
				const entry = next[key];
				const filtered = entry.emails.filter((e) => !idSet.has(e.id));
				if (filtered.length !== entry.emails.length) {
					next[key] = {
						...entry,
						emails: filtered,
						total: entry.total - (entry.emails.length - filtered.length)
					};
				}
			}
			return next;
		});
	}
};
