package auth

import (
	"context"
	"database/sql"
)

const personalAccounts = `SELECT id FROM accounts WHERE user_id = $1 AND space_id IS NULL`

// eraseStatements delete what an erased person owns, children first. Mail
// accounts and templates attached to a space belong to the space and stay.
// Sessions, API tokens, identities and push devices cascade from the users
// row the kit deletes after these run.
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
	}
}

// eraseData is the kit's EraseData hook: it runs inside the erase
// transaction, before the kit deletes the users row.
func eraseData(ctx context.Context, tx *sql.Tx, userID int64) error {
	for _, statement := range eraseStatements() {
		if _, err := tx.ExecContext(ctx, statement, userID); err != nil {
			return err
		}
	}
	return nil
}
