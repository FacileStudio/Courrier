<script lang="ts">
	import { backend } from '$lib/backend';
	import { Label } from '$lib/components/ui/label';

	type Contact = { name: string; email: string; count: number };

	let {
		value = '',
		onchange,
		label = 'To',
		id = '',
		token = '',
		accountId = 0,
		placeholder = ''
	}: {
		value: string;
		onchange: (value: string) => void;
		label?: string;
		id?: string;
		token?: string;
		accountId?: number;
		placeholder?: string;
	} = $props();

	let suggestions = $state<Contact[]>([]);
	let selectedIndex = $state(-1);
	let showDropdown = $state(false);
	let inputEl = $state<HTMLInputElement | null>(null);
	let containerEl = $state<HTMLDivElement | null>(null);
	let debounceTimer: ReturnType<typeof setTimeout> | undefined;

	function currentQuery(): string {
		const parts = value.split(',');
		return (parts[parts.length - 1] ?? '').trim();
	}

	function search(query: string) {
		if (debounceTimer) clearTimeout(debounceTimer);
		if (!query || query.length < 2 || !token || !accountId) {
			suggestions = [];
			showDropdown = false;
			return;
		}
		debounceTimer = setTimeout(async () => {
			try {
				const res = await backend.searchContacts(token, accountId, query);
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
		const parts = value.split(',').map((s) => s.trim()).filter(Boolean);
		parts.pop();
		const display = contact.name ? `${contact.name} <${contact.email}>` : contact.email;
		parts.push(display);
		const newValue = parts.join(', ') + ', ';
		onchange(newValue);
		suggestions = [];
		showDropdown = false;
		selectedIndex = -1;
		inputEl?.focus();
	}

	function handleInput(e: Event) {
		const target = e.target as HTMLInputElement;
		onchange(target.value);
		search(currentQuery());
	}

	function handleKeydown(e: KeyboardEvent) {
		if (!showDropdown) return;

		if (e.key === 'ArrowDown') {
			e.preventDefault();
			selectedIndex = Math.min(selectedIndex + 1, suggestions.length - 1);
		} else if (e.key === 'ArrowUp') {
			e.preventDefault();
			selectedIndex = Math.max(selectedIndex - 1, 0);
		} else if (e.key === 'Enter' && selectedIndex >= 0) {
			e.preventDefault();
			selectContact(suggestions[selectedIndex]);
		} else if (e.key === 'Escape') {
			showDropdown = false;
			selectedIndex = -1;
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

<div class="relative flex flex-1 items-center" bind:this={containerEl}>
	<Label for={id} class="w-16 shrink-0 text-sm text-muted-foreground">{label}</Label>
	<input
		{id}
		type="text"
		{value}
		{placeholder}
		oninput={handleInput}
		onkeydown={handleKeydown}
		onfocus={() => { if (suggestions.length > 0) showDropdown = true; }}
		bind:this={inputEl}
		autocomplete="off"
		class="flex h-9 w-full bg-transparent py-1 text-sm transition-colors placeholder:text-muted-foreground focus-visible:outline-none disabled:cursor-not-allowed disabled:opacity-50"
	/>

	{#if showDropdown && suggestions.length > 0}
		<div class="absolute left-16 top-full z-50 mt-1 w-[calc(100%-4rem)] rounded-md border bg-popover p-1 shadow-md">
			{#each suggestions as contact, i}
				<button
					type="button"
					class="flex w-full cursor-pointer items-center gap-2 rounded-sm px-2 py-1.5 text-sm transition-colors hover:bg-accent {i === selectedIndex ? 'bg-accent' : ''}"
					onmousedown={(e) => { e.preventDefault(); selectContact(contact); }}
				>
					<span class="flex flex-col items-start">
						{#if contact.name}
							<span class="font-medium">{contact.name}</span>
							<span class="text-xs text-muted-foreground">{contact.email}</span>
						{:else}
							<span>{contact.email}</span>
						{/if}
					</span>
				</button>
			{/each}
		</div>
	{/if}
</div>
