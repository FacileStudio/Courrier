package schemas

import (
	"log"
	"os"
	"testing"

	"github.com/FacileStudio/tronc/testdb"
	"gorm.io/gorm"
)

var testORM *gorm.DB

// Postgres and nothing else: these tests assert on regexp_replace, `!~` and COALESCE/NULLIF
// expressions that ship to production, and SQLite would build a different schema from the
// struct tags and then agree with itself about it.
func TestMain(m *testing.M) {
	url, configured := testdb.URL()
	if !configured {
		testdb.Announce("createdb courrier_test  # or point at any scratch database")
	} else {
		db, err := testdb.Open(url, testdb.Config{Prefix: "courrier_test", Migrate: Migrate})
		if err != nil {
			log.Fatalf("testdb: %v", err)
		}
		testORM = db
	}
	os.Exit(m.Run())
}

func openTestDatabase(t *testing.T) *gorm.DB {
	t.Helper()
	if testORM == nil {
		t.Skip(testdb.SkipReason("createdb courrier_test  # or point at any scratch database"))
	}
	if err := testORM.Exec(`DELETE FROM users`).Error; err != nil {
		t.Fatalf("clear users: %v", err)
	}
	return testORM
}
