package mail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
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

func buildMessage(fromEmail, fromName string, toAddrs, ccAddrs []string, subject, bodyText, bodyHTML, inReplyTo string, references []string, attachments []AttachmentUpload) ([]byte, error) {
	var buf bytes.Buffer

	var h gomail.Header
	h.SetDate(time.Now())
	h.SetSubject(subject)
	h.SetAddressList("From", []*gomail.Address{{Name: fromName, Address: fromEmail}})

	toList := make([]*gomail.Address, len(toAddrs))
	for i, addr := range toAddrs {
		toList[i] = &gomail.Address{Address: addr}
	}
	h.SetAddressList("To", toList)

	if len(ccAddrs) > 0 {
		ccList := make([]*gomail.Address, len(ccAddrs))
		for i, addr := range ccAddrs {
			ccList[i] = &gomail.Address{Address: addr}
		}
		h.SetAddressList("Cc", ccList)
	}

	if err := h.GenerateMessageIDWithHostname("courrier.local"); err != nil {
		return nil, err
	}

	if inReplyTo != "" {
		h.SetMsgIDList("In-Reply-To", []string{inReplyTo})
	}
	if len(references) > 0 {
		h.SetMsgIDList("References", references)
	}

	if len(attachments) == 0 {
		if bodyHTML == "" {
			w, err := gomail.CreateSingleInlineWriter(&buf, h)
			if err != nil {
				return nil, err
			}
			if _, err := io.WriteString(w, bodyText); err != nil {
				return nil, err
			}
			if err := w.Close(); err != nil {
				return nil, err
			}
			return buf.Bytes(), nil
		}

		if bodyText == "" {
			bodyText = stripHTMLTags(bodyHTML)
		}

		iw, err := gomail.CreateInlineWriter(&buf, h)
		if err != nil {
			return nil, err
		}
		if err := writeAlternativeParts(iw, bodyText, bodyHTML); err != nil {
			return nil, err
		}
		if err := iw.Close(); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}

	if bodyText == "" && bodyHTML != "" {
		bodyText = stripHTMLTags(bodyHTML)
	}

	mw, err := gomail.CreateWriter(&buf, h)
	if err != nil {
		return nil, err
	}

	if bodyHTML != "" {
		altPart, err := mw.CreateInline()
		if err != nil {
			return nil, err
		}
		if err := writeAlternativeParts(altPart, bodyText, bodyHTML); err != nil {
			return nil, err
		}
		if err := altPart.Close(); err != nil {
			return nil, err
		}
	} else {
		var textHeader gomail.InlineHeader
		textHeader.SetContentType("text/plain", map[string]string{"charset": "utf-8"})
		tw, err := mw.CreateSingleInline(textHeader)
		if err != nil {
			return nil, err
		}
		if _, err := io.WriteString(tw, bodyText); err != nil {
			return nil, err
		}
		if err := tw.Close(); err != nil {
			return nil, err
		}
	}

	for _, att := range attachments {
		var attHeader gomail.AttachmentHeader
		attHeader.SetContentType(att.MimeType, nil)
		attHeader.SetFilename(att.Filename)
		aw, err := mw.CreateAttachment(attHeader)
		if err != nil {
			return nil, err
		}
		if _, err := aw.Write(att.Data); err != nil {
			return nil, err
		}
		if err := aw.Close(); err != nil {
			return nil, err
		}
	}

	if err := mw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
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
