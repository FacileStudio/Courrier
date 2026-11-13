package auth

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/FacileStudio/porte/kit"
)

// columns maps porte/kit onto Courrier's users table. The SSO photo lives in
// oidc_picture_url; avatar_url is an unread legacy column. Courrier has no
// admin flag. created_at takes the database default 00002_column_defaults
// sets; where a rollback stripped it, kit.NewStore adopts the column instead.
func columns() kit.Columns {
	return kit.Columns{Picture: "oidc_picture_url"}
}

const fillNewAccountSQL = `UPDATE users SET
	avatar_url = '', avatar_source = '', avatar_upload_path = '',
	oidc_access_token = '', oidc_refresh_token = '', password_hash = '',
	oidc_token_expiry = '0001-01-01 00:00:00+00', profile_synced_at = '0001-01-01 00:00:00+00'
	WHERE id = $1`

// fillNewAccount writes the values GORM wrote on every account it created, so
// an account the kit creates reads exactly like an older one. oidc_subject
// stays NULL: it is unique, so an empty string would collide on the second
// account. 00002_column_defaults writes the same values as defaults; this
// keeps them on a database where a rollback to the AutoMigrate release
// stripped those defaults.
func fillNewAccount(ctx context.Context, tx *sql.Tx, userID int64) error {
	_, err := tx.ExecContext(ctx, fillNewAccountSQL, userID)
	return err
}

// NewUserStore builds porte/kit's user store over Courrier's users table: the
// sign-in upsert, the identity lookup behind RequireAuth, and the Registre
// erasure with eraseData for the rows a person owns.
func NewUserStore(ctx context.Context, db *sql.DB, logger *slog.Logger) (*kit.Store, error) {
	return kit.NewStore(ctx, kit.Store{
		DB:        db,
		Columns:   columns(),
		OnCreate:  fillNewAccount,
		EraseData: eraseData,
		Logger:    logger,
	})
}
