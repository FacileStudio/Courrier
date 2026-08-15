package mail

import (
	"regexp"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
)

var msgIDToken = regexp.MustCompile(`<[^<>]+>`)

// normID reduces a Message-ID to a comparable form: angle brackets stripped,
// trimmed, lowercased. Both sides of every comparison go through this so the
// envelope's "<a@b>" matches a References entry written the same way.
func normID(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "<>")
	s = strings.TrimSpace(s)
	return strings.ToLower(s)
}

// parseRefIDs extracts the Message-IDs from a References or In-Reply-To header,
// oldest-first, normalized. It prefers angle-bracketed tokens (the RFC form)
// and falls back to splitting on whitespace/commas for sloppy senders.
func parseRefIDs(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	var out []string
	if matches := msgIDToken.FindAllString(s, -1); len(matches) > 0 {
		for _, m := range matches {
			if id := normID(m); id != "" {
				out = append(out, id)
			}
		}
		return out
	}
	for _, f := range strings.FieldsFunc(s, func(r rune) bool {
		return r == ' ' || r == '\t' || r == ',' || r == '\n' || r == '\r'
	}) {
		if id := normID(f); id != "" {
			out = append(out, id)
		}
	}
	return out
}

// collectThreadIDs gathers every Message-ID related to a message — its own, its
// References chain, and its In-Reply-To — deduped, preserving order.
func collectThreadIDs(messageID, inReplyTo, references string) []string {
	seen := map[string]bool{}
	var out []string
	for _, src := range []string{messageID, references, inReplyTo} {
		for _, id := range parseRefIDs(src) {
			if !seen[id] {
				seen[id] = true
				out = append(out, id)
			}
		}
	}
	return out
}

// newThreadID picks the stable identifier for a brand-new conversation: the
// root ancestor from References (oldest-first), else the direct parent, else
// the message's own id. Using the root keeps the id consistent even when
// replies arrive before their parents.
func newThreadID(messageID, inReplyTo, references string) string {
	if refs := parseRefIDs(references); len(refs) > 0 {
		return refs[0]
	}
	if id := normID(inReplyTo); id != "" {
		return id
	}
	return normID(messageID)
}

// assignThreadID resolves the thread a message belongs to and records every
// related Message-ID against it. Incremental union-find: zero known threads
// starts a new one, one joins it, several means this message bridges them and
// they are merged (smallest id wins). Runs in a transaction so a merge is
// atomic. Returns "" only when the message carries no usable ids at all.
//
// Concurrent syncs of one account run sequentially in practice (folders sync
// one at a time), so read-committed is sufficient; a bridging message
// self-heals any split. If the transaction fails, return "" rather than a
// thread id with no backing links.
func assignThreadID(db *gorm.DB, accountID int64, messageID, inReplyTo, references string) string {
	related := collectThreadIDs(messageID, inReplyTo, references)
	if len(related) == 0 {
		return ""
	}

	var threadID string
	if err := db.Transaction(func(tx *gorm.DB) error {
		var existing []string
		tx.Model(&schemas.ThreadLink{}).
			Where("account_id = ? AND message_id IN ?", accountID, related).
			Distinct().
			Pluck("thread_id", &existing)

		switch len(existing) {
		case 0:
			threadID = newThreadID(messageID, inReplyTo, references)
		case 1:
			threadID = existing[0]
		default:
			threadID = existing[0]
			for _, t := range existing {
				if t < threadID {
					threadID = t
				}
			}
			var losers []string
			for _, t := range existing {
				if t != threadID {
					losers = append(losers, t)
				}
			}
			if len(losers) > 0 {
				tx.Model(&schemas.Email{}).
					Where("account_id = ? AND thread_id IN ?", accountID, losers).
					Update("thread_id", threadID)
				tx.Model(&schemas.ThreadLink{}).
					Where("account_id = ? AND thread_id IN ?", accountID, losers).
					Update("thread_id", threadID)
			}
		}

		if threadID == "" {
			return nil
		}

		links := make([]schemas.ThreadLink, len(related))
		for i, id := range related {
			links[i] = schemas.ThreadLink{AccountID: accountID, MessageID: id, ThreadID: threadID}
		}
		return tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "account_id"}, {Name: "message_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"thread_id"}),
		}).Create(&links).Error
	}); err != nil {
		return ""
	}

	return threadID
}

// BackfillThreads assigns thread ids to any emails stored before threading
// existed. Idempotent: once thread_id is populated a row is skipped on
// subsequent runs. Processing order is irrelevant — union-find anchors each
// thread on its References root regardless of arrival order. Selects only the
// header columns and works in primary-key batches so a large legacy mailbox
// never loads message bodies into memory; meant to run in the background.
func BackfillThreads(db *gorm.DB) error {
	var batch []schemas.Email
	return db.
		Model(&schemas.Email{}).
		Select("id", "account_id", "message_id", "in_reply_to", "references").
		Where("thread_id = '' OR thread_id IS NULL").
		FindInBatches(&batch, 500, func(tx *gorm.DB, _ int) error {
			for _, e := range batch {
				tid := assignThreadID(db, e.AccountID, e.MessageID, e.InReplyTo, e.References)
				if tid != "" {
					db.Model(&schemas.Email{}).Where("id = ?", e.ID).Update("thread_id", tid)
				}
			}
			return nil
		}).Error
}
