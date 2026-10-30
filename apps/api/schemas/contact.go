package schemas

import "time"

// Contact is an address the account has written to, counted from its sent
// mail. Email is stored lowercased and is unique per account.
type Contact struct {
	ID         int64     `gorm:"column:id;primaryKey"`
	AccountID  int64     `gorm:"column:account_id;uniqueIndex:idx_contacts_account_email"`
	SpaceID    *string   `gorm:"column:space_id;index"`
	Email      string    `gorm:"column:email;uniqueIndex:idx_contacts_account_email"`
	Name       string    `gorm:"column:name"`
	SentCount  int       `gorm:"column:sent_count;default:0"`
	LastSentAt time.Time `gorm:"column:last_sent_at;index"`
}

func (Contact) TableName() string { return "contacts" }
