package schemas

import (
	"github.com/FacileStudio/tronc/mcp"
	"gorm.io/gorm"
)

// ensurePendingActions creates tronc/mcp's mcp_actions table, which holds
// the sensitive actions an MCP client asked for and a person has not
// confirmed yet. Idempotent. The mcp_pending_actions table of v0.14.0 is
// no longer read; it is left in place.
func ensurePendingActions(db *gorm.DB) error {
	return db.Exec(mcp.PendingSchema).Error
}
