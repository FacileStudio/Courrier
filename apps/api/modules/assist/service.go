package assist

import (
	"context"
	"errors"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/FacileStudio/Courrier/apps/api/internal/bulle"
	troncerrors "github.com/FacileStudio/tronc/errors"
)

const (
	maxText      = 20000
	maxContext   = 40000
	maxNote      = 500
	answerTokens = 1500
)

// Request is one helper call from the composer or a thread.
type Request struct {
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

// Status says whether the helper can run for this user and, when it cannot, why.
type Status struct {
	Available bool     `json:"available"`
	Reason    string   `json:"reason,omitempty"`
	Actions   []string `json:"actions"`
	BulleURL  string   `json:"bulle_url,omitempty"`
}

// Service runs helper actions through Bulle. A nil client means BULLE_URL or
// Registre is not configured, and every call answers "off".
type Service struct {
	bulle    *bulle.Client
	bulleURL string
}

// NewService builds the helper over a Bulle client, which may be nil, and the
// Bulle URL the client links to for "Ask Bulle".
func NewService(client *bulle.Client, bulleURL string) *Service {
	return &Service{bulle: client, bulleURL: bulleURL}
}

// Status probes Bulle for the user's model offer.
func (s *Service) Status(ctx context.Context, userID int64) Status {
	if s.bulle == nil {
		return Status{Reason: "off", Actions: Actions}
	}
	if _, err := s.bulle.Offer(ctx, userID); err != nil {
		return Status{Reason: reason(err), Actions: Actions, BulleURL: s.bulleURL}
	}
	return Status{Available: true, Actions: Actions, BulleURL: s.bulleURL}
}

// Run validates the request and returns Bulle's text.
func (s *Service) Run(ctx context.Context, userID int64, req Request) (Response, error) {
	if err := validate(&req); err != nil {
		return Response{}, err
	}
	if s.bulle == nil {
		return Response{}, troncerrors.Failed("off: the Bulle helper is not configured on this server")
	}
	offer, err := s.bulle.Offer(ctx, userID)
	if err != nil {
		return Response{}, failure(err)
	}
	system, user := buildPrompt(req, offer.Persona)
	text, err := s.bulle.Complete(ctx, userID, offer, bulle.Prompt{System: system, User: user, MaxTokens: answerTokens})
	if err != nil {
		return Response{}, failure(err)
	}
	return Response{Text: text}, nil
}

func validate(req *Request) error {
	req.Text, req.Context = strings.TrimSpace(req.Text), strings.TrimSpace(req.Context)
	req.Language, req.Instructions = strings.TrimSpace(req.Language), strings.TrimSpace(req.Instructions)
	switch {
	case !slices.Contains(Actions, req.Action):
		return troncerrors.Invalid("unknown action")
	case utf8.RuneCountInString(req.Text) > maxText, utf8.RuneCountInString(req.Context) > maxContext, utf8.RuneCountInString(req.Instructions) > maxNote:
		return troncerrors.TooLarge("text too long for the helper")
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
	switch {
	case errors.Is(err, bulle.ErrNotLinked):
		return "not_linked"
	case errors.Is(err, bulle.ErrNoKey):
		return "no_key"
	case errors.Is(err, bulle.ErrUnsupported):
		return "unsupported"
	}
	return "unreachable"
}

func failure(err error) error {
	r := reason(err)
	if r == "unreachable" {
		return troncerrors.Unavailable("unreachable: Bulle did not answer")
	}
	return troncerrors.Failed(r + ": " + err.Error())
}
