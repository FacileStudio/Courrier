<script lang="ts">
	import { getContext } from 'svelte';
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
	import { formatDate, slots } from '$lib/i18n/format';
	import { t } from '$lib/i18n/index.svelte';
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
	let touched = $state(false);

	/*
	 * The layout kicks off an OIDC profile sync in the background, so the user object can be
	 * replaced a beat after this page mounts. Re-seed from it while the fields are untouched —
	 * never once they are, or the sync lands mid-sentence and eats what was being typed.
	 */
	$effect(() => {
		const fresh = app.user;
		if (touched || !fresh) return;
		name = fresh.name ?? '';
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
			: formatDate(date, { month: 'long', year: 'numeric' });
	});

	const meta = $derived(
		[
			memberSince ? { label: t('profile.memberSince'), value: memberSince } : null,
			spaceStore.active ? { label: t('profile.activeSpace'), value: spaceStore.active.name } : null,
			{ label: t('profile.signIn'), value: t('profile.signInValue') }
		].filter((entry): entry is { label: string; value: string } => entry !== null)
	);

	const rejections = $derived<Record<string, string>>({
		type: t('profile.rejectType'),
		size: t('profile.rejectSize'),
		count: t('profile.rejectCount')
	});

	async function saveIdentity() {
		saving = true;
		saved = false;
		saveError = '';
		try {
			app.setUser(await backend.updateMe({ name: name.trim() }));
			touched = false;
			saved = true;
		} catch (err) {
			saveError = err instanceof Error ? err.message : t('profile.saveFailed');
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
			const message = err instanceof Error ? err.message : t('profile.uploadFailed');
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
			avatarError = err instanceof Error ? err.message : t('profile.removeFailed');
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
		name={name.trim() || app.user?.email || t('common.you')}
		email={app.user?.email}
		avatar={app.user?.avatar_url || undefined}
		role={spaceStore.active?.role}
		{meta}
	/>

	<SettingsSection title={t('profile.identity')} description={t('profile.identityBody')}>
		<SettingsRow
			label={t('profile.displayName')}
			description={t('profile.displayNameBody')}
			for="profile-name"
			stacked
		>
			<Input
				id="profile-name"
				bind:value={name}
				maxlength={80}
				placeholder={t('profile.namePlaceholder')}
				class="sm:w-80"
				oninput={() => (touched = true)}
			/>
		</SettingsRow>

		<SettingsRow
			label={t('common.email')}
			description={t('profile.emailBody')}
			for="profile-email"
			stacked
		>
			<Input
				id="profile-email"
				value={app.user?.email ?? ''}
				type="email"
				disabled
				class="sm:w-80"
			/>
		</SettingsRow>

		{#if saveError}
			<Alert tone="danger">{saveError}</Alert>
		{:else if saved}
			<Alert tone="success">{t('profile.saved')}</Alert>
		{/if}

		<Button
			icon={icons.check}
			class="self-start"
			disabled={saving || name.trim().length === 0}
			onclick={saveIdentity}
		>
			{saving ? t('common.saving') : t('profile.save')}
		</Button>
	</SettingsSection>

	<!--
		Two sources, one derived value: an SSO photo always wins, so while one exists the
		upload is not merely outranked — the endpoint refuses it. Showing the dropzone anyway
		would spend a 4xx to say what the copy can say for free.
	-->
	<SettingsSection
		title={t('profile.photo')}
		description={t('profile.photoBody')}
		bare
	>
		{#if fromSSO}
			<Alert tone="info" title={t('profile.photoFromFacile')}>
				{#each slots(t('profile.photoFromFacileBody')) as part, i (i)}
					{#if typeof part === 'string'}
						{part}
					{:else}
						<a href="https://sso.facile.studio/account" class="underline">{t('profile.facileAccount')}</a>
					{/if}
				{/each}
			</Alert>
		{:else}
			<Dropzone
				accept="image/png,image/jpeg,image/gif,image/webp"
				maxSize={5 * 1024 * 1024}
				label={t('profile.drop')}
				hint={t('profile.dropHint')}
				onFiles={upload}
				onReject={(rejected) =>
					(avatarError = t('profile.notUploaded', { reason: rejections[rejected[0].reason] ?? t('profile.rejectOther') }))}
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
					{clearing ? t('profile.removing') : t('profile.removePhoto')}
				</Button>
			{/if}
		{/if}
	</SettingsSection>

	<SettingsSection
		title={t('profile.session')}
		description={t('profile.sessionBody')}
	>
		<SettingsRow label={t('shell.signOut')} description={t('profile.signOutBody')}>
			<Button variant="outline" icon={icons.logout} disabled={loggingOut} onclick={logOut}>
				{loggingOut ? t('profile.signingOut') : t('shell.signOut')}
			</Button>
		</SettingsRow>
	</SettingsSection>
</div>
