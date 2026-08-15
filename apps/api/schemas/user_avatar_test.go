package schemas

import "testing"

func TestAvatarPrecedence(t *testing.T) {
	const porte = "https://porte.facile.studio/media/user-avatars/x.png"

	cases := []struct {
		name       string
		user       User
		wantURL    string
		wantOrigin string
	}{
		{"Porte photo wins over an upload", User{OIDCPictureURL: porte, AvatarUploadPath: "avatars/user-3-1.png"}, porte, "oidc"},
		{"upload is the fallback", User{AvatarUploadPath: "avatars/user-3-1.png"}, "/api/files/avatars/user-3-1.png", "upload"},
		{"only Porte", User{OIDCPictureURL: porte}, porte, "oidc"},
		{"neither, so the client draws initials", User{}, "", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.user.Avatar(); got != tc.wantURL {
				t.Errorf("Avatar() = %q, want %q", got, tc.wantURL)
			}
			if got := tc.user.AvatarOrigin(); got != tc.wantOrigin {
				t.Errorf("AvatarOrigin() = %q, want %q", got, tc.wantOrigin)
			}
		})
	}
}

// AvatarSelectExpr exists for the joins that read a user's picture without loading the row,
// so the two spellings of the same rule have to agree. This is the test that fails when
// someone edits one and forgets the other.
func TestAvatarSelectExprMatchesAvatar(t *testing.T) {
	orm := openTestDatabase(t)

	users := []User{
		{Email: "both@example.com", OIDCPictureURL: "https://porte.facile.studio/media/user-avatars/a.png", AvatarUploadPath: "avatars/user-1-1.png"},
		{Email: "upload@example.com", AvatarUploadPath: "avatars/user-2-1.png"},
		{Email: "oidc@example.com", OIDCPictureURL: "https://porte.facile.studio/media/user-avatars/b.png"},
		{Email: "neither@example.com"},
	}
	for i := range users {
		if err := orm.Create(&users[i]).Error; err != nil {
			t.Fatalf("create %s: %v", users[i].Email, err)
		}
	}

	for _, want := range users {
		var got string
		if err := orm.Model(&User{}).
			Select(AvatarSelectExpr).
			Where("users.id = ?", want.ID).
			Scan(&got).Error; err != nil {
			t.Fatalf("select for %s: %v", want.Email, err)
		}
		if got != want.Avatar() {
			t.Errorf("%s: SQL gave %q, Avatar() gave %q", want.Email, got, want.Avatar())
		}
	}
}

// Rows 2 and 5 are why this test exists. Row 2 holds an uploaded avatar with avatar_source
// empty, because it predates that column, and a backfill keyed on avatar_source = 'upload'
// drops its picture without a word. Row 5 holds the data: blob the old sync stored verbatim
// for a user Authentik has no photo for — left in place it would read as "has an SSO photo"
// and suppress the upload fallback forever.
//
// The row that carries both keeps its file and still renders the Porte photo;
// the initials-blob row is now free to fall back to an upload.
func TestBackfillAvatarSources(t *testing.T) {
	orm := openTestDatabase(t)

	rows := []struct {
		email      string
		url        string
		source     string
		oidc       string
		wantUpload string
		wantOIDC   string
	}{
		{"oidc-copy@example.com", "/api/files/avatars/oidc-1-178006.png", "oidc", "https://porte.facile.studio/media/user-avatars/a.png", "", "https://porte.facile.studio/media/user-avatars/a.png"},
		{"legacy-upload@example.com", "/api/files/avatars/user-2-177802.jpg", "", "", "avatars/user-2-177802.jpg", ""},
		{"upload-and-sso@example.com", "/api/files/avatars/user-3-178096.jpg", "upload", "https://porte.facile.studio/media/user-avatars/b.jpeg", "avatars/user-3-178096.jpg", "https://porte.facile.studio/media/user-avatars/b.jpeg"},
		{"no-avatar@example.com", "", "", "", "", ""},
		{"initials-blob@example.com", "", "", "data:image/svg+xml;base64,PHN2Zz48L3N2Zz4=", "", ""},
		{"pre-prefix-upload@example.com", "/files/avatars/user-6-177000.png", "", "", "avatars/user-6-177000.png", ""},
	}
	for _, row := range rows {
		if err := orm.Exec(
			`INSERT INTO users (email, password_hash, avatar_url, avatar_source, oidc_picture_url) VALUES (?, 'hash', ?, ?, ?)`,
			row.email, row.url, row.source, row.oidc).Error; err != nil {
			t.Fatalf("insert %s: %v", row.email, err)
		}
	}

	if err := backfillAvatarSources(orm); err != nil {
		t.Fatalf("backfill: %v", err)
	}

	for _, row := range rows {
		var got User
		if err := orm.Where("email = ?", row.email).First(&got).Error; err != nil {
			t.Fatalf("read %s: %v", row.email, err)
		}
		if got.AvatarUploadPath != row.wantUpload {
			t.Errorf("%s: avatar_upload_path = %q, want %q", row.email, got.AvatarUploadPath, row.wantUpload)
		}
		if got.OIDCPictureURL != row.wantOIDC {
			t.Errorf("%s: oidc_picture_url = %q, want %q", row.email, got.OIDCPictureURL, row.wantOIDC)
		}
	}

	var both User
	if err := orm.Where("email = ?", "upload-and-sso@example.com").First(&both).Error; err != nil {
		t.Fatalf("read both: %v", err)
	}
	if both.Avatar() != "https://porte.facile.studio/media/user-avatars/b.jpeg" {
		t.Errorf("SSO photo should win, got %q", both.Avatar())
	}

	var blob User
	if err := orm.Where("email = ?", "initials-blob@example.com").First(&blob).Error; err != nil {
		t.Fatalf("read blob: %v", err)
	}
	if blob.AvatarOrigin() != "" {
		t.Errorf("stale data: URI still counts as an SSO photo, origin = %q", blob.AvatarOrigin())
	}
}

// The upload path is the only file this app still owns, so the migration must not create a
// row that points at somebody else's file.
func TestBackfillIsIdempotent(t *testing.T) {
	orm := openTestDatabase(t)

	if err := orm.Exec(
		`INSERT INTO users (email, password_hash, avatar_url, avatar_source, oidc_picture_url) VALUES (?, 'hash', ?, '', '')`,
		"twice@example.com", "/api/files/avatars/user-9-1.png").Error; err != nil {
		t.Fatalf("insert: %v", err)
	}
	for i := 0; i < 2; i++ {
		if err := backfillAvatarSources(orm); err != nil {
			t.Fatalf("backfill run %d: %v", i, err)
		}
	}

	var got User
	if err := orm.Where("email = ?", "twice@example.com").First(&got).Error; err != nil {
		t.Fatalf("read: %v", err)
	}
	if got.AvatarUploadPath != "avatars/user-9-1.png" {
		t.Errorf("avatar_upload_path = %q after two runs", got.AvatarUploadPath)
	}
}
