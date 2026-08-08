<script lang="ts">
	import {
		Alert,
		Button,
		ConfirmModal,
		Drawer,
		EmptyState,
		Field,
		Input,
		SettingsSection,
		Spinner,
		Table,
		Textarea,
		icons
	} from '@facile/muse';
	import { backend, type EmailTemplate } from '$lib/backend';
	import { spaceStore } from '$lib/stores/space.svelte';

	let templates = $state<EmailTemplate[]>([]);
	let loadError = $state('');

	let editorOpen = $state(false);
	let editing = $state<EmailTemplate | null>(null);
	let name = $state('');
	let subject = $state('');
	let body = $state('');
	let saving = $state(false);
	let saveError = $state('');

	let removeTarget = $state<EmailTemplate | null>(null);
	let removeOpen = $state(false);
	let removeError = $state('');

	$effect(() => {
		spaceStore.active; // reload when the active space changes
		void load();
	});

	async function load() {
		try {
			templates = (await backend.listTemplates(spaceStore.active?.id)).templates;
			loadError = '';
		} catch (err) {
			templates = [];
			loadError = err instanceof Error ? err.message : 'Could not load your templates';
		}
	}

	function open(template: EmailTemplate | null) {
		editing = template;
		name = template?.name ?? '';
		subject = template?.subject ?? '';
		body = template?.body_text || template?.body_html || '';
		saveError = '';
		editorOpen = true;
	}

	async function save(event: Event) {
		event.preventDefault();
		saving = true;
		saveError = '';
		try {
			const data = {
				name: name.trim(),
				subject: subject.trim(),
				body_html: body,
				body_text: body.replace(/<[^>]*>/g, '')
			};
			if (editing) {
				await backend.updateTemplate(editing.id, data);
			} else {
				await backend.createTemplate({ ...data, space_id: spaceStore.active?.id });
			}
			editorOpen = false;
			await load();
		} catch (err) {
			saveError = err instanceof Error ? err.message : 'Could not save the template';
		}
		saving = false;
	}

	async function remove() {
		const target = removeTarget;
		if (!target) return;
		removeError = '';
		try {
			await backend.deleteTemplate(target.id);
			removeTarget = null;
			await load();
		} catch (err) {
			removeError = err instanceof Error ? err.message : 'Could not delete the template';
		}
	}

	function when(raw: string) {
		const date = new Date(raw);
		return Number.isNaN(date.getTime()) ? '—' : date.toLocaleDateString();
	}
</script>

<div class="flex flex-col gap-10">
	<SettingsSection
		title="Email templates"
		description="Canned messages you can drop into a compose window instead of retyping them."
		bare
	>
		{#snippet actions()}
			<Button icon={icons.plus} onclick={() => open(null)}>New template</Button>
		{/snippet}

		{#if loadError}
			<Alert tone="danger">{loadError}</Alert>
		{/if}

		{#if templates.length === 0}
			<EmptyState
				icon={icons.edit}
				title="No templates yet"
				description="Anything you type more than twice — an onboarding reply, a quote follow-up — belongs here."
			>
				<Button icon={icons.plus} onclick={() => open(null)}>Write your first one</Button>
			</EmptyState>
		{:else}
			<Table>
				<thead>
					<tr>
						<th>Name</th>
						<th>Subject</th>
						<th>Updated</th>
						<th class="text-right">Actions</th>
					</tr>
				</thead>
				<tbody>
					{#each templates as template (template.id)}
						<tr>
							<td class="font-medium text-fc-fg">{template.name}</td>
							<td class="text-fc-fg-muted">{template.subject || '—'}</td>
							<td class="whitespace-nowrap text-fc-fg-muted">{when(template.updated_at)}</td>
							<td>
								<div class="flex items-center justify-end gap-1">
									<Button
										variant="ghost"
										size="sm"
										icon={icons.edit}
										onclick={() => open(template)}
									>
										Edit
									</Button>
									<Button
										variant="ghost-danger"
										size="sm"
										icon={icons.remove}
										onclick={() => {
											removeTarget = template;
											removeError = '';
											removeOpen = true;
										}}
									>
										Delete
									</Button>
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</Table>
		{/if}

		{#if removeError}
			<Alert tone="danger">{removeError}</Alert>
		{/if}
	</SettingsSection>
</div>

<Drawer bind:open={editorOpen} title={editing ? `Edit ${editing.name}` : 'New template'}>
	<form class="flex flex-col gap-4" onsubmit={save}>
		<Field label="Name" helper="How you will find it in the compose window.">
			<Input bind:value={name} placeholder="Quote follow-up" required disabled={saving} />
		</Field>

		<Field label="Subject line" helper="Optional — leave it empty to keep whatever you typed.">
			<Input bind:value={subject} placeholder="Following up on our quote" disabled={saving} />
		</Field>

		<Field label="Body">
			<Textarea bind:value={body} rows={10} placeholder="Hi there," disabled={saving} />
		</Field>

		{#if saveError}
			<Alert tone="danger">{saveError}</Alert>
		{/if}

		<div class="flex justify-end gap-2 pt-1">
			<Button type="button" variant="ghost" disabled={saving} onclick={() => (editorOpen = false)}>
				Cancel
			</Button>
			<Button type="submit" disabled={saving || name.trim().length === 0}>
				{#if saving}<Spinner size="sm" />{/if}
				{saving ? 'Saving…' : editing ? 'Save template' : 'Create template'}
			</Button>
		</div>
	</form>
</Drawer>

<ConfirmModal
	bind:open={removeOpen}
	tone="danger"
	title="Delete this template?"
	description={`"${removeTarget?.name ?? ''}" disappears from every compose window in this space, for everyone who shares it. Messages already sent from it are untouched, and there is no undo.`}
	confirmLabel="Delete template"
	cancelLabel="Keep it"
	onConfirm={remove}
/>
