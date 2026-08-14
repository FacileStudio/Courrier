<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { backend } from '$lib/backend';

	let ready = $state(false);
	let ssoOnly = $state(false);

	onMount(async () => {
		try {
			await backend.me();
			goto('/mail');
			return;
		} catch {}

		try {
			const cfg = await fetch(`${backend.baseUrl}/api/auth/config`, {
				credentials: 'include'
			}).then((r) => r.json());
			ssoOnly = cfg.sso_only ?? false;
		} catch {}

		ready = true;
	});
</script>

<svelte:head>
	<title>Courrier — Self-hosted email for creative studios</title>
	<meta
		name="description"
		content="A self-hosted email client for creative studios. Connect your IMAP accounts, read and send from a clean interface."
	/>
</svelte:head>

{#if ready}
	<div class="min-h-screen bg-background text-foreground">
		<header class="border-b border-border">
			<div class="mx-auto flex max-w-5xl items-center justify-between px-6 py-4">
				<div class="flex h-14 items-center gap-3">
					<iconify-icon
						icon="solar:letter-bold-duotone"
						width="28"
						height="28"
						class="size-7 text-foreground"
					></iconify-icon>
					<span class="text-2xl font-bold font-heading tracking-tight">Courrier</span>
				</div>
				<div class="flex items-center gap-2">
					<a
						href="/login"
						class="inline-flex h-9 items-center justify-center rounded-md px-4 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground"
					>
						Log in
					</a>
					<a
						href={ssoOnly ? '/login' : '/login?tab=register'}
						class="inline-flex h-9 items-center justify-center rounded-md bg-primary px-4 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
					>
						Sign in with Facile
					</a>
				</div>
			</div>
		</header>

		<main>
			<section class="mx-auto max-w-5xl px-6 py-24 text-center">
				<h1 class="text-4xl font-bold font-heading tracking-tight leading-tight">
					Your mail.<br />Your server.
				</h1>
				<p class="mx-auto mt-6 max-w-xl text-lg text-muted-foreground leading-relaxed">
					Courrier is a self-hosted email client. Connect your IMAP accounts, read and send from a
					clean interface — no cloud, no tracking, no compromise.
				</p>
				<div class="mt-10 flex justify-center gap-3">
					<a
						href={ssoOnly ? '/login' : '/login?tab=register'}
						class="inline-flex h-11 items-center justify-center rounded-md bg-primary px-6 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
					>
						Sign in with Facile
						<iconify-icon
							icon="solar:arrow-right-linear"
							width="16"
							height="16"
							class="ml-2 size-4"
						></iconify-icon>
					</a>
					<a
						href="/login"
						class="inline-flex h-11 items-center justify-center rounded-md border border-border bg-background px-6 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground"
					>
						Log in
					</a>
				</div>
			</section>

			<div class="mx-auto max-w-5xl"><div class="h-px bg-border"></div></div>

			<section class="mx-auto max-w-5xl px-6 py-20">
				<div class="grid gap-6 md:grid-cols-3">
					<div class="rounded-lg border border-border p-6">
						<div
							class="mb-3 flex size-10 items-center justify-center rounded-md border border-border"
						>
							<iconify-icon
								icon="solar:inbox-in-linear"
								width="20"
								height="20"
								class="size-5"
							></iconify-icon>
						</div>
						<h3 class="text-base font-semibold">Any IMAP account</h3>
						<p class="mt-1.5 text-sm text-muted-foreground">
							Point Courrier at your IMAP and SMTP servers. Credentials are encrypted at rest.
						</p>
					</div>

					<div class="rounded-lg border border-border p-6">
						<div
							class="mb-3 flex size-10 items-center justify-center rounded-md border border-border"
						>
							<iconify-icon
								icon="solar:users-group-rounded-linear"
								width="20"
								height="20"
								class="size-5"
							></iconify-icon>
						</div>
						<h3 class="text-base font-semibold">Shared spaces</h3>
						<p class="mt-1.5 text-sm text-muted-foreground">
							Group accounts into spaces and invite the people who need to read the same inbox.
						</p>
					</div>

					<div class="rounded-lg border border-border p-6">
						<div
							class="mb-3 flex size-10 items-center justify-center rounded-md border border-border"
						>
							<iconify-icon
								icon="solar:lock-linear"
								width="20"
								height="20"
								class="size-5"
							></iconify-icon>
						</div>
						<h3 class="text-base font-semibold">Self-hosted</h3>
						<p class="mt-1.5 text-sm text-muted-foreground">
							One binary and a Postgres database. Your mail stays on your server.
						</p>
					</div>
				</div>
			</section>

			<div class="mx-auto max-w-5xl"><div class="h-px bg-border"></div></div>

			<section class="mx-auto max-w-5xl px-6 py-20 text-center">
				<h2 class="text-3xl font-bold font-heading tracking-tight">
					{ssoOnly ? 'Ready to sign in?' : 'Ready to start?'}
				</h2>
				<p class="mt-4 text-muted-foreground">
					{ssoOnly
						? 'Use your Facile SSO to access Courrier.'
						: 'Free to use. Self-hosted. No credit card required.'}
				</p>
				<a
					href={ssoOnly ? '/login' : '/login?tab=register'}
					class="mt-8 inline-flex h-11 items-center justify-center rounded-md bg-primary px-6 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
				>
					Sign in with Facile
				</a>
			</section>
		</main>

		<footer class="border-t border-border text-center">
			<div class="mx-auto max-w-5xl px-6 py-6 text-sm text-muted-foreground">
				© {new Date().getFullYear()} Courrier by <a
					href="https://facile.studio"
					class="font-semibold underline hover:cursor-pointer">Facile.</a
				>
			</div>
		</footer>
	</div>
{/if}
