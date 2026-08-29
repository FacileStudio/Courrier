package users

import (
	"context"
	stderrors "errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/tronc/errors"
)

// changePassword picks between porte's two password writes.
//
// They are two calls and not one because only one of them is safe to make
// with nothing but a session: SetPassword gives a first password to an account
// that has none, and porte refuses it with ErrPasswordSet once there is one.
// Replacing a password goes through ChangePassword, which confirms the current
// one — OWASP ASVS puts that at L1 (v4 2.1.6, v5 6.2.3) — then ends the
// account's other logins and rotates this caller's session through w.
//
// A blank current password counts as none given, so the answer is the one that
// says what is missing rather than "invalid credentials".
func (controller *Controller) changePassword(context context.Context, w http.ResponseWriter, request *http.Request, userID string, req *UpdateRequest, password string) error {
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return errors.Internal("failed to parse user id", err)
	}

	current := ""
	if req.CurrentPassword != nil {
		current = strings.TrimSpace(*req.CurrentPassword)
	}
	if current != "" {
		_, _, err := controller.service.tokens.ChangePassword(context, w, request, id, current, password)
		return passwordError(err)
	}

	err = controller.service.tokens.SetPassword(context, id, password)
	if stderrors.Is(err, porte.ErrPasswordSet) {
		return errors.Invalid("current_password is required to change an existing password")
	}
	return passwordError(err)
}

// passwordError turns porte's sentinels into the suite's error envelope.
//
// tronc answers 500 to anything it cannot recognise as an *errors.Error, so a
// wrong current password would otherwise read as an outage rather than as a
// refusal the caller can act on.
func passwordError(err error) error {
	switch {
	case err == nil:
		return nil
	case stderrors.Is(err, porte.ErrWrongPassword):
		return errors.Unauthorized("current password is incorrect")
	case stderrors.Is(err, porte.ErrNoPassword):
		return errors.Invalid("this account has no password to change")
	case stderrors.Is(err, porte.ErrWeakPassword):
		return errors.Invalid("password is too short")
	default:
		return err
	}
}
