package spaces

import (
	"context"
	stderrors "errors"
	"io"
	"log/slog"
	"maps"
	"testing"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/porte/directory"
	"github.com/FacileStudio/tronc/errors"
)

const testIssuer = "https://sso.test"

// directoryFixture is one Courrier space before Registre owned it: Alice owns
// it and signed in through Registre, Bob is a member who never did, Carol is a
// password-only account Registre has never heard of, and Dave has a Registre
// identity but is no longer in the space there.
func directoryFixture(t *testing.T) (*Directory, string) {
	t.Helper()
	_, spaceID, aliceID := newSpace(t)
	identity(t, aliceID, "alice")
	bob := user(t, "bob@facile.studio")
	carol := user(t, "carol@facile.studio")
	dave := user(t, "dave@facile.studio")
	identity(t, dave, "dave")
	for _, id := range []int64{bob, carol, dave} {
		member := schemas.SpaceMember{SpaceID: spaceID, UserID: id, Role: schemas.SpaceRoleMember}
		if err := testORM.Create(&member).Error; err != nil {
			t.Fatal(err)
		}
	}
	return NewDirectory(testORM, directory.Client{Issuer: testIssuer}, slog.New(slog.NewTextHandler(io.Discard, nil))), spaceID
}

func user(t *testing.T, email string) int64 {
	t.Helper()
	row := schemas.User{Email: email}
	if err := testORM.Create(&row).Error; err != nil {
		t.Fatal(err)
	}
	return row.ID
}

func identity(t *testing.T, userID int64, subject string) {
	t.Helper()
	if err := testORM.Exec(`INSERT INTO porte_identities (user_id, provider, subject) VALUES (?, ?, ?)`,
		userID, testIssuer, subject).Error; err != nil {
		t.Fatal(err)
	}
}

// snapshot is Registre's view: the space renamed, Alice still owner, Bob
// promoted, Erin new to Courrier, Dave gone.
func snapshot() directory.Snapshot {
	members := []directory.Member{
		{Subject: "alice", Email: "owner@facile.studio", Role: directory.RoleOwner},
		{Subject: "bob", Email: "Bob@facile.studio", Role: directory.RoleAdmin},
		{Subject: "erin", Email: "erin@facile.studio", Name: "Erin", Picture: testIssuer + "/avatars/erin"},
	}
	return directory.Snapshot{Version: 7, Spaces: []directory.Space{{ID: "r-studio", Name: "STUDIO", Members: members}}}
}

// assertRoles compares the space's membership, as email to role, with want.
func assertRoles(t *testing.T, spaceID string, want map[string]string) {
	t.Helper()
	var rows []struct{ Email, Role string }
	if err := testORM.Raw(`SELECT u.email, m.role FROM space_members m JOIN users u ON u.id = m.user_id
		WHERE m.space_id = ?`, spaceID).Scan(&rows).Error; err != nil {
		t.Fatal(err)
	}
	got := map[string]string{}
	for _, row := range rows {
		got[row.Email] = row.Role
	}
	if !maps.Equal(got, want) {
		t.Fatalf("members = %v, want %v", got, want)
	}
}

func TestApplyLinksProvisionsAndRemoves(t *testing.T) {
	dir, spaceID := directoryFixture(t)
	if err := dir.Apply(context.Background(), snapshot()); err != nil {
		t.Fatal(err)
	}

	var space schemas.Space
	testORM.First(&space, "id = ?", spaceID)
	if space.RegistreSpaceID == nil || *space.RegistreSpaceID != "r-studio" || space.Name != "STUDIO" {
		t.Fatalf("space was not linked by name: %+v", space)
	}
	assertRoles(t, spaceID, map[string]string{
		"owner@facile.studio": "owner", "bob@facile.studio": "admin",
		"erin@facile.studio": "member", "carol@facile.studio": "member",
	})
	var erin schemas.User
	testORM.First(&erin, "email = ?", "erin@facile.studio")
	if erin.OIDCPictureURL != testIssuer+"/avatars/erin" || erin.Name != "Erin" {
		t.Errorf("erin was not pre-provisioned with her profile: %+v", erin)
	}
	var linked int64
	testORM.Raw(`SELECT count(*) FROM porte_identities WHERE provider = ? AND subject IN ('bob', 'erin')`,
		testIssuer).Scan(&linked)
	if linked != 2 {
		t.Errorf("%d of bob and erin have a porte identity, want 2", linked)
	}
}

// The syncer reapplies on every boot, interval and nudge, so a snapshot that
// already holds must leave every row exactly as it was.
func TestApplyTwiceChangesNothing(t *testing.T) {
	dir, spaceID := directoryFixture(t)
	if err := dir.Apply(context.Background(), snapshot()); err != nil {
		t.Fatal(err)
	}
	before := state(t, spaceID)
	if err := dir.Apply(context.Background(), snapshot()); err != nil {
		t.Fatal(err)
	}
	if after := state(t, spaceID); after != before {
		t.Fatalf("second apply changed state:\n%s\n%s", before, after)
	}
}

// state is every row Apply can touch, flattened for comparison.
func state(t *testing.T, spaceID string) string {
	t.Helper()
	var out string
	if err := testORM.Raw(`SELECT
		(SELECT string_agg(id || name || updated_at::text, ',' ORDER BY id) FROM spaces) ||
		(SELECT string_agg(id || user_id::text || role || joined_at::text, ',' ORDER BY id)
		   FROM space_members WHERE space_id = ?) ||
		(SELECT count(*)::text FROM users) || (SELECT count(*)::text FROM porte_identities)`,
		spaceID).Scan(&out).Error; err != nil {
		t.Fatal(err)
	}
	return out
}

func TestALinkedSpaceRefusesLocalEdits(t *testing.T) {
	dir, spaceID := directoryFixture(t)
	if err := dir.Apply(context.Background(), snapshot()); err != nil {
		t.Fatal(err)
	}
	var owner schemas.User
	testORM.First(&owner, "email = ?", "owner@facile.studio")
	service := NewService(testORM).WithRegistre(testIssuer + "/")

	_, err := service.AddMember(context.Background(), owner.ID, spaceID, AddMemberRequest{Email: "carol@facile.studio"})
	var envelope *errors.Error
	if !stderrors.As(err, &envelope) || envelope.Code != "already_exists" {
		t.Fatalf("expected a 409 naming Registre, got %v", err)
	}
	if envelope.Message != "this space is managed in Registre: "+testIssuer+"/spaces/r-studio" {
		t.Errorf("message = %q", envelope.Message)
	}
	if _, _, err := service.Create(context.Background(), owner.ID, CreateSpaceRequest{Name: "New"}); err == nil {
		t.Error("creating a space succeeded while Registre owns spaces")
	}
}
