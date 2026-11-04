package mail

import (
	"context"
	stderrors "errors"
	"strings"
	"testing"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/testdb"
)

// TestTemplateSpaceIsCheckedAndCanonical stores a template's space under its
// canonical id, and refuses a space the caller is not in or a spelling that
// is not an id.
func TestTemplateSpaceIsCheckedAndCanonical(t *testing.T) {
	if testORM == nil {
		t.Skip(testdb.SkipReason("createdb courrier_test  # or point at any scratch database"))
	}
	testORM.Exec("DELETE FROM users WHERE email LIKE 'tmpl-%'")
	member := schemas.User{Email: "tmpl-member-" + t.Name() + "@facile.studio"}
	outsider := schemas.User{Email: "tmpl-outsider-" + t.Name() + "@facile.studio"}
	space := schemas.Space{Name: "Templates"}
	if testORM.Create(&member).Error != nil || testORM.Create(&outsider).Error != nil || testORM.Create(&space).Error != nil ||
		testORM.Create(&schemas.SpaceMember{SpaceID: space.ID, UserID: member.ID}).Error != nil {
		t.Fatal("fixture rows")
	}
	service := NewService(testORM, nil)
	ctx := context.Background()
	upper := strings.ToUpper(space.ID)
	tmpl, err := service.CreateTemplate(ctx, member.ID, EmailTemplateRequest{Name: "t", SpaceID: &upper})
	if err != nil || tmpl.SpaceID == nil || *tmpl.SpaceID != space.ID {
		t.Fatalf("member template: %v, stored %v", err, tmpl.SpaceID)
	}
	expectCode(t, "outsider", "permission_denied", func() error {
		_, err := service.CreateTemplate(ctx, outsider.ID, EmailTemplateRequest{Name: "t", SpaceID: &space.ID})
		return err
	})
	braced := "{" + space.ID + "}"
	expectCode(t, "braced", "invalid_argument", func() error {
		_, err := service.ListTemplates(ctx, member.ID, &braced)
		return err
	})
}

// expectCode fails unless call answers an envelope error with code.
func expectCode(t *testing.T, what, code string, call func() error) {
	t.Helper()
	var appErr *errors.Error
	if err := call(); !stderrors.As(err, &appErr) || appErr.Code != code {
		t.Errorf("%s: %v, want %s", what, err, code)
	}
}
