package mail

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/FacileStudio/tronc/testdb"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
)

type triageEnv struct {
	svc                       *Service
	a, b, c                   schemas.Account
	inboxA, sentA, inboxOther int64
}

func newFolder(t *testing.T, accountID int64, kind string) int64 {
	t.Helper()
	f := schemas.Folder{AccountID: accountID, Path: kind, Name: kind, Type: kind}
	if err := testORM.Create(&f).Error; err != nil {
		t.Fatal(err)
	}
	return f.ID
}

type mailKey struct {
	Account, Folder int64
	Thread          string
}

func addMail(t *testing.T, k mailKey, from string, date time.Time, e schemas.Email) {
	t.Helper()
	e.AccountID, e.FolderID, e.ThreadID, e.FromAddress, e.Date = k.Account, k.Folder, k.Thread, from, date
	e.MessageID = fmt.Sprintf("<%s-%d@x.test>", k.Thread, date.UnixNano())
	if err := testORM.Create(&e).Error; err != nil {
		t.Fatal(err)
	}
}

func newTriageEnv(t *testing.T) triageEnv {
	t.Helper()
	if testORM == nil {
		t.Skip(testdb.SkipReason("createdb courrier_test  # or point at any scratch database"))
	}
	for _, table := range []string{"emails", "folders", "contacts", "accounts"} {
		testORM.Exec("DELETE FROM " + table)
	}
	env := triageEnv{svc: NewService(testORM, nil)}
	env.a = schemas.Account{UserID: 1, Name: "A", Email: "a@example.test"}
	env.b = schemas.Account{UserID: 1, Name: "B", Email: "b@example.test"}
	env.c = schemas.Account{UserID: 2, Name: "C", Email: "c@example.test"}
	for _, acc := range []*schemas.Account{&env.a, &env.b, &env.c} {
		testORM.Create(acc)
	}
	env.inboxA = newFolder(t, env.a.ID, schemas.FolderTypeInbox)
	env.sentA = newFolder(t, env.a.ID, schemas.FolderTypeSent)
	env.inboxOther = newFolder(t, env.c.ID, schemas.FolderTypeInbox)
	return env
}

func classOf(resp TriageResponse) map[string]string {
	out := map[string]string{}
	for _, it := range resp.Items {
		out[it.ThreadID] = it.Class + "/" + it.Reason
	}
	return out
}

func seedClassOrder(t *testing.T, env triageEnv) {
	t.Helper()
	now := time.Now()
	addMail(t, mailKey{env.a.ID, env.inboxA, "t1"}, "noreply@shop.test", now, schemas.Email{})
	addMail(t, mailKey{env.a.ID, env.sentA, "t1"}, "a@example.test", now, schemas.Email{IsRead: true})
	addMail(t, mailKey{env.a.ID, env.inboxA, "t2"}, "friend@x.test", now, schemas.Email{})
	addMail(t, mailKey{env.a.ID, env.sentA, "t2"}, "a@example.test", now, schemas.Email{IsRead: true})
	addMail(t, mailKey{env.a.ID, env.inboxA, "t3"}, "bob@x.test", now, schemas.Email{})
	addMail(t, mailKey{env.a.ID, env.inboxA, "t4"}, "stranger@x.test", now, schemas.Email{})
	addMail(t, mailKey{env.a.ID, env.inboxA, "t5"}, "news@brand.test", now, schemas.Email{Category: "bulk", BulkReason: "list-unsubscribe"})
	addMail(t, mailKey{env.c.ID, env.inboxOther, "t6"}, "x@x.test", now, schemas.Email{})
	upsertContacts(context.Background(), testORM, env.b.ID, nil, []recipient{{Email: "bob@x.test"}}, now)
}
