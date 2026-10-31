package mail

import (
	"context"
	"strings"
	"unicode"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"
)

const maxFolderNameLength = 100

func folderDelimiter(client *imapclient.Client) string {
	boxes, err := client.List("", "", &imap.ListOptions{}).Collect()
	if err == nil {
		for _, box := range boxes {
			if box.Delim != 0 {
				return string(box.Delim)
			}
		}
	}
	return "/"
}

// validateFolderName rejects names IMAP cannot hold as one top-level mailbox
// (empty, the hierarchy delimiter, LIST wildcards, control characters) and
// names that detectFolderType would file under a system folder.
func validateFolderName(name, delim string) error {
	if name == "" || len(name) > maxFolderNameLength {
		return errors.Invalid("a folder name is 1 to 100 characters")
	}
	if strings.Contains(name, delim) || strings.ContainsAny(name, "*%/\\") || strings.HasPrefix(name, ".") {
		return errors.Invalid("a folder name cannot contain " + delim + " * % / \\ or start with a dot")
	}
	if strings.IndexFunc(name, unicode.IsControl) >= 0 {
		return errors.Invalid("a folder name cannot contain control characters")
	}
	if detectFolderType(&imap.ListData{Mailbox: name}) != schemas.FolderTypeCustom {
		return errors.Invalid("that name is reserved for a system folder (inbox, sent, drafts, trash, junk, archive)")
	}
	return nil
}

func (s *Service) pathTaken(ctx context.Context, accountID int64, path string, exceptID int64) bool {
	var count int64
	s.orm.WithContext(ctx).Model(&schemas.Folder{}).
		Where("account_id = ? AND lower(path) = lower(?) AND id <> ?", accountID, path, exceptID).Count(&count)
	return count > 0
}

func siblingPath(path, name, delim string) string {
	if i := strings.LastIndex(path, delim); i >= 0 {
		return path[:i+len(delim)] + name
	}
	return name
}
