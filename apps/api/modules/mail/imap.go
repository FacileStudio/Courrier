package mail

import (
	"bufio"
	"bytes"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"io"
	"net/textproto"
	"net/url"
	"strings"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	gomessage "github.com/emersion/go-message"
	// Registers decoders for non-UTF-8 charsets (iso-8859-1, windows-1252, …)
	// by setting gomessage.CharsetReader; without this, decodePartContent and
	// go-imap envelope parsing fall back to raw bytes and mangle accents.
	_ "github.com/emersion/go-message/charset"

	"api/schemas"
)

// referencesSection fetches just the References header alongside envelopes — the
// envelope exposes In-Reply-To but not References, which is the load-bearing
// header for threading. Fetching only this field keeps the sync cheap.
var referencesSection = &imap.FetchItemBodySection{
	Specifier:    imap.PartSpecifierHeader,
	HeaderFields: []string{"References"},
	Peek:         true,
}

// messageReferences pulls the (whitespace-collapsed) References header out of a
// fetched message, or "" when the sender omitted it.
func messageReferences(msg *imapclient.FetchMessageBuffer) string {
	raw := msg.FindBodySection(referencesSection)
	if len(raw) == 0 {
		return ""
	}
	hdr, _ := textproto.NewReader(bufio.NewReader(bytes.NewReader(raw))).ReadMIMEHeader()
	return strings.Join(strings.Fields(hdr.Get("References")), " ")
}

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
		BodySection: []*imap.FetchItemBodySection{referencesSection},
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

type textPartRef struct {
	path     []int
	encoding string
	charset  string
}

func selectTextParts(bs imap.BodyStructure) (textPlain, htmlPart *textPartRef) {
	if bs == nil {
		return nil, nil
	}
	if _, ok := bs.(*imap.BodyStructureSinglePart); ok {
		return nil, nil
	}
	bs.Walk(func(path []int, part imap.BodyStructure) bool {
		sp, ok := part.(*imap.BodyStructureSinglePart)
		if !ok {
			return true
		}
		if sp.Filename() != "" {
			return true
		}
		if disp := sp.Disposition(); disp != nil && strings.EqualFold(disp.Value, "attachment") {
			return true
		}
		newRef := func() *textPartRef {
			p := make([]int, len(path))
			copy(p, path)
			return &textPartRef{path: p, encoding: sp.Encoding, charset: sp.Params["charset"]}
		}
		switch sp.MediaType() {
		case "text/plain":
			if textPlain == nil {
				textPlain = newRef()
			}
		case "text/html":
			if htmlPart == nil {
				htmlPart = newRef()
			}
		}
		return true
	})
	return textPlain, htmlPart
}

func decodePartContent(data []byte, encoding, charset string) string {
	decoded, err := decodeTransferEncoding(data, encoding)
	if err != nil {
		decoded = data
	}
	cs := strings.ToLower(strings.TrimSpace(charset))
	if cs == "" || cs == "utf-8" || cs == "us-ascii" {
		return string(decoded)
	}
	if gomessage.CharsetReader != nil {
		if r, err := gomessage.CharsetReader(cs, bytes.NewReader(decoded)); err == nil {
			if converted, err := io.ReadAll(r); err == nil {
				return string(converted)
			}
		}
	}
	return string(decoded)
}

func fetchMessageBodyParts(client *imapclient.Client, mailbox string, uid imap.UID, bs imap.BodyStructure) (string, string, error) {
	textRef, htmlRef := selectTextParts(bs)
	if textRef == nil && htmlRef == nil {
		return fetchMessageBody(client, mailbox, uid)
	}

	if _, err := client.Select(mailbox, nil).Wait(); err != nil {
		return "", "", fmt.Errorf("SELECT %q failed: %w", mailbox, err)
	}

	var sections []*imap.FetchItemBodySection
	var textSection, htmlSection *imap.FetchItemBodySection
	if textRef != nil {
		textSection = &imap.FetchItemBodySection{Part: textRef.path, Peek: true}
		sections = append(sections, textSection)
	}
	if htmlRef != nil {
		htmlSection = &imap.FetchItemBodySection{Part: htmlRef.path, Peek: true}
		sections = append(sections, htmlSection)
	}

	uidSet := imap.UIDSetNum(uid)
	fetchCmd := client.Fetch(uidSet, &imap.FetchOptions{
		BodySection: sections,
	})
	msgs, err := fetchCmd.Collect()
	if err != nil {
		return "", "", fmt.Errorf("FETCH body parts failed: %w", err)
	}
	if len(msgs) == 0 {
		return "", "", fmt.Errorf("message not found")
	}

	msg := msgs[0]
	var textBody, htmlBody string
	if textSection != nil {
		if data := msg.FindBodySection(textSection); data != nil {
			textBody = decodePartContent(data, textRef.encoding, textRef.charset)
		}
	}
	if htmlSection != nil {
		if data := msg.FindBodySection(htmlSection); data != nil {
			htmlBody = decodePartContent(data, htmlRef.encoding, htmlRef.charset)
		}
	}

	if textBody == "" && htmlBody == "" {
		return fetchMessageBody(client, mailbox, uid)
	}
	return textBody, htmlBody, nil
}

