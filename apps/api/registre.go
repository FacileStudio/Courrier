package main

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/FacileStudio/Courrier/apps/api/internal/env"
	"github.com/FacileStudio/Courrier/apps/api/internal/plans"
	"github.com/FacileStudio/Courrier/apps/api/modules/auth"
	"github.com/FacileStudio/Courrier/apps/api/modules/spaces"
	"github.com/FacileStudio/porte/directory"
	portepg "github.com/FacileStudio/porte/pg"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// newDirectory builds the Registre space sync over the OIDC client Courrier
// already signs users in with, feeding every snapshot to the plan gate too,
// or nil when SSO is off and Courrier manages its own spaces.
func newDirectory(db *gorm.DB, appEnv env.Config, gate *plans.Gate, appLogger *slog.Logger) *spaces.Directory {
	if appEnv.OIDC == nil {
		return nil
	}
	return spaces.NewDirectory(db, directory.Client{
		Issuer:       appEnv.OIDC.Issuer,
		ClientID:     appEnv.OIDC.ClientID,
		ClientSecret: appEnv.OIDC.ClientSecret,
	}, appLogger).WithEntitlements(gate.Apply)
}

// newDeletions builds the GDPR deletion receiver over the same client
// credentials as the snapshot sync, or nil when SSO is off.
func newDeletions(db *gorm.DB, appEnv env.Config, appLogger *slog.Logger) (*spaces.Deletions, error) {
	if appEnv.OIDC == nil {
		return nil, nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	client := directory.Client{
		Issuer:       appEnv.OIDC.Issuer,
		ClientID:     appEnv.OIDC.ClientID,
		ClientSecret: appEnv.OIDC.ClientSecret,
	}
	return spaces.NewDeletions(db, client, portepg.New(sqlDB).EraseIdentity, appEnv.StorageDir, appLogger), nil
}

// newPlans builds the plan gate and restores the plan values saved before the
// last restart, or nil when SSO is off: without Registre no plan applies.
func newPlans(ctx context.Context, db *gorm.DB, appEnv env.Config, appLogger *slog.Logger) (*plans.Gate, error) {
	if appEnv.OIDC == nil {
		return nil, nil
	}
	gate := plans.New(db, appEnv.OIDC.Issuer, appLogger)
	return gate, gate.Load(ctx)
}

// registerSpaces mounts the spaces routes and, with SSO on, hands spaces to
// Registre and serves its nudge at directory.SyncPath and its signed deletion
// notice at directory.DeletionsPath inside the /api subtree. Both sit outside
// every auth middleware on purpose: Registre posts anonymously, the nudge only
// schedules an authenticated pull, and the deletion handler verifies the
// notice's signature itself.
func registerSpaces(r chi.Router, db *gorm.DB, appEnv env.Config, dir *spaces.Directory, deletions *spaces.Deletions, authService *auth.Service) {
	service := spaces.NewService(db)
	spaces.RegisterRoutes(r, service, authService)
	if dir == nil {
		return
	}
	service.WithRegistre(appEnv.OIDC.Issuer)
	r.Method(http.MethodPost, strings.TrimPrefix(directory.SyncPath, "/api"), dir.Handler())
	if deletions != nil {
		r.Method(http.MethodPost, strings.TrimPrefix(directory.DeletionsPath, "/api"), deletions.Handler())
	}
}
