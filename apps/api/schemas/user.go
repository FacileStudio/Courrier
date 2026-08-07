package schemas

import "time"

type User struct {
	ID               int64     `gorm:"column:id;primaryKey"`
	Email            string    `gorm:"column:email;uniqueIndex"`
	Name             string    `gorm:"column:name"`
	AvatarURL        string    `gorm:"column:avatar_url"`
	AvatarSource     string    `gorm:"column:avatar_source"`
	AvatarUploadPath string    `gorm:"column:avatar_upload_path"`
	OIDCPictureURL   string    `gorm:"column:oidc_picture_url"`
	OIDCSubject      *string   `gorm:"column:oidc_subject;uniqueIndex"`
	OIDCAccessToken  string    `gorm:"column:oidc_access_token"`
	OIDCRefreshToken string    `gorm:"column:oidc_refresh_token"`
	OIDCTokenExpiry  time.Time `gorm:"column:oidc_token_expiry"`
	ProfileSyncedAt  time.Time `gorm:"column:profile_synced_at"`
	PasswordHash     string    `gorm:"column:password_hash"`
	CreatedAt        time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (User) TableName() string { return "users" }

// AvatarFilePrefix is where this app serves the storage volume. It is /api/files/ rather
// than the suite's usual /files/ because Courrier mounts its whole API under /api, and the
// client concatenates nothing — whatever Avatar returns is used as the src verbatim.
const AvatarFilePrefix = "/api/files/"

// AvatarSelectExpr is Avatar() as SQL, for the joins that read a user's picture without
// loading the row. It has to stay in step with Avatar below — hence both being here,
// one above the other, rather than one in Go and one buried in a Select string.
const AvatarSelectExpr = `COALESCE(NULLIF(users.oidc_picture_url, ''), ` +
	`NULLIF('` + AvatarFilePrefix + `' || COALESCE(users.avatar_upload_path, ''), '` + AvatarFilePrefix + `'), '')`

// Avatar is the picture to render. It is derived from the two sources rather than stored
// alongside them: a photo set in Porte always wins, an upload shows only when the IdP
// offers none, and because nothing is written there is no third value that can drift out
// of agreement with the two that matter.
func (u User) Avatar() string {
	if u.OIDCPictureURL != "" {
		return u.OIDCPictureURL
	}
	if u.AvatarUploadPath != "" {
		return AvatarFilePrefix + u.AvatarUploadPath
	}
	return ""
}

// AvatarOrigin names where Avatar came from, so the client can say *why* uploading is
// unavailable instead of just greying the button out.
func (u User) AvatarOrigin() string {
	switch {
	case u.OIDCPictureURL != "":
		return "oidc"
	case u.AvatarUploadPath != "":
		return "upload"
	default:
		return ""
	}
}
