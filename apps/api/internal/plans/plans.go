// Package plans enforces the plan a Registre space pays for. It feeds every
// directory snapshot to a porte/forfait Store, keeps the Store's state in one
// database row so a restart never fails open again, and answers whether a new
// mail account fits the plan of the scope it is created in.
package plans

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync/atomic"

	"github.com/FacileStudio/porte/directory"
	"github.com/FacileStudio/porte/forfait"

	"gorm.io/gorm"
)

// Mailboxes is the feature that caps how many mail accounts one scope holds.
const Mailboxes = "courrier.mailboxes"

// Gate is Courrier's view of the plans Registre serves. A nil Gate, or one
// that never received an entitlements snapshot, allows everything.
type Gate struct {
	orm    *gorm.DB
	issuer string
	logger *slog.Logger
	store  *forfait.Store
	active atomic.Bool
}

// New builds the Gate for the OIDC issuer Courrier signs users in with. The
// issuer is both the provider of the porte identities a personal scope is read
// from and the origin of every upgrade link.
func New(orm *gorm.DB, issuer string, logger *slog.Logger) *Gate {
	store := forfait.NewStore(forfait.Feature{Key: Mailboxes, Kind: forfait.Limit, Default: json.RawMessage("1")})
	store.Issuer = issuer
	store.Logger = logger
	return &Gate{orm: orm, issuer: issuer, logger: logger, store: store}
}

// Load restores the state the last Apply saved. Call it at boot, before the
// first pull. A row that cannot be imported is logged and the Gate fails
// closed at the feature defaults until the next snapshot: a row exists, so
// entitlements were authoritative once, and failing open would lift every
// limit.
func (g *Gate) Load(ctx context.Context) error {
	var documents []string
	if err := g.orm.WithContext(ctx).Raw(`SELECT document FROM forfait_state WHERE id = 1`).
		Scan(&documents).Error; err != nil {
		return err
	}
	if len(documents) == 0 {
		return nil
	}
	if err := g.store.Import([]byte(documents[0])); err != nil {
		g.logger.Error("plans: saved entitlements could not be restored, enforcing defaults", slog.Any("error", err))
		g.store.Apply(directory.Snapshot{EntitlementsRevision: 1})
		g.active.Store(true)
		return nil
	}
	g.refresh([]byte(documents[0]))
	return nil
}

// Apply feeds one directory snapshot to the Store and saves the result. It is
// chained into the directory Syncer's Apply, so a failed save makes the
// Syncer pull and save again.
func (g *Gate) Apply(ctx context.Context, snapshot directory.Snapshot) error {
	g.store.Apply(snapshot)
	document, err := g.store.Export()
	if err != nil {
		return err
	}
	if err := g.orm.WithContext(ctx).Exec(`INSERT INTO forfait_state (id, document, updated_at)
		VALUES (1, ?, now())
		ON CONFLICT (id) DO UPDATE SET document = EXCLUDED.document, updated_at = now()`,
		string(document)).Error; err != nil {
		return err
	}
	g.refresh(document)
	return nil
}

// Active reports whether entitlements are authoritative: an entitlements
// snapshot was applied, now or before the last restart.
func (g *Gate) Active() bool {
	return g != nil && g.active.Load()
}

// refresh marks the Gate active once an exported document carries a revision.
func (g *Gate) refresh(document []byte) {
	var state struct {
		Revision int64 `json:"entitlements_revision"`
	}
	if json.Unmarshal(document, &state) == nil && state.Revision > 0 {
		g.active.Store(true)
	}
}
