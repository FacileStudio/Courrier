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
	import { formatDate } from '$lib/i18n/format';
	import { t } from '$lib/i18n/index.svelte';
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
		// Reload when the active space changes.
		spaceStore.active;
		void load();
	});

	async function load() {
		try {
			templates = (await backend.listTemplates(spaceStore.active?.id)).templates;
			loadError = '';
		} catch (err) {
			templates = [];
			loadError = err instanceof Error ? err.message : t('templates.loadFailed');
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
			saveError = err instanceof Error ? err.message : t('templates.saveFailed');
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
			removeError = err instanceof Error ? err.message : t('templates.deleteFailed');
		}
	}

	function when(raw: string) {
		const date = new Date(raw);
		return Number.isNaN(date.getTime()) ? '—' : formatDate(date);
	}
</script>

<div class="flex flex-col gap-10">
	<SettingsSection
		title={t('templates.title')}
		description={t('templates.body')}
		bare
	>
		{#snippet actions()}
			<Button icon={icons.plus} onclick={() => open(null)}>{t('templates.new')}</Button>
		{/snippet}

		{#if loadError}
			<Alert tone="danger">{loadError}</Alert>
		{/if}

		{#if templates.length === 0}
			<EmptyState
				icon={icons.edit}
				title={t('templates.empty')}
				description={t('templates.emptyBody')}
			>
				<Button icon={icons.plus} onclick={() => open(null)}>{t('templates.writeFirst')}</Button>
			</EmptyState>
		{:else}
			<Table>
				<thead>
					<tr>
						<th>{t('templates.colName')}</th>
						<th>{t('templates.colSubject')}</th>
						<th>{t('templates.colUpdated')}</th>
						<th class="text-right">{t('templates.colActions')}</th>
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
										{t('common.edit')}
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
										{t('common.delete')}
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

<Drawer bind:open={editorOpen} title={editing ? t('templates.editTitle', { name: editing.name }) : t('templates.new')}>
	<form class="flex flex-col gap-4" onsubmit={save}>
		<Field label={t('templates.name')} helper={t('templates.nameHelp')}>
			<Input bind:value={name} placeholder={t('templates.namePlaceholder')} required disabled={saving} />
		</Field>

		<Field label={t('templates.subject')} helper={t('templates.subjectHelp')}>
			<Input bind:value={subject} placeholder={t('templates.subjectPlaceholder')} disabled={saving} />
		</Field>

		<Field label={t('templates.bodyLabel')}>
			<Textarea bind:value={body} rows={10} placeholder={t('templates.bodyPlaceholder')} disabled={saving} />
		</Field>

		{#if saveError}
			<Alert tone="danger">{saveError}</Alert>
		{/if}

		<div class="flex justify-end gap-2 pt-1">
			<Button type="button" variant="ghost" disabled={saving} onclick={() => (editorOpen = false)}>
				{t('common.cancel')}
			</Button>
			<Button type="submit" disabled={saving || name.trim().length === 0}>
				{#if saving}<Spinner size="sm" />{/if}
				{saving ? t('common.saving') : editing ? t('templates.save') : t('templates.create')}
			</Button>
		</div>
	</form>
</Drawer>

<ConfirmModal
	bind:open={removeOpen}
	tone="danger"
	title={t('templates.deleteTitle')}
	description={t('templates.deleteBody', { name: removeTarget?.name ?? '' })}
	confirmLabel={t('templates.deleteConfirm')}
	cancelLabel={t('accounts.keepIt')}
	onConfirm={remove}
/>
