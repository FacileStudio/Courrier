package main

import (
	"context"

	"github.com/FacileStudio/Courrier/apps/api/internal/env"
	"github.com/FacileStudio/Courrier/apps/api/internal/plans"
	"github.com/FacileStudio/Courrier/apps/api/modules/assist"
	"github.com/FacileStudio/Courrier/apps/api/modules/mail"

	"gorm.io/gorm"
)

// newAssist builds the Bulle writing helper over the Bulle link. It stays off
// (every call answers "off") until BULLE_URL and the Registre client are both
// configured.
func newAssist(db *gorm.DB, appEnv env.Config, link bulleLink, gate *plans.Gate, mailService *mail.Service) *assist.Service {
	if link.client == nil {
		return assist.NewService(assist.Config{})
	}
	return assist.NewService(assist.Config{
		Bulle:      link.client,
		Linked:     link.linked,
		Mailbox:    mailService.MailboxSpace,
		AIAccess:   aiAccess(db, gate),
		BulleURL:   appEnv.BulleURL,
		UpgradeURL: upgradeURL(appEnv),
	})
}

// aiAccess reads whether the user's personal plan, or the plan of the
// Registre space link the mailbox belongs to, includes ai.access: the rule
// Bulle applies to the space_id Courrier sends it.
func aiAccess(db *gorm.DB, gate *plans.Gate) func(context.Context, int64, string) (bool, error) {
	return func(ctx context.Context, userID int64, link string) (bool, error) {
		return gate.AIAccessIn(db.WithContext(ctx), userID, link)
	}
}
