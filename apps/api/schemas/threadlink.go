package schemas

// ThreadLink maps every Message-ID ever referenced by an account's mail —
// including parents that have not arrived yet — to the thread it belongs to.
// It is the union-find substrate that lets threading survive out-of-order
// delivery and References truncation: a late parent finds its children, and a
// reply that bridges two partial chains merges them.
type ThreadLink struct {
	AccountID int64  `gorm:"column:account_id;primaryKey"`
	MessageID string `gorm:"column:message_id;primaryKey"`
	ThreadID  string `gorm:"column:thread_id;index"`
}

func (ThreadLink) TableName() string { return "thread_links" }
