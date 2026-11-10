package mcp

import (
	"context"
	"errors"

	mcpkit "github.com/FacileStudio/tronc/mcp"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// errBoth is the refusal for a send that names a draft and a message.
var errBoth = errors.New("give draft_id or the message, not both. Change the draft with courrier_draft first")

// send is courrier_send.
func (t *toolset) send(ctx context.Context, req *sdk.CallToolRequest, in draftInput) (*sdk.CallToolResult, actionOutput, error) {
	if mcpkit.Retrying(req) {
		return t.resume(ctx, req, kindSend)
	}
	who, err := callerOf(req)
	if err != nil {
		return nil, actionOutput{}, err
	}
	draftID, err := t.draftToSend(ctx, who, in)
	if err != nil {
		return nil, actionOutput{}, err
	}
	return t.request(ctx, req, who, kindSend, draftRef{AccountID: in.AccountID, DraftID: draftID})
}

// deleteDraft is courrier_delete_draft.
func (t *toolset) deleteDraft(ctx context.Context, req *sdk.CallToolRequest, in draftRef) (*sdk.CallToolResult, actionOutput, error) {
	if mcpkit.Retrying(req) {
		return t.resume(ctx, req, kindDeleteDraft)
	}
	who, err := callerOf(req)
	if err != nil {
		return nil, actionOutput{}, err
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
	return saved.ID, err
}

// request asks the person to confirm an action on the draft, in the client
// or through the confirm link. Either way nothing has happened yet.
func (t *toolset) request(ctx context.Context, req *sdk.CallToolRequest, who caller, kind string, ref draftRef) (*sdk.CallToolResult, actionOutput, error) {
	draft, err := t.mail.Draft(ctx, who.userID, ref.AccountID, ref.DraftID)
	if err != nil {
		return nil, actionOutput{}, mcpkit.Failure(err, hintDraft)
	}
	result, outcome, err := t.confirm.Request(ctx, req, kind, target(ref), question(kind, draft))
	if err != nil {
		return nil, actionOutput{}, mcpkit.Failure(err, hintDraft)
	}
	if result != nil {
		return result, actionOutput{}, nil
	}
	return nil, actionOutput{Status: outcome.Status, DraftID: ref.DraftID, ConfirmURL: outcome.ConfirmURL, Message: waiting(kind)}, nil
}

// resume handles the retry that carries the person's answer.
func (t *toolset) resume(ctx context.Context, req *sdk.CallToolRequest, kind string) (*sdk.CallToolResult, actionOutput, error) {
	outcome, err := t.confirm.Resume(ctx, req, kind)
	if err != nil {
		return nil, actionOutput{}, mcpkit.Failure(err, "Call the tool again to ask the user afresh.")
	}
	if outcome.Status == mcpkit.StatusDeclined {
		return nil, actionOutput{Status: mcpkit.StatusDeclined, Message: "The user declined. Nothing happened, and the draft is kept."}, nil
	}
	_, ref, _ := refOf(outcome.Action)
	return nil, actionOutput{Status: doneStatus(kind), DraftID: ref.DraftID, Message: doneMessage(kind)}, nil
}
