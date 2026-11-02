# Changelog

All notable changes to Courrier are listed here. The format follows Keep a Changelog, and versions follow semver.

## [Unreleased]

## [0.1.1] - 2026-11-02

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
