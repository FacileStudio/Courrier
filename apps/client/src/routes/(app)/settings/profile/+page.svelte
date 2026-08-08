<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import {
		Alert,
		Button,
		Dropzone,
		Input,
		ProfileCard,
		SettingsRow,
		SettingsSection,
		UploadProgress,
		icons
	} from '@facile/muse';
	import { backend, type UserProfile } from '$lib/backend';
	import { spaceStore } from '$lib/stores/space.svelte';

	type Upload = {
		id: string;
		name: string;
		size?: number;
		progress: number;
		status: 'pending' | 'uploading' | 'done' | 'error';
		error?: string;
	};

	const app = getContext<{
		user: UserProfile | null;
		setUser: (user: UserProfile) => void;
	}>('app');

	let name = $state(app.user?.name ?? '');
	let email = $state(app.user?.email ?? '');
	let touched = $state(false);
	let ssoManaged = $state(false);

	/*
	 * The layout kicks off an OIDC profile sync in the background, so the user object can be
	 * replaced a beat after this page mounts. Re-seed from it while the fields are untouched —
	 * never once they are, or the sync lands mid-sentence and eats what was being typed.
	 */
	$effect(() => {
		const fresh = app.user;
		if (touched || !fresh) return;
		name = fresh.name ?? '';
		email = fresh.email ?? '';
	});
	let saving = $state(false);
	let saved = $state(false);
	let saveError = $state('');

	let uploads = $state<Upload[]>([]);
	let avatarError = $state('');
	let clearing = $state(false);
	let uploadSeq = 0;

	let loggingOut = $state(false);

	const avatarSource = $derived(app.user?.avatar_source ?? '');
	const fromSSO = $derived(avatarSource === 'oidc');

	const memberSince = $derived.by(() => {
		const raw = app.user?.created_at;
		if (!raw) return '';
		const date = new Date(raw);
		return Number.isNaN(date.getTime())
			? ''
			: date.toLocaleDateString(undefined, { month: 'long', year: 'numeric' });
	});

	const meta = $derived(
		[
			memberSince ? { label: 'Member since', value: memberSince } : null,
			spaceStore.active ? { label: 'Active space', value: spaceStore.active.name } : null,
			{ label: 'Sign-in', value: ssoManaged ? 'Facile SSO' : 'Password on this instance' }
		].filter((entry): entry is { label: string; value: string } => entry !== null)
	);

	const rejections: Record<string, string> = {
		type: 'that file is not a PNG, JPEG, GIF or WebP image',
		size: 'that file is larger than 5 MB',
		count: 'only one photo at a time'
	};

	onMount(async () => {
		try {
			const config = await backend.authConfig();
			ssoManaged = config.oidc_enabled;
		} catch {
			ssoManaged = false;
		}
	});

	async function saveIdentity() {
		saving = true;
		saved = false;
		saveError = '';
		try {
			const payload: { name?: string; email?: string } = { name: name.trim() };
			if (!ssoManaged && email.trim() !== app.user?.email) payload.email = email.trim();
			app.setUser(await backend.updateMe(payload));
			touched = false;
			saved = true;
		} catch (err) {
			saveError = err instanceof Error ? err.message : 'Could not save your details';
		}
		saving = false;
	}

	async function upload(files: File[]) {
		const file = files[0];
		if (!file) return;
		avatarError = '';
		const id = `avatar-${++uploadSeq}`;
		uploads = [{ id, name: file.name, size: file.size, progress: 0, status: 'uploading' }];
		try {
			const user = await backend.uploadAvatar(file, (percent) => {
				const item = uploads.find((entry) => entry.id === id);
				if (item) item.progress = percent;
			});
			const item = uploads.find((entry) => entry.id === id);
			if (item) {
				item.progress = 100;
				item.status = 'done';
			}
			app.setUser(user);
		} catch (err) {
			const item = uploads.find((entry) => entry.id === id);
			const message = err instanceof Error ? err.message : 'Upload failed';
			if (item) {
				item.status = 'error';
				item.error = message;
			}
			avatarError = message;
		}
	}

	async function removePhoto() {
		clearing = true;
		avatarError = '';
		try {
			app.setUser(await backend.clearAvatar());
			uploads = [];
		} catch (err) {
			avatarError = err instanceof Error ? err.message : 'Could not remove the photo';
		}
		clearing = false;
	}

	async function logOut() {
		loggingOut = true;
		try {
			await backend.logout();
		} catch {
			/* The cookie is gone either way — never trap someone on a page they asked to leave. */
		}
		await goto('/login');
	}
