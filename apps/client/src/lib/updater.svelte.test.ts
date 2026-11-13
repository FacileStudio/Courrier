import { beforeAll, beforeEach, describe, expect, it, vi } from 'vitest';
import pkg from '../../package.json';

const mockIsTauri = vi.fn(() => false);
const mockCheck = vi.fn();
const mockRelaunch = vi.fn();

vi.mock('@tauri-apps/api/core', () => ({
	isTauri: () => mockIsTauri()
}));

vi.mock('@tauri-apps/plugin-updater', () => ({
	check: () => mockCheck()
}));

vi.mock('@tauri-apps/plugin-process', () => ({
	relaunch: () => mockRelaunch()
}));

describe('UpdaterStore', () => {
	beforeAll(async () => {
		await import('./updater.svelte');
	}, 60_000);

	beforeEach(() => {
		vi.clearAllMocks();
		mockIsTauri.mockReturnValue(false);
	});

	it('initializes with default values', async () => {
		const { UpdaterStore } = await import('./updater.svelte');
		const store = new UpdaterStore();
		expect(store.state).toBe('idle');
		expect(store.currentVersion).toBe(pkg.version);
		expect(store.availableVersion).toBeNull();
		expect(store.downloadProgress).toBe(0);
		expect(store.errorMessage).toBeNull();
	});

	it('marks upToDate outside tauri on manual check', async () => {
		const { UpdaterStore } = await import('./updater.svelte');
		const store = new UpdaterStore();
		await store.checkNow(false);
		expect(store.state).toBe('upToDate');
	});

	it('stays idle outside tauri on auto check', async () => {
		const { UpdaterStore } = await import('./updater.svelte');
		const store = new UpdaterStore();
		await store.checkNow(true);
		expect(store.state).toBe('idle');
	});

	it('transitions to available when an update is found', async () => {
		mockIsTauri.mockReturnValue(true);
		mockCheck.mockResolvedValueOnce({
			version: '0.8.0',
			currentVersion: pkg.version
		});

		const { UpdaterStore } = await import('./updater.svelte');
		const store = new UpdaterStore();
		await store.checkNow(false);

		expect(store.state).toBe('available');
		expect(store.availableVersion).toBe('0.8.0');
	});

	it('transitions to upToDate when no update is returned', async () => {
		mockIsTauri.mockReturnValue(true);
		mockCheck.mockResolvedValueOnce(null);

		const { UpdaterStore } = await import('./updater.svelte');
		const store = new UpdaterStore();
		await store.checkNow(false);

		expect(store.state).toBe('upToDate');
		expect(store.availableVersion).toBeNull();
	});

	it('records error on manual check failure', async () => {
		mockIsTauri.mockReturnValue(true);
		mockCheck.mockRejectedValueOnce(new Error('Network offline'));

		const { UpdaterStore } = await import('./updater.svelte');
		const store = new UpdaterStore();
		await store.checkNow(false);

		expect(store.state).toBe('error');
		expect(store.errorMessage).toBe('Network offline');
	});

	it('reverts to idle on auto check failure', async () => {
		mockIsTauri.mockReturnValue(true);
		mockCheck.mockRejectedValueOnce(new Error('Network offline'));

		const { UpdaterStore } = await import('./updater.svelte');
		const store = new UpdaterStore();
		await store.checkNow(true);

		expect(store.state).toBe('idle');
		expect(store.errorMessage).toBeNull();
	});

	it('completes download and tracks progress events', async () => {
		mockIsTauri.mockReturnValue(true);

		let progressCallback: ((event: unknown) => void) | undefined;
		const downloadAndInstall = vi.fn((cb: (event: unknown) => void) => {
			progressCallback = cb;
			progressCallback({
				event: 'Started',
				data: { contentLength: 1000 }
			});
			progressCallback({
				event: 'Progress',
				data: { chunkLength: 500 }
			});
			progressCallback({
				event: 'Finished'
			});
			return Promise.resolve();
		});

		mockCheck.mockResolvedValueOnce({
			version: '0.8.0',
			currentVersion: pkg.version,
			downloadAndInstall
		});

		const { UpdaterStore } = await import('./updater.svelte');
		const store = new UpdaterStore();
		await store.checkNow(false);
		await store.installNow();

		expect(downloadAndInstall).toHaveBeenCalledTimes(1);
		expect(store.downloadProgress).toBe(100);
		expect(store.state).toBe('ready');
	});

	it('records error when download fails', async () => {
		mockIsTauri.mockReturnValue(true);
		mockCheck.mockResolvedValueOnce({
			version: '0.8.0',
			currentVersion: pkg.version,
			downloadAndInstall: vi.fn().mockRejectedValueOnce(new Error('Disk full'))
		});

		const { UpdaterStore } = await import('./updater.svelte');
		const store = new UpdaterStore();
		await store.checkNow(false);
		await store.installNow();

		expect(store.state).toBe('error');
		expect(store.errorMessage).toBe('Disk full');
	});

	it('restarts application when ready in tauri', async () => {
		mockIsTauri.mockReturnValue(true);
		mockRelaunch.mockResolvedValueOnce(undefined);

		const { UpdaterStore } = await import('./updater.svelte');
		const store = new UpdaterStore();
		await store.restartNow();

		expect(mockRelaunch).toHaveBeenCalledTimes(1);
	});
});
