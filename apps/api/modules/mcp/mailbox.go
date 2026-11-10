package mcp

import (
	"context"

	"github.com/FacileStudio/Courrier/apps/api/modules/mail"
	"github.com/FacileStudio/Courrier/apps/api/schemas"
)

// Finder is the read side of the mailbox the tools search.
type Finder interface {
	VisibleAccounts(ctx context.Context, userID int64) ([]schemas.Account, error)
	SearchMessages(ctx context.Context, userID int64, q mail.MessageQuery) ([]schemas.Email, error)
	FoldersByID(ctx context.Context, accountIDs []int64) (map[int64]schemas.Folder, error)
}

// Reader loads one message or a whole conversation.
type Reader interface {
	GetEmail(ctx context.Context, userID, accountID, emailID int64) (schemas.Email, error)
	GetThread(ctx context.Context, userID, accountID int64, threadID string) ([]schemas.Email, error)
}

// Drafter writes and removes drafts.
type Drafter interface {
	Draft(ctx context.Context, userID, accountID, draftID int64) (schemas.Email, error)
	SaveDraft(ctx context.Context, userID, accountID int64, req mail.SendRequest) (schemas.Email, error)
	UpdateDraft(ctx context.Context, userID, accountID, emailID int64, req mail.DraftRequest) (schemas.Email, error)
	DeleteDraft(ctx context.Context, userID, accountID, emailID int64) error
}

// Organizer changes flags and folders, and sends drafts.
type Organizer interface {
	SetReadState(ctx context.Context, userID, accountID int64, emailIDs []int64, isRead bool) error
	SetStarState(ctx context.Context, userID, accountID int64, emailIDs []int64, starred bool) error
	MoveEmailsTo(ctx context.Context, userID, accountID int64, emailIDs []int64, req mail.MoveEmailsRequest) error
	SendDraft(ctx context.Context, userID, accountID, draftID int64) error
}

// Mailbox is everything the tools do with mail. *mail.Service implements it,
// and the tests replace it with an in-memory one.
type Mailbox interface {
	Finder
	Reader
	Drafter
	Organizer
}
