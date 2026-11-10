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

// answering returns an elicitation handler that answers action, and
// records the questions.
func answering(action string, asked *[]string) func(context.Context, *sdk.ElicitRequest) (*sdk.ElicitResult, error) {
	return func(_ context.Context, req *sdk.ElicitRequest) (*sdk.ElicitResult, error) {
		*asked = append(*asked, req.Params.Message)
		return &sdk.ElicitResult{Action: action, Content: map[string]any{"confirm": action == "accept"}}, nil
	}
}

// TestConfirmedSendGoesOutOnce checks the MRTR round trip: the client asks,
// the person accepts, the mail goes out once and the draft is gone.
func TestConfirmedSendGoesOutOnce(t *testing.T) {
	f := newFixture(t)
	var asked []string
	session := f.connect(t, f.token(t), clientOptions{Elicit: answering("accept", &asked)})
	var out actionOutput
	decode(t, call(t, session, toolSend, compose(f.accountID)), &out)
	if out.Status != "sent" || f.mail.sentCount() != 1 || len(asked) != 1 || !strings.Contains(asked[0], `"Quote"`) {
		t.Fatalf("want one confirmed send, got status %q, %d sent, questions %q", out.Status, f.mail.sentCount(), asked)
	}
}
