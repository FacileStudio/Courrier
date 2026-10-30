package mail

import (
	"context"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"
)

const maxTriageRows = 5000

// TriageQuery filters the cross-account triage listing.
type TriageQuery struct {
	UnreadOnly bool
	Classes    map[string]bool
	Since      time.Time
	AccountID  *int64
	Page       int
	Limit      int
}

// TriageItem is one inbox conversation with its resolved class.
type TriageItem struct {
	AccountID    int64     `json:"account_id"`
	Account      string    `json:"account"`
	ThreadID     string    `json:"thread_id"`
	EmailID      int64     `json:"email_id"`
	Subject      string    `json:"subject"`
	FromAddress  string    `json:"from_address"`
	FromName     string    `json:"from_name"`
	Date         time.Time `json:"date"`
	MessageCount int       `json:"message_count"`
	UnreadCount  int       `json:"unread_count"`
	Class        string    `json:"class"`
	Reason       string    `json:"reason"`
}

// TriageResponse is the paginated, classified inbox of every account.
type TriageResponse struct {
	GeneratedAt time.Time      `json:"generated_at"`
	Total       int            `json:"total"`
	Counts      map[string]int `json:"counts"`
	Items       []TriageItem   `json:"items"`
	Page        int            `json:"page"`
	Limit       int            `json:"limit"`
}

func (s *Service) triageAccounts(ctx context.Context, userID int64, only *int64) (map[int64]string, error) {
	ids, err := s.userAccountIDs(ctx, userID)
	if err != nil {
		return nil, err
	}
	if only != nil {
		ids = filterID(ids, *only)
		if len(ids) == 0 {
			return nil, errors.NotFound("account not found")
		}
	}
	var accounts []schemas.Account
	if err := s.orm.WithContext(ctx).Select("id", "email").Where("id IN ?", ids).Find(&accounts).Error; err != nil {
		return nil, errors.Internal("failed to load accounts", err)
	}
	emails := make(map[int64]string, len(accounts))
	for _, a := range accounts {
		emails[a.ID] = a.Email
	}
	return emails, nil
}

func filterID(ids []int64, want int64) []int64 {
	for _, id := range ids {
		if id == want {
			return []int64{id}
		}
	}
	return nil
}

func (s *Service) classifyInbox(ctx context.Context, userID int64, accounts map[int64]string, q TriageQuery) ([]TriageItem, error) {
	ids := make([]int64, 0, len(accounts))
	for id := range accounts {
		ids = append(ids, id)
	}
	var folderIDs []int64
	err := s.orm.WithContext(ctx).Model(&schemas.Folder{}).
		Where("account_id IN ? AND type = ?", ids, schemas.FolderTypeInbox).Pluck("id", &folderIDs).Error
	if err != nil {
		return nil, errors.Internal("failed to list inbox folders", err)
	}
	filter := conversationFilter{FolderIDs: folderIDs, UnreadOnly: q.UnreadOnly, Since: q.Since}
	convs, err := s.queryConversations(ctx, filter, maxTriageRows, 0)
	if err != nil {
		return nil, err
	}
	return s.triageItems(ctx, userID, ids, accounts, convs)
}

func (s *Service) triageItems(ctx context.Context, userID int64, ids []int64, accounts map[int64]string, convs []Conversation) ([]TriageItem, error) {
	sent, err := s.sentThreads(ctx, ids, convs)
	if err != nil {
		return nil, err
	}
	senders := make([]string, 0, len(convs))
	for _, c := range convs {
		senders = append(senders, c.Email.FromAddress)
	}
	known, err := s.KnownEmails(ctx, userID, senders)
	if err != nil {
		return nil, err
	}
	items := make([]TriageItem, len(convs))
	for i, c := range convs {
		class, reason := resolveClass(c.Email, sent, known)
		items[i] = newTriageItem(c, accounts[c.Email.AccountID], class, reason)
	}
	return items, nil
}

func newTriageItem(c Conversation, account, class, reason string) TriageItem {
	e := c.Email
	return TriageItem{
		AccountID: e.AccountID, Account: account, ThreadID: e.ThreadID, EmailID: e.ID,
		Subject: e.Subject, FromAddress: e.FromAddress, FromName: e.FromName, Date: e.Date.UTC(),
		MessageCount: c.MessageCount, UnreadCount: c.UnreadCount, Class: class, Reason: reason,
	}
}

func (q TriageQuery) wants(class string) bool {
	return len(q.Classes) == 0 || q.Classes[class]
}

// Triage lists the unread-or-all inbox conversations of every account the
// user can see, each tagged known, unknown_human or bulk. Counts cover the
// whole unread, since and account filtered set whatever the class filter, so a
// client can report how much bulk it skipped; total and items follow the class
// filter, and items are the requested page.
func (s *Service) Triage(ctx context.Context, userID int64, q TriageQuery) (TriageResponse, error) {
	accounts, err := s.triageAccounts(ctx, userID, q.AccountID)
	if err != nil {
		return TriageResponse{}, err
	}
	all, err := s.classifyInbox(ctx, userID, accounts, q)
	if err != nil {
		return TriageResponse{}, err
	}
	counts := map[string]int{}
	for _, class := range triageClasses {
		counts[class] = 0
	}
	filtered := make([]TriageItem, 0, len(all))
	for _, it := range all {
		counts[it.Class]++
		if q.wants(it.Class) {
			filtered = append(filtered, it)
		}
	}
	return TriageResponse{
		GeneratedAt: time.Now().UTC(), Total: len(filtered), Counts: counts,
		Items: pageOf(filtered, q.Page, q.Limit), Page: q.Page, Limit: q.Limit,
	}, nil
}

func pageOf(items []TriageItem, page, limit int) []TriageItem {
	start := (page - 1) * limit
	if start >= len(items) {
		return []TriageItem{}
	}
	return items[start:min(start+limit, len(items))]
}
