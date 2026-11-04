package spaces

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/porte/directory"
)

// TestEntitlementsRideTheSamePull runs the chained apply on every snapshot
// the Syncer pulls, keeps syncing spaces when it fails, and reports its error
// so the Syncer pulls again.
func TestEntitlementsRideTheSamePull(t *testing.T) {
	dir, spaceID := directoryFixture(t)
	failure := stderrors.New("save failed")
	var seen int64
	dir.WithEntitlements(func(_ context.Context, snap directory.Snapshot) error {
		seen = snap.Version
		return failure
	})
	if err := dir.syncer.Apply(context.Background(), snapshot()); !stderrors.Is(err, failure) {
		t.Fatalf("chained apply = %v, want the entitlements error", err)
	}
	var space schemas.Space
	testORM.First(&space, "id = ?", spaceID)
	if seen != snapshot().Version || space.Name != "STUDIO" {
		t.Fatalf("entitlements saw version %d, space %+v", seen, space)
	}
}
