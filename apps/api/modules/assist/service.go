package assist

import (
	"context"
	stderrors "errors"
	"net/http"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/FacileStudio/porte/exchange"
	"github.com/FacileStudio/tronc/bulle"
	troncerrors "github.com/FacileStudio/tronc/errors"
)

const (
	maxText      = 20000
	maxContext   = 40000
	maxNote      = 500
	maxInput     = 64 << 10
	maxBody      = 96 << 10
	answerTokens = 1500
	planMessage  = "The Bulle helper needs a plan with AI access."
	tooLong      = "text too long for the helper"
)

// Request is one helper call from the composer or a thread. AccountID names
// the mailbox the text comes from, so a space mailbox runs on the space's
// plan when the user's own has no AI.
type Request struct {
	AccountID    int64  `json:"account_id,omitempty"`
	Action       string `json:"action"`
	Text         string `json:"text,omitempty"`
	Context      string `json:"context,omitempty"`
	Language     string `json:"language,omitempty"`
	Instructions string `json:"instructions,omitempty"`
}

// Response is the helper's text.
type Response struct {
	Text string `json:"text"`
}

// Status says whether the helper can run for this user and, when it cannot,
// why. UpgradeURL is where to get a plan with AI access, set with the
// plan_required reason.
type Status struct {
	Available  bool     `json:"available"`
	Reason     string   `json:"reason,omitempty"`
	Actions    []string `json:"actions"`
	BulleURL   string   `json:"bulle_url,omitempty"`
	UpgradeURL string   `json:"upgrade_url,omitempty"`
}

// Completer is Bulle's completion call, as tronc/bulle's Client makes it for
// the user on the context.
type Completer interface {
	Complete(ctx context.Context, req bulle.CompleteRequest) (bulle.Completion, error)
}

// Config wires the helper. A nil Bulle means BULLE_URL or Registre is not
// configured, and every call answers "off". Linked returns the user's Bulle
// token, or why there is none. Mailbox checks the user may use a mail
// account and returns its space's Registre id, empty for a personal one;
// AIAccess reads whether the user's plan, or that space's, includes AI.
type Config struct {
	Bulle      Completer
	Linked     bulle.TokenSource
	Mailbox    func(ctx context.Context, userID, accountID int64) (string, error)
	AIAccess   func(ctx context.Context, userID int64, link string) (bool, error)
	BulleURL   string
	UpgradeURL string
}

// Service runs helper actions through Bulle's /api/ai/complete.
type Service struct {
	cfg Config
}

// NewService builds the helper over cfg.
func NewService(cfg Config) *Service {
	return &Service{cfg: cfg}
}

// Status checks the plan of the user, or of the space of mail account
// accountID (0 for none), then that Courrier can act for them in Bulle.
func (s *Service) Status(ctx context.Context, userID, accountID int64) Status {
	if s.cfg.Bulle == nil {
		return Status{Reason: "off", Actions: Actions}
	}
	status := Status{Actions: Actions, BulleURL: s.cfg.BulleURL}
	if _, err := s.admit(ctx, userID, accountID); err != nil {
		status.Reason = reason(err)
		if status.Reason == "plan_required" {
			status.UpgradeURL = s.cfg.UpgradeURL
		}
		return status
	}
	if _, err := s.cfg.Linked(ctx); err != nil {
		status.Reason = reason(err)
		return status
	}
	status.Available = true
	return status
}

// Run validates the request and returns Bulle's text.
func (s *Service) Run(ctx context.Context, userID int64, req Request) (Response, error) {
	if err := validate(&req); err != nil {
		return Response{}, err
	}
	if s.cfg.Bulle == nil {
		return Response{}, troncerrors.Failed("off: the Bulle helper is not configured on this server")
	}
	link, err := s.admit(ctx, userID, req.AccountID)
	if err != nil {
		return Response{}, s.failure(err)
	}
	call := buildPrompt(req)
	call.SpaceID = link
	if tooLarge(call) {
		return Response{}, troncerrors.TooLarge(tooLong)
	}
	out, err := s.cfg.Bulle.Complete(ctx, call)
	if err != nil {
		return Response{}, s.failure(err)
	}
	text := strings.TrimSpace(out.Text)
	if text == "" {
		return Response{}, troncerrors.Unavailable("unreachable: Bulle did not answer")
	}
	return Response{Text: text}, nil
}

