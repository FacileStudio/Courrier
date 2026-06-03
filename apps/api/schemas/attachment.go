package schemas

type Attachment struct {
	ID       int64  `gorm:"column:id;primaryKey"`
	EmailID  int64  `gorm:"column:email_id;index"`
	Filename string `gorm:"column:filename"`
	MimeType string `gorm:"column:mime_type"`
	Size     int64  `gorm:"column:size"`
	CID      string `gorm:"column:cid"`
	PartID   string `gorm:"column:part_id"`
}

func (Attachment) TableName() string { return "attachments" }
