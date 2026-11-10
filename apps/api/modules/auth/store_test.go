package auth

import (
	"context"
	"io"
	"log"
	"log/slog"
	"os"
	"testing"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/porte/kit"
	"github.com/FacileStudio/tronc/testdb"

	"gorm.io/gorm"
)

const (
	testBootstrap = "docker compose up -d db   # or any PostgreSQL 16"
	testIssuer    = "https://sso.test"
)

var testORM *gorm.DB

func TestMain(m *testing.M) {
	url, configured := testdb.URL()
	if !configured {
		testdb.Announce(testBootstrap)
	} else {
		db, err := testdb.Open(url, testdb.Config{Prefix: "courrier_auth_test", Migrate: schemas.Migrate})
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
		t.Skip(testdb.SkipReason(testBootstrap))
	}
	if err := testORM.Exec(`DELETE FROM users`).Error; err != nil {
		t.Fatalf("clear users: %v", err)
	}
	return testORM
}

func newStore(t *testing.T, db *gorm.DB) *kit.Store {
	t.Helper()
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db: %v", err)
	}
	store, err := NewUserStore(context.Background(), sqlDB, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatalf("NewUserStore: %v", err)
	}
	return store
}

func claims(subject, email string, verified bool) porte.Claims {
	return porte.Claims{Provider: testIssuer, Subject: subject, Email: email, EmailVerified: verified, Locale: "de"}
}

func link(t *testing.T, db *gorm.DB, userID int64, subject string) {
	t.Helper()
	err := db.Exec(`INSERT INTO porte_identities (user_id, provider, subject) VALUES (?, ?, ?)`, userID, testIssuer, subject).Error
	if err != nil {
		t.Fatalf("link identity: %v", err)
	}
}

// An account created before Registre owned sign-in has no federated identity.
// Its owner's first Facile sign-in must land on that same row, keep its data,
// and then resolve by subject. The session locale is stored on the way.
func TestFirstFacileSignInLinksAnExistingAccountByVerifiedEmail(t *testing.T) {
	db := openTestDatabase(t)
	store := newStore(t, db)
	ctx := context.Background()

	existing := schemas.User{Email: "camille@facile.studio", Name: "Camille"}
	if err := db.Create(&existing).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}

	id, err := store.UpsertFromOIDC(ctx, claims("sub-1", "Camille@Facile.Studio", true))
	if err != nil {
		t.Fatalf("first sign-in: %v", err)
	}
	if id != existing.ID {
		t.Fatalf("first sign-in created account %d, want the existing %d", id, existing.ID)
	}
	link(t, db, id, "sub-1")

	again, err := store.UpsertFromOIDC(ctx, claims("sub-1", "camille@facile.studio", true))
	if err != nil || again != existing.ID {
		t.Fatalf("second sign-in = %d, %v; want %d", again, err, existing.ID)
	}

	var user schemas.User
	if err := db.First(&user, id).Error; err != nil {
		t.Fatalf("reload: %v", err)
	}
	if user.Locale != "de" {
		t.Errorf("locale = %q, want the sign-in's de", user.Locale)
	}
}

func TestAnUnverifiedAddressNeverLinksToAnExistingAccount(t *testing.T) {
	db := openTestDatabase(t)
	store := newStore(t, db)

	if err := db.Create(&schemas.User{Email: "noah@facile.studio", Name: "Noah"}).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}
	if _, err := store.UpsertFromOIDC(context.Background(), claims("sub-2", "noah@facile.studio", false)); err == nil {
		t.Fatal("an unverified address linked to an existing account")
	}
}

func TestARecycledAddressIsNotHandedToANewSubject(t *testing.T) {
	db := openTestDatabase(t)
	store := newStore(t, db)

	owner := schemas.User{Email: "lea@facile.studio", Name: "Lea"}
	if err := db.Create(&owner).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}
	link(t, db, owner.ID, "sub-owner")

	if _, err := store.UpsertFromOIDC(context.Background(), claims("sub-newcomer", "lea@facile.studio", true)); err == nil {
		t.Fatal("a second subject of the same issuer took over a linked account")
	}
}
