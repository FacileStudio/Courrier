package main

import (
	"context"

	"github.com/FacileStudio/Courrier/apps/api/internal/env"
	"github.com/FacileStudio/Courrier/apps/api/internal/plans"
	"github.com/FacileStudio/Courrier/apps/api/modules/assist"

	"gorm.io/gorm"
)

// newAssist builds the Bulle writing helper over the Bulle link. It stays off
// (every call answers "off") until BULLE_URL and the Registre client are both
// configured.
func newAssist(db *gorm.DB, appEnv env.Config, link bulleLink, gate *plans.Gate) *assist.Service {
	if link.client == nil {
		return assist.NewService(assist.Config{})
	}
	return assist.NewService(assist.Config{
		Bulle:      link.client,
		Linked:     link.linked,
		AIAccess:   personalAIAccess(db, gate),
		BulleURL:   appEnv.BulleURL,
		UpgradeURL: upgradeURL(appEnv),
	})
}

// personalAIAccess reads whether the user's personal plan includes ai.access,
// the plan Bulle checks while Courrier sends it no space.
func personalAIAccess(db *gorm.DB, gate *plans.Gate) func(context.Context, int64) (bool, error) {
	return func(ctx context.Context, userID int64) (bool, error) {
		return gate.PersonalAIAccess(db.WithContext(ctx), userID)
	}
}
