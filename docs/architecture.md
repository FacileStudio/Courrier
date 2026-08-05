# Courrier — Architecture

How a request reaches the database, how mail gets in and out, how conversation threading
works, and what the tables look like.

## Runtime topology

```
Internet ──▶ Traefik ──▶ Go binary (:4000) ──┬──▶ /health /ready     liveness, readiness
                                              ├──▶ /api/*            six modules
                                              ├──▶ /api/files/*      STORAGE_DIR
                                              └──▶ /*                SPA catch-all, gzipped
                                                              │
                                                        Postgres 16
                                                        + pg_trgm

  IMAP server ◀──── fetch, flags, APPEND ────  mail module  ──── SMTP send ────▶ SMTP server
```

One process, one container, one Traefik router. The SvelteKit client is built to static
files at image build time and served by the same Go binary through `tronc`'s `spa`
package.

## The `/api` subtree

Every route the application serves — including the avatar file server — is registered
inside a single `router.Route("/api", ...)` block in `apps/api/main.go`. The six modules
declare relative prefixes: `/auth`, `/accounts`, the mail routes, `/users`, `/settings`,
`/spaces`.

That subtree is load-bearing. Because the SPA catch-all is mounted last at `/*`, anything
not claimed under `/api` falls through and gets `index.html` with a 200 — which surfaces in
the browser as an unrelated JSON parse error rather than the 404 it should be. Only the
health probes live outside it, and only because `tronc/health.Mount` registers them at both
the root and under `/api` so the same probe answers regardless of what the edge forwards.

`routing_test.go` guards all of this: `TestUnknownAPIRoutesDoNotFallThroughToTheSPA`,
`TestPublicAPIURLsAreUnchanged` (a literal list of public URLs walked out of the router),
`TestHealthRoutesSurviveTheAPISubtree`, `TestFilesRouteStillServesFromStorage`, and
`TestSPAFallbackStillServesNonAPIPaths`.

**Avatars are served from `/api/files/`, not `/files/`.** They moved when the API was
consolidated under `/api`, and `schemas.backfillAvatarURLs` rewrites any row still holding
the old `/files/` prefix on every boot.

## Components

| Piece | Where | What it does |
|---|---|---|
| Entrypoint | `apps/api/main.go` | Config, DB, migrations, thread backfill, credential migrations, router, graceful shutdown |
| Router | `tronc/httpx.NewRouter` | Chi router with the suite's logging, recovery, and CORS middleware |
| Modules | `apps/api/modules/*` | `auth`, `accounts`, `mail`, `users`, `settings`, `spaces` |
| Schemas | `apps/api/schemas` | Twelve GORM models, `AutoMigrate`, the trigram indexes, and the avatar backfill |
| Internal | `apps/api/internal/*` | `database`, `middleware`, `env`, `crypto`, `authcrypto`, `authcontext`, `resourcetoken`, `oidcavatar`, `documentation` |
| Client | `apps/client` | SvelteKit 5 with `adapter-static` and `fallback: index.html` |

`modules/mail` is the biggest one and splits into `imap.go` (fetch and flag sync),
`smtp.go` (delivery), `attachments.go`, `threading.go`, `service.go`, and `router.go`.

## Middleware

Applied globally: the `tronc/httpx` stack (request logger, panic recovery, CORS with
`AllowCredentials: true`), then `middleware.SecurityHeaders`. Gzip is applied to the SPA
handler alone. Rate limits are per-route and tight on the expensive paths:

| Route | Limit |
|---|---|
| `POST /api/auth/register` | 3 per minute |
| `POST /api/auth/login` | 10 per minute |
| `POST /api/accounts/{id}/mail/sync` | 5 per minute |
| `POST /api/accounts/{id}/mail/folders/{id}/sync` | 10 per minute |
| `POST /api/accounts/{id}/mail/send` | 10 per minute |
| `POST /api/accounts/{id}/mail/emails/bulk-action` | 30 per minute |
| `POST /api/mail/test-connection` | 5 per minute |
| `/api/templates/*` | 30 per minute |

## Authentication

