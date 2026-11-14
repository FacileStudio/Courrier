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
	if st := s.Status(context.Background(), 1); st.Available || st.Reason != "off" {
		t.Fatalf("status = %+v", st)
	}
	if _, err := s.Run(context.Background(), 1, Request{Action: "fix", Text: "x"}); err == nil || !strings.Contains(err.Error(), "off") {
		t.Fatalf("run = %v", err)
	}
}

type fakeBulle struct {
	calls int
	out   bulle.Completion
	err   error
}

func (f *fakeBulle) Complete(context.Context, bulle.CompleteRequest) (bulle.Completion, error) {
	f.calls++
	return f.out, f.err
}

func service(fake *fakeBulle, access bool, linked error) *Service {
	return NewService(Config{
		Bulle:      fake,
		Linked:     func(context.Context) (string, error) { return "token", linked },
		AIAccess:   func(context.Context, int64) (bool, error) { return access, nil },
		UpgradeURL: "https://sso.example.test/spaces",
	})
}

func TestStatusFollowsThePlanAndTheLink(t *testing.T) {
	ctx := context.Background()
	free := service(&fakeBulle{}, false, nil).Status(ctx, 1)
	if free.Available || free.Reason != "plan_required" || free.UpgradeURL != "https://sso.example.test/spaces" {
		t.Fatalf("free = %+v", free)
	}
	unlinked := service(&fakeBulle{}, true, exchange.ErrReconnect).Status(ctx, 1)
	if unlinked.Available || unlinked.Reason != "not_linked" {
		t.Fatalf("unlinked = %+v", unlinked)
	}
	if paid := service(&fakeBulle{}, true, nil).Status(ctx, 1); !paid.Available {
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
