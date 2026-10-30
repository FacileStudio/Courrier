package mail

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"
)

const maxDeepSyncSent = 1000

func imapRecipients(groups ...[]imap.Address) []recipient {
	out := []recipient{}
	for _, group := range groups {
		for _, a := range group {
			out = append(out, recipient{Email: a.Addr(), Name: a.Name})
		}
	}
	return out
}

func sendRecipients(groups ...[]string) []recipient {
	out := []recipient{}
	for _, group := range groups {
		for _, a := range parseAddresses(group) {
			out = append(out, recipient{Email: a.Address, Name: a.Name})
		}
	}
	return out
}

func (s *Service) recordSyncedSent(ctx context.Context, folder schemas.Folder, spaceID *string, env *imap.Envelope) {
	if folder.Type != schemas.FolderTypeSent {
		return
	}
	if env.MessageID != "" {
		var copies int64
		s.orm.WithContext(ctx).Model(&schemas.Email{}).Where("folder_id = ? AND message_id = ?", folder.ID, env.MessageID).Count(&copies)
		if copies > 1 {
			return
		}
	}
	s.recordContacts(ctx, folder.AccountID, spaceID, imapRecipients(env.To, env.Cc), env.Date)
}

func (s *Service) recordContacts(ctx context.Context, accountID int64, spaceID *string, recipients []recipient, sentAt time.Time) {
	if err := upsertContacts(ctx, s.orm, accountID, spaceID, recipients, sentAt); err != nil {
		slog.Warn("contacts upsert failed", "account_id", accountID, "error", err)
	}
}

func fetchSentEnvelopes(client *imapclient.Client, mailbox string, limit int) ([]*imapclient.FetchMessageBuffer, error) {
	selectData, err := client.Select(mailbox, &imap.SelectOptions{ReadOnly: true}).Wait()
	if err != nil {
		return nil, fmt.Errorf("SELECT %q failed: %w", mailbox, err)
	}
	if selectData.NumMessages == 0 {
		return nil, nil
	}
	from := uint32(1)
	if selectData.NumMessages > uint32(limit) {
		from = selectData.NumMessages - uint32(limit) + 1
	}
	var seqSet imap.SeqSet
	seqSet.AddRange(from, selectData.NumMessages)
	return client.Fetch(seqSet, &imap.FetchOptions{Envelope: true}).Collect()
}

func readSentEnvelopes(account schemas.Account, path string, limit int) ([]*imapclient.FetchMessageBuffer, error) {
	client, err := connectIMAP(account.IMAPHost, account.IMAPPort, account.IMAPUser, account.IMAPPassword)
	if err != nil {
		return nil, errors.Failed(err.Error())
	}
	defer func() {
		client.Logout().Wait()
		client.Close()
	}()
	msgs, err := fetchSentEnvelopes(client, path, limit)
	if err != nil {
		return nil, errors.Internal("failed to fetch sent envelopes", err)
	}
	return msgs, nil
}

// deepSyncSent reads up to the newest limit envelopes of the account's sent
// folder over IMAP, never above 1000, and merges their recipients into the
// contacts table. It stores no Email rows and can be run any number of times:
// counts only ever rise to what the envelopes show.
func (s *Service) deepSyncSent(ctx context.Context, userID, accountID int64, limit int) (int, error) {
	if limit < 1 || limit > maxDeepSyncSent {
		limit = maxDeepSyncSent
	}
	account, err := s.getAccountDecrypted(ctx, userID, accountID)
	if err != nil {
		return 0, err
	}
	var folder schemas.Folder
	if s.orm.WithContext(ctx).Where("account_id = ? AND type = ?", accountID, schemas.FolderTypeSent).First(&folder).Error != nil {
		return 0, errors.NotFound("sent folder not found; sync the account first")
	}
	msgs, err := readSentEnvelopes(account, folder.Path, limit)
	if err != nil {
		return 0, err
	}
	entries := map[string]*contactEntry{}
	for _, msg := range msgs {
		if msg.Envelope != nil {
			addRecipients(entries, account.Email, imapRecipients(msg.Envelope.To, msg.Envelope.Cc), msg.Envelope.Date)
		}
	}
	if err := saveContacts(ctx, s.orm, accountID, account.SpaceID, entries, true); err != nil {
		return 0, errors.Internal("failed to save contacts", err)
	}
	return len(msgs), nil
}
