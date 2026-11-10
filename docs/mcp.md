# Courrier — MCP server

Courrier serves a remote MCP server at `https://courrier.facile.studio/mcp`. Any MCP client
that speaks OAuth (Claude, ChatGPT, Cursor, VS Code, MCP Inspector, Bulle) can search, read,
draft, organize and send the person's mail through it. It is free for every account.

## Connect a client

Paste `https://courrier.facile.studio/mcp` into the client. The client reads the
protected resource metadata, signs the person in at Registre (`https://sso.facile.studio`),
shows Registre's consent screen, and gets a 10-minute token bound to this URL.

- Claude Code: `claude mcp add --transport http courrier https://courrier.facile.studio/mcp`, then `/mcp` to sign in.
- claude.ai: Settings, Connectors, Add custom connector, URL `https://courrier.facile.studio/mcp`. Leave the
  OAuth client id and secret empty: Registre registers Claude through its client ID metadata document.

The person revokes a client in Registre under Account, Connected apps.

## How it works

| Piece | Where |
|---|---|
| Endpoint | `POST /mcp`, go-sdk v1.8.0 streamable HTTP, stateless, spec 2026-07-28 (older clients negotiate 2025-11-25). GET and DELETE answer 405 |
| Metadata | `GET /.well-known/oauth-protected-resource/mcp` (RFC 9728): `resource` is the URL above, the only authorization server is Registre |
| Tokens | `porte/mcpauth`: RS256 from Registre's JWKS, `aud` must be the URL above, `(iss, sub)` maps to the Courrier account. No token is forwarded anywhere |
| Scopes | `courrier.read` for the read tools, `courrier.write` for the write tools. A read-only token calling a write tool gets 403 `insufficient_scope` |
| Origin | A request whose `Origin` is not Courrier's own origin or a `CORS_ALLOWED_ORIGINS` entry gets 403 |
| Rate limit | 120 calls a minute per person, then 429 with `Retry-After` |
| Audit | Every write logs `mcp write` with the tool, user, client id, actor, account and outcome, never content. Journal collects it |

## Tools

| Tool | What it does | Hints |
|---|---|---|
| `courrier_list_accounts` | The accounts the person can use, with folders in detailed format | read-only |
| `courrier_search_mail` | Messages newest first, across accounts or one, by words, folder, unread, date; cursor pages | read-only, open world |
| `courrier_read_thread` | A conversation, oldest first; concise trims quotes and long bodies | read-only, open world |
| `courrier_draft` | Create or replace a draft, optionally as a reply | write |
| `courrier_send` | Send a draft, or compose and send; always confirmed by the person | destructive, open world |
| `courrier_organize` | Read, unread, star, unstar, archive, inbox, trash, move; never a permanent delete | write, idempotent |
| `courrier_delete_draft` | Delete a draft for good; always confirmed by the person | destructive |

Read tools take `response_format` (`concise` or `detailed`). Lists page with `cursor` and
`next_cursor`. A result stays under about 25,000 tokens: a longer one is cut and says how to
continue. Business errors are `isError` results that name the tool or argument to fix.

## Confirmations

`courrier_send` and `courrier_delete_draft` never run on the model's word alone.

1. The tool stores a pending action (`mcp_actions`, tronc/mcp's table) holding a digest of the draft.
2. A client on 2026-07-28 that declares form elicitation gets an `input_required` result
   asking the person. Its retry carries the answer and a `requestState` signed with a key
   derived from `ENCRYPTION_KEY`, naming the pending action, the person and a 15-minute expiry.
3. Any other client, including every client on 2025-11-25 (a stateless server cannot call
   back), gets `status: pending_confirmation` and a `confirm_url`:
   `https://courrier.facile.studio/pending?action=<id>`. Nothing happens until the person
   confirms on that page. The link stays open for 7 days.

Either way the action runs once: the first confirmation claims it. If the draft changed after
the request, nothing is sent and the action closes. A delegated (Bulle) token cannot use the
confirm page.

The transport, Origin check, scope gate, paging helpers and the whole confirmation flow come
from `github.com/FacileStudio/tronc/mcp` (v0.21.0); `apps/api/modules/mcp` keeps only the mail
tools, what a send or a draft delete runs, and the confirm page's routes. Its
`TestConformance` runs tronc's `mcptest` kit against the endpoint.
