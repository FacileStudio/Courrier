package users

// User is the profile shape returned to the client.
type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	AvatarURL    string `json:"avatar_url"`
	AvatarSource string `json:"avatar_source"`
	CreatedAt    string `json:"created_at"`
}

// MeResponse wraps the signed-in user's own profile.
type MeResponse struct {
	User User `json:"user"`
}

// ListResponse wraps a page of users.
type ListResponse struct {
	Users []User `json:"users"`
}

// UpdateRequest is the body of PATCH /users/me; every field is optional.
//
// CurrentPassword is what turns Password from "set a first password" into
// "replace the one there is". An account that already has a password is
// refused without it, so a stolen session cannot quietly become a stolen
// account.
type UpdateRequest struct {
	Name            *string `json:"name"`
	Email           *string `json:"email"`
	Password        *string `json:"password"`
	CurrentPassword *string `json:"current_password"`
}

// ApiTokenResponse carries a freshly minted API token plus its metadata.
type ApiTokenResponse struct {
	Token     string `json:"token"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

// ApiTokenStatusResponse reports whether the user has an API token and its
// metadata when they do.
type ApiTokenStatusResponse struct {
	HasToken  bool   `json:"has_token"`
	Name      string `json:"name,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

// CreateApiTokenRequest is the body of POST /users/me/api-token.
type CreateApiTokenRequest struct {
	Name string `json:"name"`
}
