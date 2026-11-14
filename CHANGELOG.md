# Changelog

All notable changes to Courrier are listed here. The format follows Keep a Changelog, and versions follow semver.

## [Unreleased]

## [0.21.2] - 2026-11-14

### Security

- Browser session mark fix: porte v0.20.1 and porte/mcpauth v0.3.1 (vendored). porte v0.20.0 marked browser sessions with a `b-` prefix that about 1 in 4096 plain CLI or API tokens also started with, so such a token sent as a cookie could pass the MCP confirm page's own-session check. The mark is now `b~`, which a plain token never contains. Sign in once more before confirming an MCP action.

## [0.21.1] - 2026-11-14

### Changed

- build(api): tronc v0.23.1 (was v0.23.0). An MCP confirmation repeat now matches only when the OAuth client id is the same too, and Bulle errors carry `Retry-After`.

## [0.21.0] - 2026-11-14

### Changed

- fix(api): the writing helper calls Bulle's `/api/ai/complete` through `tronc/bulle` on the token exchange the Ideas routes use. `internal/bulle`, its second RFC 8693 exchange and the `/api/llm` proxy calls are gone. The helper checks `ai.access` on the personal plan: `GET /api/assist` answers `plan_required` with `upgrade_url`, and a Bulle 402 is a 402 with the upgrade link instead of "Bulle did not answer". The `no_key` and `unsupported` reasons are gone. Bulle's persona no longer prefixes the helper's prompt.
- fix(api): the inbox triage of mail the headers could not sort goes through Bulle `/api/ai/complete` with a JSON schema, as the mailbox owner, one call per pass, and only when the owner's plan has `ai.access`. Each call holds mail from one sender domain only, so a mail can steer the verdicts of its own domain's mail and no other. Triage spends at most 40 Bulle calls per owner and UTC day, backs off from 5 minutes doubling to 2 hours after a Bulle failure, and pauses until the next UTC day after a 429, so it never eats the daily AI cap the person shares with the other apps. A space mailbox is triaged only while its owner is still a member of the space. Without it only the header and sender rules sort mail. New verdicts carry `bulk_reason` `bulle`; older ones keep `jev`. `TYPESAFE_API_KEY` and `TYPESAFE_URL` are no longer read.
- fix(api): the confirm page API accepts only the person's own browser session (`mcpauth.RequireOwnSession`): a bearer, CLI or API token is refused. People sign in once after the upgrade before they can confirm.
- fix(api): the helper and the triage measure their input once encoded, as Bulle does, so a text full of quotes or angle brackets stays under Bulle's 96 KiB body limit, and a Bulle 413 shows as "text too long for the helper".
- fix(client): the confirm page tells a person whose session cannot confirm (an old cookie, the mobile or desktop app) to sign in again in the browser, with a sign-in button, instead of the raw error.
- fix(api): `courrier_draft` is annotated destructive, since replacing a draft by `draft_id` overwrites it.
- build(api): tronc v0.23.0, porte v0.20.0, porte/mcpauth v0.3.0, porte/forfait v0.1.3, tronc/testdb v0.1.2. The MCP state key comes from `mcpkit.StateKey` (same bytes), the confirm page mounts with `PendingPage.Mount`, pending-action texts come from the kit, and `mcp_actions` rows resolved or expired over 30 days ago are pruned hourly.
- build: mise pins Go 1.26.11, as go.mod and the Dockerfile do.

## [0.20.0] - 2026-11-14

### Added

- feat(client): the landing shows a looping inbox demo built with `@facile/mise-en-scene` v0.2.1: a mail arrives, opens as a thread and gets a reply. Phones get a compact layout.
- feat(client): the landing download section lists only the visitor's platform (muse `DownloadSection mode="device"`), with a link to `/downloads` for the rest.

### Changed

- build(client): muse v0.23.0 (was v0.21.0). Titles and body text use Magnetik, and the landing footer starts at the left edge on phones.
- feat(client): `/downloads` uses muse `DownloadSection` and `LandingHeader`/`LandingFooter`, with the same release URLs plus the CLI install commands.
- feat(client): shorter landing copy in en, fr, es and de, and a fourth feature block for the inbox tabs.

## [0.19.7] - 2026-11-14

### Added

- feat(client): the web client loads Vision's tracker (`https://vision.facile.studio/s.js`, with `data-mask`) on `courrier.facile.studio` only. The API adds `https://vision.facile.studio` to `script-src`; the loader's hash comes with the other inline scripts of the build.

## [0.19.6] - 2026-11-14

### Fixed

- fix(deploy): drop the swarm-only `update_config` and `restart_policy` from the compose services, so `restart: unless-stopped` applies under plain docker compose.

## [0.19.5] - 2026-11-13

