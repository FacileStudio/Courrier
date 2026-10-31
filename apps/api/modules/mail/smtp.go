package mail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
	"time"

	gomail "github.com/emersion/go-message/mail"
)

func sendSTARTTLS(addr, host, user, password, from string, to []string, msg []byte) error {
	c, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("SMTP dial failed: %w", err)
	}
	defer c.Close()

	if err := c.Hello("localhost"); err != nil {
		return err
	}
	if ok, _ := c.Extension("STARTTLS"); ok {
		if err := c.StartTLS(&tls.Config{ServerName: host}); err != nil {
			return fmt.Errorf("STARTTLS failed: %w", err)
		}
	}
	auth := smtp.PlainAuth("", user, password, host)
	if err := c.Auth(auth); err != nil {
		return fmt.Errorf("SMTP auth failed: %w", err)
	}
	if err := c.Mail(from); err != nil {
		return err
	}
	for _, rcpt := range to {
		if err := c.Rcpt(rcpt); err != nil {
			return fmt.Errorf("RCPT %s failed: %w", rcpt, err)
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(msg); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return c.Quit()
}

func sendImplicitTLS(addr, host, user, password, from string, to []string, msg []byte) error {
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 10 * time.Second}, "tcp", addr, &tls.Config{ServerName: host})
	if err != nil {
		return fmt.Errorf("TLS dial failed: %w", err)
	}
	defer conn.Close()

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("SMTP client create failed: %w", err)
	}
	defer c.Close()

	auth := smtp.PlainAuth("", user, password, host)
	if err := c.Auth(auth); err != nil {
		return fmt.Errorf("SMTP auth failed: %w", err)
	}
	if err := c.Mail(from); err != nil {
		return err
	}
	for _, rcpt := range to {
		if err := c.Rcpt(rcpt); err != nil {
			return fmt.Errorf("RCPT %s failed: %w", rcpt, err)
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(msg); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return c.Quit()
}

// buildMessage builds the RFC822 message and returns its bytes plus the
// generated Message-ID, so callers can persist it and thread outbound mail
// against the replies it later receives.
func buildMessage(spec messageSpec) ([]byte, string, error) {
	var buf bytes.Buffer

	var h gomail.Header
	h.SetDate(time.Now())
	h.SetSubject(spec.Subject)
	h.SetAddressList("From", []*gomail.Address{{Name: spec.FromName, Address: spec.FromEmail}})

	if to := parseAddresses(spec.To); len(to) > 0 {
		h.SetAddressList("To", to)
	}

	if len(spec.Cc) > 0 {
		h.SetAddressList("Cc", parseAddresses(spec.Cc))
	}
	if bcc := parseAddresses(spec.Bcc); len(bcc) > 0 {
		h.SetAddressList("Bcc", bcc)
	}

	if spec.MessageID != "" {
		h.Set("Message-Id", spec.MessageID)
	} else if err := h.GenerateMessageIDWithHostname("courrier.local"); err != nil {
		return nil, "", err
	}
	messageID := h.Get("Message-Id")

	if spec.InReplyTo != "" {
		h.SetMsgIDList("In-Reply-To", []string{spec.InReplyTo})
	}
	if len(spec.References) > 0 {
		h.SetMsgIDList("References", spec.References)
	}

	if len(spec.Attachments) == 0 {
		if spec.BodyHTML == "" {
			w, err := gomail.CreateSingleInlineWriter(&buf, h)
			if err != nil {
				return nil, "", err
			}
			if _, err := io.WriteString(w, spec.Body); err != nil {
				return nil, "", err
			}
			if err := w.Close(); err != nil {
				return nil, "", err
			}
			return buf.Bytes(), messageID, nil
		}

		if spec.Body == "" {
			spec.Body = stripHTMLTags(spec.BodyHTML)
		}

		iw, err := gomail.CreateInlineWriter(&buf, h)
		if err != nil {
			return nil, "", err
		}
		if err := writeAlternativeParts(iw, spec.Body, spec.BodyHTML); err != nil {
			return nil, "", err
		}
		if err := iw.Close(); err != nil {
			return nil, "", err
		}
		return buf.Bytes(), messageID, nil
	}

	if spec.Body == "" && spec.BodyHTML != "" {
		spec.Body = stripHTMLTags(spec.BodyHTML)
	}

	mw, err := gomail.CreateWriter(&buf, h)
	if err != nil {
		return nil, "", err
	}

	if spec.BodyHTML != "" {
		altPart, err := mw.CreateInline()
		if err != nil {
			return nil, "", err
		}
		if err := writeAlternativeParts(altPart, spec.Body, spec.BodyHTML); err != nil {
			return nil, "", err
		}
		if err := altPart.Close(); err != nil {
			return nil, "", err
		}
	} else {
		var textHeader gomail.InlineHeader
		textHeader.SetContentType("text/plain", map[string]string{"charset": "utf-8"})
		tw, err := mw.CreateSingleInline(textHeader)
		if err != nil {
			return nil, "", err
		}
		if _, err := io.WriteString(tw, spec.Body); err != nil {
			return nil, "", err
		}
		if err := tw.Close(); err != nil {
			return nil, "", err
		}
	}

	for i := range spec.Attachments {
		var attHeader gomail.AttachmentHeader
		attHeader.SetContentType(spec.Attachments[i].MimeType, nil)
		attHeader.SetFilename(spec.Attachments[i].Filename)
		aw, err := mw.CreateAttachment(attHeader)
		if err != nil {
			return nil, "", err
		}
		reader, err := spec.Attachments[i].Reader()
		if err != nil {
			return nil, "", err
		}
		if _, err := io.Copy(aw, reader); err != nil {
			reader.Close()
			return nil, "", err
		}
		reader.Close()
		if err := aw.Close(); err != nil {
			return nil, "", err
		}
	}

	if err := mw.Close(); err != nil {
		return nil, "", err
	}

	return buf.Bytes(), messageID, nil
}

func writeAlternativeParts(iw *gomail.InlineWriter, bodyText, bodyHTML string) error {
	var textHeader gomail.InlineHeader
	textHeader.SetContentType("text/plain", map[string]string{"charset": "utf-8"})
	tw, err := iw.CreatePart(textHeader)
	if err != nil {
		return err
	}
	if _, err := io.WriteString(tw, bodyText); err != nil {
		return err
	}
	if err := tw.Close(); err != nil {
		return err
	}

	var htmlHeader gomail.InlineHeader
	htmlHeader.SetContentType("text/html", map[string]string{"charset": "utf-8"})
	hw, err := iw.CreatePart(htmlHeader)
	if err != nil {
		return err
	}
	if _, err := io.WriteString(hw, bodyHTML); err != nil {
		return err
	}
	return hw.Close()
}

func stripHTMLTags(s string) string {
	var result strings.Builder
	inTag := false
	for _, r := range s {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case !inTag:
			result.WriteRune(r)
		}
	}
	return result.String()
}

