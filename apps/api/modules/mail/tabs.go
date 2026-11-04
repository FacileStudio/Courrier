package mail

import (
	"context"
	"net/http"
	"strings"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"
	"gorm.io/gorm/clause"
)

const (
	categoryPrimary       = ""
	categoryNewsletter    = "bulk"
	categoryNotification  = "auto"
	maxSenderAddressBytes = 320
)

var tabs = map[string]string{
	"primary":       categoryPrimary,
	"newsletters":   categoryNewsletter,
	"notifications": categoryNotification,
}

// FolderQuery is one page of a folder's conversation list, optionally narrowed
// to an inbox tab.
type FolderQuery struct {
	FolderType string
	Page       int
	Limit      int
	UnreadOnly bool
	Category   *string
}

// tabCategory maps an inbox tab to the email category it lists. An empty tab
// lists every conversation.
func tabCategory(tab string) (*string, error) {
	if tab == "" || tab == "all" {
		return nil, nil
	}
	category, ok := tabs[tab]
	if !ok {
		return nil, errors.Invalid("tab must be primary, newsletters, notifications or all")
	}
	return &category, nil
}

func searchWords(query string) []string {
	words := strings.Fields(query)
	if len(words) == 0 {
		return []string{query}
	}
	return words
}

// SenderRuleRequest pins a sender to an inbox tab, or with tab "auto" drops
// the pin and lets the classifier decide again.
type SenderRuleRequest struct {
	Address string `json:"address"`
	Tab     string `json:"tab"`
}

// SetSenderRule records which tab every message from one address lands in.
func (s *Service) SetSenderRule(ctx context.Context, userID, accountID int64, req SenderRuleRequest) error {
	if _, err := s.getAccount(ctx, userID, accountID); err != nil {
		return err
	}
	address := strings.ToLower(strings.TrimSpace(req.Address))
	if address == "" || len(address) > maxSenderAddressBytes || !strings.Contains(address, "@") {
		return errors.Invalid("address must be an email address")
	}
	db := s.orm.WithContext(ctx)
	if req.Tab == "auto" {
		return wrapInternal(db.Where("account_id = ? AND address = ?", accountID, address).Delete(&schemas.SenderRule{}).Error)
	}
	category, ok := tabs[req.Tab]
	if !ok {
		return errors.Invalid("tab must be primary, newsletters, notifications or auto")
	}
	rule := schemas.SenderRule{AccountID: accountID, Address: address, Category: category}
	return wrapInternal(db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "account_id"}, {Name: "address"}},
		DoUpdates: clause.AssignmentColumns([]string{"category", "updated_at"}),
	}).Create(&rule).Error)
}

func wrapInternal(err error) error {
	if err != nil {
		return errors.Internal("failed to save the sender rule", err)
	}
	return nil
}

func (s *Service) handleSenderRule(w http.ResponseWriter, req *http.Request) {
	userID, accountID, err := s.accountScope(req)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	var body SenderRuleRequest
	if err := httpjson.DecodeJSON(w, req, &body); err != nil {
		httpjson.WriteError(w, err)
		return
	}
	if err := s.SetSenderRule(req.Context(), userID, accountID, body); err != nil {
		httpjson.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