### Changed

- build(client): muse v0.21.0 (was v0.18.0). Icon-only controls in muse components now show a muse tooltip.

### Fixed

- fix(client): every icon-only control shows a muse tooltip, and no control uses the native `title`: the mail and compose toolbars, the editor marks, the desktop title bar, the chip remove buttons, the intro dots, the Bulle assist trigger and the account check and default buttons.
- test(client): the updater test loads its module before its timed tests, so a slow cold import no longer fails it.

## [0.19.3] - 2026-11-13

### Changed

- chore(deploy): healthchecks run every 30 s with a 1 s start interval.

## [0.19.2] - 2026-11-13

### Changed

- fix(api): `courrier.mailboxes` defaults to -1 (unlimited) instead of 1. A person or space with no entitlements in the Registre snapshot (outside every Registre space, or before the first push) can now add any number of mail accounts, matching the free plan in Ardoise.

## [0.19.0] - 2026-11-13

### Changed

- refactor(client): API calls go through `@facile/client` v0.2.0. The web build uses `createApi` (relative `/api`, the session cookie, `X-Facile-CSRF` on writes, the tronc error envelope); the installed apps use `@facile/client/native`'s `createNativeApi` (the bearer token to production, no cookie, no CSRF header, the token forgotten on a 401). Every call keeps its method, path, headers and body; Courrier's `ApiError` keeps `retryAfter` and `planLimit`, read from the kernel's error.
- feat(client): on the web, a 401 from any call other than the session probe (`me()`, which is quiet) sends the person to `/login`, as across the suite. The layouts still handle a signed-out session themselves, so the `/pending` confirm-link return path is unchanged.
- refactor(client): the theme store is `@facile/client/theme`'s `createTheme('courrier.theme', { resolveSystem: true })`, under the same storage key, still writing `.dark` or `.light` in system mode.
- build(client): `bun run dev` proxies `/api` to `127.0.0.1:4000`, so development is same-origin like production. `VITE_API_BASE_URL` now only points the installed apps at another API.

## [0.18.0] - 2026-11-13

### Security

- feat(client): the installed apps keep their bearer token in the platform's secure store instead of `@capacitor/preferences`, which is plain storage (SharedPreferences on Android, `localStorage` in the desktop webview). Android uses `capacitor-secure-storage-plugin` 0.13.0 (Keystore); the desktop app uses the OS keychain through the `keyring` crate (Keychain, Credential Manager, Secret Service) behind the `token_get`, `token_set` and `token_clear` commands. The web build is unchanged.
- feat(client): on first launch the app moves a token an older release kept in Preferences into the secure store and deletes the Preferences copy, so nobody is signed out. If the secure store refuses the write, the old copy stays and keeps working. Signing out clears both.
- build(desktop): the Linux build installs `libdbus-1-dev` for the keychain.

## [0.17.0] - 2026-11-13

### Changed

- feat(api): the schema is ordered goose SQL under `apps/api/migrations/`, applied at boot by `tronc/migrate` v0.1.1 behind a Postgres advisory lock. GORM `AutoMigrate` and the hand-written `CREATE TABLE IF NOT EXISTS` (push devices, plan state, deletion state, the MCP actions table, porte's tables) are gone from the boot path; the GORM structs stay as the query layer. `scripts/check.sh` fails on an `.AutoMigrate(` call outside a test file.
- feat(api): `00001_baseline.sql` is production's schema as `pg_dump --schema-only` printed it on 2026-11-13, leftover tables and columns included. A database built before this release (Courrier's tables, no applied migration) is stamped at version 1 on its first boot, so the baseline is recorded without running. A database built from the baseline dumps identical to production. `pg_trgm` is created in `public` and its operator class is qualified, so a fresh install that cannot create the extension now fails its first boot instead of silently skipping the search indexes.
- feat(api): the data steps that ran at every boot still run at every boot, after the migrations and in the same order (`schemas.Reconcile`): NULL email class columns to `''`, the avatar URL and source backfills, the porte identity re-key, the pre-porte session and API token carry-over, and the password and OIDC subject adoption with the Registre picture backfill. `crypto.MigrateAccountPasswords`, `crypto.MigrateOIDCTokens` and the thread backfill are unchanged.
- build: `github.com/FacileStudio/tronc/migrate` v0.1.1 with `pressly/goose/v3` v3.27.3; goose raises `segmentio/asm` from v1.1.3 to v1.2.1. Re-vendored.

### Added

- feat(api): `docker run <image> migrate status` (also `up`, `version`, `baseline <n>`) answers without a shell.
- feat(api): `00002_column_defaults.sql` gives the timestamps GORM stamped at insert a `now()` default, and the users columns porte/kit does not write the values `fillNewAccount` writes (`''` and `0001-01-01 00:00:00+00`), so a row inserted outside GORM no longer gets NULL there. It has a Down.
- test: the GORM structs are checked against the migrated schema, and `tronc/mcp.PendingSchema` and porte/pg's identity, session and login code tables are pinned by drift tests.