</script>

<div class="flex flex-col gap-10">
	<ProfileCard
		name={name.trim() || app.user?.email || 'You'}
		email={app.user?.email}
		avatar={app.user?.avatar_url || undefined}
		role={spaceStore.active?.role}
		{meta}
	/>

	<SettingsSection title="Identity" description="How Courrier names you when you send mail.">
		<SettingsRow
			label="Display name"
			description="Used as the From name on messages you send and on your sidebar card."
			for="profile-name"
			stacked
		>
			<Input
				id="profile-name"
				bind:value={name}
				maxlength={80}
				placeholder="Your name"
				class="sm:w-80"
				oninput={() => (touched = true)}
			/>
		</SettingsRow>

		<SettingsRow
			label="Email"
			description={ssoManaged
				? 'Managed by Facile SSO — change it at porte.facile.studio and it syncs back here.'
				: 'The address you sign in with. Mail accounts have their own addresses.'}
			for="profile-email"
			stacked
		>
			<Input
				id="profile-email"
				bind:value={email}
				type="email"
				disabled={ssoManaged}
				class="sm:w-80"
				oninput={() => (touched = true)}
			/>
		</SettingsRow>

		{#if saveError}
			<Alert tone="danger">{saveError}</Alert>
		{:else if saved}
			<Alert tone="success">Saved.</Alert>
		{/if}

		<Button
			icon={icons.check}
			class="self-start"
			disabled={saving || name.trim().length === 0}
			onclick={saveIdentity}
		>
			{saving ? 'Saving…' : 'Save changes'}
		</Button>
	</SettingsSection>

	<!--
		Two sources, one derived value: an SSO photo always wins, so while one exists the
		upload is not merely outranked — the endpoint refuses it. Showing the dropzone anyway
		would spend a 4xx to say what the copy can say for free.
	-->
	<SettingsSection
		title="Photo"
		description="Shown on your sidebar card. Square images crop best."
		bare
	>
		{#if fromSSO}
			<Alert tone="info" title="Your photo comes from single sign-on">
				Change it at porte.facile.studio and it lands here within a few minutes. Courrier keeps
				no copy, so there is nothing to upload or delete on this side.
			</Alert>
		{:else}
			<Dropzone
				accept="image/png,image/jpeg,image/gif,image/webp"
				maxSize={5 * 1024 * 1024}
				label="Drop a photo here"
				hint="PNG, JPEG, GIF or WebP — up to 5 MB"
				onFiles={upload}
				onReject={(rejected) =>
					(avatarError = `Not uploaded — ${rejections[rejected[0].reason] ?? 'that file was rejected'}.`)}
			/>

			{#if avatarError}
				<Alert tone="danger">{avatarError}</Alert>
			{/if}

			{#if uploads.length > 0}
				<UploadProgress items={uploads} showTotal={false} onCancel={() => (uploads = [])} />
			{/if}

			{#if avatarSource === 'upload'}
				<Button
					variant="outline"
					icon={icons.remove}
					class="self-start"
					disabled={clearing}
					onclick={removePhoto}
				>
					{clearing ? 'Removing…' : 'Remove photo'}
				</Button>
			{/if}
		{/if}
	</SettingsSection>

	<SettingsSection
		title="Session"
		description={ssoManaged
			? 'Signed in through Facile SSO at porte.facile.studio.'
			: 'Signed in with a password on this instance.'}
	>
		<SettingsRow
			label="Log out"
			description="Ends this session in this browser. Other browsers stay signed in."
		>
			<Button variant="outline" icon={icons.logout} disabled={loggingOut} onclick={logOut}>
				{loggingOut ? 'Logging out…' : 'Log out'}
			</Button>
		</SettingsRow>
	</SettingsSection>
</div>
