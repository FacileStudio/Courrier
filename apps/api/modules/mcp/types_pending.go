package mcp

import "time"

// PendingItem is one action waiting on the confirm page.
type PendingItem struct {
	ID           string        `json:"id"`
	Kind         string        `json:"kind"`
	ClientID     string        `json:"client_id"`
	CreatedAt    time.Time     `json:"created_at"`
	ExpiresAt    time.Time     `json:"expires_at"`
	AccountID    int64         `json:"account_id"`
	AccountEmail string        `json:"account_email"`
	State        string        `json:"state"`
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

// PendingList is the body of GET /api/mcp/pending.
type PendingList struct {
	Actions []PendingItem `json:"actions"`
}

// PendingResult is the body answered when an action is confirmed or dismissed.
type PendingResult struct {
	Status string `json:"status"`
}

// EmptyRequest documents a route that takes no body.
type EmptyRequest struct{}
