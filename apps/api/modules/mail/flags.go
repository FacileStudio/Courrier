package mail

import (
	"context"

	"github.com/emersion/go-imap/v2"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"
)

// SetReadState marks many emails read or unread at once: one DB update,
// refreshed folder unread counts, and a single IMAP connection that issues one
// STORE per folder. Opening a conversation and bulk mark-read both use it, so
// a big thread no longer costs one IMAP login per message.
func (s *Service) SetReadState(ctx context.Context, userID, accountID int64, emailIDs []int64, isRead bool) error {
	return s.setFlagState(ctx, userID, accountID, emailIDs, flagChange{"is_read", imap.FlagSeen, isRead})
}

// SetStarState stars or unstars many emails with the same single DB update and
// single IMAP connection as SetReadState.
func (s *Service) SetStarState(ctx context.Context, userID, accountID int64, emailIDs []int64, starred bool) error {
	return s.setFlagState(ctx, userID, accountID, emailIDs, flagChange{"is_starred", imap.FlagFlagged, starred})
}

type flagChange struct {
	column string
	flag   imap.Flag
	on     bool
}

// setFlagState writes one boolean column for the given emails and mirrors it
// to the IMAP flag in the background. Rows with no IMAP uid are skipped on the
// server side: they are locally created draft or sent rows with no message to
// flag.
func (s *Service) setFlagState(ctx context.Context, userID, accountID int64, emailIDs []int64, change flagChange) error {
	account, err := s.getAccountDecrypted(ctx, userID, accountID)
	if err != nil {
		return err
	}
	if len(emailIDs) == 0 {
		return nil
	}

	var emails []schemas.Email
	if err := s.orm.WithContext(ctx).Select("id", "folder_id", "imap_uid").
		Where("id IN ? AND account_id = ?", emailIDs, accountID).Find(&emails).Error; err != nil {
		return errors.Internal("failed to load emails", err)
	}
	if len(emails) == 0 {
		return nil
	}
	s.orm.WithContext(ctx).Model(&schemas.Email{}).
		Where("id IN ? AND account_id = ?", emailIDs, accountID).Update(change.column, change.on)

	go s.storeFlagOnServer(account, s.serverUIDs(ctx, emails, change.column), change)
	return nil
}

func (s *Service) serverUIDs(ctx context.Context, emails []schemas.Email, column string) map[int64][]imap.UID {
	uidsByFolder := map[int64][]imap.UID{}
	refreshed := map[int64]bool{}
	for _, e := range emails {
		if column == "is_read" && !refreshed[e.FolderID] {
			refreshed[e.FolderID] = true
			s.refreshUnread(ctx, e.FolderID)
		}
		if e.IMAPUID != 0 {
			uidsByFolder[e.FolderID] = append(uidsByFolder[e.FolderID], imap.UID(e.IMAPUID))
		}
	}
	return uidsByFolder
}

func (s *Service) refreshUnread(ctx context.Context, folderID int64) {
	var unread int64
	s.orm.WithContext(ctx).Model(&schemas.Email{}).Where("folder_id = ? AND is_read = false", folderID).Count(&unread)
	s.orm.WithContext(ctx).Model(&schemas.Folder{}).Where("id = ?", folderID).Update("unread_count", unread)
}

func (s *Service) storeFlagOnServer(account schemas.Account, uidsByFolder map[int64][]imap.UID, change flagChange) {
	client, err := s.connect(account)
	if err != nil {
		return
	}
	defer closeIMAP(client)

	var folders []schemas.Folder
	s.orm.Where("account_id = ?", account.ID).Find(&folders)
	paths := map[int64]string{}
	for _, f := range folders {
		paths[f.ID] = f.Path
	}

	op := imap.StoreFlagsAdd
	if !change.on {
		op = imap.StoreFlagsDel
	}
	for folderID, uids := range uidsByFolder {
		if path := paths[folderID]; path != "" {
			storeFlagsBulk(client, path, uids, op, []imap.Flag{change.flag})
		}
	}
}