### Rollback

- Rolling back to v0.16.x runs `AutoMigrate` again, which drops the defaults `00002_column_defaults.sql` added while `goose_db_version` still records version 2. Before deploying this release again after such a rollback, run `DELETE FROM goose_db_version WHERE version_id = 2`. Rehearsed: old, new, old, new boots on one database all succeed and keep every row.

## [0.16.5] - 2026-11-12

### Security

- fix(auth): an account re-created for a person Registre erased is erased again on the next deletion push or poll, without a restart (porte v0.19.0, `directory.Deletions.Holds`). Before, a bearer token minted before the erasure, or an MCP first call, could bring the account back and it stayed until the next restart.
- fix(auth): sign-in, bearer and MCP provisioning, and profile sync refuse an identity Courrier already erased (`oidc.Deps.Erased`), and the Registre space sync no longer pre-provisions an erased member from a stale snapshot.
- build: porte v0.19.0 and porte/directory v0.5.0, which also raise go-chi/chi to v5.5.0, golang.org/x/crypto to v0.58.0 and go-jose to v4.1.9. Re-vendored.

## [0.16.4] - 2026-11-11

### Fixed

- fix(brand): favicons, touch and manifest icons, the link preview image and `site.webmanifest` come from muse v0.19.2's brand generator, and `app.html` carries its head block: the catalog title "Courrier · Your email, on your server", its description, theme colours and Open Graph and Twitter tags. The landing page no longer overrides the title and description, so every locale shows the catalog's.
- feat(brand): the sidebar, landing header, footer and sign-in screen show the catalog icon `solar:letter-opened-bold-duotone`.

## [0.16.3] - 2026-11-11

### Fixed

