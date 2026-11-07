package auth

// ResourceTokenResponse is the body of GET /auth/resource-token.
type ResourceTokenResponse struct {
	Token string `json:"token"`
}

// Data is the profile payload a user carries through the auth flow.
type Data struct {
	Email string `json:"email"`
}

func (d *Data) GetEmail() string {
	if d == nil {
		return ""
	}
	return d.Email
}
