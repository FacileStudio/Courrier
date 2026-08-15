package mail

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/emersion/go-imap/v2"

	"github.com/FacileStudio/Courrier/apps/api/internal/crypto"
	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"

	"gorm.io/gorm"
)

// Service is the mail module's data access: it owns email sync, threading,
// folders, attachments and sending, and decrypts the account credentials it
// works with.
type Service struct {
	orm           *gorm.DB
	encryptionKey []byte
}

// NewService builds a mail Service over the database and the app's encryption
// key.
func NewService(orm *gorm.DB, encryptionKey []byte) *Service {
	return &Service{orm: orm, encryptionKey: encryptionKey}
}

func (s *Service) decryptPassword(cipher string) (string, error) {
	if len(s.encryptionKey) == 0 || cipher == "" {
		return cipher, nil
	}
	return crypto.Decrypt(cipher, s.encryptionKey)
}

func (s *Service) getAccount(ctx context.Context, userID, accountID int64) (schemas.Account, error) {
	var account schemas.Account
	if err := s.orm.WithContext(ctx).Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		return schemas.Account{}, errors.NotFound("account not found")
	}
	return account, nil
}

func (s *Service) getAccountDecrypted(ctx context.Context, userID, accountID int64) (schemas.Account, error) {
	account, err := s.getAccount(ctx, userID, accountID)
	if err != nil {
		return account, err
	}
	account.IMAPPassword, err = s.decryptPassword(account.IMAPPassword)
	if err != nil {
		return schemas.Account{}, errors.Internal("failed to decrypt IMAP password", err)
	}
	account.SMTPPassword, err = s.decryptPassword(account.SMTPPassword)
	if err != nil {
		return schemas.Account{}, errors.Internal("failed to decrypt SMTP password", err)
	}
	return account, nil
}

/*
CheckAccount answers "can this saved account still talk to its servers", which
TestConnection cannot: that one needs the passwords in the request body, and the stored
ones never leave the server. It opens each leg and closes it again — no mailbox listing,
no message fetch, nothing written.

Both legs always run. Stopping at the first failure is what TestConnection does, and it
hides the answer people actually need while fixing an account: which half is broken.
*/
func (s *Service) CheckAccount(ctx context.Context, userID, accountID int64) (CheckResult, error) {
	account, err := s.getAccountDecrypted(ctx, userID, accountID)
	if err != nil {
		return CheckResult{}, err
	}

	var result CheckResult

	if account.IMAPHost != "" {
		port := account.IMAPPort
		if port == 0 {
			port = 993
		}
		result.IMAP.Configured = true
		client, err := connectIMAP(account.IMAPHost, port, account.IMAPUser, account.IMAPPassword)
		if err != nil {
			result.IMAP.Error = err.Error()
		} else {
			result.IMAP.OK = true
			client.Logout().Wait()
			client.Close()
		}
	}

	if account.SMTPHost != "" {
		port := account.SMTPPort
		if port == 0 {
			port = 587
		}
		result.SMTP.Configured = true
		if err := testSMTP(account.SMTPHost, port, account.SMTPUser, account.SMTPPassword); err != nil {
			result.SMTP.Error = err.Error()
		} else {
			result.SMTP.OK = true
		}
	}

	return result, nil
}

func (s *Service) SyncAccount(ctx context.Context, userID, accountID int64) error {
	account, err := s.getAccountDecrypted(ctx, userID, accountID)
	if err != nil {
		return err
	}

	client, err := connectIMAP(account.IMAPHost, account.IMAPPort, account.IMAPUser, account.IMAPPassword)
	if err != nil {
		return errors.Failed(err.Error())
	}
	defer func() {
		client.Logout().Wait()
		client.Close()
	}()

	mailboxes, err := listMailboxes(client)
	if err != nil {
		return errors.Internal("failed to list folders", err)
	}

	for _, mbox := range mailboxes {
		if isNoSelect(mbox) {
			continue
		}

		folderType := detectFolderType(mbox)
		name := folderDisplayName(mbox)

		var folder schemas.Folder
		result := s.orm.WithContext(ctx).Where("account_id = ? AND path = ?", accountID, mbox.Mailbox).First(&folder)
		if result.Error != nil {
			folder = schemas.Folder{
				AccountID: accountID,
				SpaceID:   account.SpaceID,
				Path:      mbox.Mailbox,
				Name:      name,
				Type:      folderType,
			}
			s.orm.WithContext(ctx).Create(&folder)
		} else {
			s.orm.WithContext(ctx).Model(&folder).Updates(map[string]any{
				"name": name,
				"type": folderType,
			})
		}

		selectData, err := client.Select(mbox.Mailbox, &imap.SelectOptions{ReadOnly: true}).Wait()
		if err != nil {
			continue
		}

		s.orm.WithContext(ctx).Model(&folder).Updates(map[string]any{
			"total_count":  selectData.NumMessages,
			"uid_validity": selectData.UIDValidity,
		})
	}

	return nil
}

