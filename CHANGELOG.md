# Changelog

All notable changes to Courrier are listed here. The format follows Keep a Changelog, and versions follow semver.

## [Unreleased]

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
