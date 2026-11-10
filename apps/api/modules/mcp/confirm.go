package mcp

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"
	mcpkit "github.com/FacileStudio/tronc/mcp"
)

// Pending action kinds.
const (
	kindSend        = "send"
	kindDeleteDraft = "delete_draft"
)

// errTarget is the refusal for a pending action whose target is not a draft
// reference.
var errTarget = errors.NotFound("this action names no draft")

// kinds are the two sensitive actions: both fingerprint the draft, so an
// edit after the question closes the action instead of sending it.
func (t *toolset) kinds() map[string]mcpkit.Kind {
	return map[string]mcpkit.Kind{
		kindSend:        {Tool: toolSend, Fingerprint: t.fingerprint, Run: t.sendDraft},
		kindDeleteDraft: {Tool: toolDeleteDraft, Fingerprint: t.fingerprint, Run: t.removeDraft},
	}
}

// target is the pending action target of a draft.
func target(ref draftRef) string {
	return fmt.Sprintf("%d/%d", ref.AccountID, ref.DraftID)
}

// refOf reads the draft an action targets, and the person it acts for.
func refOf(a mcpkit.Action) (int64, draftRef, error) {
	userID, errUser := strconv.ParseInt(a.UserID, 10, 64)
	account, draft, _ := strings.Cut(a.Target, "/")
	accountID, errAccount := strconv.ParseInt(account, 10, 64)
	draftID, errDraft := strconv.ParseInt(draft, 10, 64)
	if errUser != nil || errAccount != nil || errDraft != nil {
		return 0, draftRef{}, errTarget
	}
	return userID, draftRef{AccountID: accountID, DraftID: draftID}, nil
}

// loadDraft loads the draft an action targets.
func (t *toolset) loadDraft(ctx context.Context, a mcpkit.Action) (schemas.Email, error) {
	userID, ref, err := refOf(a)
	if err != nil {
		return schemas.Email{}, err
	}
	return t.mail.Draft(ctx, userID, ref.AccountID, ref.DraftID)
}

// fingerprint is the digest of the targeted draft as it is now.
func (t *toolset) fingerprint(ctx context.Context, a mcpkit.Action) (string, error) {
	draft, err := t.loadDraft(ctx, a)
	if err != nil {
		return "", err
	}
	return digest(draft), nil
}

// sendDraft sends the targeted draft.
func (t *toolset) sendDraft(ctx context.Context, a mcpkit.Action) error {
	userID, ref, err := refOf(a)
	if err != nil {
		return err
	}
	return t.mail.SendDraft(ctx, userID, ref.AccountID, ref.DraftID)
}

// removeDraft deletes the targeted draft.
func (t *toolset) removeDraft(ctx context.Context, a mcpkit.Action) error {
	userID, ref, err := refOf(a)
	if err != nil {
		return err
	}
	return t.mail.DeleteDraft(ctx, userID, ref.AccountID, ref.DraftID)
}

// digest fingerprints what a draft would send.
func digest(draft schemas.Email) string {
	fields := []string{
		draft.ToAddresses, draft.CcAddresses, draft.BccAddresses, draft.Subject,
		draft.BodyText, draft.BodyHTML, draft.InReplyTo, fmt.Sprint(draft.HasAttachments),
	}
	sum := sha256.Sum256([]byte(strings.Join(fields, "\x00")))
	return hex.EncodeToString(sum[:16])
}
