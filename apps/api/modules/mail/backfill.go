package mail

import (
	"context"
	"time"

	"github.com/emersion/go-imap/v2"

	"github.com/FacileStudio/Courrier/apps/api/schemas"

	"gorm.io/gorm"
)

const (
	backfillWindow = 60 * 24 * time.Hour
	backfillBatch  = 50
	backfillPause  = 250 * time.Millisecond
)

type headerFetcher interface {
	Headers(mailbox string, uids []imap.UID) (map[imap.UID]Signals, error)
	Close()
}

type backfillJob struct {
	orm     *gorm.DB
	fetcher headerFetcher
	folder  schemas.Folder
	cutoff  time.Time
	pause   time.Duration
}

func (j backfillJob) nextBatch(afterID int64) ([]schemas.Email, error) {
	var rows []schemas.Email
	err := j.orm.Select("id", "imap_uid", "from_address").
		Where("folder_id = ? AND COALESCE(category, '') = '' AND COALESCE(bulk_reason, '') = '' AND date >= ? AND imap_uid > 0 AND id > ?",
			j.folder.ID, j.cutoff, afterID).
		Order("id ASC").Limit(backfillBatch).Find(&rows).Error
	return rows, err
}

func (j backfillJob) classify(rows []schemas.Email) error {
	uids := make([]imap.UID, len(rows))
	for i, r := range rows {
		uids[i] = imap.UID(r.IMAPUID)
	}
	found, err := j.fetcher.Headers(j.folder.Path, uids)
	if err != nil {
		return err
	}
	for _, r := range rows {
		category, reason := classifyForStore(found[imap.UID(r.IMAPUID)], r.FromAddress)
		update := map[string]any{"category": category, "bulk_reason": reason}
		if err := j.orm.Model(&schemas.Email{}).Where("id = ?", r.ID).Updates(update).Error; err != nil {
			return err
		}
	}
	return nil
}

func (j backfillJob) wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(j.pause):
		return nil
	}
}

func (j backfillJob) run(ctx context.Context) (int, error) {
	done, afterID := 0, int64(0)
	for {
		rows, err := j.nextBatch(afterID)
		if err != nil || len(rows) == 0 {
			return done, err
		}
		if err := j.classify(rows); err != nil {
			return done, err
		}
		done += len(rows)
		afterID = rows[len(rows)-1].ID
		if len(rows) < backfillBatch {
			return done, nil
		}
		if err := j.wait(ctx); err != nil {
			return done, err
		}
	}
}
