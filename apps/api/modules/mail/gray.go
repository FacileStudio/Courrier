package mail

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/tronc/bulle"
)

const (
	grayFloor      = 0.8
	grayBatch      = 25
	grayWindow     = 200
	grayBodyBytes  = 1500
	grayInputBytes = 60 << 10
	grayWireBytes  = 88 << 10
	grayTokens     = 2048
	grayJudgedNone = "judged"
	grayJudgedBy   = "bulle"
	grayPassTime   = 2 * time.Minute
)

var grayChoices = map[string]string{
	"primary":      categoryPrimary,
	"newsletter":   categoryNewsletter,
	"notification": categoryNotification,
}

const graySystem = "You sort email into inbox tabs for a mail client."

const grayPrompt = "The input is a JSON list of emails that reached a person's inbox. Their headers did not mark them as bulk, " +
	"and the person never wrote to their senders. For each email, decide which inbox tab it belongs in and how sure you are, " +
	"from 0 to 1.\n" +
	"- primary: a person wrote it to the recipient, or it needs the recipient's attention or an answer: a client, colleague, " +
	"teacher, partner, a personal message, an invoice or a contract.\n" +
	"- newsletter: marketing, a newsletter, a promotion, a digest, product news, or any message sent to many people at once.\n" +
	"- notification: an automated message from a service: an alert, a receipt, a sign-in code, a CI or monitoring report, " +
	"a shipping update.\n" +
	"Answer one verdict per email, with the email's n."

var graySchema = json.RawMessage(`{"type":"object","additionalProperties":false,"required":["verdicts"],` +
	`"properties":{"verdicts":{"type":"array","items":{"type":"object","additionalProperties":false,` +
	`"required":["n","tab","confidence"],"properties":{"n":{"type":"integer"},` +
	`"tab":{"type":"string","enum":["primary","newsletter","notification"]},` +
	`"confidence":{"type":"number","minimum":0,"maximum":1}}}}}}`)

// Completer is Bulle's completion call, as tronc/bulle's Client makes it for
// the user on the context.
type Completer interface {
	Complete(ctx context.Context, req bulle.CompleteRequest) (bulle.Completion, error)
}

// Triage sorts the inbox mail the headers could not classify through Bulle,
// for owners whose plan, or whose mailbox's space plan, includes ai.access.
// AIAccess reads those plans; link is the space's Registre id, empty for a
// personal mailbox.
type Triage struct {
	Bulle    Completer
	AIAccess func(ctx context.Context, userID int64, link string) (bool, error)
}

// WithTriage turns on the AI triage after inbox syncs and returns the
// service. Without it only the header and sender rules sort mail.
func (s *Service) WithTriage(triage Triage) *Service {
	s.triage = &triage
	s.gate = &triageGate{}
	return s
}

