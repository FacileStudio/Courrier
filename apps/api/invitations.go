package main

import (
	"log/slog"

	"github.com/FacileStudio/Courrier/apps/api/internal/env"
	"github.com/FacileStudio/Courrier/apps/api/internal/middleware"
	"github.com/FacileStudio/Courrier/apps/api/modules/auth"
	"github.com/FacileStudio/porte/directory"
	"github.com/FacileStudio/porte/kit"

	"github.com/go-chi/chi/v5"
)

// mountInvitations serves GET /api/me/invitations behind the session when
// SSO is on: the signed-in person's pending Registre space invitations for
// the user menu, asked with the directory client the space sync uses.
func mountInvitations(router chi.Router, appEnv env.Config, svc services, authService *auth.Service, appLogger *slog.Logger) error {
	if appEnv.OIDC == nil {
		return nil
	}
	return kit.MountInvitations(router.With(middleware.RequireAuth(authService)), kit.InvitationsConfig{
		Directory: directory.Client{
			Issuer:       appEnv.OIDC.Issuer,
			ClientID:     appEnv.OIDC.ClientID,
			ClientSecret: appEnv.OIDC.ClientSecret,
		},
		Identities: svc.identity,
		Logger:     appLogger,
	})
}
