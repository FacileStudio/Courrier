package middleware

import (
	"net/http"
	"strings"

	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/porte/kit"
	"github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"
)

// Authenticator is the auth service: porte's session middleware, plus the
// kit resolver that turns the user id porte authenticated into the identity
// the rest of Courrier reads.
//
// It stays one parameter so every module router keeps calling
// middleware.RequireAuth(authService) unchanged.
type Authenticator interface {
	kit.Authenticator
	kit.Resolver
}

// managesCredentials reports whether the route mints, lists or revokes API
// tokens, or edits the profile that holds the sign-in password and address. A
// delegated suite token (porte sets Actor from its act.sub claim) is refused on
// these: a five-minute delegation must not turn itself into a credential that
// outlives it.
func managesCredentials(method, path string) bool {
	path = strings.TrimSuffix(strings.TrimPrefix(path, "/api"), "/")
	if path == "/users/me" {
		return method != http.MethodGet
	}
	return strings.HasPrefix(path, "/users/me/api-token") || strings.HasPrefix(path, "/users/me/password")
}

// delegationGuard runs porte's middleware and then refuses a delegated token
// on the credential routes, before the kit reads the users row.
type delegationGuard struct {
	Authenticator
}

// RequireAuth wraps next in porte's middleware and the delegation check.
func (guard delegationGuard) RequireAuth(next http.Handler) http.Handler {
	return guard.Authenticator.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		authenticated, ok := porte.From(request.Context())
		if ok && authenticated.Actor != "" && managesCredentials(request.Method, request.URL.Path) {
			httpjson.WriteError(w, errors.Forbidden("a delegated token cannot manage credentials"))
			return
		}
		next.ServeHTTP(w, request)
	}))
}

// RequireAuth is porte/kit's middleware with Courrier's delegation rule:
// porte verifies the cookie or bearer credential, a delegated token is
// refused on the credential routes, and the kit hydrates kit.Identity from
// the users table. A session whose user is gone answers 401.
func RequireAuth(auth Authenticator) func(http.Handler) http.Handler {
	return kit.RequireAuth(delegationGuard{auth}, auth)
}

// Identity returns the caller RequireAuth stored on the request. It panics
// on a route mounted outside RequireAuth, which is a wiring bug.
func Identity(request *http.Request) kit.Identity {
	identity, ok := kit.FromContext(request.Context())
	if !ok {
		panic("middleware: missing identity in context")
	}
	return identity
}

// UserID is the caller's user id. It panics like Identity outside RequireAuth.
func UserID(request *http.Request) int64 {
	return Identity(request).UserID
}
