package mcp

// accountsInput is the input of courrier_list_accounts.
type accountsInput struct {
	ResponseFormat string `json:"response_format,omitempty" jsonschema:"concise (default) or detailed, which adds each account's folders with unread counts"`
}

// accountsOutput lists the mail accounts a person can use.
type accountsOutput struct {
	Accounts []accountSummary `json:"accounts"`
}

// accountSummary is one mail account.
type accountSummary struct {
	AccountID int64           `json:"account_id"`
	Email     string          `json:"email"`
	Name      string          `json:"name"`
	Shared    bool            `json:"shared"`
	Default   bool            `json:"default"`
	Folders   []folderSummary `json:"folders,omitempty"`
}

// folderSummary is one folder of an account, in detailed output.
type folderSummary struct {
	FolderID int64  `json:"folder_id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Unread   int    `json:"unread"`
	Total    int    `json:"total"`
}

// searchInput is the input of courrier_search_mail.
type searchInput struct {
	Query          string `json:"query,omitempty" jsonschema:"words that must all appear in the subject, sender or body. Omit to list the newest mail"`
	AccountID      int64  `json:"account_id,omitempty" jsonschema:"search one account only. Omit to search every account"`
	Folder         string `json:"folder,omitempty" jsonschema:"inbox, sent, drafts, archive, trash, junk or any. Omit to search everything except trash and junk"`
	UnreadOnly     bool   `json:"unread_only,omitempty" jsonschema:"only unread messages"`
	Since          string `json:"since,omitempty" jsonschema:"only mail received on or after this date, as YYYY-MM-DD or RFC 3339"`
	Limit          int    `json:"limit,omitempty" jsonschema:"messages per page, 1 to 50. Default 20"`
	Cursor         string `json:"cursor,omitempty" jsonschema:"next_cursor from the previous page of the same search"`
	ResponseFormat string `json:"response_format,omitempty" jsonschema:"concise (default) or detailed, which adds recipients, folder and flags"`
}

// searchOutput is one page of messages, newest first.
type searchOutput struct {
	Messages   []messageSummary `json:"messages"`
	NextCursor string           `json:"next_cursor,omitempty"`
	Note       string           `json:"note,omitempty"`
}

// messageSummary is one message in a result list.
type messageSummary struct {
	MessageID int64  `json:"message_id"`
	AccountID int64  `json:"account_id"`
	From      string `json:"from"`
	Subject   string `json:"subject"`
	Date      string `json:"date"`
	Snippet   string `json:"snippet"`
	Unread    bool   `json:"unread"`
	messageDetail
}

// messageDetail is what detailed output adds to a message.
type messageDetail struct {
	To             []string `json:"to,omitempty"`
	Cc             []string `json:"cc,omitempty"`
	Folder         string   `json:"folder,omitempty"`
	Starred        bool     `json:"starred,omitempty"`
	HasAttachments bool     `json:"has_attachments,omitempty"`
}

// threadInput is the input of courrier_read_thread.
type threadInput struct {
	AccountID      int64  `json:"account_id" jsonschema:"the account_id of the message, from courrier_search_mail"`
	MessageID      int64  `json:"message_id" jsonschema:"any message_id of the conversation"`
	Cursor         string `json:"cursor,omitempty" jsonschema:"next_cursor from the previous page of this conversation"`
	ResponseFormat string `json:"response_format,omitempty" jsonschema:"concise (default) trims quoted replies and long bodies; detailed keeps them, up to 20000 characters a message"`
}

// threadOutput is one page of a conversation, oldest message first.
type threadOutput struct {
	Subject      string          `json:"subject"`
	MessageCount int             `json:"message_count"`
	Messages     []threadMessage `json:"messages"`
	NextCursor   string          `json:"next_cursor,omitempty"`
	Note         string          `json:"note,omitempty"`
}

// threadMessage is one message of a conversation with its body.
type threadMessage struct {
	MessageID      int64    `json:"message_id"`
	From           string   `json:"from"`
	To             []string `json:"to"`
	Cc             []string `json:"cc,omitempty"`
	Date           string   `json:"date"`
	Body           string   `json:"body"`
	BodyTruncated  bool     `json:"body_truncated,omitempty"`
	Unread         bool     `json:"unread,omitempty"`
	HasAttachments bool     `json:"has_attachments,omitempty"`
}
