package mail

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/emersion/go-imap/v2"

	"github.com/FacileStudio/Courrier/apps/api/internal/crypto"
	"github.com/FacileStudio/Courrier/apps/api/internal/push"
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
	backfilling   sync.Map
	cooldowns     sync.Map
	dial          imapDialer
	deliver       smtpDeliverer
	notifier      *push.Hub
	folderLocks   sync.Map
}

// NewService builds a mail Service over the database and the app's encryption
// key.
func NewService(orm *gorm.DB, encryptionKey []byte) *Service {
	return &Service{orm: orm, encryptionKey: encryptionKey, dial: connectIMAP, deliver: deliverSMTP}
}

func (s *Service) decryptPassword(cipher string) (string, error) {
	if len(s.encryptionKey) == 0 || cipher == "" {
		return cipher, nil
	}
	return crypto.Decrypt(cipher, s.encryptionKey)
}

func (s *Service) requireSpaceMember(ctx context.Context, spaceID string, userID int64) error {
	var count int64
	if err := s.orm.WithContext(ctx).Model(&schemas.SpaceMember{}).
		Where("space_id = ? AND user_id = ?", spaceID, userID).
		Count(&count).Error; err != nil {
		return errors.Internal("failed to check space membership", err)
	}
	if count == 0 {
		return errors.Forbidden("not a member of this space")
	}
	return nil
}

