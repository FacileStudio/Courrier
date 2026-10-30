package mail

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/FacileStudio/tronc/errors"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
)

func TestTriageClassOrder(t *testing.T) {
	env := newTriageEnv(t)
	seedClassOrder(t, env)

	resp, err := env.svc.Triage(context.Background(), 1, TriageQuery{Page: 1, Limit: 50})
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{
		"t1": "bulk/noreply-sender", "t2": "known/in-sent-thread", "t3": "known/sent-to",
		"t4": "unknown_human/no-signal", "t5": "bulk/list-unsubscribe",
	}
	got := classOf(resp)
	if len(got) != len(want) {
		t.Fatalf("want %d items, got %v", len(want), got)
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("%s: want %s, got %s", k, v, got[k])
		}
	}
	if resp.Counts["known"] != 2 || resp.Counts["bulk"] != 2 || resp.Counts["unknown_human"] != 1 {
		t.Fatalf("counts wrong: %v", resp.Counts)
	}
}

func TestTriageCountsCoverFullSetAcrossPages(t *testing.T) {
	env := newTriageEnv(t)
	now := time.Now()
	for i := 0; i < 7; i++ {
		from := fmt.Sprintf("person%d@x.test", i)
		if i%2 == 0 {
			from = fmt.Sprintf("noreply@s%d.test", i)
		}
		addMail(t, mailKey{env.a.ID, env.inboxA, fmt.Sprintf("p%d", i)}, from, now.Add(-time.Duration(i)*time.Minute), schemas.Email{})
	}
	seen := map[string]bool{}
	for page := 1; page <= 3; page++ {
		resp, err := env.svc.Triage(context.Background(), 1, TriageQuery{Page: page, Limit: 3})
		if err != nil {
			t.Fatal(err)
		}
		if resp.Total != 7 || resp.Counts["bulk"] != 4 || resp.Counts["unknown_human"] != 3 || resp.Counts["known"] != 0 {
			t.Fatalf("page %d total/counts wrong: %d %v", page, resp.Total, resp.Counts)
		}
		for _, it := range resp.Items {
			seen[it.ThreadID] = true
		}
	}
	if len(seen) != 7 {
		t.Fatalf("pages must cover 7 distinct threads, got %d", len(seen))
	}
	only, _ := env.svc.Triage(context.Background(), 1, TriageQuery{Page: 1, Limit: 50, Classes: map[string]bool{"unknown_human": true}})
	if only.Total != 3 || len(only.Items) != 3 || only.Counts["bulk"] != 4 {
		t.Fatalf("class filter must narrow items but not counts: %d %v", only.Total, only.Counts)
	}
}

func TestTriageSinceUnreadAndAccount(t *testing.T) {
	env := newTriageEnv(t)
	now := time.Now()
	addMail(t, mailKey{env.a.ID, env.inboxA, "new"}, "p@x.test", now.Add(-10*time.Minute), schemas.Email{})
	addMail(t, mailKey{env.a.ID, env.inboxA, "old"}, "q@x.test", now.Add(-3*24*time.Hour), schemas.Email{})
	addMail(t, mailKey{env.a.ID, env.inboxA, "read"}, "r@x.test", now.Add(-5*time.Minute), schemas.Email{IsRead: true})
	since, _ := parseSince("2h", now)
	resp, _ := env.svc.Triage(context.Background(), 1, TriageQuery{Page: 1, Limit: 50, Since: since})
	if resp.Total != 2 {
		t.Fatalf("since 2h must keep new and read, got %d", resp.Total)
	}
	resp, _ = env.svc.Triage(context.Background(), 1, TriageQuery{Page: 1, Limit: 50, Since: since, UnreadOnly: true})
	if resp.Total != 1 || resp.Items[0].ThreadID != "new" {
		t.Fatalf("unread since 2h must keep only new: %+v", resp.Items)
	}
	resp, _ = env.svc.Triage(context.Background(), 1, TriageQuery{Page: 1, Limit: 50, AccountID: &env.b.ID})
	if resp.Total != 0 || len(resp.Items) != 0 || len(resp.Counts) != 3 {
		t.Fatalf("account b has no mail but all count keys: %+v", resp)
	}
}

func TestTriageBadAccountIs404(t *testing.T) {
	env := newTriageEnv(t)
	_, err := env.svc.Triage(context.Background(), 1, TriageQuery{Page: 1, Limit: 50, AccountID: &env.c.ID})
	if err == nil || errors.Status(err) != 404 {
		t.Fatalf("another user's account must be 404, got %v", err)
	}
	missing := int64(999999)
	if _, err := env.svc.Triage(context.Background(), 1, TriageQuery{Page: 1, Limit: 50, AccountID: &missing}); errors.Status(err) != 404 {
		t.Fatalf("unknown account must be 404, got %v", err)
	}
}

func TestParseSinceGrammar(t *testing.T) {
	now := time.Date(2026, 10, 30, 12, 0, 0, 0, time.UTC)
	cases := map[string]time.Time{
		"":                     {},
		"30m":                  now.Add(-30 * time.Minute),
		"2h":                   now.Add(-2 * time.Hour),
		"1d":                   now.Add(-24 * time.Hour),
		"1w":                   now.Add(-7 * 24 * time.Hour),
		"2026-10-29T10:00:00Z": time.Date(2026, 10, 29, 10, 0, 0, 0, time.UTC),
	}
	for in, want := range cases {
		got, err := parseSince(in, now)
		if err != nil || !got.Equal(want) {
			t.Errorf("%q: want %v, got %v (%v)", in, want, got, err)
		}
	}
	for _, bad := range []string{"2x", "h", "-1h", "yesterday"} {
		if _, err := parseSince(bad, now); err == nil {
			t.Errorf("%q must be rejected", bad)
		}
	}
}
