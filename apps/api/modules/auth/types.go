package auth

// RegisterRequest is the body of POST /auth/register.
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest is the body of POST /auth/login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse is the {user_id, token} shape auth has always returned.
type AuthResponse struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
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
