package users

import (
	"context"
	stderrors "errors"
	"net/http"
	"testing"

	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/porte/session"
	"github.com/FacileStudio/tronc/errors"
)

// recordingAuth stands in for the auth service and records which of porte's
// two password writes the controller picked. The choice is the whole logic
// under test: porte refuses SetPassword once an account has a password, so
// sending a replacement down that path answers 409 to every caller.
type recordingAuth struct {
	setCalls    int
	changeCalls int
	current     string
	next        string
	setErr      error
}

func (auth *recordingAuth) Issue(context.Context, int64, string) (string, porte.Session, error) {
	return "", porte.Session{}, nil
}

func (auth *recordingAuth) Sessions() *session.Manager { return nil }

func (auth *recordingAuth) SetPassword(_ context.Context, _ int64, password string) error {
	auth.setCalls++
	auth.next = password
	return auth.setErr
}

func (auth *recordingAuth) ChangePassword(_ context.Context, _ http.ResponseWriter, _ *http.Request, _ int64, current, next string) (string, int64, error) {
	auth.changeCalls++
	auth.current, auth.next = current, next
	return "rotated", 2, nil
}

func changePasswordFor(t *testing.T, auth *recordingAuth, request *UpdateRequest, password string) error {
	t.Helper()
	service := NewService(nil, t.TempDir(), auth)
	return service.controller.changePassword(context.Background(), nil, nil, "7", request, password)
}

// A body carrying only a password is somebody giving an account its first one,
// which is the only write a bare session may make.
func TestAFirstPasswordGoesToSetPassword(t *testing.T) {
	auth := &recordingAuth{}
	if err := changePasswordFor(t, auth, &UpdateRequest{}, "correct horse battery"); err != nil {
		t.Fatalf("set a first password: %v", err)
	}
	if auth.setCalls != 1 || auth.changeCalls != 0 {
		t.Fatalf("expected one SetPassword and no ChangePassword, got %d and %d", auth.setCalls, auth.changeCalls)
	}
	if auth.next != "correct horse battery" {
		t.Fatalf("the password was not passed through: %q", auth.next)
	}
}

// The current password is what turns the request into a replacement, and it is
// confirmed by porte rather than by this app.
func TestACurrentPasswordGoesToChangePassword(t *testing.T) {
	auth := &recordingAuth{}
	current := "  old passphrase  "
	request := &UpdateRequest{CurrentPassword: &current}
	if err := changePasswordFor(t, auth, request, "new passphrase"); err != nil {
		t.Fatalf("change the password: %v", err)
	}
	if auth.changeCalls != 1 || auth.setCalls != 0 {
		t.Fatalf("expected one ChangePassword and no SetPassword, got %d and %d", auth.changeCalls, auth.setCalls)
	}
	if auth.current != "old passphrase" || auth.next != "new passphrase" {
		t.Fatalf("the passwords were not passed through: current=%q next=%q", auth.current, auth.next)
	}
}

// porte's answer to a replacement sent down the first-password path is
// ErrPasswordSet, and it maps to a 400 naming the missing field. Letting it
// through as porte's 409 tells the caller they lost a race rather than that
// they left something out.
func TestAReplacementWithoutTheCurrentPasswordIsRefusedByName(t *testing.T) {
	auth := &recordingAuth{setErr: porte.ErrPasswordSet}
	blank := "   "
	request := &UpdateRequest{CurrentPassword: &blank}

	err := changePasswordFor(t, auth, request, "new passphrase")
	if auth.changeCalls != 0 {
		t.Fatalf("a blank current password reached ChangePassword: %d calls", auth.changeCalls)
	}

	var envelope *errors.Error
	if !stderrors.As(err, &envelope) {
		t.Fatalf("expected a tronc error envelope, got %v", err)
	}
	if envelope.Code != "invalid_argument" {
		t.Fatalf("expected invalid_argument, got %q", envelope.Code)
	}
}
