package plans

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/FacileStudio/porte/forfait"
	"github.com/FacileStudio/tronc/errors"
)

// TestFailOpenBeforeFirstSnapshot allows every creation, unlinked spaces
// included, until an entitlements revision arrives, and a revision-0
// snapshot does not end it.
func TestFailOpenBeforeFirstSnapshot(t *testing.T) {
	requireDB(t)
	gate := newGate()
	alice := newUser(t, "alice@facile.studio", "alice")
	local := newSpace(t, "Local", "")
	addAccounts(t, alice, nil, 5)
	apply(t, gate, entitled(0, nil, map[string]int{"alice": 1}))
	if gate.Active() {
		t.Fatal("a revision-0 snapshot activated enforcement")
	}
	if err := check(gate, alice, nil); err != nil {
		t.Fatalf("fail-open personal check = %v", err)
	}
	if err := check(gate, alice, &local); err != nil {
		t.Fatalf("fail-open unlinked space check = %v", err)
	}
	var nilGate *Gate
	if err := check(nilGate, alice, nil); err != nil {
		t.Fatalf("nil gate check = %v", err)
	}
}

// TestLimitPerScope counts each space and each person on its own: a space at
// its ceiling does not block another, and personal accounts never count
// toward a space.
func TestLimitPerScope(t *testing.T) {
	requireDB(t)
	gate := newGate()
	alice := newUser(t, "alice@facile.studio", "alice")
	bob := newUser(t, "bob@facile.studio", "bob")
	pro := newSpace(t, "Pro", "spc_pro")
	free := newSpace(t, "Free", "spc_free")
	apply(t, gate, entitled(1, map[string]int{"spc_pro": 2}, map[string]int{"alice": 3}))
	addAccounts(t, alice, &pro, 1)
	addAccounts(t, alice, nil, 2)
	expectAllowed(t, check(gate, alice, &pro), "pro space with room")
	expectAllowed(t, check(gate, alice, &free), "free space empty")
	addAccounts(t, alice, &free, 5)
	expectAllowed(t, check(gate, bob, &free), "free space at its unlimited default")
	addAccounts(t, alice, &pro, 1)
	expectLimit(t, check(gate, alice, &pro), 2, testIssuer+"/spaces/spc_pro", "pro space at its ceiling")
	expectAllowed(t, check(gate, alice, nil), "alice personal under 3")
	expectAllowed(t, check(gate, bob, nil), "bob personal empty")
	addAccounts(t, bob, nil, 5)
	expectAllowed(t, check(gate, bob, nil), "bob personal at the unlimited default")
	password := newUser(t, "carol@facile.studio", "")
	addAccounts(t, password, nil, 5)
	expectAllowed(t, check(gate, password, nil), "user without identity")
}

// TestUnlinkedSpaceRefused refuses a space Registre does not know once
// enforcement is active, instead of reading the creator's personal plan.
func TestUnlinkedSpaceRefused(t *testing.T) {
	requireDB(t)
	gate := newGate()
	alice := newUser(t, "alice@facile.studio", "alice")
	local := newSpace(t, "Local", "")
	apply(t, gate, entitled(1, nil, map[string]int{"alice": -1}))
	err := check(gate, alice, &local)
	var appErr *errors.Error
	if !stderrors.As(err, &appErr) || appErr.Code != "failed_precondition" {
		t.Fatalf("unlinked space check = %v, want failed_precondition", err)
	}
	missing := "00000000-0000-0000-0000-000000000000"
	if err := check(gate, alice, &missing); !stderrors.As(err, &appErr) || appErr.Code != "not_found" {
		t.Fatalf("missing space check = %v, want not_found", err)
	}
}

// TestStateSurvivesRestart reads the saved row in a fresh Gate, which then
// enforces the last known values without any pull.
func TestStateSurvivesRestart(t *testing.T) {
	requireDB(t)
	alice := newUser(t, "alice@facile.studio", "alice")
	pro := newSpace(t, "Pro", "spc_pro")
	apply(t, newGate(), entitled(4, map[string]int{"spc_pro": 1}, nil))
	addAccounts(t, alice, &pro, 1)
	restarted := newGate()
	if err := restarted.Load(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !restarted.Active() {
		t.Fatal("restored state did not activate enforcement")
	}
	expectLimit(t, check(restarted, alice, &pro), 1, testIssuer+"/spaces/spc_pro", "after restart")
	apply(t, restarted, entitled(3, map[string]int{"spc_pro": 5}, nil))
	expectLimit(t, check(restarted, alice, &pro), 1, testIssuer+"/spaces/spc_pro", "older revision after restart")
}

// TestCorruptRowFailsClosed enforces the defaults when the saved row cannot
// be imported, which leave mailboxes unlimited, and the next real snapshot
// replaces them.
func TestCorruptRowFailsClosed(t *testing.T) {
	requireDB(t)
	if err := testORM.Exec(`INSERT INTO forfait_state (id, document) VALUES (1, '{"version":9}')`).Error; err != nil {
		t.Fatal(err)
	}
	alice := newUser(t, "alice@facile.studio", "alice")
	addAccounts(t, alice, nil, 1)
	gate := newGate()
	if err := gate.Load(context.Background()); err != nil || !gate.Active() {
		t.Fatalf("Load = %v, active = %v", err, gate.Active())
	}
	expectAllowed(t, check(gate, alice, nil), "corrupt row at the unlimited default")
	apply(t, gate, entitled(1, nil, map[string]int{"alice": 1}))
	expectLimit(t, check(gate, alice, nil), 1, testIssuer+"/spaces", "first snapshot after a corrupt row")
}

// expectAllowed fails when err is not nil.
func expectAllowed(t *testing.T, err error, what string) {
	t.Helper()
	if err != nil {
		t.Errorf("%s: %v, want allowed", what, err)
	}
}

// expectLimit fails unless err is a LimitError with that ceiling and link.
func expectLimit(t *testing.T, err error, limit int64, upgradeURL, what string) {
	t.Helper()
	var refusal *forfait.LimitError
	if !stderrors.As(err, &refusal) || refusal.Key != Mailboxes || refusal.Limit != limit || refusal.UpgradeURL != upgradeURL {
		t.Errorf("%s: %v (%+v), want limit %d at %s", what, err, refusal, limit, upgradeURL)
	}
}
