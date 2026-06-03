package schemas

type AppSetting struct {
	ID              int    `gorm:"primaryKey"`
	EncryptionKey   string `gorm:"column:encryption_key;not null;default:''"`
	NookPoolURL     string `gorm:"column:nook_pool_url;not null;default:''" json:"nook_pool_url"`
	NookPoolSecret  string `gorm:"column:nook_pool_secret;not null;default:''" json:"nook_pool_secret"`
	NookPoolEnabled bool   `gorm:"column:nook_pool_enabled;not null;default:false" json:"nook_pool_enabled"`
}