// admit refuses, before any mail leaves for Bulle, a user who may not use
// mail account accountID, or whose plan and whose mailbox's space plan both
// lack ai.access. It returns the space's Registre id Bulle checks.
func (s *Service) admit(ctx context.Context, userID, accountID int64) (string, error) {
	link := ""
	if accountID != 0 && s.cfg.Mailbox != nil {
		var err error
		if link, err = s.cfg.Mailbox(ctx, userID, accountID); err != nil {
			return "", err
		}
	}
	allowed, err := s.cfg.AIAccess(ctx, userID, link)
	if err != nil {
		return "", err
	}
	if !allowed {
		return "", &bulle.PlanRequired{Message: planMessage}
	}
	return link, nil
}

func validate(req *Request) error {
	req.Text, req.Context = strings.TrimSpace(req.Text), strings.TrimSpace(req.Context)
	req.Language, req.Instructions = strings.TrimSpace(req.Language), strings.TrimSpace(req.Instructions)
	switch {
	case !slices.Contains(Actions, req.Action):
		return troncerrors.Invalid("unknown action")
	case utf8.RuneCountInString(req.Text) > maxText, utf8.RuneCountInString(req.Context) > maxContext, utf8.RuneCountInString(req.Instructions) > maxNote:
		return troncerrors.TooLarge(tooLong)
	case req.Action == "translate" && (req.Language == "" || len(req.Language) > 40):
		return troncerrors.Invalid("translate needs a language")
	case (req.Action == "reply" || req.Action == "summarize") && req.Context == "":
		return troncerrors.Invalid(req.Action + " needs the thread as context")
	case req.Action != "reply" && req.Action != "summarize" && req.Text == "":
		return troncerrors.Invalid("text is empty")
	}
	return nil
}

func reason(err error) string {
	var plan *bulle.PlanRequired
	var refused *exchange.RefusedError
	switch {
	case stderrors.As(err, &plan):
		return "plan_required"
	case stderrors.Is(err, exchange.ErrReconnect), stderrors.Is(err, exchange.ErrNoUser), stderrors.As(err, &refused):
		return "not_linked"
	}
	return "unreachable"
}

// failure maps a Bulle or exchange error to the answer the client shows: a
// plan refusal as a 402 with the upgrade link, a dead Registre grant as the
// reconnect message, the daily cap as a 429, a body Bulle found too large
// as a 413.
func (s *Service) failure(err error) error {
	var plan *bulle.PlanRequired
	var refused *bulle.StatusError
	var envelope *troncerrors.Error
	switch {
	case stderrors.As(err, &plan):
		return troncerrors.PaymentRequired(orDefault(plan.Message, planMessage)).WithUpgradeURL(orDefault(plan.UpgradeURL, s.cfg.UpgradeURL))
	case stderrors.As(err, &refused) && refused.Status == http.StatusTooManyRequests:
		return troncerrors.RateLimited(orDefault(refused.Message, "Bulle's daily limit is reached. Try again tomorrow."))
	case stderrors.As(err, &refused) && refused.Status == http.StatusRequestEntityTooLarge:
		return troncerrors.TooLarge(tooLong)
	case reason(err) == "not_linked":
		return troncerrors.Failed("not_linked: " + exchange.ErrReconnect.Message)
	case stderrors.As(err, &envelope):
		return envelope
	}
	return troncerrors.Unavailable("unreachable: Bulle did not answer")
}

func orDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
