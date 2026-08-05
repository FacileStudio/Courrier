# Courrier

Self-hosted email client for creative studios. Single Go binary serves both the API and the static SvelteKit frontend. PostgreSQL for storage.

## Tech Stack

| Layer    | Stack                                                        |
| -------- | ------------------------------------------------------------ |
| API      | Go 1.24, Chi router, GORM, PostgreSQL 16                    |
| Client   | SvelteKit 5 (Svelte 5 runes), Tailwind CSS 4, shadcn-svelte |
| Build    | `go build` (vendored deps), Bun                             |
| Deploy   | Docker Compose (single container: distroless Go binary)      |
| Auth     | Session tokens + optional OIDC SSO, `SSO_ONLY` mode         |
| Email    | IMAP (go-imap/v2), SMTP (go-smtp), MIME (go-message)        |

## Project Structure

```
Dockerfile                Unified multi-stage build (bun build + go build -> distroless)
docker-compose.yml        Two services: db + courrier
.env.example              Root-level env template (production)
apps/
  api/                    Go backend
    main.go               Entrypoint: env, DB, migrations, router, static file serving, graceful shutdown
    modules/              Domain modules (auth, accounts, users, settings, spaces)
    internal/             Shared infra (database, middleware, logger, env, errors, etc.)
    schemas/              GORM models and migrations (auto-run on startup)
    vendor/               Vendored Go dependencies
  client/                 SvelteKit frontend
    src/
      routes/             SvelteKit file-based routing
        (app)/            Authenticated layout group (mail, settings, profile)
        login/            Login page
      lib/
        backend.ts        API client (fetch wrapper)
        components/       App components + shadcn-svelte ui/ primitives
```

## Commands

### API (`apps/api/`)

```sh
cp .env.example .env
go run .                    # Dev server on :4000
go test ./...               # Run tests
go build -o bin/api .       # Production binary
```

### Client (`apps/client/`)

```sh
bun install                 # Install dependencies
bun run dev                 # Dev server on :5173 (needs VITE_API_BASE_URL in .env)
bun run build               # Static build to build/
bun run preview             # Preview production build
bun run check               # Svelte type checking
```

### Full Stack (Docker)

```sh
cp .env.example .env
docker compose up --build           # Everything on :4000
docker compose up db -d             # Just PostgreSQL for local dev
```

## Environment Variables

Core variables (see `.env.example` for full list):

- `DATABASE_URL` -- PostgreSQL connection string (**required**, no default)
- `PORT` -- API port (default `8080`; the compose service and `.env.example` pin it to `4000`)
- `LOG_LEVEL` -- `debug`, `info`, `warn`, `error`
- `STORAGE_DIR` -- File storage (default `./data`)
- `APP_ENV` -- `development`, `staging`, `production` (default `development`)
- `ENCRYPTION_KEY` -- **required**; 32+ char key for encrypting IMAP/SMTP credentials at rest. The API exits 1 without it
- `CLIENT_DIR` -- Path to SvelteKit build output (default `./client`); the Go binary serves these files as the frontend
- `CORS_ALLOWED_ORIGINS` -- Comma-separated allowed CORS origins (optional; only needed if deploying the client separately from the API). Legacy names `DOMAINS`, `ALLOWED_ORIGINS`, `DOMAIN`, `CORS_ORIGINS`, `TRUSTED_ORIGINS`, `CLIENT_ORIGIN` are still read as fallbacks
- `OIDC_*` -- OpenID Connect config (optional)
- `SSO_ONLY` -- Hide password auth when `true`

Client dev only (in `apps/client/.env`):

- `VITE_API_BASE_URL` -- API URL for local dev (set to `http://localhost:4000`); not needed in production since everything is same-origin

## Key Endpoints

- `GET /health`, `GET /ready` -- Health and readiness probes (outside `/api`)
- `GET /api/files/*` -- Static file serving (avatars)
- Auth: `/api/auth/register`, `/api/auth/login`, `/api/auth/config`, `/api/auth/oidc/*`
- Accounts: `/api/accounts` (CRUD for IMAP/SMTP mail accounts)
- Mail: `/api/accounts/{id}/mail/sync`, `/api/accounts/{id}/mail/folders`, `/api/accounts/{id}/mail/folders/{type}/emails`, `/api/accounts/{id}/mail/emails/{id}`, `/api/accounts/{id}/mail/send`
- Connection test: `POST /api/mail/test-connection`
- Users: `/api/users/me`, `/api/users`
- Settings: `/api/settings/`
- Spaces: `/api/spaces` (CRUD), `/api/spaces/{id}/members` (member management), `/api/spaces/{id}/leave`

Every API route is registered inside a single `router.Route("/api", ...)` in `main.go`; the
modules declare relative prefixes (`/auth`, `/accounts`, ...). That subtree is what makes an
unknown `/api/*` path return 404 instead of falling through to the SPA catch-all.

## Conventions

- **Single binary architecture** -- In production, the Go binary serves the SvelteKit static build from `CLIENT_DIR`. No Nginx, no CORS, no separate client container.
- **Go modules are vendored** -- run `go mod vendor` after changing dependencies.
- **Migrations run on startup** via `schemas.Migrate(db)`. No separate migration tool.
- **Client uses static adapter** -- output is plain HTML/CSS/JS served by the Go binary in production.
- **shadcn-svelte** provides the UI component primitives in `src/lib/components/ui/`.
- **Svelte 5 runes** are enforced (`$state`, `$props`, `$derived`, `$effect`).
- **IMAP/SMTP credentials** are encrypted at rest using AES-GCM with `ENCRYPTION_KEY`.
- Avatar uploads are stored on disk at `STORAGE_DIR/avatars/` and served under `/files/`.
- **Local dev** still uses two processes: `go run .` (API on :4000) + `bun run dev` (client on :5173 with `VITE_API_BASE_URL=http://localhost:4000`).
- **One Dockerfile** at the repo root builds the whole app (bun build + go build -> distroless). There are no per-app Dockerfiles; local dev runs the two processes directly.
