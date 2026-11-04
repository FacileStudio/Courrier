package accounts

import (
	"context"
	stderrors "errors"
	"strings"
	"testing"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/porte/forfait"
	"github.com/FacileStudio/tronc/errors"
)

// TestEmptySpaceIDIsPersonal stores space_id "" as the personal scope, so it
// counts toward the personal plan instead of no plan at all.
func TestEmptySpaceIDIsPersonal(t *testing.T) {
	service, userID, _ := planFixture(t)
	ctx := context.Background()
	empty := ""
	for i := range 4 {
		account, err := service.Create(ctx, userID, mailbox(&empty))
		if i == 0 && (err != nil || account.SpaceID != nil) {
			t.Fatalf("first personal account: %v, space %v", err, account.SpaceID)
		}
		if i > 0 && !stderrors.Is(err, forfait.ErrLimit) {
			t.Fatalf("personal account %d over the default of 1: %v", i, err)
		}
	}
}

// TestSpaceIDSpellings counts every accepted spelling of one space together
// and refuses the others outright, so no spelling starts a fresh count.
func TestSpaceIDSpellings(t *testing.T) {
	service, userID, spaceID := planFixture(t)
	ctx := context.Background()
	upper := strings.ToUpper(spaceID)
	account, err := service.Create(ctx, userID, mailbox(&upper))
	if err != nil || account.SpaceID == nil || *account.SpaceID != spaceID {
		t.Fatalf("uppercase spelling: %v, stored %v, want %s", err, account.SpaceID, spaceID)
	}
	for _, spelling := range []string{spaceID, upper} {
		if _, err := service.Create(ctx, userID, mailbox(&spelling)); !stderrors.Is(err, forfait.ErrLimit) {
			t.Errorf("%q over the ceiling: %v, want ErrLimit", spelling, err)
		}
	}
	for _, spelling := range []string{"{" + spaceID + "}", strings.ReplaceAll(spaceID, "-", ""), " " + spaceID} {
		_, err := service.Create(ctx, userID, mailbox(&spelling))
		var appErr *errors.Error
		if !stderrors.As(err, &appErr) || appErr.Code != "invalid_argument" {
			t.Errorf("%q: %v, want invalid_argument", spelling, err)
		}
		if _, err := service.List(ctx, userID, &spelling); err == nil {
			t.Errorf("List accepted %q", spelling)
		}
	}
	var stored int64
	testORM.Model(&schemas.Account{}).Where("space_id = ?", spaceID).Count(&stored)
	if listed, _ := service.List(ctx, userID, &upper); stored != 1 || len(listed) != 1 {
		t.Fatalf("space holds %d accounts, uppercase list shows %d, want 1 and 1", stored, len(listed))
	}
}
