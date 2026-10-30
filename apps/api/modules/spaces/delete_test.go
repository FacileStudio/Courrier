package spaces

import (
	"context"
	"testing"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
)

func TestDeleteRemovesSpaceContacts(t *testing.T) {
	service, spaceID, ownerID := newSpace(t)
	testORM.Exec("DELETE FROM contacts")
	acc := schemas.Account{UserID: ownerID, Name: "A", Email: "a@example.test", SpaceID: &spaceID}
	if err := testORM.Create(&acc).Error; err != nil {
		t.Fatal(err)
	}
	contact := schemas.Contact{AccountID: acc.ID, SpaceID: &spaceID, Email: "x@x.test", LastSentAt: time.Now()}
	if err := testORM.Create(&contact).Error; err != nil {
		t.Fatal(err)
	}
	if err := service.Delete(context.Background(), ownerID, spaceID); err != nil {
		t.Fatal(err)
	}
	var n int64
	testORM.Model(&schemas.Contact{}).Where("space_id = ?", spaceID).Count(&n)
	if n != 0 {
		t.Fatalf("%d space contacts survived the space", n)
	}
}
