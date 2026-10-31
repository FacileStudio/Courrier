package mail

import (
	"fmt"

	"github.com/emersion/go-imap/v2/imapclient"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
)

type imapDialer func(host string, port int, user, password string) (*imapclient.Client, error)

type smtpDeliverer func(account schemas.Account, recipients []string, msg []byte) error

// connect opens an authenticated IMAP session for the account through the
// service's dialer, which tests replace to reach an in-memory server.
func (s *Service) connect(account schemas.Account) (*imapclient.Client, error) {
	return s.dial(account.IMAPHost, account.IMAPPort, account.IMAPUser, account.IMAPPassword)
}

// closeIMAP logs out and closes a session opened by connect.
func closeIMAP(client *imapclient.Client) {
	client.Logout().Wait()
	client.Close()
}

// deliverSMTP hands a finished message to the account's SMTP server, using
// implicit TLS on port 465 and STARTTLS everywhere else.
func deliverSMTP(account schemas.Account, recipients []string, msg []byte) error {
	addr := fmt.Sprintf("%s:%d", account.SMTPHost, account.SMTPPort)
	if account.SMTPPort == 465 {
		return sendImplicitTLS(addr, account.SMTPHost, account.SMTPUser, account.SMTPPassword, account.Email, recipients, msg)
	}
	return sendSTARTTLS(addr, account.SMTPHost, account.SMTPUser, account.SMTPPassword, account.Email, recipients, msg)
}
