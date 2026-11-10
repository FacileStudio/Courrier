package plans

import (
	"github.com/FacileStudio/porte/forfait"

	"gorm.io/gorm"
)

// SpaceAIAccess reports whether the Registre space's plan includes ai.access.
func (g *Gate) SpaceAIAccess(link string) bool {
	return g.store.Enabled(forfait.Scope{SpaceID: link}, AIAccess)
}

// PersonalAIAccess reports whether the user's personal plan includes
// ai.access.
func (g *Gate) PersonalAIAccess(db *gorm.DB, userID int64) (bool, error) {
	subject, err := g.personalSubject(db, userID)
	if err != nil {
		return false, err
	}
	return g.store.Enabled(forfait.Scope{Subject: subject}, AIAccess), nil
}
