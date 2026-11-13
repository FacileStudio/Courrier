# Courrier

Self-hosted email client for creative studios. Single Go binary serves both the API and the static SvelteKit frontend. PostgreSQL for storage.

## Tech Stack

| Layer    | Stack                                                        |
| -------- | ------------------------------------------------------------ |
| API      | Go 1.24, Chi router, GORM, PostgreSQL 16                    |
| Client   | SvelteKit 5 (Svelte 5 runes), Tailwind CSS 4, shadcn-svelte |
| Build    | `go build` (vendored deps), Bun                             |
| Deploy   | Docker Compose (single container: distroless Go binary)      |
| Auth     | Sessions and API tokens over OIDC sign-in through Registre |
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
    migrations/           Ordered goose SQL, embedded and applied at boot by tronc/migrate
    schemas/              GORM models (query layer only) and the idempotent boot data steps
    vendor/               Vendored Go dependencies
  desktop/                Tauri 2 desktop shell (own Cargo workspace) around ../client/build
  client/                 SvelteKit frontend (android/ is the Capacitor app)
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

Client dev only (in `apps/client/.env`):

- `VITE_API_BASE_URL` -- API URL for local dev (set to `http://localhost:4000`); not needed in production since everything is same-origin

## Key Endpoints

- `GET /health`, `GET /ready` -- Health and readiness probes (outside `/api`)
- `GET /api/files/*` -- Static file serving (avatars)
- Auth: `/api/auth/config`, `/api/auth/oidc/*` (no register or login route: Registre owns sign-in)
- Accounts: `/api/accounts` (CRUD for IMAP/SMTP mail accounts)
- Mail: `/api/accounts/{id}/mail/sync`, `/api/accounts/{id}/mail/folders`, `/api/accounts/{id}/mail/folders/{type}/emails`, `/api/accounts/{id}/mail/emails/{id}`, `/api/accounts/{id}/mail/send`, `/api/accounts/{id}/mail/drafts[/{id}]`, `/api/accounts/{id}/mail/emails/move`, `/api/accounts/{id}/mail/folders[/{id}]` (create, rename, delete)
- Push: `/api/push`, `/api/push/test`, `/api/push/devices`, `/api/push/stream` (desktop SSE)
- Connection test: `POST /api/mail/test-connection`
- Users: `/api/users/me`, `/api/users`
- Settings: `/api/settings/`
- Spaces: `/api/spaces` (CRUD), `/api/spaces/{id}/members` (member management), `/api/spaces/{id}/leave`
- Registre nudge: `POST /api/registre/sync` (anonymous, answers 202, mounted only with OIDC on)
- MCP: `POST /mcp` and `/.well-known/oauth-protected-resource/mcp` at the root, plus the confirm page API `/api/mcp/pending` (mounted only with OIDC on and `MCP_RESOURCE` set; see docs/mcp.md)

Every API route is registered inside a single `router.Route("/api", ...)` in `main.go`; the
modules declare relative prefixes (`/auth`, `/accounts`, ...). That subtree is what makes an
unknown `/api/*` path return 404 instead of falling through to the SPA catch-all.

## Conventions

- **Single binary architecture** -- In production, the Go binary serves the SvelteKit static build from `CLIENT_DIR`. No Nginx, no CORS, no separate client container.
- **Go modules are vendored** -- run `go mod vendor` after changing dependencies.
- **Migrations are ordered goose SQL** in `apps/api/migrations/`, embedded in the binary and applied at boot by
  `tronc/migrate` behind an advisory lock, then `schemas.Reconcile` runs the idempotent data steps. GORM is a
  query layer: `AutoMigrate` never runs and `scripts/check.sh` fails on one outside a test. A schema change is a
  new numbered file with `-- +goose Up` and `-- +goose Down`, plus the struct in `schemas/` by hand.
  `00001_baseline.sql` is production's schema as AutoMigrate left it; a database from before v0.17.0 is stamped
  at version 1 on its first boot. `docker run <image> migrate status` (also `up`, `version`, `baseline <n>`)
  answers without a shell.
