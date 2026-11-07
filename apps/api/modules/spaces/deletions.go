package spaces

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/FacileStudio/porte/directory"

	"gorm.io/gorm"
)

// Deletions erases the people Registre erased. It polls Registre's deletion
// list with the same client credentials as the snapshot sync and accepts
// Registre's signed push on directory.DeletionsPath.
type Deletions struct {
	orm        *gorm.DB
	storageDir string
	receiver   *directory.Deletions
}

// NewDeletions builds the receiver over Courrier's tables. erase is porte's
// own rows eraser, pg.Store.EraseIdentity.
func NewDeletions(
	orm *gorm.DB, client directory.Client, erase func(context.Context, string, string) error,
	storageDir string, logger *slog.Logger,
) *Deletions {
	d := &Deletions{orm: orm, storageDir: storageDir}
	d.receiver = &directory.Deletions{
		Client:          client,
		OnPersonDeleted: d.ErasePerson,
		Erase:           erase,
		State:           deletionState{orm: orm},
		Logger:          logger,
	}
	return d
}

// Run polls until ctx ends. Start it in exactly one replica.
func (d *Deletions) Run(ctx context.Context) {
	d.receiver.Run(ctx)
}

// Handler verifies and applies Registre's signed user.deleted push.
func (d *Deletions) Handler() http.Handler {
	return d.receiver.Handler()
}

const (
	personalAccounts = `SELECT id FROM accounts WHERE user_id = @id AND space_id IS NULL`

	deleteOwnedAttachments = `DELETE FROM attachments WHERE email_id IN (
		SELECT id FROM emails WHERE account_id IN (` + personalAccounts + `))`
	deleteOwnedEmails  = `DELETE FROM emails WHERE account_id IN (` + personalAccounts + `)`
	deleteOwnedFolders = `DELETE FROM folders WHERE account_id IN (` + personalAccounts + `)`
	deleteOwnedLinks   = `DELETE FROM thread_links WHERE account_id IN (` + personalAccounts + `)`
	deleteOwnedRules   = `DELETE FROM sender_rules WHERE account_id IN (` + personalAccounts + `)`
	deleteOwnedContact = `DELETE FROM contacts WHERE account_id IN (` + personalAccounts + `)`
	deleteOwnedAccount = `DELETE FROM accounts WHERE user_id = @id AND space_id IS NULL`
	deleteTemplates    = `DELETE FROM email_templates WHERE user_id = @id AND space_id IS NULL`
	deleteMemberships  = `DELETE FROM space_members WHERE user_id = @id`
	deleteAccount      = `DELETE FROM users WHERE id = @id`
)

// ErasePerson deletes the account behind a Registre identity and what it owns:
// its personal mail accounts with their folders, cached messages, contacts and
// rules, its personal templates, its memberships, its devices and sessions.
// Mail accounts and templates attached to a space belong to the space and
// stay. A subject that is already gone is a no-op, so a replayed notice
// changes nothing.
func (d *Deletions) ErasePerson(ctx context.Context, provider, subject string) error {
	var account struct {
		ID         int64
		UploadPath string
	}
	err := d.orm.WithContext(ctx).Raw(
		`SELECT u.id, coalesce(u.avatar_upload_path, '') AS upload_path FROM users u
		 JOIN porte_identities i ON i.user_id = u.id WHERE i.provider = ? AND i.subject = ?`,
		provider, subject).Scan(&account).Error
	if err != nil || account.ID == 0 {
		return err
	}
	err = d.orm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, statement := range []string{
			deleteOwnedAttachments, deleteOwnedEmails, deleteOwnedFolders, deleteOwnedLinks,
			deleteOwnedRules, deleteOwnedContact, deleteOwnedAccount, deleteTemplates,
			deleteMemberships, deleteAccount,
		} {
			if err := tx.Exec(statement, map[string]any{"id": account.ID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err == nil && account.UploadPath != "" {
		removeUpload(d.storageDir, account.UploadPath)
	}
	return err
}

func removeUpload(storageDir, relative string) {
	path := filepath.Join(storageDir, filepath.Clean("/"+relative))
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		slog.Warn("failed to remove an erased person's avatar", slog.Any("error", err))
	}
}

// deletionState keeps the receiver's cursor and handled subjects in one row.
type deletionState struct {
	orm *gorm.DB
}

// Load returns the saved state, or nil before the first save.
func (s deletionState) Load(ctx context.Context) ([]byte, error) {
	var rows []struct{ State []byte }
	err := s.orm.WithContext(ctx).Raw(`SELECT state FROM directory_deletion_state WHERE id = 1`).Scan(&rows).Error
	if err != nil || len(rows) == 0 {
		return nil, err
	}
	return rows[0].State, nil
}

// Save replaces the saved state.
func (s deletionState) Save(ctx context.Context, state []byte) error {
	return s.orm.WithContext(ctx).Exec(
		`INSERT INTO directory_deletion_state (id, state, updated_at) VALUES (1, ?, now())
		 ON CONFLICT (id) DO UPDATE SET state = EXCLUDED.state, updated_at = now()`, state).Error
}
