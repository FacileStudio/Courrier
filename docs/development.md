# Courrier — Development

Getting a local Courrier running, the tasks that exist, and the quality gate that runs
before every push.

## Prerequisites

| Tool | Version | Why |
|---|---|---|
| Go | 1.24 | `apps/api/go.mod` declares `go 1.24.0`; `mise.toml` pins `go = "1.24"` |
| Bun | any recent | Client package manager and dev server |
| Docker | any recent | Postgres 16 for local development |
| mise | optional | Runs the tasks below; each is a one-line shell command you can run by hand |

You will also want a real IMAP/SMTP account to point Courrier at. There is no mock mail
server in this repo, and most of the interesting behavior only shows up against a real one.

## Setup

```sh
mise run hooks
mise run install
docker compose up db -d
```

`mise run hooks` points `core.hooksPath` at `.githooks`, which wires up the pre-push
quality gate. `mise run install` runs `bun install --frozen-lockfile` in `apps/client`.

The `db` service requires `POSTGRES_USER` and `POSTGRES_PASSWORD` in your `.env` — Compose
declares them with `:?`, so it refuses to start without them. Copy `.env.example` first.

## Running

The API takes configuration from the process environment — **there is no `.env` loading in
the Go code.** `apps/api/.env.example` documents what to set, but copying it to
`apps/api/.env` does nothing on its own. Export what you need:

```sh
cd apps/api
DATABASE_URL=postgres://postgres:postgres@localhost:5432/courrier?sslmode=disable \
  ENCRYPTION_KEY="$(openssl rand -base64 32)" \
  PORT=4000 \
  CORS_ALLOWED_ORIGINS=http://localhost:5173 \
  go run .
```

Use a **stable** `ENCRYPTION_KEY` once you have added a mail account: a fresh key on every
start makes the previously stored IMAP and SMTP passwords undecryptable.

The client, in another terminal:

```sh
cd apps/client
cp .env.example .env
bun run dev
```

The API listens on `:4000` and the client on `:5173`. There is no Vite proxy: the client
calls the API cross-origin at `VITE_API_BASE_URL`, which is why the API needs the dev
origin in its CORS list. CORS runs with `AllowCredentials: true` so the session cookie
survives the hop.

Migrations run at startup via `schemas.Migrate`, followed by the trigram index creation,
the avatar-URL backfill, the credential and OIDC token migrations, and — in a background
goroutine — `mail.BackfillThreads`.

## Tasks

| Task | Command | What it does |
|---|---|---|
| `mise run install` | `bun install --frozen-lockfile` in `apps/client` | Client dependencies |
| `mise run check` | `sh ./scripts/check.sh` | The full quality gate: Go, then the client |
| `mise run check-go` | `sh ./scripts/check.sh --go-only` | Go half only |
| `mise run format` | `sh ./scripts/check.sh --format` | `go fmt ./...`, rewriting files in place |
| `mise run hooks` | `git config core.hooksPath .githooks` | Enables the tracked hooks in this clone |

Client scripts, run from `apps/client`: `bun run dev`, `bun run build`, `bun run preview`,
`bun run check` (`svelte-check` against `tsconfig.json`).

## The quality gate

`scripts/check.sh` is the gate, and `.githooks/pre-push` does nothing but exec it. It
reports and never rewrites, except under `--format`:

1. `gofmt -l .` over `apps/api`, ignoring `vendor/`. Any listed file fails the gate.
2. `go vet ./...`
3. `go test ./...`
4. `bun run check` in `apps/client`, skipped with a warning if `bun` is not on `PATH`.

Two deliberate details worth knowing before you "fix" the script:

- **It is not invoked through mise.** `mise run` resolves every tool in the merged config
  before running any task body, so one broken tool in your global mise config would take
  the gate down with it. The hook calls `sh` directly.
- **It resolves the toolchain from `GOROOT`.** mise exports `GOROOT` for the pinned
  version but can leave an unrelated `go` earlier on `PATH`; mixing them produces
  `compile: version "X" does not match go tool version "Y"`.

Bypass once with `git push --no-verify`.

## Tests

```sh
cd apps/api
go test ./...
```

| File | What it covers |
|---|---|
| `routing_test.go` | Unknown `/api` paths do not reach the SPA, the public URL list is unchanged, health survives the `/api` subtree, `/api/files/` serves from storage, the SPA still answers non-API paths |
| `modules/mail/threading_test.go` | `parseRefIDs`, `collectThreadIDs`, `newThreadID` |
| `modules/mail/cid_route_test.go` | CID decoding and normalization, transfer-encoding decoding, content-type resolution, inline image resolution, content-disposition formatting |
| `modules/users/service_test.go` | Avatar files land in `STORAGE_DIR` and only managed avatars are deleted |
| `internal/crypto/crypto_test.go` | AES-GCM round trip, empty string, nonce uniqueness, wrong key, deterministic key derivation |

`routing_test.go` carries a literal `want` list of public API URLs and walks the router to
check each is still registered. **Add to it when you add a public endpoint** — that list is
the contract between the client and the router.

There is no client-side test setup.

## Dependencies

Go dependencies are vendored in `apps/api/vendor`, and the Docker build runs
`go build -mod=vendor`. After changing anything in `go.mod`:

```sh
cd apps/api
go mod tidy
go mod vendor
```

Forgetting `go mod vendor` produces a build that works locally and fails in the image.

## Conventions

- Every route goes under `/api`. The SPA catch-all is mounted last at `/*`, so a path
  registered outside `/api` would be shadowed by it — or worse, an unregistered `/api` path
  would fall through and return HTML with a 200.
- Each module keeps the same layout: `router.go` wires paths, `service.go` holds the logic
  and GORM calls, `types.go` holds the request and response structs. `auth`, `users`, and
  `settings` also have a `controller.go` and a `documentation.go`.
- Attachment bytes are never copied into Postgres. `attachments.part_id` records the IMAP
  part and downloads stream from the server on demand — keep it that way.
- Anything touching thread assignment goes through `assignThreadID`, which owns the
  transaction and the merge. Writing `emails.thread_id` directly will desynchronize
  `thread_links`.
- Svelte 5 runes are enforced through `dynamicCompileOptions` in `svelte.config.js`. UI
  primitives come from shadcn-svelte in `src/lib/components/ui/`.
