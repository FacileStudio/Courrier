package accounts

type CreateAccountRequest struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	IMAPHost     string `json:"imap_host"`
	IMAPPort     int    `json:"imap_port"`
	IMAPUser     string `json:"imap_user"`
	IMAPPassword string `json:"imap_password"`
	SMTPHost     string `json:"smtp_host"`
	SMTPPort     int    `json:"smtp_port"`
	SMTPUser     string `json:"smtp_user"`
	SMTPPassword string `json:"smtp_password"`
	Signature    string `json:"signature"`
	IsDefault    bool   `json:"is_default"`
}

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
