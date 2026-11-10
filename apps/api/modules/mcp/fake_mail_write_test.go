package mcp

import (
	"cmp"
	"context"
	"slices"

	"github.com/FacileStudio/Courrier/apps/api/modules/mail"
	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"
)

// GetEmail returns a message of the account.
func (f *fakeMail) GetEmail(_ context.Context, _, accountID, emailID int64) (schemas.Email, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	email, ok := f.emails[emailID]
	if !ok || accountID != f.account.ID {
		return schemas.Email{}, errors.NotFound("email not found")
	}
	return email, nil
}

// GetThread returns a conversation, oldest first.
func (f *fakeMail) GetThread(_ context.Context, _, _ int64, threadID string) ([]schemas.Email, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var thread []schemas.Email
	for _, email := range f.emails {
		if email.ThreadID == threadID {
			thread = append(thread, email)
		}
	}
	slices.SortFunc(thread, func(a, b schemas.Email) int { return cmp.Or(a.Date.Compare(b.Date), cmp.Compare(a.ID, b.ID)) })
	return thread, nil
}

// Draft returns a draft of the account.
func (f *fakeMail) Draft(ctx context.Context, userID, accountID, draftID int64) (schemas.Email, error) {
	email, err := f.GetEmail(ctx, userID, accountID, draftID)
	if err != nil || email.FolderID != draftsID {
		return schemas.Email{}, errors.NotFound("draft not found")
	}
	return email, nil
}

// SaveDraft stores a new draft.
func (f *fakeMail) SaveDraft(_ context.Context, _, accountID int64, req mail.SendRequest) (schemas.Email, error) {
	if accountID != f.account.ID {
		return schemas.Email{}, errors.NotFound("account not found")
	}
	return f.add(draftOf(req)), nil
}

// UpdateDraft replaces a draft in place.
func (f *fakeMail) UpdateDraft(ctx context.Context, userID, accountID, emailID int64, req mail.DraftRequest) (schemas.Email, error) {
	if _, err := f.Draft(ctx, userID, accountID, emailID); err != nil {
		return schemas.Email{}, err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	draft := draftOf(req.SendRequest)
	draft.ID, draft.AccountID = emailID, accountID
	f.emails[emailID] = draft
	return draft, nil
}

// DeleteDraft removes a draft.
func (f *fakeMail) DeleteDraft(ctx context.Context, userID, accountID, emailID int64) error {
	if _, err := f.Draft(ctx, userID, accountID, emailID); err != nil {
		return err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.emails, emailID)
	f.removed[emailID]++
	return nil
}

// SendDraft records the draft as sent and removes it.
func (f *fakeMail) SendDraft(ctx context.Context, userID, accountID, draftID int64) error {
	draft, err := f.Draft(ctx, userID, accountID, draftID)
	if err != nil {
		return err
	}
	f.mu.Lock()
	f.sent = append(f.sent, mail.DraftMessage(draft))
	f.mu.Unlock()
	return f.DeleteDraft(ctx, userID, accountID, draftID)
}
