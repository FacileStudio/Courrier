package mail

import (
	"bytes"
	"io"
	"os"
)

// AttachmentUpload is a file being sent with an email, either held in memory
// or staged on disk at FilePath.
type AttachmentUpload struct {
	Filename string
	MimeType string
	Data     []byte
	FilePath string
}

func (a *AttachmentUpload) Reader() (io.ReadCloser, error) {
	if a.Data != nil {
		return io.NopCloser(bytes.NewReader(a.Data)), nil
	}
	return os.Open(a.FilePath)
}

func (a *AttachmentUpload) Cleanup() {
	if a.FilePath != "" {
		os.Remove(a.FilePath)
	}
}

// SendRequest is the body of POST /accounts/{id}/mail/send.
type SendRequest struct {
	To          []string           `json:"to"`
	Cc          []string           `json:"cc"`
	Subject     string             `json:"subject"`
	Body        string             `json:"body"`
	BodyHTML    string             `json:"body_html"`
	InReplyTo   string             `json:"in_reply_to"`
	References  []string           `json:"references"`
	Attachments []AttachmentUpload `json:"-"`
}

// TestConnectionRequest is the body of POST /api/mail/test-connection.
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

// CheckLeg is one half of a connection check. `Configured` false means the account has no
// host for that protocol, which is a different answer from "it did not work".
type CheckLeg struct {
	Configured bool   `json:"configured"`
	OK         bool   `json:"ok"`
	Error      string `json:"error,omitempty"`
}

// CheckResult is the outcome of a connection check, one leg per protocol.
type CheckResult struct {
	IMAP CheckLeg `json:"imap"`
	SMTP CheckLeg `json:"smtp"`
}

// FolderResponse is a mail folder as returned to the client.
type FolderResponse struct {
	ID          int64  `json:"id"`
	AccountID   int64  `json:"account_id"`
	Path        string `json:"path"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	UnreadCount int    `json:"unread_count"`
	TotalCount  int    `json:"total_count"`
}

// EmailResponse is an email (or an expanded conversation) as returned to the
// client.
type EmailResponse struct {
	ID             int64                `json:"id"`
	AccountID      int64                `json:"account_id"`
	FolderID       int64                `json:"folder_id"`
	MessageID      string               `json:"message_id"`
	ThreadID       string               `json:"thread_id,omitempty"`
	Subject        string               `json:"subject"`
	FromAddress    string               `json:"from_address"`
	FromName       string               `json:"from_name"`
	ToAddresses    []AddressResponse    `json:"to_addresses"`
	CcAddresses    []AddressResponse    `json:"cc_addresses"`
	Date           string               `json:"date"`
	BodyText       string               `json:"body_text,omitempty"`
	BodyHTML       string               `json:"body_html,omitempty"`
	IsRead         bool                 `json:"is_read"`
	IsStarred      bool                 `json:"is_starred"`
	HasAttachments bool                 `json:"has_attachments"`
	Attachments    []AttachmentResponse `json:"attachments,omitempty"`
	InReplyTo      string               `json:"in_reply_to,omitempty"`
	References     string               `json:"references,omitempty"`
	// Conversation aggregates, present only on folder-list (collapsed) rows.
	MessageCount int     `json:"message_count,omitempty"`
	UnreadCount  int     `json:"unread_count,omitempty"`
	EmailIDs     []int64 `json:"email_ids,omitempty"`
}

// AttachmentResponse is a stored email attachment as returned to the client.
type AttachmentResponse struct {
	ID       int64  `json:"id"`
	Filename string `json:"filename"`
	MimeType string `json:"mime_type"`
	Size     int64  `json:"size"`
}

// AddressResponse is a single mail address on an email.
type AddressResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UpdateEmailRequest is the body of the endpoint that flips an email's read or
// starred state.
type UpdateEmailRequest struct {
	IsRead    *bool `json:"is_read"`
	IsStarred *bool `json:"is_starred"`
}

// ContactResult is one auto-complete contact with its usage count.
type ContactResult struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Count int    `json:"count"`
}

// BulkActionRequest is the body of the bulk-apply endpoint for a set of
// emails.
type BulkActionRequest struct {
	EmailIDs []int64 `json:"email_ids"`
	Action   string  `json:"action"`
}

// EmailTemplateRequest is the body for creating or updating an email
// template.
type EmailTemplateRequest struct {
	Name     string  `json:"name"`
	Subject  string  `json:"subject"`
	BodyHTML string  `json:"body_html"`
	BodyText string  `json:"body_text"`
	SpaceID  *string `json:"space_id"`
}

// EmailTemplateResponse is an email template as returned to the client.
type EmailTemplateResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Subject   string `json:"subject"`
	BodyHTML  string `json:"body_html"`
	BodyText  string `json:"body_text"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type SyncResponse struct {
	Synced bool `json:"synced"`
}

type FolderListResponse struct {
	Folders []FolderResponse `json:"folders"`
}

type ConversationListResponse struct {
	Emails []EmailResponse `json:"emails"`
	Total  int64           `json:"total"`
	Page   int             `json:"page"`
	Limit  int             `json:"limit"`
}

type ThreadResponse struct {
	Emails []EmailResponse `json:"emails"`
}

type SearchResponse struct {
	Emails []EmailResponse `json:"emails"`
	Total  int64           `json:"total"`
	Page   int             `json:"page"`
	Limit  int             `json:"limit"`
}

type ContactListResponse struct {
	Contacts []ContactResult `json:"contacts"`
}

type SendResponse struct {
	Sent bool `json:"sent"`
}

type DraftResponse struct {
	ID int64 `json:"id"`
}

type DeleteDraftResponse struct {
	Deleted bool `json:"deleted"`
}

type BulkActionResponse struct {
	OK bool `json:"ok"`
}

type TestConnectionResponse struct {
	OK bool `json:"ok"`
}

type TemplateListResponse struct {
	Templates []EmailTemplateResponse `json:"templates"`
}

type DeleteTemplateResponse struct {
	Deleted bool `json:"deleted"`
}
