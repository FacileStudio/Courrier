package authcontext

import "context"

// Identity is the authenticated caller as the rest of the app reads it: the
// user id porte resolved and the email looked up from the users row.
type Identity struct {
	UserID string
	Email  string
}

type contextKey struct{}

// WithIdentity stores the caller's identity in the request context so
// handlers mounted behind RequireAuth can retrieve it.
func WithIdentity(parentContext context.Context, identity Identity) context.Context {
	return context.WithValue(parentContext, contextKey{}, identity)
}

// IdentityFromContext returns the identity stored in the context, and whether
// one was present.
func IdentityFromContext(parentContext context.Context) (Identity, bool) {
	identity, ok := parentContext.Value(contextKey{}).(Identity)
	return identity, ok
}

// MustIdentity returns the context identity or panics. Call it only on routes
// that are guaranteed to run behind RequireAuth.
func MustIdentity(ctx context.Context) Identity {
	identity, ok := IdentityFromContext(ctx)
	if !ok {
		panic("authcontext: missing identity in context")
	}
	return identity
}
