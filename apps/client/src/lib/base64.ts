/** Encodes bytes as base64, or base64url without padding, the way Capsule does. */
export function toBase64(bytes: Uint8Array, urlSafe = false): string {
	let binary = '';
	for (let i = 0; i < bytes.length; i += 8192) binary += String.fromCharCode(...bytes.subarray(i, i + 8192));
	const b64 = btoa(binary);
	return urlSafe ? b64.replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '') : b64;
}

/** Decodes base64 or base64url. */
export function fromBase64(str: string): Uint8Array<ArrayBuffer> {
	const binary = atob(str.replace(/-/g, '+').replace(/_/g, '/'));
	const bytes = new Uint8Array(binary.length);
	for (let i = 0; i < binary.length; i++) bytes[i] = binary.charCodeAt(i);
	return bytes;
}
