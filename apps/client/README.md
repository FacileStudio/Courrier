# Client

SvelteKit frontend for Courrier.

## Responsibilities

- Login, registration, and OIDC entry flow
- Mail inbox, compose, folder navigation
- Account settings and profile management

## Run locally

```sh
bun install
bun run dev
```

Default dev URL: `http://localhost:5173`

## Scripts

```sh
bun run dev
bun run build
bun run preview
bun run check
```

## Configuration

- `VITE_API_BASE_URL`: points the installed apps (desktop, Android) at another API in development. The web build ignores it and calls `/api` on its own origin; `bun run dev` proxies `/api` to `127.0.0.1:4000`.

The production image is built from the repository's root `Dockerfile`; the Go binary serves this static build.
