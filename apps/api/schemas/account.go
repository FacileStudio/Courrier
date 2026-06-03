package schemas

import "time"

type Account struct {
	ID           int64     `gorm:"column:id;primaryKey"`
	UserID       int64     `gorm:"column:user_id;index"`
	Name         string    `gorm:"column:name"`
	Email        string    `gorm:"column:email"`
	IMAPHost     string    `gorm:"column:imap_host"`
	IMAPPort     int       `gorm:"column:imap_port;default:993"`
	IMAPUser     string    `gorm:"column:imap_user"`
	IMAPPassword string    `gorm:"column:imap_password"`
	SMTPHost     string    `gorm:"column:smtp_host"`
	SMTPPort     int       `gorm:"column:smtp_port;default:587"`
	SMTPUser     string    `gorm:"column:smtp_user"`
	SMTPPassword string    `gorm:"column:smtp_password"`
	Signature    string    `gorm:"column:signature;type:text"`
	IsDefault    bool      `gorm:"column:is_default;default:false"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Account) TableName() string { return "accounts" }
