package mcp

import (
	"context"
	"errors"
	"fmt"

	mcpkit "github.com/FacileStudio/tronc/mcp"

	"github.com/FacileStudio/Courrier/apps/api/modules/mail"
	"github.com/FacileStudio/Courrier/apps/api/schemas"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// maxOrganize caps the messages one courrier_organize call changes.
const maxOrganize = 100

// errAction lists the actions courrier_organize accepts.
var errAction = errors.New("action must be mark_read, mark_unread, star, unstar, archive, inbox, trash or move")

// moves maps the folder actions to the system folder they move to.
var moves = map[string]string{
	"archive": schemas.FolderTypeArchive, "inbox": schemas.FolderTypeInbox, "trash": schemas.FolderTypeTrash,
}

// organize is courrier_organize.
func (t *toolset) organize(ctx context.Context, req *sdk.CallToolRequest, in organizeInput) (*sdk.CallToolResult, actionOutput, error) {
	who, err := callerOf(req)
	if err != nil {
		return nil, actionOutput{}, err
	}
	if len(in.MessageIDs) == 0 || len(in.MessageIDs) > maxOrganize {
		return nil, actionOutput{}, fmt.Errorf("message_ids must hold 1 to %d ids", maxOrganize)
	}
	err = t.apply(ctx, who.userID, in)
	if err != nil {
		return nil, actionOutput{}, err
	}
	return nil, actionOutput{
		Status: mcpkit.StatusDone, Count: len(in.MessageIDs),
		Message: fmt.Sprintf("%s applied to %d messages.", in.Action, len(in.MessageIDs)),
	}, nil
}

// apply runs one organize action.
func (t *toolset) apply(ctx context.Context, userID int64, in organizeInput) error {
	ids := in.MessageIDs
	var err error
	switch in.Action {
	case "mark_read", "mark_unread":
		err = t.mail.SetReadState(ctx, userID, in.AccountID, ids, in.Action == "mark_read")
	case "star", "unstar":
		err = t.mail.SetStarState(ctx, userID, in.AccountID, ids, in.Action == "star")
	case "archive", "inbox", "trash":
		err = t.mail.MoveEmailsTo(ctx, userID, in.AccountID, ids, mail.MoveEmailsRequest{FolderType: moves[in.Action]})
	case "move":
		if in.FolderID == 0 {
			return errors.New("action move needs folder_id. List folders with courrier_list_accounts in detailed format")
		}
		err = t.mail.MoveEmailsTo(ctx, userID, in.AccountID, ids, mail.MoveEmailsRequest{FolderID: &in.FolderID})
	default:
		return errAction
	}
	return mcpkit.Failure(err, hintMessage)
}
