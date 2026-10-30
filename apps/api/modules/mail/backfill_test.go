package mail

import (
	"context"
	"testing"
	"time"

	"github.com/emersion/go-imap/v2"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
)

type fakeFetcher struct {
	headers map[imap.UID]Signals
	calls   int
	seen    int
}

func (f *fakeFetcher) Headers(_ string, uids []imap.UID) (map[imap.UID]Signals, error) {
	f.calls++
	f.seen += len(uids)
	return f.headers, nil
}

func (f *fakeFetcher) Close() {}

func backfillEnv(t *testing.T) (backfillJob, *fakeFetcher, triageEnv) {
	t.Helper()
	env := newTriageEnv(t)
	fake := &fakeFetcher{headers: map[imap.UID]Signals{}}
	var folder schemas.Folder
	testORM.First(&folder, env.inboxA)
	job := backfillJob{orm: testORM, fetcher: fake, folder: folder, cutoff: time.Now().Add(-backfillWindow)}
	return job, fake, env
}

func addUID(t *testing.T, env triageEnv, from string, age time.Duration, e schemas.Email) {
	t.Helper()
	addMail(t, mailKey{env.a.ID, env.inboxA, "u"}, from, time.Now().Add(-age), e)
}

func TestBackfillClassifiesAndIsIdempotent(t *testing.T) {
	job, fake, env := backfillEnv(t)
	addUID(t, env, "news@brand.test", time.Hour, schemas.Email{IMAPUID: 1})
	addUID(t, env, "friend@x.test", time.Hour, schemas.Email{IMAPUID: 2})
	fake.headers[1] = Signals{ListUnsubscribe: "<mailto:u@brand.test>"}
	n, err := job.run(context.Background())
	if err != nil || n != 2 {
		t.Fatalf("first run: n=%d err=%v", n, err)
	}
	var rows []schemas.Email
	testORM.Order("imap_uid").Find(&rows, "folder_id = ?", env.inboxA)
	if rows[0].Category != "bulk" || rows[0].BulkReason != "list-unsubscribe" {
		t.Fatalf("row 1 not bulk: %+v", rows[0])
	}
	if rows[1].Category != "" || rows[1].BulkReason != bulkReasonNone {
		t.Fatalf("row 2 must carry the none marker: %+v", rows[1])
	}
	calls := fake.calls
	if n, _ = job.run(context.Background()); n != 0 || fake.calls != calls {
		t.Fatalf("second run must not fetch: n=%d calls=%d->%d", n, calls, fake.calls)
	}
}

func TestBackfillSkipsClassifiedAndStopsAtSixtyDays(t *testing.T) {
	job, fake, env := backfillEnv(t)
	addUID(t, env, "a@x.test", time.Hour, schemas.Email{IMAPUID: 1, Category: "bulk", BulkReason: "list-id"})
	addUID(t, env, "b@x.test", time.Hour, schemas.Email{IMAPUID: 2, BulkReason: bulkReasonNone})
	addUID(t, env, "c@x.test", 59*24*time.Hour, schemas.Email{IMAPUID: 3})
	addUID(t, env, "d@x.test", 61*24*time.Hour, schemas.Email{IMAPUID: 4})
	n, err := job.run(context.Background())
	if err != nil || n != 1 || fake.seen != 1 {
		t.Fatalf("only the 59 day row must be fetched: n=%d seen=%d err=%v", n, fake.seen, err)
	}
	var old schemas.Email
	testORM.First(&old, "imap_uid = ?", 4)
	if old.BulkReason != "" {
		t.Fatalf("row older than 60 days must stay untouched: %+v", old)
	}
}

func TestBackfillBatchesOfFifty(t *testing.T) {
	job, fake, env := backfillEnv(t)
	for i := 1; i <= 120; i++ {
		addUID(t, env, "p@x.test", time.Hour, schemas.Email{IMAPUID: uint32(i)})
	}
	n, err := job.run(context.Background())
	if err != nil || n != 120 || fake.calls != 3 {
		t.Fatalf("want 120 rows in 3 batches, got n=%d calls=%d err=%v", n, fake.calls, err)
	}
}
