package assist

import (
	"context"
	stderrors "errors"
	"strings"
	"testing"

	"github.com/FacileStudio/porte/exchange"
	"github.com/FacileStudio/tronc/bulle"
	troncerrors "github.com/FacileStudio/tronc/errors"
)

func TestValidateAndPrompt(t *testing.T) {
	cases := map[string]Request{
		"unknown action":       {Action: "delete", Text: "x"},
		"empty text":           {Action: "shorten"},
		"translate, no target": {Action: "translate", Text: "x"},
		"reply, no thread":     {Action: "reply"},
	}
	for name, req := range cases {
		if err := validate(&req); err == nil {
			t.Errorf("%s: want an error", name)
		}
	}
	req := Request{Action: "translate", Text: "Bonjour", Language: "English", Instructions: "keep it short"}
	if err := validate(&req); err != nil {
		t.Fatal(err)
	}
	call := buildPrompt(req)
	if call.System != baseSystem || !strings.Contains(call.Prompt, "into English") || !strings.Contains(call.Prompt, "keep it short") ||
		call.Input != "<text>\nBonjour\n</text>" || call.SpaceID != "" {
		t.Fatalf("call = %+v", call)
	}
}

func TestOffServiceSaysSo(t *testing.T) {
	s := NewService(Config{})
	if st := s.Status(context.Background(), 1, 0); st.Available || st.Reason != "off" {
		t.Fatalf("status = %+v", st)
	}
	if _, err := s.Run(context.Background(), 1, Request{Action: "fix", Text: "x"}); err == nil || !strings.Contains(err.Error(), "off") {
		t.Fatalf("run = %v", err)
	}
}

type fakeBulle struct {
	calls   int
	spaceID string
	out     bulle.Completion
	err     error
}

func (f *fakeBulle) Complete(_ context.Context, req bulle.CompleteRequest) (bulle.Completion, error) {
	f.calls++
	f.spaceID = req.SpaceID
	return f.out, f.err
}

func service(fake *fakeBulle, access bool, linked error) *Service {
	return NewService(Config{
		Bulle:      fake,
		Linked:     func(context.Context) (string, error) { return "token", linked },
		AIAccess:   func(context.Context, int64, string) (bool, error) { return access, nil },
		UpgradeURL: "https://sso.example.test/spaces",
	})
}

func TestStatusFollowsThePlanAndTheLink(t *testing.T) {
	ctx := context.Background()
	free := service(&fakeBulle{}, false, nil).Status(ctx, 1, 0)
	if free.Available || free.Reason != "plan_required" || free.UpgradeURL != "https://sso.example.test/spaces" {
		t.Fatalf("free = %+v", free)
	}
	unlinked := service(&fakeBulle{}, true, exchange.ErrReconnect).Status(ctx, 1, 0)
	if unlinked.Available || unlinked.Reason != "not_linked" {
		t.Fatalf("unlinked = %+v", unlinked)
	}
	if paid := service(&fakeBulle{}, true, nil).Status(ctx, 1, 0); !paid.Available {
		t.Fatalf("paid = %+v", paid)
	}
}

func TestRunMapsBulleAnswers(t *testing.T) {
	ctx := context.Background()
	req := Request{Action: "fix", Text: "helo"}
	free := &fakeBulle{}
	if _, err := service(free, false, nil).Run(ctx, 1, req); !isPlanError(err, "https://sso.example.test/spaces") || free.calls != 0 {
		t.Fatalf("free: err = %v, calls = %d", err, free.calls)
	}
	refused := &fakeBulle{err: &bulle.PlanRequired{Message: "Upgrade", UpgradeURL: "https://bulle.example.test/plans"}}
	if _, err := service(refused, true, nil).Run(ctx, 1, req); !isPlanError(err, "https://bulle.example.test/plans") {
		t.Fatalf("refused: err = %v", err)
	}
	capped := &fakeBulle{err: &bulle.StatusError{Status: 429, Message: "at most 200 completions per day"}}
	if _, err := service(capped, true, nil).Run(ctx, 1, req); codeOf(err) != "rate_limited" {
		t.Fatalf("capped: err = %v", err)
	}
	down := &fakeBulle{err: bulle.ErrUnreachable}
	if _, err := service(down, true, nil).Run(ctx, 1, req); codeOf(err) != "unavailable" {
		t.Fatalf("down: err = %v", err)
	}
	ok := &fakeBulle{out: bulle.Completion{Text: " hello \n"}}
	if out, err := service(ok, true, nil).Run(ctx, 1, req); err != nil || out.Text != "hello" {
		t.Fatalf("ok: out = %+v, err = %v", out, err)
	}
}

func isPlanError(err error, upgradeURL string) bool {
	var envelope *troncerrors.Error
	return stderrors.As(err, &envelope) && envelope.Code == "payment_required" && envelope.UpgradeURL == upgradeURL
}

func codeOf(err error) string {
	var envelope *troncerrors.Error
	if stderrors.As(err, &envelope) {
		return envelope.Code
	}
	return ""
}

func TestRunMeasuresTheEncodedRequest(t *testing.T) {
	quotes := &fakeBulle{}
	req := Request{Action: "reply", Context: strings.Repeat(`"`, 39000), Text: strings.Repeat("<", 19000)}
	if len(buildPrompt(req).Input) > maxInput {
		t.Fatalf("the input alone is already past the cap")
	}
	if _, err := service(quotes, true, nil).Run(context.Background(), 1, req); codeOf(err) != "resource_exhausted" || quotes.calls != 0 {
		t.Fatalf("err = %v, calls = %d, want a local 413", err, quotes.calls)
	}
	huge := &fakeBulle{err: &bulle.StatusError{Status: 413, Message: "request body too large"}}
	if _, err := service(huge, true, nil).Run(context.Background(), 1, Request{Action: "fix", Text: "helo"}); codeOf(err) != "resource_exhausted" {
		t.Fatalf("huge: err = %v", err)
	}
}

func TestSpaceMailboxRunsOnTheSpacePlan(t *testing.T) {
	ctx := context.Background()
	fake := &fakeBulle{out: bulle.Completion{Text: "hello"}}
	s := NewService(Config{
		Bulle:  fake,
		Linked: func(context.Context) (string, error) { return "token", nil },
		Mailbox: func(_ context.Context, _ int64, accountID int64) (string, error) {
			switch accountID {
			case 7:
				return "reg-space", nil
			case 8:
				return "", nil
			}
			return "", troncerrors.NotFound("account not found")
		},
		AIAccess: func(_ context.Context, _ int64, link string) (bool, error) { return link == "reg-space", nil },
	})
	if _, err := s.Run(ctx, 1, Request{AccountID: 7, Action: "fix", Text: "helo"}); err != nil || fake.calls != 1 || fake.spaceID != "reg-space" {
		t.Fatalf("space mailbox: err = %v, calls = %d, space = %q", err, fake.calls, fake.spaceID)
	}
	if st := s.Status(ctx, 1, 7); !st.Available {
		t.Fatalf("space mailbox status = %+v", st)
	}
	for _, account := range []int64{0, 8} {
		if _, err := s.Run(ctx, 1, Request{AccountID: account, Action: "fix", Text: "helo"}); codeOf(err) != "payment_required" || fake.calls != 1 {
			t.Fatalf("personal mailbox %d: err = %v, calls = %d", account, err, fake.calls)
		}
	}
	if _, err := s.Run(ctx, 1, Request{AccountID: 9, Action: "fix", Text: "helo"}); codeOf(err) != "not_found" || fake.calls != 1 {
		t.Fatalf("foreign mailbox: err = %v, calls = %d", err, fake.calls)
	}
}
