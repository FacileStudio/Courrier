package mail

import (
	"context"

	"github.com/emersion/go-imap/v2"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"
)

// SaveDraft stores a new draft in the Drafts folder, locally and on the IMAP
// server. A draft with attachments must reach the server, because the server
// copy is where the attachment bytes live; one without keeps the local row
// even when the server is unreachable.
func (s *Service) SaveDraft(ctx context.Context, userID, accountID int64, req SendRequest) (schemas.Email, error) {
	account, err := s.getAccountDecrypted(ctx, userID, accountID)
	if err != nil {
		return schemas.Email{}, err
	}
	folder, err := s.draftsFolder(ctx, accountID)
	if err != nil {
		return schemas.Email{}, errors.NotFound("drafts folder not found, sync the account first")
	}
	return s.storeDraft(ctx, draftWrite{account: account, folder: folder, req: DraftRequest{SendRequest: req}})
}

// UpdateDraft replaces the content of a draft in place. The row keeps its id
// and its Message-ID; the server copy is replaced by a fresh one. Attachments
// stay unless listed in RemoveAttachments, and uploads are added.
func (s *Service) UpdateDraft(ctx context.Context, userID, accountID, emailID int64, req DraftRequest) (schemas.Email, error) {
	account, err := s.getAccountDecrypted(ctx, userID, accountID)
	if err != nil {
		return schemas.Email{}, err
	}
	folder, existing, err := s.findDraft(ctx, accountID, emailID)
	if err != nil {
		return schemas.Email{}, err
	}
	return s.storeDraft(ctx, draftWrite{account: account, folder: folder, existing: &existing, req: req})
}

// DeleteDraft removes a draft from the cache and from the server's Drafts
// mailbox, so a later sync does not bring it back.
func (s *Service) DeleteDraft(ctx context.Context, userID, accountID, emailID int64) error {
	account, err := s.getAccountDecrypted(ctx, userID, accountID)
	if err != nil {
		return err
	}
	folder, email, err := s.findDraft(ctx, accountID, emailID)
	if err != nil {
		return err
	}
	if email.IMAPUID != 0 {
		client, err := s.connect(account)
		if err != nil {
			return errors.Failed(err.Error())
		}
		defer closeIMAP(client)
		if err := expungeUIDs(client, folder.Path, []imap.UID{imap.UID(email.IMAPUID)}); err != nil {
			return errors.Failed(err.Error())
		}
	}
	s.orm.WithContext(ctx).Where("email_id = ?", email.ID).Delete(&schemas.Attachment{})
	s.orm.WithContext(ctx).Delete(&email)
	return nil
}

func (s *Service) draftsFolder(ctx context.Context, accountID int64) (schemas.Folder, error) {
	var folder schemas.Folder
	err := s.orm.WithContext(ctx).Where("account_id = ? AND type = ?", accountID, schemas.FolderTypeDrafts).First(&folder).Error
	return folder, err
}

func (s *Service) findDraft(ctx context.Context, accountID, emailID int64) (schemas.Folder, schemas.Email, error) {
	folder, err := s.draftsFolder(ctx, accountID)
	if err != nil {
		return folder, schemas.Email{}, errors.NotFound("drafts folder not found")
	}
	var email schemas.Email
	if err := s.orm.WithContext(ctx).Where("id = ? AND account_id = ? AND folder_id = ?", emailID, accountID, folder.ID).First(&email).Error; err != nil {
		return folder, email, errors.NotFound("draft not found")
	}
	return folder, email, nil
}
