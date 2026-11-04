package mail

import (
	"net/http"

	documentation "github.com/FacileStudio/Courrier/apps/api/internal/documentation"
)

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
	Routes: append([]documentation.Route{
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
			Auth:         "bearer token required",
			PathParams:   []documentation.Field{accountID},
			ResponseBody: CheckResult{},
			Errors:       []documentation.Error{unauthenticated, accountNotFound, rateLimited},
		},
		{
			Method:       "POST",
			Path:         "/accounts/{accountId}/mail/sync",
			Summary:      "Sync the account's folder list",
			Description:  "Connects over IMAP, lists the mailboxes and maps them to folder types. Rate limited to 5 per minute.",
			Auth:         "bearer token required",
			PathParams:   []documentation.Field{accountID},
			ResponseBody: SyncResponse{},
			Errors:       []documentation.Error{unauthenticated, accountNotFound, rateLimited},
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
			ResponseBody: SyncResponse{},
			Errors:       []documentation.Error{unauthenticated, accountNotFound, rateLimited},
		},
		{
			Method:       "GET",
			Path:         "/accounts/{accountId}/mail/folders",
			Summary:      "List folders",
			Description:  "Returns the account's synced folders with unread and total counts.",
			Auth:         "bearer token required",
			PathParams:   []documentation.Field{accountID},
			ResponseBody: FolderListResponse{},
			Errors:       []documentation.Error{unauthenticated, accountNotFound},
		},
		{
			Method:      "GET",
			Path:        "/accounts/{accountId}/mail/folders/{folderType}/emails",
			Summary:     "List a folder's conversations",
			Description: "Rows are collapsed conversations, not individual messages. Takes page (default 1), limit (default 50, capped at 100), unread=true and tab. A sender rule overrides the classifier's category.",
			Auth:        "bearer token required",
			PathParams: []documentation.Field{
				accountID,
				{Name: "folderType", Type: "string", Description: "Folder type: inbox, sent, drafts, trash, junk, archive or custom."},
			},
			QueryParams: []documentation.Field{
				{Name: "page", Type: "int", Description: "Page number."},
				{Name: "limit", Type: "int", Description: "Items per page."},
				{Name: "unread", Type: "bool", Description: "Unread emails only."},
				{Name: "tab", Type: "string", Description: "primary, newsletters, notifications, or all (default)."},
			},
			ResponseBody: ConversationListResponse{},
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
			ResponseBody: ThreadResponse{},
			Errors:       []documentation.Error{unauthenticated, accountNotFound},
		},
		{
			Method:      "GET",
			Path:        "/accounts/{accountId}/mail/search",
			Summary:     "Search messages",
			Description: "Runs q against the trigram indexes on subject, sender name, sender address and body text. Takes page and limit (default 30, capped at 100). A blank q returns an empty result.",
			Auth:        "bearer token required",
			PathParams:  []documentation.Field{accountID},
			QueryParams: []documentation.Field{
				{Name: "q", Type: "string", Description: "Search query."},
				{Name: "page", Type: "int", Description: "Page number."},
				{Name: "limit", Type: "int", Description: "Items per page."},
			},
			ResponseBody: SearchResponse{},
			Errors:       []documentation.Error{unauthenticated, accountNotFound},
		},
		{
			Method:       "GET",
			Path:         "/accounts/{accountId}/mail/emails/{emailId}",
			Summary:      "Get one message",
			Description:  "Returns a single message with its attachment metadata.",
			Auth:         "bearer token required",
			PathParams:   []documentation.Field{accountID, emailID},
			ResponseBody: EmailResponse{},
			Errors: []documentation.Error{
				unauthenticated,
				{Status: 404, Code: "not_found", Description: "No such message in this account."},
			},
		},
		{
			Method:      "GET",
			Path:        "/accounts/{accountId}/mail/contacts",
			Summary:     "Search seen contacts",
			Description: "With q, returns addresses seen in the account's mail with how often each appears (auto-complete). Without q, returns an empty contacts list; use GET /mail/contacts with account for the paginated list.",
			Auth:        "bearer token required",
			PathParams:  []documentation.Field{accountID},
			QueryParams: []documentation.Field{
				{Name: "q", Type: "string", Description: "Auto-complete query. Empty returns no contacts."},
			},
			ResponseBody: ContactListResponse{},
			Errors:       []documentation.Error{unauthenticated, accountNotFound},
		},
		{
			Method:      "GET",
			Path:        "/mail/contacts",
			Summary:     "List contacts across accounts",
			Description: "Returns the addresses the user has written to from any account they can see, merged by address, newest first. Rate limited to 60 per minute.",
			Auth:        "bearer token required",
			QueryParams: []documentation.Field{
				{Name: "account", Type: "int", Description: "Restrict to one account. 404 when the user cannot see it."},
				{Name: "page", Type: "int", Description: "Page number, capped at 10000."},
				{Name: "limit", Type: "int", Description: "Items per page, default 50, capped at 100."},
			},
			ResponseBody: ContactsResponse{},
			Errors:       []documentation.Error{unauthenticated},
		},
		{
			Method:       "POST",
			Path:         "/accounts/{accountId}/mail/contacts/backfill",
			Summary:      "Backfill contacts from the sent folder",
			Description:  "Reads up to the newest 1000 envelopes of the sent folder over IMAP, without bodies and without storing messages, to fill the contacts table. Safe to repeat. Shares the per-account guard of the mail backfill: a busy account returns running true. Rate limited to 2 per minute.",
			Auth:         "bearer token required",
			PathParams:   []documentation.Field{accountID},
			ResponseBody: BackfillContactsResponse{},
			Errors:       []documentation.Error{unauthenticated, accountNotFound, rateLimited},
		},
		{
			Method:      "GET",
			Path:        "/mail/triage",
			Summary:     "Classified inbox across accounts",
			Description: "Lists the inbox conversations of every account the user can see (personal accounts and the accounts of their spaces), newest first, each tagged known, unknown_human or bulk with a reason. Bulk signals win over known. counts covers the whole set matching unread, since and account, whatever the class filter, and always carries the three keys; total, items and paging follow the class filter. At most 5000 conversations are listed; truncated is true when more match, and counts and total then cover only the newest 5000. Rate limited to 60 per minute.",
			Auth:        "bearer token required",
			QueryParams: []documentation.Field{
				{Name: "unread", Type: "bool", Description: "true keeps conversations with an unread message."},
				{Name: "class", Type: "string", Description: "Comma list of known, unknown_human, bulk. Default all."},
				{Name: "since", Type: "string", Description: "RFC3339 time or a duration such as 30m, 2h, 1d, 1w, up to 3650 units."},
				{Name: "account", Type: "int", Description: "Restrict to one account. 404 when the user cannot see it."},
				{Name: "page", Type: "int", Description: "Page number, capped at 10000."},
				{Name: "limit", Type: "int", Description: "Items per page, default 50, capped at 100."},
			},
			ResponseBody: TriageResponse{},
			Errors: []documentation.Error{
				unauthenticated,
				{Status: 400, Code: "invalid_argument", Description: "Unknown class, bad since or bad account id."},
				accountNotFound,
			},
		},
		{
			Method:       "POST",
			Path:         "/accounts/{accountId}/mail/backfill",
			Summary:      "Classify cached mail and deep-sync sent",
			Description:  "Fetches the headers of cached inbox rows from the last 60 days that have no category yet, in batches of 50, then reads up to 1000 sent envelopes into contacts. Idempotent and resumable. Only one run per account is active; a second call returns running true. A sync also starts it in the background when unclassified rows exist. Rate limited to 2 per minute.",
			Auth:         "bearer token required",
			PathParams:   []documentation.Field{accountID},
			ResponseBody: BackfillResult{},
			Errors:       []documentation.Error{unauthenticated, accountNotFound, rateLimited},
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
			QueryParams: []documentation.Field{
				{Name: "token", Type: "string", Description: "5-minute resource HMAC token."},
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
			RequestBody:  UpdateEmailRequest{},
			ResponseBody: EmailResponse{},
			Errors: []documentation.Error{
				{Status: 400, Code: "invalid_argument", Description: "Invalid JSON body."},
				unauthenticated,
				{Status: 404, Code: "not_found", Description: "No such message in this account."},
			},
		},
		{
			Method:       "POST",
			Path:         "/accounts/{accountId}/mail/send",
			Summary:      "Send a message",
			Description:  "Sends over SMTP and appends the result to the account's Sent folder. Accepts JSON or multipart/form-data with repeated attachments fields, capped at 25 MB; in multipart, to, cc, bcc and references are comma-separated. At least one of to, cc and bcc is required. Bcc recipients are added to the SMTP envelope only: the delivered message carries no Bcc header, while the copy in Sent and the cached row record it as bcc_addresses. Rate limited to 10 per minute.",
			Auth:         "bearer token required",
			PathParams:   []documentation.Field{accountID},
			RequestBody:  SendRequest{},
			ResponseBody: SendResponse{},
			Errors: []documentation.Error{
				{Status: 400, Code: "invalid_argument", Description: "Invalid body or missing recipient."},
				unauthenticated,
				accountNotFound,
				{Status: 413, Code: "resource_exhausted", Description: "Attachments exceed 25 MB."},
				rateLimited,
			},
		},
		{
			Method:       "POST",
			Path:         "/accounts/{accountId}/mail/drafts",
			Summary:      "Save a draft",
			Description:  "Stores the message and appends it to the account's Drafts folder. Returns the new message id. Accepts JSON or multipart/form-data with repeated attachments fields, capped at 25 MB; the Drafts copy keeps the Bcc header. A draft with attachments must reach the IMAP server, which holds the attachment bytes, and fails with 412 otherwise.",
			Auth:         "bearer token required",
			PathParams:   []documentation.Field{accountID},
			RequestBody:  SendRequest{},
			ResponseBody: DraftResponse{},
			Errors: []documentation.Error{
				{Status: 400, Code: "invalid_argument", Description: "Invalid body."},
				unauthenticated,
				accountNotFound,
				{Status: 412, Code: "failed_precondition", Description: "The draft has attachments and the IMAP server could not store it."},
				{Status: 413, Code: "resource_exhausted", Description: "Attachments exceed 25 MB."},
			},
		},
		{
			Method:       "DELETE",
			Path:         "/accounts/{accountId}/mail/drafts/{emailId}",
			Summary:      "Delete a draft",
			Description:  "Removes the draft locally and from the account's Drafts folder.",
			Auth:         "bearer token required",
			PathParams:   []documentation.Field{accountID, emailID},
			ResponseBody: DeleteDraftResponse{},
			Errors: []documentation.Error{
				unauthenticated,
				{Status: 404, Code: "not_found", Description: "No such draft."},
			},
		},
		{
			Method:       "POST",
			Path:         "/accounts/{accountId}/mail/emails/bulk-action",
			Summary:      "Act on several messages",
			Description:  "Applies delete, archive, mark_read, mark_unread, star or unstar to up to 200 message ids. Deleting a message already in Trash removes it from the server. Rate limited to 30 per minute.",
			Auth:         "bearer token required",
			PathParams:   []documentation.Field{accountID},
			RequestBody:  BulkActionRequest{},
			ResponseBody: BulkActionResponse{},
			Errors: []documentation.Error{
				{Status: 400, Code: "invalid_argument", Description: "Empty email_ids, more than 200 ids, or an unknown action."},
				unauthenticated,
				accountNotFound,
				rateLimited,
			},
		},
		{
			Method:       "POST",
			Path:         "/mail/test-connection",
			Summary:      "Test IMAP and SMTP credentials",
			Description:  "Dials both servers for real, so an error here is the honest report on whether an account will work. Registered outside /accounts because it exists to validate credentials before an account is created. Rate limited to 5 per minute.",
			Auth:         "bearer token required",
			RequestBody:  TestConnectionRequest{},
			ResponseBody: TestConnectionResponse{},
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
			QueryParams:  []documentation.Field{{Name: "space_id", Type: "string", Description: "Space identifier."}},
			ResponseBody: TemplateListResponse{},
			Errors:       []documentation.Error{unauthenticated, rateLimited},
		},
		{
			Method:       "POST",
			Path:         "/templates",
			Summary:      "Create a message template",
			Description:  "Stores a reusable message, optionally scoped to a space.",
			Auth:         "bearer token required",
			RequestBody:  EmailTemplateRequest{},
			ResponseBody: EmailTemplateResponse{},
			Status:       http.StatusCreated,
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
			RequestBody:  EmailTemplateRequest{},
			ResponseBody: EmailTemplateResponse{},
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
			ResponseBody: DeleteTemplateResponse{},
			Errors: []documentation.Error{
				unauthenticated,
				{Status: 404, Code: "not_found", Description: "No such template."},
				rateLimited,
			},
		},
	}, writeRoutes...),
}
