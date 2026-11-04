package mail

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
)

const (
	grayQuestion   = "inbox_tab"
	grayFloor      = 0.8
	grayBatch      = 50
	grayBodyBytes  = 1500
	grayJudgedNone = "judged"
	grayJudgedBy   = "jev"
	grayPassTime   = 2 * time.Minute
)

var grayChoices = map[string]string{
	"primary":      categoryPrimary,
	"newsletter":   categoryNewsletter,
	"notification": categoryNotification,
}

func grayAsk() jevQuestion {
	return jevQuestion{
		Type: "choice",
		Instructions: "The state is one email that reached a person's inbox. Its headers did not mark it as bulk. " +
			"Decide which inbox tab it belongs in.",
		Criteria: map[string]string{
			"primary":      "a person wrote it to the recipient, or it needs the recipient's attention or an answer: a client, colleague, teacher, partner, a personal message, an invoice or a contract",
			"newsletter":   "marketing, a newsletter, a promotion, a digest, product news, or any message sent to many people at once",
			"notification": "an automated message from a service: an alert, a receipt, a sign-in code, a CI or monitoring report, a shipping update",
		},
		Options: []string{"primary", "newsletter", "notification"},
	}
}

type grayState struct {
	From    string `json:"from"`
	Name    string `json:"name"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// judgeGrayZone asks jev where the inbox messages the headers could not sort
// belong: no list headers, no sender rule, and a sender the user has never
// written to (contacts holds everyone the user sent mail to). Only a confident answer moves a message; every judged message
// is marked so it is asked about once.
func (s *Service) judgeGrayZone(folderID int64) {
	if s.judge == nil {
		return
	}
	if _, busy := s.judging.LoadOrStore(folderID, true); busy {
		return
	}
	defer s.judging.Delete(folderID)
	ctx, cancel := context.WithTimeout(context.Background(), grayPassTime)
	defer cancel()
	for _, email := range s.grayCandidates(ctx, folderID) {
		category, reason := s.grayVerdict(ctx, email)
		s.orm.WithContext(ctx).Model(&schemas.Email{}).Where("id = ?", email.ID).
			Updates(map[string]any{"category": category, "bulk_reason": reason})
	}
}

func (s *Service) grayVerdict(ctx context.Context, email schemas.Email) (string, string) {
	answer, err := s.judge.choose(ctx, grayState{
		From: email.FromAddress, Name: email.FromName, Subject: email.Subject, Body: clampText(email.BodyText, grayBodyBytes),
	}, grayQuestion, grayAsk())
	if err != nil {
		slog.Warn("jev inbox judge", "email", email.ID, "error", err)
		return categoryPrimary, grayJudgedNone
	}
	category, ok := grayChoices[answer.Choice]
	if !ok || answer.Confidence < grayFloor || category == categoryPrimary {
		return categoryPrimary, grayJudgedNone
	}
	return category, grayJudgedBy
}

func (s *Service) grayCandidates(ctx context.Context, folderID int64) []schemas.Email {
	var emails []schemas.Email
	s.orm.WithContext(ctx).Raw(`
		SELECT e.* FROM emails e
		JOIN folders f ON f.id = e.folder_id AND f.type = ?
		WHERE e.folder_id = ? AND e.category = '' AND e.bulk_reason = ?
		  AND NOT EXISTS (SELECT 1 FROM sender_rules r WHERE r.account_id = e.account_id AND r.address = LOWER(e.from_address))
		  AND NOT EXISTS (SELECT 1 FROM contacts c WHERE c.account_id = e.account_id AND LOWER(c.email) = LOWER(e.from_address))
		ORDER BY e.date DESC
		LIMIT ?`, schemas.FolderTypeInbox, folderID, bulkReasonNone, grayBatch).Scan(&emails)
	return emails
}

func clampText(text string, limit int) string {
	text = strings.TrimSpace(text)
	if len(text) <= limit {
		return text
	}
	return text[:limit]
}
