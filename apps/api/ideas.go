package main

import (
	"log/slog"
	"strings"

	"github.com/FacileStudio/Courrier/apps/api/internal/env"
	"github.com/FacileStudio/Courrier/apps/api/internal/ideas"
	"github.com/FacileStudio/Courrier/apps/api/internal/middleware"
	"github.com/FacileStudio/Courrier/apps/api/modules/auth"

	"github.com/FacileStudio/porte/exchange"
	"github.com/FacileStudio/tronc/bulle"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// newIdeasOptions wires the Ideas routes to Bulle: the app acts as the
// signed-in user through a Registre token exchange for audience bulle, over
// the identity store that seals the Registre tokens. It answers nil when SSO
// or BULLE_URL is off.
func newIdeasOptions(appEnv env.Config, identities exchange.Identities, appLogger *slog.Logger) (*bulle.Options, error) {
	if appEnv.OIDC == nil || appEnv.BulleURL == "" {
		return nil, nil
	}
	exchanger, err := exchange.New(exchange.Config{
		Issuer:       appEnv.OIDC.Issuer,
		ClientID:     appEnv.OIDC.ClientID,
		ClientSecret: appEnv.OIDC.ClientSecret,
		Identities:   identities,
		Logger:       appLogger,
	})
	if err != nil {
		return nil, err
	}
	source := exchanger.Source(exchange.Target{Audience: "bulle", Scopes: []string{"bulle.read", "bulle.write"}})
	return &bulle.Options{
		URL:         appEnv.BulleURL,
		TokenSource: source.Token,
		Forget:      source.Forget,
		UpgradeURL:  strings.TrimRight(appEnv.OIDC.Issuer, "/") + "/spaces",
	}, nil
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