`middleware.RequireAuth` resolves a token from the `session` cookie first, then the
`Authorization` header, then a `?token=` query parameter, which logs a deprecation
warning. Register and login both set the cookie and return the token in the body; the
cookie is `HttpOnly`, `SameSite=Lax`, `Path=/`, with a 30-day `Max-Age`, and `Secure`
whenever the request arrived over TLS or carried `X-Forwarded-Proto: https`.

Sessions live 30 days (`auth.SessionTTL`) and are stored hashed in `sessions`.
`api_tokens` holds the long-lived per-user alternative.

### Resource tokens

Inline images inside an email body are loaded by the browser as plain `<img src>`
requests, which carry no `Authorization` header and, on a cross-origin dev setup, no
cookie either. `internal/resourcetoken` solves that with a short-lived bearer:

- `GET /api/auth/resource-token` issues an HMAC-SHA256 token over `userID:expiry`, valid
  for **5 minutes**, base64url-encoded.
- The CID image route accepts it as `?token=`, verifies the signature and expiry, and only
  falls back to the normal cookie/header path if the resource token is absent or invalid.

The HMAC secret is derived from `ENCRYPTION_KEY` — `sha256("courrier-resource-token:" +
key)` — so it is a separate secret from the credential-encryption key even though both
come from the same input. With `ENCRYPTION_KEY` unset the route is registered but returns
500; in practice the API refuses to start without it anyway.

### OIDC

Additive. `env.Load` builds an `OIDCConfig` only when `OIDC_ISSUER` is set, and then
requires `OIDC_CLIENT_ID`, `OIDC_CLIENT_SECRET`, and `OIDC_REDIRECT_URL`. When configured,
`GET /api/auth/oidc`, `GET /api/auth/oidc/callback`, and `POST /api/auth/sync-profile`
appear. `OIDC_SUCCESS_URL` defaults to the first entry of `CORS_ALLOWED_ORIGINS`.
`SSO_ONLY=true` removes register and login from the router, so they 404 rather than 403.

### Encryption at rest

`crypto.DeriveKey` takes the SHA-256 of `ENCRYPTION_KEY` and uses it as an AES-256-GCM
key. Two things are encrypted with it: the IMAP and SMTP passwords on every account row,
and the OIDC access and refresh tokens on every user row. Both have a boot-time migration
(`MigrateAccountPasswords`, `MigrateOIDCTokens`) that encrypts any plaintext value found,
so the feature back-fills itself. Both are best-effort: a failure logs a warning and
startup continues.

Changing `ENCRYPTION_KEY` after the fact makes existing credentials undecryptable. There
is no re-key path.

## Mail flow

**Sync.** `POST .../mail/sync` lists the account's IMAP folders and maps each to one of
`inbox`, `sent`, `drafts`, `trash`, `junk`, `archive`, or `custom`, recording
`uid_validity` alongside the counts. `POST .../folders/{id}/sync` fetches messages for one
folder, parses MIME with go-message, and stores headers, both body parts, and attachment
metadata. Attachment bytes are not copied into Postgres — `attachments.part_id` records
the IMAP part, and downloads stream from the server on demand.

**Send.** `POST .../mail/send` accepts JSON or `multipart/form-data` (25 MB cap for the
multipart parse, attachments spooled to temp files). Delivery is `net/smtp` with `PLAIN`
auth. Port `465` gets implicit TLS; **every other port** gets a plain dial followed by
`STARTTLS` when the server advertises it — and, notably, continues unencrypted when it
does not. After a successful send the message is appended to the account's Sent folder
over IMAP; saving a draft appends to Drafts the same way.

IMAP connects with implicit TLS only (`imapclient.DialTLS`), so a server that requires
`STARTTLS` on port 143 will not work.

**Flags.** Read and starred state written through the API is pushed back to IMAP as
message flags, and bulk delete, archive, and read-state changes go through the same path.

## Conversation threading

`modules/mail/threading.go` implements incremental union-find over Message-IDs. The
substrate is `thread_links`, which maps every Message-ID an account has ever *referenced*
— including parents that have not arrived yet — to a thread id.

