package mcp

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Pending action kinds.
const (
	kindSend        = "send"
	kindDeleteDraft = "delete_draft"
)

// Outcomes a pending action is resolved with.
const (
	outcomeRunning   = "running"
	outcomeDone      = "done"
	outcomeDeclined  = "declined"
	outcomeDismissed = "dismissed"
)

// errGone is the refusal for an action that ran, was dismissed, expired or
// was never this person's.
var errGone = errors.NotFound("this action was already confirmed, dismissed or expired")

// pendingStore keeps the sensitive actions waiting for a person, in
// mcp_pending_actions.
type pendingStore struct {
	orm *gorm.DB
	ttl time.Duration
}

// create stores a new pending action under a fresh random id.
func (p pendingStore) create(ctx context.Context, action schemas.PendingAction) (schemas.PendingAction, error) {
	token := make([]byte, 18)
	if _, err := rand.Read(token); err != nil {
		return action, errors.Internal("failed to create a pending action", err)
	}
	action.ID = base64.RawURLEncoding.EncodeToString(token)
	action.ExpiresAt = time.Now().Add(p.ttl)
	if err := p.orm.WithContext(ctx).Create(&action).Error; err != nil {
		return action, errors.Internal("failed to create a pending action", err)
	}
	return action, nil
}

// open lists userID's unresolved, unexpired actions, newest first.
func (p pendingStore) open(ctx context.Context, userID int64) ([]schemas.PendingAction, error) {
	var actions []schemas.PendingAction
	err := p.orm.WithContext(ctx).Where("user_id = ? AND resolved_at IS NULL AND expires_at > now()", userID).
		Order("created_at DESC, id").Find(&actions).Error
	if err != nil {
		return nil, errors.Internal("failed to list pending actions", err)
	}
	return actions, nil
}

// claim marks one of userID's open actions as running and returns it. Only
// one claim of an action can succeed, which makes every confirmation, from
// the client or from the confirm page, single use.
func (p pendingStore) claim(ctx context.Context, userID int64, id string) (schemas.PendingAction, error) {
	var claimed []schemas.PendingAction
	err := p.orm.WithContext(ctx).Model(&claimed).Clauses(clause.Returning{}).
		Where("id = ? AND user_id = ? AND resolved_at IS NULL AND expires_at > now()", id, userID).
		Updates(map[string]any{"resolved_at": gorm.Expr("now()"), "outcome": outcomeRunning}).Error
	if err != nil {
		return schemas.PendingAction{}, errors.Internal("failed to claim a pending action", err)
	}
	if len(claimed) == 0 {
		return schemas.PendingAction{}, errGone
	}
	return claimed[0], nil
}

// finish records the outcome of a claimed action.
func (p pendingStore) finish(ctx context.Context, id, outcome string) {
	p.orm.WithContext(ctx).Model(&schemas.PendingAction{}).Where("id = ?", id).Update("outcome", outcome)
}

// release reopens a claimed action whose run failed, so the person can try
// again.
func (p pendingStore) release(ctx context.Context, id string) {
	p.orm.WithContext(ctx).Model(&schemas.PendingAction{}).Where("id = ? AND outcome = ?", id, outcomeRunning).
		Updates(map[string]any{"resolved_at": nil, "outcome": ""})
}

// resolve closes one of userID's open actions without running it.
func (p pendingStore) resolve(ctx context.Context, userID int64, id, outcome string) error {
	result := p.orm.WithContext(ctx).Model(&schemas.PendingAction{}).
		Where("id = ? AND user_id = ? AND resolved_at IS NULL", id, userID).
		Updates(map[string]any{"resolved_at": gorm.Expr("now()"), "outcome": outcome})
	if result.Error != nil {
		return errors.Internal("failed to resolve a pending action", result.Error)
	}
	if result.RowsAffected == 0 {
		return errGone
	}
	return nil
}
