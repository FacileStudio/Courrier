package mail

type SendRequest struct {
	To      []string `json:"to"`
	Cc      []string `json:"cc"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

type TestConnectionRequest struct {
	IMAPHost     string `json:"imap_host"`
	IMAPPort     int    `json:"imap_port"`
	IMAPUser     string `json:"imap_user"`
	IMAPPassword string `json:"imap_password"`
	SMTPHost     string `json:"smtp_host"`
	SMTPPort     int    `json:"smtp_port"`
	SMTPUser     string `json:"smtp_user"`
	SMTPPassword string `json:"smtp_password"`
}

type FolderResponse struct {
	ID          int64  `json:"id"`
	AccountID   int64  `json:"account_id"`
	Path        string `json:"path"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	UnreadCount int    `json:"unread_count"`
	TotalCount  int    `json:"total_count"`
}

type EmailResponse struct {
	ID             int64             `json:"id"`
	AccountID      int64             `json:"account_id"`
	FolderID       int64             `json:"folder_id"`
	MessageID      string            `json:"message_id"`
	Subject        string            `json:"subject"`
	FromAddress    string            `json:"from_address"`
	FromName       string            `json:"from_name"`
	ToAddresses    []AddressResponse `json:"to_addresses"`
	CcAddresses    []AddressResponse `json:"cc_addresses"`
	Date           string            `json:"date"`
	BodyText       string            `json:"body_text,omitempty"`
	BodyHTML       string            `json:"body_html,omitempty"`
	IsRead         bool              `json:"is_read"`
	IsStarred      bool              `json:"is_starred"`
	HasAttachments bool              `json:"has_attachments"`
}

type AddressResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateEmailRequest struct {
	IsRead    *bool `json:"is_read"`
	IsStarred *bool `json:"is_starred"`
}
