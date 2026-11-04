package schemas

import "time"

// SenderRule pins every message from one address to a category in one
// account, overriding the header classifier and the model. Category uses the
// emails.category values: "" (primary), "bulk" (newsletters), "auto" (notifications).
type SenderRule struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	AccountID int64     `gorm:"column:account_id;not null;uniqueIndex:idx_sender_rules_account_address"`
	Account   *Account  `gorm:"constraint:OnDelete:CASCADE"`
	Address   string    `gorm:"column:address;not null;uniqueIndex:idx_sender_rules_account_address"`
	Category  string    `gorm:"column:category;not null;default:''"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (SenderRule) TableName() string { return "sender_rules" }
