package spaces

import (
	"context"
	"testing"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	portepg "github.com/FacileStudio/porte/pg"

	"gorm.io/gorm"
)

func count(t *testing.T, db *gorm.DB, query string, args ...any) int64 {
	t.Helper()
	var n int64
	if err := db.Raw(query, args...).Scan(&n).Error; err != nil {
		t.Fatalf("count: %v", err)
	}
	return n
}

func mustCreate(t *testing.T, row any) {
	t.Helper()
	if err := testORM.Create(row).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}
}

type erasureFixture struct {
	gone, other int64
	deletions   *Deletions
}

func seedMailbox(t *testing.T, owner int64, spaceID *string, label string) schemas.Account {
	t.Helper()
	account := schemas.Account{UserID: owner, SpaceID: spaceID, Name: label, Email: label + "@example.test"}
	mustCreate(t, &account)
	folder := schemas.Folder{AccountID: account.ID, SpaceID: spaceID, Name: "INBOX"}
	mustCreate(t, &folder)
	email := schemas.Email{AccountID: account.ID, SpaceID: spaceID, FolderID: folder.ID, MessageID: label, Subject: label}
	mustCreate(t, &email)
	mustCreate(t, &schemas.Attachment{EmailID: email.ID, Filename: label + ".pdf"})
	mustCreate(t, &schemas.Contact{AccountID: account.ID, SpaceID: spaceID, Email: "c@" + label + ".test"})
	mustCreate(t, &schemas.ThreadLink{AccountID: account.ID, MessageID: label, ThreadID: label})
	mustCreate(t, &schemas.SenderRule{AccountID: account.ID, Address: "s@" + label + ".test", Category: "newsletters"})
	return account
}

func seedErasure(t *testing.T) erasureFixture {
	t.Helper()
	_, spaceID, _ := newSpace(t)
	for _, table := range []string{"attachments", "emails", "folders", "contacts", "thread_links", "sender_rules", "accounts", "email_templates"} {
		if err := testORM.Exec(`DELETE FROM ` + table).Error; err != nil {
			t.Fatalf("clear %s: %v", table, err)
		}
	}
	gone := user(t, "gone@facile.studio")
	identity(t, gone, "sub-gone")
	other := user(t, "other@facile.studio")
	identity(t, other, "sub-other")
	mustCreate(t, &schemas.SpaceMember{SpaceID: spaceID, UserID: gone, Role: schemas.SpaceRoleMember})

	seedMailbox(t, gone, nil, "private")
	seedMailbox(t, gone, &spaceID, "team")
	seedMailbox(t, other, nil, "theirs")
	mustCreate(t, &schemas.EmailTemplate{UserID: gone, Name: "private-template"})
	mustCreate(t, &schemas.EmailTemplate{UserID: gone, SpaceID: &spaceID, Name: "team-template"})
	mustCreate(t, &schemas.EmailTemplate{UserID: other, Name: "their-template"})

	seedCredentials(t, gone)
	seedCredentials(t, other)
	deletions := newErasureDeletions(t)
	return erasureFixture{gone: gone, other: other, deletions: deletions}
}

// A subject Registre lists as erased, handled by the same receiver the
// deletion poller runs, loses the account and every personal mail account
// with what hangs off it. A mailbox attached to a space stays, other people's
// data is untouched, and a replay is a no-op.
func TestErasePersonDeletesPersonalDataAndKeepsSpaceData(t *testing.T) {
	f := seedErasure(t)
	pollErased(t, f.deletions)
	if err := f.deletions.ErasePerson(context.Background(), testIssuer, "sub-gone"); err != nil {
		t.Fatalf("replayed erase: %v", err)
	}

	checks := map[string]int64{
		`SELECT count(*) FROM users WHERE id = ?`:                                     0,
		`SELECT count(*) FROM space_members WHERE user_id = ?`:                        0,
		`SELECT count(*) FROM accounts WHERE user_id = ? AND space_id IS NULL`:        0,
		`SELECT count(*) FROM email_templates WHERE user_id = ? AND space_id IS NULL`: 0,
		`SELECT count(*) FROM porte_identities WHERE user_id = ?`:                     0,
	}
	for query, want := range checks {
		if got := count(t, testORM, query, f.gone); got != want {
			t.Errorf("%s = %d, want %d", query, got, want)
		}
	}
	for _, table := range []string{"folders", "emails", "contacts", "thread_links", "sender_rules"} {
		query := `SELECT count(*) FROM ` + table + ` WHERE account_id IN (SELECT id FROM accounts WHERE name = ?)`
		if got := count(t, testORM, query, "private"); got != 0 {
			t.Errorf("%s of the erased mailbox = %d, want 0", table, got)
		}
		if got := count(t, testORM, query, "team"); got != 1 {
			t.Errorf("%s of the space mailbox = %d, want 1", table, got)
		}
		if got := count(t, testORM, query, "theirs"); got != 1 {
			t.Errorf("%s of another person's mailbox = %d, want 1", table, got)
		}
	}
	if got := count(t, testORM, `SELECT count(*) FROM attachments`); got != 2 {
		t.Errorf("attachments left = %d, want 2 (space and other)", got)
	}
	for name, want := range map[string]int64{"team-template": 1, "their-template": 1, "private-template": 0} {
		if got := count(t, testORM, `SELECT count(*) FROM email_templates WHERE name = ?`, name); got != want {
			t.Errorf("template %q = %d, want %d", name, got, want)
		}
	}
	if got := count(t, testORM, `SELECT count(*) FROM users WHERE id = ?`, f.other); got != 1 {
		t.Errorf("another account was touched")
	}
}

func TestErasePersonForAnUnknownSubjectChangesNothing(t *testing.T) {
	f := seedErasure(t)
	if err := f.deletions.ErasePerson(context.Background(), testIssuer, "nobody"); err != nil {
		t.Fatalf("erase: %v", err)
	}
	if got := count(t, testORM, `SELECT count(*) FROM users WHERE id = ?`, f.gone); got != 1 {
		t.Errorf("an unknown subject deleted an account")
	}
}

func TestDeletionStateSurvivesARestart(t *testing.T) {
	newSpace(t)
	state := deletionState{orm: testORM}
	if err := testORM.Exec(`DELETE FROM directory_deletion_state`).Error; err != nil {
		t.Fatal(err)
	}

	if saved, err := state.Load(context.Background()); err != nil || saved != nil {
		t.Fatalf("empty Load = %q, %v; want nil", saved, err)
	}
	for _, blob := range []string{"one", "two"} {
		if err := state.Save(context.Background(), []byte(blob)); err != nil {
			t.Fatalf("save: %v", err)
		}
	}
	if saved, err := state.Load(context.Background()); err != nil || string(saved) != "two" {
		t.Fatalf("Load = %q, %v; want two", saved, err)
	}
}

func TestPorteEraseIdentityWorksWithoutAMagicLinksTable(t *testing.T) {
	f := seedErasure(t)
	sqlDB, err := testORM.DB()
	if err != nil {
		t.Fatal(err)
	}
	if err := portepg.New(sqlDB).EraseIdentity(context.Background(), testIssuer, "sub-gone"); err != nil {
		t.Fatalf("erase identity: %v", err)
	}
	if got := count(t, testORM, `SELECT count(*) FROM porte_identities WHERE user_id = ?`, f.gone); got != 0 {
		t.Errorf("identity rows left = %d, want 0", got)
	}
}