func (s *Service) SyncFolderEmails(ctx context.Context, userID, accountID, folderID int64) error {
	account, err := s.getAccountDecrypted(ctx, userID, accountID)
	if err != nil {
		return err
	}

	var folder schemas.Folder
	if err := s.orm.WithContext(ctx).Where("id = ? AND account_id = ?", folderID, accountID).First(&folder).Error; err != nil {
		return errors.NotFound("folder not found")
	}

	client, err := connectIMAP(account.IMAPHost, account.IMAPPort, account.IMAPUser, account.IMAPPassword)
	if err != nil {
		return errors.Failed(err.Error())
	}
	defer func() {
		client.Logout().Wait()
		client.Close()
	}()

	msgs, selectData, err := fetchEnvelopes(client, folder.Path, 100)
	if err != nil {
		return errors.Internal("failed to fetch emails", err)
	}

	if selectData != nil && selectData.UIDValidity != folder.UIDValidity && folder.UIDValidity != 0 {
		s.orm.WithContext(ctx).Where("folder_id = ?", folderID).Delete(&schemas.Email{})
	}

	if selectData != nil {
		s.orm.WithContext(ctx).Model(&folder).Updates(map[string]any{
			"total_count":  selectData.NumMessages,
			"uid_validity": selectData.UIDValidity,
		})
	}

	for _, msg := range msgs {
		if msg.Envelope == nil {
			continue
		}

		var existing schemas.Email
		if s.orm.WithContext(ctx).Where("folder_id = ? AND imap_uid = ?", folderID, msg.UID).First(&existing).Error == nil {
			isRead := containsFlag(msg.Flags, imap.FlagSeen)
			isStarred := containsFlag(msg.Flags, imap.FlagFlagged)
			if existing.IsRead != isRead || existing.IsStarred != isStarred {
				s.orm.WithContext(ctx).Model(&existing).Updates(map[string]any{
					"is_read":    isRead,
					"is_starred": isStarred,
				})
			}
			continue
		}

		env := msg.Envelope
		attachments := extractAttachments(msg.BodyStructure, 0)
		references := messageReferences(msg)
		inReplyTo := strings.Join(env.InReplyTo, " ")

		email := schemas.Email{
			AccountID:      accountID,
			SpaceID:        account.SpaceID,
			FolderID:       folderID,
			MessageID:      env.MessageID,
			ThreadID:       assignThreadID(s.orm.WithContext(ctx), accountID, env.MessageID, inReplyTo, references),
			Subject:        env.Subject,
			FromAddress:    firstAddr(env.From),
			FromName:       firstName(env.From),
			ToAddresses:    imapAddressesToJSON(env.To),
			CcAddresses:    imapAddressesToJSON(env.Cc),
			Date:           env.Date,
			IsRead:         containsFlag(msg.Flags, imap.FlagSeen),
			IsStarred:      containsFlag(msg.Flags, imap.FlagFlagged),
			HasAttachments: len(attachments) > 0,
			InReplyTo:      inReplyTo,
			References:     references,
			IMAPUID:        uint32(msg.UID),
		}

		s.orm.WithContext(ctx).Create(&email)

		if len(attachments) > 0 {
			for i := range attachments {
				attachments[i].EmailID = email.ID
			}
			s.orm.WithContext(ctx).Create(&attachments)
		}
	}

	const prefetchLimit = 20
	start := len(msgs) - prefetchLimit
	if start < 0 {
		start = 0
	}
	for i := len(msgs) - 1; i >= start; i-- {
		msg := msgs[i]
		if msg.Envelope == nil {
			continue
		}
		var row schemas.Email
		if s.orm.WithContext(ctx).Where("folder_id = ? AND imap_uid = ?", folderID, msg.UID).First(&row).Error != nil {
			continue
		}
		if row.BodyText != "" || row.BodyHTML != "" {
			continue
		}
		textBody, htmlBody, err := fetchMessageBodyParts(client, folder.Path, msg.UID, msg.BodyStructure)
		if err != nil || (textBody == "" && htmlBody == "") {
			continue
		}
		s.orm.WithContext(ctx).Model(&row).Updates(map[string]any{
			"body_text": textBody,
			"body_html": htmlBody,
		})
	}

	var unreadCount int64
	s.orm.WithContext(ctx).Model(&schemas.Email{}).Where("folder_id = ? AND is_read = false", folderID).Count(&unreadCount)
	s.orm.WithContext(ctx).Model(&folder).Update("unread_count", unreadCount)

	return nil
}

func (s *Service) GetFolders(ctx context.Context, userID, accountID int64) ([]schemas.Folder, error) {
	if _, err := s.getAccount(ctx, userID, accountID); err != nil {
		return nil, err
	}

	var folders []schemas.Folder
	if err := s.orm.WithContext(ctx).Where("account_id = ?", accountID).Order("type ASC, name ASC").Find(&folders).Error; err != nil {
		return nil, errors.Internal("failed to list folders", err)
	}
	return folders, nil
}

// Conversation is one collapsed thread within a folder: the latest message
// stands in for the group, plus thread-level aggregates the list view needs.
type Conversation struct {
	Email        schemas.Email
	MessageCount int
	UnreadCount  int
	EmailIDs     []int64
}

// conversationRow is the raw scan target: the representative email columns plus
// the window aggregates computed over each conversation partition.
type conversationRow struct {
	schemas.Email
	ConvCount   int    `gorm:"column:conv_count"`
	ConvUnread  int    `gorm:"column:conv_unread"`
	ConvStarred bool   `gorm:"column:conv_starred"`
	ConvAttach  bool   `gorm:"column:conv_attach"`
	ConvIDs     string `gorm:"column:conv_ids"`
}

// conversationPartition groups a folder's messages into threads. Messages with
// no thread id each form their own singleton (keyed by id), so unrelated mail is
// never lumped together.
const conversationPartition = "COALESCE(NULLIF(thread_id, ''), 'msg:' || id::text)"

