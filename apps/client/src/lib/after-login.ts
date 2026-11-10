const KEY = 'courrier.after-login';

/** Keeps the page a signed-out person opened, such as a confirm link, for after sign-in. */
export function rememberReturn(url: URL) {
	if (!url.pathname.startsWith('/pending')) return;
	try {
		sessionStorage.setItem(KEY, url.pathname + url.search);
	} catch {
		return;
	}
}

/** Returns the remembered page once, as a same-origin path, or null. */
export function takeReturn(): string | null {
	try {
		const path = sessionStorage.getItem(KEY);
		sessionStorage.removeItem(KEY);
		return path && path.startsWith('/') && !path.startsWith('//') ? path : null;
	} catch {
		return null;
	}
}
