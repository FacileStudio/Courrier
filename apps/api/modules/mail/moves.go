package mail

import (
	"context"

	"github.com/emersion/go-imap/v2/imapclient"
	"gorm.io/gorm"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"
)

// MoveEmails moves emails to the account's folder of the given type. Delete
// and archive are moves to trash and archive.
func (s *Service) MoveEmails(ctx context.Context, userID, accountID int64, emailIDs []int64, destFolderType string) error {
	return s.MoveEmailsTo(ctx, userID, accountID, emailIDs, MoveEmailsRequest{FolderType: destFolderType})
}

// MoveEmailsTo moves emails to the folder the request names, by id or by type.
// A type must name a system folder, because "custom" matches many.
func (s *Service) MoveEmailsTo(ctx context.Context, userID, accountID int64, emailIDs []int64, req MoveEmailsRequest) error {
	account, err := s.getAccountDecrypted(ctx, userID, accountID)
	if err != nil {
		return err
	}
	dest, err := s.resolveDestination(ctx, accountID, req)
	if err != nil {
		return err
	}

	var emails []schemas.Email
	if err := s.orm.WithContext(ctx).Where("id IN ? AND account_id = ?", emailIDs, accountID).Find(&emails).Error; err != nil {
		return errors.Internal("failed to find emails", err)
	}
	if len(emails) == 0 {
		return errors.NotFound("no emails found")
	}

	bySource := map[int64][]schemas.Email{}
	for _, e := range emails {
		if e.FolderID != dest.ID {
			bySource[e.FolderID] = append(bySource[e.FolderID], e)
		}
	}
	return s.moveGroups(ctx, account, bySource, dest)
}

func (s *Service) resolveDestination(ctx context.Context, accountID int64, req MoveEmailsRequest) (schemas.Folder, error) {
	var dest schemas.Folder
	query := s.orm.WithContext(ctx).Where("account_id = ?", accountID)
	switch {
	case req.FolderID != nil:
		query = query.Where("id = ?", *req.FolderID)
	case req.FolderType == "" || req.FolderType == schemas.FolderTypeCustom:
		return dest, errors.Invalid("name the destination with folder_id, or a system folder_type")
	default:
		query = query.Where("type = ?", req.FolderType)
	}
	if err := query.First(&dest).Error; err != nil {
		return dest, errors.NotFound("destination folder not found")
	}
	return dest, nil
}

func (s *Service) moveGroups(ctx context.Context, account schemas.Account, bySource map[int64][]schemas.Email, dest schemas.Folder) error {
	var client *imapclient.Client
	defer func() {
		if client != nil {
			closeIMAP(client)
		}
	}()

	for sourceID, group := range bySource {
		var source schemas.Folder
		if s.orm.WithContext(ctx).Where("id = ?", sourceID).First(&source).Error != nil {
			continue
		}
		newUIDs, err := s.moveOnServer(account, &client, source, dest, group)
		if err != nil {
			return errors.Failed(err.Error())
		}
		s.applyMove(ctx, group, source, dest, newUIDs)
	}
	return nil
}

func (s *Service) moveOnServer(account schemas.Account, client **imapclient.Client, source, dest schemas.Folder, group []schemas.Email) (map[uint32]uint32, error) {
	if !hasServerCopies(group) {
		return map[uint32]uint32{}, nil
	}
	if *client == nil {
		c, err := s.connect(account)
		if err != nil {
			return nil, err
		}
		*client = c
	}
	return moveUIDs(*client, source.Path, uidsOf(group), dest.Path)
}

func (s *Service) applyMove(ctx context.Context, group []schemas.Email, source, dest schemas.Folder, newUIDs map[uint32]uint32) {
	for _, e := range group {
		updates := map[string]any{"folder_id": dest.ID}
		if e.IMAPUID != 0 {
			updates["imap_uid"] = newUIDs[e.IMAPUID]
		}
		s.orm.WithContext(ctx).Model(&schemas.Email{}).Where("id = ?", e.ID).Updates(updates)
	}
	n := len(group)
	s.orm.WithContext(ctx).Model(&schemas.Folder{}).Where("id = ?", source.ID).
		UpdateColumn("total_count", gorm.Expr("GREATEST(total_count - ?, 0)", n))
	s.orm.WithContext(ctx).Model(&schemas.Folder{}).Where("id = ?", dest.ID).
		UpdateColumn("total_count", gorm.Expr("total_count + ?", n))
	s.refreshUnread(ctx, source.ID)
	s.refreshUnread(ctx, dest.ID)
}

// purgeFromTrash permanently deletes emails already in Trash: the server copy
// first, so a later sync cannot bring them back, then the cached rows.
func (s *Service) purgeFromTrash(ctx context.Context, account schemas.Account, trash schemas.Folder, group []schemas.Email) error {
	if len(group) == 0 {
		return nil
	}
	if err := s.expungeFromServer(account, trash, group); err != nil {
		return errors.Failed(err.Error())
	}
	ids := make([]int64, len(group))
	for i, e := range group {
		ids[i] = e.ID
	}
	s.orm.WithContext(ctx).Where("id IN ?", ids).Delete(&schemas.Email{})
	s.orm.WithContext(ctx).Where("email_id IN ?", ids).Delete(&schemas.Attachment{})
	return nil
}
