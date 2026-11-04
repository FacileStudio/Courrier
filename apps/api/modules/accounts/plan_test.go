package accounts

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"io"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FacileStudio/Courrier/apps/api/internal/plans"
	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/porte/directory"
	"github.com/FacileStudio/porte/forfait"
	"github.com/FacileStudio/tronc/httpjson"
	"github.com/FacileStudio/tronc/testdb"
)

const planIssuer = "https://sso.test"

// planFixture is a user who is a member of a Registre space whose plan allows
// one mailbox, with the gate already holding that plan.
func planFixture(t *testing.T) (*Service, int64, string) {
	t.Helper()
	if testORM == nil {
		t.Skip(testdb.SkipReason("createdb courrier_test  # or point at any scratch database"))
	}
	for _, table := range []string{"forfait_state", "contacts", "accounts", "space_members", "spaces"} {
		testORM.Exec("DELETE FROM " + table)
	}
	testORM.Exec("DELETE FROM users WHERE email LIKE 'plan-%'")
	user := schemas.User{Email: "plan-" + t.Name() + "@facile.studio"}
	link := "spc_" + strings.ToLower(t.Name())
	space := schemas.Space{Name: t.Name(), RegistreSpaceID: &link}
	if testORM.Create(&user).Error != nil || testORM.Create(&space).Error != nil ||
		testORM.Create(&schemas.SpaceMember{SpaceID: space.ID, UserID: user.ID, Role: schemas.SpaceRoleOwner}).Error != nil {
		t.Fatal("fixture rows")
	}
	gate := plans.New(testORM, planIssuer, slog.New(slog.NewTextHandler(io.Discard, nil)))
	snap := directory.Snapshot{Version: 1, EntitlementsRevision: 1, Spaces: []directory.Space{
		{ID: link, Entitlements: map[string]json.RawMessage{plans.Mailboxes: json.RawMessage("1")}},
	}}
	if err := gate.Apply(context.Background(), snap); err != nil {
		t.Fatal(err)
	}
	return NewService(testORM, nil, gate), user.ID, space.ID
}

func mailbox(spaceID *string) CreateAccountRequest {
	return CreateAccountRequest{Name: "Box", Email: "box@facile.studio", IMAPHost: "imap.test", SMTPHost: "smtp.test", SpaceID: spaceID}
}

// TestCreateRefusedOverLimit adds up to the ceiling, then refuses without
// writing anything.
func TestCreateRefusedOverLimit(t *testing.T) {
	service, userID, spaceID := planFixture(t)
	ctx := context.Background()
	if _, err := service.Create(ctx, userID, mailbox(&spaceID)); err != nil {
		t.Fatalf("first mailbox: %v", err)
	}
	_, err := service.Create(ctx, userID, mailbox(&spaceID))
	if !stderrors.Is(err, forfait.ErrLimit) {
		t.Fatalf("second mailbox: %v, want ErrLimit", err)
	}
	var count int64
	testORM.Model(&schemas.Account{}).Where("space_id = ?", spaceID).Count(&count)
	if count != 1 {
		t.Fatalf("space holds %d accounts after the refusal, want 1", count)
	}
}

// TestOverLimitAccountsKeepWorking reads, edits and deletes accounts of a
// scope already past its ceiling.
func TestOverLimitAccountsKeepWorking(t *testing.T) {
	service, userID, spaceID := planFixture(t)
	ctx := context.Background()
	held := []schemas.Account{{UserID: userID, SpaceID: &spaceID, Name: "a"}, {UserID: userID, SpaceID: &spaceID, Name: "b"}}
	if err := testORM.Create(&held).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := service.List(ctx, userID, &spaceID); err != nil {
		t.Fatalf("list: %v", err)
	}
	renamed := "renamed"
	if _, err := service.Update(ctx, userID, held[0].ID, UpdateAccountRequest{Name: &renamed}); err != nil {
		t.Fatalf("update: %v", err)
	}
	if err := service.Delete(ctx, userID, held[1].ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

// TestMoveIntoSpaceRefused pins that an update cannot carry a space: the
// strict decoder rejects the field, so no account changes scope past the
// plan check.
func TestMoveIntoSpaceRefused(t *testing.T) {
	request := httptest.NewRequest("PUT", "/accounts/1", strings.NewReader(`{"name":"x","space_id":"spc"}`))
	var body UpdateAccountRequest
	if err := httpjson.DecodeJSON(httptest.NewRecorder(), request, &body); err == nil {
		t.Fatal("an update carrying space_id was accepted")
	}
}

// TestConcurrentCreatesAtTheCeiling lets exactly one of several racing
// creations take the last mailbox, because the check locks the space row.
func TestConcurrentCreatesAtTheCeiling(t *testing.T) {
	service, userID, spaceID := planFixture(t)
	results := make(chan error, 6)
	for range cap(results) {
		go func() {
			_, err := service.Create(context.Background(), userID, mailbox(&spaceID))
			results <- err
		}()
	}
	created := 0
	for range cap(results) {
		if err := <-results; err == nil {
			created++
		} else if !stderrors.Is(err, forfait.ErrLimit) {
			t.Errorf("racing create: %v", err)
		}
	}
	if created != 1 {
		t.Fatalf("%d racing creations passed a ceiling of 1", created)
	}
}
