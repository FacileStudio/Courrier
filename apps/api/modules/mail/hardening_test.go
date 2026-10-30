package mail

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/FacileStudio/tronc/errors"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
)

func TestBackfillSelectsRowsWithNullColumns(t *testing.T) {
	job, fake, env := backfillEnv(t)
	addUID(t, env, "news@brand.test", time.Hour, schemas.Email{IMAPUID: 1})
	addUID(t, env, "friend@x.test", time.Hour, schemas.Email{IMAPUID: 2})
	testORM.Exec("UPDATE emails SET category = NULL, bulk_reason = NULL")
	fake.headers[1] = Signals{ListUnsubscribe: "<mailto:u@brand.test>"}
	if n := env.svc.unclassifiedCount(context.Background(), env.a.ID); n != 2 {
		t.Fatalf("NULL rows must count as unclassified, got %d", n)
	}
	n, err := job.run(context.Background())
	if err != nil || n != 2 {
		t.Fatalf("NULL rows must be classified: n=%d err=%v", n, err)
	}
	var row schemas.Email
	testORM.First(&row, "imap_uid = ?", 1)
	if row.Category != "bulk" {
		t.Fatalf("row not classified: %+v", row)
	}
}

func TestUnclassifiedCountMatchesBackfillSelection(t *testing.T) {
	env := newTriageEnv(t)
	addUID(t, env, "a@x.test", time.Hour, schemas.Email{IMAPUID: 1})
	addUID(t, env, "b@x.test", time.Hour, schemas.Email{IMAPUID: 0})
	addUID(t, env, "c@x.test", 61*24*time.Hour, schemas.Email{IMAPUID: 3})
	addMail(t, mailKey{env.a.ID, env.sentA, "s"}, "d@x.test", time.Now(), schemas.Email{IMAPUID: 4})
	if n := env.svc.unclassifiedCount(context.Background(), env.a.ID); n != 1 {
		t.Fatalf("only the recent inbox row with a UID can be classified, got %d", n)
	}
}

func TestCooldownAfterFailedBackgroundRun(t *testing.T) {
	env := newTriageEnv(t)
	if env.svc.cooling(env.a.ID) {
		t.Fatal("no cool-down before any run")
	}
	env.svc.runBackgroundBackfill(999, env.a.ID)
	if !env.svc.cooling(env.a.ID) {
		t.Fatal("a failed run must start the cool-down")
	}
}

func TestBackfillContactsSharesTheGuard(t *testing.T) {
	env := newTriageEnv(t)
	release, ok := env.svc.acquireBackfill(env.a.ID)
	if !ok {
		t.Fatal("guard must be free")
	}
	defer release()
	res, err := env.svc.BackfillContacts(context.Background(), 1, env.a.ID)
	if err != nil || !res.Running {
		t.Fatalf("busy account must report running: %+v %v", res, err)
	}
}

func TestTriageReportsTruncation(t *testing.T) {
	env := newTriageEnv(t)
	err := testORM.Exec(`INSERT INTO emails (account_id, folder_id, thread_id, message_id, from_address, date, is_read, category, bulk_reason)
		SELECT ?, ?, 't' || g, 'm' || g, 'p@x.test', now(), false, '', '' FROM generate_series(1, ?) g`,
		env.a.ID, env.inboxA, maxTriageRows+1).Error
	if err != nil {
		t.Fatal(err)
	}
	resp, err := env.svc.Triage(context.Background(), 1, TriageQuery{Page: 1, Limit: 10})
	if err != nil || !resp.Truncated || resp.Total != maxTriageRows {
		t.Fatalf("want truncated at %d: truncated=%v total=%d err=%v", maxTriageRows, resp.Truncated, resp.Total, err)
	}
	testORM.Exec("DELETE FROM emails WHERE thread_id = 't1'")
	resp, _ = env.svc.Triage(context.Background(), 1, TriageQuery{Page: 1, Limit: 10})
	if resp.Truncated {
		t.Fatal("exactly the cap is not truncated")
	}
}

func TestHugePageDoesNotPanic(t *testing.T) {
	env := newTriageEnv(t)
	seedClassOrder(t, env)
	req := httptest.NewRequest(http.MethodGet, "/api/mail/triage?page=184467440737095518&limit=100", nil)
	q, err := parseTriageQuery(req)
	if err != nil || q.Page != maxPage {
		t.Fatalf("page must clamp: %d %v", q.Page, err)
	}
	if resp, err := env.svc.Triage(context.Background(), 1, q); err != nil || len(resp.Items) != 0 {
		t.Fatalf("far page is empty: %v %v", resp.Items, err)
	}
	rec := httptest.NewRecorder()
	env.svc.serveContacts(rec, httptest.NewRequest(http.MethodGet, "/x?page=184467440737095518&limit=9999999999999", nil), 1, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("contacts huge page: %d %s", rec.Code, rec.Body)
	}
}

func TestSinceBound(t *testing.T) {
	now := time.Now()
	if _, err := parseSince("3650d", now); err != nil {
		t.Fatal(err)
	}
	for _, raw := range []string{"3651d", "99999999999d", "999999999999999999999h"} {
		if _, err := parseSince(raw, now); err == nil {
			t.Fatalf("%s must be rejected", raw)
		}
	}
}

func TestInvisibleAccountsAreUniformNotFound(t *testing.T) {
	env := newTriageEnv(t)
	ctx := context.Background()
	space := schemas.Space{Name: "S"}
	testORM.Create(&space)
	shared := schemas.Account{UserID: 2, Name: "Shared", Email: "s@example.test", SpaceID: &space.ID}
	testORM.Create(&shared)
	upsertContacts(ctx, testORM, env.c.ID, nil, []recipient{{Email: "secret@x.test"}}, time.Now())
	for _, id := range []int64{env.c.ID, shared.ID, 987654} {
		checks := map[string]error{}
		_, checks["triage"] = env.svc.Triage(ctx, 1, TriageQuery{Page: 1, Limit: 5, AccountID: &id})
		_, _, checks["contacts"] = env.svc.ListContacts(ctx, 1, &id, 1, 5)
		_, checks["backfill"] = env.svc.BackfillMail(ctx, 1, id)
		_, checks["contacts backfill"] = env.svc.BackfillContacts(ctx, 1, id)
		for name, err := range checks {
			if errors.Status(err) != http.StatusNotFound {
				t.Errorf("%s on account %d: want 404, got %v", name, id, err)
			}
		}
	}
	rows, total, err := env.svc.ListContacts(ctx, 1, nil, 1, 5)
	if err != nil || total != 0 || len(rows) != 0 {
		t.Fatalf("user 1 must not see user 2 contacts: %v %d %v", rows, total, err)
	}
}

func TestListContactsAccountFilter(t *testing.T) {
	env := newTriageEnv(t)
	ctx := context.Background()
	upsertContacts(ctx, testORM, env.a.ID, nil, []recipient{{Email: "one@x.test"}}, time.Now())
	upsertContacts(ctx, testORM, env.b.ID, nil, []recipient{{Email: "two@x.test"}}, time.Now())
	id := env.a.ID
	rows, total, err := env.svc.ListContacts(ctx, 1, &id, 1, 5)
	if err != nil || total != 1 || fmt.Sprint(rows[0].Email) != "one@x.test" {
		t.Fatalf("account filter ignored: %v %d %v", rows, total, err)
	}
}
