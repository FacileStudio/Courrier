package mail

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"

	"github.com/go-chi/chi/v5"
)

const backgroundBackfillTimeout = 15 * time.Minute

// BackfillResult reports what one backfill run touched.
type BackfillResult struct {
	Classified int  `json:"classified"`
	Scanned    int  `json:"scanned"`
	Running    bool `json:"running,omitempty"`
}

type imapHeaders struct {
	client *imapclient.Client
}

func (h imapHeaders) Headers(mailbox string, uids []imap.UID) (map[imap.UID]Signals, error) {
	if _, err := h.client.Select(mailbox, &imap.SelectOptions{ReadOnly: true}).Wait(); err != nil {
		return nil, err
	}
	opts := &imap.FetchOptions{UID: true, BodySection: []*imap.FetchItemBodySection{referencesSection}}
	msgs, err := h.client.Fetch(imap.UIDSetNum(uids...), opts).Collect()
	if err != nil {
		return nil, err
	}
	out := make(map[imap.UID]Signals, len(msgs))
	for _, m := range msgs {
		out[m.UID] = bulkSignals(m)
	}
	return out, nil
}

func (h imapHeaders) Close() {
	h.client.Logout().Wait()
	h.client.Close()
}

func (s *Service) unclassifiedCount(ctx context.Context, accountID int64) int64 {
	var n int64
	s.orm.WithContext(ctx).Table("emails e").Joins("JOIN folders f ON f.id = e.folder_id").
		Where("e.account_id = ? AND f.type = ? AND e.category = '' AND e.bulk_reason = '' AND e.date >= ?",
			accountID, schemas.FolderTypeInbox, time.Now().Add(-backfillWindow)).Count(&n)
	return n
}

func (s *Service) backfillInboxes(ctx context.Context, account schemas.Account) (int, error) {
	var folders []schemas.Folder
	err := s.orm.WithContext(ctx).Where("account_id = ? AND type = ?", account.ID, schemas.FolderTypeInbox).Find(&folders).Error
	if err != nil || len(folders) == 0 {
		return 0, err
	}
	client, err := connectIMAP(account.IMAPHost, account.IMAPPort, account.IMAPUser, account.IMAPPassword)
	if err != nil {
		return 0, errors.Failed(err.Error())
	}
	fetcher := imapHeaders{client: client}
	defer fetcher.Close()
	total := 0
	for _, folder := range folders {
		job := backfillJob{orm: s.orm.WithContext(ctx), fetcher: fetcher, folder: folder,
			cutoff: time.Now().Add(-backfillWindow), pause: backfillPause}
		n, err := job.run(ctx)
		total += n
		if err != nil {
			return total, err
		}
	}
	return total, nil
}

func (s *Service) ensureFolders(ctx context.Context, userID, accountID int64) error {
	var n int64
	s.orm.WithContext(ctx).Model(&schemas.Folder{}).
		Where("account_id = ? AND type IN ?", accountID, []string{schemas.FolderTypeInbox, schemas.FolderTypeSent}).Count(&n)
	if n >= 2 {
		return nil
	}
	return s.SyncAccount(ctx, userID, accountID)
}

// BackfillMail classifies cached inbox rows from the last 60 days that have no
// category yet, then deep-syncs the sent folder into contacts. Only one run per
// account is active at a time; a second call reports running and does nothing.
func (s *Service) BackfillMail(ctx context.Context, userID, accountID int64) (BackfillResult, error) {
	if _, loaded := s.backfilling.LoadOrStore(accountID, true); loaded {
		return BackfillResult{Running: true}, nil
	}
	defer s.backfilling.Delete(accountID)
	if err := s.ensureFolders(ctx, userID, accountID); err != nil {
		return BackfillResult{}, err
	}
	account, err := s.getAccountDecrypted(ctx, userID, accountID)
	if err != nil {
		return BackfillResult{}, err
	}
	classified, err := s.backfillInboxes(ctx, account)
	if err != nil {
		return BackfillResult{Classified: classified}, err
	}
	scanned, err := s.deepSyncSent(ctx, userID, accountID, maxDeepSyncSent)
	return BackfillResult{Classified: classified, Scanned: scanned}, err
}

func (s *Service) backfillAfterSync(userID, accountID int64) {
	if s.unclassifiedCount(context.Background(), accountID) == 0 {
		return
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), backgroundBackfillTimeout)
		defer cancel()
		res, err := s.BackfillMail(ctx, userID, accountID)
		if err != nil {
			slog.Warn("mail backfill failed", "account_id", accountID, "error", err)
			return
		}
		slog.Info("mail backfill finished", "account_id", accountID, "classified", res.Classified, "scanned", res.Scanned)
	}()
}

func (s *Service) handleBackfill(w http.ResponseWriter, req *http.Request) {
	uid, err := requestUserID(req)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	accountID, err := strconv.ParseInt(chi.URLParam(req, "accountId"), 10, 64)
	if err != nil {
		httpjson.WriteError(w, errors.Invalid("invalid account id"))
		return
	}
	res, err := s.BackfillMail(req.Context(), uid, accountID)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, res)
}
