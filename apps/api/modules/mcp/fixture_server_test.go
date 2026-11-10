package mcp

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FacileStudio/porte/mcpauth"
	"github.com/go-chi/chi/v5"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// fixture is a running MCP endpoint over a fake mailbox.
type fixture struct {
	srv       *httptest.Server
	reg       *registre
	mail      *fakeMail
	service   *Service
	userID    int64
	accountID int64
	resource  string
}

// newFixture starts Courrier's MCP routes behind a fake Registre.
func newFixture(t *testing.T) *fixture {
	t.Helper()
	return newFixtureOn(t, nil, nil)
}

// newFixtureOn is newFixture on a given listener, nil for any free port,
// with front wrapping the routes when it is set.
func newFixtureOn(t *testing.T, listener net.Listener, front func(*fixture, http.Handler) http.Handler) *fixture {
	t.Helper()
	f := &fixture{reg: newRegistre(t)}
	f.userID, f.accountID = seedPerson(t)
	f.mail = newFakeMail(f.accountID)
	var handler http.Handler
	f.srv = httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { handler.ServeHTTP(w, r) }))
	if listener != nil {
		f.srv.Listener = listener
	}
	f.srv.Start()
	t.Cleanup(f.srv.Close)
	f.resource = f.srv.URL + "/mcp"
	verifier, err := mcpauth.Verifier(context.Background(), mcpauth.Config{
		Issuer: f.reg.srv.URL, Resource: f.resource, ScopePrefix: scopePrefix, Identities: identities{userID: f.userID},
	})
	if err != nil {
		t.Fatal(err)
	}
	f.service, err = NewService(testORM, f.mail, Config{
		Resource: f.resource, Issuer: f.reg.srv.URL, Version: "test", StateKey: []byte("state-key"), Verifier: verifier,
	})
	if err != nil {
		t.Fatal(err)
	}
	router := chi.NewRouter()
	RegisterRoutes(router, f.service)
	handler = router
	if front != nil {
		handler = front(f, router)
	}
	return f
}

// token is a third-party token for this endpoint with both scopes.
func (f *fixture) token(t *testing.T) string {
	return f.reg.mint(t, f.resource, nil)
}

// bearer adds an access token to every request.
type bearer string

// RoundTrip sends r with the token.
func (b bearer) RoundTrip(r *http.Request) (*http.Response, error) {
	r = r.Clone(r.Context())
	r.Header.Set("Authorization", "Bearer "+string(b))
	return http.DefaultTransport.RoundTrip(r)
}

// clientOptions shape the test client: whether it can elicit, and which
// protocol version it speaks.
type clientOptions struct {
	elicit   func(context.Context, *sdk.ElicitRequest) (*sdk.ElicitResult, error)
	version  string
	noRetrys bool
}

// connect opens an SDK client session with token.
func (f *fixture) connect(t *testing.T, token string, opts clientOptions) *sdk.ClientSession {
	t.Helper()
	client := sdk.NewClient(&sdk.Implementation{Name: "test", Version: "1"}, &sdk.ClientOptions{
		ElicitationHandler: opts.elicit,
		MultiRoundTrip:     &sdk.MultiRoundTripOptions{Disabled: opts.noRetrys},
	})
	transport := &sdk.StreamableClientTransport{Endpoint: f.resource, HTTPClient: &http.Client{Transport: bearer(token)}}
	session, err := client.Connect(context.Background(), transport, &sdk.ClientSessionOptions{ProtocolVersion: opts.version})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = session.Close() })
	return session
}

// call runs one tool and fails the test on a protocol error.
func call(t *testing.T, session *sdk.ClientSession, name string, args any) *sdk.CallToolResult {
	t.Helper()
	result, err := session.CallTool(context.Background(), &sdk.CallToolParams{Name: name, Arguments: args})
	if err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	return result
}
