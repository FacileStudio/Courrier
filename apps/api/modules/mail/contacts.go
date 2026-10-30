package mail

import (
	"context"
	"strings"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type recipient struct {
	Email string
	Name  string
}

type contactEntry struct {
	Name  string
	Count int
	Last  time.Time
}

// ContactRow is one address the user has written to, merged across accounts.
type ContactRow struct {
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	SentCount  int       `json:"sent_count"`
	LastSentAt time.Time `json:"last_sent_at"`
}

func upsertContacts(ctx context.Context, orm *gorm.DB, accountID int64, spaceID *string, recipients []recipient, sentAt time.Time) error {
	var account schemas.Account
	if err := orm.WithContext(ctx).Select("email").Where("id = ?", accountID).First(&account).Error; err != nil {
		return err
	}
	entries := map[string]*contactEntry{}
	addRecipients(entries, account.Email, recipients, sentAt)
	return saveContacts(ctx, orm, accountID, spaceID, entries, false)
}

func addRecipients(entries map[string]*contactEntry, own string, recipients []recipient, sentAt time.Time) {
	own = strings.ToLower(strings.TrimSpace(own))
	seen := map[string]bool{}
	for _, r := range recipients {
		email := strings.ToLower(strings.TrimSpace(r.Email))
		if email == "" || email == own || seen[email] {
			continue
		}
		seen[email] = true
		entry, ok := entries[email]
		if !ok {
			entry = &contactEntry{}
			entries[email] = entry
		}
		entry.Count++
		if entry.Name == "" {
			entry.Name = strings.TrimSpace(r.Name)
		}
		if sentAt.After(entry.Last) {
			entry.Last = sentAt
		}
	}
}

func saveContacts(ctx context.Context, orm *gorm.DB, accountID int64, spaceID *string, entries map[string]*contactEntry, absolute bool) error {
	if len(entries) == 0 {
		return nil
	}
	rows := make([]schemas.Contact, 0, len(entries))
	for email, e := range entries {
		rows = append(rows, schemas.Contact{AccountID: accountID, SpaceID: spaceID, Email: email, Name: e.Name, SentCount: e.Count, LastSentAt: e.Last})
	}
	count := "contacts.sent_count + excluded.sent_count"
	if absolute {
		count = "GREATEST(contacts.sent_count, excluded.sent_count)"
	}
	return orm.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "account_id"}, {Name: "email"}},
		DoUpdates: clause.Assignments(map[string]any{
			"sent_count":   gorm.Expr(count),
			"last_sent_at": gorm.Expr("GREATEST(contacts.last_sent_at, excluded.last_sent_at)"),
			"name":         gorm.Expr("COALESCE(NULLIF(contacts.name, ''), excluded.name)"),
		}),
	}).CreateInBatches(&rows, 200).Error
}

func (s *Service) userAccountIDs(ctx context.Context, userID int64) ([]int64, error) {
	var ids []int64
	err := s.orm.WithContext(ctx).Model(&schemas.Account{}).
		Where("(user_id = ? AND space_id IS NULL) OR space_id IN (SELECT space_id::text FROM space_members WHERE user_id = ?)", userID, userID).
		Pluck("id", &ids).Error
	if err != nil {
		return nil, errors.Internal("failed to list accounts", err)
	}
	return ids, nil
}

// ListContacts returns the addresses the user has written to, newest first.
// A nil accountID covers every account the user can see, merging an address
// seen on several accounts into one row.
func (s *Service) ListContacts(ctx context.Context, userID int64, accountID *int64, page, limit int) ([]ContactRow, int64, error) {
	ids := []int64{}
	if accountID != nil {
		if _, err := s.getAccount(ctx, userID, *accountID); err != nil {
			return nil, 0, err
		}
		ids = append(ids, *accountID)
	} else {
		all, err := s.userAccountIDs(ctx, userID)
		if err != nil {
			return nil, 0, err
		}
		ids = all
	}
	merged := s.orm.WithContext(ctx).Table("contacts").
		Select("email, MAX(name) AS name, SUM(sent_count) AS sent_count, MAX(last_sent_at) AS last_sent_at").
		Where("account_id IN ?", ids).Group("email")
	var total int64
	if err := s.orm.WithContext(ctx).Table("(?) AS merged", merged).Count(&total).Error; err != nil {
		return nil, 0, errors.Internal("failed to count contacts", err)
	}
	rows := []ContactRow{}
	err := merged.Order("last_sent_at DESC, email ASC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error
	if err != nil {
		return nil, 0, errors.Internal("failed to list contacts", err)
	}
	return rows, total, nil
}

// KnownEmails returns which of the given addresses the user has written to
// from any account they can see. Keys are lowercased.
func (s *Service) KnownEmails(ctx context.Context, userID int64, emails []string) (map[string]bool, error) {
	known := map[string]bool{}
	ids, err := s.userAccountIDs(ctx, userID)
	if err != nil || len(emails) == 0 || len(ids) == 0 {
		return known, err
	}
	lowered := make([]string, len(emails))
	for i, e := range emails {
		lowered[i] = strings.ToLower(strings.TrimSpace(e))
	}
	var found []string
	err = s.orm.WithContext(ctx).Model(&schemas.Contact{}).
		Where("account_id IN ? AND email IN ?", ids, lowered).Distinct().Pluck("email", &found).Error
	for _, e := range found {
		known[e] = true
	}
	return known, err
}
