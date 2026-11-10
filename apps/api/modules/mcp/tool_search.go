package mcp

import (
	"context"
	"errors"
	"slices"
	"time"

	mcpkit "github.com/FacileStudio/tronc/mcp"

	"github.com/FacileStudio/Courrier/apps/api/modules/mail"
	"github.com/FacileStudio/Courrier/apps/api/schemas"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Page sizes of courrier_search_mail.
const (
	defaultSearchLimit = 20
	maxSearchLimit     = 50
)

// searchFolders are the folder values courrier_search_mail accepts.
var searchFolders = []string{
	schemas.FolderTypeInbox, schemas.FolderTypeSent, schemas.FolderTypeDrafts, schemas.FolderTypeArchive,
	schemas.FolderTypeTrash, schemas.FolderTypeJunk, mail.FolderAny,
}

// searchPosition is the cursor of courrier_search_mail: the last message of
// the previous page.
type searchPosition struct {
	Date int64 `json:"d"`
	ID   int64 `json:"i"`
}

// searchMail is courrier_search_mail.
func (t *toolset) searchMail(ctx context.Context, req *sdk.CallToolRequest, in searchInput) (*sdk.CallToolResult, searchOutput, error) {
	who, err := callerOf(req)
	if err != nil {
		return nil, searchOutput{}, err
	}
	query, full, err := messageQuery(in)
	if err != nil {
		return nil, searchOutput{}, err
	}
	emails, err := t.mail.SearchMessages(ctx, who.userID, query)
	if err != nil {
		return nil, searchOutput{}, mcpkit.Failure(err, hintAccount)
	}
	folders, err := t.mail.FoldersByID(ctx, accountIDs(emails))
	if err != nil {
		return nil, searchOutput{}, mcpkit.Failure(err, "")
	}
	return nil, searchPage(emails, query.Limit-1, folders, full), nil
}

// messageQuery validates the input and turns it into the mail query, asking
// for one message more than the page so the page knows whether more follow.
func messageQuery(in searchInput) (mail.MessageQuery, bool, error) {
	full, err := mcpkit.Detailed(in.ResponseFormat)
	if err != nil {
		return mail.MessageQuery{}, false, err
	}
	limit := in.Limit
	if limit == 0 {
		limit = defaultSearchLimit
	}
	if limit < 1 || limit > maxSearchLimit {
		return mail.MessageQuery{}, false, errors.New("limit must be between 1 and 50")
	}
	if in.Folder != "" && !slices.Contains(searchFolders, in.Folder) {
		return mail.MessageQuery{}, false, errors.New("folder must be one of inbox, sent, drafts, archive, trash, junk or any")
	}
	query := mail.MessageQuery{Text: in.Query, FolderType: in.Folder, UnreadOnly: in.UnreadOnly, Limit: limit + 1}
	if in.AccountID != 0 {
		query.AccountIDs = []int64{in.AccountID}
	}
	if query.Since, err = parseSince(in.Since); err != nil {
		return query, false, err
	}
	query.Before, err = searchCursor(in.Cursor)
	return query, full, err
}

// searchPage renders up to limit messages and, when more matched or the
// budget cut the page short, the cursor that resumes after the last one kept.
func searchPage(emails []schemas.Email, limit int, folders map[int64]schemas.Folder, full bool) searchOutput {
	summaries := make([]messageSummary, 0, len(emails))
	for _, email := range emails {
		summaries = append(summaries, summarize(email, folders, full))
	}
	more := len(summaries) > limit
	keep := mcpkit.FitCount(min(limit, len(summaries)), func(n int) int {
		return mcpkit.EncodedSize(searchOutput{Messages: summaries[:n], NextCursor: "next-cursor-placeholder"})
	})
	out := searchOutput{Messages: summaries[:keep]}
	if keep < min(limit, len(summaries)) {
		out.Note = "The page was cut to stay under the size limit. Narrow the search, or continue with next_cursor."
	}
	if (more || keep < len(summaries)) && keep > 0 {
		last := emails[keep-1]
		out.NextCursor = mcpkit.EncodeCursor(searchPosition{Date: last.Date.UnixNano(), ID: last.ID})
	}
	if len(out.Messages) == 0 {
		out.Note = "No message matched. Try fewer words, another folder, or folder \"any\"."
	}
	return out
}

// parseSince reads a date filter.
func parseSince(raw string) (time.Time, error) {
	if raw == "" {
		return time.Time{}, nil
	}
	for _, layout := range []string{time.DateOnly, time.RFC3339} {
		if since, err := time.Parse(layout, raw); err == nil {
			return since, nil
		}
	}
	return time.Time{}, errors.New("since must be a date like 2026-11-01 or an RFC 3339 time")
}

// searchCursor reads the position a search cursor names.
func searchCursor(cursor string) (*mail.MessageKey, error) {
	if cursor == "" {
		return nil, nil
	}
	var position searchPosition
	if err := mcpkit.DecodeCursor(cursor, &position); err != nil || position.ID == 0 {
		return nil, mcpkit.ErrCursor
	}
	return &mail.MessageKey{Date: time.Unix(0, position.Date), ID: position.ID}, nil
}

// accountIDs lists the distinct accounts of emails.
func accountIDs(emails []schemas.Email) []int64 {
	ids := []int64{}
	for _, email := range emails {
		if !slices.Contains(ids, email.AccountID) {
			ids = append(ids, email.AccountID)
		}
	}
	return ids
}
