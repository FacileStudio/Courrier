package mail

import documentation "github.com/FacileStudio/Courrier/apps/api/internal/documentation"

var (
	accountID = documentation.Field{Name: "accountId", Type: "int", Description: "Mail account identifier."}
	emailID   = documentation.Field{Name: "emailId", Type: "int", Description: "Stored message identifier."}

	unauthenticated = documentation.Error{Status: 401, Code: "unauthenticated", Description: "No valid session cookie, bearer token or API token."}
	accountNotFound = documentation.Error{Status: 404, Code: "not_found", Description: "No such account for this user."}
	rateLimited     = documentation.Error{Status: 429, Code: "resource_exhausted", Description: "Rate limit exceeded."}
)

var Documentation = documentation.Module{
	Name:        "mail",
	Description: "Mailbox synchronisation, reading, sending, drafts and reusable templates. Every mail route is scoped to one account.",
	Routes: []documentation.Route{
		{
			Method:  "POST",
			Path:    "/accounts/{accountId}/mail/check",
			Summary: "Check the account's stored credentials",
			Description: "Signs in to IMAP and to SMTP with the stored passwords and disconnects again — " +
				"no mailbox listing, no message fetch, nothing written. Both protocols are always " +
				"attempted and reported separately, so \"IMAP works, SMTP does not\" is expressible. " +
				"A refused or rejected handshake is still a 200: the body carries `configured`, `ok` " +
				"and `error` per protocol. Rate limited to 10 per minute. Unlike test-connection this " +
				"needs no passwords in the request, which is why it is the only probe usable on a " +
				"saved account.",
			Auth:       "bearer token required",
			PathParams: []documentation.Field{accountID},
			Errors:     []documentation.Error{unauthenticated, accountNotFound, rateLimited},
		},
		{
			Method:      "POST",
			Path:        "/accounts/{accountId}/mail/sync",
			Summary:     "Sync the account's folder list",
			Description: "Connects over IMAP, lists the mailboxes and maps them to folder types. Rate limited to 5 per minute.",
			Auth:        "bearer token required",
			PathParams:  []documentation.Field{accountID},
			Errors:      []documentation.Error{unauthenticated, accountNotFound, rateLimited},
		},
		{
			Method:      "POST",
			Path:        "/accounts/{accountId}/mail/folders/{folderId}/sync",
			Summary:     "Sync one folder's messages",
			Description: "Fetches messages for a single folder over IMAP. Rate limited to 10 per minute.",
			Auth:        "bearer token required",
			PathParams: []documentation.Field{
				accountID,
				{Name: "folderId", Type: "int", Description: "Folder identifier."},
			},
			Errors: []documentation.Error{unauthenticated, accountNotFound, rateLimited},
		},
		{
			Method:       "GET",
			Path:         "/accounts/{accountId}/mail/folders",
			Summary:      "List folders",
			Description:  "Returns the account's synced folders with unread and total counts.",
			Auth:         "bearer token required",
			PathParams:   []documentation.Field{accountID},
			ResponseBody: "[]FolderResponse",
			Errors:       []documentation.Error{unauthenticated, accountNotFound},
		},
		{
			Method:      "GET",
			Path:        "/accounts/{accountId}/mail/folders/{folderType}/emails",
			Summary:     "List a folder's conversations",
			Description: "Rows are collapsed conversations, not individual messages. Takes page (default 1), limit (default 50, capped at 100) and unread=true.",
			Auth:        "bearer token required",
			PathParams: []documentation.Field{
				accountID,
				{Name: "folderType", Type: "string", Description: "Folder type: inbox, sent, drafts, trash, junk, archive or custom."},
			},
			ResponseBody: "[]EmailResponse",
			Errors:       []documentation.Error{unauthenticated, accountNotFound},
		},
		{
			Method:      "GET",
			Path:        "/accounts/{accountId}/mail/threads/{threadId}",
			Summary:     "Get a conversation",
			Description: "Returns every message in one thread. The thread id is a Message-ID and is percent-decoded before use.",
			Auth:        "bearer token required",
			PathParams: []documentation.Field{
				accountID,
				{Name: "threadId", Type: "string", Description: "Percent-encoded Message-ID of the conversation."},
			},
			ResponseBody: "[]EmailResponse",
			Errors:       []documentation.Error{unauthenticated, accountNotFound},
		},
		{
			Method:       "GET",
			Path:         "/accounts/{accountId}/mail/search",
			Summary:      "Search messages",
			Description:  "Runs q against the trigram indexes on subject, sender name, sender address and body text. Takes page and limit (default 30, capped at 100). A blank q returns an empty result.",
			Auth:         "bearer token required",
			PathParams:   []documentation.Field{accountID},
			ResponseBody: "[]EmailResponse",
			Errors:       []documentation.Error{unauthenticated, accountNotFound},
		},
		{
			Method:       "GET",
			Path:         "/accounts/{accountId}/mail/emails/{emailId}",
			Summary:      "Get one message",
			Description:  "Returns a single message with its attachment metadata.",
			Auth:         "bearer token required",
			PathParams:   []documentation.Field{accountID, emailID},
			ResponseBody: "EmailResponse",
			Errors: []documentation.Error{
				unauthenticated,
				{Status: 404, Code: "not_found", Description: "No such message in this account."},
			},
		},
		{
			Method:       "GET",
			Path:         "/accounts/{accountId}/mail/contacts",
			Summary:      "List seen contacts",
			Description:  "Returns addresses seen in the account's mail with how often each appears. Takes q; a blank q returns an empty result.",
			Auth:         "bearer token required",
			PathParams:   []documentation.Field{accountID},
			ResponseBody: "[]ContactResult",
			Errors:       []documentation.Error{unauthenticated, accountNotFound},
		},
		{
			Method:      "GET",
			Path:        "/accounts/{accountId}/mail/emails/{emailId}/attachments/{attachmentId}/download",
			Summary:     "Download an attachment",
			Description: "Streams the attachment, fetched from IMAP on demand rather than stored.",
			Auth:        "bearer token required",
			PathParams: []documentation.Field{
				accountID,
				emailID,
				{Name: "attachmentId", Type: "int", Description: "Attachment identifier."},
			},
			Errors: []documentation.Error{
				unauthenticated,
				{Status: 404, Code: "not_found", Description: "No such attachment."},
			},
		},
		{
			Method:      "GET",
			Path:        "/accounts/{accountId}/mail/emails/{emailId}/cid/{cid}",
			Summary:     "Fetch an inline image",
			Description: "Serves a cid: part referenced by an HTML body. A browser sends no Authorization header for an <img>, so this route also accepts ?token= carrying a resource token from GET /auth/resource-token.",
			Auth:        "resource token or bearer token",
			PathParams: []documentation.Field{
				accountID,
				emailID,
				{Name: "cid", Type: "string", Description: "Percent-encoded Content-ID of the inline part."},
			},
			Errors: []documentation.Error{
				unauthenticated,
				{Status: 404, Code: "not_found", Description: "No part with that Content-ID."},
			},
		},
		{
			Method:       "PATCH",
			Path:         "/accounts/{accountId}/mail/emails/{emailId}",
			Summary:      "Flag a message",
			Description:  "Sets is_read and/or is_starred; both fields are optional.",
			Auth:         "bearer token required",
			PathParams:   []documentation.Field{accountID, emailID},
			RequestBody:  "UpdateEmailRequest",
			ResponseBody: "EmailResponse",
			Errors: []documentation.Error{
				{Status: 400, Code: "invalid_argument", Description: "Invalid JSON body."},
				unauthenticated,
				{Status: 404, Code: "not_found", Description: "No such message in this account."},
			},
		},
		{
			Method:      "POST",
			Path:        "/accounts/{accountId}/mail/send",
			Summary:     "Send a message",
			Description: "Sends over SMTP and appends the result to the account's Sent folder. Accepts JSON or multipart/form-data with repeated attachments fields, capped at 25 MB. Rate limited to 10 per minute.",
			Auth:        "bearer token required",
			PathParams:  []documentation.Field{accountID},
			RequestBody: "SendRequest",
			Errors: []documentation.Error{
				{Status: 400, Code: "invalid_argument", Description: "Invalid body or missing recipient."},
				unauthenticated,
				accountNotFound,
				{Status: 413, Code: "resource_exhausted", Description: "Attachments exceed 25 MB."},
				rateLimited,
			},
		},
		{
			Method:      "POST",
			Path:        "/accounts/{accountId}/mail/drafts",
			Summary:     "Save a draft",
			Description: "Stores the message and appends it to the account's Drafts folder. Returns the new message id.",
			Auth:        "bearer token required",
			PathParams:  []documentation.Field{accountID},
			RequestBody: "SendRequest",
			Errors: []documentation.Error{
				{Status: 400, Code: "invalid_argument", Description: "Invalid body."},
				unauthenticated,
				accountNotFound,
			},
		},
		{
			Method:      "DELETE",
			Path:        "/accounts/{accountId}/mail/drafts/{emailId}",
			Summary:     "Delete a draft",
			Description: "Removes the draft locally and from the account's Drafts folder.",
			Auth:        "bearer token required",
			PathParams:  []documentation.Field{accountID, emailID},
			Errors: []documentation.Error{
				unauthenticated,
				{Status: 404, Code: "not_found", Description: "No such draft."},
			},
		},
		{
			Method:      "POST",
			Path:        "/accounts/{accountId}/mail/emails/bulk-action",
			Summary:     "Act on several messages",
			Description: "Applies delete, archive, mark_read or mark_unread to up to 200 message ids. Rate limited to 30 per minute.",
			Auth:        "bearer token required",
			PathParams:  []documentation.Field{accountID},
			RequestBody: "BulkActionRequest",
			Errors: []documentation.Error{
				{Status: 400, Code: "invalid_argument", Description: "Empty email_ids, more than 200 ids, or an unknown action."},
				unauthenticated,
				accountNotFound,
				rateLimited,
			},
		},
		{
			Method:      "POST",
			Path:        "/mail/test-connection",
			Summary:     "Test IMAP and SMTP credentials",
			Description: "Dials both servers for real, so an error here is the honest report on whether an account will work. Registered outside /accounts because it exists to validate credentials before an account is created. Rate limited to 5 per minute.",
			Auth:        "bearer token required",
			RequestBody: "TestConnectionRequest",
			Errors: []documentation.Error{
				{Status: 400, Code: "invalid_argument", Description: "Invalid JSON body, or the connection was refused."},
				unauthenticated,
				rateLimited,
			},
		},
		{
			Method:       "GET",
			Path:         "/templates",
			Summary:      "List message templates",
			Description:  "Returns the caller's reusable messages. Accepts a space_id query parameter.",
			Auth:         "bearer token required",
			ResponseBody: "[]EmailTemplateResponse",
			Errors:       []documentation.Error{unauthenticated, rateLimited},
		},
		{
			Method:       "POST",
			Path:         "/templates",
			Summary:      "Create a message template",
			Description:  "Stores a reusable message, optionally scoped to a space.",
			Auth:         "bearer token required",
			RequestBody:  "EmailTemplateRequest",
			ResponseBody: "EmailTemplateResponse",
			Errors: []documentation.Error{
				{Status: 400, Code: "invalid_argument", Description: "Invalid JSON body or missing name."},
				unauthenticated,
				rateLimited,
			},
		},
		{
			Method:      "PUT",
			Path:        "/templates/{templateId}",
			Summary:     "Update a message template",
			Description: "Replaces a template's name, subject and bodies.",
			Auth:        "bearer token required",
			PathParams: []documentation.Field{
				{Name: "templateId", Type: "int", Description: "Template identifier."},
			},
			RequestBody:  "EmailTemplateRequest",
			ResponseBody: "EmailTemplateResponse",
			Errors: []documentation.Error{
				{Status: 400, Code: "invalid_argument", Description: "Invalid JSON body."},
				unauthenticated,
				{Status: 404, Code: "not_found", Description: "No such template."},
				rateLimited,
			},
		},
		{
			Method:      "DELETE",
			Path:        "/templates/{templateId}",
			Summary:     "Delete a message template",
			Description: "Removes one reusable message.",
			Auth:        "bearer token required",
			PathParams: []documentation.Field{
				{Name: "templateId", Type: "int", Description: "Template identifier."},
			},
			Errors: []documentation.Error{
				unauthenticated,
				{Status: 404, Code: "not_found", Description: "No such template."},
				rateLimited,
			},
		},
	},
}
