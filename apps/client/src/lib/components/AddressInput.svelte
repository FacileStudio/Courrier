<script lang="ts">
	import { backend } from '$lib/backend';
	import { Badge, getFieldContext, icons } from '@facile/muse';

	type Contact = { name: string; email: string; count: number };

	let {
		value = '',
		onchange,
		id = '',
		accountId = 0,
		placeholder = ''
	}: {
		value: string;
		onchange: (value: string) => void;
		id?: string;
		accountId?: number;
		placeholder?: string;
	} = $props();

	/* Inside a muse `Field`, adopt the id it generated so its <label> points at this input. */
	const field = getFieldContext();
	const controlId = $derived(id || field?.().id);

	let suggestions = $state<Contact[]>([]);
	let selectedIndex = $state(-1);
	let showDropdown = $state(false);
	let inputEl = $state<HTMLInputElement | null>(null);
	let containerEl = $state<HTMLDivElement | null>(null);
	let debounceTimer: ReturnType<typeof setTimeout> | undefined;

	/*
	 * The value stays what every caller already expects: one comma-separated string. What
	 * changes is how it reads — every completed address is a chip, and the trailing token is
	 * the one still being typed, which is exactly the split the contact search already used.
	 */
	const parts = $derived(value.split(','));
	const chips = $derived(
		parts
			.slice(0, -1)
			.map((s) => s.trim())
			.filter(Boolean)
	);
	const draft = $derived((parts[parts.length - 1] ?? '').replace(/^\s+/, ''));

	function commit(nextChips: string[], nextDraft: string) {
		onchange([...nextChips, nextDraft].join(', '));
	}

	function search(query: string) {
		if (debounceTimer) clearTimeout(debounceTimer);
		if (!query || query.length < 2 || !accountId) {
			suggestions = [];
			showDropdown = false;
			return;
		}
		debounceTimer = setTimeout(async () => {
			try {
				const res = await backend.searchContacts(accountId, query);
				suggestions = res.contacts ?? [];
				selectedIndex = -1;
				showDropdown = suggestions.length > 0;
			} catch {
				suggestions = [];
				showDropdown = false;
			}
		}, 200);
	}

	function selectContact(contact: Contact) {
		commit([...chips, contact.email], '');
		suggestions = [];
		showDropdown = false;
		selectedIndex = -1;
		inputEl?.focus();
	}

	function chipifyDraft(): boolean {
		const trimmed = draft.trim();
		if (!trimmed) return false;
		commit([...chips, trimmed], '');
		suggestions = [];
		showDropdown = false;
		selectedIndex = -1;
		return true;
	}

	function removeChip(index: number) {
		commit(
			chips.filter((_, i) => i !== index),
			draft
		);
		inputEl?.focus();
	}

	function handleInput(e: Event) {
		const target = e.target as HTMLInputElement;
		commit(chips, target.value);
		search(target.value.trim());
	}

	function handleKeydown(e: KeyboardEvent) {
		if (showDropdown) {
			if (e.key === 'ArrowDown') {
				e.preventDefault();
				selectedIndex = Math.min(selectedIndex + 1, suggestions.length - 1);
				return;
			}
			if (e.key === 'ArrowUp') {
				e.preventDefault();
				selectedIndex = Math.max(selectedIndex - 1, 0);
				return;
			}
			if (e.key === 'Enter' && selectedIndex >= 0) {
				e.preventDefault();
				selectContact(suggestions[selectedIndex]);
				return;
			}
			if (e.key === 'Escape') {
				showDropdown = false;
				selectedIndex = -1;
				return;
			}
		}

		if (e.key === 'Enter' || e.key === ',' || e.key === 'Tab') {
			if (chipifyDraft() && e.key !== 'Tab') e.preventDefault();
			else if (e.key === ',') e.preventDefault();
			return;
		}

		if (e.key === 'Backspace' && draft === '' && chips.length > 0) {
			e.preventDefault();
			removeChip(chips.length - 1);
		}
	}

	function handleClickOutside(e: MouseEvent) {
		if (containerEl && !containerEl.contains(e.target as Node)) {
			showDropdown = false;
			selectedIndex = -1;
		}
	}

	$effect(() => {
		document.addEventListener('mousedown', handleClickOutside);
		return () => document.removeEventListener('mousedown', handleClickOutside);
	});
</script>

<!--
	muse has no chips/token input, so the shell is hand-built from the same tokens `Input`
	uses (min-h-11, rounded-fc-md, border-fc-border, bg-fc-bg, the focus ring) and the chips
	themselves are muse `Badge`s.
-->
<div class="relative" bind:this={containerEl}>
	<div
		class="flex min-h-11 w-full flex-wrap items-center gap-1.5 rounded-fc-md border border-fc-border bg-fc-bg px-2 py-1.5 focus-within:outline-2 focus-within:outline-offset-2 focus-within:outline-fc-ring"
	>
		{#each chips as chip, i (`${chip}-${i}`)}
			<Badge tone="neutral" class="max-w-full pr-1">
				<span class="truncate">{chip}</span>
				<button
					type="button"
					aria-label="Remove {chip}"
					class="flex size-4 shrink-0 items-center justify-center rounded-fc-pill text-fc-fg-muted transition-colors hover:bg-fc-component hover:text-fc-fg focus-visible:outline-2 focus-visible:outline-offset-1 focus-visible:outline-fc-ring"
					onclick={() => removeChip(i)}
				>
					<iconify-icon icon={icons.close} width="12" height="12" class="block size-3"></iconify-icon>
				</button>
			</Badge>
		{/each}
		<input
			id={controlId}
			type="text"
			value={draft}
			placeholder={chips.length === 0 ? placeholder : ''}
			oninput={handleInput}
			onkeydown={handleKeydown}
			onblur={() => chipifyDraft()}
			onfocus={() => {
				if (suggestions.length > 0) showDropdown = true;
			}}
			bind:this={inputEl}
			autocomplete="off"
			class="min-w-32 flex-1 bg-transparent px-1 text-fc-sm text-fc-fg placeholder:text-fc-fg-muted focus-visible:outline-none"
		/>
	</div>

	{#if showDropdown && suggestions.length > 0}
		<div
			class="absolute left-0 top-full z-50 mt-1 w-full rounded-fc-md border border-fc-border bg-fc-component p-1 shadow-lg"
			role="listbox"
		>
			{#each suggestions as contact, i (contact.email)}
				<button
					type="button"
					role="option"
					aria-selected={i === selectedIndex}
					class="flex w-full flex-col items-start gap-0.5 rounded-fc-sm px-2.5 py-2 text-left transition-colors hover:bg-fc-surface {i ===
					selectedIndex
						? 'bg-fc-surface'
						: ''}"
					onmousedown={(e) => {
						e.preventDefault();
						selectContact(contact);
					}}
				>
					{#if contact.name}
						<span class="truncate text-fc-sm font-medium text-fc-fg">{contact.name}</span>
						<span class="truncate text-fc-xs text-fc-fg-muted">{contact.email}</span>
					{:else}
						<span class="truncate text-fc-sm text-fc-fg">{contact.email}</span>
					{/if}
				</button>
			{/each}
		</div>
	{/if}
</div>