For each message, `collectThreadIDs` gathers its own `Message-ID` plus every id in
`References` and `In-Reply-To`, normalized (angle brackets stripped, trimmed, lowercased)
so `<a@b>` and `a@b` compare equal. Then, inside one transaction:

| Distinct known threads | What happens |
|---|---|
| 0 | A new thread. The id is the `References` root if there is one, else `In-Reply-To`, else the message's own id — using the root keeps the id stable when replies arrive before parents |
| 1 | The message joins it |
| 2 or more | The message **bridges** them. The lexicographically smallest id wins and every `emails` and `thread_links` row on the losers is repointed to it |

Every related id is then upserted into `thread_links` against the winning thread. The
transaction is what makes a merge atomic; a failure returns an empty thread id rather than
one with no backing links.

This is why threading survives the two things that break naive `In-Reply-To` chaining:
out-of-order delivery, because a late parent finds its children through the links table,
and truncated `References` headers, because a later message carrying both halves merges the
two partial chains.

`mail.BackfillThreads` runs in a goroutine at startup and assigns thread ids to rows stored
before threading existed. It is idempotent — a populated `thread_id` is skipped — and
selects only header columns in batches of 500, so a large legacy mailbox never loads
message bodies into memory.

## Data model

| Table | Key columns | Notes |
|---|---|---|
| `users` | `id`, `email` unique, `name`, `password_hash` | Plus `avatar_url`, `avatar_source`, and the encrypted OIDC tokens |
| `sessions` | `token` PK (hashed), `user_id`, `expires_at` | 30-day sessions |
| `api_tokens` | `token` PK (hashed), `user_id`, `name` | One long-lived token per user |
| `app_settings` | single row, `encryption_key` | Surfaced through `/api/settings` as a boolean only |
| `spaces` | `id` UUID PK (`gen_random_uuid()`), `name` | |
| `space_members` | `space_id`, `user_id`, `role` | Roles are `owner`, `admin`, `member`; cascades on delete |
| `accounts` | `id`, `user_id`, `space_id`, `email`, IMAP and SMTP host/port/user/password, `signature`, `is_default` | Both passwords are AES-GCM ciphertext |
| `folders` | `id`, `account_id`, `path`, `type`, `unread_count`, `total_count`, `uid_validity` | `type` is one of the seven folder kinds |
| `emails` | `id`, `account_id`, `folder_id`, `message_id`, `thread_id`, `date`, `imap_uid`, `body_text`, `body_html` | Also `in_reply_to`, `references`, `is_read`, `is_starred`, `has_attachments`, `facile_id` |
| `thread_links` | composite PK `(account_id, message_id)`, `thread_id` | The union-find substrate |
| `attachments` | `id`, `email_id`, `filename`, `mime_type`, `size`, `part_id` | Metadata only; bytes stay on the IMAP server |
| `email_templates` | `id`, `user_id`, `space_id`, `name`, `subject`, `body_html`, `body_text` | |

`schemas.Migrate` runs `AutoMigrate` over all twelve, then `ensureSearchIndexes`, then the
avatar backfill.

### Search indexes

`ensureSearchIndexes` creates the `pg_trgm` extension and seven GIN and B-tree indexes:
trigram indexes on `emails.subject`, `from_name`, `from_address`, and `body_text`, plus
`(account_id, thread_id)` and `(account_id, folder_id, thread_id, date DESC, id DESC)` on
`emails` and `(account_id, thread_id)` on `thread_links`. Every statement is
`IF NOT EXISTS`, and the whole function returns early and silently if `CREATE EXTENSION`
fails — a database user without the privilege gets a working app with slow search rather
than a failed boot.

## Cross-app integration

Logs ship to Journal when both `JOURNAL_URL` and `JOURNAL_TOKEN` are set; the Journal SDK
wraps the `slog` handler. Courrier is **not** wired into the Nook Pool — there is no `pool`
or `enveloppe` dependency in `go.mod`, though `emails.facile_id` exists as the
cross-app identity column for when it is. Identity federation, when enabled, goes through
Authentik over standard OIDC.
