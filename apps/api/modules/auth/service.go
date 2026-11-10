package auth

import (
	"context"
	"net/http"

	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/porte/kit"
	"github.com/FacileStudio/porte/session"
)

// Service is what is left of Courrier's authentication after porte took the
// credential: the profile lookup the rest of the app reads, and the named API
// tokens.
type Service struct {
	sessions *session.Manager
	users    kit.Resolver
}

// NewService builds the auth Service over the porte session manager and the
// kit store that resolves a user id to its identity.
func NewService(sessions *session.Manager, users kit.Resolver) *Service {
	return &Service{sessions: sessions, users: users}
}

// RequireAuth is porte's session middleware, re-exported so the module routers
// keep passing this one service to middleware.RequireAuth.
func (service *Service) RequireAuth(next http.Handler) http.Handler {
	return service.sessions.RequireAuth(next)
}

// Identity turns the user id porte authenticated into the identity the rest
// of Courrier reads, through porte/kit's store. A session that outlived its
// user answers 401 invalid auth token.
func (service *Service) Identity(ctx context.Context, userID int64) (kit.Identity, error) {
	return service.users.Identity(ctx, userID)
}

// Issue mints a named API token: a porte session with a label and no expiry,
// which is what the separate api_tokens table used to be.
func (service *Service) Issue(ctx context.Context, userID int64, label string) (string, porte.Session, error) {
	return service.sessions.Issue(ctx, userID, label)
}

// AuthenticateRequest resolves the caller of a route that is not mounted
// behind RequireAuth — the inline-image endpoint, which a browser reaches with
// an <img src> and therefore with a cookie and no header.
func (service *Service) AuthenticateRequest(w http.ResponseWriter, r *http.Request) (int64, error) {
	identity, err := service.sessions.Authenticate(w, r)
	if err != nil {
		return 0, err
	}
	return identity.UserID, nil
}

// Sessions exposes the manager for the modules that list or revoke tokens.
func (service *Service) Sessions() *session.Manager { return service.sessions }
