package spaces

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FacileStudio/Courrier/apps/api/modules/auth"
	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/porte/directory"
	"github.com/FacileStudio/porte/oidc"
	jwtpkg "github.com/FacileStudio/porte/oidc/jwt"
	portepg "github.com/FacileStudio/porte/pg"
	"github.com/FacileStudio/porte/session"
)

var _ porte.ErasedSubjects = (*Deletions)(nil)

// recreateGone gives sub-gone a fresh account, as a bearer token minted before
// the erasure or a stale snapshot would, and returns its user id.
func recreateGone(t *testing.T) int64 {
	t.Helper()
	again := user(t, "gone-again@facile.studio")
	identity(t, again, "sub-gone")
	return again
}

// A re-polled notice erases an account that came back for an erased subject,
// without a restart, and leaves everybody else alone.
func TestARecreatedAccountIsErasedAgainOnTheNextPoll(t *testing.T) {
	f := seedErasure(t)
	pollErased(t, f.deletions)
	again := recreateGone(t)
	others := count(t, testORM, `SELECT count(*) FROM users WHERE id <> ?`, again)

	if err := f.deletions.receiver.Poll(context.Background()); err != nil {
		t.Fatalf("re-poll: %v", err)
	}

	if n := count(t, testORM, `SELECT count(*) FROM users WHERE id = ?`, again); n != 0 {
		t.Errorf("the re-created account survived the re-poll")
	}
	if n := count(t, testORM, `SELECT count(*) FROM porte_identities WHERE subject = 'sub-gone'`); n != 0 {
		t.Errorf("sub-gone still holds %d identities", n)
	}
	if n := count(t, testORM, `SELECT count(*) FROM users`); n != others {
		t.Errorf("users = %d after the re-poll, want the %d others untouched", n, others)
	}
}

// Erased answers from the saved state, for either issuer spelling, and only
// for the subject Courrier erased.
func TestErasedAnswersFromTheSavedState(t *testing.T) {
	f := seedErasure(t)
	pollErased(t, f.deletions)
	for _, c := range []struct {
		provider, subject string
		want              bool
	}{
		{testIssuer, "sub-gone", true},
		{testIssuer + "/", "sub-gone", true},
		{testIssuer, "sub-other", false},
		{"https://elsewhere.test", "sub-gone", false},
	} {
		got, err := f.deletions.Erased(context.Background(), c.provider, c.subject)
		if err != nil || got != c.want {
			t.Errorf("Erased(%s, %s) = %v, %v, want %v", c.provider, c.subject, got, err, c.want)
		}
	}
}

// A snapshot that still lists an erased member does not pre-provision them.
func TestTheSpaceSyncSkipsAnErasedMember(t *testing.T) {
	f := seedErasure(t)
	pollErased(t, f.deletions)
	dir := NewDirectory(testORM, directory.Client{Issuer: testIssuer}, slog.New(slog.NewTextHandler(io.Discard, nil))).
		WithErased(f.deletions)
	stale := directory.Snapshot{Spaces: []directory.Space{{
		ID: "space-erased", Name: "Erased", Members: []directory.Member{
			{Subject: "sub-gone", Email: "gone@facile.studio", Role: directory.RoleMember},
		},
	}}}

	if err := dir.Apply(context.Background(), stale); err != nil {
		t.Fatalf("apply: %v", err)
	}

	if n := count(t, testORM, `SELECT count(*) FROM users WHERE lower(email) = 'gone@facile.studio'`); n != 0 {
		t.Errorf("the space sync re-created the erased member")
	}
}

// newIssuer is a Registre that only answers discovery, enough for oidc.New.
func newIssuer(t *testing.T) string {
	t.Helper()
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"issuer":                 server.URL,
			"authorization_endpoint": server.URL + "/authorize",
			"token_endpoint":         server.URL + "/token",
			"jwks_uri":               server.URL + "/keys",
		})
	}))
	t.Cleanup(server.Close)
	return server.URL
}

// The OIDC kit built as main.go builds it refuses to provision an account for
// an erased subject from a bearer token, and creates nothing.
func TestBearerProvisioningRefusesAnErasedSubject(t *testing.T) {
	f := seedErasure(t)
	pollErased(t, f.deletions)
	issuer := newIssuer(t)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	sqlDB, err := testORM.DB()
	if err != nil {
		t.Fatal(err)
	}
	store := portepg.New(sqlDB)
	users, err := auth.NewUserStore(context.Background(), sqlDB, logger)
	if err != nil {
		t.Fatal(err)
	}
	cfg := porte.Config{Issuer: issuer, ClientID: "courrier", ClientSecret: "secret", RedirectURL: "https://courrier.test/api/auth/oidc/callback", SuccessURL: "https://courrier.test"}
	sessions, err := session.New(cfg, session.Deps{Sessions: store.Sessions(), Logger: logger})
	if err != nil {
		t.Fatal(err)
	}
	erased := NewDeletions(testORM, DeletionsConfig{Client: directory.Client{Issuer: issuer}, Logger: logger})
	kit, err := oidc.New(context.Background(), cfg, oidc.Deps{
		Users: users, Identities: store.Identities(), Sessions: sessions, Codes: store.LoginCodes(), Logger: logger, Erased: erased,
	})
	if err != nil {
		t.Fatalf("oidc: %v", err)
	}
	before := count(t, testORM, `SELECT count(*) FROM users`)
	verified := true

	_, err = kit.ProvisionBearer(context.Background(), jwtpkg.Claims{
		Subject: "sub-gone", Email: "gone@facile.studio", Profile: jwtpkg.Profile{EmailVerified: &verified},
	})

	if !stderrors.Is(err, porte.ErrSubjectErased) {
		t.Fatalf("ProvisionBearer = %v, want a refusal wrapping porte.ErrSubjectErased", err)
	}
	if n := count(t, testORM, `SELECT count(*) FROM users`); n != before {
		t.Errorf("users = %d after a refused provisioning, want %d", n, before)
	}
}