// GetConversationsByFolderType returns one row per conversation in a folder,
// newest-first, paginated by conversation. Collapsing happens in SQL via window
// functions so pagination stays correct even when a thread spans pages.
func (s *Service) GetConversationsByFolderType(ctx context.Context, userID, accountID int64, folderType string, page, limit int, unreadOnly bool) ([]Conversation, int64, error) {
	if _, err := s.getAccount(ctx, userID, accountID); err != nil {
		return nil, 0, err
	}

	var folder schemas.Folder
	if err := s.orm.WithContext(ctx).Where("account_id = ? AND type = ?", accountID, folderType).First(&folder).Error; err != nil {
		return nil, 0, errors.NotFound(fmt.Sprintf("folder type %q not found", folderType))
	}

	repFilter, havingFilter := "", ""
	if unreadOnly {
		repFilter = " AND agg.conv_unread > 0"
		havingFilter = " HAVING COUNT(*) FILTER (WHERE NOT is_read) > 0"
	}

	var total int64
	countSQL := "SELECT COUNT(*) FROM (SELECT 1 FROM emails WHERE account_id = ? AND folder_id = ? GROUP BY " +
		conversationPartition + havingFilter + ") c"
	if err := s.orm.WithContext(ctx).Raw(countSQL, accountID, folder.ID).Scan(&total).Error; err != nil {
		return nil, 0, errors.Internal("failed to count conversations", err)
	}

	listSQL := `
		WITH scoped AS (
			SELECT id, account_id, folder_id, message_id, thread_id, subject,
			       from_address, from_name, to_addresses, cc_addresses, date,
			       is_read, is_starred, has_attachments, in_reply_to, "references",
			       ` + conversationPartition + ` AS tkey
			FROM emails
			WHERE account_id = ? AND folder_id = ?
		),
		agg AS (
			SELECT tkey,
			       COUNT(*)                            AS conv_count,
			       COUNT(*) FILTER (WHERE NOT is_read) AS conv_unread,
			       BOOL_OR(is_starred)                 AS conv_starred,
			       BOOL_OR(has_attachments)            AS conv_attach,
			       STRING_AGG(id::text, ',')           AS conv_ids,
			       MAX(date)                           AS conv_last
			FROM scoped
			GROUP BY tkey
		),
		rep AS (
			SELECT DISTINCT ON (tkey)
			       id, account_id, folder_id, message_id, thread_id, subject,
			       from_address, from_name, to_addresses, cc_addresses, date,
			       in_reply_to, "references", tkey
			FROM scoped
			ORDER BY tkey, date DESC, id DESC
		)
		SELECT rep.id, rep.account_id, rep.folder_id, rep.message_id, rep.thread_id, rep.subject,
		       rep.from_address, rep.from_name, rep.to_addresses, rep.cc_addresses, rep.date,
		       rep.in_reply_to, rep."references",
		       agg.conv_count, agg.conv_unread, agg.conv_starred, agg.conv_attach, agg.conv_ids
		FROM agg JOIN rep USING (tkey)
		WHERE TRUE` + repFilter + `
		ORDER BY agg.conv_last DESC, rep.id DESC
		LIMIT ? OFFSET ?`

	var rows []conversationRow
	if err := s.orm.WithContext(ctx).Raw(listSQL, accountID, folder.ID, limit, (page-1)*limit).Scan(&rows).Error; err != nil {
		return nil, 0, errors.Internal("failed to list conversations", err)
	}

	conversations := make([]Conversation, len(rows))
	for i, r := range rows {
		e := r.Email
		e.IsRead = r.ConvUnread == 0
		e.IsStarred = r.ConvStarred
		e.HasAttachments = r.ConvAttach
		conversations[i] = Conversation{
			Email:        e,
			MessageCount: r.ConvCount,
			UnreadCount:  r.ConvUnread,
			EmailIDs:     parseCSVInts(r.ConvIDs),
		}
	}
	return conversations, total, nil
}

func (s *Service) SearchEmails(ctx context.Context, userID, accountID int64, query string, page, limit int) ([]schemas.Email, int64, error) {
	if _, err := s.getAccount(ctx, userID, accountID); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	like := "%" + escapeLikePattern(query) + "%"

	q := s.orm.WithContext(ctx).Model(&schemas.Email{}).Where("account_id = ?", accountID).
		Where("subject ILIKE ? OR from_name ILIKE ? OR from_address ILIKE ? OR body_text ILIKE ?", like, like, like, like)

	var excludedIDs []int64
	s.orm.WithContext(ctx).Model(&schemas.Folder{}).
		Where("account_id = ? AND type IN ?", accountID, []string{schemas.FolderTypeTrash, schemas.FolderTypeJunk}).
		Pluck("id", &excludedIDs)
	if len(excludedIDs) > 0 {
		q = q.Where("folder_id NOT IN ?", excludedIDs)
	}

	var total int64
	q.Count(&total)

	var emails []schemas.Email
	if err := q.Order("date DESC").Offset(offset).Limit(limit).Find(&emails).Error; err != nil {
		return nil, 0, errors.Internal("failed to search emails", err)
	}
	return emails, total, nil
}

func (s *Service) GetEmail(ctx context.Context, userID, accountID, emailID int64) (schemas.Email, error) {
	if _, err := s.getAccount(ctx, userID, accountID); err != nil {
		return schemas.Email{}, err
	}

	var email schemas.Email
	if err := s.orm.WithContext(ctx).Where("id = ? AND account_id = ?", emailID, accountID).First(&email).Error; err != nil {
		return schemas.Email{}, errors.NotFound("email not found")
	}

	if email.BodyText == "" && email.BodyHTML == "" {
		account, err := s.getAccountDecrypted(ctx, userID, accountID)
		if err != nil {
			return email, nil
		}
		var folder schemas.Folder
		if s.orm.WithContext(ctx).Where("id = ?", email.FolderID).First(&folder).Error == nil {
			client, err := connectIMAP(account.IMAPHost, account.IMAPPort, account.IMAPUser, account.IMAPPassword)
			if err == nil {
				defer func() {
					client.Logout().Wait()
					client.Close()
				}()
				textBody, htmlBody, err := fetchMessageBodySmart(client, folder.Path, imap.UID(email.IMAPUID))
				if err == nil {
					email.BodyText = textBody
					email.BodyHTML = htmlBody
					s.orm.WithContext(ctx).Model(&email).Updates(map[string]any{
						"body_text": textBody,
						"body_html": htmlBody,
					})
				}
			}
		}
	}

	return email, nil
}

func (s *Service) GetEmailWithAttachments(ctx context.Context, userID, accountID, emailID int64) (schemas.Email, []schemas.Attachment, error) {
	email, err := s.GetEmail(ctx, userID, accountID, emailID)
	if err != nil {
		return schemas.Email{}, nil, err
	}

	var attachments []schemas.Attachment
	s.orm.WithContext(ctx).Where("email_id = ?", email.ID).Find(&attachments)

	return email, attachments, nil
}

// GetThread returns every message in a conversation, oldest-first, across all
// folders of the account. Bodies are whatever is cached; the client lazy-loads
// a message's full body via GetEmail when it expands that message.
//
// The newest 500 messages are kept (huge threads are rare) and presented
// oldest-first, so the reader and "reply to latest" see the true most-recent
// message last.
func (s *Service) GetThread(ctx context.Context, userID, accountID int64, threadID string) ([]schemas.Email, error) {
	if _, err := s.getAccount(ctx, userID, accountID); err != nil {
		return nil, err
	}
	if threadID == "" {
		return nil, errors.Invalid("thread id is required")
	}

	var emails []schemas.Email
	if err := s.orm.WithContext(ctx).
		Where("account_id = ? AND thread_id = ?", accountID, threadID).
		Order("date DESC").
		Limit(500).
		Find(&emails).Error; err != nil {
		return nil, errors.Internal("failed to load thread", err)
	}
	for i, j := 0, len(emails)-1; i < j; i, j = i+1, j-1 {
		emails[i], emails[j] = emails[j], emails[i]
	}
	return emails, nil
}

