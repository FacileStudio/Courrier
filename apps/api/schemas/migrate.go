package schemas

import "gorm.io/gorm"

// Migrate brings the schema up to date and then hands authentication to
// porte. AdoptPorte runs last because it repoints foreign keys at users(id)
// and reads the columns AutoMigrate has just guaranteed exist.
func Migrate(db *gorm.DB) error {
	return MigrateWithIssuer(db, "")
}

// MigrateWithIssuer is Migrate with the OIDC issuer, which the identity
// backfill needs: porte matches an account on (provider, subject) and the
// provider is the issuer, so backfilling with a placeholder would leave every
// existing SSO user unmatched and quietly fall through to the email path.
func MigrateWithIssuer(db *gorm.DB, issuer string) error {
	if err := db.AutoMigrate(
		&User{},
		&AppSetting{},
		&Account{},
		&Folder{},
		&Email{},
		&ThreadLink{},
		&Attachment{},
		&EmailTemplate{},
		&Space{},
		&SpaceMember{},
	); err != nil {
		return err
	}

	ensureSearchIndexes(db)

	if err := backfillAvatarURLs(db); err != nil {
		return err
	}

	if err := backfillAvatarSources(db); err != nil {
		return err
	}

	return AdoptPorte(db, issuer)
}

func ensureSearchIndexes(db *gorm.DB) {
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS pg_trgm").Error; err != nil {
		return
	}

	statements := []string{
		"CREATE INDEX IF NOT EXISTS idx_emails_subject_trgm ON emails USING gin (subject gin_trgm_ops)",
		"CREATE INDEX IF NOT EXISTS idx_emails_from_name_trgm ON emails USING gin (from_name gin_trgm_ops)",
		"CREATE INDEX IF NOT EXISTS idx_emails_from_address_trgm ON emails USING gin (from_address gin_trgm_ops)",
		"CREATE INDEX IF NOT EXISTS idx_emails_body_text_trgm ON emails USING gin (body_text gin_trgm_ops)",
		"CREATE INDEX IF NOT EXISTS idx_emails_account_thread ON emails (account_id, thread_id)",
		"CREATE INDEX IF NOT EXISTS idx_emails_folder_thread_date ON emails (account_id, folder_id, thread_id, date DESC, id DESC)",
		"CREATE INDEX IF NOT EXISTS idx_thread_links_account_thread ON thread_links (account_id, thread_id)",
	}

	for _, stmt := range statements {
		db.Exec(stmt)
	}
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

	// lower() so this agrees with oidcavatar.PhotoURL, which compares the scheme case-insensitively.
	if err := db.Exec(
		`UPDATE users
		    SET oidc_picture_url = ''
		  WHERE coalesce(oidc_picture_url, '') <> ''
		    AND lower(oidc_picture_url) NOT LIKE 'https://%'`).Error; err != nil {
		return err
	}

	// A NULL here would fail to scan into the plain string the model declares.
	return db.Exec(
		`UPDATE users
		    SET avatar_upload_path = coalesce(avatar_upload_path, ''),
		        oidc_picture_url = coalesce(oidc_picture_url, '')
		  WHERE avatar_upload_path IS NULL OR oidc_picture_url IS NULL`).Error
}
