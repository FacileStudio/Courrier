package mail

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/tronc/bulle"
)

type fakeBulle struct {
	calls   int
	userID  int64
	spaceID string
	err     error
}

func (f *fakeBulle) Complete(ctx context.Context, req bulle.CompleteRequest) (bulle.Completion, error) {
	f.calls++
	f.spaceID = req.SpaceID
	if identity, ok := porte.From(ctx); ok {
		f.userID = identity.UserID
	}
	if f.err != nil {
		return bulle.Completion{}, f.err
	}
	var emails []grayEmail
	if err := json.Unmarshal([]byte(req.Input), &emails); err != nil || len(req.Schema) == 0 {
		return bulle.Completion{}, bulle.ErrUnreachable
	}
	verdicts := []grayVerdict{}
	for _, email := range emails {
		verdict := grayVerdict{N: email.N, Tab: "primary", Confidence: 0.95}
		switch {
		case strings.Contains(email.Subject, "promo"):
			verdict = grayVerdict{N: email.N, Tab: "newsletter", Confidence: 0.92}
		case strings.Contains(email.Subject, "maybe"):
			verdict = grayVerdict{N: email.N, Tab: "newsletter", Confidence: 0.55}
		}
		verdicts = append(verdicts, verdict)
	}
	answer, _ := json.Marshal(map[string]any{"verdicts": verdicts})
	return bulle.Completion{JSON: answer}, nil
}

func seedGray(t *testing.T, env triageEnv) {
	t.Helper()
	testORM.Exec("DELETE FROM sender_rules")
	testORM.Create(&schemas.Contact{AccountID: env.a.ID, Email: "friend@x.test", SentCount: 3})
	now := time.Now()
	none := schemas.Email{BulkReason: bulkReasonNone}
	seed := map[string]string{"promo week": "shop@x.test", "maybe news": "blog@x.test", "lunch?": "anna@x.test", "promo from a friend": "friend@x.test"}
	for subject, from := range seed {
		e := none
		e.Subject = subject
		addMail(t, mailKey{env.a.ID, env.inboxA, subject}, from, now, e)
	}
}

func grayState(env triageEnv) map[string]string {
	got := map[string]string{}
	var rows []schemas.Email
	testORM.Where("folder_id = ?", env.inboxA).Find(&rows)
	for _, e := range rows {
		got[e.Subject] = e.Category + "/" + e.BulkReason
	}
	return got
}

func TestGrayZoneMovesOnlyConfidentStrangers(t *testing.T) {
	env := newTriageEnv(t)
	seedGray(t, env)
	fake := &fakeBulle{}
	env.svc.WithTriage(Triage{Bulle: fake, AIAccess: func(context.Context, int64, string) (bool, error) { return true, nil }})
	env.svc.judgeGrayZone(env.inboxA)

	got := grayState(env)
	want := map[string]string{
		"promo week":          "bulk/bulle",
		"maybe news":          "/judged",
		"lunch?":              "/judged",
		"promo from a friend": "/none",
	}
	for subject, w := range want {
		if got[subject] != w {
			t.Errorf("%s = %q, want %q", subject, got[subject], w)
		}
	}
	if fake.calls != 1 || fake.userID != env.a.UserID || fake.spaceID != "" {
		t.Errorf("calls = %d as user %d in space %q, want one personal call as user %d", fake.calls, fake.userID, fake.spaceID, env.a.UserID)
	}
}

func TestGrayZoneSkipsWithoutAIAccessOrAnswer(t *testing.T) {
	env := newTriageEnv(t)
	seedGray(t, env)
	free := &fakeBulle{}
	env.svc.WithTriage(Triage{Bulle: free, AIAccess: func(context.Context, int64, string) (bool, error) { return false, nil }})
	env.svc.judgeGrayZone(env.inboxA)
	if free.calls != 0 {
		t.Fatalf("a free plan reached Bulle %d times", free.calls)
	}
	down := &fakeBulle{err: &bulle.PlanRequired{Message: "no"}}
	env.svc.WithTriage(Triage{Bulle: down, AIAccess: func(context.Context, int64, string) (bool, error) { return true, nil }})
	env.svc.judgeGrayZone(env.inboxA)
	for subject, state := range grayState(env) {
		if state != "/none" {
			t.Errorf("%s = %q, want it left for the next pass", subject, state)
		}
	}
}

func TestGrayInputStaysUnderBulleCap(t *testing.T) {
	body := strings.Repeat("é", grayBodyBytes)
	emails := make([]schemas.Email, 200)
	for i := range emails {
		emails[i] = schemas.Email{FromAddress: "x@x.test", Subject: "s", BodyText: body}
	}
	input, err := grayInput(emails)
	if err != nil || len(input) > grayInputBytes || !json.Valid([]byte(input)) {
		t.Fatalf("input is %d bytes, err %v", len(input), err)
	}
}
