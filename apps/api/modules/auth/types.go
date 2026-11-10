package auth

// ResourceTokenResponse is the body of GET /auth/resource-token.
type ResourceTokenResponse struct {
	Token string `json:"token"`
}
