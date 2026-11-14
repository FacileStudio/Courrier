package mcp

import (
	"context"
	"net/http"
	"strconv"

	"github.com/FacileStudio/Courrier/apps/api/modules/mail"
	"github.com/FacileStudio/porte/mcpauth"
	mcpkit "github.com/FacileStudio/tronc/mcp"
)

// webActor marks, in the audit, a confirmation made on the confirm page.
const webActor = "courrier-web"

// pendingUser is the signed-in person of a confirm page request. Only the
// person's own browser session confirms: a bearer token, a CLI or API token
// and a delegated token are refused, so an agent never confirms its own
// request.
func pendingUser(r *http.Request) (string, error) {
	identity, err := mcpauth.RequireOwnSession(r)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(identity.UserID, 10), nil
}

// describe adds an open action's mailbox and its draft as it is now.
func (t *toolset) describe(ctx context.Context, item mcpkit.PendingItem, a mcpkit.Action) PendingItem {
	userID, ref, _ := refOf(a)
	out := PendingItem{PendingItem: item, AccountID: ref.AccountID, Draft: t.pendingDraft(ctx, a)}
	accounts, err := t.mail.VisibleAccounts(ctx, userID)
	if err != nil {
		return out
	}
	for _, account := range accounts {
		if account.ID == ref.AccountID {
			out.AccountEmail = account.Email
		}
	}
	return out
}

// done is the status a confirmed action answers on the confirm page.
func done(_ context.Context, o mcpkit.Outcome) mcpkit.PendingResult {
	return mcpkit.PendingResult{Status: doneStatus(o.Action.Kind)}
}

// pendingDraft loads an action's draft as it is now, or nil when it is gone.
func (t *toolset) pendingDraft(ctx context.Context, action mcpkit.Action) *PendingDraft {
	draft, err := t.loadDraft(ctx, action)
	if err != nil {
		return nil
	}
	body, _ := mcpkit.Clip(mail.PlainBody(draft), detailedBodyLength)
	return &PendingDraft{
		ID: draft.ID, To: mail.AddressList(draft.ToAddresses), Cc: mail.AddressList(draft.CcAddresses),
		Bcc: mail.AddressList(draft.BccAddresses), Subject: draft.Subject, Body: body,
	}
}
