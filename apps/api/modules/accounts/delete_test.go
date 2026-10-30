package accounts

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/FacileStudio/tronc/testdb"
	"gorm.io/gorm"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
)

var testORM *gorm.DB

func TestMain(m *testing.M) {
	url, configured := testdb.URL()
	if configured {
		db, err := testdb.Open(url, testdb.Config{Prefix: "courrier_accounts_test", Migrate: schemas.Migrate})
		if err != nil {
			log.Fatalf("testdb: %v", err)
		}
		testORM = db
	}
	os.Exit(m.Run())
}

func TestDeleteRemovesContacts(t *testing.T) {
	if testORM == nil {
		t.Skip(testdb.SkipReason("createdb courrier_test  # or point at any scratch database"))
	}
	testORM.Exec("DELETE FROM contacts")
	testORM.Exec("DELETE FROM accounts")
	acc := schemas.Account{UserID: 1, Name: "A", Email: "a@example.test"}
	other := schemas.Account{UserID: 1, Name: "B", Email: "b@example.test"}
	testORM.Create(&acc)
	testORM.Create(&other)
	for _, id := range []int64{acc.ID, other.ID} {
		testORM.Create(&schemas.Contact{AccountID: id, Email: "x@x.test", LastSentAt: time.Now()})
	}
	if err := NewService(testORM, nil).Delete(context.Background(), 1, acc.ID); err != nil {
		t.Fatal(err)
	}
	var left []schemas.Contact
	testORM.Find(&left)
	if len(left) != 1 || left[0].AccountID != other.ID {
		t.Fatalf("only the other account's contacts must remain: %+v", left)
	}
}
