package mail

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/internal/push"
	"github.com/FacileStudio/Courrier/apps/api/schemas"
)

const (
	notifyMaxAge    = 24 * time.Hour
	inboxPollBudget = 90 * time.Second
	newMailPath     = "mail"
)

type watchedInbox struct {
	FolderID  int64
	AccountID int64
	UserID    int64
}

// WithNotifier makes the inbox watcher announce new mail through hub.
func (s *Service) WithNotifier(hub *push.Hub) *Service {
	s.notifier = hub
	return s
}

func (s *Service) lockFolder(folderID int64) func() {
	lock, _ := s.folderLocks.LoadOrStore(folderID, &sync.Mutex{})
	mu := lock.(*sync.Mutex)
	mu.Lock()
	return mu.Unlock
}

// WatchInboxes syncs, every interval, the inbox of each account whose owner
// has a device registered for push or a desktop app listening, and notifies the owner of the mail it
// finds. A sync the web app already ran wins: mail it cached is not new here.
// It returns when ctx ends.
func (s *Service) WatchInboxes(ctx context.Context, interval time.Duration, logger *slog.Logger) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.pollInboxes(ctx, logger)
		}
	}
}

func (s *Service) pollInboxes(ctx context.Context, logger *slog.Logger) {
	var inboxes []watchedInbox
	live := append([]int64{0}, s.notifier.LiveUsers()...)
	err := s.orm.WithContext(ctx).Raw(
		`SELECT f.id AS folder_id, a.id AS account_id, a.user_id
		   FROM folders f JOIN accounts a ON a.id = f.account_id
		  WHERE f.type = ? AND (a.user_id IN (SELECT DISTINCT user_id FROM push_devices) OR a.user_id IN ?)`,
		schemas.FolderTypeInbox, live,
	).Scan(&inboxes).Error
	if err != nil {
		logger.Warn("listing watched inboxes failed", slog.Any("error", err))
		return
	}
	for _, inbox := range inboxes {
		syncCtx, cancel := context.WithTimeout(ctx, inboxPollBudget)
		account, fresh, err := s.syncFolder(syncCtx, inbox.UserID, inbox.AccountID, inbox.FolderID)
		cancel()
		if err != nil {
			logger.Debug("inbox poll failed", slog.Int64("account_id", inbox.AccountID), slog.Any("error", err))
			continue
		}
		s.announce(ctx, account, fresh, time.Now())
	}
}

func (s *Service) announce(ctx context.Context, account schemas.Account, fresh []schemas.Email, now time.Time) {
	if n, ok := newMailNotification(account, fresh, now); ok {
		s.notifier.Notify(ctx, account.UserID, n)
	}
}

// newMailNotification summarises the mail worth a notification: unread,
// sent by a person rather than a list or a robot, and under a day old. It
// names senders and the receiving address only, never a subject or a body.
func newMailNotification(account schemas.Account, fresh []schemas.Email, now time.Time) (push.Notification, bool) {
	var senders []string
	for _, e := range fresh {
		if e.IsRead || e.Category != "" || now.Sub(e.Date) > notifyMaxAge {
			continue
		}
		senders = append(senders, senderName(e))
	}
	n := push.Notification{Body: account.Email, Path: newMailPath}
	switch len(senders) {
	case 0:
		return n, false
	case 1:
		n.Title = "New mail from " + senders[0]
	default:
		n.Title = fmt.Sprintf("%d new messages", len(senders))
		n.Body = fmt.Sprintf("From %s and others, to %s", senders[0], account.Email)
	}
	return n, true
}

func senderName(e schemas.Email) string {
	if name := strings.TrimSpace(e.FromName); name != "" {
		return name
	}
	if address := strings.TrimSpace(e.FromAddress); address != "" {
		return address
	}
	return "an unknown sender"
}
