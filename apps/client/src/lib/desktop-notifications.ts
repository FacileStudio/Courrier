import { goto } from '$app/navigation';
import { toast } from '@facile/muse';
import {
	isPermissionGranted,
	requestPermission,
	sendNotification
} from '@tauri-apps/plugin-notification';
import { authed, backend } from '$lib/backend';
import { routeFor } from '$lib/push-route';
import { createSseParser } from '$lib/sse';
import { t } from '$lib/i18n/index.svelte';

type StreamNotification = { title: string; body: string; path: string };

type StreamEnd = 'closed' | 'failed' | 'signed-out';

const MAX_BACKOFF_MS = 30000;

let listening = false;

/** Reads whether the desktop OS lets Courrier show notifications. */
export async function desktopPermission(): Promise<'granted' | 'prompt'> {
	return (await isPermissionGranted()) ? 'granted' : 'prompt';
}

/** Asks the desktop OS for permission, then starts listening. */
export async function enableDesktop(): Promise<'granted' | 'denied'> {
	if (!(await isPermissionGranted()) && (await requestPermission()) !== 'granted') return 'denied';
	listenDesktop();
	return 'granted';
}

function show(note: StreamNotification) {
	const route = routeFor(note.path);
	if (document.hasFocus()) {
		toast.info(note.title, { action: { label: t('common.open'), onClick: () => void goto(route) } });
		return;
	}
	sendNotification({ title: note.title || 'Courrier', body: note.body });
}

async function readStream(): Promise<StreamEnd> {
	const init = await authed({ headers: { Accept: 'text/event-stream' } });
	const response = await fetch(`${backend.baseUrl}/api/push/stream`, init);
	if (response.status === 401) return 'signed-out';
	if (!response.ok || !response.body) return 'failed';
	const parse = createSseParser((frame) => {
		if (frame.event !== 'notification') return;
		try {
			show(JSON.parse(frame.data) as StreamNotification);
		} catch {
			return;
		}
	});
	const decoder = new TextDecoder();
	const reader = response.body.getReader();
	for (;;) {
		const { done, value } = await reader.read();
		if (done) return 'closed';
		parse(decoder.decode(value, { stream: true }));
	}
}

/**
 * Keeps one notification stream open for this window, reconnecting with backoff. It stops on
 * a 401, so signing out ends it and the next sign-in starts a fresh one.
 */
export function listenDesktop(): void {
	if (listening) return;
	listening = true;
	let backoff = 1000;
	const loop = async () => {
		for (;;) {
			const end = await readStream().catch((): StreamEnd => 'failed');
			if (end === 'signed-out') break;
			backoff = end === 'closed' ? 1000 : Math.min(backoff * 2, MAX_BACKOFF_MS);
			await new Promise((resolve) => setTimeout(resolve, backoff));
		}
		listening = false;
	};
	void loop();
}
