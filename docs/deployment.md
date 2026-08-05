# Courrier — Deployment

How the image is built, what Compose declares, and how Courrier is routed on la ruche.

## The image

`Dockerfile` is a three-stage build producing one binary plus the built client:

1. `oven/bun:1` installs `apps/client` dependencies with `--frozen-lockfile` and runs
   `bun run build`, producing static files in `/app/build`.
2. `golang:1.24-alpine` builds the API with
   `go build -mod=vendor -trimpath -ldflags="-s -w" -o /courrier .`. Dependencies come from
   the committed `apps/api/vendor` directory, so the build never reaches the network.
3. `gcr.io/distroless/static-debian12` receives the binary at `/courrier` and the client at
   `/client`.

The runtime stage sets `ENV CLIENT_DIR=/client` explicitly. The distroless base can carry
its own `WorkingDir` (`/home/nonroot` on the `:nonroot` variant), which would make a
relative `./client` resolve there and the SPA silently not be served at all.

`EXPOSE 4000`, `ENTRYPOINT ["/courrier"]`.

## Compose topology

Two services, and no published host ports — this file is written for Dokploy and Traefik.

| Service | Image | Notes |
|---|---|---|
| `db` | `postgres:16-alpine` | `pg_isready` healthcheck every 5s, volume `courrier_db_data` |
| `courrier` | built from `Dockerfile` | `expose: 4000`, joins `default` and the external `dokploy-network`, volume `courrier_api_data` at `/data` |

`courrier` waits on `db` being healthy through `depends_on: condition: service_healthy`.

### Required variables

Three have no default and Compose refuses to start without them:

```yaml
POSTGRES_USER: ${POSTGRES_USER:?POSTGRES_USER is required}
POSTGRES_PASSWORD: ${POSTGRES_PASSWORD:?POSTGRES_PASSWORD is required}
ENCRYPTION_KEY: ${ENCRYPTION_KEY:?ENCRYPTION_KEY is required}
```

`DATABASE_URL` is assembled by Compose from the Postgres credentials, so it is not
something you set yourself in the Docker path. `PORT` and `CLIENT_DIR` are hard-coded to
`"4000"` and `/client` in the `environment:` block and are not overridable from `.env`.

`ENCRYPTION_KEY` is the one to be careful with. Rotating it after accounts exist makes
every stored IMAP and SMTP password undecryptable, and each account has to be re-entered.
Treat it as permanent per deployment.

### Healthcheck

```yaml
test: ["CMD", "/courrier", "healthcheck"]
```

The image is distroless: no shell, no `curl`, no `wget`. `tronc/healthcheck` handles this
by making the binary probe itself — `main` checks `os.Args` first and, when the argument is
`healthcheck`, requests `http://127.0.0.1:$PORT/health` with a three-second timeout and
exits 0 or 1. It targets `127.0.0.1` rather than `localhost` on purpose: in these
containers `localhost` resolves to `::1` first while the server binds `0.0.0.0`, so a
`localhost` probe fails against a perfectly healthy process.

This is also why `PORT` must stay pinned at `4000`: the self-probe reads `PORT` from the
environment, and the Traefik service port label is a separate literal that has to agree
with it.

### Traefik labels

One hostname, one service, two routers — plain HTTP redirecting to HTTPS:

```yaml
traefik.enable: "true"
traefik.docker.network: dokploy-network
traefik.http.routers.courrier-web.rule: "Host(`courrier.facile.studio`)"
traefik.http.routers.courrier-web.entrypoints: web
traefik.http.routers.courrier-web.middlewares: redirect-to-https@file
traefik.http.routers.courrier-web.service: courrier-svc
traefik.http.routers.courrier-secure.rule: "Host(`courrier.facile.studio`)"
traefik.http.routers.courrier-secure.entrypoints: websecure
traefik.http.routers.courrier-secure.tls.certresolver: letsencrypt
traefik.http.routers.courrier-secure.service: courrier-svc
traefik.http.services.courrier-svc.loadbalancer.server.port: "4000"
```

This is the suite's one-container / one-router / one-hostname rule. Do not add a second
router for `/api` — the same binary answers both the API and the SPA, and splitting them
would break the SPA fallback that every non-API path relies on.

`X-Forwarded-Proto: https` matters: it is what makes the session cookie `Secure`. A proxy
that drops it downgrades every session cookie to non-secure.

## Egress

Unlike most of the suite, Courrier makes **outbound connections to arbitrary hosts** — the
IMAP and SMTP servers of whatever accounts users add. A network policy that restricts
egress will break mail sync and sending with connection timeouts rather than anything
resembling a useful error. Ports 993 (IMAP over TLS), 465 (SMTP implicit TLS), and 587
(SMTP with `STARTTLS`) are the usual ones.

## Deploying to la ruche

Courrier runs on la ruche behind the Dokploy panel at `gare.facile.studio`. Prefer the
`dokploy` CLI over SSH plus `docker`:

```sh
dokploy compose --help
```

Environment values are set in the Dokploy project, not committed. `.env.example` is the
template — at minimum `ENCRYPTION_KEY`, `POSTGRES_USER`, and `POSTGRES_PASSWORD`.

## Migrations

There is no migration step to run. On every boot, in order:

1. `schemas.Migrate` — GORM `AutoMigrate` over the twelve models. A failure aborts startup
   and the container exits 1.
2. `ensureSearchIndexes` — `CREATE EXTENSION IF NOT EXISTS pg_trgm` plus seven
   `CREATE INDEX IF NOT EXISTS` statements. It returns early and silently if the extension
   cannot be created, so a database user without the privilege gets slow search rather than
   a failed deploy.
3. `backfillAvatarURLs` — rewrites any `users.avatar_url` still on the legacy `/files/`
   prefix to `/api/files/`. Idempotent.
4. `crypto.MigrateAccountPasswords` and `crypto.MigrateOIDCTokens` — encrypt any plaintext
   credential found. Best-effort: a failure logs a warning and startup continues.
5. `mail.BackfillThreads`, in a background goroutine — assigns thread ids to emails stored
   before threading existed. Idempotent and batched, so it does not delay serving.

The first boot after a large mailbox already exists will spend a while on step 5 in the
background. The app is serving throughout.

## Verifying a deploy

`/health` answers as soon as the process is serving and touches nothing, so a green
`/health` proves the binary started and nothing more. `/ready` pings Postgres and is the
one that tells you the database is reachable. Neither says anything about the SPA — for
that, request `/` and confirm you get the client's `index.html` rather than a 404, which is
what a missing or misplaced `CLIENT_DIR` looks like. Neither says anything about mail
either; the honest check there is `POST /api/mail/test-connection`, which dials a real
IMAP and SMTP pair and reports the failure.

## Persistent state

| Volume | Mounted at | Holds |
|---|---|---|
| `courrier_db_data` | `/var/lib/postgresql/data` | The database, including cached message bodies |
| `courrier_api_data` | `/data` | Uploaded and OIDC-fetched avatars under `/data/avatars` |

Message bodies in Postgres are a cache, not the source of truth — the IMAP server is. A
lost database costs the local index and the thread links, both of which rebuild on the
next sync. Attachment bytes were never stored locally at all.
