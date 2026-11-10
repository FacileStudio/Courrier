package mail

import (
	"context"
	"slices"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"

	"gorm.io/gorm"
)

// FolderAny matches every folder, Trash and Junk included, in MessageQuery.
const FolderAny = "any"

// snippetColumns are the message columns a result list needs, with the body
// cut to a snippet in SQL so a page never loads whole bodies.
var snippetColumns = []string{
	"id", "account_id", "folder_id", "thread_id", "subject", "from_address", "from_name",
	"to_addresses", "cc_addresses", "date", "is_read", "is_starred", "has_attachments",
	"left(body_text, 400) AS body_text",
}

// MessageKey is the position of one message in the newest-first order: its
// date, then its id to break ties.
type MessageKey struct {
	Date time.Time
	ID   int64
}

// MessageQuery selects messages across the accounts a person can see. An
// empty AccountIDs means all of them. FolderType "" skips Trash and Junk,
// FolderAny keeps them, and any other value names one system folder type.
// Before resumes after the last message of a previous page.
type MessageQuery struct {
	AccountIDs []int64
	Text       string
	FolderType string
	UnreadOnly bool
	Since      time.Time
	Before     *MessageKey
	Limit      int
}

// VisibleAccounts returns the personal accounts of userID and the accounts of
// every space they belong to, by id.
func (s *Service) VisibleAccounts(ctx context.Context, userID int64) ([]schemas.Account, error) {
	var accounts []schemas.Account
	err := s.orm.WithContext(ctx).
		Where("(user_id = ? AND space_id IS NULL) OR space_id IN (SELECT space_id::text FROM space_members WHERE user_id = ?)", userID, userID).
		Order("id ASC").Find(&accounts).Error
	if err != nil {
		return nil, errors.Internal("failed to list accounts", err)
	}
	return accounts, nil
}

// SearchMessages returns one page of messages matching q, newest first, with
// a body snippet only. The order is (date, id) descending, so a page resumed
// from Before never repeats or skips a message while new mail arrives.
func (s *Service) SearchMessages(ctx context.Context, userID int64, q MessageQuery) ([]schemas.Email, error) {
	ids, err := s.queryAccounts(ctx, userID, q.AccountIDs)
	if err != nil {
		return nil, err
	}
	query := s.orm.WithContext(ctx).Model(&schemas.Email{}).Select(snippetColumns).Where("account_id IN ?", ids)
	query = s.withFolder(query, ids, q.FolderType)
	query = withMessageFilters(query, q)
	var emails []schemas.Email
	if err := query.Order("date DESC, id DESC").Limit(q.Limit).Find(&emails).Error; err != nil {
		return nil, errors.Internal("failed to search emails", err)
	}
	return emails, nil
}

// FoldersByID returns the folders of the given accounts keyed by id, for
// naming the folder each message sits in.
func (s *Service) FoldersByID(ctx context.Context, accountIDs []int64) (map[int64]schemas.Folder, error) {
	var folders []schemas.Folder
	if err := s.orm.WithContext(ctx).Where("account_id IN ?", accountIDs).Find(&folders).Error; err != nil {
		return nil, errors.Internal("failed to list folders", err)
	}
	byID := make(map[int64]schemas.Folder, len(folders))
	for _, folder := range folders {
		byID[folder.ID] = folder
	}
	return byID, nil
}

// queryAccounts resolves the requested account ids against the visible ones:
// all of them when none is named, and not_found for any other id.
func (s *Service) queryAccounts(ctx context.Context, userID int64, requested []int64) ([]int64, error) {
	visible, err := s.userAccountIDs(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(requested) == 0 {
		return visible, nil
	}
	for _, id := range requested {
		if !slices.Contains(visible, id) {
			return nil, errors.NotFound("account not found")
		}
	}
	return requested, nil
}

// withFolder narrows query to one folder type, or keeps Trash and Junk out.
func (s *Service) withFolder(query *gorm.DB, accountIDs []int64, folderType string) *gorm.DB {
	folders := s.orm.Model(&schemas.Folder{}).Select("id").Where("account_id IN ?", accountIDs)
	switch folderType {
	case FolderAny:
		return query
	case "":
		return query.Where("folder_id NOT IN (?)",
			folders.Where("type IN ?", []string{schemas.FolderTypeTrash, schemas.FolderTypeJunk}))
	default:
		return query.Where("folder_id IN (?)", folders.Where("type = ?", folderType))
	}
}

// withMessageFilters applies the text, unread, date and page filters of q.
func withMessageFilters(query *gorm.DB, q MessageQuery) *gorm.DB {
	if q.Text != "" {
		for _, word := range searchWords(q.Text) {
			like := "%" + escapeLikePattern(word) + "%"
			query = query.Where("(subject ILIKE ? OR from_name ILIKE ? OR from_address ILIKE ? OR body_text ILIKE ?)",
				like, like, like, like)
		}
	}
	if q.UnreadOnly {
		query = query.Where("is_read = false")
	}
	if !q.Since.IsZero() {
		query = query.Where("date >= ?", q.Since)
	}
	if q.Before != nil {
		query = query.Where("(date, id) < (?, ?)", q.Before.Date, q.Before.ID)
	}
	return query
}
