<script lang="ts">
	import { onMount } from 'svelte';
	import {
		Alert,
		Button,
		ConfirmModal,
		Drawer,
		Field,
		Input,
		SecretField,
		SettingsSection,
		Spinner,
		Table,
		icons
	} from '@facile/muse';
	import { backend, type ApiTokenStatus } from '$lib/backend';

	let status = $state<ApiTokenStatus>({ has_token: false });
	let loadError = $state('');

	let createOpen = $state(false);
	let creating = $state(false);
	let createError = $state('');
	let newName = $state('');
	let issued = $state('');

	let revokeOpen = $state(false);
	let revokeError = $state('');

	const endpoint = `${backend.baseUrl || (typeof location === 'undefined' ? '' : location.origin)}/api`;

	onMount(load);

	async function load() {
		try {
			status = await backend.apiTokenStatus();
			loadError = '';
		} catch (err) {
			loadError = err instanceof Error ? err.message : 'Could not read your token';
		}
	}

	/* Reopening the drawer must never re-show the token from a previous run. */
	function openCreate() {
		issued = '';
		newName = '';
		createError = '';
		createOpen = true;
	}

	async function create(event: Event) {
		event.preventDefault();
		creating = true;
		createError = '';
		try {
			issued = (await backend.createApiToken(newName.trim())).token;
			await load();
		} catch (err) {
			createError = err instanceof Error ? err.message : 'Could not create the token';
		}
		creating = false;
	}

	async function revoke() {
		revokeError = '';
		try {
			await backend.deleteApiToken();
			await load();
		} catch (err) {
			revokeError = err instanceof Error ? err.message : 'Could not revoke the token';
		}
	}

	function when(raw?: string) {
		if (!raw) return '—';
		const date = new Date(raw);
		return Number.isNaN(date.getTime()) ? '—' : date.toLocaleDateString();
	}
</script>

<div class="flex flex-col gap-10">
	<SettingsSection
		title="Endpoint"
		description="Point a script or an HTTP client here. Not a secret — the token is."
	>
		<SecretField value={endpoint} sensitive={false} label="Base URL" />
	</SettingsSection>

	<SettingsSection
		title="Personal API token"
		description="Authenticates a script as you. Courrier keeps one per person — issuing a new token retires the old one on the spot."
		bare
	>
		{#snippet actions()}
			<Button icon={icons.plus} onclick={openCreate}>
				{status.has_token ? 'Replace token' : 'New token'}
			</Button>
		{/snippet}

		{#if loadError}
			<Alert tone="danger">{loadError}</Alert>
		{/if}

		{#if !status.has_token}
			<Alert tone="info">
				No token yet. Nothing outside the browser can talk to this instance as you until there is
				one.
			</Alert>
		{:else}
			<Table>
				<thead>
					<tr>
						<th>Name</th>
						<th>Token</th>
						<th>Created</th>
						<th class="text-right">Actions</th>
					</tr>
				</thead>
				<tbody>
					<tr>
						<td class="font-medium text-fc-fg">{status.name || 'CLI'}</td>
						<td class="font-fc-mono text-fc-xs text-fc-fg-muted">stored hashed</td>
						<td class="whitespace-nowrap text-fc-fg-muted">{when(status.created_at)}</td>
						<td class="text-right">
							<Button
								variant="ghost-danger"
								size="sm"
								icon={icons.revoke}
								onclick={() => {
									revokeError = '';
									revokeOpen = true;
								}}
							>
								Revoke
							</Button>
						</td>
					</tr>
				</tbody>
			</Table>
		{/if}

		{#if revokeError}
			<Alert tone="danger">{revokeError}</Alert>
		{/if}
	</SettingsSection>
</div>

<Drawer bind:open={createOpen} title={status.has_token && !issued ? 'Replace API token' : 'New API token'}>
	{#if issued}
		<div class="flex flex-col gap-4">
			<Alert tone="warning" title="Copy it now">
				Courrier stores only a hash of this token, so this is the one and only time it is shown.
				Lose it and the fix is issuing another one.
			</Alert>

			<!--
				The one-time token is the exception to the auto-hide rule: it starts revealed and stays
				that way, because hiding a value nobody has copied yet is theatre.
			-->
			<SecretField
				value={issued}
				visible
				autoHideMs={0}
				label="Token"
				helper="Put it in your password manager or your CI secret store."
			/>

			<div class="flex justify-end">
				<Button onclick={() => (createOpen = false)}>Done</Button>
			</div>
		</div>
	{:else}
		<form class="flex flex-col gap-4" onsubmit={create}>
			{#if status.has_token}
				<Alert tone="warning" title="This replaces your current token">
					The token named "{status.name || 'CLI'}" stops working the moment the new one is issued,
					and anything still using it starts failing.
				</Alert>
			{/if}

			<Field label="Name" helper="Where the token will live — a machine, a cron job, a script.">
				<Input bind:value={newName} placeholder="laptop-cli" required disabled={creating} />
			</Field>

			{#if createError}
				<Alert tone="danger">{createError}</Alert>
			{/if}

			<div class="flex justify-end gap-2 pt-1">
				<Button type="button" variant="ghost" disabled={creating} onclick={() => (createOpen = false)}>
					Cancel
				</Button>
				<Button type="submit" disabled={creating || newName.trim().length === 0}>
					{#if creating}<Spinner size="sm" />{/if}
					{creating ? 'Creating…' : 'Create token'}
				</Button>
			</div>
		</form>
	{/if}
</Drawer>

<ConfirmModal
	bind:open={revokeOpen}
	tone="danger"
	title="Revoke this token?"
	description="Every script, cron job and CLI still holding it starts getting 401s immediately, and it cannot be un-revoked — the only way back is issuing a new token and updating whatever used the old one."
	confirmLabel="Revoke token"
	cancelLabel="Keep it"
	onConfirm={revoke}
/>
