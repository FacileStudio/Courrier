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

- `VITE_API_BASE_URL`: API base URL used by the frontend, default `http://localhost:4000`

The production Docker build injects `VITE_API_BASE_URL` at build time and serves the static output with Nginx.
