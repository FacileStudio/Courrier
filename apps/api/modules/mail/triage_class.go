package mail

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"
)

const (
	classKnown        = "known"
	classUnknownHuman = "unknown_human"
	classBulk         = "bulk"
	bulkReasonNone    = "none"
)

var triageClasses = []string{classKnown, classUnknownHuman, classBulk}

var sinceShort = regexp.MustCompile(`^(\d+)([mhdw])$`)

var sinceUnits = map[string]time.Duration{
	"m": time.Minute,
	"h": time.Hour,
	"d": 24 * time.Hour,
	"w": 7 * 24 * time.Hour,
}

type threadKey struct {
	AccountID int64
	ThreadID  string
}

func parseSince(raw string, now time.Time) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, nil
	}
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		return t, nil
	}
	m := sinceShort.FindStringSubmatch(raw)
	if m == nil {
		return time.Time{}, errors.Invalid("since must be RFC3339 or a duration like 30m, 2h, 1d, 1w")
	}
	n, _ := strconv.Atoi(m[1])
	return now.Add(-time.Duration(n) * sinceUnits[m[2]]), nil
}

func parseClasses(raw string) (map[string]bool, error) {
	wanted := map[string]bool{}
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if part != classKnown && part != classUnknownHuman && part != classBulk {
			return nil, errors.Invalid("class must be a comma list of known, unknown_human, bulk")
		}
		wanted[part] = true
	}
	return wanted, nil
}

func visibleBulkReason(reason string) string {
	if reason == bulkReasonNone {
		return ""
	}
	return reason
}

func classifyForStore(s Signals, fromAddress string) (string, string) {
	category, reason := Classify(s, fromAddress)
	if category == "" {
		return "", bulkReasonNone
	}
	return category, reason
}

func bulkReasonOf(e schemas.Email) string {
	if e.Category == "" && e.BulkReason == "" {
		_, reason := Classify(Signals{}, e.FromAddress)
		return reason
	}
	if e.BulkReason == bulkReasonNone {
		return ""
	}
	if e.BulkReason == "" {
		return "noreply-sender"
	}
	return e.BulkReason
}

func resolveClass(e schemas.Email, sent map[threadKey]bool, known map[string]bool) (string, string) {
	if reason := bulkReasonOf(e); reason != "" {
		return classBulk, reason
	}
	if sent[threadKey{e.AccountID, e.ThreadID}] {
		return classKnown, "in-sent-thread"
	}
	if known[strings.ToLower(strings.TrimSpace(e.FromAddress))] {
		return classKnown, "sent-to"
	}
	return classUnknownHuman, "no-signal"
}

func (s *Service) sentThreads(ctx context.Context, accounts []int64, convs []Conversation) (map[threadKey]bool, error) {
	threads := []string{}
	for _, c := range convs {
		if c.Email.ThreadID != "" {
			threads = append(threads, c.Email.ThreadID)
		}
	}
	sent := map[threadKey]bool{}
	if len(threads) == 0 {
		return sent, nil
	}
	var keys []threadKey
	err := s.orm.WithContext(ctx).Table("emails e").
		Select("DISTINCT e.account_id, e.thread_id").
		Joins("JOIN folders f ON f.id = e.folder_id").
		Where("f.type = ? AND e.account_id IN ? AND e.thread_id IN ?", schemas.FolderTypeSent, accounts, threads).
		Scan(&keys).Error
	if err != nil {
		return nil, errors.Internal("failed to look up sent threads", err)
	}
	for _, k := range keys {
		sent[k] = true
	}
	return sent, nil
}
