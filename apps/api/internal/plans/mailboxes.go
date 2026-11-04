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
// cannot both pass. It rewrites a set spaceID to the space's canonical id,
// which the caller stores. It answers a *forfait.LimitError over the limit.
func (g *Gate) CheckMailbox(tx *gorm.DB, userID int64, spaceID *string) error {
	if spaceID != nil && *spaceID != "" {
		return g.checkSpace(tx, spaceID)
	}
	if !g.Active() {
		return nil
	}
	return g.checkPersonal(tx, userID)
}

// lockedSpace is a space row read under FOR UPDATE: its canonical id and its
// Registre link, empty when unlinked.
type lockedSpace struct {
	ID   string
	Link string
}

// lockSpace locks the space row and reads its canonical id and link.
func lockSpace(tx *gorm.DB, spaceID string) (lockedSpace, error) {
	var rows []lockedSpace
	if err := tx.Raw(`SELECT id::text AS id, coalesce(registre_space_id, '') AS link
		FROM spaces WHERE id = ? FOR UPDATE`, spaceID).Scan(&rows).Error; err != nil {
		return lockedSpace{}, errors.Internal("failed to read the space", err)
	}
	if len(rows) == 0 {
		return lockedSpace{}, errors.NotFound("space not found")
	}
	return rows[0], nil
}

// checkSpace counts the accounts of a linked space, under its canonical id,
// against its plan.
func (g *Gate) checkSpace(tx *gorm.DB, spaceID *string) error {
	space, err := lockSpace(tx, *spaceID)
	if err != nil {
		return err
	}
	*spaceID = space.ID
	if !g.Active() {
		return nil
	}
	if space.Link == "" {
		return unlinkedSpace()
	}
	var used int64
	if err := tx.Raw(`SELECT count(*) FROM accounts WHERE space_id = ?`, space.ID).Scan(&used).Error; err != nil {
		return errors.Internal("failed to count the space's accounts", err)
	}
	return g.store.Check(forfait.Scope{SpaceID: space.Link}, Mailboxes, used, 1)
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
