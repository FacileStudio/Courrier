package schemas

import (
	"context"
	"log/slog"

	"github.com/FacileStudio/Courrier/apps/api/migrations"
	"gorm.io/gorm"
)

// Migrate applies the migrations and then the boot data steps, with no OIDC
// issuer. The tests build their schemas through it, so they run exactly the
// SQL production runs.
func Migrate(db *gorm.DB) error {
	return MigrateWithIssuer(context.Background(), db, "", slog.Default())
}

// MigrateWithIssuer applies the ordered migrations under migrations/ and then
// Reconcile. The issuer is what the identity backfill needs: porte matches an
// account on (provider, subject) and the provider is the issuer, so
// backfilling with a placeholder would leave every existing SSO user
// unmatched and quietly fall through to the email path.
func MigrateWithIssuer(ctx context.Context, db *gorm.DB, issuer string, logger *slog.Logger) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	if err := migrations.Up(ctx, migrations.Config(sqlDB, logger)); err != nil {
		return err
	}
	return Reconcile(db, issuer)
}

// Reconcile runs the idempotent data steps that ran at every boot before the
// schema moved to migrations, in the same order. It changes rows, never the
// schema: the email class columns, the avatar columns, and the credentials
// porte reads. AdoptPorte runs last because it repoints identities at users.
func Reconcile(db *gorm.DB, issuer string) error {
	if err := backfillEmailClassColumns(db); err != nil {
		return err
	}
	if err := backfillAvatarURLs(db); err != nil {
		return err
	}
	if err := backfillAvatarSources(db); err != nil {
		return err
	}
	return AdoptPorte(db, issuer)
}

// backfillEmailClassColumns turns the NULL category and bulk_reason that
// AutoMigrate leaves on every row cached before those columns existed into the
// empty string the classifier and the backfill select on. Idempotent.
func backfillEmailClassColumns(db *gorm.DB) error {
	return db.Exec(
		`UPDATE emails
		    SET category = coalesce(category, ''),
		        bulk_reason = coalesce(bulk_reason, '')
		  WHERE category IS NULL OR bulk_reason IS NULL`).Error
}

// backfillAvatarURLs repoints avatars stored before static files moved under
// /api/files/. Rows still holding the legacy /files/ prefix 404 against the
// current route, so normalize them to the served path. Idempotent.
func backfillAvatarURLs(db *gorm.DB) error {
	return db.Model(&User{}).
		Where("avatar_url LIKE ?", "/files/%").
		Update("avatar_url", gorm.Expr("'/api/files/' || substring(avatar_url from 8)")).
		Error
}

// backfillAvatarSources moves the two real avatar sources onto the columns that now own
// them: the uploaded file onto avatar_upload_path, and a Porte photo — and only a Porte
// photo — onto oidc_picture_url. The rendered value is derived from those two by
// User.Avatar; nothing stores a third copy that can drift.
//
// The filename decides which rows are uploads, not avatar_source. That column was added
// after the upload feature, so the oldest uploaded avatars have it empty, and keying on
// avatar_source = 'upload' quietly drops their picture. persistAvatarFile has always named
// uploads "user-<id>-<nanos>" and the old OIDC download named its copies "oidc-<id>-<nanos>",
// so anything that is not an oidc- copy is somebody's upload and is kept. Both the current
// /api/files/ prefix and the legacy /files/ one are accepted, so this is correct whether or
// not backfillAvatarURLs has already run.
//
// The oidc- copies are the ones with nothing to preserve: oidc_picture_url already holds the
// URL that replaces them. They are left on the volume rather than deleted here — a migration
// that removes files has to be right the first time.
//
// The second statement is the one that is easy to miss. The old sync stored profile.Picture
// verbatim, so every user without a photo in Authentik has a data:image/svg+xml blob sitting
// in oidc_picture_url. Under the new rule that column means "there is an SSO photo", so a
// stale blob would suppress the upload fallback and render Authentik's initials instead of
// ours. Blanking it is what makes the column mean what Avatar assumes it means.
//
// avatar_url and avatar_source stay in the table, unread, until the next release drops them.
// Expanding first means a rollback is redeploying the old binary, not restoring a backup.
//
// lower() is used so the second statement agrees with kit.PictureURL,
// which compares the picture scheme case-insensitively. The trailing coalesce
// turns NULLs into empty strings, because a NULL here would fail to scan into
// the plain string the model declares.
func backfillAvatarSources(db *gorm.DB) error {
	if db.Migrator().HasColumn(&User{}, "avatar_url") {
		if err := db.Exec(
			`UPDATE users
			    SET avatar_upload_path = regexp_replace(avatar_url, '^/(api/)?files/', '')
			  WHERE coalesce(avatar_url, '') <> ''
			    AND avatar_url !~ '^/(api/)?files/avatars/oidc-'
			    AND coalesce(avatar_upload_path, '') = ''`).Error; err != nil {
			return err
		}
	}

	if err := db.Exec(
		`UPDATE users
		    SET oidc_picture_url = ''
		  WHERE coalesce(oidc_picture_url, '') <> ''
		    AND lower(oidc_picture_url) NOT LIKE 'https://%'`).Error; err != nil {
		return err
	}

	return db.Exec(
		`UPDATE users
		    SET avatar_upload_path = coalesce(avatar_upload_path, ''),
		        oidc_picture_url = coalesce(oidc_picture_url, '')
		  WHERE avatar_upload_path IS NULL OR oidc_picture_url IS NULL`).Error
}
