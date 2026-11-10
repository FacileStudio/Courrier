package mcp

// draftInput is the input of courrier_draft, and of courrier_send when it
// composes the message itself.
type draftInput struct {
	AccountID        int64    `json:"account_id" jsonschema:"the account to write from, from courrier_list_accounts"`
	DraftID          int64    `json:"draft_id,omitempty" jsonschema:"an existing draft to replace. Omit to create a new draft"`
	ReplyToMessageID int64    `json:"reply_to_message_id,omitempty" jsonschema:"answer this message: threads the draft, and fills to and subject when you leave them empty"`
	To               []string `json:"to,omitempty" jsonschema:"recipients, as email addresses or Name <email>"`
	Cc               []string `json:"cc,omitempty" jsonschema:"copy recipients"`
	Bcc              []string `json:"bcc,omitempty" jsonschema:"blind copy recipients"`
	Subject          string   `json:"subject,omitempty" jsonschema:"the subject line"`
	Body             string   `json:"body,omitempty" jsonschema:"the message, as plain text"`
}

// draftOutput confirms a saved draft.
type draftOutput struct {
	DraftID   int64    `json:"draft_id"`
	AccountID int64    `json:"account_id"`
	To        []string `json:"to"`
	Subject   string   `json:"subject"`
	ReviewURL string   `json:"review_url"`
	Next      string   `json:"next"`
}

// draftRef names one draft.
type draftRef struct {
	AccountID int64 `json:"account_id" jsonschema:"the account of the draft"`
	DraftID   int64 `json:"draft_id" jsonschema:"the draft to delete"`
}

// organizeInput is the input of courrier_organize.
type organizeInput struct {
	AccountID  int64   `json:"account_id" jsonschema:"the account of the messages"`
	MessageIDs []int64 `json:"message_ids" jsonschema:"1 to 100 message ids from the same account"`
	Action     string  `json:"action" jsonschema:"mark_read, mark_unread, star, unstar, archive, inbox (move back to the inbox), trash (move to Trash, never a permanent delete) or move"`
	FolderID   int64   `json:"folder_id,omitempty" jsonschema:"the destination for action move, from courrier_list_accounts in detailed format"`
}

// actionOutput reports what a write did, or why it did not happen yet.
type actionOutput struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	DraftID    int64  `json:"draft_id,omitempty"`
	Count      int    `json:"count,omitempty"`
	ConfirmURL string `json:"confirm_url,omitempty"`
}
