package schemas

import "time"

type Email struct {
	ID             int64     `gorm:"column:id;primaryKey"`
	AccountID      int64     `gorm:"column:account_id;index"`
	FolderID       int64     `gorm:"column:folder_id;index"`
	MessageID      string    `gorm:"column:message_id;index"`
	ThreadID       string    `gorm:"column:thread_id;index"`
	Subject        string    `gorm:"column:subject"`
	FromAddress    string    `gorm:"column:from_address"`
	FromName       string    `gorm:"column:from_name"`
	ToAddresses    string    `gorm:"column:to_addresses;type:text"`
	CcAddresses    string    `gorm:"column:cc_addresses;type:text"`
	BccAddresses   string    `gorm:"column:bcc_addresses;type:text"`
	Date           time.Time `gorm:"column:date;index"`
	BodyText       string    `gorm:"column:body_text;type:text"`
	BodyHTML       string    `gorm:"column:body_html;type:text"`
	IsRead         bool      `gorm:"column:is_read;default:false"`
	IsStarred      bool      `gorm:"column:is_starred;default:false"`
	HasAttachments bool      `gorm:"column:has_attachments;default:false"`
	InReplyTo      string    `gorm:"column:in_reply_to"`
	References     string    `gorm:"column:references;type:text"`
	IMAPUID        uint32    `gorm:"column:imap_uid"`
	FacileID       *string   `gorm:"column:facile_id;uniqueIndex" json:"facile_id,omitempty"`
	CachedAt       time.Time `gorm:"column:cached_at;autoCreateTime"`
}

func (Email) TableName() string { return "emails" }
