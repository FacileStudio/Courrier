import type { AppKey } from '$lib/i18n/messages';

const RELEASES = 'https://github.com/FacileStudio/Courrier/releases';
const LATEST = `${RELEASES}/latest/download`;

export const APK_ASSET = 'courrier-android.apk';
export const APK_URL = `${LATEST}/${APK_ASSET}`;
export const APK_SHA256_URL = `${APK_URL}.sha256`;
export const RELEASES_URL = RELEASES;
export const DESKTOP_CHECKSUMS_URL = `${LATEST}/desktop-checksums.txt`;

export type DesktopOS = 'macos' | 'windows' | 'linux';

export type DesktopBuild = { os: DesktopOS; label: string; detailKey: AppKey; url: string };

export const DESKTOP_BUILDS: readonly DesktopBuild[] = [
	{ os: 'macos', label: 'macOS', detailKey: 'downloads.detailMacos', url: `${LATEST}/Courrier-macos-universal.dmg` },
	{ os: 'windows', label: 'Windows', detailKey: 'downloads.detailWindows', url: `${LATEST}/Courrier-windows-x64-setup.exe` },
	{ os: 'linux', label: 'Linux', detailKey: 'downloads.detailLinux', url: `${LATEST}/Courrier-linux-x64.AppImage` }
];

export const DEB_URL = `${LATEST}/Courrier-linux-x64.deb`;

/** Guesses the visitor's desktop OS from the user agent, or null on a phone or an unknown system. */
export function detectOS(userAgent: string): DesktopOS | null {
	if (/android|iphone|ipad/i.test(userAgent)) return null;
	if (/mac os x|macintosh/i.test(userAgent)) return 'macos';
	if (/windows/i.test(userAgent)) return 'windows';
	if (/linux|x11/i.test(userAgent)) return 'linux';
	return null;
}
