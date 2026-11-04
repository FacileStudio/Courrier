package plans

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"os"
	"testing"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/porte/directory"
	"github.com/FacileStudio/tronc/testdb"

	"gorm.io/gorm"
)

const testIssuer = "https://sso.test"

var testORM *gorm.DB

func TestMain(m *testing.M) {
	url, configured := testdb.URL()
	if configured {
		db, err := testdb.Open(url, testdb.Config{Prefix: "courrier_plans_test", Migrate: schemas.Migrate})
		if err != nil {
			log.Fatalf("testdb: %v", err)
		}
		testORM = db
	}
	os.Exit(m.Run())
}

// requireDB skips without Postgres and empties every table these tests write.
func requireDB(t *testing.T) {
	t.Helper()
	if testORM == nil {
		t.Skip(testdb.SkipReason("createdb courrier_test  # or point at any scratch database"))
	}
	for _, table := range []string{"forfait_state", "accounts", "space_members", "spaces", "porte_identities", "users"} {
		if err := testORM.Exec("DELETE FROM " + table).Error; err != nil {
			t.Fatal(err)
		}
	}
}

// newGate is a Gate over the test database that logs nowhere.
func newGate() *Gate {
	return New(testORM, testIssuer, slog.New(slog.NewTextHandler(io.Discard, nil)))
}

// newUser creates a user and, when subject is set, their Registre identity.
func newUser(t *testing.T, email, subject string) int64 {
	t.Helper()
	row := schemas.User{Email: email}
	if err := testORM.Create(&row).Error; err != nil {
		t.Fatal(err)
	}
	if subject != "" {
		if err := testORM.Exec(`INSERT INTO porte_identities (user_id, provider, subject) VALUES (?, ?, ?)`,
			row.ID, testIssuer, subject).Error; err != nil {
			t.Fatal(err)
		}
	}
	return row.ID
}

// newSpace creates a local space, linked to registreID unless it is empty.
func newSpace(t *testing.T, name, registreID string) string {
	t.Helper()
	row := schemas.Space{Name: name}
	if registreID != "" {
		row.RegistreSpaceID = &registreID
	}
	if err := testORM.Create(&row).Error; err != nil {
		t.Fatal(err)
	}
	return row.ID
}

// addAccounts stores n accounts for userID in spaceID, or personal when nil.
func addAccounts(t *testing.T, userID int64, spaceID *string, n int) {
	t.Helper()
	for range n {
		if err := testORM.Create(&schemas.Account{UserID: userID, SpaceID: spaceID, Name: "x"}).Error; err != nil {
			t.Fatal(err)
		}
	}
}

// check runs CheckMailbox in a transaction, the way account creation does.
func check(gate *Gate, userID int64, spaceID *string) error {
	return testORM.Transaction(func(tx *gorm.DB) error {
		return gate.CheckMailbox(tx, userID, spaceID)
	})
}

// entitled is a snapshot at revision with one limit per space and person.
func entitled(revision int64, spaces, people map[string]int) directory.Snapshot {
	snap := directory.Snapshot{Version: revision + 100, EntitlementsRevision: revision}
	for id, limit := range spaces {
		snap.Spaces = append(snap.Spaces, directory.Space{ID: id, Entitlements: limitOf(limit)})
	}
	for sub, limit := range people {
		snap.People = append(snap.People, directory.Person{Subject: sub, Entitlements: limitOf(limit)})
	}
	return snap
}

// limitOf is the entitlements of a scope with that mailbox limit, or none
// when limit is 0.
func limitOf(limit int) map[string]json.RawMessage {
	if limit == 0 {
		return nil
	}
	raw, _ := json.Marshal(limit)
	return map[string]json.RawMessage{Mailboxes: raw}
}

// apply feeds snap to gate and fails the test on a save error.
func apply(t *testing.T, gate *Gate, snap directory.Snapshot) {
	t.Helper()
	if err := gate.Apply(context.Background(), snap); err != nil {
		t.Fatal(err)
	}
}
