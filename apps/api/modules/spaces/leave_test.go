package spaces

import (
	"context"
	stderrors "errors"
	"log"
	"os"
	"testing"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/testdb"

	"gorm.io/gorm"
)

const bootstrap = "createdb courrier_test  # or point at any scratch database"

var testORM *gorm.DB

func TestMain(m *testing.M) {
	url, configured := testdb.URL()
	if !configured {
		testdb.Announce(bootstrap)
	} else {
		db, err := testdb.Open(url, testdb.Config{Prefix: "courrier_test", Migrate: schemas.Migrate})
		if err != nil {
			log.Fatalf("testdb: %v", err)
		}
		testORM = db
	}
	os.Exit(m.Run())
}

func newSpace(t *testing.T) (*Service, string, int64) {
	t.Helper()
	if testORM == nil {
		t.Skip(testdb.SkipReason(bootstrap))
	}
	if err := testORM.Exec(`DELETE FROM users`).Error; err != nil {
		t.Fatalf("clear users: %v", err)
	}
	if err := testORM.Exec(`DELETE FROM spaces`).Error; err != nil {
		t.Fatalf("clear spaces: %v", err)
	}

	owner := schemas.User{Email: "owner@facile.studio", Name: "Owner"}
	if err := testORM.Create(&owner).Error; err != nil {
		t.Fatalf("create the owner: %v", err)
	}
	service := NewService(testORM)
	space, _, err := service.Create(context.Background(), owner.ID, CreateSpaceRequest{Name: "Studio"})
	if err != nil {
		t.Fatalf("create the space: %v", err)
	}
	return service, space.ID, owner.ID
}

// A space with one owner has to keep them: leaving would strand it with no
// member who can delete it or appoint a replacement.
func TestTheSoleOwnerCannotLeave(t *testing.T) {
	service, spaceID, ownerID := newSpace(t)

	err := service.Leave(context.Background(), ownerID, spaceID)
	var envelope *errors.Error
	if !stderrors.As(err, &envelope) || envelope.Code != "already_exists" {
		t.Fatalf("expected a 409-class refusal, got %v", err)
	}
}

// Two people who own a space equally both have an exit. Refusing on the rank
// rather than the count made ownership transfer the only one.
func TestAnOwnerCanLeaveWhenAnotherOwnerRemains(t *testing.T) {
	service, spaceID, ownerID := newSpace(t)

	second := schemas.User{Email: "second@facile.studio", Name: "Second"}
	if err := testORM.Create(&second).Error; err != nil {
		t.Fatalf("create the second owner: %v", err)
	}
	member := schemas.SpaceMember{SpaceID: spaceID, UserID: second.ID, Role: schemas.SpaceRoleOwner}
	if err := testORM.Create(&member).Error; err != nil {
		t.Fatalf("add the second owner: %v", err)
	}

	if err := service.Leave(context.Background(), ownerID, spaceID); err != nil {
		t.Fatalf("an owner beside another owner could not leave: %v", err)
	}

	var remaining int64
	if err := testORM.Model(&schemas.SpaceMember{}).Where("space_id = ?", spaceID).Count(&remaining).Error; err != nil {
		t.Fatalf("count the members: %v", err)
	}
	if remaining != 1 {
		t.Fatalf("expected one member left, got %d", remaining)
	}
}
