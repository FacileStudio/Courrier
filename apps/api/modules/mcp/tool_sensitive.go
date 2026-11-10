package mcp

import (
	"context"
	"errors"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// statusPending is the status of an action waiting for the confirm page.
const statusPending = "pending_confirmation"

// errBoth is the refusal for a send that names a draft and a message.
var errBoth = errors.New("give draft_id or the message, not both. Change the draft with courrier_draft first")

// send is courrier_send.
func (t *toolset) send(ctx context.Context, req *sdk.CallToolRequest, in draftInput) (*sdk.CallToolResult, actionOutput, error) {
	who, err := writer(req)
	if err != nil {
		return nil, actionOutput{}, err
	}
	if retrying(req) {
		return t.resume(ctx, req, who, kindSend)
	}
	draftID, err := t.draftToSend(ctx, who, in)
	if err != nil {
		return nil, actionOutput{}, err
	}
	return t.request(ctx, req, who, kindSend, draftRef{AccountID: in.AccountID, DraftID: draftID})
}

// deleteDraft is courrier_delete_draft.
func (t *toolset) deleteDraft(ctx context.Context, req *sdk.CallToolRequest, in draftRef) (*sdk.CallToolResult, actionOutput, error) {
	who, err := writer(req)
	if err != nil {
		return nil, actionOutput{}, err
	}
	if retrying(req) {
		return t.resume(ctx, req, who, kindDeleteDraft)
	}
	return t.request(ctx, req, who, kindDeleteDraft, in)
}

// draftToSend is the draft a send is about: the one named, or a new one
// holding the message composed in the call.
func (t *toolset) draftToSend(ctx context.Context, who caller, in draftInput) (int64, error) {
	composed := len(in.To)+len(in.Cc)+len(in.Bcc) > 0 || in.Body != "" || in.Subject != "" || in.ReplyToMessageID != 0
	if in.DraftID != 0 && composed {
		return 0, errBoth
	}
	if in.DraftID != 0 {
		return in.DraftID, nil
	}
	if in.Body == "" || (len(in.To)+len(in.Cc)+len(in.Bcc) == 0 && in.ReplyToMessageID == 0) {
		return 0, errNoDraft
	}
	saved, err := t.saveDraft(ctx, who.userID, in)
	t.audit(ctx, who, event{tool: toolDraft, outcome: outcomeOf(err), accountID: in.AccountID, targetID: saved.ID})
	return saved.ID, err
}

// request opens a pending action for the draft, then asks the client to
// confirm it, or, when the client cannot ask, returns the confirm link.
// Either way nothing has happened yet.
func (t *toolset) request(ctx context.Context, req *sdk.CallToolRequest, who caller, kind string, ref draftRef) (*sdk.CallToolResult, actionOutput, error) {
	draft, err := t.mail.Draft(ctx, who.userID, ref.AccountID, ref.DraftID)
	if err != nil {
		return nil, actionOutput{}, failure(err, hintDraft)
	}
	action, err := t.pending.create(ctx, schemas.PendingAction{
		UserID: who.userID, AccountID: ref.AccountID, EmailID: ref.DraftID, Kind: kind, Digest: digest(draft), ClientID: who.clientID,
	})
	if err != nil {
		return nil, actionOutput{}, failure(err, "")
	}
	t.audit(ctx, who, event{tool: toolOf(kind), outcome: "requested", accountID: ref.AccountID, targetID: ref.DraftID})
	if canElicit(req) {
		return t.confirm.ask(who.userID, action.ID, question(kind, draft)), actionOutput{}, nil
	}
	return nil, actionOutput{
		Status: statusPending, DraftID: ref.DraftID, ConfirmURL: t.confirmURL(action.ID), Message: waiting(kind),
	}, nil
}

// resume handles the retry that carries the person's answer.
func (t *toolset) resume(ctx context.Context, req *sdk.CallToolRequest, who caller, kind string) (*sdk.CallToolResult, actionOutput, error) {
	id, answer, err := t.confirm.answer(req, who.userID)
	if err != nil {
		return nil, actionOutput{}, err
	}
	if answer == declined {
		_ = t.pending.resolve(ctx, who.userID, id, outcomeDeclined)
		t.audit(ctx, who, event{tool: toolOf(kind), outcome: outcomeDeclined})
		return nil, actionOutput{Status: outcomeDeclined, Message: "The user declined. Nothing happened, and the draft is kept."}, nil
	}
	action, err := t.run(ctx, who, id, kind)
	if err != nil {
		return nil, actionOutput{}, failure(err, "Call the tool again to ask the user afresh.")
	}
	return nil, actionOutput{Status: doneStatus(kind), DraftID: action.EmailID, Message: doneMessage(kind)}, nil
}

// confirmURL is the Courrier page where the person confirms action id.
func (t *toolset) confirmURL(id string) string {
	return t.appURL + "/pending?action=" + id
}
