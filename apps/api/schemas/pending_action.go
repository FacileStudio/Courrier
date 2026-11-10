package schemas

import (
	"time"

	"gorm.io/gorm"
)

// pendingActionsSchema holds the sensitive actions an MCP client asked for
// and a person has not confirmed yet: sending a draft or deleting one. The id
// is a random token that appears in the confirmation link. A row is resolved
// once, by the claim that runs it or by a dismissal, and expires unconfirmed.
const pendingActionsSchema = `
CREATE TABLE IF NOT EXISTS mcp_pending_actions (
	id          text        PRIMARY KEY,
	user_id     bigint      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	account_id  bigint      NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
	email_id    bigint      NOT NULL,
	kind        text        NOT NULL,
	digest      text        NOT NULL,
	client_id   text        NOT NULL DEFAULT '',
	created_at  timestamptz NOT NULL DEFAULT now(),
	expires_at  timestamptz NOT NULL,
	resolved_at timestamptz,
	outcome     text        NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS mcp_pending_actions_open_idx
	ON mcp_pending_actions (user_id, created_at DESC) WHERE resolved_at IS NULL;
`

// PendingAction is one row of mcp_pending_actions.
type PendingAction struct {
	ID         string     `gorm:"column:id;primaryKey"`
	UserID     int64      `gorm:"column:user_id"`
	AccountID  int64      `gorm:"column:account_id"`
	EmailID    int64      `gorm:"column:email_id"`
	Kind       string     `gorm:"column:kind"`
	Digest     string     `gorm:"column:digest"`
	ClientID   string     `gorm:"column:client_id"`
	CreatedAt  time.Time  `gorm:"column:created_at;autoCreateTime"`
	ExpiresAt  time.Time  `gorm:"column:expires_at"`
	ResolvedAt *time.Time `gorm:"column:resolved_at"`
	Outcome    string     `gorm:"column:outcome"`
}

// TableName pins the GORM model to the raw table.
func (PendingAction) TableName() string { return "mcp_pending_actions" }

// ensurePendingActions creates the mcp_pending_actions table. Idempotent.
func ensurePendingActions(db *gorm.DB) error {
	return db.Exec(pendingActionsSchema).Error
}
