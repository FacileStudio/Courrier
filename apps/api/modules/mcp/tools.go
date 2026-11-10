package mcp

import (
	"errors"

	"github.com/FacileStudio/porte/mcpauth"
	mcpkit "github.com/FacileStudio/tronc/mcp"
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

// Hints appended to a not_found error, naming the tool that lists valid ids.
const (
	hintAccount = "Call courrier_list_accounts for the account ids you can use."
	hintMessage = "Call courrier_search_mail for the message ids of this account."
	hintDraft   = "Call courrier_search_mail with folder \"drafts\" for this account's draft ids, or create one with courrier_draft."
)

// errNoCaller is the refusal for a token that resolves to no Courrier user.
var errNoCaller = errors.New("this request carries no valid access token. Reconnect Courrier")

// toolset holds what the tool handlers share.
type toolset struct {
	mail    Mailbox
	confirm *mcpkit.Confirmations
	appURL  string
}

// caller is the person a tool call acts for.
type caller struct {
	userID int64
}

// callerOf reads the caller from the verified token of the request.
func callerOf(req *sdk.CallToolRequest) (caller, error) {
	identity, ok := mcpauth.Identity(req.Extra.TokenInfo)
	if !ok {
		return caller{}, errNoCaller
	}
	return caller{userID: identity.UserID}, nil
}