- **Client uses static adapter** -- output is plain HTML/CSS/JS served by the Go binary in production.
- **shadcn-svelte** provides the UI component primitives in `src/lib/components/ui/`.
- **Svelte 5 runes** are enforced (`$state`, `$props`, `$derived`, `$effect`).
- **IMAP/SMTP credentials** are encrypted at rest using AES-GCM with `ENCRYPTION_KEY`.
- **Avatars have two sources and one derived value.** `oidc_picture_url` is written by the
  OIDC sync but only when the `picture` claim starts with `https://` (Authentik returns a
  `data:image/svg+xml` placeholder of the user's initials otherwise, so a non-empty claim is
  not a photo). `avatar_upload_path` is written by upload, stored at `STORAGE_DIR/avatars/`
  and served under `/api/files/`. What the client renders is `schemas.User.Avatar()` — SSO
  photo, else the upload, else nothing — never a stored third copy. Uploading is *rejected*
  while an SSO photo exists, and clearing removes only the upload. `avatar_url` and
  `avatar_source` remain as unread columns until a later release drops them.
- **Registre owns shared spaces when OIDC is on.** `spaces.Directory` imports every local space
  once (`POST <issuer>/api/directory/import`, retried with backoff from 30 s to 5 min while
  Registre answers 404), and only then starts `porte/directory`'s Syncer: a pull at boot, every
  5 minutes and 2 s after a nudge. `Apply` links by `spaces.registre_space_id`, else by an
  unlinked space of the same `lower(name)`, else creates; resolves members by porte identity
  (issuer, sub), else email, else pre-provisions a user plus identity row. It only removes
  members who hold an identity under the issuer, so password-only accounts keep their seat.
  Creating a space, and editing a linked one, answer 409 naming the Registre page; the client
  shows "Gérer dans Registre" from `manage_url`. Mail accounts in a space stay Courrier's.
- **Plans come from Registre through porte/forfait.** `internal/plans` caps mail accounts per
  scope (`courrier.mailboxes`) and saves the store in the single-row `forfait_state` table after
  every pull. The store ignores an older entitlements revision, so after Registre or Ardoise is
  restored to an older revision, run `DELETE FROM forfait_state;` and restart Courrier: it fails
  open until the boot pull applies the restored values. Client space ids go through
  `internal/spaceid` before they touch a text `space_id` column.
- **Database-backed tests need Postgres.** `TEST_DATABASE_URL` points at a scratch database
  (`github.com/FacileStudio/tronc/testdb` gives each test binary its own schema); without it
  those tests skip loudly. SQLite is not an option here — the gate rejects the driver.
- **Local dev** still uses two processes: `go run .` (API on :4000) + `bun run dev` (client on :5173 with `VITE_API_BASE_URL=http://localhost:4000`).
- **Installed apps** -- `isNative()` is the Capacitor Android app, `isDesktop()` the Tauri app,
  `isApp()` either. Both sign in with porte's `flow=mobile` (system browser, then the
  `studio.facile.courrier://auth` deep link) and send a bearer token instead of the cookie.
  The token lives in the platform's secure store: `capacitor-secure-storage-plugin` (Android
  Keystore) and, on the desktop, the OS keychain through the `keyring` crate (`apps/desktop/src/token.rs`,
  commands `token_get`, `token_set`, `token_clear`; Linux builds need `libdbus-1-dev`). `getToken` moves a
  token an older release kept in `@capacitor/preferences` on first launch and deletes that copy only once the
  secure store holds it.
  The desktop version is read from `apps/client/package.json`; the updater key lives in tiroir
  as `COURRIER_TAURI_SIGNING_PRIVATE_KEY(_PASSWORD)` and in GitHub secrets
  `TAURI_SIGNING_PRIVATE_KEY(_PASSWORD)`. Never rotate it, or installed apps stop updating.
- **Desktop window and notifications** -- The window has no system frame on Windows and Linux
  (`TitleBar.svelte` draws the controls in a 36px drag strip); macOS keeps its traffic lights
  through `tauri.macos.conf.json`, which repeats the whole window object because arrays replace in
  the merge. `html[data-desktop]` in `app.css` shifts `h-dvh`, `min-h-dvh` and `fixed inset-0`
  down by `--titlebar-h`. The webview has no push service, so the desktop app reads
  `GET /api/push/stream` (SSE) and the inbox watcher also polls users with an open stream. The
  stream fan-out is in memory: it assumes one API replica.
- **One Dockerfile** at the repo root builds the whole app (bun build + go build -> distroless). There are no per-app Dockerfiles; local dev runs the two processes directly.
