import { render } from 'svelte/server';
import { describe, expect, it } from 'vitest';
import { APK_SHA256_URL, APK_URL, RELEASES_URL } from '$lib/downloads';
import Page from './+page.svelte';

describe('downloads page', () => {
	it('links the latest APK, its checksum and the releases page', () => {
		const html = render(Page).body;
		expect(APK_URL).toBe('https://github.com/FacileStudio/Courrier/releases/latest/download/courrier-android.apk');
		expect(html).toContain(`href="${APK_URL}"`);
		expect(html).toContain(`href="${APK_SHA256_URL}"`);
		expect(html).toContain(`href="${RELEASES_URL}"`);
		expect(html).not.toMatch(/desktop|macOS|Windows/i);
	});
});
