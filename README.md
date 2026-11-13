# Courrier

Self-hosted email client for creative studios and freelancers. Connect your existing
IMAP/SMTP accounts and read, thread, search, and send from a web interface.

One Go binary serves both the JSON API and the built SvelteKit client, so a deployment is
a single container behind a single Traefik router.

Live at [courrier.facile.studio](https://courrier.facile.studio).

## What it does

- Multiple IMAP/SMTP accounts per user, with credentials encrypted at rest
- Conversation threading that survives out-of-order delivery and truncated `References`
- Full-text search over subject, sender, and body, backed by Postgres trigram indexes
- Compose, reply, and forward with attachments, plus drafts appended back to the IMAP server
- Reusable email templates, per user or per space
- Bulk delete, archive, and read-state changes over up to 200 messages at a time
- Spaces that group users and scope accounts and templates. With Facile SSO on, Registre owns the spaces and their members, and Courrier keeps them in step
- Sign-in through your Facile account (Registre); API tokens for the CLI, scripts and mail clients

## Stack

| Layer | Tech |
|---|---|
| API | Go 1.24, Chi v5, GORM, PostgreSQL 16, [tronc](https://github.com/FacileStudio/tronc) v0.6.0 |
| Mail | go-imap v2 (beta), go-message for MIME, `net/smtp` for delivery |
| Client | SvelteKit 5 (runes), Tailwind CSS 4, shadcn-svelte, bits-ui |
| Deploy | Docker Compose, one distroless container behind Traefik |

## Quick start

`docker-compose.yml` publishes no host port and expects the external `dokploy-network`
with Traefik in front of it. It is how Courrier deploys, not how you browse it locally.

```sh
cp .env.example .env
docker compose up -d --build
```

`ENCRYPTION_KEY`, `POSTGRES_USER`, and `POSTGRES_PASSWORD` have no defaults — Compose
refuses to start without them, and the API exits 1 without `ENCRYPTION_KEY`.

### Local development

Start Postgres, then the API and the client in separate terminals.

```sh
mise run install
docker compose up db -d
```

```sh
cd apps/api
DATABASE_URL=postgres://postgres:postgres@localhost:5432/courrier?sslmode=disable \
  ENCRYPTION_KEY="$(openssl rand -base64 32)" PORT=4000 \
  CORS_ALLOWED_ORIGINS=http://localhost:5173 go run .
```

```sh
cd apps/client
cp .env.example .env
bun run dev
```

The desktop app runs against that dev server with `cd apps/desktop && mise run dev`.
On Linux it needs the WebKitGTK 4.1 development packages.

The API does not read a `.env` file — it takes everything from the process environment.
The client runs on <http://localhost:5173> and proxies `/api` to the API on :4000.

## Configuration

| Variable | What it does |
|---|---|
| `DATABASE_URL` | Postgres connection string. Required — the API exits 1 without it |
| `ENCRYPTION_KEY` | Required. Derives the AES key for IMAP/SMTP credentials and the resource-token secret |
| `PORT` | HTTP listen port. `tronc` defaults to `8080`; every deployment here pins `4000` |
| `CORS_ALLOWED_ORIGINS` | Comma-separated browser origins, only needed for a split deploy |
| `STORAGE_DIR` | Root for uploaded avatars, served back under `/api/files/` |

Full reference: [docs/configuration.md](docs/configuration.md).

## Structure

```
apps/
  api/       Go backend — modules/ (auth, accounts, mail, users, settings, spaces),
             schemas/ (GORM models), internal/ (crypto, resourcetoken, middleware)
  client/    SvelteKit 5 SPA, built into the API image and served by it; android/ wraps it
             as the Capacitor app
  desktop/   Tauri 2 shell around the same client build (mise run dev / mise run build)
scripts/     check.sh, the quality gate the pre-push hook runs
docs/        Architecture, configuration, development, deployment, API
```

## Documentation

| Doc | What's in it |
|---|---|
| [Architecture](docs/architecture.md) | Request flow, threading, data model, how the pieces fit |
| [Configuration](docs/configuration.md) | Every environment variable and default |
| [Development](docs/development.md) | Local setup, tests, the quality gate |
| [Deployment](docs/deployment.md) | Docker Compose, Dokploy, Traefik routing |
| [API](docs/api.md) | HTTP endpoints and payloads |

---

Part of the [Facile Suite](https://facile.studio) — self-hosted tools for creative studios
and freelancers. One login, zero cloud dependency.

## License

Courrier is licensed under the [GNU Affero General Public License v3.0](LICENSE) (`AGPL-3.0-only`). A commercial licence is available from Facile Studio for use that cannot meet the AGPL's terms.
