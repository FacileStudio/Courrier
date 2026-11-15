import { describe, expect, it } from 'vitest';
import { capsuleHtml, capsuleText, decryptText, encryptText, envelope } from './capsule';
import { fromBase64 } from './base64';

describe('capsule', () => {
	it('seals and opens a round trip in Capsule format', async () => {
		const { ciphertext, key } = await encryptText('the door code is 4521 ✓');
		expect(key).not.toMatch(/[+/=]/);
		expect(fromBase64(key).length).toBe(32);
		expect(fromBase64(ciphertext).length).toBeGreaterThan(12 + 16);
		expect(await decryptText(ciphertext, key)).toBe('the door code is 4521 ✓');
	});

	it('refuses the wrong key', async () => {
		const { ciphertext } = await encryptText('secret');
		const other = await encryptText('x');
		await expect(decryptText(ciphertext, other.key)).rejects.toThrow();
	});

	it('writes the envelope Capsule reads', () => {
		expect(envelope('hi')).toBe('{"__capsule":1,"text":"hi"}');
	});

	it('renders the link for HTML and plain mails', () => {
		const link = 'https://capsule.facile.studio/cap_0123456789abcdef#a-b_c';
		expect(capsuleHtml('Secret <link>:', link)).toBe(
			`<p>Secret &lt;link&gt;: <a href="${link}">${link}</a></p>`
		);
		expect(capsuleText('Secret link:', link)).toBe(`Secret link: ${link}`);
	});
});
