package mcp

import mcpkit "github.com/FacileStudio/tronc/mcp"

// PendingItem is one action waiting on the confirm page.
type PendingItem struct {
	mcpkit.PendingItem
	AccountID    int64         `json:"account_id"`
	AccountEmail string        `json:"account_email"`
	Draft        *PendingDraft `json:"draft,omitempty"`
}

// PendingDraft is the draft an action would send or delete, as it is now.
type PendingDraft struct {
	ID      int64    `json:"id"`
	To      []string `json:"to"`
	Cc      []string `json:"cc"`
	Bcc     []string `json:"bcc"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}
