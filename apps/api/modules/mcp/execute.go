package mcp

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	stderrors "errors"
	"fmt"
	"strings"

	"github.com/FacileStudio/Courrier/apps/api/modules/mail"
	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"
)

// outcomeStale resolves an action whose draft changed after the request.
const outcomeStale = "stale"

// errChanged is the refusal for a draft edited after the confirmation was
// asked: the person confirmed other content.
var errChanged = errors.Failed("the draft changed after the user was asked, so nothing happened. Call the tool again to ask about the new content")

// errWrongKind is the refusal for a confirmation replayed to another tool.
var errWrongKind = errors.Failed(errConfirmation.Error())

// run claims one of who's pending actions, runs it and records how it ended.
// A wrong kind or a failed run reopens it, a changed draft closes it. Errors
// keep their tronc codes, for the confirm page; tools pass them to failure.
func (t *toolset) run(ctx context.Context, who caller, id, kind string) (schemas.PendingAction, error) {
	action, err := t.pending.claim(ctx, who.userID, id)
	if err != nil {
		return action, err
	}
	if kind != "" && action.Kind != kind {
		t.pending.release(ctx, id)
		return action, errWrongKind
	}
	who.clientID = action.ClientID
	err = t.execute(ctx, who.userID, action)
	switch {
	case stderrors.Is(err, errChanged):
		t.pending.finish(ctx, id, outcomeStale)
	case err != nil:
		t.pending.release(ctx, id)
	default:
		t.pending.finish(ctx, id, outcomeDone)
	}
	t.audit(ctx, who, event{tool: toolOf(action.Kind), outcome: outcomeOf(err), accountID: action.AccountID, targetID: action.EmailID})
	return action, err
}

// execute sends or deletes the draft of action, if it still holds what the
// person was asked about.
func (t *toolset) execute(ctx context.Context, userID int64, action schemas.PendingAction) error {
	draft, err := t.mail.Draft(ctx, userID, action.AccountID, action.EmailID)
	if err != nil {
		return err
	}
	if digest(draft) != action.Digest {
		return errChanged
	}
	if action.Kind == kindSend {
		return t.mail.SendDraft(ctx, userID, action.AccountID, action.EmailID)
	}
	return t.mail.DeleteDraft(ctx, userID, action.AccountID, action.EmailID)
}

// digest fingerprints what a draft would send.
func digest(draft schemas.Email) string {
	fields := []string{
		draft.ToAddresses, draft.CcAddresses, draft.BccAddresses, draft.Subject,
		draft.BodyText, draft.BodyHTML, draft.InReplyTo, fmt.Sprint(draft.HasAttachments),
	}
	sum := sha256.Sum256([]byte(strings.Join(fields, "\x00")))
	return hex.EncodeToString(sum[:16])
}

// question is what the client asks the person.
func question(kind string, draft schemas.Email) string {
	recipients := strings.Join(append(append(mail.AddressList(draft.ToAddresses), mail.AddressList(draft.CcAddresses)...),
		mail.AddressList(draft.BccAddresses)...), ", ")
	if kind == kindSend {
		return fmt.Sprintf("Send the email %q to %s now?", draft.Subject, recipients)
	}
	return fmt.Sprintf("Delete the draft %q (to %s) for good?", draft.Subject, recipients)
}
