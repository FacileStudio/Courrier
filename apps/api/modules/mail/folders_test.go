package mail

import (
	"slices"
	"testing"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
)

func TestCreateFolderCreatesTheMailboxAndTheRow(t *testing.T) {
	env := newMailEnv(t)
	folder, err := env.svc.CreateFolder(ctx(), 1, env.account.ID, "  Clients ")
	if err != nil {
		t.Fatal(err)
	}
	if folder.Path != "Clients" || folder.Type != schemas.FolderTypeCustom || !slices.Contains(env.fake.mailboxes(t), "Clients") {
		t.Fatalf("folder %+v, server %v", folder, env.fake.mailboxes(t))
	}
	if _, err := env.svc.CreateFolder(ctx(), 1, env.account.ID, "clients"); err == nil {
		t.Error("a second folder with the same name must be a conflict")
	}
}

func TestCreateFolderRefusesNamesThatBreakIMAP(t *testing.T) {
	env := newMailEnv(t)
	for _, name := range []string{"", "   ", "a/b", "a*", "50%", ".hidden", "Sent Items", "Trash", "Spam", "Archive 2024", "inbox", "bad\x00"} {
		if _, err := env.svc.CreateFolder(ctx(), 1, env.account.ID, name); err == nil {
			t.Errorf("name %q must be refused", name)
		}
	}
	if got := len(env.fake.mailboxes(t)); got != 6 {
		t.Errorf("a refused name must not reach the server, mailboxes: %v", env.fake.mailboxes(t))
	}
}

/*
The in-memory test server does not rename inferior mailboxes, which RFC 3501 requires of a
real one, so the descendant check reads the cache only.
*/
func TestRenameFolderMovesDescendantsAndRefusesSystemFolders(t *testing.T) {
	env := newMailEnv(t)
	parent := env.folder(t, "Projects")
	if err := env.fake.user.Create("Projects/Acme", nil); err != nil {
		t.Fatal(err)
	}
	if err := env.svc.SyncAccount(ctx(), 1, env.account.ID); err != nil {
		t.Fatal(err)
	}
	renamed, err := env.svc.RenameFolder(ctx(), 1, env.account.ID, parent.ID, "Work")
	if err != nil {
		t.Fatal(err)
	}
	names := env.fake.mailboxes(t)
	if renamed.Path != "Work" || !slices.Contains(names, "Work") {
		t.Fatalf("renamed %+v, server %v", renamed, names)
	}
	child := env.folder(t, "Work/Acme")
	if child.Name != "Acme" {
		t.Errorf("child %+v", child)
	}
	for _, name := range []string{"Inbox", "Sent"} {
		sys := env.folder(t, map[string]string{"Inbox": "INBOX", "Sent": "Sent"}[name])
		if _, err := env.svc.RenameFolder(ctx(), 1, env.account.ID, sys.ID, "Renamed"); err == nil {
			t.Errorf("system folder %s must not be renamed", name)
		}
	}
	if _, err := env.svc.RenameFolder(ctx(), 1, env.account.ID, renamed.ID, "Drafts 2"); err == nil {
		t.Error("a reserved name must be refused on rename")
	}
}

func TestDeleteFolderNeedsForceWhenItHoldsMail(t *testing.T) {
	env := newMailEnv(t)
	env.seed(t, "Projects", "inside")
	projects := env.folder(t, "Projects")
	if err := env.svc.DeleteFolder(ctx(), 1, env.account.ID, projects.ID, false); err == nil {
		t.Fatal("a non-empty folder must need force")
	}
	if err := env.svc.DeleteFolder(ctx(), 1, env.account.ID, projects.ID, true); err != nil {
		t.Fatal(err)
	}
	var rows, folders int64
	testORM.Model(&schemas.Email{}).Where("folder_id = ?", projects.ID).Count(&rows)
	testORM.Model(&schemas.Folder{}).Where("id = ?", projects.ID).Count(&folders)
	if rows != 0 || folders != 0 || slices.Contains(env.fake.mailboxes(t), "Projects") {
		t.Errorf("rows %d, folders %d, server %v", rows, folders, env.fake.mailboxes(t))
	}
}

func TestDeleteFolderRemovesAnEmptyFolderAndRefusesSystemOnes(t *testing.T) {
	env := newMailEnv(t)
	if err := env.svc.DeleteFolder(ctx(), 1, env.account.ID, env.folder(t, "Projects").ID, false); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"INBOX", "Sent", "Drafts", "Trash", "Archive"} {
		if err := env.svc.DeleteFolder(ctx(), 1, env.account.ID, env.folder(t, name).ID, true); err == nil {
			t.Errorf("system folder %s must not be deleted", name)
		}
	}
	if len(env.fake.mailboxes(t)) != 5 {
		t.Errorf("server mailboxes %v", env.fake.mailboxes(t))
	}
}
