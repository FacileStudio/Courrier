package mcp

import (
	"context"
	"time"

	mcpkit "github.com/FacileStudio/tronc/mcp"

	"github.com/FacileStudio/Courrier/apps/api/modules/mail"
	"github.com/FacileStudio/Courrier/apps/api/schemas"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Limits of courrier_read_thread.
const (
	threadPage          = 10
	conciseBodyLength   = 2_000
	detailedBodyLength  = 20_000
	threadTruncatedNote = "Some bodies were shortened. Ask for response_format \"detailed\" to read them in full."
)

// threadPosition is the cursor of courrier_read_thread: how many messages,
// oldest first, earlier pages returned.
type threadPosition struct {
	Offset int `json:"o"`
}

// readThread is courrier_read_thread.
func (t *toolset) readThread(ctx context.Context, req *sdk.CallToolRequest, in threadInput) (*sdk.CallToolResult, threadOutput, error) {
	who, err := callerOf(req)
	if err != nil {
		return nil, threadOutput{}, err
	}
	full, err := mcpkit.Detailed(in.ResponseFormat)
	if err != nil {
		return nil, threadOutput{}, err
	}
	var position threadPosition
	if err := mcpkit.DecodeCursor(in.Cursor, &position); err != nil || position.Offset < 0 {
		return nil, threadOutput{}, mcpkit.ErrCursor
	}
	emails, err := t.conversation(ctx, who.userID, in.AccountID, in.MessageID)
	if err != nil {
		return nil, threadOutput{}, err
	}
	if position.Offset >= len(emails) {
		return nil, threadOutput{}, mcpkit.ErrCursor
	}
	return nil, t.threadPage(ctx, who.userID, emails, position.Offset, full), nil
}

// conversation loads every message of the conversation messageID is in,
// oldest first. A message without a thread is a conversation of one.
func (t *toolset) conversation(ctx context.Context, userID, accountID, messageID int64) ([]schemas.Email, error) {
	email, err := t.mail.GetEmail(ctx, userID, accountID, messageID)
	if err != nil {
		return nil, mcpkit.Failure(err, hintMessage)
	}
	if email.ThreadID == "" {
		return []schemas.Email{email}, nil
	}
	emails, err := t.mail.GetThread(ctx, userID, accountID, email.ThreadID)
	if err != nil || len(emails) == 0 {
		return []schemas.Email{email}, nil
	}
	return emails, nil
}

// threadPage renders the page of the conversation that starts at offset,
// cut to the size budget, with the cursor of the next page.
func (t *toolset) threadPage(ctx context.Context, userID int64, emails []schemas.Email, offset int, full bool) threadOutput {
	page := emails[offset:min(offset+threadPage, len(emails))]
	out := threadOutput{Subject: emails[0].Subject, MessageCount: len(emails)}
	shortened := false
	messages := make([]threadMessage, 0, len(page))
	for _, email := range page {
		message := t.threadMessage(ctx, userID, email, full)
		shortened = shortened || message.BodyTruncated
		messages = append(messages, message)
	}
	keep := max(mcpkit.FitCount(len(messages), func(n int) int {
		return mcpkit.EncodedSize(threadOutput{Subject: out.Subject, Messages: messages[:n], NextCursor: "next-cursor-placeholder"})
	}), 1)
	out.Messages = messages[:keep]
	if offset+keep < len(emails) {
		out.NextCursor = mcpkit.EncodeCursor(threadPosition{Offset: offset + keep})
	}
	if shortened && !full {
		out.Note = threadTruncatedNote
	}
	return out
}

// threadMessage renders one message, loading its body from the server when
// the cache has none.
func (t *toolset) threadMessage(ctx context.Context, userID int64, email schemas.Email, full bool) threadMessage {
	if email.BodyText == "" && email.BodyHTML == "" {
		if loaded, err := t.mail.GetEmail(ctx, userID, email.AccountID, email.ID); err == nil {
			email = loaded
		}
	}
	body, limit := mail.PlainBody(email), detailedBodyLength
	if !full {
		if trimmed := withoutQuotes(body); trimmed != "" {
			body = trimmed
		}
		limit = conciseBodyLength
	}
	body, cut := mcpkit.Clip(body, limit)
	return threadMessage{
		MessageID: email.ID, From: sender(email), To: mail.AddressList(email.ToAddresses),
		Cc: mail.AddressList(email.CcAddresses), Date: email.Date.UTC().Format(time.RFC3339),
		Body: body, BodyTruncated: cut, Unread: !email.IsRead, HasAttachments: email.HasAttachments,
	}
}
