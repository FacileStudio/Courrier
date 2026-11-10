package mail

import (
	"context"
	"strings"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"
)

// Draft returns one draft of an account the person can see.
func (s *Service) Draft(ctx context.Context, userID, accountID, draftID int64) (schemas.Email, error) {
	if _, err := s.getAccount(ctx, userID, accountID); err != nil {
		return schemas.Email{}, err
	}
	_, draft, err := s.findDraft(ctx, accountID, draftID)
	return draft, err
}

// SendDraft sends a stored draft as it is, then removes it from Drafts. A
// draft with attachments is refused, because their bytes live on the IMAP
// server and this path does not fetch them; the person sends it from the
// composer. Removing the draft is best effort once the message is out: an
// error there must not invite a second send.
func (s *Service) SendDraft(ctx context.Context, userID, accountID, draftID int64) error {
	draft, err := s.Draft(ctx, userID, accountID, draftID)
	if err != nil {
		return err
	}
	if draft.HasAttachments {
		return errors.Failed("this draft has attachments; open it in Courrier to send it")
	}
	if err := s.Send(ctx, userID, accountID, DraftMessage(draft)); err != nil {
		return err
	}
	_ = s.DeleteDraft(ctx, userID, accountID, draftID)
	return nil
}

// DraftMessage turns a stored draft back into the request that sends it.
func DraftMessage(draft schemas.Email) SendRequest {
	return SendRequest{
		To:         AddressList(draft.ToAddresses),
		Cc:         AddressList(draft.CcAddresses),
		Bcc:        AddressList(draft.BccAddresses),
		Subject:    draft.Subject,
		Body:       draft.BodyText,
		BodyHTML:   draft.BodyHTML,
		InReplyTo:  draft.InReplyTo,
		References: strings.Fields(draft.References),
	}
}

// AddressList returns the addresses stored in an address column, in order,
// as "Name <email>" when a name is known.
func AddressList(raw string) []string {
	parsed := parseAddressJSON(raw)
	out := make([]string, 0, len(parsed))
	for _, address := range parsed {
		if address.Name == "" {
			out = append(out, address.Email)
			continue
		}
		out = append(out, address.Name+" <"+address.Email+">")
	}
	return out
}

// PlainBody is the readable text of a message: the text part, or the HTML
// part with its tags stripped when there is no text part.
func PlainBody(email schemas.Email) string {
	if strings.TrimSpace(email.BodyText) != "" {
		return email.BodyText
	}
	return strings.TrimSpace(stripHTMLTags(email.BodyHTML))
}
