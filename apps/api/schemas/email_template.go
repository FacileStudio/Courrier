package schemas

import "time"

type EmailTemplate struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	UserID    int64     `gorm:"column:user_id;index"`
	Name      string    `gorm:"column:name"`
	Subject   string    `gorm:"column:subject"`
	BodyHTML  string    `gorm:"column:body_html;type:text"`
	BodyText  string    `gorm:"column:body_text;type:text"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (EmailTemplate) TableName() string { return "email_templates" }
