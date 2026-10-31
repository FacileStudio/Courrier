package mail

import (
	"bytes"
	"strings"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
)

// messageSpec is everything buildMessage needs to write one RFC822 message.
// MessageID is optional: an update keeps the Message-ID of the draft it
// replaces instead of minting a new one.
type messageSpec struct {
	FromEmail   string
	FromName    string
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	Body        string
	BodyHTML    string
	InReplyTo   string
	References  []string
	Attachments []AttachmentUpload
	MessageID   string
}

// specFor turns a send or draft request into the spec of a message written
// from the account.
func specFor(account schemas.Account, req SendRequest) messageSpec {
	return messageSpec{
		FromEmail:   account.Email,
		FromName:    account.Name,
		To:          req.To,
		Cc:          req.Cc,
		Bcc:         req.Bcc,
		Subject:     req.Subject,
		Body:        req.Body,
		BodyHTML:    req.BodyHTML,
		InReplyTo:   req.InReplyTo,
		References:  req.References,
		Attachments: req.Attachments,
	}
}

// withoutBcc drops the Bcc header from a serialized message. The copy handed
// to SMTP must not carry it: every recipient would otherwise read the blind
// list. The copy filed in Sent or Drafts keeps it.
func withoutBcc(msg []byte) []byte {
	end := bytes.Index(msg, []byte("\r\n\r\n"))
	if end < 0 {
		return msg
	}
	lines := bytes.Split(msg[:end], []byte("\r\n"))
	kept := make([][]byte, 0, len(lines))
	dropping := false
	for _, line := range lines {
		continued := len(line) > 0 && (line[0] == ' ' || line[0] == '\t')
		if !continued {
			dropping = len(line) >= 4 && strings.EqualFold(string(line[:4]), "bcc:")
		}
		if !dropping {
			kept = append(kept, line)
		}
	}
	return append(bytes.Join(kept, []byte("\r\n")), msg[end:]...)
}

// envelopeRecipients lists the SMTP RCPT addresses of a request: To, Cc and
// Bcc, each address once.
func envelopeRecipients(req SendRequest) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, group := range [][]string{req.To, req.Cc, req.Bcc} {
		for _, addr := range extractEmails(group) {
			key := strings.ToLower(addr)
			if !seen[key] {
				seen[key] = true
				out = append(out, addr)
			}
		}
	}
	return out
}

// splitList splits a comma-separated form value into trimmed, non-empty items.
func splitList(value string) []string {
	var out []string
	for _, item := range strings.Split(value, ",") {
		if item = strings.TrimSpace(item); item != "" {
			out = append(out, item)
		}
	}
	return out
}
