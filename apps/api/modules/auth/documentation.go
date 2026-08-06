package auth

import documentation "github.com/FacileStudio/Courrier/apps/api/internal/documentation"

var Documentation = documentation.Module{
	Name:        "auth",
	Description: "Authentication routes.",
	Routes: []documentation.Route{
		{
			Method:      "GET",
			Path:        "/auth/config",
			Summary:     "Describe the auth methods on offer",
			Description: "Returns sso_only and oidc_enabled, so the client knows whether to show the password form.",
			Auth:        "public",
		},
		{
			Method:       "POST",
			Path:         "/auth/register",
			Summary:      "Register a new user",
			Description:  "Creates a user account and returns an auth token.",
			Auth:         "public",
			RequestBody:  "RegisterRequest",
			ResponseBody: "AuthResponse",
			Errors: []documentation.Error{
				{Status: 400, Code: "invalid_argument", Description: "Invalid JSON body or invalid registration input."},
				{Status: 409, Code: "already_exists", Description: "A user with the same email already exists."},
				{Status: 500, Code: "internal", Description: "Unexpected server error."},
			},
		},
		{
			Method:       "POST",
			Path:         "/auth/login",
			Summary:      "Authenticate a user",
			Description:  "Authenticates credentials and returns an auth token.",
			Auth:         "public",
			RequestBody:  "LoginRequest",
			ResponseBody: "AuthResponse",
			Errors: []documentation.Error{
				{Status: 400, Code: "invalid_argument", Description: "Invalid JSON body or invalid login input."},
				{Status: 401, Code: "unauthenticated", Description: "Email or password is invalid."},
				{Status: 500, Code: "internal", Description: "Unexpected server error."},
			},
		},
		{
			Method:      "POST",
			Path:        "/auth/logout",
			Summary:     "End the session",
			Description: "Deletes the session row and clears the session cookie. Answers ok even without a valid session.",
			Auth:        "public",
		},
		{
			Method:      "GET",
			Path:        "/auth/resource-token",
			Summary:     "Mint a resource token",
			Description: "Returns a 5-minute HMAC token, used as ?token= by browser requests that cannot send an Authorization header, such as inline images.",
			Auth:        "bearer token required",
			Errors: []documentation.Error{
				{Status: 401, Code: "unauthenticated", Description: "No valid session cookie, bearer token or API token."},
			},
		},
	},
}
