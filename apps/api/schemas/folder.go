package schemas

type Folder struct {
	ID          int64  `gorm:"column:id;primaryKey"`
	AccountID   int64  `gorm:"column:account_id;index"`
	Path        string `gorm:"column:path"`
	Name        string `gorm:"column:name"`
	Type        string `gorm:"column:type;default:'custom'"`
	UnreadCount int    `gorm:"column:unread_count;default:0"`
	TotalCount  int    `gorm:"column:total_count;default:0"`
	UIDValidity uint32 `gorm:"column:uid_validity;default:0"`
}

func (Folder) TableName() string { return "folders" }

const (
	FolderTypeInbox   = "inbox"
	FolderTypeSent    = "sent"
	FolderTypeDrafts  = "drafts"
	FolderTypeTrash   = "trash"
	FolderTypeJunk    = "junk"
	FolderTypeArchive = "archive"
	FolderTypeCustom  = "custom"
)
