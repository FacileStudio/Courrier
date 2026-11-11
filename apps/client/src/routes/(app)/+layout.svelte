<script lang="ts">
	import { onMount, setContext } from 'svelte';
	import { goto } from '$app/navigation';
	import { rememberReturn, takeReturn } from '$lib/after-login';
	import { page } from '$app/state';
	import { MobileNav, SideBar, SpaceSwitcher, Topbar, getI18n, icons, ideasNavItem } from '@facile/muse';
	import { backend, ApiError, type UserProfile, type MailAccount, type Folder, type Space } from '$lib/backend';
	import { SPACES_URL } from '$lib/plan-limit';
	import { MAIL_FOLDERS } from '$lib/mail-folders';
	import { spaceStore } from '$lib/stores/space.svelte';
	import { activeAccountStore } from '$lib/stores/account.svelte';
	import { searchStore } from '$lib/stores/search.svelte';
	import { ideasStore } from '$lib/ideas.svelte';
	import AccountSwitcher from '$lib/components/AccountSwitcher.svelte';
		import { syncDevice } from '$lib/notifications';
	import { needsOnboarding, onboardingSeen } from '$lib/onboarding/steps';
	import { isApp, isNative } from '$lib/platform';
	import { rememberLocale, t } from '$lib/i18n/index.svelte';
	import { theme, type ThemePreference } from '$lib/theme.svelte';

	let { children } = $props();

	let user = $state<UserProfile | null>(null);
	let loaded = $state(false);
	let accounts = $state<MailAccount[]>([]);
	let activeAccountId = $state<number | null>(null);
	let folders = $state<Folder[]>([]);
	let spaces = $state<Space[]>([]);
	let collapsed = $state(false);
	let scroller: HTMLElement | null = $state(null);
	let ideasReady = $state(false);
	let invitations = $state<{ count: number; href: string } | null>(null);
	const i18n = getI18n();

	function setUser(nextUser: UserProfile) {
		user = nextUser;
	}

	setContext('app', {
		get user() { return user; },
		get accounts() { return accounts; },
		get activeAccountId() { return activeAccountId; },
		get folders() { return folders; },
		setUser,
		refreshAccounts,
		setActiveAccount
	});

	async function loadFolders(id: number | null) {
		if (!id) {
			folders = [];
			return;
		}
		try {
			const folderResult = await backend.getFolders(id);
			if (id === activeAccountId) folders = folderResult.folders;
		} catch {
			if (id === activeAccountId) folders = [];
		}
	}

	async function setActiveAccount(id: number) {
		if (id === activeAccountId || !accounts.some((a) => a.id === id)) return;
		activeAccountStore.set(spaceStore.active?.id, id);
		activeAccountId = id;
		await loadFolders(id);
	}

	async function refreshAccounts() {
		try {
			const result = await backend.listAccounts(spaceStore.active?.id);
			accounts = result.accounts;
			const stored = activeAccountStore.get(spaceStore.active?.id);
			const next =
				accounts.find((a) => a.id === stored) ??
				accounts.find((a) => a.is_default) ??
				accounts[0] ??
				null;
			activeAccountId = next?.id ?? null;
			await loadFolders(activeAccountId);
		} catch {
			accounts = [];
		}
	}

	async function refreshSpaces() {
		try {
			const result = await backend.listSpaces();
			spaces = result.spaces ?? [];
		} catch {
			spaces = [];
		}
	}

	onMount(async () => {
		try {
			const result = await backend.me();
			user = result.user;
			loaded = true;
			ideasReady = true;
			const back = takeReturn();
			if (back) await goto(back, { replaceState: true });
			backend.syncProfile().then(async (r) => {
				if (r.synced) {
					const fresh = await backend.me();
					user = fresh.user;
				}
			}).catch(() => {});
			await Promise.all([refreshAccounts(), refreshSpaces()]);
			const firstRun = needsOnboarding({ native: isNative(), seen: onboardingSeen(), accounts: accounts.length });
			if (!spaceStore.active && firstRun) {
				await goto('/welcome', { replaceState: true });
				return;
			}
			void backend.getPush().then(syncDevice).catch(() => undefined);
			void backend.invitations().then((answer) => (invitations = answer));
		} catch (err) {
			if (err instanceof ApiError && (err.status === 401 || err.status === 403)) {
				rememberReturn(page.url);
				goto('/login');
			}
		}
	});

	$effect(() => {
		spaceStore.active;
		if (loaded) refreshAccounts();
	});

	$effect(() => {
		const id = spaceStore.active?.id ?? null;
		if (ideasReady) void ideasStore.load(id);
	});

	/* The rail owns ⌘K now — it is the component that renders the chip, so binding the
	   shortcut anywhere else means two handlers racing for the same keypress. */
	async function openSearch() {
		if (!page.url.pathname.startsWith('/mail')) await goto('/mail');
		searchStore.requestFocus();
	}

	function selectSpace(id: string | null) {
		const next = spaces.find((s) => s.id === id);
		if (next) spaceStore.set({ id: next.id, name: next.name, role: next.role });
		else spaceStore.clear();
	}

	function unread(type: string) {
		return folders.find((f) => f.type === type)?.unread_count ?? 0;
	}

	const path = $derived(page.url.pathname);
	const isMailRoute = $derived(path === '/mail' || path.startsWith('/mail/'));

	/* <main> is the only scroller in the shell and it outlives every route, so its scrollTop
	   survives a navigation unless someone puts it back. */
	$effect(() => {
		if (path) scroller?.scrollTo({ top: 0 });
	});

	const onSettings = $derived(path === '/settings' || path.startsWith('/settings/'));

	/* Ideas leads the list in every suite app. A Courrier without Bulle answers 404 and drops it. */
	const ideasItems = $derived(
		ideasStore.available
			? [
					ideasNavItem({
						href: '/ideas',
						active: path === '/ideas',
						count: ideasStore.count,
						running: ideasStore.running,
						i18n
					})
				]
			: []
	);

	const pages = $derived([
		...ideasItems,
		{ href: '/mail/compose', label: t('shell.compose'), icon: icons.edit, active: path === '/mail/compose' },
		...MAIL_FOLDERS.filter((folder) => folder.primary).map((folder) => {
			const count = unread(folder.type);
			const label = t(folder.labelKey);
			return {
				href: folder.href,
				label: count > 0 ? `${label} · ${count}` : label,
				icon: folder.icon,
				active: path === folder.href
			};
		})
	]);

	const navUser = $derived(
		user
			? {
					name: user.name?.trim() || user.email,
					email: user.email,
					avatar: user.avatar_url || undefined,
					secondary: spaceStore.active?.name
				}
			: undefined
	);

	async function signOut() {
		await backend.logout().catch(() => undefined);
		await goto('/login', { replaceState: true });
	}

	const userMenu = $derived({
		items: isApp() ? [] : [{ label: t('shell.desktopApp'), icon: icons.download, href: '/downloads' }],
		subscription: { href: SPACES_URL },
		helpHref: 'https://facile.studio/docs/courrier',
		onSignOut: signOut,
		theme: theme.preference,
		onThemeChange: (next: ThemePreference) => theme.set(next),
		onLocaleChange: rememberLocale,
		invitations
	});

	const mobileItems = $derived([
		...ideasItems,
		{
			href: '/mail',
			label: t('shell.mail'),
			icon: icons.mail,
			active: path === '/mail' || (path.startsWith('/mail/') && path !== '/mail/compose')
		},
		{ href: '/mail/compose', label: t('shell.compose'), icon: icons.edit, active: path === '/mail/compose' }
	]);
