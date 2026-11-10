package mcp

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/FacileStudio/porte/mcpauth"
)

// listTools is a tools/call-free JSON-RPC request every token may send.
const listTools = `{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}`

// post sends one raw JSON-RPC body to the endpoint.
func (f *fixture) post(t *testing.T, token, body string, header map[string]string) *http.Response {
	t.Helper()
	req, _ := http.NewRequest(http.MethodPost, f.resource, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	for name, value := range header {
		req.Header.Set(name, value)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = resp.Body.Close() })
	return resp
}

// TestNoTokenIsAChallengeToTheMetadata checks the 401 points at the RFC 9728
// document, and that the document names this resource and Registre.
func TestNoTokenIsAChallengeToTheMetadata(t *testing.T) {
	f := newFixture(t)
	resp := f.post(t, "", listTools, nil)
	challenge := resp.Header.Get("WWW-Authenticate")
	metadataURL, _ := mcpauth.MetadataURL(f.resource)
	if resp.StatusCode != http.StatusUnauthorized || !strings.Contains(challenge, `resource_metadata="`+metadataURL+`"`) {
		t.Fatalf("want 401 naming %s, got %d %q", metadataURL, resp.StatusCode, challenge)
	}
	doc, err := http.Get(metadataURL)
	if err != nil {
		t.Fatal(err)
	}
	defer doc.Body.Close()
	var metadata struct {
		Resource             string   `json:"resource"`
		AuthorizationServers []string `json:"authorization_servers"`
		ScopesSupported      []string `json:"scopes_supported"`
	}
	_ = json.NewDecoder(doc.Body).Decode(&metadata)
	if metadata.Resource != f.resource || len(metadata.AuthorizationServers) != 1 || metadata.AuthorizationServers[0] != f.reg.srv.URL {
		t.Fatalf("metadata names the wrong resource or server: %+v", metadata)
	}
}

// TestTokenForAnotherResourceIsRefused checks audience binding.
func TestTokenForAnotherResourceIsRefused(t *testing.T) {
	f := newFixture(t)
	token := f.reg.mint(t, "https://agenda.example/mcp", nil)
	resp := f.post(t, token, listTools, nil)
	if resp.StatusCode != http.StatusUnauthorized || !strings.Contains(resp.Header.Get("WWW-Authenticate"), `error="invalid_token"`) {
		t.Fatalf("want 401 invalid_token, got %d %q", resp.StatusCode, resp.Header.Get("WWW-Authenticate"))
	}
}

// TestReadOnlyTokenCannotSend checks the 403 step-up challenge on a write
// tool, and that reads still work with the same token.
func TestReadOnlyTokenCannotSend(t *testing.T) {
	f := newFixture(t)
	token := f.reg.mint(t, f.resource, map[string]any{"scope": "courrier.read"})
	send := `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"courrier_send","arguments":{"account_id":1,"draft_id":2}}}`
	resp := f.post(t, token, send, nil)
	challenge := resp.Header.Get("WWW-Authenticate")
	if resp.StatusCode != http.StatusForbidden || !strings.Contains(challenge, `error="insufficient_scope"`) ||
		!strings.Contains(challenge, `scope="courrier.write"`) {
		t.Fatalf("want 403 insufficient_scope for courrier.write, got %d %q", resp.StatusCode, challenge)
	}
	if result := call(t, f.connect(t, token, clientOptions{}), toolListAccounts, map[string]any{}); result.IsError {
		t.Fatalf("a read-only token must still read: %+v", result.Content)
	}
}

// TestBulleExchangedTokenIsAccepted checks a delegated token Bulle got by
// token exchange with resource: same audience, plus act.
func TestBulleExchangedTokenIsAccepted(t *testing.T) {
	f := newFixture(t)
	token := f.reg.mint(t, f.resource, map[string]any{
		"client_id": mcpauth.BulleClientID, "act": map[string]string{"sub": mcpauth.BulleClientID},
	})
	result := call(t, f.connect(t, token, clientOptions{}), toolListAccounts, map[string]any{})
	if result.IsError {
		t.Fatalf("Bulle's exchanged token must work: %+v", result.Content)
	}
}

// TestForeignOriginIsRefused checks the DNS rebinding guard, and that the
// endpoint's own origin passes.
func TestForeignOriginIsRefused(t *testing.T) {
	f := newFixture(t)
	if resp := f.post(t, f.token(t), listTools, map[string]string{"Origin": "https://evil.example"}); resp.StatusCode != http.StatusForbidden {
		t.Fatalf("want 403 for a foreign origin, got %d", resp.StatusCode)
	}
	if resp := f.post(t, "", listTools, map[string]string{"Origin": f.srv.URL}); resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("the own origin must reach the token check, got %d", resp.StatusCode)
	}
}

// TestGetIsNotAllowed checks the stateless transport refuses GET.
func TestGetIsNotAllowed(t *testing.T) {
	f := newFixture(t)
	req, _ := http.NewRequest(http.MethodGet, f.resource, nil)
	req.Header.Set("Authorization", "Bearer "+f.token(t))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("want 405 on GET, got %d", resp.StatusCode)
	}
}
