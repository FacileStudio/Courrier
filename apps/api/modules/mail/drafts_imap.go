package mail

import (
	"fmt"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
)

// appendDraft files a message in the Drafts mailbox and returns the uid the
// server gave it. When the server reports no UIDPLUS data the uid is found by
// searching the Message-ID, taking the newest match.
func appendDraft(client *imapclient.Client, mailbox string, msg []byte, messageID string) (uint32, error) {
	cmd := client.Append(mailbox, int64(len(msg)), &imap.AppendOptions{
		Flags: []imap.Flag{imap.FlagDraft, imap.FlagSeen},
	})
	if _, err := cmd.Write(msg); err != nil {
		cmd.Close()
		return 0, fmt.Errorf("APPEND to %q failed: %w", mailbox, err)
	}
	if err := cmd.Close(); err != nil {
		return 0, fmt.Errorf("APPEND to %q failed: %w", mailbox, err)
	}
	data, err := cmd.Wait()
	if err != nil {
		return 0, fmt.Errorf("APPEND to %q failed: %w", mailbox, err)
	}
	if data.UID != 0 {
		return uint32(data.UID), nil
	}
	return findUIDByMessageID(client, mailbox, messageID)
}

func findUIDByMessageID(client *imapclient.Client, mailbox, messageID string) (uint32, error) {
	if _, err := client.Select(mailbox, nil).Wait(); err != nil {
		return 0, fmt.Errorf("SELECT %q failed: %w", mailbox, err)
	}
	criteria := &imap.SearchCriteria{Header: []imap.SearchCriteriaHeaderField{{Key: "Message-ID", Value: messageID}}}
	found, err := client.UIDSearch(criteria, nil).Wait()
	if err != nil {
		return 0, fmt.Errorf("SEARCH in %q failed: %w", mailbox, err)
	}
	uids := found.AllUIDs()
	if len(uids) == 0 {
		return 0, fmt.Errorf("draft not found in %q after append", mailbox)
	}
	newest := uids[0]
	for _, uid := range uids {
		if uid > newest {
			newest = uid
		}
	}
	return uint32(newest), nil
}

// fetchAttachmentRows reads the attachments of a message from its body
// structure, with the part ids the server assigns.
func fetchAttachmentRows(client *imapclient.Client, mailbox string, uid uint32, emailID int64) ([]schemas.Attachment, error) {
	if _, err := client.Select(mailbox, nil).Wait(); err != nil {
		return nil, fmt.Errorf("SELECT %q failed: %w", mailbox, err)
	}
	msgs, err := client.Fetch(imap.UIDSetNum(imap.UID(uid)), &imap.FetchOptions{
		BodyStructure: &imap.FetchItemBodyStructure{Extended: true},
	}).Collect()
	if err != nil || len(msgs) == 0 {
		return nil, fmt.Errorf("FETCH structure of uid %d failed", uid)
	}
	return extractAttachments(msgs[0].BodyStructure, emailID), nil
}

// expungeUIDs deletes messages from a mailbox. With UIDPLUS only the named
// messages go; without it the server expunges every message flagged
// \Deleted, which is the best a bare IMAP4rev1 server allows.
func expungeUIDs(client *imapclient.Client, mailbox string, uids []imap.UID) error {
	if len(uids) == 0 {
		return nil
	}
	if _, err := client.Select(mailbox, nil).Wait(); err != nil {
		return fmt.Errorf("SELECT %q failed: %w", mailbox, err)
	}
	set := imap.UIDSetNum(uids...)
	flag := client.Store(set, &imap.StoreFlags{Op: imap.StoreFlagsAdd, Silent: true, Flags: []imap.Flag{imap.FlagDeleted}}, nil)
	if err := flag.Close(); err != nil {
		return fmt.Errorf("flag deleted failed: %w", err)
	}
	if client.Caps().Has(imap.CapUIDPlus) {
		return client.UIDExpunge(set).Close()
	}
	return client.Expunge().Close()
}