</script>

{#if loaded}
	<div class="flex h-dvh w-full overflow-hidden bg-fc-page">
		<div class="hidden h-full shrink-0 flex-col gap-3 p-3 md:flex">
			<div class="flex min-h-0 flex-1 flex-col rounded-fc-lg bg-fc-component">
				<SideBar
					icon="solar:letter-opened-bold-duotone"
					title="Courrier"
					bind:collapsed
					showSearch
					onSearch={openSearch}
					{pages}
					spaces={spaces.map((s) => ({ id: s.id, name: s.name }))}
					activeSpaceId={spaceStore.active?.id ?? null}
					onSpaceSelect={selectSpace}
					manageSpacesHref="/spaces"
					user={navUser}
					userHref="/settings"
					{userMenu}
					class="min-h-0 flex-1 rounded-none bg-transparent"
				/>
				{#if accounts.length > 1}
					<div class="flex w-0 min-w-full shrink-0 flex-col px-3 pb-5">
						<AccountSwitcher rail {collapsed} />
					</div>
				{/if}
			</div>
		</div>

		<div class="flex min-w-0 flex-1 flex-col">
			<!-- Spaces live in the rail and the rail is desktop-only, so without this header a
			     phone cannot switch space at all. It sits outside <main> on purpose: pages here
			     are `h-full` panes with their own scrollers, and a header inside the scroller
			     would push them past the bottom of the viewport. -->
			{#if spaces.length > 0 || accounts.length > 1}
				<Topbar class="shrink-0 md:hidden">
					<span class="text-fc-md font-semibold text-fc-fg">Courrier</span>
					{#if spaces.length > 0}
					<div class="min-w-0 max-w-56 flex-1">
						<SpaceSwitcher
							spaces={spaces.map((s) => ({ id: s.id, name: s.name }))}
							activeId={spaceStore.active?.id ?? null}
							onSelect={selectSpace}
							manageHref="/spaces"
						/>
					</div>
					{/if}
					<div class="ml-auto shrink-0">
						<AccountSwitcher align="end" />
					</div>
				</Topbar>
			{/if}

			<!-- The one scroller in the shell. `overscroll-contain` stops a flick past either
			     end from chaining into the document and rubber-banding the whole app. -->
			<main
				bind:this={scroller}
				class={isMailRoute
					? 'flex min-w-0 flex-1 flex-col overflow-hidden'
					: 'min-w-0 flex-1 overflow-auto overscroll-contain pb-28 md:pb-0'}
			>
				{@render children()}
			</main>
		</div>

		<MobileNav items={mobileItems} user={navUser} profileHref="/settings" profileActive={onSettings} {userMenu} />
	</div>
{/if}
