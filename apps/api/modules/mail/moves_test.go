package mail

import (
	"testing"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
)

func TestMoveToACustomFolderKeepsTheIDAndTracksTheNewUID(t *testing.T) {
	env := newMailEnv(t)
	a := env.seed(t, "INBOX", "first")
	b := env.seed(t, "INBOX", "second")
	projects := env.folder(t, "Projects")

	req := MoveEmailsRequest{FolderID: &projects.ID}
	if err := env.svc.MoveEmailsTo(ctx(), 1, env.account.ID, []int64{a.ID, b.ID}, req); err != nil {
		t.Fatal(err)
	}
	if env.fake.count(t, "INBOX") != 0 || env.fake.count(t, "Projects") != 2 {
		t.Fatalf("server: inbox %d, projects %d", env.fake.count(t, "INBOX"), env.fake.count(t, "Projects"))
	}
	want := map[uint32]bool{}
	for _, uid := range env.fake.uids(t, "Projects") {
		want[uint32(uid)] = true
	}
	for _, id := range []int64{a.ID, b.ID} {
		row := reload(t, id)
		if row.FolderID != projects.ID || !want[row.IMAPUID] {
			t.Errorf("row %d: folder %d uid %d, server uids %v", id, row.FolderID, row.IMAPUID, want)
		}
	}
	if err := env.svc.SyncFolderEmails(ctx(), 1, env.account.ID, projects.ID); err != nil {
		t.Fatal(err)
	}
	var rows int64
	testORM.Model(&schemas.Email{}).Where("account_id = ?", env.account.ID).Count(&rows)
	if rows != 2 {
		t.Errorf("a sync after the move must not duplicate rows, got %d", rows)
	}
}

func TestMoveByTypeAndRefusals(t *testing.T) {
	env := newMailEnv(t)
	mail := env.seed(t, "INBOX", "to archive")
	if err := env.svc.MoveEmailsTo(ctx(), 1, env.account.ID, []int64{mail.ID}, MoveEmailsRequest{FolderType: "archive"}); err != nil {
		t.Fatal(err)
	}
	if env.fake.count(t, "Archive") != 1 || reload(t, mail.ID).FolderID != env.folder(t, "Archive").ID {
		t.Error("archive by type did not move the message")
	}
	for _, bad := range []MoveEmailsRequest{{}, {FolderType: "custom"}, {FolderType: "nope"}, {FolderID: new(int64)}} {
		if err := env.svc.MoveEmailsTo(ctx(), 1, env.account.ID, []int64{mail.ID}, bad); err == nil {
			t.Errorf("destination %+v must be refused", bad)
		}
	}
}

func TestMoveRefusesAFolderOfAnotherAccount(t *testing.T) {
	env := newMailEnv(t)
	mail := env.seed(t, "INBOX", "mine")
	other := schemas.Account{UserID: 1, Name: "Other", Email: "o@example.test"}
	testORM.Create(&other)
	foreign := schemas.Folder{AccountID: other.ID, Path: "Secret", Name: "Secret", Type: schemas.FolderTypeCustom}
	testORM.Create(&foreign)
	err := env.svc.MoveEmailsTo(ctx(), 1, env.account.ID, []int64{mail.ID}, MoveEmailsRequest{FolderID: &foreign.ID})
	if err == nil {
		t.Fatal("a folder id from another account must not resolve")
	}
}

func TestDeleteMovesToTrashThenRemovesFromTheServer(t *testing.T) {
	env := newMailEnv(t)
	mail := env.seed(t, "INBOX", "doomed")
	if err := env.svc.DeleteEmails(ctx(), 1, env.account.ID, []int64{mail.ID}); err != nil {
		t.Fatal(err)
	}
	if env.fake.count(t, "Trash") != 1 {
		t.Fatal("the first delete must land in Trash")
	}
	if err := env.svc.DeleteEmails(ctx(), 1, env.account.ID, []int64{mail.ID}); err != nil {
		t.Fatal(err)
	}
	var left int64
	testORM.Model(&schemas.Email{}).Where("id = ?", mail.ID).Count(&left)
	if left != 0 || env.fake.count(t, "Trash") != 0 {
		t.Errorf("rows left %d, trash on server %d", left, env.fake.count(t, "Trash"))
	}
}

func TestBulkStarMirrorsToTheServer(t *testing.T) {
	env := newMailEnv(t)
	a := env.seed(t, "INBOX", "s1")
	b := env.seed(t, "INBOX", "s2")
	if err := env.svc.SetStarState(ctx(), 1, env.account.ID, []int64{a.ID, b.ID}, true); err != nil {
		t.Fatal(err)
	}
	if !reload(t, a.ID).IsStarred || !reload(t, b.ID).IsStarred {
		t.Fatal("rows not starred")
	}
	waitFor(t, "flagged on server", func() bool { return env.fake.flagged(t, "INBOX") == 2 })
	if err := env.svc.SetStarState(ctx(), 1, env.account.ID, []int64{a.ID}, false); err != nil {
		t.Fatal(err)
	}
	waitFor(t, "one unflagged", func() bool { return env.fake.flagged(t, "INBOX") == 1 })
	if !reload(t, b.ID).IsStarred || reload(t, a.ID).IsStarred {
		t.Error("unstar touched the wrong row")
	}
}
