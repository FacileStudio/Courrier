# Courrier

Self-hosted email client for creative studios. Go API + SvelteKit frontend + PostgreSQL.

## Tech Stack

| Layer    | Stack                                                        |
| -------- | ------------------------------------------------------------ |
| API      | Go 1.24, Chi router, GORM, PostgreSQL 16                    |
| Client   | SvelteKit 5 (Svelte 5 runes), Tailwind CSS 4, shadcn-svelte |
| Build    | `go build` (vendored deps), Bun                             |
| Deploy   | Docker Compose (API via distroless, client via Nginx)        |
| Auth     | Session tokens + optional OIDC SSO, `SSO_ONLY` mode         |
| Email    | IMAP (go-imap/v2), SMTP (go-smtp), MIME (go-message)        |

## Project Structure

```
apps/
  api/                  Go backend
    main.go             Entrypoint: env, DB, migrations, router, graceful shutdown
    modules/            Domain modules (auth, accounts, users, settings)
    internal/           Shared infra (database, middleware, logger, env, errors, etc.)
    schemas/            GORM models and migrations (auto-run on startup)
    vendor/             Vendored Go dependencies
    Dockerfile          Multi-stage: golang:1.24-alpine -> distroless
  client/               SvelteKit frontend
    src/
      routes/           SvelteKit file-based routing
        (app)/          Authenticated layout group (mail, settings, profile)
        login/          Login page
      lib/
        backend.ts      API client (fetch wrapper)
        components/     App components + shadcn-svelte ui/ primitives
    Dockerfile          Multi-stage: oven/bun -> nginx:alpine (static adapter)
docker-compose.yml      Full stack: db + api + client
.env.example            Root-level env template (production)
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
bun run dev                 # Dev server on :5173
bun run build               # Static build to build/
bun run preview             # Preview production build
bun run check               # Svelte type checking
```

### Full Stack (Docker)

```sh
cp .env.example .env
docker compose up --build           # Everything
docker compose up db -d             # Just PostgreSQL for local dev
```

## Environment Variables

Core variables (see `.env.example` for full list):

- `DATABASE_URL` -- PostgreSQL connection string (default: local postgres)
- `DOMAINS` -- Comma-separated allowed CORS origins
- `PORT` -- API port (default `4000`)
- `LOG_LEVEL` -- `debug`, `info`, `warn`, `error`
- `STORAGE_DIR` -- File storage (default `./data`)
- `ENCRYPTION_KEY` -- 32+ char key for encrypting IMAP/SMTP credentials at rest
- `OIDC_*` -- OpenID Connect config (optional)
- `SSO_ONLY` -- Hide password auth when `true`
- `VITE_API_BASE_URL` -- Client-side API URL (build-time, default `http://localhost:4000`)

## Key Endpoints

- `GET /health`, `GET /ready` -- Health and readiness probes
- `GET /files/*` -- Static file serving (avatars)
- Auth: `/auth/register`, `/auth/login`, `/auth/config`, `/auth/oidc/*`
- Accounts: `/accounts` (CRUD for IMAP/SMTP mail accounts)
- Users: `/users/me`, `/users`
- Settings: `/settings/`

## Conventions

- **Go modules are vendored** -- run `go mod vendor` after changing dependencies.
- **Migrations run on startup** via `schemas.Migrate(db)`. No separate migration tool.
- **Client uses static adapter** -- output is plain HTML/CSS/JS served by Nginx in production.
- **shadcn-svelte** provides the UI component primitives in `src/lib/components/ui/`.
- **Svelte 5 runes** are enforced (`$state`, `$props`, `$derived`, `$effect`).
- **IMAP/SMTP credentials** are encrypted at rest using AES-GCM with `ENCRYPTION_KEY`.
- Avatar uploads are stored on disk at `STORAGE_DIR/avatars/` and served under `/files/`.
