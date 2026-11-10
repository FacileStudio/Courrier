package mail

import (
	"context"
	"testing"
	"time"

	"github.com/FacileStudio/tronc/testdb"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
)

// searchFixture stores five inbox messages sharing one date, one in Trash,
// and one in a stranger's account.
func searchFixture(t *testing.T) (*Service, schemas.Account, schemas.Account) {
	t.Helper()
	if testORM == nil {
		t.Skip(testdb.SkipReason("createdb courrier_test  # or point at any scratch database"))
	}
	for _, table := range []string{"emails", "folders", "accounts"} {
		if err := testORM.Exec("DELETE FROM " + table).Error; err != nil {
			t.Fatalf("clear %s: %v", table, err)
		}
	}
	mine := schemas.Account{UserID: 31, Name: "Mine", Email: "mine@example.test"}
	theirs := schemas.Account{UserID: 32, Name: "Theirs", Email: "theirs@example.test"}
	testORM.Create(&mine)
	testORM.Create(&theirs)
	inbox := schemas.Folder{AccountID: mine.ID, Type: schemas.FolderTypeInbox, Path: "INBOX", Name: "Inbox"}
	trash := schemas.Folder{AccountID: mine.ID, Type: schemas.FolderTypeTrash, Path: "Trash", Name: "Trash"}
	testORM.Create(&inbox)
	testORM.Create(&trash)
	same := time.Date(2026, 11, 2, 8, 0, 0, 0, time.UTC)
	for range 5 {
		testORM.Create(&schemas.Email{AccountID: mine.ID, FolderID: inbox.ID, Subject: "report", Date: same, BodyText: "quarterly numbers"})
	}
	testORM.Create(&schemas.Email{AccountID: mine.ID, FolderID: trash.ID, Subject: "report", Date: same})
	testORM.Create(&schemas.Email{AccountID: theirs.ID, Subject: "report", Date: same})
	return NewService(testORM, nil), mine, theirs
}

// TestSearchMessagesPagesByKey walks pages of two through messages that share
// a date, and checks Trash stays out unless asked for.
func TestSearchMessagesPagesByKey(t *testing.T) {
	svc, mine, _ := searchFixture(t)
	ctx := context.Background()
	q := MessageQuery{AccountIDs: []int64{mine.ID}, Text: "report", Limit: 2}
	var seen []int64
	for {
		page, err := svc.SearchMessages(ctx, mine.UserID, q)
		if err != nil {
			t.Fatal(err)
		}
		for _, email := range page {
			seen = append(seen, email.ID)
		}
		if len(page) < q.Limit {
			break
		}
		last := page[len(page)-1]
		q.Before = &MessageKey{Date: last.Date, ID: last.ID}
	}
	for i := 1; i < len(seen); i++ {
		if seen[i] >= seen[i-1] {
			t.Fatalf("want ids strictly descending, got %v", seen)
		}
	}
	if len(seen) != 5 {
		t.Fatalf("want the 5 inbox messages and no Trash, got %v", seen)
	}
	trashed, _ := svc.SearchMessages(ctx, mine.UserID, MessageQuery{FolderType: schemas.FolderTypeTrash, Limit: 10})
	if len(trashed) != 1 {
		t.Fatalf("want the one Trash message when asked, got %d", len(trashed))
	}
}

// TestSearchMessagesRefusesAnotherPersonsAccount checks the account filter
// cannot reach an account the person does not see.
func TestSearchMessagesRefusesAnotherPersonsAccount(t *testing.T) {
	svc, mine, theirs := searchFixture(t)
	if _, err := svc.SearchMessages(context.Background(), mine.UserID, MessageQuery{AccountIDs: []int64{theirs.ID}, Limit: 5}); err == nil {
		t.Fatal("another person's account must be not found")
	}
	all, _ := svc.SearchMessages(context.Background(), mine.UserID, MessageQuery{Text: "report", Limit: 50})
	for _, email := range all {
		if email.AccountID != mine.ID {
			t.Fatalf("a search over every account leaked account %d", email.AccountID)
		}
	}
}
