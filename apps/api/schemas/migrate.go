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
		&Attachment{},
		&EmailTemplate{},
		&Space{},
		&SpaceMember{},
	); err != nil {
		return err
	}

	return backfillAvatarURLs(db)
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
