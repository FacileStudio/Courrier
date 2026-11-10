package mcp

import (
	"context"
	"strings"
	"testing"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// compose is a send that writes its own message.
var compose = func(accountID int64) map[string]any {
	return map[string]any{"account_id": accountID, "to": []string{"client@example.test"}, "subject": "Quote", "body": "Here it is."}
}

// answering returns an elicitation handler that answers action, and counts
// the questions.
func answering(action string, asked *[]string) func(context.Context, *sdk.ElicitRequest) (*sdk.ElicitResult, error) {
	return func(_ context.Context, req *sdk.ElicitRequest) (*sdk.ElicitResult, error) {
		*asked = append(*asked, req.Params.Message)
		return &sdk.ElicitResult{Action: action, Content: map[string]any{"confirm": action == "accept"}}, nil
	}
}

// TestConfirmedSendGoesOutOnce checks the MRTR round trip: the client asks,
// the person accepts, the mail goes out once and the draft is gone.
func TestConfirmedSendGoesOutOnce(t *testing.T) {
	needDB(t)
	f := newFixture(t)
	var asked []string
	session := f.connect(t, f.token(t), clientOptions{elicit: answering("accept", &asked)})
	var out actionOutput
	decode(t, call(t, session, toolSend, compose(f.accountID)), &out)
	if out.Status != "sent" || f.mail.sentCount() != 1 || len(asked) != 1 || !strings.Contains(asked[0], `"Quote"`) {
		t.Fatalf("want one confirmed send, got status %q, %d sent, questions %q", out.Status, f.mail.sentCount(), asked)
	}
}

// TestDeclinedSendKeepsTheDraft checks a refusal sends nothing.
func TestDeclinedSendKeepsTheDraft(t *testing.T) {
	needDB(t)
	f := newFixture(t)
	var asked []string
	var out actionOutput
	decode(t, call(t, f.connect(t, f.token(t), clientOptions{elicit: answering("decline", &asked)}), toolSend, compose(f.accountID)), &out)
	if out.Status != outcomeDeclined || f.mail.sentCount() != 0 {
		t.Fatalf("want declined and nothing sent, got %q and %d sent", out.Status, f.mail.sentCount())
	}
}

// TestWithoutElicitationNothingIsSent checks the fallback: a draft and a
// confirm link, and nothing sent, for a client that cannot ask and for a
// pre-2026-07-28 client, which a stateless server cannot ask either.
func TestWithoutElicitationNothingIsSent(t *testing.T) {
	needDB(t)
	f := newFixture(t)
	var asked []string
	clients := map[string]clientOptions{
		"no elicitation": {},
		"2025-11-25":     {elicit: answering("accept", &asked), version: "2025-11-25"},
	}
	for name, opts := range clients {
		var out actionOutput
		decode(t, call(t, f.connect(t, f.token(t), opts), toolSend, compose(f.accountID)), &out)
		if _, err := f.mail.Draft(t.Context(), f.userID, f.accountID, out.DraftID); err != nil || out.Status != statusPending ||
			!strings.HasPrefix(out.ConfirmURL, f.srv.URL+"/pending?action=") {
			t.Fatalf("%s: want a kept draft and a confirm link, got %+v (%v)", name, out, err)
		}
	}
	if f.mail.sentCount() != 0 || len(asked) != 0 {
		t.Fatalf("nothing may be sent or asked, got %d sent, %d asked", f.mail.sentCount(), len(asked))
	}
}
