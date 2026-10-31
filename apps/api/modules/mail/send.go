package mail

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"

	"gorm.io/gorm"
)

// Send delivers a message over SMTP and files a copy in the Sent folder, both
// in the local cache and on the IMAP server. Bcc recipients receive the message
// through the SMTP envelope only: the delivered bytes carry no Bcc header, while
// the copies kept for the sender record it.
func (s *Service) Send(ctx context.Context, userID, accountID int64, req SendRequest) error {
	account, err := s.getAccountDecrypted(ctx, userID, accountID)
	if err != nil {
		return err
	}
	if len(req.To)+len(req.Cc)+len(req.Bcc) == 0 {
		return errors.Invalid("at least one recipient is required")
	}
	if req.Body == "" && req.BodyHTML == "" {
		return errors.Invalid("body is required")
	}

	msgBytes, messageID, err := buildMessage(specFor(account, req))
	if err != nil {
		return errors.Internal("failed to build message", err)
	}
	if err := s.deliver(account, envelopeRecipients(req), withoutBcc(msgBytes)); err != nil {
		return errors.Failed(fmt.Sprintf("failed to send: %s", err))
	}

	s.recordContacts(ctx, accountID, account.SpaceID, sendRecipients(req.To, req.Cc, req.Bcc), time.Now())
	s.recordSent(ctx, account, req, messageID)
	go s.appendToSent(account, msgBytes)
	return nil
}

func (s *Service) recordSent(ctx context.Context, account schemas.Account, req SendRequest, messageID string) {
	var sentFolder schemas.Folder
	if s.orm.WithContext(ctx).Where("account_id = ? AND type = ?", account.ID, schemas.FolderTypeSent).First(&sentFolder).Error != nil {
		return
	}
	references := strings.Join(req.References, " ")
	email := schemas.Email{
		AccountID:    account.ID,
		SpaceID:      account.SpaceID,
		FolderID:     sentFolder.ID,
		MessageID:    messageID,
		ThreadID:     assignThreadID(s.orm.WithContext(ctx), account.ID, messageID, req.InReplyTo, references),
		Subject:      req.Subject,
		FromAddress:  account.Email,
		FromName:     account.Name,
		ToAddresses:  stringsToAddressJSON(req.To),
		CcAddresses:  stringsToAddressJSON(req.Cc),
		BccAddresses: stringsToAddressJSON(req.Bcc),
		Date:         time.Now(),
		BodyText:     req.Body,
		BodyHTML:     req.BodyHTML,
		InReplyTo:    req.InReplyTo,
		References:   references,
		IsRead:       true,
	}
	s.orm.WithContext(ctx).Create(&email)
	s.orm.WithContext(ctx).Model(&sentFolder).UpdateColumn("total_count", gorm.Expr("total_count + 1"))
}
