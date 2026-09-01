package accounts

import (
	"net/http"

	documentation "github.com/FacileStudio/Courrier/apps/api/internal/documentation"
)

var accountID = []documentation.Field{
	{Name: "id", Type: "int", Description: "Mail account identifier."},
}

var Documentation = documentation.Module{
	Name:        "accounts",
	Description: "Mail accounts — the IMAP/SMTP connections a user has configured. Passwords are stored AES-GCM encrypted and never returned.",
	Routes: []documentation.Route{
		{
			Method:       "GET",
			Path:         "/accounts",
			Summary:      "List mail accounts",
			Description:  "Returns the caller's mail accounts. Accepts a space_id query parameter to scope the list to one space.",
			Auth:         "bearer token required",
			QueryParams:  []documentation.Field{{Name: "space_id", Type: "string", Description: "Space identifier."}},
			ResponseBody: AccountListResponse{},
			Errors: []documentation.Error{
				{Status: 401, Code: "unauthenticated", Description: "No valid session cookie, bearer token or API token."},
				{Status: 500, Code: "internal", Description: "Unexpected server error."},
			},
		},
		{
			Method:       "POST",
			Path:         "/accounts",
			Summary:      "Create a mail account",
			Description:  "Stores IMAP and SMTP credentials. IMAP defaults to port 993 and SMTP to 587.",
			Auth:         "bearer token required",
			RequestBody:  CreateAccountRequest{},
			ResponseBody: AccountResponse{},
			Status:       http.StatusCreated,
			Errors: []documentation.Error{
				{Status: 400, Code: "invalid_argument", Description: "Invalid JSON body or missing connection details."},
				{Status: 401, Code: "unauthenticated", Description: "No valid session cookie, bearer token or API token."},
				{Status: 500, Code: "internal", Description: "Unexpected server error."},
			},
		},
		{
			Method:       "GET",
			Path:         "/accounts/{id}",
			Summary:      "Get a mail account",
			Description:  "Returns one mail account, without either stored password.",
			Auth:         "bearer token required",
			PathParams:   accountID,
			ResponseBody: AccountResponse{},
			Errors: []documentation.Error{
				{Status: 400, Code: "invalid_argument", Description: "Account id is not an integer."},
				{Status: 401, Code: "unauthenticated", Description: "No valid session cookie, bearer token or API token."},
				{Status: 404, Code: "not_found", Description: "No such account for this user."},
			},
		},
		{
			Method:       "PUT",
			Path:         "/accounts/{id}",
			Summary:      "Update a mail account",
			Description:  "Updates connection details. Every field is optional; passwords are only re-encrypted when supplied.",
			Auth:         "bearer token required",
			PathParams:   accountID,
			RequestBody:  UpdateAccountRequest{},
			ResponseBody: AccountResponse{},
			Errors: []documentation.Error{
				{Status: 400, Code: "invalid_argument", Description: "Invalid JSON body or account id is not an integer."},
				{Status: 401, Code: "unauthenticated", Description: "No valid session cookie, bearer token or API token."},
				{Status: 404, Code: "not_found", Description: "No such account for this user."},
			},
		},
		{
			Method:       "DELETE",
			Path:         "/accounts/{id}",
			Summary:      "Delete a mail account",
			Description:  "Removes the account together with its synced folders and messages.",
			Auth:         "bearer token required",
			PathParams:   accountID,
			ResponseBody: DeleteAccountResponse{},
			Errors: []documentation.Error{
				{Status: 400, Code: "invalid_argument", Description: "Account id is not an integer."},
				{Status: 401, Code: "unauthenticated", Description: "No valid session cookie, bearer token or API token."},
				{Status: 404, Code: "not_found", Description: "No such account for this user."},
			},
		},
	},
}
