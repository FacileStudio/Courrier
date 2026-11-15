import { backend } from '$lib/backend';
import { fromBase64, toBase64 } from './base64';

export const CAPSULE_EXPIRIES = ['1h', '24h', '7d', '30d'] as const;
export type CapsuleExpiry = (typeof CAPSULE_EXPIRIES)[number];

/** Capsule's plaintext envelope for a text-only capsule. */
export function envelope(text: string): string {
	return JSON.stringify({ __capsule: 1, text });
}

/** Encrypts text with a fresh AES-256-GCM key: base64(iv || ciphertext), and the base64url key. */
export async function encryptText(text: string): Promise<{ ciphertext: string; key: string }> {
	const key = await crypto.subtle.generateKey({ name: 'AES-GCM', length: 256 }, true, ['encrypt', 'decrypt']);
	const iv = crypto.getRandomValues(new Uint8Array(12));
	const sealed = await crypto.subtle.encrypt({ name: 'AES-GCM', iv }, key, new TextEncoder().encode(envelope(text)));
	const combined = new Uint8Array(iv.length + sealed.byteLength);
	combined.set(iv, 0);
	combined.set(new Uint8Array(sealed), iv.length);
	const raw = new Uint8Array(await crypto.subtle.exportKey('raw', key));
	return { ciphertext: toBase64(combined), key: toBase64(raw, true) };
}

/** Decrypts a Capsule ciphertext with its base64url key and returns the enveloped text. */
export async function decryptText(ciphertext: string, key: string): Promise<string> {
	const imported = await crypto.subtle.importKey('raw', fromBase64(key), { name: 'AES-GCM', length: 256 }, false, ['decrypt']);
	const combined = fromBase64(ciphertext);
	const plain = await crypto.subtle.decrypt({ name: 'AES-GCM', iv: combined.slice(0, 12) }, imported, combined.slice(12));
	const parsed = JSON.parse(new TextDecoder().decode(plain)) as { __capsule?: number; text?: unknown };
	if (parsed.__capsule !== 1 || typeof parsed.text !== 'string') throw new Error('not a capsule envelope');
	return parsed.text;
}

/** Seals text in Capsule through Courrier and returns the full link, key in the fragment. */
export async function sealLink(text: string, expiresIn: CapsuleExpiry): Promise<string> {
	const { ciphertext, key } = await encryptText(text);
	const sealed = await backend.sealCapsule(ciphertext, expiresIn);
	return `${sealed.url}#${key}`;
}

function escapeHtml(s: string): string {
	return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
}

/** The paragraph a capsule link takes in an HTML mail. */
export function capsuleHtml(label: string, link: string): string {
	const href = escapeHtml(link);
	return `<p>${escapeHtml(label)} <a href="${href}">${href}</a></p>`;
}

/** The line a capsule link takes in a plain text mail. */
export function capsuleText(label: string, link: string): string {
	return `${label} ${link}`;
}