func (s *Service) getAccount(ctx context.Context, userID, accountID int64) (schemas.Account, error) {
	var account schemas.Account
	if err := s.orm.WithContext(ctx).Where("id = ?", accountID).First(&account).Error; err != nil {
		return schemas.Account{}, errors.NotFound("account not found")
	}
	if account.SpaceID != nil && *account.SpaceID != "" {
		if err := s.requireSpaceMember(ctx, *account.SpaceID, userID); err != nil {
			return schemas.Account{}, err
		}
	} else if account.UserID != userID {
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

	client, err := s.connect(account)
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

// SyncFolderEmails caches the newest messages of one folder from IMAP.
func (s *Service) SyncFolderEmails(ctx context.Context, userID, accountID, folderID int64) error {
	_, _, err := s.syncFolder(ctx, userID, accountID, folderID)
	return err
}

// syncFolder is SyncFolderEmails returning the account and the messages it
// cached for the first time. fresh stays empty on a folder's first sync and
// after a UIDVALIDITY reset, when every message is new to the cache but not
// to the user.
func (s *Service) syncFolder(ctx context.Context, userID, accountID, folderID int64) (schemas.Account, []schemas.Email, error) {
	unlock := s.lockFolder(folderID)
	defer unlock()

	account, err := s.getAccountDecrypted(ctx, userID, accountID)
	if err != nil {
		return account, nil, err
	}

	var folder schemas.Folder
	if err := s.orm.WithContext(ctx).Where("id = ? AND account_id = ?", folderID, accountID).First(&folder).Error; err != nil {
		return account, nil, errors.NotFound("folder not found")
	}
	var cached int64
	s.orm.WithContext(ctx).Model(&schemas.Email{}).Where("folder_id = ?", folderID).Limit(1).Count(&cached)
	var fresh []schemas.Email

	client, err := s.connect(account)
	if err != nil {
		return account, nil, errors.Failed(err.Error())
	}
	defer func() {
		client.Logout().Wait()
		client.Close()
	}()

	msgs, selectData, err := fetchEnvelopes(client, folder.Path, 100)
	if err != nil {
		return account, nil, errors.Internal("failed to fetch emails", err)
	}

	if selectData != nil && selectData.UIDValidity != folder.UIDValidity && folder.UIDValidity != 0 {
		s.orm.WithContext(ctx).Where("folder_id = ?", folderID).Delete(&schemas.Email{})
		cached = 0
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
		category, bulkReason := classifyForStore(bulkSignals(msg), firstAddr(env.From))

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
			BccAddresses:   imapAddressesToJSON(env.Bcc),
			Date:           env.Date,
			IsRead:         containsFlag(msg.Flags, imap.FlagSeen),
			IsStarred:      containsFlag(msg.Flags, imap.FlagFlagged),
			HasAttachments: len(attachments) > 0,
			InReplyTo:      inReplyTo,
			References:     references,
			Category:       category,
			BulkReason:     bulkReason,
			IMAPUID:        uint32(msg.UID),
		}

		s.orm.WithContext(ctx).Create(&email)
		s.recordSyncedSent(ctx, folder, account.SpaceID, env)
		if cached > 0 && email.ID != 0 {
			fresh = append(fresh, email)
		}

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
	s.backfillAfterSync(userID, accountID)

	return account, fresh, nil
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

	filter := conversationFilter{FolderIDs: []int64{folder.ID}, UnreadOnly: unreadOnly}
	total, err := s.countConversations(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	conversations, err := s.queryConversations(ctx, filter, limit, (page-1)*limit)
	return conversations, total, err
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
			client, err := s.connect(account)
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

	client, err := s.connect(account)
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

	client, err := s.connect(account)
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
		client, err := s.connect(account)
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

func (s *Service) DeleteEmails(ctx context.Context, userID, accountID int64, emailIDs []int64) error {
	account, err := s.getAccountDecrypted(ctx, userID, accountID)
	if err != nil {
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

	if err := s.purgeFromTrash(ctx, account, trashFolder, alreadyInTrash); err != nil {
		return err
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
	spaceID, err := s.templateSpace(ctx, userID, spaceID)
	if err != nil {
		return nil, err
	}
	var templates []schemas.EmailTemplate
	query := s.orm.WithContext(ctx)
	if spaceID != nil {
		query = query.Where("space_id = ?", *spaceID)
	} else {
		query = query.Where("user_id = ? AND space_id IS NULL", userID)
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

	spaceID, err := s.templateSpace(ctx, userID, req.SpaceID)
	if err != nil {
		return schemas.EmailTemplate{}, err
	}
	tmpl := schemas.EmailTemplate{
		UserID:   userID,
		SpaceID:  spaceID,
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

func (s *Service) getTemplate(ctx context.Context, userID, templateID int64) (schemas.EmailTemplate, error) {
	var tmpl schemas.EmailTemplate
	if err := s.orm.WithContext(ctx).Where("id = ?", templateID).First(&tmpl).Error; err != nil {
		return schemas.EmailTemplate{}, errors.NotFound("template not found")
	}
	if tmpl.SpaceID != nil && *tmpl.SpaceID != "" {
		if err := s.requireSpaceMember(ctx, *tmpl.SpaceID, userID); err != nil {
			return schemas.EmailTemplate{}, err
		}
	} else if tmpl.UserID != userID {
		return schemas.EmailTemplate{}, errors.NotFound("template not found")
	}
	return tmpl, nil
}

func (s *Service) UpdateTemplate(ctx context.Context, userID, templateID int64, req EmailTemplateRequest) (schemas.EmailTemplate, error) {
	tmpl, err := s.getTemplate(ctx, userID, templateID)
	if err != nil {
		return schemas.EmailTemplate{}, err
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
	tmpl, err := s.getTemplate(ctx, userID, templateID)
	if err != nil {
		return err
	}
	s.orm.WithContext(ctx).Delete(&tmpl)
	return nil
}

func (s *Service) appendToSent(account schemas.Account, msgBytes []byte) {
	client, err := s.connect(account)
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
		Category:       e.Category,
		BulkReason:     visibleBulkReason(e.BulkReason),
	}

	resp.ToAddresses = parseAddressJSON(e.ToAddresses)
	resp.CcAddresses = parseAddressJSON(e.CcAddresses)
	resp.BccAddresses = parseAddressJSON(e.BccAddresses)

	resp.Attachments = attachmentResponses(attachments)

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

func attachmentResponses(attachments []schemas.Attachment) []AttachmentResponse {
	if len(attachments) == 0 {
		return nil
	}
	out := make([]AttachmentResponse, len(attachments))
	for i, a := range attachments {
		out[i] = AttachmentResponse{ID: a.ID, Filename: a.Filename, MimeType: a.MimeType, Size: a.Size}
	}
	return out
}
