package mail

import (
	"context"
	"strings"
	"time"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"gorm.io/gorm"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"
)

type draftWrite struct {
	account  schemas.Account
	folder   schemas.Folder
	existing *schemas.Email
	req      DraftRequest
}

type draftFiled struct {
	messageID string
	uid       uint32
	rows      []schemas.Attachment
}

func (w draftWrite) onServer() bool {
	return w.existing != nil && w.existing.IMAPUID != 0
}

func (s *Service) storeDraft(ctx context.Context, w draftWrite) (schemas.Email, error) {
	client, connErr := s.connect(w.account)
	if connErr == nil {
		defer closeIMAP(client)
	}
	uploads, err := s.draftUploads(ctx, client, connErr, w)
	if err != nil {
		return schemas.Email{}, err
	}

	filed, err := s.fileDraft(client, connErr == nil, w, uploads)
	if err != nil && (len(uploads) > 0 || w.onServer()) {
		if filed.uid != 0 {
			expungeUIDs(client, w.folder.Path, []imap.UID{imap.UID(filed.uid)})
		}
		return schemas.Email{}, errors.Failed(err.Error())
	}
	email, err := s.persistDraft(ctx, w, filed)
	if err == nil && w.onServer() && filed.uid != 0 {
		expungeUIDs(client, w.folder.Path, []imap.UID{imap.UID(w.existing.IMAPUID)})
	}
	return email, err
}

func (s *Service) draftUploads(ctx context.Context, client *imapclient.Client, connErr error, w draftWrite) ([]AttachmentUpload, error) {
	uploads := w.req.Attachments
	if w.onServer() {
		kept, err := s.keptAttachments(ctx, client, w)
		if err != nil {
			return nil, errors.Failed(err.Error())
		}
		uploads = append(kept, uploads...)
	}
	if connErr != nil && (len(uploads) > 0 || w.onServer()) {
		return nil, errors.Failed(connErr.Error())
	}
	return uploads, nil
}

func (s *Service) fileDraft(client *imapclient.Client, connected bool, w draftWrite, uploads []AttachmentUpload) (draftFiled, error) {
	spec := specFor(w.account, w.req.SendRequest)
	spec.Attachments = uploads
	if w.existing != nil {
		spec.MessageID = w.existing.MessageID
	}
	msg, messageID, err := buildMessage(spec)
	if err != nil {
		return draftFiled{}, errors.Internal("failed to build draft message", err)
	}
	filed := draftFiled{messageID: messageID}
	if !connected {
		return filed, nil
	}
	if filed.uid, err = appendDraft(client, w.folder.Path, msg, messageID); err != nil || len(uploads) == 0 {
		return filed, err
	}
	filed.rows, err = fetchAttachmentRows(client, w.folder.Path, filed.uid, 0)
	return filed, err
}

func (s *Service) keptAttachments(ctx context.Context, client *imapclient.Client, w draftWrite) ([]AttachmentUpload, error) {
	var rows []schemas.Attachment
	query := s.orm.WithContext(ctx).Where("email_id = ?", w.existing.ID)
	if len(w.req.RemoveAttachments) > 0 {
		query = query.Where("id NOT IN ?", w.req.RemoveAttachments)
	}
	query.Find(&rows)
	if len(rows) == 0 {
		return nil, nil
	}
	if client == nil {
		return nil, errors.Failed("cannot reach the mail server to carry the attachments over")
	}
	var kept []AttachmentUpload
	for _, row := range rows {
		parts, err := parsePartID(row.PartID)
		if err != nil {
			return nil, err
		}
		data, err := fetchAttachmentPart(client, w.folder.Path, imap.UID(w.existing.IMAPUID), parts)
		if err != nil {
			return nil, err
		}
		kept = append(kept, AttachmentUpload{Filename: row.Filename, MimeType: row.MimeType, Data: data})
	}
	return kept, nil
}

func draftRow(w draftWrite, filed draftFiled) schemas.Email {
	req := w.req
	email := schemas.Email{
		AccountID:      w.account.ID,
		SpaceID:        w.account.SpaceID,
		FolderID:       w.folder.ID,
		MessageID:      filed.messageID,
		Subject:        req.Subject,
		FromAddress:    w.account.Email,
		FromName:       w.account.Name,
		ToAddresses:    stringsToAddressJSON(req.To),
		CcAddresses:    stringsToAddressJSON(req.Cc),
		BccAddresses:   stringsToAddressJSON(req.Bcc),
		Date:           time.Now(),
		BodyText:       req.Body,
		BodyHTML:       req.BodyHTML,
		InReplyTo:      req.InReplyTo,
		References:     strings.Join(req.References, " "),
		IsRead:         true,
		HasAttachments: len(filed.rows) > 0,
		IMAPUID:        filed.uid,
	}
	if old := w.existing; old != nil {
		email.ID, email.ThreadID, email.CachedAt, email.FacileID = old.ID, old.ThreadID, old.CachedAt, old.FacileID
	}
	return email
}

func (s *Service) persistDraft(ctx context.Context, w draftWrite, filed draftFiled) (schemas.Email, error) {
	email := draftRow(w, filed)
	err := s.orm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		old := w.existing
		if old == nil || old.InReplyTo != email.InReplyTo || old.References != email.References {
			email.ThreadID = assignThreadID(tx, w.account.ID, email.MessageID, email.InReplyTo, email.References)
		}
		return writeDraftRow(tx, &email, filed.rows)
	})
	if err != nil {
		return schemas.Email{}, errors.Internal("failed to save draft", err)
	}
	return email, nil
}

func writeDraftRow(tx *gorm.DB, email *schemas.Email, rows []schemas.Attachment) error {
	if email.ID != 0 {
		tx.Where("email_id = ?", email.ID).Delete(&schemas.Attachment{})
	}
	if err := tx.Save(email).Error; err != nil {
		return err
	}
	for i := range rows {
		rows[i].EmailID = email.ID
	}
	if len(rows) == 0 {
		return nil
	}
	return tx.Create(&rows).Error
}
