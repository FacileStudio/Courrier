package mcp

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FacileStudio/porte/mcpauth"
	mcpkit "github.com/FacileStudio/tronc/mcp"
	"github.com/FacileStudio/tronc/mcp/mcptest"
	"github.com/go-chi/chi/v5"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// fixture is a running MCP endpoint over a fake mailbox.
type fixture struct {
	srv       *httptest.Server
	issuer    *mcptest.Issuer
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

// pendingStore is mcp_actions on the scratch database, or memory without
// one.
func pendingStore(t *testing.T) mcpkit.Pending {
	if testORM == nil {
		return mcptest.NewPending()
	}
	sqlDB, err := testORM.DB()
	if err != nil {
		t.Fatal(err)
	}
	return mcpkit.PostgresPending{DB: sqlDB}
}

// newFixtureOn is newFixture on a given listener, nil for any free port,
// with front wrapping the routes when it is set.
func newFixtureOn(t *testing.T, listener net.Listener, front func(*fixture, http.Handler) http.Handler) *fixture {
	t.Helper()
	f := &fixture{issuer: mcptest.NewIssuer(t)}
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
		Issuer: f.issuer.URL, Resource: f.resource, ScopePrefix: "courrier", Identities: identities{userID: f.userID},
	})
	if err != nil {
		t.Fatal(err)
	}
	f.service, err = NewService(f.mail, Config{
		Resource: f.resource, Issuer: f.issuer.URL, Version: "test", StateKey: []byte("state-key"), Verifier: verifier,
		Pending: pendingStore(t),
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

// token is a third-party token for this endpoint with scopes, or both when
// none is given.
func (f *fixture) token(t testing.TB, scopes ...string) string {
	if len(scopes) == 0 {
		scopes = []string{"courrier.read", "courrier.write"}
	}
	return f.issuer.Token(t, f.resource, subject, scopes...)
}

// clientOptions shape the test client: whether it can elicit, and which
// protocol version it speaks.
type clientOptions = mcptest.ClientOptions

// connect opens an SDK client session with token.
func (f *fixture) connect(t *testing.T, token string, opts clientOptions) *sdk.ClientSession {
	t.Helper()
	return mcptest.Connect(t, f.resource, token, opts)
}

// call runs one tool and fails the test on a protocol error.
func call(t *testing.T, session *sdk.ClientSession, name string, args any) *sdk.CallToolResult {
	t.Helper()
	return mcptest.Call(t, session, name, args)
}
