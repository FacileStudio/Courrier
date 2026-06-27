<script lang="ts">
	const sections = [
		{
			href: '/profile',
			icon: 'solar:user-linear',
			title: 'Profile',
			description: 'Your name, email, avatar, and signatures'
		},
		{
			href: '/accounts',
			icon: 'solar:mailbox-linear',
			title: 'Accounts',
			description: 'Manage your IMAP and SMTP mail accounts'
		},
		{
			href: '/templates',
			icon: 'solar:documents-linear',
			title: 'Templates',
			description: 'Reusable email templates to speed up composing'
		},
		{
			href: '/spaces',
			icon: 'solar:users-group-rounded-linear',
			title: 'Spaces',
			description: 'Collaborate with your team in shared spaces'
		}
	];
</script>

<svelte:head>
	<title>Settings — Courrier</title>
</svelte:head>

<div class="flex h-full flex-col">
	<div class="border-b border-border px-4 py-4 md:px-8 md:py-5">
		<h1 class="text-lg font-semibold">Settings</h1>
		<p class="mt-0.5 text-sm text-muted-foreground">
			Manage your profile, accounts, templates, and spaces.
		</p>
	</div>

	<div class="flex-1 overflow-auto p-4 md:p-8">
		<div class="grid max-w-2xl gap-3">
			{#each sections as section (section.href)}
				<a
					href={section.href}
					class="flex items-center gap-4 rounded-lg border border-border p-4 transition-colors hover:bg-muted/50"
				>
					<div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-primary/10">
						<iconify-icon icon={section.icon} width="20" class="text-primary"></iconify-icon>
					</div>
					<div class="min-w-0 flex-1">
						<p class="truncate text-sm font-medium">{section.title}</p>
						<p class="truncate text-xs text-muted-foreground">{section.description}</p>
					</div>
					<iconify-icon
						icon="solar:alt-arrow-right-linear"
						width="16"
						class="shrink-0 text-muted-foreground"
					></iconify-icon>
				</a>
			{/each}
		</div>
	</div>
</div>