func fetchMessageBodySmart(client *imapclient.Client, mailbox string, uid imap.UID) (string, string, error) {
	if _, err := client.Select(mailbox, nil).Wait(); err != nil {
		return fetchMessageBody(client, mailbox, uid)
	}

	uidSet := imap.UIDSetNum(uid)
	fetchCmd := client.Fetch(uidSet, &imap.FetchOptions{
		BodyStructure: &imap.FetchItemBodyStructure{Extended: true},
	})
	msgs, err := fetchCmd.Collect()
	if err != nil || len(msgs) == 0 || msgs[0].BodyStructure == nil {
		return fetchMessageBody(client, mailbox, uid)
	}

	textBody, htmlBody, err := fetchMessageBodyParts(client, mailbox, uid, msgs[0].BodyStructure)
	if err != nil {
		return fetchMessageBody(client, mailbox, uid)
	}
	return textBody, htmlBody, nil
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

// extractAttachments lists the downloadable parts of a message: those with a
// filename or an explicit attachment disposition. Inline parts referenced only
// by Content-ID are served on demand by ServeCIDImage and are not listed here.
func extractAttachments(bs imap.BodyStructure, emailID int64) []schemas.Attachment {
	if bs == nil {
		return nil
	}
	var attachments []schemas.Attachment
	bs.Walk(func(path []int, part imap.BodyStructure) bool {
		sp, ok := part.(*imap.BodyStructureSinglePart)
		if !ok {
			return true
		}

		filename := sp.Filename()
		disp := sp.Disposition()
		if filename == "" && (disp == nil || !strings.EqualFold(disp.Value, "attachment")) {
			return true
		}
		if filename == "" {
			filename = "unnamed"
		}

		partNums := make([]string, len(path))
		for i, n := range path {
			partNums[i] = fmt.Sprintf("%d", n)
		}

		attachments = append(attachments, schemas.Attachment{
			EmailID:  emailID,
			Filename: filename,
			MimeType: sp.MediaType(),
			Size:     int64(sp.Size),
			PartID:   strings.Join(partNums, "."),
		})
		return true
	})
	return attachments
}

func fetchAttachmentPart(client *imapclient.Client, mailbox string, uid imap.UID, partNums []int) ([]byte, error) {
	_, err := client.Select(mailbox, nil).Wait()
	if err != nil {
		return nil, fmt.Errorf("SELECT %q failed: %w", mailbox, err)
	}

	uidSet := imap.UIDSetNum(uid)
	bodySection := &imap.FetchItemBodySection{
		Part: partNums,
		Peek: true,
	}
	fetchCmd := client.Fetch(uidSet, &imap.FetchOptions{
		BodyStructure: &imap.FetchItemBodyStructure{Extended: true},
		BodySection:   []*imap.FetchItemBodySection{bodySection},
	})
	msgs, err := fetchCmd.Collect()
	if err != nil {
		return nil, fmt.Errorf("FETCH part failed: %w", err)
	}
	if len(msgs) == 0 {
		return nil, fmt.Errorf("message not found")
	}

	data := msgs[0].FindBodySection(bodySection)
	if data == nil {
		return nil, fmt.Errorf("body section not found")
	}

	return decodeTransferEncoding(data, partEncoding(msgs[0].BodyStructure, partNums))
}

const maxInlineImageSize = 25 << 20

// fetchFullMessage retrieves the entire raw RFC822 message (BODY.PEEK[]).
func fetchFullMessage(client *imapclient.Client, mailbox string, uid imap.UID) ([]byte, error) {
	if _, err := client.Select(mailbox, nil).Wait(); err != nil {
		return nil, fmt.Errorf("SELECT %q failed: %w", mailbox, err)
	}

	uidSet := imap.UIDSetNum(uid)
	section := &imap.FetchItemBodySection{Peek: true}
	fetchCmd := client.Fetch(uidSet, &imap.FetchOptions{
		BodySection: []*imap.FetchItemBodySection{section},
	})
	msgs, err := fetchCmd.Collect()
	if err != nil {
		return nil, fmt.Errorf("FETCH message failed: %w", err)
	}
	if len(msgs) == 0 {
		return nil, fmt.Errorf("message not found")
	}

	data := msgs[0].FindBodySection(section)
	if data == nil {
		return nil, fmt.Errorf("message body not found")
	}
	return data, nil
}

// normalizeCID reduces an HTML cid: token or a MIME Content-ID header value to
// a comparable form per RFC 2392: drop the cid: prefix, the angle brackets that
// both go-imap and go-message preserve, surrounding whitespace, and %hh escapes.
func normalizeCID(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "cid:")
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "<>")
	s = strings.TrimSpace(s)
	if decoded, err := url.PathUnescape(s); err == nil {
		s = decoded
	}
	return strings.TrimSpace(s)
}

