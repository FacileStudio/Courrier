package main

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/FacileStudio/Courrier/apps/api/internal/env"
	"github.com/FacileStudio/Courrier/apps/api/modules/auth"
	"github.com/FacileStudio/Courrier/apps/api/modules/mcp"
	"github.com/FacileStudio/porte/mcpauth"
	"github.com/FacileStudio/tronc/apiref"
	"github.com/FacileStudio/tronc/health"
	mcpkit "github.com/FacileStudio/tronc/mcp"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// mountRoot mounts the routes outside /api: health, version, the API
// reference, and the MCP endpoint with its metadata when MCP is on.
func mountRoot(router chi.Router, dbCheck health.Check, svc services) {
	health.Mount(router, dbCheck)
	router.Get("/version", versionHandler)
	apiref.Mount(router, referenceConfig())
	if svc.mcp != nil {
		mcp.RegisterRoutes(router, svc.mcp)
	}
}

// mountPending mounts the confirm page's API on the /api router when MCP is
// on.
func mountPending(router chi.Router, svc services, authService *auth.Service) {
	if svc.mcp != nil {
		mcp.RegisterPendingRoutes(router, svc.mcp, authService)
	}
}

// newMCP builds the MCP server when MCP_RESOURCE is set and SSO is on, and
// nothing otherwise. The verifier runs Registre discovery here, so a wrong
// issuer fails at boot like the OIDC kit does. A person's first MCP call
// creates their account through the OIDC kit, as a first web login would.
func newMCP(ctx context.Context, db *gorm.DB, appEnv env.Config, svc services, appLogger *slog.Logger) (*mcp.Service, error) {
	if appEnv.OIDC == nil || appEnv.MCP.Resource == "" {
		return nil, nil
	}
	verifier, err := mcpauth.Verifier(ctx, mcpauth.Config{
		Issuer:      appEnv.OIDC.Issuer,
		Resource:    appEnv.MCP.Resource,
		ScopePrefix: "courrier",
		Identities:  svc.identity,
		Provision:   true,
		Accounts:    svc.kit,
		Logger:      appLogger,
	})
	if err != nil {
		return nil, err
	}
	appLogger.Info("mcp server on", slog.String("resource", appEnv.MCP.Resource))
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	return mcp.NewService(svc.mail, mcp.Config{
		Resource: appEnv.MCP.Resource,
		Issuer:   appEnv.OIDC.Issuer,
		Version:  version,
		Origins:  appEnv.CORSAllowedOrigins,
		StateKey: appEnv.MCP.StateKey,
		Verifier: verifier,
		Pending:  mcpkit.PostgresPending{DB: sqlDB},
		Logger:   appLogger,
	})
}

// keepPendingPruned deletes, hourly until ctx ends, the confirmations
// resolved or expired more than 30 days ago, when MCP is on.
func keepPendingPruned(ctx context.Context, sqlDB *sql.DB, svc services, appLogger *slog.Logger) {
	if svc.mcp != nil {
		go mcpkit.KeepPruned(ctx, mcpkit.PostgresPending{DB: sqlDB}, 0, 0, appLogger)
	}
}
