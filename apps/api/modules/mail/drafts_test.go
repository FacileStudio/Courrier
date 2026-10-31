package mail

import (
	"strings"
	"testing"

	"github.com/emersion/go-imap/v2"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
)

func saveDraft(t *testing.T, env *mailEnv, req SendRequest) schemas.Email {
	t.Helper()
	draft, err := env.svc.SaveDraft(ctx(), 1, env.account.ID, req)
	if err != nil {
		t.Fatal(err)
	}
	return draft
}

func reload(t *testing.T, id int64) schemas.Email {
	t.Helper()
	var email schemas.Email
	if err := testORM.First(&email, id).Error; err != nil {
		t.Fatal(err)
	}
	return email
}

func TestSaveDraftRecordsTheServerUIDAndBcc(t *testing.T) {
	env := newMailEnv(t)
	draft := saveDraft(t, env, SendRequest{To: []string{"to@x.test"}, Bcc: []string{"hidden@x.test"}, Subject: "d1", Body: "text"})
	if draft.IMAPUID == 0 {
		t.Fatal("the draft row must carry the uid of its Drafts copy")
	}
	if got := env.fake.uids(t, "Drafts"); len(got) != 1 || uint32(got[0]) != draft.IMAPUID {
		t.Fatalf("server uids %v, row uid %d", got, draft.IMAPUID)
	}
	if !strings.Contains(draft.BccAddresses, "hidden@x.test") {
		t.Errorf("bcc_addresses = %s", draft.BccAddresses)
	}
}

func TestUpdateDraftKeepsTheIDAndReplacesTheServerCopy(t *testing.T) {
	env := newMailEnv(t)
	draft := saveDraft(t, env, SendRequest{To: []string{"to@x.test"}, Subject: "before", Body: "old"})
	updated, err := env.svc.UpdateDraft(ctx(), 1, env.account.ID, draft.ID,
		DraftRequest{SendRequest: SendRequest{To: []string{"to@x.test"}, Bcc: []string{"b@x.test"}, Subject: "after", Body: "new"}})
	if err != nil {
		t.Fatal(err)
	}
	row := reload(t, draft.ID)
	if updated.ID != draft.ID || row.Subject != "after" || row.BodyText != "new" || row.MessageID != draft.MessageID {
		t.Fatalf("row after update: %+v", row)
	}
	if got := env.fake.uids(t, "Drafts"); len(got) != 1 || uint32(got[0]) != row.IMAPUID || row.IMAPUID == draft.IMAPUID {
		t.Fatalf("server uids %v, row uid %d (was %d)", got, row.IMAPUID, draft.IMAPUID)
	}
	var count int64
	testORM.Model(&schemas.Email{}).Where("subject IN ?", []string{"before", "after"}).Count(&count)
	if count != 1 {
		t.Errorf("%d rows for the draft, want 1", count)
	}
}

func TestDraftAttachmentsSurviveUpdatesAndCanBeRemoved(t *testing.T) {
	env := newMailEnv(t)
	file := AttachmentUpload{Filename: "a.txt", MimeType: "text/plain", Data: []byte("alpha")}
	draft, err := env.svc.SaveDraft(ctx(), 1, env.account.ID, SendRequest{To: []string{"t@x.test"}, Subject: "att", Body: "b", Attachments: []AttachmentUpload{file}})
	if err != nil {
		t.Fatal(err)
	}
	second := AttachmentUpload{Filename: "b.txt", MimeType: "text/plain", Data: []byte("beta")}
	req := DraftRequest{SendRequest: SendRequest{To: []string{"t@x.test"}, Subject: "att2", Body: "b", Attachments: []AttachmentUpload{second}}}
	if _, err := env.svc.UpdateDraft(ctx(), 1, env.account.ID, draft.ID, req); err != nil {
		t.Fatal(err)
	}
	assertDraftFiles(t, env, draft.ID, map[string]string{"a.txt": "alpha", "b.txt": "beta"})

	var rows []schemas.Attachment
	testORM.Where("email_id = ? AND filename = ?", draft.ID, "a.txt").Find(&rows)
	drop := DraftRequest{SendRequest: SendRequest{To: []string{"t@x.test"}, Subject: "att3", Body: "b"}, RemoveAttachments: []int64{rows[0].ID}}
	if _, err := env.svc.UpdateDraft(ctx(), 1, env.account.ID, draft.ID, drop); err != nil {
		t.Fatal(err)
	}
	assertDraftFiles(t, env, draft.ID, map[string]string{"b.txt": "beta"})
}

func assertDraftFiles(t *testing.T, env *mailEnv, id int64, want map[string]string) {
	t.Helper()
	row := reload(t, id)
	var rows []schemas.Attachment
	testORM.Where("email_id = ?", id).Find(&rows)
	if len(rows) != len(want) || row.HasAttachments != (len(want) > 0) {
		t.Fatalf("attachment rows %+v, want files %v", rows, want)
	}
	client := env.fake.client(t)
	for _, a := range rows {
		parts, _ := parsePartID(a.PartID)
		data, err := fetchAttachmentPart(client, "Drafts", imap.UID(row.IMAPUID), parts)
		if err != nil || string(data) != want[a.Filename] {
			t.Errorf("%s: %q, %v (want %q)", a.Filename, data, err, want[a.Filename])
		}
	}
}

func TestDeleteDraftRemovesTheServerCopyAndItsAttachments(t *testing.T) {
	env := newMailEnv(t)
	file := AttachmentUpload{Filename: "a.txt", MimeType: "text/plain", Data: []byte("alpha")}
	draft := saveDraft(t, env, SendRequest{To: []string{"t@x.test"}, Subject: "gone", Body: "b", Attachments: []AttachmentUpload{file}})
	if err := env.svc.DeleteDraft(ctx(), 1, env.account.ID, draft.ID); err != nil {
		t.Fatal(err)
	}
	var left int64
	testORM.Model(&schemas.Attachment{}).Where("email_id = ?", draft.ID).Count(&left)
	if left != 0 || env.fake.count(t, "Drafts") != 0 {
		t.Errorf("attachments left %d, server drafts %d", left, env.fake.count(t, "Drafts"))
	}
	if err := env.svc.DeleteDraft(ctx(), 1, env.account.ID, draft.ID); err == nil {
		t.Error("deleting twice must be a 404")
	}
}

func TestUpdateDraftRefusesAMessageThatIsNotADraft(t *testing.T) {
	env := newMailEnv(t)
	inbox := env.seed(t, "INBOX", "not a draft")
	if _, err := env.svc.UpdateDraft(ctx(), 1, env.account.ID, inbox.ID, DraftRequest{}); err == nil {
		t.Fatal("an inbox message is not a draft")
	}
}

func TestUpdateDraftFailsCleanlyWhenTheServerIsDown(t *testing.T) {
	env := newMailEnv(t)
	draft := saveDraft(t, env, SendRequest{To: []string{"t@x.test"}, Subject: "keep", Body: "b"})
	env.svc.dial = connectIMAPRefused
	_, err := env.svc.UpdateDraft(ctx(), 1, env.account.ID, draft.ID, DraftRequest{SendRequest: SendRequest{Subject: "lost", Body: "b"}})
	if err == nil {
		t.Fatal("an unreachable server must fail the update")
	}
	if reload(t, draft.ID).Subject != "keep" {
		t.Error("a failed update must leave the draft unchanged")
	}
}
