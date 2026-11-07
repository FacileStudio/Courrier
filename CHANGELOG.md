# Changelog

All notable changes to Courrier are listed here. The format follows Keep a Changelog, and versions follow semver.

## [Unreleased]

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
