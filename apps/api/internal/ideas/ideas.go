// Package ideas serves Bulle's Ideas API to Courrier's own client.
//
// The routes proxy Bulle with a token exchanged for the signed-in user. The
// package decides the two things Bulle cannot know: which Registre space is
// active, taken from a space_id the user is a member of and never sent to
// Bulle as the client wrote it, and whether that space or person has ai.access.
package ideas

import (
	"context"
	"net/http"

	"github.com/FacileStudio/Courrier/apps/api/internal/plans"
	"github.com/FacileStudio/Courrier/apps/api/internal/spaceid"
	"github.com/FacileStudio/porte/kit"

	"github.com/FacileStudio/tronc/bulle"
	"github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// App is Courrier's name in Bulle's idea enum.
const App = "courrier"

type scopeKey struct{}

type scope struct {
	space    string
	aiAccess bool
}

// Gate resolves the space and the plan of every Ideas request.
type Gate struct {
	db    *gorm.DB
	plans *plans.Gate
}

// New builds the gate over the database and the plan gate.
func New(db *gorm.DB, gate *plans.Gate) *Gate {
	return &Gate{db: db, plans: gate}
}

// Mount serves /ideas on router, which must already require authentication.
// opts supplies the Bulle URL, the token source and the upgrade link; Mount
// sets the app, the space and the AI check.
func (g *Gate) Mount(router chi.Router, opts bulle.Options) error {
	opts.App = App
	opts.Space = func(r *http.Request) string { return scopeOf(r).space }
	opts.AIAccess = func(r *http.Request) bool { return scopeOf(r).aiAccess }
	router.Use(g.resolve)
	return bulle.MountIdeas(router, opts)
}

func scopeOf(r *http.Request) scope {
	resolved, _ := r.Context().Value(scopeKey{}).(scope)
	return resolved
}

func (g *Gate) resolve(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resolved, err := g.scope(r)
		if err != nil {
			httpjson.WriteError(w, err)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), scopeKey{}, resolved)))
	})
}

func (g *Gate) scope(r *http.Request) (scope, error) {
	identity, ok := kit.FromContext(r.Context())
	if !ok {
		return scope{}, errors.Unauthorized("missing auth")
	}
	userID := identity.UserID
	id, err := spaceid.Query(r.URL.Query().Get("space_id"))
	if err != nil {
		return scope{}, err
	}
	if id != nil {
		return g.spaceScope(r.Context(), *id, userID)
	}
	access, err := g.plans.PersonalAIAccess(g.db.WithContext(r.Context()), userID)
	if err != nil {
		return scope{}, err
	}
	return scope{aiAccess: access}, nil
}

func (g *Gate) spaceScope(ctx context.Context, id string, userID int64) (scope, error) {
	var links []string
	err := g.db.WithContext(ctx).Raw(`SELECT coalesce(s.registre_space_id, '') FROM spaces s
		JOIN space_members m ON m.space_id = s.id
		WHERE s.id = ? AND m.user_id = ?`, id, userID).Scan(&links).Error
	if err != nil {
		return scope{}, errors.Internal("failed to load the space", err)
	}
	if len(links) == 0 {
		return scope{}, errors.NotFound("space not found")
	}
	if links[0] == "" {
		return scope{}, errors.Failed("this space is not linked to Registre yet, so Bulle cannot suggest ideas for it")
	}
	return scope{space: links[0], aiAccess: g.plans.SpaceAIAccess(links[0])}, nil
}
