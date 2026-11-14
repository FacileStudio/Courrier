package mcp

import (
	documentation "github.com/FacileStudio/Courrier/apps/api/internal/documentation"
	mcpkit "github.com/FacileStudio/tronc/mcp"
)

// Documentation describes the confirm page routes for the API reference. The
// MCP endpoint itself is outside /api: POST /mcp speaks the Model Context
// Protocol to any MCP client holding a Registre access token for this
// server, and /.well-known/oauth-protected-resource/mcp describes how to get
// one.
var Documentation = documentation.Module{
	Name: "mcp",
	Description: "Sensitive actions MCP clients asked for, waiting for the person to confirm them: sending a " +
		"draft or deleting one. A client that cannot ask the person itself hands them a link to the pending page, " +
		"which uses these routes. The MCP endpoint is POST /mcp, outside /api.",
	Routes: pendingRoutes(),
}

// pendingRoutes documents the three confirm page routes.
func pendingRoutes() []documentation.Route {
	return mcpkit.PendingRoutes[PendingItem](mcpkit.PendingDoc{
		Auth: "the person's own browser session cookie (no bearer, CLI or API token)",
		List: "Open actions, newest first, with their draft as it is now. state is ready, changed (the " +
			"draft was edited after the request, so it cannot be confirmed) or missing.",
		Confirm: "Sends or deletes the draft, once. A draft edited since the request is refused with 412 " +
			"and the action closes.",
		Dismiss: "Closes the action without running it. The draft stays in Drafts.",
		Changed: "The draft changed, or has attachments.",
	})
}
