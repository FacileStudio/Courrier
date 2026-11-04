package plans

import (
	"github.com/FacileStudio/porte/forfait"
	"github.com/FacileStudio/tronc/errors"

	"gorm.io/gorm"
)

// unlinkedSpace refuses a new account in a space Registre does not know yet:
// its plan cannot be read, and the creator's personal plan must never govern
// a shared space.
func unlinkedSpace() error {
	return errors.Failed("this space is not linked to Registre yet, so it cannot take a new mail account")
}

// CheckMailbox reports whether one more mail account fits the plan of the
// scope it would be created in: the space when spaceID is set, else the
// user's personal scope. Call it inside the transaction that creates the
// account: it locks the space or user row, so two requests at the ceiling
// cannot both pass. It answers a *forfait.LimitError over the limit.
func (g *Gate) CheckMailbox(tx *gorm.DB, userID int64, spaceID *string) error {
	if !g.Active() {
		return nil
	}
	if spaceID != nil && *spaceID != "" {
		return g.checkSpace(tx, *spaceID)
	}
	return g.checkPersonal(tx, userID)
}

// checkSpace counts the accounts of a linked space against its plan.
func (g *Gate) checkSpace(tx *gorm.DB, spaceID string) error {
	var links []struct{ Link string }
	if err := tx.Raw(`SELECT coalesce(registre_space_id, '') AS link FROM spaces WHERE id = ? FOR UPDATE`, spaceID).
		Scan(&links).Error; err != nil {
		return errors.Internal("failed to read the space", err)
	}
	if len(links) == 0 {
		return errors.NotFound("space not found")
	}
	if links[0].Link == "" {
		return unlinkedSpace()
	}
	var used int64
	if err := tx.Raw(`SELECT count(*) FROM accounts WHERE space_id = ?`, spaceID).Scan(&used).Error; err != nil {
		return errors.Internal("failed to count the space's accounts", err)
	}
	return g.store.Check(forfait.Scope{SpaceID: links[0].Link}, Mailboxes, used, 1)
}

// checkPersonal counts the user's accounts outside any space against the plan
// of their Registre subject. A user with no Registre identity reads the free
// defaults.
func (g *Gate) checkPersonal(tx *gorm.DB, userID int64) error {
	var locked []int64
	if err := tx.Raw(`SELECT id FROM users WHERE id = ? FOR UPDATE`, userID).Scan(&locked).Error; err != nil {
		return errors.Internal("failed to read the user", err)
	}
	var subjects []string
	if err := tx.Raw(`SELECT subject FROM porte_identities WHERE user_id = ? AND provider = ?`,
		userID, g.issuer).Scan(&subjects).Error; err != nil {
		return errors.Internal("failed to read the user's identity", err)
	}
	var used int64
	if err := tx.Raw(`SELECT count(*) FROM accounts WHERE user_id = ? AND space_id IS NULL`, userID).
		Scan(&used).Error; err != nil {
		return errors.Internal("failed to count the user's accounts", err)
	}
	scope := forfait.Scope{}
	if len(subjects) > 0 {
		scope.Subject = subjects[0]
	}
	return g.store.Check(scope, Mailboxes, used, 1)
}
