import { render } from 'svelte/server';
import { describe, expect, it } from 'vitest';
import { APK_SHA256_URL, APK_URL, DESKTOP_BUILDS, DESKTOP_CHECKSUMS_URL, RELEASES_URL, detectOS } from '$lib/downloads';
import Page from './+page.svelte';

describe('downloads page', () => {
	it('links the latest APK, every desktop installer, the checksums and the releases page', () => {
		const html = render(Page).body;
		expect(APK_URL).toBe('https://github.com/FacileStudio/Courrier/releases/latest/download/courrier-android.apk');
		for (const url of [APK_URL, APK_SHA256_URL, RELEASES_URL, DESKTOP_CHECKSUMS_URL]) {
			expect(html).toContain(`href="${url}"`);
		}
		for (const build of DESKTOP_BUILDS) expect(html).toContain(`href="${build.url}"`);
	});
});

describe('detectOS', () => {
	it('names the desktop system and ignores phones', () => {
		expect(detectOS('Mozilla/5.0 (Macintosh; Intel Mac OS X 14_5)')).toBe('macos');
		expect(detectOS('Mozilla/5.0 (Windows NT 10.0; Win64; x64)')).toBe('windows');
		expect(detectOS('Mozilla/5.0 (X11; Linux x86_64)')).toBe('linux');
		expect(detectOS('Mozilla/5.0 (Linux; Android 14; Pixel 8)')).toBeNull();
		expect(detectOS('Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X)')).toBeNull();
	});
});
