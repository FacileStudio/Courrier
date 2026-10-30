package mail

import "testing"

func TestClassify(t *testing.T) {
	cases := []struct {
		name     string
		signals  Signals
		from     string
		category string
		reason   string
	}{
		{"railway", Signals{ListUnsubscribe: "<https://railway.app/unsub>", ListID: "<updates.railway.app>"}, "hello@notify.railway.app", "bulk", "list-unsubscribe"},
		{"ouigo", Signals{ListUnsubscribe: "<mailto:unsub@ouigo.com>", Precedence: "bulk"}, "ne-pas-repondre@ouigo.com", "bulk", "list-unsubscribe"},
		{"list id only", Signals{ListID: "<dev.lists.example.org>"}, "alice@example.org", "bulk", "list-id"},
		{"precedence list", Signals{Precedence: "list"}, "alice@example.org", "bulk", "precedence-bulk"},
		{"precedence junk", Signals{Precedence: "junk"}, "alice@example.org", "bulk", "precedence-bulk"},
		{"auto replied", Signals{AutoSubmitted: "auto-replied"}, "alice@example.org", "auto", "auto-submitted"},
		{"auto-submitted no", Signals{AutoSubmitted: "no"}, "alice@example.org", "", ""},
		{"noreply sender", Signals{}, "No-Reply@bank.example", "bulk", "noreply-sender"},
		{"mailer-daemon", Signals{}, "MAILER-DAEMON@mx.example", "bulk", "noreply-sender"},
		{"header beats sender", Signals{ListID: "<x.example>"}, "noreply@x.example", "bulk", "list-id"},
		{"plain human", Signals{}, "flavie@example.fr", "", ""},
		{"human with news prefix", Signals{}, "newsome@example.fr", "", ""},
		{"no at sign", Signals{}, "weird", "", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			category, reason := Classify(c.signals, c.from)
			if category != c.category || reason != c.reason {
				t.Fatalf("got (%q, %q), want (%q, %q)", category, reason, c.category, c.reason)
			}
		})
	}
}
