package mail

import (
	"context"
	"slices"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"
)

const backfillCooldown = 15 * time.Minute

func (s *Service) acquireBackfill(accountID int64) (func(), bool) {
	if _, loaded := s.backfilling.LoadOrStore(accountID, true); loaded {
		return nil, false
	}
	return func() { s.backfilling.Delete(accountID) }, true
}

func (s *Service) startCooldown(accountID int64) {
	s.cooldowns.Store(accountID, time.Now().Add(backfillCooldown))
}

func (s *Service) cooling(accountID int64) bool {
	until, ok := s.cooldowns.Load(accountID)
	return ok && time.Now().Before(until.(time.Time))
}

func (s *Service) visibleAccount(ctx context.Context, userID, accountID int64) error {
	ids, err := s.userAccountIDs(ctx, userID)
	if err != nil {
		return err
	}
	if !slices.Contains(ids, accountID) {
		return errors.NotFound("account not found")
	}
	return nil
}

func (s *Service) unclassifiedCount(ctx context.Context, accountID int64) int64 {
	var n int64
	s.orm.WithContext(ctx).Table("emails e").Joins("JOIN folders f ON f.id = e.folder_id").
		Where("e.account_id = ? AND f.type = ? AND e.imap_uid > 0 AND e.date >= ? AND COALESCE(e.category, '') = '' AND COALESCE(e.bulk_reason, '') = ''",
			accountID, schemas.FolderTypeInbox, time.Now().Add(-backfillWindow)).Count(&n)
	return n
}

// BackfillContacts deep-syncs the sent folder into contacts under the same
// per-account guard as BackfillMail, so two backfills never share an account's
// IMAP session. A busy account reports running and does nothing.
func (s *Service) BackfillContacts(ctx context.Context, userID, accountID int64) (BackfillContactsResponse, error) {
	if err := s.visibleAccount(ctx, userID, accountID); err != nil {
		return BackfillContactsResponse{}, err
	}
	release, ok := s.acquireBackfill(accountID)
	if !ok {
		return BackfillContactsResponse{Running: true}, nil
	}
	defer release()
	scanned, err := s.deepSyncSent(ctx, userID, accountID, maxDeepSyncSent)
	return BackfillContactsResponse{Scanned: scanned}, err
}