func (s *Service) DownloadAttachment(w http.ResponseWriter, req *http.Request, userID, accountID, emailID, attachmentID int64) {
	ctx := req.Context()

	account, err := s.getAccountDecrypted(ctx, userID, accountID)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}

	var email schemas.Email
	if err := s.orm.WithContext(ctx).Where("id = ? AND account_id = ?", emailID, accountID).First(&email).Error; err != nil {
		httpjson.WriteError(w, errors.NotFound("email not found"))
		return
	}

	var attachment schemas.Attachment
	if err := s.orm.WithContext(ctx).Where("id = ? AND email_id = ?", attachmentID, emailID).First(&attachment).Error; err != nil {
		httpjson.WriteError(w, errors.NotFound("attachment not found"))
		return
	}

	partNums, err := parsePartID(attachment.PartID)
	if err != nil {
		httpjson.WriteError(w, errors.Internal("invalid part ID", err))
		return
	}

	var folder schemas.Folder
	if err := s.orm.WithContext(ctx).Where("id = ?", email.FolderID).First(&folder).Error; err != nil {
		httpjson.WriteError(w, errors.NotFound("folder not found"))
		return
	}

	client, err := connectIMAP(account.IMAPHost, account.IMAPPort, account.IMAPUser, account.IMAPPassword)
	if err != nil {
		httpjson.WriteError(w, errors.Failed(err.Error()))
		return
	}
	defer func() {
		client.Logout().Wait()
		client.Close()
	}()

	data, err := fetchAttachmentPart(client, folder.Path, imap.UID(email.IMAPUID), partNums)
	if err != nil {
		httpjson.WriteError(w, errors.Internal("failed to fetch attachment", err))
		return
	}

	etag := attachmentETag(email.IMAPUID, attachment.PartID)
	writeAttachmentResponse(w, req, data, attachment.MimeType, attachment.Filename, false, etag)
}

func (s *Service) ServeCIDImage(w http.ResponseWriter, req *http.Request, userID, accountID, emailID int64, cid string) {
	ctx := req.Context()

	account, err := s.getAccountDecrypted(ctx, userID, accountID)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}

	var email schemas.Email
	if err := s.orm.WithContext(ctx).Where("id = ? AND account_id = ?", emailID, accountID).First(&email).Error; err != nil {
		httpjson.WriteError(w, errors.NotFound("email not found"))
		return
	}

	var folder schemas.Folder
	if err := s.orm.WithContext(ctx).Where("id = ?", email.FolderID).First(&folder).Error; err != nil {
		httpjson.WriteError(w, errors.NotFound("folder not found"))
		return
	}

	client, err := connectIMAP(account.IMAPHost, account.IMAPPort, account.IMAPUser, account.IMAPPassword)
	if err != nil {
		httpjson.WriteError(w, errors.Failed(err.Error()))
		return
	}
	defer func() {
		client.Logout().Wait()
		client.Close()
	}()

	raw, err := fetchFullMessage(client, folder.Path, imap.UID(email.IMAPUID))
	if err != nil {
		httpjson.WriteError(w, errors.Internal("failed to fetch message", err))
		return
	}

	data, contentType, filename, ok := resolveInlineImage(raw, cid)
	if !ok {
		httpjson.WriteError(w, errors.NotFound("inline image not found"))
		return
	}
	if filename == "" {
		filename = normalizeCID(cid)
	}

	etag := attachmentETag(email.IMAPUID, normalizeCID(cid))
	writeAttachmentResponse(w, req, data, contentType, filename, true, etag)
}

func (s *Service) UpdateEmail(ctx context.Context, userID, accountID, emailID int64, req UpdateEmailRequest) (schemas.Email, error) {
	account, err := s.getAccountDecrypted(ctx, userID, accountID)
	if err != nil {
		return schemas.Email{}, err
	}

	var email schemas.Email
	if err := s.orm.WithContext(ctx).Where("id = ? AND account_id = ?", emailID, accountID).First(&email).Error; err != nil {
		return schemas.Email{}, errors.NotFound("email not found")
	}

	updates := map[string]any{}
	if req.IsRead != nil {
		updates["is_read"] = *req.IsRead
	}
	if req.IsStarred != nil {
		updates["is_starred"] = *req.IsStarred
	}

	if len(updates) > 0 {
		s.orm.WithContext(ctx).Model(&email).Updates(updates)
		if req.IsRead != nil {
			email.IsRead = *req.IsRead
		}
		if req.IsStarred != nil {
			email.IsStarred = *req.IsStarred
		}
	}

	go func() {
		var folder schemas.Folder
		if s.orm.WithContext(context.Background()).Where("id = ?", email.FolderID).First(&folder).Error != nil {
			return
		}
		client, err := connectIMAP(account.IMAPHost, account.IMAPPort, account.IMAPUser, account.IMAPPassword)
		if err != nil {
			return
		}
		defer func() {
			client.Logout().Wait()
			client.Close()
		}()

		if req.IsRead != nil {
			op := imap.StoreFlagsAdd
			if !*req.IsRead {
				op = imap.StoreFlagsDel
			}
			storeFlags(client, folder.Path, imap.UID(email.IMAPUID), op, []imap.Flag{imap.FlagSeen})
		}
		if req.IsStarred != nil {
			op := imap.StoreFlagsAdd
			if !*req.IsStarred {
				op = imap.StoreFlagsDel
			}
			storeFlags(client, folder.Path, imap.UID(email.IMAPUID), op, []imap.Flag{imap.FlagFlagged})
		}
	}()

	return email, nil
}

