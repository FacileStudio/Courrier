import { render } from 'svelte/server';
import { describe, expect, it } from 'vitest';
import { APK_SHA256_URL, APK_URL, DESKTOP_CHECKSUMS_URL, DOWNLOADS, RELEASES_URL } from '$lib/downloads';
import Page from './+page.svelte';

describe('downloads page', () => {
	it('links the latest APK, every desktop installer, the checksums and the releases page', () => {
		const html = render(Page).body;
		expect(APK_URL).toBe('https://github.com/FacileStudio/Courrier/releases/latest/download/courrier-android.apk');
		for (const url of [APK_URL, APK_SHA256_URL, RELEASES_URL, DESKTOP_CHECKSUMS_URL]) {
			expect(html).toContain(`href="${url}"`);
		}
		for (const entry of DOWNLOADS) expect(html).toContain(`href="${entry.href}"`);
		expect(html).toContain('facile install courrier');
	});
});
