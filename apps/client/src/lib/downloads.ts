import type { CliCommand, DownloadEntry } from '@facile/muse';

const RELEASES = 'https://github.com/FacileStudio/Courrier/releases';
const LATEST = `${RELEASES}/latest/download`;

export const APK_URL = `${LATEST}/courrier-android.apk`;
export const APK_SHA256_URL = `${APK_URL}.sha256`;
export const RELEASES_URL = RELEASES;
export const DESKTOP_CHECKSUMS_URL = `${LATEST}/desktop-checksums.txt`;

export const DOWNLOADS: readonly DownloadEntry[] = [
	{ platform: 'macos', href: `${LATEST}/Courrier-macos-universal.dmg`, label: '.dmg', arch: 'universal' },
	{ platform: 'windows', href: `${LATEST}/Courrier-windows-x64-setup.exe`, label: '.exe', arch: 'x64' },
	{ platform: 'linux', href: `${LATEST}/Courrier-linux-x64.AppImage`, label: 'AppImage', arch: 'x64' },
	{ platform: 'linux', href: `${LATEST}/Courrier-linux-x64.deb`, label: '.deb', arch: 'x64' },
	{ platform: 'android', href: APK_URL, label: 'APK', checksumHref: APK_SHA256_URL }
];

export const CLI_COMMANDS: readonly CliCommand[] = [
	{ label: 'install.sh', command: 'curl -fsSL https://raw.githubusercontent.com/FacileStudio/courrier-cli/main/install.sh | bash' },
	{ label: 'facile', command: 'facile install courrier' }
];
