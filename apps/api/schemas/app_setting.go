package schemas

// AppSetting is the single-row table holding app-wide settings.
type AppSetting struct {
	ID            int    `gorm:"primaryKey"`
	EncryptionKey string `gorm:"column:encryption_key;not null;default:''"`
}
