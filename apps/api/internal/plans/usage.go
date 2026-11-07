package plans

import (
	"github.com/FacileStudio/porte/forfait"

	"gorm.io/gorm"
)

// Usage is how many mail accounts a scope holds against what its plan allows.
// Limit is -1 when the scope is unlimited or no entitlements apply yet.
type Usage struct {
	Used  int64
	Limit int64
}

// MailboxUsage reads, without locking, the accounts the scope holds and the
// ceiling its plan sets: the space when spaceID is set, else the user's
// personal scope. A space not linked to Registre reads a ceiling of 0, the
// same as CheckMailbox refusing it. It never blocks a read: an over-limit
// scope still sees its numbers.
func (g *Gate) MailboxUsage(db *gorm.DB, userID int64, spaceID *string) (Usage, error) {
	if spaceID != nil && *spaceID != "" {
		return g.spaceUsage(db, *spaceID)
	}
	return g.personalUsage(db, userID)
}

func (g *Gate) spaceUsage(db *gorm.DB, spaceID string) (Usage, error) {
	space, err := lockSpace(db, spaceID)
	if err != nil {
		return Usage{}, err
	}
	usage := Usage{Limit: -1}
	if err := db.Raw(`SELECT count(*) FROM accounts WHERE space_id = ?`, space.ID).Scan(&usage.Used).Error; err != nil {
		return Usage{}, err
	}
	if !g.Active() {
		return usage, nil
	}
	if space.Link == "" {
		usage.Limit = 0
		return usage, nil
	}
	usage.Limit = g.store.Limit(forfait.Scope{SpaceID: space.Link}, Mailboxes)
	return usage, nil
}

func (g *Gate) personalUsage(db *gorm.DB, userID int64) (Usage, error) {
	usage := Usage{Limit: -1}
	if err := db.Raw(`SELECT count(*) FROM accounts WHERE user_id = ? AND space_id IS NULL`, userID).
		Scan(&usage.Used).Error; err != nil {
		return Usage{}, err
	}
	if !g.Active() {
		return usage, nil
	}
	subject, err := g.personalSubject(db, userID)
	if err != nil {
		return Usage{}, err
	}
	usage.Limit = g.store.Limit(forfait.Scope{Subject: subject}, Mailboxes)
	return usage, nil
}
