package mcp

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/FacileStudio/Courrier/apps/api/modules/mail"
	"github.com/FacileStudio/Courrier/apps/api/schemas"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// maxRecipients caps to, cc and bcc together.
const maxRecipients = 50

// draft is courrier_draft.
func (t *toolset) draft(ctx context.Context, req *sdk.CallToolRequest, in draftInput) (*sdk.CallToolResult, draftOutput, error) {
	who, err := writer(req)
	if err != nil {
		return nil, draftOutput{}, err
	}
	saved, err := t.saveDraft(ctx, who.userID, in)
	t.audit(ctx, who, event{tool: toolDraft, outcome: outcomeOf(err), accountID: in.AccountID, targetID: saved.ID})
	if err != nil {
		return nil, draftOutput{}, err
	}
	return nil, draftOutput{
		DraftID: saved.ID, AccountID: in.AccountID, To: mail.AddressList(saved.ToAddresses), Subject: saved.Subject,
		ReviewURL: t.draftURL(in.AccountID, saved.ID),
		Next:      "Nothing was sent. Call courrier_send with this draft_id to send it, after the user agrees.",
	}, nil
}

// saveDraft creates the draft in, or replaces draft in.DraftID with it.
func (t *toolset) saveDraft(ctx context.Context, userID int64, in draftInput) (schemas.Email, error) {
	message, err := t.compose(ctx, userID, in)
	if err != nil {
		return schemas.Email{}, err
	}
	if in.DraftID == 0 {
		saved, err := t.mail.SaveDraft(ctx, userID, in.AccountID, message)
		return saved, failure(err, hintAccount)
	}
	saved, err := t.mail.UpdateDraft(ctx, userID, in.AccountID, in.DraftID, mail.DraftRequest{SendRequest: message})
	return saved, failure(err, hintDraft)
}

// compose builds the message of a draft input, threading it under the
// message it replies to.
func (t *toolset) compose(ctx context.Context, userID int64, in draftInput) (mail.SendRequest, error) {
	message := mail.SendRequest{To: in.To, Cc: in.Cc, Bcc: in.Bcc, Subject: in.Subject, Body: in.Body}
	if len(in.To)+len(in.Cc)+len(in.Bcc) > maxRecipients {
		return message, fmt.Errorf("a message takes at most %d recipients", maxRecipients)
	}
	if in.ReplyToMessageID == 0 {
		return message, nil
	}
	original, err := t.mail.GetEmail(ctx, userID, in.AccountID, in.ReplyToMessageID)
	if err != nil {
		return message, failure(err, hintMessage)
	}
	return replyTo(message, original), nil
}

// replyTo threads message under original, and fills its recipient and
// subject when they are empty.
func replyTo(message mail.SendRequest, original schemas.Email) mail.SendRequest {
	message.InReplyTo = original.MessageID
	message.References = append(strings.Fields(original.References), original.MessageID)
	if len(message.To) == 0 {
		message.To = []string{original.FromAddress}
	}
	if message.Subject == "" {
		message.Subject = original.Subject
		if !strings.HasPrefix(strings.ToLower(original.Subject), "re:") {
			message.Subject = "Re: " + original.Subject
		}
	}
	return message
}

// draftURL opens a draft in the Courrier composer.
func (t *toolset) draftURL(accountID, draftID int64) string {
	query := url.Values{"accountId": {fmt.Sprint(accountID)}, "draft": {fmt.Sprint(draftID)}}
	return t.appURL + "/mail/compose?" + query.Encode()
}

// errNoDraft is the refusal for a send with neither a draft nor a message.
var errNoDraft = errors.New("give draft_id, or to and body to compose the message here")
