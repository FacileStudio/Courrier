package plans

import (
	"github.com/FacileStudio/porte/forfait"

	"gorm.io/gorm"
)

// SpaceAIAccess reports whether the Registre space's plan includes ai.access.
func (g *Gate) SpaceAIAccess(link string) bool {
	return g.store.Enabled(forfait.Scope{SpaceID: link}, AIAccess)
}

// AIAccessIn reports whether the user's personal plan, or the plan of the
// Registre space link the work is in, includes ai.access: Bulle grants AI on
// the same rule. An empty link is the personal scope; the caller has already
// checked that the user belongs to the space.
func (g *Gate) AIAccessIn(db *gorm.DB, userID int64, link string) (bool, error) {
	if link != "" && g.SpaceAIAccess(link) {
		return true, nil
	}
	return g.PersonalAIAccess(db, userID)
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
