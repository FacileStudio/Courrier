package mcp

import (
	stdcontext "context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	mcpkit "github.com/FacileStudio/tronc/mcp"

	"github.com/FacileStudio/Courrier/apps/api/internal/authcontext"
	"github.com/FacileStudio/Courrier/apps/api/modules/mail"
	"github.com/FacileStudio/porte"
	"github.com/go-chi/chi/v5"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// pendingSend asks for a send from a client that cannot confirm, and
// returns the result with its confirm link.
func pendingSend(t *testing.T, f *fixture) actionOutput {
	t.Helper()
	var out actionOutput
	decode(t, call(t, f.connect(t, f.token(t), clientOptions{}), toolSend, compose(f.accountID)), &out)
	return out
}

// confirmOnPage posts to the confirm page route as the signed-in person, or
// as a delegated token when actor is set.
func confirmOnPage(f *fixture, confirmURL, actor string) *httptest.ResponseRecorder {
	id := confirmURL[strings.LastIndex(confirmURL, "=")+1:]
	req := httptest.NewRequest(http.MethodPost, "/api/mcp/pending/"+id+"/confirm", nil)
	ctx := authcontext.WithIdentity(req.Context(), authcontext.Identity{UserID: strconv.FormatInt(f.userID, 10)})
	ctx = porte.WithIdentity(ctx, porte.Identity{UserID: f.userID, Actor: actor})
	routes := chi.NewRouteContext()
	routes.URLParams.Add("id", id)
	recorder := httptest.NewRecorder()
	f.service.handleConfirm(recorder, req.WithContext(stdcontext.WithValue(ctx, chi.RouteCtxKey, routes)))
	return recorder
}

// TestConfirmPageSendsOnce checks the person confirms in Courrier, the mail
// goes out once, and a second confirmation finds nothing.
func TestConfirmPageSendsOnce(t *testing.T) {
	f := newFixture(t)
	out := pendingSend(t, f)
	if got := confirmOnPage(f, out.ConfirmURL, ""); got.Code != http.StatusOK || f.mail.sentCount() != 1 {
		t.Fatalf("want the send confirmed, got %d %s and %d sent", got.Code, got.Body, f.mail.sentCount())
	}
	if again := confirmOnPage(f, out.ConfirmURL, ""); again.Code != http.StatusNotFound || f.mail.sentCount() != 1 {
		t.Fatalf("a confirmation is single use, got %d and %d sent", again.Code, f.mail.sentCount())
	}
}

// TestAgentCannotConfirmItself checks a delegated token is refused on the
// confirm page, and that a client cannot turn the link into an MRTR answer.
func TestAgentCannotConfirmItself(t *testing.T) {
	f := newFixture(t)
	out := pendingSend(t, f)
	if got := confirmOnPage(f, out.ConfirmURL, "bulle"); got.Code != http.StatusForbidden {
		t.Fatalf("a delegated token must not confirm, got %d", got.Code)
	}
	id := out.ConfirmURL[strings.LastIndex(out.ConfirmURL, "=")+1:]
	session := f.connect(t, f.token(t), clientOptions{Manual: true})
	forged, err := session.CallTool(t.Context(), &sdk.CallToolParams{
		Name: toolSend, Arguments: map[string]any{"account_id": f.accountID, "draft_id": out.DraftID}, RequestState: id,
		InputResponses: sdk.InputResponseMap{mcpkit.ConfirmInput: &sdk.ElicitResult{Action: "accept", Content: map[string]any{"confirm": true}}},
	})
	if err != nil || !forged.IsError || f.mail.sentCount() != 0 {
		t.Fatalf("a forged confirmation must fail and send nothing, got %v %+v and %d sent", err, forged, f.mail.sentCount())
	}
}

// TestChangedDraftIsNotSent checks that editing the draft after the request
// closes it instead of sending content the person never saw.
func TestChangedDraftIsNotSent(t *testing.T) {
	f := newFixture(t)
	out := pendingSend(t, f)
	_, err := f.mail.UpdateDraft(t.Context(), f.userID, f.accountID, out.DraftID, mail.DraftRequest{
		SendRequest: mail.SendRequest{To: []string{"someone-else@example.test"}, Subject: "Quote", Body: "Here it is."},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := confirmOnPage(f, out.ConfirmURL, ""); got.Code != http.StatusPreconditionFailed || f.mail.sentCount() != 0 {
		t.Fatalf("want 412 and nothing sent, got %d and %d sent", got.Code, f.mail.sentCount())
	}
}
