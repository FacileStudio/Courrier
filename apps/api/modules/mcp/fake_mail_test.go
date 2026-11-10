package mcp

import (
	"cmp"
	"context"
	"encoding/json"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/modules/mail"
	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"
)

// Folder ids of the fake account.
const (
	inboxID  int64 = 1
	draftsID int64 = 2
)

// fakeMail is an in-memory Mailbox over one account, recording sends.
type fakeMail struct {
	mu      sync.Mutex
	account schemas.Account
	emails  map[int64]schemas.Email
	nextID  int64
	sent    []mail.SendRequest
	removed map[int64]int
}

// newFakeMail holds one account and no mail.
func newFakeMail(accountID int64) *fakeMail {
	return &fakeMail{
		account: schemas.Account{ID: accountID, Email: "me@example.test", Name: "Me", IsDefault: true},
		emails:  map[int64]schemas.Email{}, nextID: 100, removed: map[int64]int{},
	}
}

// add stores a message in the inbox and returns it.
func (f *fakeMail) add(email schemas.Email) schemas.Email {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.nextID++
	email.ID, email.AccountID = f.nextID, f.account.ID
	if email.FolderID == 0 {
		email.FolderID = inboxID
	}
	f.emails[email.ID] = email
	return email
}

// sentCount is how many messages went out.
func (f *fakeMail) sentCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.sent)
}

// removals is how many times draft id was sent or deleted.
func (f *fakeMail) removals(id int64) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.removed[id]
}

// VisibleAccounts returns the one account.
func (f *fakeMail) VisibleAccounts(context.Context, int64) ([]schemas.Account, error) {
	return []schemas.Account{f.account}, nil
}

// FoldersByID returns the inbox and the drafts folder.
func (f *fakeMail) FoldersByID(context.Context, []int64) (map[int64]schemas.Folder, error) {
	return map[int64]schemas.Folder{
		inboxID:  {ID: inboxID, AccountID: f.account.ID, Type: schemas.FolderTypeInbox, Name: "Inbox"},
		draftsID: {ID: draftsID, AccountID: f.account.ID, Type: schemas.FolderTypeDrafts, Name: "Drafts"},
	}, nil
}

// SearchMessages filters by text and folder, newest first, after Before.
func (f *fakeMail) SearchMessages(_ context.Context, _ int64, q mail.MessageQuery) ([]schemas.Email, error) {
	if len(q.AccountIDs) > 0 && q.AccountIDs[0] != f.account.ID {
		return nil, errors.NotFound("account not found")
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	var found []schemas.Email
	for _, email := range f.emails {
		if matches(email, q) {
			found = append(found, email)
		}
	}
	slices.SortFunc(found, func(a, b schemas.Email) int {
		return cmp.Or(b.Date.Compare(a.Date), cmp.Compare(b.ID, a.ID))
	})
	return found[:min(q.Limit, len(found))], nil
}

// matches applies the query filters the fake supports.
func matches(email schemas.Email, q mail.MessageQuery) bool {
	if q.Text != "" && !strings.Contains(strings.ToLower(email.Subject), strings.ToLower(q.Text)) {
		return false
	}
	if q.FolderType == schemas.FolderTypeDrafts && email.FolderID != draftsID {
		return false
	}
	if q.Before == nil {
		return true
	}
	return email.Date.Before(q.Before.Date) || (email.Date.Equal(q.Before.Date) && email.ID < q.Before.ID)
}

// addressJSON stores addresses the way the mail module does.
func addressJSON(addresses []string) string {
	type entry struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	entries := make([]entry, len(addresses))
	for i, address := range addresses {
		entries[i] = entry{Email: address}
	}
	raw, _ := json.Marshal(entries)
	return string(raw)
}

// draftOf turns a send request into a stored draft.
func draftOf(req mail.SendRequest) schemas.Email {
	return schemas.Email{
		FolderID: draftsID, Subject: req.Subject, BodyText: req.Body, Date: time.Now(), IsRead: true,
		ToAddresses: addressJSON(req.To), CcAddresses: addressJSON(req.Cc), BccAddresses: addressJSON(req.Bcc),
		InReplyTo: req.InReplyTo, References: strings.Join(req.References, " "),
	}
}
