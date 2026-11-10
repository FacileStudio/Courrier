package mcp

import (
	"context"
	"net/http"
	"strconv"

	"github.com/FacileStudio/Courrier/apps/api/internal/authcontext"
	"github.com/FacileStudio/Courrier/apps/api/modules/mail"
	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"
	"github.com/go-chi/chi/v5"
)

// webActor marks, in the audit, a confirmation made on the confirm page.
const webActor = "courrier-web"

// pendingUser is the signed-in person of a confirm page request. A delegated
// token is refused: an agent must never confirm its own request.
func pendingUser(r *http.Request) (int64, error) {
	if identity, ok := porte.From(r.Context()); ok && identity.Actor != "" {
		return 0, errors.Forbidden("a delegated token cannot confirm actions")
	}
	userID, err := strconv.ParseInt(authcontext.MustIdentity(r.Context()).UserID, 10, 64)
	if err != nil {
		return 0, errors.Invalid("invalid user id")
	}
	return userID, nil
}

// handleList answers the person's open actions with their drafts.
func (s *Service) handleList(w http.ResponseWriter, r *http.Request) {
	userID, err := pendingUser(r)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	list, err := s.tools.pendingList(r.Context(), userID)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, list)
}

// handleConfirm runs one action.
func (s *Service) handleConfirm(w http.ResponseWriter, r *http.Request) {
	userID, err := pendingUser(r)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	action, err := s.tools.run(r.Context(), caller{userID: userID, actor: webActor}, chi.URLParam(r, "id"), "")
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, PendingResult{Status: doneStatus(action.Kind)})
}

// handleDismiss closes one action without running it. The draft stays.
func (s *Service) handleDismiss(w http.ResponseWriter, r *http.Request) {
	userID, err := pendingUser(r)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	if err := s.tools.pending.resolve(r.Context(), userID, chi.URLParam(r, "id"), outcomeDismissed); err != nil {
		httpjson.WriteError(w, err)
		return
	}
	s.tools.audit(r.Context(), caller{userID: userID, actor: webActor}, event{tool: "pending", outcome: outcomeDismissed})
	httpjson.WriteJSON(w, http.StatusOK, PendingResult{Status: outcomeDismissed})
}

// pendingList renders userID's open actions with the current content of
// their drafts.
func (t *toolset) pendingList(ctx context.Context, userID int64) (PendingList, error) {
	actions, err := t.pending.open(ctx, userID)
	if err != nil {
		return PendingList{}, err
	}
	accounts, err := t.mail.VisibleAccounts(ctx, userID)
	if err != nil {
		return PendingList{}, err
	}
	emails := map[int64]string{}
	for _, account := range accounts {
		emails[account.ID] = account.Email
	}
	list := PendingList{Actions: make([]PendingItem, 0, len(actions))}
	for _, action := range actions {
		item := PendingItem{
			ID: action.ID, Kind: action.Kind, ClientID: action.ClientID, CreatedAt: action.CreatedAt,
			ExpiresAt: action.ExpiresAt, AccountID: action.AccountID, AccountEmail: emails[action.AccountID],
		}
		item.State, item.Draft = t.pendingDraft(ctx, userID, action)
		list.Actions = append(list.Actions, item)
	}
	return list, nil
}

// pendingDraft loads an action's draft and says whether it still holds what
// the person was asked about.
func (t *toolset) pendingDraft(ctx context.Context, userID int64, action schemas.PendingAction) (string, *PendingDraft) {
	draft, err := t.mail.Draft(ctx, userID, action.AccountID, action.EmailID)
	if err != nil {
		return pendingMissing, nil
	}
	body, _ := clip(mail.PlainBody(draft), detailedBodyLength)
	view := &PendingDraft{
		ID: draft.ID, To: mail.AddressList(draft.ToAddresses), Cc: mail.AddressList(draft.CcAddresses),
		Bcc: mail.AddressList(draft.BccAddresses), Subject: draft.Subject, Body: body,
	}
	if digest(draft) != action.Digest {
		return pendingChanged, view
	}
	return pendingReady, view
}