// SetReadState marks many emails read/unread at once: one DB update, refreshed
// folder unread counts, and a single IMAP connection that issues one STORE per
// folder. This is what bulk mark-read and opening a conversation use, so a big
// thread no longer triggers one IMAP login per message.
//
// Rows with no IMAP uid are skipped: they are locally-created draft or sent
// rows with no server message to flag.
func (s *Service) SetReadState(ctx context.Context, userID, accountID int64, emailIDs []int64, isRead bool) error {
	account, err := s.getAccountDecrypted(ctx, userID, accountID)
	if err != nil {
		return err
	}
	if len(emailIDs) == 0 {
		return nil
	}

	var emails []schemas.Email
	if err := s.orm.WithContext(ctx).
		Select("id", "folder_id", "imap_uid").
		Where("id IN ? AND account_id = ?", emailIDs, accountID).
		Find(&emails).Error; err != nil {
		return errors.Internal("failed to load emails", err)
	}
	if len(emails) == 0 {
		return nil
	}

	s.orm.WithContext(ctx).Model(&schemas.Email{}).
		Where("id IN ? AND account_id = ?", emailIDs, accountID).
		Update("is_read", isRead)

	affected := map[int64]bool{}
	for _, e := range emails {
		affected[e.FolderID] = true
	}
	for folderID := range affected {
		var unread int64
		s.orm.WithContext(ctx).Model(&schemas.Email{}).Where("folder_id = ? AND is_read = false", folderID).Count(&unread)
		s.orm.WithContext(ctx).Model(&schemas.Folder{}).Where("id = ?", folderID).Update("unread_count", unread)
	}

	uidsByFolder := map[int64][]imap.UID{}
	for _, e := range emails {
		if e.IMAPUID == 0 {
			continue
		}
		uidsByFolder[e.FolderID] = append(uidsByFolder[e.FolderID], imap.UID(e.IMAPUID))
	}

	go func() {
		client, err := connectIMAP(account.IMAPHost, account.IMAPPort, account.IMAPUser, account.IMAPPassword)
		if err != nil {
			return
		}
		defer func() {
			client.Logout().Wait()
			client.Close()
		}()

		paths := map[int64]string{}
		var folders []schemas.Folder
		s.orm.WithContext(context.Background()).Where("account_id = ?", accountID).Find(&folders)
		for _, f := range folders {
			paths[f.ID] = f.Path
		}

		op := imap.StoreFlagsAdd
		if !isRead {
			op = imap.StoreFlagsDel
		}
		for folderID, uids := range uidsByFolder {
			if path := paths[folderID]; path != "" {
				storeFlagsBulk(client, path, uids, op, []imap.Flag{imap.FlagSeen})
			}
		}
	}()

	return nil
}

func (s *Service) Send(ctx context.Context, userID, accountID int64, req SendRequest) error {
	account, err := s.getAccountDecrypted(ctx, userID, accountID)
	if err != nil {
		return err
	}

	if len(req.To) == 0 {
		return errors.Invalid("at least one recipient is required")
	}
	if req.Body == "" && req.BodyHTML == "" {
		return errors.Invalid("body is required")
	}

	msgBytes, messageID, err := buildMessage(
		account.Email,
		account.Name,
		req.To,
		req.Cc,
		req.Subject,
		req.Body,
		req.BodyHTML,
		req.InReplyTo,
		req.References,
		req.Attachments,
	)
	if err != nil {
		return errors.Internal("failed to build message", err)
	}

	allRecipients := make([]string, 0, len(req.To)+len(req.Cc))
	allRecipients = append(allRecipients, extractEmails(req.To)...)
	allRecipients = append(allRecipients, extractEmails(req.Cc)...)

	addr := fmt.Sprintf("%s:%d", account.SMTPHost, account.SMTPPort)
	if account.SMTPPort == 465 {
		err = sendImplicitTLS(addr, account.SMTPHost, account.SMTPUser, account.SMTPPassword, account.Email, allRecipients, msgBytes)
	} else {
		err = sendSTARTTLS(addr, account.SMTPHost, account.SMTPUser, account.SMTPPassword, account.Email, allRecipients, msgBytes)
	}
	if err != nil {
		return errors.Failed(fmt.Sprintf("failed to send: %s", err))
	}

	var sentFolder schemas.Folder
	if s.orm.WithContext(ctx).Where("account_id = ? AND type = ?", accountID, schemas.FolderTypeSent).First(&sentFolder).Error == nil {
		references := strings.Join(req.References, " ")
		email := schemas.Email{
			AccountID:   accountID,
			SpaceID:     account.SpaceID,
			FolderID:    sentFolder.ID,
			MessageID:   messageID,
			ThreadID:    assignThreadID(s.orm.WithContext(ctx), accountID, messageID, req.InReplyTo, references),
			Subject:     req.Subject,
			FromAddress: account.Email,
			FromName:    account.Name,
			ToAddresses: stringsToAddressJSON(req.To),
			CcAddresses: stringsToAddressJSON(req.Cc),
			Date:        time.Now(),
			BodyText:    req.Body,
			BodyHTML:    req.BodyHTML,
			InReplyTo:   req.InReplyTo,
			References:  references,
			IsRead:      true,
		}
		s.orm.WithContext(ctx).Create(&email)
		s.orm.WithContext(ctx).Model(&sentFolder).UpdateColumn("total_count", gorm.Expr("total_count + 1"))
	}

	go s.appendToSent(account, msgBytes)

	return nil
}

