package mail

import (
	"testing"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
)

func TestNewMailNotificationSkipsBulkReadAndOldMail(t *testing.T) {
	now := time.Date(2026, 11, 3, 12, 0, 0, 0, time.UTC)
	account := schemas.Account{Email: "hello@facile.studio"}
	fresh := []schemas.Email{
		{FromName: "Ada", FromAddress: "ada@example.com", Subject: "secret plans", Date: now.Add(-time.Minute)},
		{FromName: "Shop", Category: "bulk", Date: now},
		{FromName: "Old", Date: now.Add(-48 * time.Hour)},
		{FromName: "Seen", IsRead: true, Date: now},
	}
	n, ok := newMailNotification(account, fresh, now)
	if !ok || n.Title != "New mail from Ada" || n.Body != "hello@facile.studio" || n.Path != "mail" {
		t.Fatalf("got %+v ok=%v", n, ok)
	}

	fresh = append(fresh, schemas.Email{FromAddress: "bob@example.com", Date: now})
	n, ok = newMailNotification(account, fresh, now)
	if !ok || n.Title != "2 new messages" || n.Body != "From Ada and others, to hello@facile.studio" {
		t.Fatalf("got %+v ok=%v", n, ok)
	}

	if _, ok := newMailNotification(account, fresh[1:4], now); ok {
		t.Fatal("bulk, old and read mail must not notify")
	}
}

func TestSyncFolderReportsOnlyMailNewToAWarmCache(t *testing.T) {
	env := newMailEnv(t)
	inbox := env.folder(t, "INBOX")
	env.fake.put(t, "INBOX", "From: Ada <ada@x.test>\r\nSubject: one\r\nMessage-ID: <one@x.test>\r\n\r\nbody\r\n")
	if _, fresh, err := env.svc.syncFolder(ctx(), 1, env.account.ID, inbox.ID); err != nil || len(fresh) != 0 {
		t.Fatalf("first sync of a cold inbox: fresh=%d err=%v", len(fresh), err)
	}
	env.fake.put(t, "INBOX", "From: Bob <bob@x.test>\r\nSubject: two\r\nMessage-ID: <two@x.test>\r\n\r\nbody\r\n")
	_, fresh, err := env.svc.syncFolder(ctx(), 1, env.account.ID, inbox.ID)
	if err != nil || len(fresh) != 1 || fresh[0].FromName != "Bob" {
		t.Fatalf("second sync: fresh=%+v err=%v", fresh, err)
	}
}
