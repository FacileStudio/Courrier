# Courrier

Self-hosted email client for creative studios. Part of the [Facile Suite](https://github.com/FacileStudio).

Connect your IMAP/SMTP accounts, read and send email from a modern web interface. Multi-user, team-ready, integrated with the Facile ecosystem.

## Quick Start

```bash
cp .env.example .env
# Edit .env — set ENCRYPTION_KEY and DOMAINS at minimum
docker compose up --build
```

- **Frontend**: http://localhost:3000
- **API**: http://localhost:4000
- **API health**: http://localhost:4000/health

## Development

```bash
# Start PostgreSQL
docker compose up db -d

# API (Go)
cd apps/api
cp .env.example .env
go run .

# Client (SvelteKit)
cd apps/client
bun install
bun run dev
```

## Stack

- **API**: Go 1.24, Chi, GORM, PostgreSQL 16
- **Client**: SvelteKit 5, Svelte 5 runes, Tailwind CSS 4, shadcn-svelte
- **Email**: go-imap/v2 (IMAP), go-smtp (SMTP), go-message (MIME)
- **Auth**: Session tokens + optional OIDC SSO
- **Deploy**: Docker Compose (distroless Go + Nginx)

## License

MIT
