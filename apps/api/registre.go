package main

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/FacileStudio/Courrier/apps/api/internal/env"
	"github.com/FacileStudio/Courrier/apps/api/modules/auth"
	"github.com/FacileStudio/Courrier/apps/api/modules/spaces"
	"github.com/FacileStudio/porte/directory"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// newDirectory builds the Registre space sync over the OIDC client Courrier
// already signs users in with, or nil when SSO is off and Courrier manages its
// own spaces.
func newDirectory(db *gorm.DB, appEnv env.Config, appLogger *slog.Logger) *spaces.Directory {
	if appEnv.OIDC == nil {
		return nil
	}
	return spaces.NewDirectory(db, directory.Client{
		Issuer:       appEnv.OIDC.Issuer,
		ClientID:     appEnv.OIDC.ClientID,
		ClientSecret: appEnv.OIDC.ClientSecret,
	}, appLogger)
}

// registerSpaces mounts the spaces routes and, with SSO on, hands spaces to
// Registre and serves its nudge at directory.SyncPath inside the /api subtree.
// The nudge sits outside every auth middleware on purpose: Registre posts
// anonymously and the nudge only schedules an authenticated pull.
func registerSpaces(r chi.Router, db *gorm.DB, appEnv env.Config, dir *spaces.Directory, authService *auth.Service) {
	service := spaces.NewService(db)
	spaces.RegisterRoutes(r, service, authService)
	if dir == nil {
		return
	}
	service.WithRegistre(appEnv.OIDC.Issuer)
	r.Method(http.MethodPost, strings.TrimPrefix(directory.SyncPath, "/api"), dir.Handler())
}
