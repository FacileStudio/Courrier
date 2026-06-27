package schemas

import "gorm.io/gorm"

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&User{},
		&Session{},
		&AppSetting{},
		&ApiToken{},
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

	return backfillAvatarURLs(db)
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