- build: clear the douane high and medium findings. Client: vite 7.3.7 to 7.3.8, @sveltejs/kit 2.70.4, svelte 5.57.3, @tiptap/* 3.33.1 (prosemirror-view 1.42.3+), and refreshed transitive deps (devalue, postcss, picomatch, nanoid, source-map-js). `overrides` pin cookie ^0.7.2 and esbuild ^0.28.1. API: go-chi/chi v5.3.1 to v5.4.0, re-vendored.
- build: drop the stale `apps/client/package-lock.json`. Bun is the only installer (Dockerfile, CI, `scripts/check.sh`), and the npm lockfile still pinned vite 7.3.1.
- build: left open: glib 0.18 in the desktop app (the fix is glib 0.20, which Tauri 2's gtk 0.18 stack cannot take), the unmaintained zbus/zvariant and proc-macro-error crates (no fix), @sveltejs/kit findings fixed only in 3.0.0 (low) and deepmerge (no fix).

## [0.16.2] - 2026-11-11

### Security

- `GET /api/users` and `GET /api/users/{id}` show only the caller and the people it shares a space with. Any other id answers 404. Before, every signed-in account read every account's name and email. Accounts whose email is on the new `COURRIER_ADMINS` list (comma-separated, empty by default) still see every account.
- `POST /api/spaces/{id}/members` takes a `user_id` only for someone the caller already shares a space with. Anyone else is added by exact `email`. Before, an owner could add guessed ids to its own space and read their addresses in the member list.
- The member picker lists the people you already share a space with and adds anyone else by exact email address.

## [0.16.1] - 2026-11-11

### Changed

- refactor(api): the MCP module runs on `tronc/mcp` v0.22.0 (`NewApp`, `PendingPage`, `PendingRoutes`, `Digest`, `ErasePending`). The local server setup, the three `/api/mcp/pending` handlers and their list, result and empty-request types are gone. Tool names and schemas are unchanged, and open pending actions still confirm (same digest). `GET /api/mcp/pending` items gain `summary`; the delegated-token 403 reads "only the person's own session can confirm actions" (code still `permission_denied`); the reference documents the `id` path param and an optional `url` in the result.
- build: tronc v0.22.0, Go 1.26.11 (`go` directive and the `golang:1.26.11-alpine` builder) for the stdlib advisories GO-2026-6848 to GO-2026-6859.

## [0.16.0] - 2026-11-10

### Added

- User menu: pending Registre space invitations, with a count and a link, in the sidebar and the phone bar. `GET /api/me/invitations` (porte/kit) answers `{"count", "href"?}`, never fails a page, caches per user for 60 seconds, and is mounted only with SSO on.

### Changed

- Licence: Courrier is now licensed under `AGPL-3.0-only` (see `LICENSE`); before this, the repository carried no licence.
- refactor(api): the sign-in glue runs on `porte/kit`: its user store over Courrier's `users` table, its middleware and request identity, and its `EraseAccount`. `internal/authcontext`, `internal/oidcavatar` and the hand-written user store are gone. Routes, status codes and error bodies are unchanged. By owner decision, a sign-in now rewrites the stored email only when Registre marks it verified. Accounts the kit creates are filled exactly like the older ones (only `oidc_subject` stays NULL), pinned by a test.
- build: tronc v0.21.0, porte v0.18.2 (`directory` v0.4.0), muse v0.18.0, Go 1.26.10 (stdlib criticals), `golang.org/x/crypto` v0.56.0, `golang.org/x/text` v0.41.0, `go-jose` v4.1.6, dompurify 3.4.17. The Dockerfile builds on `golang:1.26.10-alpine` and copies `go.mod` explicitly.
- ci: filet v0.19.1 with `preset: facile`, and a douane job (`-fail critical -scope prod`). `douane.yml` permits GO-2026-5932 (x/crypto `openpgp`, not compiled in) until 2027-02-07.
- The pre-push hook moves from `.githooks` to `lefthook.yml`: `scripts/check.sh`, then `filet gate .`. Run `mise run hooks` once in each clone.

### Fixed

- Erasure: when Registre erases a person, Courrier now also deletes their open MCP actions (`mcp_actions`), and deletes sessions, API tokens and push devices explicitly. A test drives the deletion poller against a fake Registre and asserts no row with a `user_id` still references the person, apart from mailboxes and templates that belong to a space.
- Client: the bearer token the web app kept in `localStorage` before the cookie session (`courrier.token`) is deleted on load. The installed apps keep theirs.
- api: FCM credentials are read with `CredentialsFromJSONWithType` and must be a service account; two unused functions are removed.

## [0.15.1] - 2026-11-10

### Fixed

- MCP: a person's first call to `/mcp` now creates their Courrier account, exactly as a first web sign-in would, so a new Registre user can connect an MCP client without ever opening Courrier on the web. Before, a valid token for a person Courrier had never seen answered `401 invalid_token`. porte/mcpauth v0.2.0 (`Provision: true`) runs the OIDC kit's own login write path; it refuses a token without a verified email. Registre v0.10.3 carries the profile on resource-bound tokens.

### Changed

- build: porte v0.15.1 to v0.16.0 and porte/mcpauth v0.1.0 to v0.2.0.

## [0.15.0] - 2026-11-10

### Added

- Ideas: a lightbulb entry first in the sidebar and the phone bar, with a badge for new ideas, and an `/ideas` page where Bulle suggests which mails to answer, follow up or sort in the active space. An idea runs in place: Courrier polls the run, shows an approval card for each held call, then the result with its links and a link to continue in Bulle. Ideas are paid (`ai.access`): a free plan sees the upsell and Registre's plans link.
- API: `GET /api/ideas`, `POST /api/ideas/generate`, `POST /api/ideas/{id}/run`, `POST /api/ideas/{id}/dismiss`, `GET /api/ideas/runs/{run}` and `POST /api/ideas/runs/{run}/approvals/{call}`, mounted when SSO and `BULLE_URL` are set. They call Bulle as the signed-in user through a Registre token exchange (audience `bulle`, scopes `bulle.read` and `bulle.write`). `space_id` must be a space the caller belongs to and that Registre knows. Courrier declares `ai.access` as a plan feature and answers 402 with the upgrade link, or `ai_access: false` on the list, without calling Bulle when the plan lacks it.

### Changed

- porte to v0.15.1 (`porte/exchange`). Client: `@facile/muse` v0.16.1 to v0.17.1.

## [0.14.1] - 2026-11-10

### Changed

- refactor(api): the MCP server runs on `tronc/mcp` v0.19.0, the kit extracted from this module. The transport, Origin check, scope gate, paging helpers, confirmation flow and audit now live in tronc; `modules/mcp` keeps the mail tools. Tools, arguments, statuses and the confirm page are unchanged.
- Pending actions move to tronc's `mcp_actions` table, created at boot (additive). `mcp_pending_actions` is no longer read and is left in place; it held no open action on 2026-11-10.
- Write audit records come from tronc: one per write tool call and one per confirmation event, with the target in `target`. A draft save or an organize call no longer logs its account id.
- tests: `TestConformance` runs `mcptest.Run` (401 challenge, metadata, 405, Origin, annotations, schemas, order, write scope, and the confirmation path of `courrier_send` and `courrier_delete_draft`). Confirmation tests run without a database on the in-memory store.

## [0.14.0] - 2026-11-10

### Added

- A remote MCP server at `/mcp` for any MCP client, free for every account (docs/mcp.md). Registre is the authorization server: tokens are bound to `MCP_RESOURCE`, scoped `courrier.read` and `courrier.write`, and verified by `porte/mcpauth` v0.1.0. go-sdk v1.8.0, stateless, spec 2026-07-28.
- Seven tools: `courrier_list_accounts`, `courrier_search_mail`, `courrier_read_thread`, `courrier_draft`, `courrier_send`, `courrier_organize`, `courrier_delete_draft`, each annotated, with concise or detailed output, cursor pages and a 25,000-token cap.
- Sending and deleting a draft need the person's confirmation: MRTR `input_required` when the client can ask, otherwise a pending action and a link to the new Pending actions page (`/pending`). Every write is logged with the client id.
- The `mcp_pending_actions` table (created at boot, additive).

### Changed

- porte to v0.14.0 and porte/forfait to v0.1.1: a higher entitlements revision now wins over a lower version, so a Registre restored from backup no longer freezes the plan values.
- The per-address rate limiter takes a key function and locks its counters.

## [0.13.4] - 2026-11-07

### Changed

- build(client): bump muse to v0.16.1, replace the deprecated `outline` button variant with `secondary`, and the `danger` and `default` icon buttons with `ghost-danger` and `secondary`.
- build(api): `GET /version` is baked at build time from `apps/client/package.json` by `scripts/version.sh`. The `APP_VERSION` build arg is gone from the Dockerfile, compose file and docs, so a stale Dokploy env can no longer show the wrong version. The response is unchanged.

## [0.13.3] - 2026-11-07

### Changed

- build(client): bump muse to v0.14.4 (from v0.14.2). Switching the language in the user menu re-renders the whole UI at once, with no reload. A single 100% DonutChart slice draws as a full ring.

## [0.13.2] - 2026-11-07

### Security

- Registre tokens are sealed at rest. porte to v0.13.0: `porte_identities.access_token` and `refresh_token`
  are encrypted with AES-256-GCM under `PORTE_TOKEN_KEY`. Existing rows are sealed once, in the background,
  after boot, and the count is logged. An invalid `PORTE_TOKEN_KEY` stops the boot. With the key unset nothing changes.
- Losing `PORTE_TOKEN_KEY` makes the stored tokens unreadable: people sign in again, nothing else is lost.
  Do not roll back to a build older than porte v0.13.0 after the first reseal.

## [0.13.1] - 2026-11-07

### Added

- A Subscription section in the settings: the mail account allowance of the active space (or of your personal plan), how many accounts are used, and links to manage or upgrade the plan on Facile. It reads the new `GET /accounts/usage`, which answers `{used, limit}` (`limit` is -1 when unlimited).
- The user menu gains a Subscription row that opens the plans on Facile (`https://sso.facile.studio/spaces`).

### Changed

- build(client): bump muse to v0.14.2. The API tokens section drops the expiry field and the last-used column: Courrier tokens do not expire and do not record their use.

## [0.13.0] - 2026-11-07

### Changed

- `@facile/muse` v0.13.0 to v0.14.1. The sidebar user card and the mobile avatar open muse's user menu: Facile account, Downloads (browser only), Settings, language, theme, help and Sign out. The hand-made `SidebarFooter` popover and its drawn rule are gone, and so is the cog.
- Settings use muse's `SettingsLayout` in the order Profile, Preferences, Notifications, Accounts, Templates, API tokens, Advanced (plus Updates on the desktop app). Profile is read-only and links to the Facile account, since Registre owns the identity. The Appearance tab became Preferences (language and theme), and the in-settings Sign out is in the user menu. The API token page uses `TokensSection` and the Advanced danger zone uses `DangerZone`.
- `/settings/appearance` redirects to `/settings/preferences` and `/settings/api` to `/settings/tokens`.
- Choosing a theme no longer fails when the browser blocks local storage.

### Removed

- The local name edit, photo upload and the profile meta rows, and the unused client calls and strings behind them.

## [0.12.1] - 2026-11-07

### Fixed

- `porte` v0.12.0 to v0.12.1. porte's identity eraser no longer fails on a database without the `porte_magic_links` table, which Courrier does not create.

## [0.12.0] - 2026-11-07

### Added

- Sign-in goes through Registre only. The login page is muse's `AuthScreen`, and the landing page uses `LandingHeader`, `LandingFooter` and `SignInWithFacile` in place of the hand-made top bar, footer and buttons. The desktop and Android sign-in button is `SignInWithFacile` too.
- The whole interface is available in French, English, Spanish and German, in the language of the Registre profile (`users.locale`, read from the OIDC `locale` claim at every sign-in), then the language picked on a public page, then the browser's. Dates, times, sizes and counts use `Intl` in that language. The spaces pages, which were French only, are translated. Spanish and German need native review.
- GDPR deletion. When Registre erases a person, Courrier deletes their account, personal mail accounts with the folders, messages, attachments records, contacts, thread links and sender rules cached from them, personal templates, memberships and avatar file. Mail accounts and templates attached to a space stay. `POST /api/registre/deletions` takes Registre's signed notice, and the receiver also polls Registre. A new table, `directory_deletion_state`, keeps its cursor.
- `GET /version` answers the version set by the `APP_VERSION` build argument.
- `.github/workflows/ci.yml` runs gofmt, vet, build and tests against Postgres, the client check, tests and build, and filet. It replaces `filet.yml`.

### Changed

- `porte` v0.8.0 to v0.12.0, `porte/directory` v0.2.0 to v0.3.0, `tronc` v0.16.0 to v0.17.1, `@facile/muse` v0.11.0 to v0.13.0. Upgrade notes applied: the two binding columns on `porte_login_codes`, `users.locale`, porte's account linking rules in the user store (link by verified email only, refuse an account that already holds another subject of the same issuer), tronc's security headers, `go 1.26` in `mise.toml`, and muse v0.13's borderless forms and custom `Select`.
- Security headers come from tronc. The local `SecurityHeaders` middleware is gone. The Content-Security-Policy stays report-only; framing is denied and enforced. The client's inline script hashes are added, and `connect-src` allows `https://api-ardoise.facile.studio` (plans) and `https://api.iconify.design` (icons).
- `OIDC_ISSUER` is required when `APP_ENV=production`.
- `PATCH /api/users/me` takes `name` only. The email comes from Registre, and a signed-in person can no longer park an account on someone else's address.
- Form controls, menus and dividers lose their borders and rings, and the deprecated `fc-border` token is gone.

### Removed

- `POST /api/auth/register`, `POST /api/auth/login`, the password change on `PATCH /api/users/me`, the register and password screens, and `SSO_ONLY`. API tokens (Settings > API), the CLI's device and loopback login, the Android app's code exchange, delegated bearer JWTs and the mail accounts' own IMAP and SMTP passwords are unchanged.

## [0.11.0] - 2026-11-06

### Added

- The landing page has a "Pricing" section: muse's `PricingPlans` shows the live Facile plans from Ardoise's catalog with `courrier.mailboxes` first. It renders nothing while the catalog is empty or unreachable.

### Changed

- `@facile/muse` v0.8.1 to v0.11.0. The root layout calls `provideLocale(() => 'en')`, so muse's own strings (pricing toggle, plan cards) are English like the rest of the shell and `<html lang>` stays `en`. The spaces pages keep their explicit French text.
- The landing page's "Get the desktop and Android apps" link is muse's `DownloadCta`, with the macOS, Windows, Linux and Android icons. It still goes to `/downloads`.

## [0.10.1] - 2026-11-06

### Fixed

- The personal mail account limit now reads the right person's plan when a user holds several Registre identities, for example a stale row from an earlier Authentik. Courrier takes the subject the last directory snapshot lists among its people. If the snapshot lists none, it takes the most recently synced identity. Before, it took whichever row Postgres returned first.

## [0.10.0] - 2026-11-04

### Added

- Inbox tabs: Primary, Newsletters and Notifications (`tab` on the folder list), filed by the header classifier.
- `PUT /api/accounts/{id}/mail/sender-rules` and a "Move sender to" picker pin every message from one address to a tab.
- After each inbox sync, mail the headers could not sort, from senders the user never wrote to, gets one jev decision (`TYPESAFE_API_KEY`); only a confident newsletter or notification verdict moves it.
- The Bulle helper in the composer: a blob button in the toolbar and on any selection to rewrite, shorten, expand, fix, change tone, translate or draft a reply, with Replace, Insert below, Try again and Discard.
- Threads get Summarize and an Ask Bulle link. `GET`/`POST /api/assist` run through the user's own Bulle model via a Registre token exchange (`BULLE_URL`); Courrier holds no model key.

### Fixed

- Search matches every word anywhere, so "Lisa Nay" finds a sender stored as "NAY Lisa".

## [0.9.0] - 2026-11-04

### Added

- Plans cap mail accounts per scope through `courrier.mailboxes` (porte/forfait v0.1.0, default 1). A space's accounts count against its Registre space's plan, and personal accounts against the owner's Registre subject. Only adding an account is refused; accounts already over the limit keep reading, syncing and sending.
- A refused creation answers HTTP 402 with `{"error":{"code":"payment_required","message","upgrade_url","key","limit"}}`, and the account drawer shows an upgrade notice linking to Registre. Only an https or same-origin `upgrade_url` becomes a link.
- The plan values ride the Registre directory pull and are saved in the new `forfait_state` table after every pull, so a restart keeps enforcing the last known values. Until the first entitlements snapshot ever arrives, every creation is allowed and a warning is logged.

### Changed

- Once plans are enforced, a space not linked to Registre refuses new mail accounts with 412, because its plan cannot be read.
- If the saved plan row cannot be read at boot, Courrier enforces the free defaults until the next pull instead of allowing everything.
- Space ids from clients are canonicalised (lowercased) before they are stored or filtered on, for accounts and templates. `""` means the personal scope, and braced or unhyphenated spellings are refused with 400. At boot Courrier logs, without changing them, rows whose `space_id` is not canonical.
- porte/directory v0.2.0.

### Security

- Creating a template in a space now requires membership of that space.

## [0.8.1] - 2026-11-03

### Changed

- Spaces use the share-circle icon, through muse v0.8.1.

## [0.8.0] - 2026-11-03

### Added

- Desktop self-update system with a dedicated updates settings page, status badges, download progress bar, and action triggers.
- Root layout checks for desktop updates in background on startup.

## [0.7.1] - 2026-11-03

### Fixed

- The desktop app account menu no longer shows a download option for the mobile app.

## [0.7.0] - 2026-11-03

### Added

- The desktop app shows native system notifications for new mail, on by default. Its webview has no push service, so it listens on `GET /api/push/stream`, a server-sent event stream, and the inbox watcher also polls the inboxes of users with an open stream. When the window has focus the notification shows as a toast instead.

### Changed

- The web client uses muse v0.8.0 and drops decorative borders on cards, panels, menus, headers and toolbars. Surfaces separate by background tone and spacing. Inputs keep their borders.
- The desktop window has no system title bar. On Windows and Linux it draws its own close, minimize and maximize controls in a slim drag strip; on macOS it keeps the native traffic lights over the content.
- `GET /api/push` always reports `enabled: true`, and `POST /api/push/test` no longer answers 409: the desktop stream works without FCM or Web Push.

## [0.6.0] - 2026-11-03

### Added

- Registre owns shared spaces when Facile SSO is on. At boot Courrier hands Registre its existing spaces and members once, retrying with backoff while Registre is not deployed, and then keeps every space, member and role equal to Registre's: at boot, every 5 minutes, and within seconds of Registre's nudge on `POST /api/registre/sync`. A member added in Registre gets a Courrier account before their first sign-in. Mail accounts attached to a space stay Courrier's own.
- `spaces.registre_space_id` links a space to its Registre twin. The column and its unique index are added at boot.

### Changed

- With SSO on, the spaces pages show "Gérer dans Registre" in place of creating a space, renaming it, adding, removing or re-ranking members and deleting it. Member lists stay, read-only, with pictures.
- With SSO on, `POST /api/spaces` and every edit of a linked space (`PUT`/`DELETE /api/spaces/{id}`, its members routes and `leave`) answer 409 with the Registre page to use. Password-only instances manage spaces as before.

## [0.5.1] - 2026-11-03

### Fixed

- The account menu items stay on one line: they read "Desktop app" and "Android app".
- The sidebar user card sits inside the sidebar, with one background, instead of a second card below it.

## [0.5.0] - 2026-11-03

### Added

- A Tauri 2 desktop app (`apps/desktop`) for macOS, Windows and Linux. It signs in through the same Facile SSO flow as the Android app (system browser, then the `studio.facile.courrier://auth` link), keeps a bearer token, and updates itself from the GitHub release's `latest.json`.
- `desktop.yml` builds the signed installers when a release is published and attaches `Courrier-macos-universal.dmg`, `Courrier-windows-x64-setup.exe`, `Courrier-linux-x64.AppImage`, `Courrier-linux-x64.deb`, `desktop-checksums.txt` and `latest.json`.
- The downloads page lists the desktop installers, with the visitor's system highlighted.

### Changed

- The sidebar user card opens a menu with Settings, the desktop app download and Sign out, in place of the settings cog. The desktop app shows the Android download there instead.
- The API allows the desktop app origins `tauri://localhost` and `http://tauri.localhost` while SSO is on.
- The landing page links read "Get the apps".

### Fixed

- Inline (`cid:`) images load in the Android and desktop apps, which have no session cookie: the thread view now adds a short-lived resource token to their URLs.

## [0.4.0] - 2026-11-03

### Added

- Space members carry `avatar_url`, and the space overview and members pages show each member's picture, with the initial when there is none.

### Changed

- Every federated user gets the identity provider's avatar URL (`<OIDC_ISSUER>/avatars/<subject>`) at boot when none is stored, so a Facile SSO picture shows without signing in again. A changed picture shows on the next page load.
- The sidebar card renders the picture with muse `Avatar`, so a missing picture falls back to the initial instead of a broken image.
- The profile page sends photo changes to the Facile account page at sso.facile.studio.
- muse is pinned to v0.6.5.

## [0.3.0] - 2026-11-03

### Added

- Push notifications for new inbox mail: the Android app through Firebase Cloud Messaging (`FCM_SERVICE_ACCOUNT`), browsers through Web Push (`VAPID_PUBLIC_KEY`, `VAPID_PRIVATE_KEY`, `VAPID_SUBJECT`). A watcher syncs the inbox of every account whose owner registered a device every `MAIL_POLL_INTERVAL` (default `2m`). Bulk and automated mail never notifies, and a notification names the sender, never the subject or the body.
- `GET /api/push`, `POST /api/push/test`, `POST /api/push/devices` and `DELETE /api/push/devices`, with the `push_devices` table.
- The Android app can sign in through porte's mobile flow (`OIDC_MOBILE_REDIRECT`, default `studio.facile.courrier://auth`). The API allows the app origin `https://localhost` while SSO is on.
- The client builds as a Capacitor 8 Android app (`studio.facile.courrier`): an intro carousel, sign-in through the system browser, and a bearer token in place of the session cookie.
- A first-run flow at `/welcome`: connect a mailbox or pick the one that opens first, turn on notifications, done. It runs once per device in the app, and on the web only for someone with no mailbox.
- Settings > Notifications turns push on for this phone or browser and sends a test.
- A public `/downloads` page with the Android APK, its SHA-256 checksum and install steps, linked from the landing page header, hero and footer. The Android app hides it.
- `.github/workflows/release.yml` builds the signed APK when a release is published and attaches `courrier-android.apk` and its `.sha256`.

### Changed

- porte v0.8.0.

### Fixed

- A delegated suite token (one Registre minted for another app acting for the user) can no longer mint or revoke API tokens or edit the sign-in profile.

## [0.2.7] - 2026-11-02

### Changed

- The mailbox switcher sits at the bottom of the main sidebar. The settings card is alone in the separate card below it.

## [0.2.6] - 2026-11-02

### Changed

- The mailbox switcher and the settings card sit in their own card under the sidebar, switcher first.

## [0.2.5] - 2026-11-02

### Changed

- The bottom of the sidebar shows the mailbox switcher first, then the settings card, both inside the sidebar.

## [0.2.4] - 2026-11-02

### Changed

- Compose is the first row of the sidebar again.
- The account switcher sits under the sidebar and shows the email address. It shrinks to the avatar when the sidebar collapses.

### Removed

- The floating action pill at the top of mail views. It covered the search bar and the mail titles.

## [0.2.3] - 2026-11-02

### Fixed

- Settings > Advanced no longer says ENCRYPTION_KEY is missing when it is set. It used to read an unused database column instead of the server environment.

### Removed

- `PUT /api/settings`. It only wrote that unused column.

## [0.2.2] - 2026-11-02

### Fixed

- The action pill floats over the mail view instead of pushing it down.

### Changed

- The account switcher shows the account's email address instead of its name.

## [0.2.1] - 2026-11-02

### Changed

- Compose and the account switcher sit in a floating pill at the top centre of mail views on desktop. The bottom-right compose button is gone.
- On phones, Compose is back in the bottom nav and the account switcher is in the top bar.
- The sidebar shows four folders: Inbox, Sent, Drafts and Archive. Junk and Trash are in the folder dropdown above the mail list, now on desktop too.

### Fixed

- The account switcher no longer crowds the mail list header.

## [0.2.0] - 2026-11-02

### Added

- Scoped delegated access. With `OIDC_MACHINE_AUDIENCE` set to Courrier's client id and `OIDC_SCOPE_PREFIX=courrier`, a short-lived token minted by Registre for another suite service opens the API. `courrier.read` allows GET, HEAD and OPTIONS. Anything that changes data needs `courrier.write`. Both variables are empty by default, which leaves behaviour unchanged.

### Changed

- porte v0.5.2 to v0.7.0. The Go directive moves to 1.26, matching the builder image.

## [0.1.0] - 2026-11-02

First tagged release.

### Added

- Account switcher in the mail list header. The choice is per browser and per space, and never changes the server default.
- From picker in compose when more than one account exists. Replies send from the account that received the mail.
- Floating compose button at the bottom right of mail views.

### Changed

- Account settings use one card per account instead of a table. Remove moved into the edit drawer, and "Make default" is now "Set as default".
- Compose left the sidebar and the mobile nav.

### Fixed

- HTML emails no longer render white text on a white background in dark mode. Bodies render in a sandboxed light frame, which also stops an email's styles from leaking into the app.
- Drafts saved under one account are deleted from that account after switching the From address.
