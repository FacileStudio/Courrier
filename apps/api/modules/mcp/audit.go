package mcp

import (
	"context"
	"log/slog"
)

// event is one audited write: the tool, how it ended, and the account and
// draft or message it touched.
type event struct {
	tool      string
	outcome   string
	accountID int64
	targetID  int64
}

// audit logs one write: who asked, through which client, and how it ended.
// It never logs content or recipients. The logger ships to Journal.
func (t *toolset) audit(ctx context.Context, who caller, e event) {
	t.logger.InfoContext(ctx, "mcp write",
		slog.String("tool", e.tool),
		slog.Int64("user_id", who.userID),
		slog.String("client_id", who.clientID),
		slog.String("actor", who.actor),
		slog.Int64("account_id", e.accountID),
		slog.Int64("target_id", e.targetID),
		slog.String("outcome", e.outcome),
	)
}

// outcomeOf is the audit outcome of a write that returned err.
func outcomeOf(err error) string {
	if err != nil {
		return "failed"
	}
	return outcomeDone
}
