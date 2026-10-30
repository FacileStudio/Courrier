package mail

import (
	"context"
	"math"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"
)

// Conversation is one collapsed thread within a folder: the latest message
// stands in for the group, plus thread-level aggregates the list view needs.
type Conversation struct {
	Email        schemas.Email
	MessageCount int
	UnreadCount  int
	EmailIDs     []int64
}

type conversationRow struct {
	schemas.Email
	ConvCount   int    `gorm:"column:conv_count"`
	ConvUnread  int    `gorm:"column:conv_unread"`
	ConvStarred bool   `gorm:"column:conv_starred"`
	ConvAttach  bool   `gorm:"column:conv_attach"`
	ConvIDs     string `gorm:"column:conv_ids"`
}

type conversationFilter struct {
	FolderIDs  []int64
	UnreadOnly bool
	Since      time.Time
}

const conversationPartition = "COALESCE(NULLIF(thread_id, ''), 'msg:' || id::text)"

const countConversationsSQL = `SELECT COUNT(*) FROM (SELECT 1 FROM emails WHERE folder_id IN ? GROUP BY ` +
	conversationPartition + ` HAVING (? = FALSE OR COUNT(*) FILTER (WHERE NOT is_read) > 0) AND MAX(date) >= ?) c`

const listConversationsSQL = `
	WITH scoped AS (
		SELECT id, account_id, folder_id, message_id, thread_id, subject,
		       from_address, from_name, to_addresses, cc_addresses, date,
		       is_read, is_starred, has_attachments, in_reply_to, "references",
		       category, bulk_reason, ` + conversationPartition + ` AS tkey
		FROM emails
		WHERE folder_id IN ?
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
		       in_reply_to, "references", category, bulk_reason, tkey
		FROM scoped
		ORDER BY tkey, date DESC, id DESC
	)
	SELECT rep.id, rep.account_id, rep.folder_id, rep.message_id, rep.thread_id, rep.subject,
	       rep.from_address, rep.from_name, rep.to_addresses, rep.cc_addresses, rep.date,
	       rep.in_reply_to, rep."references", rep.category, rep.bulk_reason,
	       agg.conv_count, agg.conv_unread, agg.conv_starred, agg.conv_attach, agg.conv_ids
	FROM agg JOIN rep USING (tkey)
	WHERE (? = FALSE OR agg.conv_unread > 0) AND agg.conv_last >= ?
	ORDER BY agg.conv_last DESC, rep.id DESC
	LIMIT ? OFFSET ?`

func (s *Service) countConversations(ctx context.Context, f conversationFilter) (int64, error) {
	var total int64
	err := s.orm.WithContext(ctx).Raw(countConversationsSQL, f.FolderIDs, f.UnreadOnly, f.Since).Scan(&total).Error
	if err != nil {
		return 0, errors.Internal("failed to count conversations", err)
	}
	return total, nil
}

func (s *Service) queryConversations(ctx context.Context, f conversationFilter, limit, offset int) ([]Conversation, error) {
	if limit <= 0 {
		limit = math.MaxInt32
	}
	var rows []conversationRow
	err := s.orm.WithContext(ctx).Raw(listConversationsSQL, f.FolderIDs, f.UnreadOnly, f.Since, limit, offset).Scan(&rows).Error
	if err != nil {
		return nil, errors.Internal("failed to list conversations", err)
	}
	conversations := make([]Conversation, len(rows))
	for i, r := range rows {
		e := r.Email
		e.IsRead = r.ConvUnread == 0
		e.IsStarred = r.ConvStarred
		e.HasAttachments = r.ConvAttach
		conversations[i] = Conversation{Email: e, MessageCount: r.ConvCount, UnreadCount: r.ConvUnread, EmailIDs: parseCSVInts(r.ConvIDs)}
	}
	return conversations, nil
}
