package auth

import (
	"net/http"

	documentation "github.com/FacileStudio/Courrier/apps/api/internal/documentation"
	"github.com/FacileStudio/porte"
)

var Documentation = documentation.Module{
	Name:        "auth",
	Description: "Authentication routes.",
	Routes: []documentation.Route{
		{
			Method:       "GET",
			Path:         "/auth/config",
			Summary:      "Describe the auth methods on offer",
			Description:  "Returns sso_only and oidc_enabled, so the client knows whether to show the password form.",
			Auth:         "public",
			ResponseBody: porte.ConfigResponse{},
		},
		{
			Method:       "POST",
			Path:         "/auth/register",
			Summary:      "Register a new user",
			Description:  "Creates a user account and returns an auth token.",
			Auth:         "public",
			RequestBody:  RegisterRequest{},
			ResponseBody: AuthResponse{},
			Status:       http.StatusCreated,
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
			RequestBody:  LoginRequest{},
			ResponseBody: AuthResponse{},
			Errors: []documentation.Error{
				{Status: 400, Code: "invalid_argument", Description: "Invalid JSON body or invalid login input."},
				{Status: 401, Code: "unauthenticated", Description: "Email or password is invalid."},
				{Status: 500, Code: "internal", Description: "Unexpected server error."},
			},
		},
		{
			Method:       "POST",
			Path:         "/auth/logout",
			Summary:      "End the session",
			Description:  "Deletes the session row and clears the session cookie. Answers ok even without a valid session.",
			Auth:         "public",
			ResponseBody: porte.LogoutResponse{},
		},
		{
			Method:      "GET",
			Path:        "/auth/oidc",
			Summary:     "Start the SSO login",
			Description: "Redirects to the identity provider with PKCE and a nonce. Registered only when OIDC_ISSUER is set. Add ?flow=cli for the CLI handoff.",
			Auth:        "public",
		},
		{
			Method:      "GET",
			Path:        "/auth/oidc/callback",
			Summary:     "Complete the SSO login",
			Description: "Verifies the ID token, upserts the account and sets the session cookie. A failure is a redirect to the login page carrying ?error=, not a JSON body — a refusal and a success are both a 302, so read Location and not the status.",
			Auth:        "public",
		},
		{
			Method:       "POST",
			Path:         "/auth/oidc/exchange",
			Summary:      "Trade a CLI login code for a session token",
			Description:  "The other half of ?flow=cli. The code is single use, hashed at rest and valid for 60 seconds.",
			Auth:         "public",
			RequestBody:  porte.ExchangeRequest{},
			ResponseBody: porte.ExchangeResponse{},
		},
		{
			Method:       "POST",
			Path:         "/auth/sync-profile",
			Summary:      "Refresh the profile from the provider",
			Description:  "Calls UserInfo with the stored refresh token and updates the name and photo. Rate-limited server-side.",
			Auth:         "session cookie or bearer token required",
			ResponseBody: porte.SyncProfileResponse{},
		},
		{
			Method:       "POST",
			Path:         "/auth/backchannel-logout",
			Summary:      "Revoke sessions on the provider's behalf",
			Description:  "Called by the identity provider, not the client. Validates the logout token and deletes that user's sessions.",
			Auth:         "signed logout token",
			ResponseBody: porte.LogoutResponse{},
		},
		{
			Method:       "GET",
			Path:         "/auth/resource-token",
			Summary:      "Mint a resource token",
			Description:  "Returns a 5-minute HMAC token, used as ?token= by browser requests that cannot send an Authorization header, such as inline images.",
			Auth:         "bearer token required",
			ResponseBody: ResourceTokenResponse{},
			Errors: []documentation.Error{
				{Status: 401, Code: "unauthenticated", Description: "No valid session cookie, bearer token or API token."},
			},
		},
	},
}
