package mcp

import (
	"context"

	"github.com/FacileStudio/Courrier/apps/api/modules/mail"
)

// SetReadState flips is_read on the given messages.
func (f *fakeMail) SetReadState(_ context.Context, _, _ int64, ids []int64, read bool) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, id := range ids {
		if email, ok := f.emails[id]; ok {
			email.IsRead = read
			f.emails[id] = email
		}
	}
	return nil
}

// SetStarState flips is_starred on the given messages.
func (f *fakeMail) SetStarState(_ context.Context, _, _ int64, ids []int64, starred bool) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, id := range ids {
		if email, ok := f.emails[id]; ok {
			email.IsStarred = starred
			f.emails[id] = email
		}
	}
	return nil
}

// MoveEmailsTo accepts every move.
func (f *fakeMail) MoveEmailsTo(context.Context, int64, int64, []int64, mail.MoveEmailsRequest) error {
	return nil
}
