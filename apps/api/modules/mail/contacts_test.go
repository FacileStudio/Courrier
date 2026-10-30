package mail

import (
	"context"
	"testing"
	"time"

	"github.com/FacileStudio/tronc/testdb"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
)

func contactFixture(t *testing.T) (*Service, schemas.Account, schemas.Account) {
	t.Helper()
	if testORM == nil {
		t.Skip(testdb.SkipReason("createdb courrier_test  # or point at any scratch database"))
	}
	for _, table := range []string{"contacts", "accounts"} {
		if err := testORM.Exec("DELETE FROM " + table).Error; err != nil {
			t.Fatalf("clear %s: %v", table, err)
		}
	}
	a := schemas.Account{UserID: 1, Name: "A", Email: "Me@Example.test"}
	b := schemas.Account{UserID: 1, Name: "B", Email: "other@example.test"}
	c := schemas.Account{UserID: 2, Name: "C", Email: "stranger@example.test"}
	for _, acc := range []*schemas.Account{&a, &b, &c} {
		if err := testORM.Create(acc).Error; err != nil {
			t.Fatalf("create account: %v", err)
		}
	}
	return NewService(testORM, nil), a, c
}

func TestUpsertContactsCountsAndNormalizes(t *testing.T) {
	_, a, _ := contactFixture(t)
	ctx := context.Background()
	old := time.Now().Add(-time.Hour)
	recs := []recipient{{Email: "Bob@X.test", Name: "Bob"}, {Email: "carol@x.test"}, {Email: "me@example.test"}}
	if err := upsertContacts(ctx, testORM, a.ID, nil, recs, old); err != nil {
		t.Fatal(err)
	}
	newer := time.Now()
	if err := upsertContacts(ctx, testORM, a.ID, nil, []recipient{{Email: "bob@x.test"}}, newer); err != nil {
		t.Fatal(err)
	}
	var rows []schemas.Contact
	testORM.Where("account_id = ?", a.ID).Order("email").Find(&rows)
	if len(rows) != 2 {
		t.Fatalf("want 2 contacts (own address skipped), got %d", len(rows))
	}
	if rows[0].Email != "bob@x.test" || rows[0].SentCount != 2 || rows[0].Name != "Bob" {
		t.Fatalf("bob row wrong: %+v", rows[0])
	}
	if rows[0].LastSentAt.Before(newer.Add(-time.Second)) {
		t.Fatalf("last_sent_at not the newest: %v", rows[0].LastSentAt)
	}
}

func TestSaveContactsAbsoluteIsIdempotent(t *testing.T) {
	_, a, _ := contactFixture(t)
	ctx := context.Background()
	for i := 0; i < 2; i++ {
		entries := map[string]*contactEntry{}
		addRecipients(entries, a.Email, []recipient{{Email: "d@x.test"}}, time.Now())
		addRecipients(entries, a.Email, []recipient{{Email: "d@x.test"}}, time.Now())
		if err := saveContacts(ctx, testORM, a.ID, nil, entries, true); err != nil {
			t.Fatal(err)
		}
	}
	var row schemas.Contact
	testORM.Where("account_id = ? AND email = ?", a.ID, "d@x.test").First(&row)
	if row.SentCount != 2 {
		t.Fatalf("want 2 after repeated deep sync, got %d", row.SentCount)
	}
}

func TestListContactsScopesAcrossAccounts(t *testing.T) {
	s, a, c := contactFixture(t)
	ctx := context.Background()
	var b schemas.Account
	testORM.Where("name = ?", "B").First(&b)
	now := time.Now()
	if err := upsertContacts(ctx, testORM, a.ID, nil, []recipient{{Email: "shared@x.test"}, {Email: "only-a@x.test"}}, now.Add(-time.Hour)); err != nil {
		t.Fatal(err)
	}
	if err := upsertContacts(ctx, testORM, b.ID, nil, []recipient{{Email: "shared@x.test"}}, now); err != nil {
		t.Fatal(err)
	}
	if err := upsertContacts(ctx, testORM, c.ID, nil, []recipient{{Email: "secret@x.test"}}, now); err != nil {
		t.Fatal(err)
	}

	rows, total, err := s.ListContacts(ctx, 1, nil, 1, 50)
	if err != nil || total != 2 || len(rows) != 2 {
		t.Fatalf("user 1 all accounts: total=%d rows=%d err=%v", total, len(rows), err)
	}
	if rows[0].Email != "shared@x.test" || rows[0].SentCount != 2 {
		t.Fatalf("shared row not merged first: %+v", rows[0])
	}
	if _, _, err := s.ListContacts(ctx, 1, &c.ID, 1, 50); err == nil {
		t.Fatal("user 1 must not list another user's account")
	}
	one, total, _ := s.ListContacts(ctx, 1, &a.ID, 1, 1)
	if total != 2 || len(one) != 1 {
		t.Fatalf("single account paging: total=%d rows=%d", total, len(one))
	}
	known, _ := s.KnownEmails(ctx, 1, []string{"SHARED@x.test", "secret@x.test"})
	if !known["shared@x.test"] || known["secret@x.test"] {
		t.Fatalf("known emails wrong: %v", known)
	}
}
