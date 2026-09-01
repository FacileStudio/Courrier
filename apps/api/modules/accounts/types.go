package accounts

// CreateAccountRequest is the body of POST /accounts.
type CreateAccountRequest struct {
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	IMAPHost     string  `json:"imap_host"`
	IMAPPort     int     `json:"imap_port"`
	IMAPUser     string  `json:"imap_user"`
	IMAPPassword string  `json:"imap_password"`
	SMTPHost     string  `json:"smtp_host"`
	SMTPPort     int     `json:"smtp_port"`
	SMTPUser     string  `json:"smtp_user"`
	SMTPPassword string  `json:"smtp_password"`
	Signature    string  `json:"signature"`
	IsDefault    bool    `json:"is_default"`
	SpaceID      *string `json:"space_id"`
}

// UpdateAccountRequest is the body of PUT /accounts/{id}; every field is a
// pointer so only the fields present in the body change.
type UpdateAccountRequest struct {
	Name         *string `json:"name"`
	Email        *string `json:"email"`
	IMAPHost     *string `json:"imap_host"`
	IMAPPort     *int    `json:"imap_port"`
	IMAPUser     *string `json:"imap_user"`
	IMAPPassword *string `json:"imap_password"`
	SMTPHost     *string `json:"smtp_host"`
	SMTPPort     *int    `json:"smtp_port"`
	SMTPUser     *string `json:"smtp_user"`
	SMTPPassword *string `json:"smtp_password"`
	Signature    *string `json:"signature"`
	IsDefault    *bool   `json:"is_default"`
}

// AccountResponse is the account shape returned to the client, with secrets
// omitted.
type AccountResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	IMAPHost  string `json:"imap_host"`
	IMAPPort  int    `json:"imap_port"`
	IMAPUser  string `json:"imap_user"`
	SMTPHost  string `json:"smtp_host"`
	SMTPPort  int    `json:"smtp_port"`
	SMTPUser  string `json:"smtp_user"`
	Signature string `json:"signature"`
	IsDefault bool   `json:"is_default"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// AccountListResponse wraps the list of mail accounts.
type AccountListResponse struct {
	Accounts []AccountResponse `json:"accounts"`
}

// DeleteAccountResponse is returned when an account is deleted.
type DeleteAccountResponse struct {
	Deleted bool `json:"deleted"`
}
