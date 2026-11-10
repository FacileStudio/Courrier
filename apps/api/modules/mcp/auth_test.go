package mcp

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/FacileStudio/Courrier/apps/api/modules/mail"
	"github.com/FacileStudio/porte/mcpauth"
	"github.com/FacileStudio/tronc/mcp/mcptest"
)

// listTools is a tools/call-free JSON-RPC request every token may send.
const listTools = `{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}`

// TestConformance runs tronc's conformance kit on the endpoint, with send
// and delete draft as the sensitive tools.
func TestConformance(t *testing.T) {
	f := newFixture(t)
	mcptest.Run(t, mcptest.Config{
		URL:   f.resource,
		Token: f.token,
		Confirm: []mcptest.ConfirmCase{
			{Tool: toolSend, Setup: f.freshDraft},
			{Tool: toolDeleteDraft, Setup: f.freshDraft},
		},
	})
}

// freshDraft saves a new draft and counts how often it is sent or deleted.
func (f *fixture) freshDraft(t testing.TB) (map[string]any, func() int) {
	draft, err := f.mail.SaveDraft(t.Context(), f.userID, f.accountID, mail.SendRequest{
		To: []string{"client@example.test"}, Subject: fmt.Sprintf("Quote %s", t.Name()), Body: "Here it is.",
	})
	if err != nil {
		t.Fatal(err)
	}
	return map[string]any{"account_id": f.accountID, "draft_id": draft.ID}, func() int { return f.mail.removals(draft.ID) }
}

// post sends one raw JSON-RPC body to the endpoint.
func (f *fixture) post(t *testing.T, token, body string) *http.Response {
	t.Helper()
	req, _ := http.NewRequest(http.MethodPost, f.resource, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = resp.Body.Close() })
	return resp
}

// TestTokenForAnotherResourceIsRefused checks audience binding.
func TestTokenForAnotherResourceIsRefused(t *testing.T) {
	f := newFixture(t)
	token := f.issuer.Token(t, "https://agenda.example/mcp", subject, "courrier.read", "courrier.write")
	resp := f.post(t, token, listTools)
	if resp.StatusCode != http.StatusUnauthorized || !strings.Contains(resp.Header.Get("WWW-Authenticate"), `error="invalid_token"`) {
		t.Fatalf("want 401 invalid_token, got %d %q", resp.StatusCode, resp.Header.Get("WWW-Authenticate"))
	}
}

// TestBulleExchangedTokenIsAccepted checks a delegated token Bulle got by
// token exchange with resource: same audience, plus act.
func TestBulleExchangedTokenIsAccepted(t *testing.T) {
	f := newFixture(t)
	token := f.issuer.Mint(t, map[string]any{
		"aud": []string{f.resource}, "sub": subject, "scope": "courrier.read courrier.write",
		"client_id": mcpauth.BulleClientID, "act": map[string]string{"sub": mcpauth.BulleClientID},
	})
	result := call(t, f.connect(t, token, clientOptions{}), toolListAccounts, map[string]any{})
	if result.IsError {
		t.Fatalf("Bulle's exchanged token must work: %+v", result.Content)
	}
}
