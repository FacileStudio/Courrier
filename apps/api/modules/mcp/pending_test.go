package mcp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mcpkit "github.com/FacileStudio/tronc/mcp"

	"github.com/FacileStudio/Courrier/apps/api/modules/mail"
	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/porte/kit"
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
	return onPage(f, http.MethodPost, "/"+id+"/confirm", actor)
}

// onPage serves one confirm page request as the signed-in person, or as a
// delegated token when actor is set.
func onPage(f *fixture, method, path, actor string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	ctx := kit.WithIdentity(req.Context(), kit.Identity{UserID: f.userID})
	ctx = porte.WithIdentity(ctx, porte.Identity{UserID: f.userID, Actor: actor})
	router := chi.NewRouter()
	f.service.page.Routes(router)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req.WithContext(ctx))
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

// TestPendingListDescribesDraft checks the confirm page lists the open
// send with its mailbox and its draft as it is now.
func TestPendingListDescribesDraft(t *testing.T) {
	f := newFixture(t)
	out := pendingSend(t, f)
	got := onPage(f, http.MethodGet, "/", "")
	var list mcpkit.PendingList[PendingItem]
	if err := json.Unmarshal(got.Body.Bytes(), &list); err != nil || got.Code != http.StatusOK || len(list.Actions) != 1 {
		t.Fatalf("want one open action, got %d %s", got.Code, got.Body)
	}
	item := list.Actions[0]
	if item.Kind != kindSend || item.AccountID != f.accountID || item.AccountEmail == "" || item.Draft == nil || item.Draft.ID != out.DraftID || item.State != "ready" {
		t.Fatalf("want the send described with its mailbox and draft, got %+v", item)
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
