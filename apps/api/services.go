package main

import (
	"context"
	"log/slog"

	"github.com/FacileStudio/Courrier/apps/api/internal/env"
	"github.com/FacileStudio/Courrier/apps/api/internal/plans"
	"github.com/FacileStudio/Courrier/apps/api/internal/push"
	"github.com/FacileStudio/Courrier/apps/api/modules/mail"
	"github.com/FacileStudio/Courrier/apps/api/modules/spaces"

	"github.com/FacileStudio/porte/local"
	"github.com/FacileStudio/porte/oidc"
	"github.com/FacileStudio/porte/session"

	"gorm.io/gorm"
)

// services are the long-lived pieces run builds once: porte's session manager,
// password login and OIDC kit, the mail service the watcher shares with the
// HTTP routes, the push hub, and the Registre space sync and plan gate when
// SSO is on.
type services struct {
	sessions  *session.Manager
	passwords *local.Kit
	kit       *oidc.Kit
	mail      *mail.Service
	push      *push.Hub
	directory *spaces.Directory
	plans     *plans.Gate
}

// buildServices builds porte, then the mail service with its push hub, and,
// with SSO on, the plan gate restored from its saved state and the Registre
// sync that feeds it.
func buildServices(ctx context.Context, db *gorm.DB, appEnv env.Config, appLogger *slog.Logger) (services, error) {
	svc, err := buildAuth(ctx, db, appEnv, appLogger)
	if err != nil {
		return services{}, err
	}
	svc.push = newPushHub(db, appEnv.Push, appLogger)
	svc.mail = mail.NewService(db, appEnv.EncryptionKey).WithNotifier(svc.push)
	if svc.plans, err = newPlans(ctx, db, appEnv, appLogger); err != nil {
		return services{}, err
	}
	svc.directory = newDirectory(db, appEnv, svc.plans, appLogger)
	return svc, nil
}
