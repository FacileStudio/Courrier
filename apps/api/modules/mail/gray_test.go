package mail

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
)

func fakeJev(t *testing.T) *jevClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			State grayState `json:"state"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		answer := jevAnswer{Choice: "primary", Confidence: 0.95}
		switch {
		case strings.Contains(body.State.Subject, "promo"):
			answer = jevAnswer{Choice: "newsletter", Confidence: 0.92}
		case strings.Contains(body.State.Subject, "maybe"):
			answer = jevAnswer{Choice: "newsletter", Confidence: 0.55}
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"answers": map[string]jevAnswer{grayQuestion: answer}})
	}))
	t.Cleanup(server.Close)
	return &jevClient{url: server.URL, apiKey: "k", http: server.Client()}
}

func TestGrayZoneMovesOnlyConfidentStrangers(t *testing.T) {
	env := newTriageEnv(t)
	testORM.Exec("DELETE FROM sender_rules")
	env.svc.judge = fakeJev(t)
	testORM.Create(&schemas.Contact{AccountID: env.a.ID, Email: "friend@x.test", SentCount: 3})
	now := time.Now()
	none := schemas.Email{BulkReason: bulkReasonNone}
	seed := map[string]string{"promo week": "shop@x.test", "maybe news": "blog@x.test", "lunch?": "anna@x.test", "promo from a friend": "friend@x.test"}
	for subject, from := range seed {
		e := none
		e.Subject = subject
		addMail(t, mailKey{env.a.ID, env.inboxA, subject}, from, now, e)
	}
	env.svc.judgeGrayZone(env.inboxA)

	got := map[string]string{}
	var rows []schemas.Email
	testORM.Where("folder_id = ?", env.inboxA).Find(&rows)
	for _, e := range rows {
		got[e.Subject] = e.Category + "/" + e.BulkReason
	}
	want := map[string]string{
		"promo week":          "bulk/jev",
		"maybe news":          "/judged",
		"lunch?":              "/judged",
		"promo from a friend": "/none",
	}
	for subject, w := range want {
		if got[subject] != w {
			t.Errorf("%s = %q, want %q", subject, got[subject], w)
		}
	}
}
