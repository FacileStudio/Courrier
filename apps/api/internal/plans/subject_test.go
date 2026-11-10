package plans

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/FacileStudio/porte/directory"
)

// identityAt adds a Registre identity to userID, synced hoursAgo hours ago.
func identityAt(t *testing.T, userID int64, subject string, hoursAgo int) {
	t.Helper()
	if err := testORM.Exec(`INSERT INTO porte_identities (user_id, provider, subject, synced_at, created_at)
		VALUES (?, ?, ?, now() - make_interval(hours => ?), now() - interval '30 days')`,
		userID, testIssuer, subject, hoursAgo).Error; err != nil {
		t.Fatal(err)
	}
}

// TestPersonalPlanPrefersTheSubjectRegistreLists reads the plan of the
// identity the snapshot knows, even when a stale identity synced later, and
// keeps that choice after a restart.
func TestPersonalPlanPrefersTheSubjectRegistreLists(t *testing.T) {
	requireDB(t)
	alice := newUser(t, "alice@facile.studio", "")
	identityAt(t, alice, "alice", 48)
	identityAt(t, alice, "stale", 1)
	addAccounts(t, alice, nil, 2)
	apply(t, newGate(), entitled(1, nil, map[string]int{"alice": 3}))
	restarted := newGate()
	if err := restarted.Load(context.Background()); err != nil {
		t.Fatal(err)
	}
	expectAllowed(t, check(restarted, alice, nil), "alice's plan of 3 through her listed identity")
	addAccounts(t, alice, nil, 1)
	expectLimit(t, check(restarted, alice, nil), 3, testIssuer+"/spaces", "alice at her plan of 3")
}

// TestPersonalSubjectFallsBackToTheMostRecent picks the most recently synced
// identity when the snapshot lists none of them, the same one every time.
func TestPersonalSubjectFallsBackToTheMostRecent(t *testing.T) {
	requireDB(t)
	bob := newUser(t, "bob@facile.studio", "")
	identityAt(t, bob, "b-old", 48)
	identityAt(t, bob, "b-new", 1)
	gate := newGate()
	apply(t, gate, entitled(1, nil, map[string]int{"someone": 5}))
	for range 3 {
		subject, err := gate.personalSubject(testORM, bob)
		if err != nil || subject != "b-new" {
			t.Fatalf("personalSubject = %q, %v, want b-new", subject, err)
		}
	}
}

// TestSnapshotVersionGoingBackIsIgnored pins porte/forfait v0.1.1: a
// snapshot with a lower version is dropped at the same revision, but a higher
// revision wins whatever its version, so a Registre restored from backup
// moves Courrier to the restored values on the next pull.
func TestSnapshotVersionGoingBackIsIgnored(t *testing.T) {
	requireDB(t)
	alice := newUser(t, "alice@facile.studio", "alice")
	addAccounts(t, alice, nil, 1)
	gate := newGate()
	apply(t, gate, versioned(10, 2, 1))
	apply(t, gate, versioned(5, 2, 5))
	expectLimit(t, check(gate, alice, nil), 1, testIssuer+"/spaces", "after an older version")
	apply(t, gate, versioned(5, 3, 5))
	expectAllowed(t, check(gate, alice, nil), "after a newer revision")
}

// versioned is a snapshot giving alice a personal mailbox limit.
func versioned(version, revision int64, limit int) directory.Snapshot {
	raw, _ := json.Marshal(limit)
	return directory.Snapshot{Version: version, EntitlementsRevision: revision, People: []directory.Person{
		{Subject: "alice", Entitlements: map[string]json.RawMessage{Mailboxes: raw}},
	}}
}
