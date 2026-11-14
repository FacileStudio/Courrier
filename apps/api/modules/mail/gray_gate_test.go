package mail

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/tronc/bulle"
)

func TestTriageGateBacksOffAndPausesUntilTomorrow(t *testing.T) {
	now := time.Date(2026, 11, 14, 10, 0, 0, 0, time.UTC)
	gate := &triageGate{now: func() time.Time { return now }}
	steps := []struct {
		after time.Duration
		want  bool
		err   error
	}{
		{0, true, bulle.ErrUnreachable},
		{4 * time.Minute, false, nil},
		{time.Minute, true, bulle.ErrUnreachable},
		{9 * time.Minute, false, nil},
		{time.Minute, true, nil},
		{0, true, &bulle.StatusError{Status: 429}},
		{12 * time.Hour, false, nil},
		{2 * time.Hour, true, nil},
	}
	for i, step := range steps {
		now = now.Add(step.after)
		got := gate.take(7)
		if got != step.want {
			t.Fatalf("step %d at %s: take = %v, want %v", i, now.Format(time.TimeOnly), got, step.want)
		}
		if got {
			gate.settle(7, step.err)
		}
	}
}

func TestTriageGateKeepsADailyBudget(t *testing.T) {
	now := time.Date(2026, 11, 14, 23, 0, 0, 0, time.UTC)
	gate := &triageGate{now: func() time.Time { return now }}
	for i := range grayDailyCalls {
		if !gate.take(7) {
			t.Fatalf("call %d refused inside the budget", i)
		}
		gate.settle(7, nil)
	}
	if gate.take(7) {
		t.Fatal("a call past the daily budget went through")
	}
	if !gate.take(8) {
		t.Fatal("another owner shares the budget")
	}
	now = now.Add(time.Hour)
	if !gate.take(7) {
		t.Fatal("the budget did not reset on the next UTC day")
	}
}

func TestBatchesHoldOneSenderDomain(t *testing.T) {
	emails := []schemas.Email{{FromAddress: "a@Shop.test"}, {FromAddress: "b@evil.test"}, {FromAddress: "c@shop.test"}, {FromAddress: "d@shop.test"}}
	got := sameSenderDomain(emails, 2)
	if len(got) != 2 || got[0].FromAddress != "a@Shop.test" || got[1].FromAddress != "c@shop.test" {
		t.Fatalf("batch = %+v", got)
	}
}

func TestGrayInputStaysUnderBulleBodyOnceEncoded(t *testing.T) {
	emails := make([]schemas.Email, grayBatch)
	for i := range emails {
		emails[i] = schemas.Email{FromAddress: "x@x.test", Subject: "s", BodyText: strings.Repeat(`"<`, grayBodyBytes)}
	}
	input, err := grayInput(emails)
	if err != nil {
		t.Fatal(err)
	}
	request, _ := json.Marshal(bulle.CompleteRequest{System: graySystem, Prompt: grayPrompt, Input: input, Schema: graySchema, MaxTokens: grayTokens})
	if len(input) > grayInputBytes || len(request) > 96<<10 || len(input) < 2 {
		t.Fatalf("input %d bytes, request %d bytes", len(input), len(request))
	}
}

func TestGrayZoneFailureReachesBulleOnce(t *testing.T) {
	env := newTriageEnv(t)
	seedGray(t, env)
	down := &fakeBulle{err: bulle.ErrUnreachable}
	env.svc.WithTriage(Triage{Bulle: down, AIAccess: func(context.Context, int64) (bool, error) { return true, nil }})
	env.svc.judgeGrayZone(env.inboxA)
	env.svc.judgeGrayZone(env.inboxA)
	if down.calls != 1 {
		t.Fatalf("a failing Bulle was called %d times in a row, want 1", down.calls)
	}
}

func TestSpaceMailboxTriageNeedsAMember(t *testing.T) {
	env := newTriageEnv(t)
	seedGray(t, env)
	testORM.Exec("DELETE FROM users WHERE email LIKE 'triage-owner-%'")
	space := schemas.Space{Name: "Triage " + t.Name()}
	owner := schemas.User{Email: "triage-owner-" + t.Name() + "@facile.studio"}
	if testORM.Create(&space).Error != nil || testORM.Create(&owner).Error != nil {
		t.Fatal("seed space")
	}
	testORM.Model(&schemas.Account{}).Where("id = ?", env.a.ID).Updates(map[string]any{"space_id": space.ID, "user_id": owner.ID})
	fake := &fakeBulle{}
	env.svc.WithTriage(Triage{Bulle: fake, AIAccess: func(context.Context, int64) (bool, error) { return true, nil }})
	env.svc.judgeGrayZone(env.inboxA)
	if fake.calls != 0 {
		t.Fatalf("a former member's mailbox reached Bulle %d times", fake.calls)
	}
	testORM.Create(&schemas.SpaceMember{SpaceID: space.ID, UserID: owner.ID})
	env.svc.judgeGrayZone(env.inboxA)
	if fake.calls != 1 || fake.userID != owner.ID {
		t.Fatalf("member: calls = %d as user %d", fake.calls, fake.userID)
	}
}
