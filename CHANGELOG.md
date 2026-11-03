# Changelog

All notable changes to Courrier are listed here. The format follows Keep a Changelog, and versions follow semver.

## [Unreleased]

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
