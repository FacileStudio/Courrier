package main

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FacileStudio/porte"
	portepg "github.com/FacileStudio/porte/pg"
	"github.com/FacileStudio/porte/session"
	"github.com/go-chi/chi/v5"

	"github.com/FacileStudio/Courrier/apps/api/internal/env"
	"github.com/FacileStudio/Courrier/apps/api/internal/middleware"
	"github.com/FacileStudio/Courrier/apps/api/modules/accounts"
	"github.com/FacileStudio/Courrier/apps/api/modules/auth"
	"github.com/FacileStudio/Courrier/apps/api/modules/mail"
)

const (
	readToken     = "aaa.bbb.read"
	writeToken    = "aaa.bbb.write"
	noScopeToken  = "aaa.bbb.none"
	bearerSubject = "user-7"
)

type scopeVerifier struct{}

func (scopeVerifier) VerifyJWT(_ context.Context, raw string) (porte.Identity, error) {
	identity := porte.Identity{UserID: 7, Subject: bearerSubject, Actor: "bulle"}
	switch raw {
	case readToken:
		identity.Scopes = []string{"courrier.read"}
	case writeToken:
		identity.Scopes = []string{"courrier.write"}
	case noScopeToken:
	default:
		return porte.Identity{}, errors.New("unknown token")
	}
	return identity, nil
}

type stubAuthenticator struct{ sessions *session.Manager }

func (s stubAuthenticator) RequireAuth(next http.Handler) http.Handler {
	return s.sessions.RequireAuth(next)
}

func (stubAuthenticator) IdentityForUser(context.Context, int64) (string, string, error) {
	return "7", "user@example.test", nil
}

func bearerSessions(t *testing.T) *session.Manager {
	t.Helper()
	cfg := env.Config{OIDC: &env.OIDCConfig{
		Issuer:          "https://sso.example.test",
		ClientID:        "courrier-client",
		ClientSecret:    "secret",
		RedirectURL:     "https://courrier.example.test/api/auth/oidc/callback",
		MachineAudience: "courrier-client",
		ScopePrefix:     "courrier",
	}}.Porte()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	sessions, err := session.New(cfg, session.Deps{Sessions: portepg.New(nil).Sessions(), Logger: logger})
	if err != nil {
		t.Fatal(err)
	}
	return sessions.WithJWT(scopeVerifier{})
}

func bearerRequest(method, path, token string) *http.Request {
	request := httptest.NewRequest(method, path, nil)
	request.Header.Set("Authorization", "Bearer "+token)
	return request
}

func scopeRouters(t *testing.T) (real, readable http.Handler) {
	t.Helper()
	sessions := bearerSessions(t)
	authService := auth.NewService(nil, sessions, slog.New(slog.NewTextHandler(io.Discard, nil)))

	realRoutes := chi.NewRouter()
	realRoutes.Route("/api", func(r chi.Router) {
		accounts.RegisterRoutes(r, nil, authService)
		mail.RegisterRoutes(r, nil, authService, nil)
	})

	readRoutes := chi.NewRouter()
	readRoutes.With(middleware.RequireAuth(stubAuthenticator{sessions})).
		Get("/api/accounts/{accountId}/mail/search", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	return realRoutes, readRoutes
}

func TestReadScopeReadsAndSendIsForbidden(t *testing.T) {
	realRoutes, readable := scopeRouters(t)

	cases := []struct {
		name   string
		router http.Handler
		method string
		path   string
		token  string
		want   int
	}{
		{"read scope reads", readable, http.MethodGet, "/api/accounts/1/mail/search?q=x", readToken, http.StatusOK},
		{"write scope reads", readable, http.MethodGet, "/api/accounts/1/mail/search?q=x", writeToken, http.StatusOK},
		{"read scope cannot send", realRoutes, http.MethodPost, "/api/accounts/1/mail/send", readToken, http.StatusForbidden},
		{"read scope cannot move", realRoutes, http.MethodPatch, "/api/accounts/1/mail/emails/2", readToken, http.StatusForbidden},
		{"no scope cannot read", readable, http.MethodGet, "/api/accounts/1/mail/search?q=x", noScopeToken, http.StatusForbidden},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			tc.router.ServeHTTP(recorder, bearerRequest(tc.method, tc.path, tc.token))
			if recorder.Code != tc.want {
				t.Fatalf("%s %s = %d, want %d: %s", tc.method, tc.path, recorder.Code, tc.want, recorder.Body.String())
			}
		})
	}
}