type grayEmail struct {
	N       int    `json:"n"`
	From    string `json:"from"`
	Name    string `json:"name"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type grayVerdict struct {
	N          int     `json:"n"`
	Tab        string  `json:"tab"`
	Confidence float64 `json:"confidence"`
}

// judgeGrayZone asks Bulle, as the mailbox owner, where the inbox messages
// the headers could not sort belong: no list headers, no sender rule, and a
// sender the user has never written to (contacts holds everyone the user sent
// mail to). It runs only when the owner's plan includes ai.access, in one
// call per pass over mail from one sender domain, so a mail can only steer
// the verdicts of mail from its own domain. The owner's triage gate bounds
// the calls: a daily budget, and a backoff after a failure. Only a
// confident answer moves a message; every message with a verdict is marked
// so it is asked about once, and a failed or skipped call leaves them for
// the non-AI rules and a later pass.
func (s *Service) judgeGrayZone(folderID int64) {
	if s.triage == nil {
		return
	}
	if _, busy := s.judging.LoadOrStore(folderID, true); busy {
		return
	}
	defer s.judging.Delete(folderID)
	ctx, cancel := context.WithTimeout(context.Background(), grayPassTime)
	defer cancel()
	owner, link, ok := s.triageOwner(ctx, folderID)
	if !ok {
		return
	}
	emails := s.grayCandidates(ctx, folderID)
	if len(emails) == 0 || !s.gate.take(owner) {
		return
	}
	verdicts, err := s.askBulle(porte.WithIdentity(ctx, porte.Identity{UserID: owner}), emails, link)
	s.gate.settle(owner, err)
	if err != nil {
		slog.Warn("bulle inbox triage", "folder", folderID, "error", err)
		return
	}
	s.fileVerdicts(ctx, emails, verdicts)
}

// fileVerdicts stores each email's verdict, by its position in emails.
func (s *Service) fileVerdicts(ctx context.Context, emails []schemas.Email, verdicts map[int]grayVerdict) {
	for i, email := range emails {
		verdict, found := verdicts[i]
		if !found {
			continue
		}
		category, reason := grayOutcome(verdict)
		s.orm.WithContext(ctx).Model(&schemas.Email{}).Where("id = ?", email.ID).
			Updates(map[string]any{"category": category, "bulk_reason": reason})
	}
}

// triageOwner is the user the folder's mail account belongs to and, for a
// space mailbox, the space's Registre id, when the owner is still a member
// of the space and their plan or the space's includes ai.access.
func (s *Service) triageOwner(ctx context.Context, folderID int64) (int64, string, bool) {
	var owners []struct {
		UserID  int64
		SpaceID *string
		Link    string
	}
	if err := s.orm.WithContext(ctx).Raw(`SELECT a.user_id, a.space_id, coalesce(sp.registre_space_id, '') AS link
		FROM folders f JOIN accounts a ON a.id = f.account_id LEFT JOIN spaces sp ON sp.id::text = a.space_id
		WHERE f.id = ?`, folderID).Scan(&owners).Error; err != nil || len(owners) == 0 {
		return 0, "", false
	}
	owner := owners[0]
	if owner.SpaceID != nil && *owner.SpaceID != "" && s.requireSpaceMember(ctx, *owner.SpaceID, owner.UserID) != nil {
		return 0, "", false
	}
	allowed, err := s.triage.AIAccess(ctx, owner.UserID, owner.Link)
	if err != nil {
		slog.Warn("bulle inbox triage: plan check", "folder", folderID, "error", err)
		return 0, "", false
	}
	return owner.UserID, owner.Link, allowed
}

// askBulle classifies emails in one completion, for the Registre space link
// when the mailbox has one, and returns the verdicts by position in emails.
func (s *Service) askBulle(ctx context.Context, emails []schemas.Email, link string) (map[int]grayVerdict, error) {
	input, err := grayInput(emails)
	if err != nil {
		return nil, err
	}
	out, err := s.triage.Bulle.Complete(ctx, bulle.CompleteRequest{
		System: graySystem, Prompt: grayPrompt, Input: input, Schema: graySchema, MaxTokens: grayTokens, SpaceID: link,
	})
	if err != nil {
		return nil, err
	}
	var answer struct {
		Verdicts []grayVerdict `json:"verdicts"`
	}
	if err := json.Unmarshal(out.JSON, &answer); err != nil {
		return nil, err
	}
	verdicts := make(map[int]grayVerdict, len(answer.Verdicts))
	for _, verdict := range answer.Verdicts {
		if verdict.N >= 0 && verdict.N < len(emails) {
			verdicts[verdict.N] = verdict
		}
	}
	return verdicts, nil
}

// grayOutcome is the category and reason a verdict files a message under:
// only a confident newsletter or notification verdict moves it.
func grayOutcome(verdict grayVerdict) (string, string) {
	category, ok := grayChoices[verdict.Tab]
	if !ok || verdict.Confidence < grayFloor || category == categoryPrimary {
		return categoryPrimary, grayJudgedNone
	}
	return category, grayJudgedBy
}

// grayCandidates is the newest unjudged inbox mail from strangers, all from
// the sender domain of the newest one.
func (s *Service) grayCandidates(ctx context.Context, folderID int64) []schemas.Email {
	var emails []schemas.Email
	s.orm.WithContext(ctx).Raw(`
		SELECT e.* FROM emails e
		JOIN folders f ON f.id = e.folder_id AND f.type = ?
		WHERE e.folder_id = ? AND e.category = '' AND e.bulk_reason = ?
		  AND NOT EXISTS (SELECT 1 FROM sender_rules r WHERE r.account_id = e.account_id AND r.address = LOWER(e.from_address))
		  AND NOT EXISTS (SELECT 1 FROM contacts c WHERE c.account_id = e.account_id AND LOWER(c.email) = LOWER(e.from_address))
		ORDER BY e.date DESC
		LIMIT ?`, schemas.FolderTypeInbox, folderID, bulkReasonNone, grayWindow).Scan(&emails)
	return sameSenderDomain(emails, grayBatch)
}