func testSMTP(host string, port int, user, password string) error {
	addr := fmt.Sprintf("%s:%d", host, port)

	if port == 465 {
		conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 10 * time.Second}, "tcp", addr, &tls.Config{ServerName: host})
		if err != nil {
			return fmt.Errorf("TLS dial failed: %w", err)
		}
		defer conn.Close()
		c, err := smtp.NewClient(conn, host)
		if err != nil {
			return fmt.Errorf("SMTP client create failed: %w", err)
		}
		defer c.Close()
		auth := smtp.PlainAuth("", user, password, host)
		if err := c.Auth(auth); err != nil {
			return fmt.Errorf("SMTP auth failed: %w", err)
		}
		return c.Quit()
	}

	c, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("SMTP dial failed: %w", err)
	}
	defer c.Close()
	if err := c.Hello("localhost"); err != nil {
		return err
	}
	if ok, _ := c.Extension("STARTTLS"); ok {
		if err := c.StartTLS(&tls.Config{ServerName: host}); err != nil {
			return fmt.Errorf("STARTTLS failed: %w", err)
		}
	}
	auth := smtp.PlainAuth("", user, password, host)
	if err := c.Auth(auth); err != nil {
		return fmt.Errorf("SMTP auth failed: %w", err)
	}
	return c.Quit()
}

func parseAddresses(addrs []string) []*gomail.Address {
	out := make([]*gomail.Address, 0, len(addrs))
	for _, a := range addrs {
		a = strings.TrimSpace(a)
		if a == "" {
			continue
		}
		if parsed, err := mail.ParseAddress(a); err == nil {
			out = append(out, &gomail.Address{Name: parsed.Name, Address: parsed.Address})
		} else {
			out = append(out, &gomail.Address{Address: a})
		}
	}
	return out
}

func extractEmails(addrs []string) []string {
	out := make([]string, 0, len(addrs))
	for _, a := range addrs {
		a = strings.TrimSpace(a)
		if a == "" {
			continue
		}
		if parsed, err := mail.ParseAddress(a); err == nil {
			out = append(out, parsed.Address)
		} else {
			out = append(out, a)
		}
	}
	return out
}
