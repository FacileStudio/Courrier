package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/porte/kit"
)

type stubAuth struct{ identity porte.Identity }

func (s stubAuth) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(porte.WithIdentity(r.Context(), s.identity)))
	})
}

func (stubAuth) Identity(_ context.Context, userID int64) (kit.Identity, error) {
	return kit.Identity{UserID: userID, Email: "a@example.com"}, nil
}

func status(identity porte.Identity, method, path string) int {
	handler := RequireAuth(stubAuth{identity: identity})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(method, path, nil))
	return recorder.Code
}

func TestDelegatedTokenCannotManageCredentials(t *testing.T) {
	delegated := porte.Identity{UserID: 1, Actor: "bulle"}
	session := porte.Identity{UserID: 1, SessionID: 7}
	cases := []struct {
		identity     porte.Identity
		method, path string
		want         int
	}{
		{delegated, http.MethodPost, "/api/users/me/api-token", http.StatusForbidden},
		{delegated, http.MethodDelete, "/api/users/me/api-token", http.StatusForbidden},
		{delegated, http.MethodPatch, "/api/users/me", http.StatusForbidden},
		{delegated, http.MethodGet, "/api/users/me", http.StatusNoContent},
		{delegated, http.MethodPost, "/api/accounts/1/mail/send", http.StatusNoContent},
		{session, http.MethodPost, "/api/users/me/api-token", http.StatusNoContent},
		{session, http.MethodPatch, "/api/users/me", http.StatusNoContent},
	}
	for _, c := range cases {
		if got := status(c.identity, c.method, c.path); got != c.want {
			t.Errorf("%s %s actor=%q: got %d, want %d", c.method, c.path, c.identity.Actor, got, c.want)
		}
	}
}
