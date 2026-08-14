<script lang="ts">
	import { onMount, setContext } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { MobileNav, SideBar, SpaceSwitcher, Topbar, icons } from '@facile/muse';
	import { backend, ApiError, type UserProfile, type MailAccount, type Folder, type Space } from '$lib/backend';
	import { MAIL_FOLDERS } from '$lib/mail-folders';
	import { spaceStore } from '$lib/stores/space.svelte';
	import { searchStore } from '$lib/stores/search.svelte';

	let { children } = $props();

	let user = $state<UserProfile | null>(null);
	let loaded = $state(false);
	let accounts = $state<MailAccount[]>([]);
	let defaultAccountId = $state<number | null>(null);
	let folders = $state<Folder[]>([]);
	let spaces = $state<Space[]>([]);
	let collapsed = $state(false);
	let scroller: HTMLElement | null = $state(null);

	function setUser(nextUser: UserProfile) {
		user = nextUser;
	}

	setContext('app', {
		get user() { return user; },
		get accounts() { return accounts; },
		get defaultAccountId() { return defaultAccountId; },
		get folders() { return folders; },
		setUser,
		refreshAccounts
	});

	async function refreshAccounts() {
		try {
			const result = await backend.listAccounts(spaceStore.active?.id);
			accounts = result.accounts;
			const def = accounts.find((a) => a.is_default) ?? accounts[0] ?? null;
			defaultAccountId = def?.id ?? null;

			if (defaultAccountId) {
				try {
					const folderResult = await backend.getFolders(defaultAccountId);
					folders = folderResult.folders;
				} catch {
					folders = [];
				}
			}
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
			backend.syncProfile().then(async (r) => {
				if (r.synced) {
					const fresh = await backend.me();
					user = fresh.user;
				}
			}).catch(() => {});
			await Promise.all([refreshAccounts(), refreshSpaces()]);
		} catch (err) {
			if (err instanceof ApiError && (err.status === 401 || err.status === 403)) {
				goto('/login');
			}
		}
	});

	$effect(() => {
		spaceStore.active;
		if (loaded) refreshAccounts();
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

	/* <main> is the only scroller in the shell and it outlives every route, so its scrollTop
	   survives a navigation unless someone puts it back. */
	$effect(() => {
		if (path) scroller?.scrollTo({ top: 0 });
	});

	const onSettings = $derived(path === '/settings' || path.startsWith('/settings/'));

	/*
	 * No Settings row: it hangs off the user card at the bottom of the rail. Compose is a
	 * route, so it is a nav link like the folders rather than a hand-shaped button, and the
	 * unread count rides in the label because `pages` has no badge slot.
	 */
	const pages = $derived([
		{ href: '/mail/compose', label: 'Compose', icon: icons.edit, active: path === '/mail/compose' },
		...MAIL_FOLDERS.map((folder) => {
			const count = unread(folder.type);
			return {
				href: folder.href,
				label: count > 0 ? `${folder.label} · ${count}` : folder.label,
				icon: folder.icon,
				active: path === folder.href
			};
		})
	]);

	const navUser = $derived(
		user ? { name: user.name?.trim() || user.email, avatar: user.avatar_url || undefined } : undefined
	);

	/* Mail and Compose only. Every target is a fixed 44px square, so a phone at 360px has
	   room for four plus the avatar — folders live behind the switcher in the mail header and
	   spaces behind the topbar switcher rather than eating rows here. */
	const mobileItems = $derived([
		{
			href: '/mail',
			label: 'Mail',
			icon: icons.mail,
			active: path === '/mail' || (path.startsWith('/mail/') && path !== '/mail/compose')
		},
		{ href: '/mail/compose', label: 'Compose', icon: icons.edit, active: path === '/mail/compose' }
	]);
</script>

{#if loaded}
	<div class="flex h-dvh w-full overflow-hidden bg-fc-page">
		<div class="hidden h-full shrink-0 p-3 md:block">
			<SideBar
				icon="solar:letter-bold-duotone"
				title="Courrier"
				bind:collapsed
				showSearch
				onSearch={openSearch}
				{pages}
				spaces={spaces.map((s) => ({ id: s.id, name: s.name }))}
				activeSpaceId={spaceStore.active?.id ?? null}
				onSpaceSelect={selectSpace}
				manageSpacesHref="/spaces"
				personalSpaceLabel="Personal"
				manageSpacesLabel="Manage spaces"
				user={navUser}
				userHref="/settings"
				userActive={onSettings}
				class="h-full"
			/>
		</div>

		<div class="flex min-w-0 flex-1 flex-col">
			<!-- Spaces live in the rail and the rail is desktop-only, so without this header a
			     phone cannot switch space at all. It sits outside <main> on purpose: pages here
			     are `h-full` panes with their own scrollers, and a header inside the scroller
			     would push them past the bottom of the viewport. -->
			{#if spaces.length > 0}
				<Topbar class="shrink-0 md:hidden">
					<span class="text-fc-md font-semibold text-fc-fg">Courrier</span>
					<div class="min-w-0 max-w-56 flex-1">
						<SpaceSwitcher
							spaces={spaces.map((s) => ({ id: s.id, name: s.name }))}
							activeId={spaceStore.active?.id ?? null}
							onSelect={selectSpace}
							personalLabel="Personal"
							manageHref="/spaces"
						/>
					</div>
				</Topbar>
			{/if}

			<!-- The one scroller in the shell. `overscroll-contain` stops a flick past either
			     end from chaining into the document and rubber-banding the whole app. -->
			<main bind:this={scroller} class="min-w-0 flex-1 overflow-auto overscroll-contain pb-28 md:pb-0">
				{@render children()}
			</main>
		</div>

		<MobileNav items={mobileItems} user={navUser} profileHref="/settings" profileActive={onSettings} />
	</div>
{/if}