func (s *Service) SaveDraft(ctx context.Context, userID, accountID int64, req SendRequest) (schemas.Email, error) {
	account, err := s.getAccountDecrypted(ctx, userID, accountID)
	if err != nil {
		return schemas.Email{}, err
	}

	msgBytes, messageID, err := buildMessage(
		account.Email,
		account.Name,
		req.To,
		req.Cc,
		req.Subject,
		req.Body,
		req.BodyHTML,
		req.InReplyTo,
		req.References,
		nil,
	)
	if err != nil {
		return schemas.Email{}, errors.Internal("failed to build draft message", err)
	}

	var draftsFolder schemas.Folder
	if err := s.orm.WithContext(ctx).Where("account_id = ? AND type = ?", accountID, schemas.FolderTypeDrafts).First(&draftsFolder).Error; err != nil {
		return schemas.Email{}, errors.NotFound("drafts folder not found — sync account first")
	}

	references := strings.Join(req.References, " ")
	email := schemas.Email{
		AccountID:   accountID,
		SpaceID:     account.SpaceID,
		FolderID:    draftsFolder.ID,
		MessageID:   messageID,
		ThreadID:    assignThreadID(s.orm.WithContext(ctx), accountID, messageID, req.InReplyTo, references),
		Subject:     req.Subject,
		FromAddress: account.Email,
		FromName:    account.Name,
		ToAddresses: stringsToAddressJSON(req.To),
		CcAddresses: stringsToAddressJSON(req.Cc),
		Date:        time.Now(),
		BodyText:    req.Body,
		BodyHTML:    req.BodyHTML,
		InReplyTo:   req.InReplyTo,
		References:  references,
		IsRead:      true,
	}
	if err := s.orm.WithContext(ctx).Create(&email).Error; err != nil {
		return schemas.Email{}, errors.Internal("failed to save draft", err)
	}

	go s.appendToDrafts(account, msgBytes)

	return email, nil
}

func (s *Service) DeleteDraft(ctx context.Context, userID, accountID, emailID int64) error {
	if _, err := s.getAccount(ctx, userID, accountID); err != nil {
		return err
	}

	var draftsFolder schemas.Folder
	if err := s.orm.WithContext(ctx).Where("account_id = ? AND type = ?", accountID, schemas.FolderTypeDrafts).First(&draftsFolder).Error; err != nil {
		return errors.NotFound("drafts folder not found")
	}

	var email schemas.Email
	if err := s.orm.WithContext(ctx).Where("id = ? AND account_id = ? AND folder_id = ?", emailID, accountID, draftsFolder.ID).First(&email).Error; err != nil {
		return errors.NotFound("draft not found")
	}

	s.orm.WithContext(ctx).Delete(&email)
	return nil
}

func (s *Service) MoveEmails(ctx context.Context, userID, accountID int64, emailIDs []int64, destFolderType string) error {
	account, err := s.getAccountDecrypted(ctx, userID, accountID)
	if err != nil {
		return err
	}

	var destFolder schemas.Folder
	if err := s.orm.WithContext(ctx).Where("account_id = ? AND type = ?", accountID, destFolderType).First(&destFolder).Error; err != nil {
		return errors.NotFound(fmt.Sprintf("destination folder %q not found", destFolderType))
	}

	var emails []schemas.Email
	if err := s.orm.WithContext(ctx).Where("id IN ? AND account_id = ?", emailIDs, accountID).Find(&emails).Error; err != nil {
		return errors.Internal("failed to find emails", err)
	}

	if len(emails) == 0 {
		return errors.NotFound("no emails found")
	}

	srcFolderIDs := map[int64]bool{}
	for _, e := range emails {
		srcFolderIDs[e.FolderID] = true
	}

	s.orm.WithContext(ctx).Model(&schemas.Email{}).Where("id IN ? AND account_id = ?", emailIDs, accountID).Update("folder_id", destFolder.ID)

	for folderID := range srcFolderIDs {
		var count int64
		s.orm.WithContext(ctx).Model(&schemas.Email{}).Where("folder_id = ? AND is_read = false", folderID).Count(&count)
		s.orm.WithContext(ctx).Model(&schemas.Folder{}).Where("id = ?", folderID).Update("unread_count", count)
	}
	var destUnread int64
	s.orm.WithContext(ctx).Model(&schemas.Email{}).Where("folder_id = ? AND is_read = false", destFolder.ID).Count(&destUnread)
	s.orm.WithContext(ctx).Model(&destFolder).Update("unread_count", destUnread)

	go func() {
		client, err := connectIMAP(account.IMAPHost, account.IMAPPort, account.IMAPUser, account.IMAPPassword)
		if err != nil {
			return
		}
		defer func() {
			client.Logout().Wait()
			client.Close()
		}()

		mailboxes, err := listMailboxes(client)
		if err != nil {
			return
		}
		var destMailbox string
		for _, mbox := range mailboxes {
			if detectFolderType(mbox) == destFolderType {
				destMailbox = mbox.Mailbox
				break
			}
		}
		if destMailbox == "" {
			return
		}

		folderPathCache := map[int64]string{}
		for _, e := range emails {
			srcPath, ok := folderPathCache[e.FolderID]
			if !ok {
				var f schemas.Folder
				if s.orm.WithContext(context.Background()).Where("id = ?", e.FolderID).First(&f).Error == nil {
					srcPath = f.Path
					folderPathCache[e.FolderID] = srcPath
				}
			}
			if srcPath != "" {
				moveMessage(client, srcPath, imap.UID(e.IMAPUID), destMailbox)
			}
		}
	}()

	return nil
}

func (s *Service) DeleteEmails(ctx context.Context, userID, accountID int64, emailIDs []int64) error {
	if _, err := s.getAccount(ctx, userID, accountID); err != nil {
		return err
	}

	var trashFolder schemas.Folder
	if err := s.orm.WithContext(ctx).Where("account_id = ? AND type = ?", accountID, schemas.FolderTypeTrash).First(&trashFolder).Error; err != nil {
		return errors.NotFound("trash folder not found — sync account first")
	}

	var emails []schemas.Email
	if err := s.orm.WithContext(ctx).Where("id IN ? AND account_id = ?", emailIDs, accountID).Find(&emails).Error; err != nil {
		return errors.Internal("failed to find emails", err)
	}

	alreadyInTrash := []schemas.Email{}
	notInTrash := []int64{}
	for _, e := range emails {
		if e.FolderID == trashFolder.ID {
			alreadyInTrash = append(alreadyInTrash, e)
		} else {
			notInTrash = append(notInTrash, e.ID)
		}
	}

	if len(alreadyInTrash) > 0 {
		ids := make([]int64, len(alreadyInTrash))
		for i, e := range alreadyInTrash {
			ids[i] = e.ID
		}
		s.orm.WithContext(ctx).Where("id IN ?", ids).Delete(&schemas.Email{})
		s.orm.WithContext(ctx).Where("email_id IN ?", ids).Delete(&schemas.Attachment{})
	}

	if len(notInTrash) > 0 {
		return s.MoveEmails(ctx, userID, accountID, notInTrash, schemas.FolderTypeTrash)
	}

	return nil
}

