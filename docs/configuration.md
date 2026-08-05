# Courrier — Configuration

Every environment variable the API actually reads, taken from `apps/api/internal/env` and
`tronc/env`.

**Courrier does not load a `.env` file.** There is no `godotenv` call in `env.Load` —
configuration comes from the process environment only. In Docker that is what
`docker-compose.yml` provides through its `environment:` block. Running `go run .` by hand
means exporting the variables yourself; `apps/api/.env.example` documents what to export
but copying it to `apps/api/.env` will not do anything on its own.

## Core

These come from `tronc/env`'s `LoadCore`, shared by every Go app in the suite.

| Variable | Required | Default | What it does |
|---|---|---|---|
| `DATABASE_URL` | yes | — | Postgres connection string |
| `PORT` | no | `8080` | HTTP listen port. Must be a valid TCP port or startup fails |
| `APP_ENV` | no | `development` | `development`, `staging`, or `production`. Never gates security behavior |
| `LOG_LEVEL` | no | `info` | `debug`, `info`, `warn`, or `error` |
| `CORS_ALLOWED_ORIGINS` | no | — | Comma-separated browser origins, scheme included |
| `JOURNAL_URL` | no | — | Journal ingest URL for log shipping |
| `JOURNAL_TOKEN` | no | — | Journal per-app key. Logs ship only when both this and `JOURNAL_URL` are set |

**Two traps.**

`DATABASE_URL` is required — a missing or blank value makes startup return an error and
the process exit 1. There is no fallback DSN.

**`PORT` defaults to `8080` but Courrier is pinned to `4000`** in `docker-compose.yml`
(`PORT: "4000"`, not overridable from `.env`), both `.env.example` files, the Dockerfile's
`EXPOSE`, and the Traefik service port. The pin must stay: the Traefik
`loadbalancer.server.port` label and the healthcheck's self-probe both assume it. Change
one and you must change all of them.

CORS is only needed when the client is served from a different origin than the API. In the
single-container deployment the SPA is served by the same binary, so the list can be empty.

### CORS name fallbacks

`troncenv.CORSOrigins` reads the first of these that is set, in order:
`CORS_ALLOWED_ORIGINS`, `ALLOWED_ORIGINS`, `DOMAINS`, `DOMAIN`, `CORS_ORIGINS`,
`TRUSTED_ORIGINS`, `CLIENT_ORIGIN`.

## Encryption

| Variable | Required | Default | What it does |
|---|---|---|---|
| `ENCRYPTION_KEY` | **yes** | — | The single secret behind credential encryption and resource tokens |

This is the one variable Courrier requires that most suite apps do not. `env.Load` calls
`troncenv.Required("ENCRYPTION_KEY")`, so a missing or blank value exits 1 before anything
else happens. `docker-compose.yml` enforces it too, with
`${ENCRYPTION_KEY:?ENCRYPTION_KEY is required}`.

It feeds two separate derivations:

- `crypto.DeriveKey` — `sha256(key)`, used as the AES-256-GCM key for IMAP and SMTP
  passwords on `accounts` and for OIDC tokens on `users`.
- `resourcetoken.DeriveSecret` — `sha256("courrier-resource-token:" + key)`, the HMAC
  secret for the 5-minute inline-image tokens.

Any string works — the SHA-256 produces the 32 bytes regardless, so the "32+ characters"
note in `.env.example` is a strength suggestion, not a length requirement. Generate one
with `openssl rand -base64 32`.

**Changing it after accounts exist makes their stored credentials undecryptable**, and
every affected account has to be re-entered. There is no re-key path.

## Storage and client

| Variable | Required | Default | What it does |
|---|---|---|---|
| `STORAGE_DIR` | no | `./data` | Root for uploaded files. `STORAGE_DIR/avatars` is created at startup and served under `/api/files/` |
| `CLIENT_DIR` | no | `./client` | Directory holding the built SPA. The image sets `/client` explicitly |

`CLIENT_DIR` is read by `tronc/spa`. If the directory has no `index.html`, the SPA
catch-all is not mounted and the binary serves the API alone.

Note the `/api/files/` prefix: avatars moved under `/api` when the API was consolidated
there, and `schemas.backfillAvatarURLs` rewrites any row still holding the legacy `/files/`
prefix on every boot.

## Auth and OIDC

| Variable | Required | Default | What it does |
|---|---|---|---|
| `SSO_ONLY` | no | `false` | When true, the register and login routes are not registered at all |
| `OIDC_ISSUER` | no | — | Discovery URL, for example `https://porte.facile.studio/application/o/courrier/` |
| `OIDC_CLIENT_ID` | with issuer | — | Client id from the provider |
| `OIDC_CLIENT_SECRET` | with issuer | — | Client secret from the provider |
| `OIDC_REDIRECT_URL` | with issuer | — | Must point at `/api/auth/oidc/callback` on the public hostname |
| `OIDC_SUCCESS_URL` | no | first `CORS_ALLOWED_ORIGINS` entry | Where the callback redirects after a successful login |

Setting `OIDC_ISSUER` without all three of the client id, secret, and redirect URL fails
startup with an explicit error.

## Compose-only

`POSTGRES_USER` and `POSTGRES_PASSWORD` are **required** by `docker-compose.yml` —
declared as `${POSTGRES_USER:?...}`, so Compose refuses to run without them.
`POSTGRES_DB` defaults to `courrier`. The API never reads any of the three; Compose
assembles them into the `DATABASE_URL` it passes to the container.

## Client

`apps/client/src/lib/backend.ts` reads `VITE_API_BASE_URL` at build time and falls back to
the empty string — meaning same-origin. The production image never sets it, because the Go
binary serves the SPA. In development, `apps/client/.env.example` sets it to
`http://localhost:4000`, and the API then needs `http://localhost:5173` in its CORS origin
list. CORS is configured with `AllowCredentials: true` so the session cookie survives the
cross-origin hop.
