package main

import (
	"context"
	"github.com/FacileStudio/Courrier/apps/api/modules/assist"
	"log/slog"

	"github.com/FacileStudio/Courrier/apps/api/internal/env"
	"github.com/FacileStudio/Courrier/apps/api/internal/plans"
	"github.com/FacileStudio/Courrier/apps/api/internal/push"
	"github.com/FacileStudio/Courrier/apps/api/modules/mail"
	"github.com/FacileStudio/Courrier/apps/api/modules/spaces"
	"github.com/FacileStudio/Courrier/apps/api/schemas"

	"github.com/FacileStudio/porte/oidc"
	"github.com/FacileStudio/porte/session"

	"gorm.io/gorm"
)

// services are the long-lived pieces run builds once: porte's session manager
// and OIDC kit, the mail service the watcher shares with the HTTP routes, the
// push hub, and the Registre space sync, deletion receiver and plan gate when
// SSO is on.
type services struct {
	sessions  *session.Manager
	kit       *oidc.Kit
	mail      *mail.Service
	push      *push.Hub
	directory *spaces.Directory
	deletions *spaces.Deletions
	plans     *plans.Gate
	assist    *assist.Service
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
	if svc.deletions, err = newDeletions(db, appEnv, appLogger); err != nil {
		return services{}, err
	}
	reportSpaceIDs(db, appLogger)
	return svc, nil
}

// reportSpaceIDs logs rows whose space_id is not canonical, which count
// toward no plan and match no space filter. It changes nothing.
func reportSpaceIDs(db *gorm.DB, appLogger *slog.Logger) {
	found, err := schemas.NonCanonicalSpaceIDs(db)
	if err != nil {
		appLogger.Warn("space_id audit failed", slog.Any("error", err))
		return
	}
	if len(found) > 0 {
		appLogger.Warn("rows with a non-canonical space_id, left as they are", slog.Any("tables", found))
	}
}