func (s *Service) ArchiveEmails(ctx context.Context, userID, accountID int64, emailIDs []int64) error {
	return s.MoveEmails(ctx, userID, accountID, emailIDs, schemas.FolderTypeArchive)
}

func (s *Service) ListTemplates(ctx context.Context, userID int64, spaceID *string) ([]schemas.EmailTemplate, error) {
	var templates []schemas.EmailTemplate
	query := s.orm.WithContext(ctx).Where("user_id = ?", userID)
	if spaceID != nil {
		query = query.Where("space_id = ?", *spaceID)
	} else {
		query = query.Where("space_id IS NULL")
	}
	if err := query.Order("updated_at DESC").Find(&templates).Error; err != nil {
		return nil, errors.Internal("failed to list templates", err)
	}
	return templates, nil
}

func (s *Service) CreateTemplate(ctx context.Context, userID int64, req EmailTemplateRequest) (schemas.EmailTemplate, error) {
	if req.Name == "" {
		return schemas.EmailTemplate{}, errors.Invalid("template name is required")
	}

	var count int64
	s.orm.WithContext(ctx).Model(&schemas.EmailTemplate{}).Where("user_id = ?", userID).Count(&count)
	if count >= 50 {
		return schemas.EmailTemplate{}, errors.Invalid("template limit reached (max 50)")
	}

	tmpl := schemas.EmailTemplate{
		UserID:   userID,
		SpaceID:  req.SpaceID,
		Name:     req.Name,
		Subject:  req.Subject,
		BodyHTML: req.BodyHTML,
		BodyText: req.BodyText,
	}
	if err := s.orm.WithContext(ctx).Create(&tmpl).Error; err != nil {
		return schemas.EmailTemplate{}, errors.Internal("failed to create template", err)
	}
	return tmpl, nil
}

