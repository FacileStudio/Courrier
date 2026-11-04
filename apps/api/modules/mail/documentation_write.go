package mail

import (
	documentation "github.com/FacileStudio/Courrier/apps/api/internal/documentation"
)

var folderID = documentation.Field{Name: "folderId", Type: "int", Description: "Folder identifier."}

var writeRoutes = []documentation.Route{
	{
		Method:  "PUT",
		Path:    "/accounts/{accountId}/mail/sender-rules",
		Summary: "Pin a sender to an inbox tab",
		Description: "Every message from address in this account lists under tab (primary, newsletters, " +
			"notifications) from now on, whatever its headers say. tab auto removes the pin. Returns 204.",
		Auth:        "bearer token required",
		PathParams:  []documentation.Field{accountID},
		RequestBody: SenderRuleRequest{},
		Errors: []documentation.Error{
			{Status: 400, Code: "invalid_argument", Description: "Not an email address, or an unknown tab."},
			unauthenticated,
			{Status: 404, Code: "not_found", Description: "No such account."},
		},
	},
	{
		Method:  "PUT",
		Path:    "/accounts/{accountId}/mail/drafts/{emailId}",
		Summary: "Update a draft in place",
		Description: "Replaces the draft's recipients, subject and bodies with the request's and keeps its id " +
			"and Message-ID, so references to it stay valid. The Drafts copy on the IMAP server is replaced. " +
			"Existing attachments are kept unless their ids are listed in remove_attachments; files sent in " +
			"the attachments field of a multipart body are added (remove_attachments is then a comma-separated " +
			"string). Returns the same message id.",
		Auth:         "bearer token required",
		PathParams:   []documentation.Field{accountID, emailID},
		RequestBody:  DraftRequest{},
		ResponseBody: DraftResponse{},
		Errors: []documentation.Error{
			{Status: 400, Code: "invalid_argument", Description: "Invalid body."},
			unauthenticated,
			{Status: 404, Code: "not_found", Description: "No such draft."},
			{Status: 412, Code: "failed_precondition", Description: "The IMAP server could not store the new copy; the draft is unchanged."},
			{Status: 413, Code: "resource_exhausted", Description: "Attachments exceed 25 MB."},
		},
	},
	{
		Method:  "POST",
		Path:    "/accounts/{accountId}/mail/emails/move",
		Summary: "Move messages to a folder",
		Description: "Moves up to 200 messages into the folder named by folder_id, or by folder_type for a " +
			"system folder (inbox, sent, drafts, trash, junk, archive). The move runs on the IMAP server " +
			"first and the cached rows follow, keeping their ids and taking the uid the server assigns. " +
			"Messages already in the destination are left alone. Rate limited to 30 per minute.",
		Auth:         "bearer token required",
		PathParams:   []documentation.Field{accountID},
		RequestBody:  MoveEmailsRequest{},
		ResponseBody: BulkActionResponse{},
		Errors: []documentation.Error{
			{Status: 400, Code: "invalid_argument", Description: "Empty email_ids, more than 200 ids, or no usable destination."},
			unauthenticated,
			{Status: 404, Code: "not_found", Description: "No such account, destination folder or messages."},
			{Status: 412, Code: "failed_precondition", Description: "The IMAP server refused the move."},
			rateLimited,
		},
	},
	{
		Method:  "POST",
		Path:    "/accounts/{accountId}/mail/folders",
		Summary: "Create a folder",
		Description: "Creates a top-level mailbox on the IMAP server and caches it as a custom folder. The " +
			"name is 1 to 100 characters, without the server's hierarchy delimiter, / \\ * %, control " +
			"characters or a leading dot, and must not read as a system folder (inbox, sent, drafts, " +
			"trash, junk, archive): the next sync would file it under that type. Rate limited to 10 per minute.",
		Auth:         "bearer token required",
		PathParams:   []documentation.Field{accountID},
		RequestBody:  FolderNameRequest{},
		ResponseBody: FolderResponse{},
		Errors: []documentation.Error{
			{Status: 400, Code: "invalid_argument", Description: "Invalid or reserved name."},
			unauthenticated,
			accountNotFound,
			{Status: 409, Code: "already_exists", Description: "A folder with that name exists."},
			{Status: 412, Code: "failed_precondition", Description: "The IMAP server refused to create it."},
			rateLimited,
		},
	},
	{
		Method:  "PATCH",
		Path:    "/accounts/{accountId}/mail/folders/{folderId}",
		Summary: "Rename a folder",
		Description: "Renames a custom folder on the IMAP server and in the cache; the paths of its " +
			"descendants follow. System folders are refused with 400. Rate limited to 10 per minute.",
		Auth:         "bearer token required",
		PathParams:   []documentation.Field{accountID, folderID},
		RequestBody:  FolderNameRequest{},
		ResponseBody: FolderResponse{},
		Errors: []documentation.Error{
			{Status: 400, Code: "invalid_argument", Description: "Invalid or reserved name, or a system folder."},
			unauthenticated,
			{Status: 404, Code: "not_found", Description: "No such folder in this account."},
			{Status: 409, Code: "already_exists", Description: "A folder with that name exists."},
			{Status: 412, Code: "failed_precondition", Description: "The IMAP server refused the rename."},
			rateLimited,
		},
	},
	{
		Method:  "DELETE",
		Path:    "/accounts/{accountId}/mail/folders/{folderId}",
		Summary: "Delete a folder",
		Description: "Deletes a custom folder from the IMAP server and the cache. IMAP deletes the messages " +
			"inside with it, so a folder that holds messages answers 409 unless the force=true query " +
			"parameter is set. System folders are refused with 400. Rate limited to 10 per minute.",
		Auth:       "bearer token required",
		PathParams: []documentation.Field{accountID, folderID},
		QueryParams: []documentation.Field{
			{Name: "force", Type: "bool", Description: "Delete the folder even when it holds messages."},
		},
		ResponseBody: DeleteFolderResponse{},
		Errors: []documentation.Error{
			{Status: 400, Code: "invalid_argument", Description: "A system folder."},
			unauthenticated,
			{Status: 404, Code: "not_found", Description: "No such folder in this account."},
			{Status: 409, Code: "already_exists", Description: "The folder holds messages and force is not set."},
			{Status: 412, Code: "failed_precondition", Description: "The IMAP server refused the delete."},
			rateLimited,
		},
	},
}
