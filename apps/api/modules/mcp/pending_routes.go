package mcp

import (
	"context"
	"net/http"
	"strconv"

	"github.com/FacileStudio/Courrier/apps/api/internal/middleware"
	"github.com/FacileStudio/Courrier/apps/api/modules/mail"
	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"
	mcpkit "github.com/FacileStudio/tronc/mcp"
	"github.com/go-chi/chi/v5"
)

// webActor marks, in the audit, a confirmation made on the confirm page.
const webActor = "courrier-web"

// pendingUser is the signed-in person of a confirm page request. A delegated
// token is refused: an agent must never confirm its own request.
func pendingUser(r *http.Request) (mcpkit.Caller, error) {
	if identity, ok := porte.From(r.Context()); ok && identity.Actor != "" {
		return mcpkit.Caller{}, errors.Forbidden("a delegated token cannot confirm actions")
	}
	return mcpkit.Caller{UserID: strconv.FormatInt(middleware.UserID(r), 10), Actor: webActor}, nil
}

// handleList answers the person's open actions with their drafts.
func (s *Service) handleList(w http.ResponseWriter, r *http.Request) {
	who, err := pendingUser(r)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	list, err := s.tools.pendingList(r.Context(), who.UserID)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, list)
}

// handleConfirm runs one action.
func (s *Service) handleConfirm(w http.ResponseWriter, r *http.Request) {
	who, err := pendingUser(r)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	outcome, err := s.tools.confirm.Confirm(r.Context(), who, chi.URLParam(r, "id"))
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, PendingResult{Status: doneStatus(outcome.Action.Kind)})
}

// handleDismiss closes one action without running it. The draft stays.
func (s *Service) handleDismiss(w http.ResponseWriter, r *http.Request) {
	who, err := pendingUser(r)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	if err := s.tools.confirm.Dismiss(r.Context(), who, chi.URLParam(r, "id")); err != nil {
		httpjson.WriteError(w, err)
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, PendingResult{Status: mcpkit.OutcomeDismissed})
}

// pendingList renders the person's open actions with the current content of
// their drafts.
func (t *toolset) pendingList(ctx context.Context, userID string) (PendingList, error) {
	actions, err := t.confirm.Open(ctx, userID)
	if err != nil {
		return PendingList{}, err
	}
	id, _ := strconv.ParseInt(userID, 10, 64)
	accounts, err := t.mail.VisibleAccounts(ctx, id)
	if err != nil {
		return PendingList{}, err
	}
	emails := map[int64]string{}
	for _, account := range accounts {
		emails[account.ID] = account.Email
	}
	list := PendingList{Actions: make([]PendingItem, 0, len(actions))}
	for _, action := range actions {
		_, ref, _ := refOf(action)
		item := PendingItem{
			ID: action.ID, Kind: action.Kind, ClientID: action.ClientID, CreatedAt: action.CreatedAt,
			ExpiresAt: action.ExpiresAt, AccountID: ref.AccountID, AccountEmail: emails[ref.AccountID],
			State: t.confirm.State(ctx, action), Draft: t.pendingDraft(ctx, action),
		}
		list.Actions = append(list.Actions, item)
	}
	return list, nil
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
