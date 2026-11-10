package mcp

import (
	"context"
	"crypto/rand"
	"log"
	"os"
	"testing"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/tronc/testdb"
	"gorm.io/gorm"
)

// subject is the Registre sub of the test person.
const subject = "person-sub"

// testORM is the scratch database, nil when TEST_DATABASE_URL is unset.
var testORM *gorm.DB

// TestMain opens the scratch database when one is configured.
func TestMain(m *testing.M) {
	if url, configured := testdb.URL(); configured {
		db, err := testdb.Open(url, testdb.Config{Prefix: "courrier_mcp_test", Migrate: schemas.Migrate})
		if err != nil {
			log.Fatalf("testdb: %v", err)
		}
		testORM = db
	}
	os.Exit(m.Run())
}

// identities maps the test subject to userID and nothing else.
type identities struct {
	porte.IdentityStore
	userID int64
}

// Find resolves only the test subject.
func (s identities) Find(_ context.Context, provider, sub string) (porte.StoredIdentity, error) {
	if sub != subject {
		return porte.StoredIdentity{}, porte.ErrNotFound
	}
	return porte.StoredIdentity{UserID: s.userID, Provider: provider, Subject: sub}, nil
}

// seedPerson creates the user and account rows pending actions point at,
// or returns fixed ids without a database.
func seedPerson(t *testing.T) (userID, accountID int64) {
	t.Helper()
	if testORM == nil {
		return 7, 70
	}
	user := schemas.User{Email: t.Name() + "-" + rand.Text() + "@example.test"}
	if err := testORM.Create(&user).Error; err != nil {
		t.Fatal(err)
	}
	account := schemas.Account{UserID: user.ID, Email: "me@example.test", Name: "Me"}
	if err := testORM.Create(&account).Error; err != nil {
		t.Fatal(err)
	}
	return user.ID, account.ID
}
