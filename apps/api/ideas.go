package main

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/internal/env"
	"github.com/FacileStudio/Courrier/apps/api/internal/ideas"
	"github.com/FacileStudio/Courrier/apps/api/internal/middleware"
	"github.com/FacileStudio/Courrier/apps/api/modules/auth"

	"github.com/FacileStudio/porte/exchange"
	"github.com/FacileStudio/tronc/bulle"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// completeTimeout bounds one Bulle completion: the writing helper and the
// inbox triage wait for a model, which the 20 s default does not cover.
const completeTimeout = 90 * time.Second

// bulleLink is how Courrier reaches Bulle as the signed-in user, over one
// Registre token exchange for audience bulle: the Ideas routes, the client
// the writing helper and the inbox triage complete through, and the token
// source that says whether the user is linked.
type bulleLink struct {
	ideas  *bulle.Options
	client *bulle.Client
	linked bulle.TokenSource
}

// newBulle builds the Bulle link over the identity store that seals the
// Registre tokens. It is empty when SSO or BULLE_URL is off.
func newBulle(appEnv env.Config, identities exchange.Identities, appLogger *slog.Logger) (bulleLink, error) {
	if appEnv.OIDC == nil || appEnv.BulleURL == "" {
		return bulleLink{}, nil
	}
	source, err := bulleSource(appEnv, identities, appLogger)
	if err != nil {
		return bulleLink{}, err
	}
	client, err := bulle.New(bulle.Config{
		URL:         appEnv.BulleURL,
		TokenSource: source.Token,
		Forget:      source.Forget,
		HTTPClient:  &http.Client{Timeout: completeTimeout},
	})
	if err != nil {
		return bulleLink{}, err
	}
	return bulleLink{
		ideas: &bulle.Options{
			URL:         appEnv.BulleURL,
			TokenSource: source.Token,
			Forget:      source.Forget,
			UpgradeURL:  upgradeURL(appEnv),
		},
		client: client,
		linked: source.Token,
	}, nil
}

// bulleSource exchanges the user's Registre session for a Bulle token
// (audience bulle, scopes bulle.read and bulle.write).
func bulleSource(appEnv env.Config, identities exchange.Identities, appLogger *slog.Logger) (exchange.Source, error) {
	exchanger, err := exchange.New(exchange.Config{
		Issuer:       appEnv.OIDC.Issuer,
		ClientID:     appEnv.OIDC.ClientID,
		ClientSecret: appEnv.OIDC.ClientSecret,
		Identities:   identities,
		Logger:       appLogger,
	})
	if err != nil {
		return exchange.Source{}, err
	}
	return exchanger.Source(exchange.Target{Audience: "bulle", Scopes: []string{"bulle.read", "bulle.write"}}), nil
}

// upgradeURL is Registre's plans page, where a person gets ai.access.
func upgradeURL(appEnv env.Config) string {
	return strings.TrimRight(appEnv.OIDC.Issuer, "/") + "/spaces"
}

// mountIdeas mounts the Ideas routes behind authentication on the /api
// router when they are configured.
func mountIdeas(router chi.Router, db *gorm.DB, svc services, authService *auth.Service) error {
	if svc.ideas == nil {
		return nil
	}
	var mountErr error
	router.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(authService))
		mountErr = ideas.New(db, svc.plans).Mount(r, *svc.ideas)
	})
	return mountErr
}
