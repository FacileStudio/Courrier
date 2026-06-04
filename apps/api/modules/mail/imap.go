package mail

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	gomessage "github.com/emersion/go-message"
	"github.com/emersion/go-message/mail"

	"api/schemas"
)

func connectIMAP(host string, port int, user, password string) (*imapclient.Client, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	client, err := imapclient.DialTLS(addr, nil)
	if err != nil {
		return nil, fmt.Errorf("IMAP dial failed: %w", err)
	}
	if err := client.Login(user, password).Wait(); err != nil {
		client.Close()
		return nil, fmt.Errorf("IMAP login failed: %w", err)
	}
	return client, nil
}

func listMailboxes(client *imapclient.Client) ([]*imap.ListData, error) {
	cmd := client.List("", "*", &imap.ListOptions{
		ReturnSpecialUse: true,
	})
	mailboxes, err := cmd.Collect()
	if err != nil {
		cmd2 := client.List("", "*", nil)
		mailboxes, err = cmd2.Collect()
		if err != nil {
			return nil, fmt.Errorf("LIST failed: %w", err)
		}
	}
	return mailboxes, nil
}

func detectFolderType(mbox *imap.ListData) string {
	for _, attr := range mbox.Attrs {
		switch attr {
		case imap.MailboxAttrDrafts:
			return schemas.FolderTypeDrafts
		case imap.MailboxAttrSent:
			return schemas.FolderTypeSent
		case imap.MailboxAttrTrash:
			return schemas.FolderTypeTrash
		case imap.MailboxAttrJunk:
			return schemas.FolderTypeJunk
		case imap.MailboxAttrArchive:
			return schemas.FolderTypeArchive
		}
	}

	lower := strings.ToLower(mbox.Mailbox)
	switch {
	case lower == "inbox":
		return schemas.FolderTypeInbox
	case strings.Contains(lower, "sent"):
		return schemas.FolderTypeSent
	case strings.Contains(lower, "draft"):
		return schemas.FolderTypeDrafts
	case strings.Contains(lower, "trash") || strings.Contains(lower, "deleted"):
		return schemas.FolderTypeTrash
	case strings.Contains(lower, "junk") || strings.Contains(lower, "spam"):
		return schemas.FolderTypeJunk
	case strings.Contains(lower, "archive"):
		return schemas.FolderTypeArchive
	default:
		return schemas.FolderTypeCustom
	}
}

func folderDisplayName(mbox *imap.ListData) string {
	parts := strings.Split(mbox.Mailbox, string(mbox.Delim))
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return mbox.Mailbox
}

func fetchEnvelopes(client *imapclient.Client, mailbox string, limit int) ([]*imapclient.FetchMessageBuffer, *imap.SelectData, error) {
	selectData, err := client.Select(mailbox, nil).Wait()
	if err != nil {
		return nil, nil, fmt.Errorf("SELECT %q failed: %w", mailbox, err)
	}

	if selectData.NumMessages == 0 {
		return nil, selectData, nil
	}

	from := uint32(1)
	if selectData.NumMessages > uint32(limit) {
		from = selectData.NumMessages - uint32(limit) + 1
	}

	var seqSet imap.SeqSet
	seqSet.AddRange(from, selectData.NumMessages)

	fetchCmd := client.Fetch(seqSet, &imap.FetchOptions{
		Envelope: true,
		Flags:    true,
		UID:      true,
		BodyStructure: &imap.FetchItemBodyStructure{
			Extended: false,
		},
	})
	msgs, err := fetchCmd.Collect()
	if err != nil {
		return nil, selectData, fmt.Errorf("FETCH envelopes failed: %w", err)
	}
	return msgs, selectData, nil
}

func fetchMessageBody(client *imapclient.Client, mailbox string, uid imap.UID) (string, string, error) {
	_, err := client.Select(mailbox, nil).Wait()
	if err != nil {
		return "", "", fmt.Errorf("SELECT %q failed: %w", mailbox, err)
	}

	uidSet := imap.UIDSetNum(uid)
	bodySection := &imap.FetchItemBodySection{
		Specifier: imap.PartSpecifierNone,
		Peek:      true,
	}
	fetchCmd := client.Fetch(uidSet, &imap.FetchOptions{
		BodySection: []*imap.FetchItemBodySection{bodySection},
	})
	msgs, err := fetchCmd.Collect()
	if err != nil {
		return "", "", fmt.Errorf("FETCH body failed: %w", err)
	}
	if len(msgs) == 0 {
		return "", "", fmt.Errorf("message not found")
	}

	msg := msgs[0]
	if len(msg.BodySection) == 0 {
		return "", "", fmt.Errorf("no body section in response")
	}

	raw := msg.BodySection[0].Bytes
	return parseMessageBody(raw)
}

func parseMessageBody(raw []byte) (string, string, error) {
	entity, err := gomessage.Read(strings.NewReader(string(raw)))
	if err != nil {
		return string(raw), "", nil
	}

	var textBody, htmlBody string

	if mr := entity.MultipartReader(); mr != nil {
		collectParts(mr, &textBody, &htmlBody)
	} else {
		mediaType, _, _ := entity.Header.ContentType()
		body, _ := io.ReadAll(entity.Body)
		switch {
		case strings.HasPrefix(mediaType, "text/html"):
			htmlBody = string(body)
		default:
			textBody = string(body)
		}
	}

	return textBody, htmlBody, nil
}

func collectParts(mr gomessage.MultipartReader, textBody, htmlBody *string) {
	for {
		part, err := mr.NextPart()
		if err != nil {
			break
		}

		if nested := part.MultipartReader(); nested != nil {
			collectParts(nested, textBody, htmlBody)
			continue
		}

		mediaType, _, _ := part.Header.ContentType()
		body, err := io.ReadAll(part.Body)
		if err != nil {
			continue
		}

		switch {
		case strings.HasPrefix(mediaType, "text/plain") && *textBody == "":
			*textBody = string(body)
		case strings.HasPrefix(mediaType, "text/html") && *htmlBody == "":
			*htmlBody = string(body)
		}
	}
}

func hasAttachments(bs imap.BodyStructure) bool {
	if bs == nil {
		return false
	}
	found := false
	bs.Walk(func(path []int, part imap.BodyStructure) bool {
		if sp, ok := part.(*imap.BodyStructureSinglePart); ok {
			if sp.Filename() != "" {
				found = true
				return false
			}
			disp := sp.Disposition()
			if disp != nil && strings.EqualFold(disp.Value, "attachment") {
				found = true
				return false
			}
		}
		return true
	})
	return found
}

func imapAddressesToJSON(addrs []imap.Address) string {
	if len(addrs) == 0 {
		return "[]"
	}
	type addrEntry struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	entries := make([]addrEntry, len(addrs))
	for i, a := range addrs {
		entries[i] = addrEntry{Name: a.Name, Email: a.Addr()}
	}
	b, _ := json.Marshal(entries)
	return string(b)
}

func storeFlags(client *imapclient.Client, mailbox string, uid imap.UID, op imap.StoreFlagsOp, flags []imap.Flag) error {
	_, err := client.Select(mailbox, nil).Wait()
	if err != nil {
		return fmt.Errorf("SELECT %q failed: %w", mailbox, err)
	}

	uidSet := imap.UIDSetNum(uid)
	storeCmd := client.Store(uidSet, &imap.StoreFlags{
		Op:     op,
		Silent: true,
		Flags:  flags,
	}, nil)
	return storeCmd.Close()
}

var _ = mail.Header{}
