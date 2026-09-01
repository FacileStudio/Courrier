package users

import (
	"net/http"

	documentation "github.com/FacileStudio/Courrier/apps/api/internal/documentation"
)

var Documentation = documentation.Module{
	Name:        "users",
	Description: "User listing plus current-user retrieval and update routes.",
	Routes: []documentation.Route{
		{
			Method:       "GET",
			Path:         "/users",
			Summary:      "List users",
			Description:  "Returns all authenticated users with profile metadata.",
			Auth:         "bearer token required",
			ResponseBody: ListResponse{},
			Errors: []documentation.Error{
				{Status: 401, Code: "unauthenticated", Description: "Authorization header is missing or invalid."},
				{Status: 500, Code: "internal", Description: "Unexpected server error."},
			},
		},
		{
			Method:       "GET",
			Path:         "/users/me",
			Summary:      "Return the current user",
			Description:  "Returns the authenticated user with profile metadata.",
			Auth:         "bearer token required",
			ResponseBody: MeResponse{},
			Errors: []documentation.Error{
				{Status: 401, Code: "unauthenticated", Description: "Authorization header is missing or invalid."},
				{Status: 500, Code: "internal", Description: "Unexpected server error."},
			},
		},
		{
			Method:       "PATCH",
			Path:         "/users/me",
			Summary:      "Update the current user",
			Description:  "Updates the authenticated user's name, email, and/or password. Changing a password that already exists requires current_password; without it the request is refused. A successful change ends the account's other logins, keeps named API tokens, and rotates this session's cookie.",
			Auth:         "bearer token required",
			RequestBody:  UpdateRequest{},
			ResponseBody: MeResponse{},
			Errors: []documentation.Error{
				{Status: 400, Code: "invalid_argument", Description: "Invalid JSON body, invalid update input, a password too short, a missing current_password on an account that has one, or no password to change."},
				{Status: 401, Code: "unauthenticated", Description: "Authorization header is missing or invalid, or current_password is wrong."},
				{Status: 404, Code: "not_found", Description: "The authenticated user no longer exists."},
				{Status: 409, Code: "already_exists", Description: "A user with the same email already exists."},
				{Status: 500, Code: "internal", Description: "Unexpected server error."},
			},
		},
		{
			Method:      "GET",
			Path:        "/users/{id}",
			Summary:     "Return one user",
			Description: "Returns a single user's profile metadata by ID.",
			Auth:        "bearer token required",
			PathParams: []documentation.Field{
				{Name: "id", Type: "string", Description: "User ID."},
			},
			ResponseBody: MeResponse{},
			Errors: []documentation.Error{
				{Status: 401, Code: "unauthenticated", Description: "Authorization header is missing or invalid."},
				{Status: 404, Code: "not_found", Description: "No user with that ID."},
				{Status: 500, Code: "internal", Description: "Unexpected server error."},
			},
		},
		{
			Method:       "POST",
			Path:         "/users/me/avatar",
			Summary:      "Upload the current user's avatar",
			Description:  "Stores a new avatar file for the authenticated user and returns the updated profile. Uploading is the fallback for users with no photo in SSO; it is rejected while one is set there.",
			Auth:         "bearer token required",
			ResponseBody: MeResponse{},
			Errors: []documentation.Error{
				{Status: 400, Code: "invalid_argument", Description: "Missing file, unsupported image type, or the photo is managed in SSO."},
				{Status: 401, Code: "unauthenticated", Description: "Authorization header is missing or invalid."},
				{Status: 404, Code: "not_found", Description: "The authenticated user no longer exists."},
				{Status: 413, Code: "resource_exhausted", Description: "Avatar file is too large."},
				{Status: 500, Code: "internal", Description: "Unexpected server error."},
			},
		},
		{
			Method:       "DELETE",
			Path:         "/users/me/avatar",
			Summary:      "Remove the current user's avatar",
			Description:  "Deletes the uploaded avatar file and returns the updated profile. A photo coming from SSO is not touched.",
			Auth:         "bearer token required",
			ResponseBody: MeResponse{},
			Errors: []documentation.Error{
				{Status: 401, Code: "unauthenticated", Description: "Authorization header is missing or invalid."},
				{Status: 404, Code: "not_found", Description: "The authenticated user no longer exists."},
				{Status: 500, Code: "internal", Description: "Unexpected server error."},
			},
		},
		{
			Method:       "GET",
			Path:         "/users/me/api-token",
			Summary:      "Report whether an API token exists",
			Description:  "Returns the current token's metadata without ever returning the token itself.",
			Auth:         "bearer token required",
			ResponseBody: ApiTokenStatusResponse{},
			Errors: []documentation.Error{
				{Status: 401, Code: "unauthenticated", Description: "Authorization header is missing or invalid."},
				{Status: 500, Code: "internal", Description: "Unexpected server error."},
			},
		},
		{
			Method:       "POST",
			Path:         "/users/me/api-token",
			Summary:      "Create an API token",
			Description:  "Issues a new API token and returns it once — it cannot be read back afterwards.",
			Auth:         "bearer token required",
			RequestBody:  CreateApiTokenRequest{},
			ResponseBody: ApiTokenResponse{},
			Status:       http.StatusCreated,
			Errors: []documentation.Error{
				{Status: 400, Code: "invalid_argument", Description: "Invalid JSON body or missing token name."},
				{Status: 401, Code: "unauthenticated", Description: "Authorization header is missing or invalid."},
				{Status: 500, Code: "internal", Description: "Unexpected server error."},
			},
		},
		{
			Method:       "DELETE",
			Path:         "/users/me/api-token",
			Summary:      "Revoke the current API token",
			Description:  "Deletes the authenticated user's API token.",
			Auth:         "bearer token required",
			ResponseBody: DeleteApiTokenResponse{},
			Errors: []documentation.Error{
				{Status: 401, Code: "unauthenticated", Description: "Authorization header is missing or invalid."},
				{Status: 500, Code: "internal", Description: "Unexpected server error."},
			},
		},
	},
}
