package mcp

import (
	"errors"
	"log/slog"

	"github.com/FacileStudio/porte/mcpauth"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Tool names. Every name carries the app prefix, so a client holding many
// servers can tell them apart.
const (
	toolListAccounts = "courrier_list_accounts"
	toolSearch       = "courrier_search_mail"
	toolReadThread   = "courrier_read_thread"
	toolDraft        = "courrier_draft"
	toolSend         = "courrier_send"
	toolOrganize     = "courrier_organize"
	toolDeleteDraft  = "courrier_delete_draft"
)

// writeTools need the courrier.write scope.
var writeTools = map[string]bool{toolDraft: true, toolSend: true, toolOrganize: true, toolDeleteDraft: true}

// errNeedsWrite is the refusal inside a write tool for a read-only token.
var errNeedsWrite = errors.New("this tool needs the courrier.write scope. Reconnect Courrier and grant write access")

// toolset holds what the tool handlers share.
type toolset struct {
	mail    Mailbox
	pending pendingStore
	confirm confirmer
	appURL  string
	logger  *slog.Logger
}

// annotations builds a tool's hints. Clients treat them as untrusted UX
// hints; the scope check and the confirmations are the actual gate.
func annotations(title string, readOnly, destructive, idempotent, openWorld bool) *sdk.ToolAnnotations {
	return &sdk.ToolAnnotations{
		Title: title, ReadOnlyHint: readOnly, DestructiveHint: &destructive,
		IdempotentHint: idempotent, OpenWorldHint: &openWorld,
	}
}

// register adds the curated tools. tools/list returns them sorted by name.
func (t *toolset) register(server *sdk.Server) {
	t.registerReads(server)
	t.registerWrites(server)
}

// caller is the person a tool call acts for, and the client acting.
type caller struct {
	userID   int64
	clientID string
	actor    string
	canWrite bool
}

// callerOf reads the caller from the verified token of the request.
func callerOf(req *sdk.CallToolRequest) (caller, error) {
	info := req.Extra.TokenInfo
	identity, ok := mcpauth.Identity(info)
	if !ok {
		return caller{}, errors.New("this request carries no valid access token. Reconnect Courrier")
	}
	clientID, _ := info.Extra[mcpauth.ExtraClientID].(string)
	return caller{
		userID: identity.UserID, clientID: clientID, actor: identity.Actor,
		canWrite: mcpauth.HasScope(info, scopePrefix+".write"),
	}, nil
}

// writer is callerOf for a write tool: a token without the write scope is
// refused, even if the HTTP gate let it through.
func writer(req *sdk.CallToolRequest) (caller, error) {
	who, err := callerOf(req)
	if err == nil && !who.canWrite {
		err = errNeedsWrite
	}
	return who, err
}