// resolveInlineImage parses a raw message and returns the decoded bytes of the
// part referenced by a cid: URL. It walks the whole MIME tree (not just
// multipart/related) and matches by Content-ID first, then falls back to
// filename / Content-Type name / Content-Location for senders that reference
// inline images loosely (Outlook, Word "save as web page").
func resolveInlineImage(raw []byte, cid string) (data []byte, contentType, filename string, ok bool) {
	want := normalizeCID(cid)
	if want == "" {
		return nil, "", "", false
	}

	if d, ct, fn, found := walkInlineParts(raw, func(h *gomessage.Header) bool {
		return strings.EqualFold(normalizeCID(h.Get("Content-Id")), want)
	}); found {
		return d, ct, fn, true
	}

	return walkInlineParts(raw, func(h *gomessage.Header) bool {
		_, ctParams, _ := h.ContentType()
		_, dispParams, _ := h.ContentDisposition()
		return strings.EqualFold(dispParams["filename"], want) ||
			strings.EqualFold(ctParams["name"], want) ||
			strings.EqualFold(strings.TrimSpace(h.Get("Content-Location")), want)
	})
}

func walkInlineParts(raw []byte, match func(*gomessage.Header) bool) ([]byte, string, string, bool) {
	entity, _ := gomessage.Read(bytes.NewReader(raw))
	if entity == nil {
		return nil, "", "", false
	}

	var (
		outData []byte
		outType string
		outName string
		found   bool
	)
	stop := stderrors.New("stop")

	_ = entity.Walk(func(path []int, part *gomessage.Entity, werr error) error {
		if part == nil {
			return nil
		}
		mediaType, params, _ := part.Header.ContentType()
		if strings.HasPrefix(mediaType, "multipart/") {
			return nil
		}
		if !match(&part.Header) {
			return nil
		}

		body, _ := io.ReadAll(io.LimitReader(part.Body, maxInlineImageSize))
		_, dispParams, _ := part.Header.ContentDisposition()
		name := dispParams["filename"]
		if name == "" {
			name = params["name"]
		}

		outData, outType, outName, found = body, mediaType, name, true
		return stop
	})

	return outData, outType, outName, found
}

func partEncoding(bs imap.BodyStructure, partNums []int) string {
	if bs == nil {
		return ""
	}
	var encoding string
	bs.Walk(func(path []int, part imap.BodyStructure) bool {
		sp, ok := part.(*imap.BodyStructureSinglePart)
		if ok && intSliceEqual(path, partNums) {
			encoding = sp.Encoding
			return false
		}
		return true
	})
	return encoding
}

func intSliceEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func parsePartID(partID string) ([]int, error) {
	parts := strings.Split(partID, ".")
	nums := make([]int, len(parts))
	for i, p := range parts {
		n, err := fmt.Sscanf(p, "%d", &nums[i])
		if err != nil || n != 1 {
			return nil, fmt.Errorf("invalid part id: %s", partID)
		}
	}
	return nums, nil
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

// storeFlagsBulk applies a flag change to many messages in one mailbox with a
// single SELECT + STORE, instead of one round-trip (and connection) per message.
func storeFlagsBulk(client *imapclient.Client, mailbox string, uids []imap.UID, op imap.StoreFlagsOp, flags []imap.Flag) error {
	if len(uids) == 0 {
		return nil
	}
	if _, err := client.Select(mailbox, nil).Wait(); err != nil {
		return fmt.Errorf("SELECT %q failed: %w", mailbox, err)
	}
	storeCmd := client.Store(imap.UIDSetNum(uids...), &imap.StoreFlags{
		Op:     op,
		Silent: true,
		Flags:  flags,
	}, nil)
	return storeCmd.Close()
}

func moveMessage(client *imapclient.Client, srcMailbox string, uid imap.UID, destMailbox string) error {
	_, err := client.Select(srcMailbox, nil).Wait()
	if err != nil {
		return fmt.Errorf("SELECT %q failed: %w", srcMailbox, err)
	}

	uidSet := imap.UIDSetNum(uid)
	moveCmd := client.Move(uidSet, destMailbox)
	_, err = moveCmd.Wait()
	if err != nil {
		copyCmd := client.Copy(uidSet, destMailbox)
		if _, copyErr := copyCmd.Wait(); copyErr != nil {
			return fmt.Errorf("COPY to %q failed: %w", destMailbox, copyErr)
		}
		storeCmd := client.Store(uidSet, &imap.StoreFlags{
			Op:     imap.StoreFlagsAdd,
			Silent: true,
			Flags:  []imap.Flag{imap.FlagDeleted},
		}, nil)
		if storeErr := storeCmd.Close(); storeErr != nil {
			return fmt.Errorf("flag deleted failed: %w", storeErr)
		}
		if expungeErr := client.Expunge().Close(); expungeErr != nil {
			return fmt.Errorf("expunge failed: %w", expungeErr)
		}
	}
	return nil
}
