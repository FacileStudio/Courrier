# Courrier

Self-hosted email client for creative studios. Part of the [Facile Suite](https://github.com/FacileStudio).

Connect your IMAP/SMTP accounts, read and send email from a modern web interface. Multi-user, team-ready, integrated with the Facile ecosystem.

A single Go binary serves both the API and the static SvelteKit frontend.

## Quick Start

```bash
cp .env.example .env
# Edit .env — set ENCRYPTION_KEY at minimum
docker compose up --build
```

Open http://localhost:4000 — the API and frontend are both served from there.

## Development

```bash
# Start PostgreSQL
docker compose up db -d

# Terminal 1 — API (Go)
cd apps/api
cp .env.example .env
go run .

# Terminal 2 — Client (SvelteKit)
cd apps/client
cp .env.example .env    # Sets VITE_API_BASE_URL=http://localhost:4000
bun install
bun run dev
```

The client dev server runs on http://localhost:5173 and proxies API calls to http://localhost:4000.

## Stack

- **API + Server**: Go 1.24, Chi, GORM, PostgreSQL 16 (single binary serves API + frontend)
- **Client**: SvelteKit 5, Svelte 5 runes, Tailwind CSS 4, shadcn-svelte
- **Email**: go-imap/v2 (IMAP), go-smtp (SMTP), go-message (MIME)
- **Auth**: Session tokens + optional OIDC SSO
- **Deploy**: Docker Compose (single distroless container + PostgreSQL)

## License

MIT
