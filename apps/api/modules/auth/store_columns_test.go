package auth

import (
	"context"
	"encoding/json"
	"slices"
	"testing"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
)

// The kit's insert leaves these users columns out. Every one but oidc_subject
// is filled by fillNewAccount; a new nullable column fails this test instead
// of shipping NULL.
func TestUnwrittenColumnsArePinned(t *testing.T) {
	store := newStore(t, openTestDatabase(t))
	want := []string{
		"avatar_url", "avatar_source", "avatar_upload_path", "oidc_subject", "oidc_access_token",
		"oidc_refresh_token", "oidc_token_expiry", "profile_synced_at", "password_hash",
	}
	got := store.Unwritten()
	slices.Sort(got)
	slices.Sort(want)
	if !slices.Equal(got, want) {
		t.Fatalf("Unwritten() = %v, want %v", got, want)
	}
}

func nullColumns(t *testing.T, userID int64) []string {
	t.Helper()
	var names []string
	err := testORM.Raw(`SELECT key FROM users u, jsonb_each(to_jsonb(u)) WHERE u.id = ? AND value = 'null'::jsonb ORDER BY key`, userID).
		Scan(&names).Error
	if err != nil {
		t.Fatalf("null columns: %v", err)
	}
	return names
}

func row(t *testing.T, userID int64) map[string]any {
	t.Helper()
	var raw string
	if err := testORM.Raw(`SELECT to_jsonb(u)::text FROM users u WHERE u.id = ?`, userID).Scan(&raw).Error; err != nil {
		t.Fatalf("row: %v", err)
	}
	out := map[string]any{}
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("row: %v", err)
	}
	return out
}

// An account the kit creates at sign-in has NULL in exactly the columns an
// account GORM created has NULL in: only oidc_subject.
func TestANewAccountHasNoNullAnOlderAccountFilled(t *testing.T) {
	db := openTestDatabase(t)
	store := newStore(t, db)

	older := schemas.User{Email: "older@facile.studio", Name: "Older"}
	if err := db.Create(&older).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}
	id, err := store.UpsertFromOIDC(context.Background(), claims("sub-new", "new@facile.studio", true))
	if err != nil {
		t.Fatalf("sign-in: %v", err)
	}

	got, want := nullColumns(t, id), nullColumns(t, older.ID)
	if !slices.Equal(got, want) || !slices.Equal(got, []string{"oidc_subject"}) {
		t.Fatalf("new account NULL columns = %v, older account = %v, want [oidc_subject]", got, want)
	}
	gormRow, kitRow := row(t, older.ID), row(t, id)
	for _, column := range []string{
		"avatar_url", "avatar_source", "avatar_upload_path", "oidc_access_token",
		"oidc_refresh_token", "oidc_token_expiry", "profile_synced_at", "password_hash",
	} {
		if gormRow[column] != kitRow[column] {
			t.Errorf("%s = %v, GORM wrote %v", column, kitRow[column], gormRow[column])
		}
	}
	var user schemas.User
	if err := db.First(&user, id).Error; err != nil {
		t.Fatalf("reload: %v", err)
	}
	if user.CreatedAt.IsZero() || user.Locale != "de" || user.Email != "new@facile.studio" {
		t.Fatalf("new account = %+v", user)
	}
}