func (s *Service) UpdateTemplate(ctx context.Context, userID, templateID int64, req EmailTemplateRequest) (schemas.EmailTemplate, error) {
	var tmpl schemas.EmailTemplate
	if err := s.orm.WithContext(ctx).Where("id = ? AND user_id = ?", templateID, userID).First(&tmpl).Error; err != nil {
		return schemas.EmailTemplate{}, errors.NotFound("template not found")
	}

	updates := map[string]any{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	updates["subject"] = req.Subject
	updates["body_html"] = req.BodyHTML
	updates["body_text"] = req.BodyText

	if err := s.orm.WithContext(ctx).Model(&tmpl).Updates(updates).Error; err != nil {
		return schemas.EmailTemplate{}, errors.Internal("failed to update template", err)
	}

	s.orm.WithContext(ctx).Where("id = ?", templateID).First(&tmpl)
	return tmpl, nil
}

func (s *Service) DeleteTemplate(ctx context.Context, userID, templateID int64) error {
	var tmpl schemas.EmailTemplate
	if err := s.orm.WithContext(ctx).Where("id = ? AND user_id = ?", templateID, userID).First(&tmpl).Error; err != nil {
		return errors.NotFound("template not found")
	}
	s.orm.WithContext(ctx).Delete(&tmpl)
	return nil
}

func (s *Service) appendToDrafts(account schemas.Account, msgBytes []byte) {
	client, err := connectIMAP(account.IMAPHost, account.IMAPPort, account.IMAPUser, account.IMAPPassword)
	if err != nil {
		return
	}
	defer func() {
		client.Logout().Wait()
		client.Close()
	}()

	var draftsFolder string
	mailboxes, err := listMailboxes(client)
	if err != nil {
		return
	}
	for _, mbox := range mailboxes {
		if detectFolderType(mbox) == schemas.FolderTypeDrafts {
			draftsFolder = mbox.Mailbox
			break
		}
	}
	if draftsFolder == "" {
		return
	}

	appendCmd := client.Append(draftsFolder, int64(len(msgBytes)), &imap.AppendOptions{
		Flags: []imap.Flag{imap.FlagDraft, imap.FlagSeen},
	})
	if _, err := appendCmd.Write(msgBytes); err != nil {
		appendCmd.Close()
		return
	}
	appendCmd.Close()
	appendCmd.Wait()
}

func (s *Service) appendToSent(account schemas.Account, msgBytes []byte) {
	client, err := connectIMAP(account.IMAPHost, account.IMAPPort, account.IMAPUser, account.IMAPPassword)
	if err != nil {
		return
	}
	defer func() {
		client.Logout().Wait()
		client.Close()
	}()

	var sentFolder string
	mailboxes, err := listMailboxes(client)
	if err != nil {
		return
	}
	for _, mbox := range mailboxes {
		if detectFolderType(mbox) == schemas.FolderTypeSent {
			sentFolder = mbox.Mailbox
			break
		}
	}
	if sentFolder == "" {
		return
	}

	appendCmd := client.Append(sentFolder, int64(len(msgBytes)), &imap.AppendOptions{
		Flags: []imap.Flag{imap.FlagSeen},
	})
	if _, err := appendCmd.Write(msgBytes); err != nil {
		appendCmd.Close()
		return
	}
	appendCmd.Close()
	appendCmd.Wait()
}

func (s *Service) TestConnection(ctx context.Context, req TestConnectionRequest) error {
	if req.IMAPHost != "" {
		port := req.IMAPPort
		if port == 0 {
			port = 993
		}
		client, err := connectIMAP(req.IMAPHost, port, req.IMAPUser, req.IMAPPassword)
		if err != nil {
			return errors.Failed(fmt.Sprintf("IMAP: %s", err))
		}
		client.Logout().Wait()
		client.Close()
	}

	if req.SMTPHost != "" {
		port := req.SMTPPort
		if port == 0 {
			port = 587
		}
		if err := testSMTP(req.SMTPHost, port, req.SMTPUser, req.SMTPPassword); err != nil {
			return errors.Failed(fmt.Sprintf("SMTP: %s", err))
		}
	}

	return nil
}

func (s *Service) SearchContacts(ctx context.Context, userID, accountID int64, query string) ([]ContactResult, error) {
	if _, err := s.getAccount(ctx, userID, accountID); err != nil {
		return nil, err
	}

	pattern := "%" + escapeLikePattern(strings.ToLower(query)) + "%"

	type fromRow struct {
		FromAddress string
		FromName    string
		Cnt         int
	}
	var fromRows []fromRow
	s.orm.WithContext(ctx).
		Model(&schemas.Email{}).
		Select("from_address, from_name, COUNT(*) as cnt").
		Where("account_id = ? AND (LOWER(from_address) LIKE ? OR LOWER(from_name) LIKE ?)", accountID, pattern, pattern).
		Group("from_address, from_name").
		Order("cnt DESC").
		Limit(20).
		Find(&fromRows)

	seen := map[string]*ContactResult{}
	for _, r := range fromRows {
		key := strings.ToLower(r.FromAddress)
		if key == "" {
			continue
		}
		if existing, ok := seen[key]; ok {
			existing.Count += r.Cnt
			if existing.Name == "" && r.FromName != "" {
				existing.Name = r.FromName
			}
		} else {
			seen[key] = &ContactResult{
				Name:  r.FromName,
				Email: r.FromAddress,
				Count: r.Cnt,
			}
		}
	}

	var recipientEmails []schemas.Email
	s.orm.WithContext(ctx).
		Where("account_id = ? AND (LOWER(to_addresses) LIKE ? OR LOWER(cc_addresses) LIKE ?)", accountID, pattern, pattern).
		Select("to_addresses, cc_addresses").
		Limit(100).
		Find(&recipientEmails)

	for _, email := range recipientEmails {
		for _, raw := range []string{email.ToAddresses, email.CcAddresses} {
			addrs := parseAddressJSON(raw)
			for _, addr := range addrs {
				if addr.Email == "" {
					continue
				}
				lower := strings.ToLower(addr.Email)
				if !strings.Contains(lower, strings.ToLower(query)) && !strings.Contains(strings.ToLower(addr.Name), strings.ToLower(query)) {
					continue
				}
				if existing, ok := seen[lower]; ok {
					existing.Count++
					if existing.Name == "" && addr.Name != "" {
						existing.Name = addr.Name
					}
				} else {
					seen[lower] = &ContactResult{
						Name:  addr.Name,
						Email: addr.Email,
						Count: 1,
					}
				}
			}
		}
	}

	results := make([]ContactResult, 0, len(seen))
	for _, c := range seen {
		results = append(results, *c)
	}

	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].Count > results[i].Count {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	if len(results) > 10 {
		results = results[:10]
	}

	return results, nil
}

func isNoSelect(mbox *imap.ListData) bool {
	for _, attr := range mbox.Attrs {
		if attr == imap.MailboxAttrNoSelect {
			return true
		}
	}
	return false
}

func containsFlag(flags []imap.Flag, target imap.Flag) bool {
	for _, f := range flags {
		if f == target {
			return true
		}
	}
	return false
}

func firstAddr(addrs []imap.Address) string {
	if len(addrs) == 0 {
		return ""
	}
	return addrs[0].Addr()
}

func firstName(addrs []imap.Address) string {
	if len(addrs) == 0 {
		return ""
	}
	return addrs[0].Name
}

func emailToResponse(e schemas.Email, attachments ...schemas.Attachment) EmailResponse {
	resp := EmailResponse{
		ID:             e.ID,
		AccountID:      e.AccountID,
		FolderID:       e.FolderID,
		MessageID:      e.MessageID,
		ThreadID:       e.ThreadID,
		Subject:        e.Subject,
		FromAddress:    e.FromAddress,
		FromName:       e.FromName,
		Date:           e.Date.Format(time.RFC3339),
		BodyText:       e.BodyText,
		BodyHTML:       e.BodyHTML,
		IsRead:         e.IsRead,
		IsStarred:      e.IsStarred,
		HasAttachments: e.HasAttachments,
		InReplyTo:      e.InReplyTo,
		References:     e.References,
	}

	resp.ToAddresses = parseAddressJSON(e.ToAddresses)
	resp.CcAddresses = parseAddressJSON(e.CcAddresses)

	if len(attachments) > 0 {
		resp.Attachments = make([]AttachmentResponse, len(attachments))
		for i, a := range attachments {
			resp.Attachments[i] = AttachmentResponse{
				ID:       a.ID,
				Filename: a.Filename,
				MimeType: a.MimeType,
				Size:     a.Size,
			}
		}
	}

	return resp
}

func stringsToAddressJSON(addrs []string) string {
	if len(addrs) == 0 {
		return "[]"
	}
	type entry struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	items := make([]entry, len(addrs))
	for i, addr := range addrs {
		if parsed, err := mail.ParseAddress(strings.TrimSpace(addr)); err == nil {
			items[i] = entry{Name: parsed.Name, Email: parsed.Address}
		} else {
			items[i] = entry{Email: addr}
		}
	}
	data, _ := json.Marshal(items)
	return string(data)
}

func parseAddressJSON(raw string) []AddressResponse {
	if raw == "" || raw == "[]" {
		return []AddressResponse{}
	}
	var addrs []AddressResponse
	if err := json.Unmarshal([]byte(raw), &addrs); err != nil {
		return []AddressResponse{}
	}
	return addrs
}

func parseCSVInts(s string) []int64 {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]int64, 0, len(parts))
	for _, p := range parts {
		if n, err := strconv.ParseInt(strings.TrimSpace(p), 10, 64); err == nil {
			out = append(out, n)
		}
	}
	return out
}

func escapeLikePattern(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "%", "\\%")
	s = strings.ReplaceAll(s, "_", "\\_")
	return s
}

func folderToResponse(f schemas.Folder) FolderResponse {
	return FolderResponse{
		ID:          f.ID,
		AccountID:   f.AccountID,
		Path:        f.Path,
		Name:        f.Name,
		Type:        f.Type,
		UnreadCount: f.UnreadCount,
		TotalCount:  f.TotalCount,
	}
}
