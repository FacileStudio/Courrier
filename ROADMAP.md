# Courrier — Roadmap

Self-hosted email client for the Facile Suite.

---

## Phase 1 — Foundation

Core infrastructure. Register, log in, add IMAP/SMTP accounts, see folders.

- [x] Project scaffold (Go API + SvelteKit client from Sablier base)
- [x] Auth system (registration, login, OIDC SSO)
- [x] Account CRUD (IMAP/SMTP credentials, encrypted at rest)
- [x] Database schemas (accounts, folders, emails, attachments)
- [x] IMAP connection manager (go-imap/v2, per-account connections)
- [x] Folder sync (list IMAP folders, detect type via SPECIAL-USE + name fallback, cache metadata)
- [x] IMAP/SMTP credential validation (test connection endpoint)

## Phase 2 — Read Mail

Fetch, cache, and display emails. The core email reading experience.

- [x] Email envelope fetch (subjects, senders, dates — without full bodies)
- [x] Email body fetch on demand (lazy-load text/html, cache in PostgreSQL)
- [x] Email list view (grouped by folder, unread indicator, date formatting)
- [x] Email detail view (safe HTML rendering via DOMPurify sanitization)
- [x] Plain text fallback when no HTML body
- [ ] Inline image resolution (CID → fetched content)
- [ ] Attachment list display (filename, size, type icon)
- [ ] Attachment download (lazy-fetch from IMAP by part ID)
- [x] Read/unread toggle (sync back to IMAP via \Seen flag)
- [x] Star/flag toggle (sync via \Flagged)
- [ ] Pagination / infinite scroll for large folders
- [x] Email caching strategy (IMAP UID validity, incremental sync)

## Phase 3 — Compose & Send

Write and send emails. Reply, forward, drafts.

- [ ] Compose view with rich text editor (Tiptap)
- [ ] Plain text compose mode
- [ ] Reply with quoted body (RFC 3676 format=flowed or HTML quoting)
- [ ] Reply-all with correct recipient handling
- [ ] Forward with original body + attachments
- [ ] In-Reply-To / References headers for proper threading
- [ ] File attachments (upload → attach to outgoing MIME)
- [ ] Signature per account (auto-append on compose)
- [ ] Draft auto-save to IMAP Drafts folder
- [x] SMTP send via net/smtp (STARTTLS + implicit TLS)
- [ ] Sent mail copy to IMAP Sent folder
- [ ] Address autocomplete from recent contacts / address book

## Phase 4 — Threading & Search

Group related emails, find anything fast.

- [ ] RFC 5256 threading (References + In-Reply-To + subject fallback)
- [ ] Thread view (collapsed conversation, expand inline)
- [ ] IMAP server-side search (SEARCH command)
- [ ] Local full-text search (PostgreSQL tsvector on cached emails)
- [ ] Search by sender, subject, date range, has:attachment
- [ ] Search results view with highlighting

## Phase 5 — Real-time

Push notifications for new mail. No manual refresh.

- [ ] IMAP IDLE listener (one connection per account, long-poll for new mail)
- [ ] SSE or WebSocket push to frontend (new email notification)
- [ ] Unread badge update in sidebar (live count)
- [ ] Desktop notification (browser Notification API, opt-in)
- [ ] Background sync manager (reconnect on disconnect, retry with backoff)
- [ ] Connection health monitoring (detect stale IMAP connections)

## Phase 6 — Multi-Account & Polish

Smooth multi-account experience. Keyboard shortcuts. Mobile.

- [ ] Account switcher in sidebar (show all accounts, unified or per-account view)
- [ ] Unified inbox (merge all accounts into one view, color-coded by account)
- [ ] Keyboard shortcuts (vim-style: j/k navigate, o open, r reply, c compose, / search)
- [ ] Move/copy between folders (drag or action menu)
- [ ] Bulk actions (select multiple → archive, delete, mark read)
- [ ] Folder management (create, rename, delete IMAP folders)
- [ ] Mobile responsive layout (drawer for sidebar, stacked panels)
- [ ] Resizable panels (paneforge — sidebar + list + detail, persisted sizes)
- [ ] Dark mode (already supported via Tailwind, needs testing)
- [ ] Loading states and skeleton screens

## Phase 7 — Team Features

Shared inboxes, collaboration, assignment.

- [ ] Shared team inboxes (info@, contact@, billing@ — multiple users access one account)
- [ ] Internal notes on email threads (visible to team, not sent)
- [ ] Thread assignment (assign to team member, shows in their view)
- [ ] Snooze / remind later (hide thread, resurface at chosen time)
- [ ] Labels / tags (user-defined, per-thread, filterable)
- [ ] Canned responses (templated replies for common questions)
- [ ] Collision detection (show when teammate is viewing/replying to same thread)

## Phase 8 — Suite Integration (Nook)

Connect Courrier to the Facile ecosystem via Nook event bus.

- [ ] Nook Pool client (emit events, receive events)
- [ ] `FacileEmail` canonical object shape
- [ ] `FacileContact` canonical object shape
- [ ] Courrier → Opus: link email thread to project/task
- [ ] Courrier → Glouton: incoming email from unknown sender → suggest lead creation
- [ ] Courrier → Ardoise: detect invoice-related emails, surface in billing context
- [ ] Courrier → Plume: signing-related emails tagged automatically
- [ ] Courrier → Perception: email analytics tap (volume, response times, patterns)
- [ ] Centralized outbox: other Facile tools send emails through Courrier's SMTP
- [ ] Cross-app contact sync (email addresses matched across tools)

## Phase 9 — Advanced

Nice-to-haves for power users and larger deployments.

- [ ] OAuth2 for Gmail / Outlook (XOAUTH2 IMAP/SMTP auth)
- [ ] PGP/GPG encryption and signing (OpenPGP.js on client, gpgme on server)
- [ ] Email rules / filters (auto-label, auto-move, auto-archive based on conditions)
- [ ] Calendar invites (iCalendar parsing, accept/decline, sync to external calendar)
- [ ] Contact management (address book with groups, import/export vCard)
- [ ] Email analytics dashboard (response time, volume trends, busiest hours)
- [ ] S/MIME certificate support
- [ ] JMAP protocol support (alongside IMAP, for servers that support it)
- [ ] Offline mode (service worker, cached emails available without connection)
- [ ] Browser extension (quick compose from any page)

---

## Non-Goals

Things Courrier intentionally does NOT do:

- **Run a mail server.** Courrier is a client. It connects to IMAP/SMTP servers you already run (Postfix, Dovecot, Stalwart, Gmail, etc.). Running an MTA is a separate project.
- **Replace Vero.** Vero is a personal Rust TUI mail client. Courrier is a web-based team tool. Different audiences, different architectures.
- **Marketing email / newsletters.** Courrier is for transactional and conversational email, not bulk sends. Use Plume or a dedicated tool for campaigns.
- **Be a helpdesk.** Shared inboxes and thread assignment overlap with helpdesk territory, but Courrier is an email client first. No ticket numbers, SLA tracking, or customer portals.
