package users

import (
	"net/http"
	"strings"
	"testing"

	"github.com/FacileStudio/porte"
)

const mePath = "/users/me"

// A body carrying only a password is somebody giving an account its first one,
// which is the only password write a bare session may make.
func TestAFirstPasswordIsAcceptedFromABareSession(t *testing.T) {
	h := newHarness(t)
	_, token := h.federated("federated@facile.studio")

	recorder := h.send(http.MethodPatch, mePath, `{"password":"correct horse battery"}`, token)
	if recorder.Code != http.StatusOK {
		t.Fatalf("set a first password: %d %s", recorder.Code, recorder.Body.String())
	}

	login := h.send(http.MethodPost, porte.RouteLoginLocal,
		`{"email":"federated@facile.studio","password":"correct horse battery"}`, "")
	if login.Code != http.StatusOK {
		t.Fatalf("the first password does not sign in: %d %s", login.Code, login.Body.String())
	}
}

// The current password is what turns the request into a replacement. porte
// confirms it, rotates the caller's session and writes the new cookie itself.
func TestChangingAPasswordConfirmsTheCurrentOneAndRotatesTheCookie(t *testing.T) {
	h := newHarness(t)
	_, token := h.register("owner@facile.studio", "old passphrase")

	recorder := h.send(http.MethodPatch, mePath,
		`{"password":"new passphrase","current_password":"old passphrase"}`, token)
	if recorder.Code != http.StatusOK {
		t.Fatalf("change the password: %d %s", recorder.Code, recorder.Body.String())
	}
	if h.rotatedToken(recorder) == token {
		t.Fatalf("the session was not rotated")
	}

	login := h.send(http.MethodPost, porte.RouteLoginLocal,
		`{"email":"owner@facile.studio","password":"new passphrase"}`, "")
	if login.Code != http.StatusOK {
		t.Fatalf("the new password does not sign in: %d %s", login.Code, login.Body.String())
	}
}

// porte's ErrWrongPassword is a refusal the caller can act on, not an outage:
// tronc answers 500 to any error it cannot recognise, so the mapping is the
// difference between 401 and a pager.
func TestAWrongCurrentPasswordIsRefusedWithUnauthorized(t *testing.T) {
	h := newHarness(t)
	_, token := h.register("owner@facile.studio", "old passphrase")

	recorder := h.send(http.MethodPatch, mePath,
		`{"password":"new passphrase","current_password":"not the old one"}`, token)
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d %s", recorder.Code, recorder.Body.String())
	}
	if code := errorCode(t, recorder); code != "unauthenticated" {
		t.Fatalf("expected unauthenticated, got %q", code)
	}
	if message := errorMessage(t, recorder); !strings.Contains(message, "current password") {
		t.Fatalf("porte's empty message reached the caller: %q", message)
	}

	login := h.send(http.MethodPost, porte.RouteLoginLocal,
		`{"email":"owner@facile.studio","password":"old passphrase"}`, "")
	if login.Code != http.StatusOK {
		t.Fatalf("the refused change replaced the password anyway: %d", login.Code)
	}
}

// A replacement sent down the first-password path is porte's 409. The caller
// omitted a field rather than losing a race, so the answer names the field.
func TestAReplacementWithoutTheCurrentPasswordIsRefusedByName(t *testing.T) {
	h := newHarness(t)
	_, token := h.register("owner@facile.studio", "old passphrase")

	recorder := h.send(http.MethodPatch, mePath, `{"password":"new passphrase"}`, token)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d %s", recorder.Code, recorder.Body.String())
	}
	if code := errorCode(t, recorder); code != "invalid_argument" {
		t.Fatalf("expected invalid_argument, got %q", code)
	}
	if message := errorMessage(t, recorder); !strings.Contains(message, "current_password") {
		t.Fatalf("the answer does not name the missing field: %q", message)
	}
}

// A password's replacement must not leave credentials minted by the old one
// alive, and must leave the browser that made the change signed in.
func TestChangingAPasswordEndsTheOtherLoginsButNotTheCallers(t *testing.T) {
	h := newHarness(t)
	userID, token := h.register("owner@facile.studio", "old passphrase")
	other := h.login(userID)

	recorder := h.send(http.MethodPatch, mePath,
		`{"password":"new passphrase","current_password":"old passphrase"}`, token)
	if recorder.Code != http.StatusOK {
		t.Fatalf("change the password: %d %s", recorder.Code, recorder.Body.String())
	}

	if stale := h.send(http.MethodGet, mePath, "", other); stale.Code != http.StatusUnauthorized {
		t.Fatalf("the other browser is still signed in: %d %s", stale.Code, stale.Body.String())
	}
	if dead := h.send(http.MethodGet, mePath, "", token); dead.Code != http.StatusUnauthorized {
		t.Fatalf("the replaced session still authenticates: %d", dead.Code)
	}
	if live := h.send(http.MethodGet, mePath, "", h.rotatedToken(recorder)); live.Code != http.StatusOK {
		t.Fatalf("the caller was signed out by their own change: %d %s", live.Code, live.Body.String())
	}
}
