package middleware

import (
	"context"
	"net/http"

	"github.com/FacileStudio/Courrier/apps/api/internal/authcontext"
	"github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"
)

type Authenticator interface {
	Authenticate(context context.Context, authorization string) (string, any, error)
}

// extractToken reads the session credential from the cookie first, then the
// Authorization header.
//
// It deliberately does not read a ?token= query parameter. That path existed
// with a deprecation warning and nothing used it — the client sends
// credentials: 'include' on every request — but a credential in a URL is
// copied into access logs, Referer headers, browser history and any proxy in
// between, and the warning made it no less true.
func extractToken(r *http.Request) string {
	if c, err := r.Cookie("session"); err == nil && c.Value != "" {
		return c.Value
	}

	return r.Header.Get("Authorization")
}

func RequireAuth(authService Authenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			token := extractToken(request)
			if token == "" {
				httpjson.WriteError(w, errors.Unauthorized("missing auth"))
				return
			}

			userID, rawData, err := authService.Authenticate(request.Context(), token)
			if err != nil {
				httpjson.WriteError(w, err)
				return
			}
			data, ok := rawData.(interface{ GetEmail() string })
			if !ok || data == nil {
				httpjson.WriteError(w, errors.Unauthorized("missing auth"))
				return
			}

			authContext := authcontext.WithIdentity(request.Context(), authcontext.Identity{
				UserID: userID,
				Email:  data.GetEmail(),
			})
			next.ServeHTTP(w, request.WithContext(authContext))
		})
	}
}
