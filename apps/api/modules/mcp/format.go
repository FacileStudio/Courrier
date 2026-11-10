package mcp

import (
	"cmp"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/modules/mail"
	"github.com/FacileStudio/Courrier/apps/api/schemas"
)

// snippetLength caps the snippet of a message in a result list.
const snippetLength = 200

// summarize renders one message of a result list.
func summarize(email schemas.Email, folders map[int64]schemas.Folder, full bool) messageSummary {
	snippet, _ := clip(strings.Join(strings.Fields(email.BodyText), " "), snippetLength)
	summary := messageSummary{
		MessageID: email.ID, AccountID: email.AccountID, From: sender(email), Subject: email.Subject,
		Date: email.Date.UTC().Format(time.RFC3339), Snippet: snippet, Unread: !email.IsRead,
	}
	if full {
		summary.messageDetail = messageDetail{
			To: mail.AddressList(email.ToAddresses), Cc: mail.AddressList(email.CcAddresses),
			Folder: folderLabel(folders[email.FolderID]), Starred: email.IsStarred, HasAttachments: email.HasAttachments,
		}
	}
	return summary
}

// sender is "Name <address>", or the bare address when the name is unknown.
func sender(email schemas.Email) string {
	if email.FromName == "" || email.FromName == email.FromAddress {
		return email.FromAddress
	}
	return email.FromName + " <" + email.FromAddress + ">"
}

// folderLabel names a folder by its type, or by its name for the person's
// own folders.
func folderLabel(folder schemas.Folder) string {
	if folder.Type == schemas.FolderTypeCustom || folder.Type == "" {
		return folder.Name
	}
	return folder.Type
}

// sortedFolders orders folders by account, type, name, then id, so detailed
// account listings come out the same every time.
func sortedFolders(folders map[int64]schemas.Folder) []schemas.Folder {
	return slices.SortedFunc(maps.Values(folders), func(a, b schemas.Folder) int {
		return cmp.Or(cmp.Compare(a.AccountID, b.AccountID), cmp.Compare(a.Type, b.Type),
			cmp.Compare(a.Name, b.Name), cmp.Compare(a.ID, b.ID))
	})
}

// withoutQuotes drops the quoted part of a reply: lines starting with ">",
// and everything from an "On ... wrote:" line onwards.
func withoutQuotes(body string) string {
	lines := strings.Split(body, "\n")
	kept := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasSuffix(trimmed, "wrote:") || strings.HasSuffix(trimmed, "a écrit :") {
			break
		}
		if !strings.HasPrefix(trimmed, ">") {
			kept = append(kept, line)
		}
	}
	return strings.TrimSpace(strings.Join(kept, "\n"))
}
