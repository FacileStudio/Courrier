package auth

import (
	"context"
	"database/sql"
	"strconv"
)

const personalAccounts = `SELECT id FROM accounts WHERE user_id = $1 AND space_id IS NULL`

// eraseStatements delete what an erased person owns, children first. Mail
// accounts and templates attached to a space belong to the space and stay.
// Sessions and API tokens (porte_sessions) and push devices also cascade from
// the users row, and are deleted here so the list names every table. porte
// deletes the identity rows itself; the v0.14.0 mcp_pending_actions cascade.
func eraseStatements() []string {
	return []string{
		`DELETE FROM attachments WHERE email_id IN (SELECT id FROM emails WHERE account_id IN (` + personalAccounts + `))`,
		`DELETE FROM emails WHERE account_id IN (` + personalAccounts + `)`,
		`DELETE FROM folders WHERE account_id IN (` + personalAccounts + `)`,
		`DELETE FROM thread_links WHERE account_id IN (` + personalAccounts + `)`,
		`DELETE FROM sender_rules WHERE account_id IN (` + personalAccounts + `)`,
		`DELETE FROM contacts WHERE account_id IN (` + personalAccounts + `)`,
		`DELETE FROM accounts WHERE user_id = $1 AND space_id IS NULL`,
		`DELETE FROM email_templates WHERE user_id = $1 AND space_id IS NULL`,
		`DELETE FROM space_members WHERE user_id = $1`,
		`DELETE FROM push_devices WHERE user_id = $1`,
		`DELETE FROM porte_sessions WHERE user_id = $1`,
	}
}

// eraseData is the kit's EraseData hook: it runs inside the erase
// transaction, before the kit deletes the users row. tronc/mcp keys
// mcp_actions by a text user id, so that table takes the id as text.
func eraseData(ctx context.Context, tx *sql.Tx, userID int64) error {
	for _, statement := range eraseStatements() {
		if _, err := tx.ExecContext(ctx, statement, userID); err != nil {
			return err
		}
	}
	_, err := tx.ExecContext(ctx, `DELETE FROM mcp_actions WHERE user_id = $1`, strconv.FormatInt(userID, 10))
	return err
}
