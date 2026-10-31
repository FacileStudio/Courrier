package mail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/FacileStudio/Courrier/apps/api/internal/authcontext"
	"github.com/FacileStudio/Courrier/apps/api/schemas"

	"github.com/go-chi/chi/v5"
)

func (e *mailEnv) call(userID int, method, path, contentType string, body []byte) *httptest.ResponseRecorder {
	router := chi.NewRouter()
	router.Route("/accounts/{accountId}/mail", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				id := authcontext.Identity{UserID: strconv.Itoa(userID)}
				next.ServeHTTP(w, req.WithContext(authcontext.WithIdentity(req.Context(), id)))
			})
		})
		registerWriteRoutes(r, e.svc)
	})
	req := httptest.NewRequest(method, fmt.Sprintf("/accounts/%d/mail%s", e.account.ID, path), bytes.NewReader(body))
	req.Header.Set("Content-Type", contentType)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func (e *mailEnv) json(t *testing.T, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	raw, _ := json.Marshal(body)
	return e.call(1, method, path, "application/json", raw)
}

func decodeBody(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var out map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("body %q: %v", rec.Body.String(), err)
	}
	return out
}

func expectStatus(t *testing.T, rec *httptest.ResponseRecorder, want int) {
	t.Helper()
	if rec.Code != want {
		t.Fatalf("status %d, want %d: %s", rec.Code, want, rec.Body.String())
	}
}

func TestFolderRoutesCoverCreateRenameDelete(t *testing.T) {
	env := newMailEnv(t)
	rec := env.json(t, "POST", "/folders", map[string]string{"name": "Clients"})
	expectStatus(t, rec, http.StatusCreated)
	id := int64(decodeBody(t, rec)["id"].(float64))

	expectStatus(t, env.json(t, "POST", "/folders", map[string]string{"name": "Clients"}), http.StatusConflict)
	expectStatus(t, env.json(t, "POST", "/folders", map[string]string{"name": "Trash"}), http.StatusBadRequest)
	expectStatus(t, env.json(t, "POST", "/folders", map[string]any{"name": "x", "extra": 1}), http.StatusBadRequest)

	rec = env.json(t, "PATCH", fmt.Sprintf("/folders/%d", id), map[string]string{"name": "Customers"})
	expectStatus(t, rec, http.StatusOK)
	if decodeBody(t, rec)["path"] != "Customers" {
		t.Errorf("rename body %s", rec.Body.String())
	}
	expectStatus(t, env.json(t, "PATCH", "/folders/abc", map[string]string{"name": "x"}), http.StatusBadRequest)
	expectStatus(t, env.json(t, "PATCH", "/folders/999999", map[string]string{"name": "x"}), http.StatusNotFound)
	expectStatus(t, env.json(t, "DELETE", fmt.Sprintf("/folders/%d", id), nil), http.StatusOK)
	expectStatus(t, env.json(t, "DELETE", fmt.Sprintf("/folders/%d", id), nil), http.StatusNotFound)
}

func TestDeleteFolderRouteHonoursForce(t *testing.T) {
	env := newMailEnv(t)
	env.seed(t, "Projects", "inside")
	path := fmt.Sprintf("/folders/%d", env.folder(t, "Projects").ID)
	expectStatus(t, env.json(t, "DELETE", path, nil), http.StatusConflict)
	expectStatus(t, env.json(t, "DELETE", path+"?force=true", nil), http.StatusOK)
}

func TestMoveRouteValidatesAndMoves(t *testing.T) {
	env := newMailEnv(t)
	mail := env.seed(t, "INBOX", "route move")
	projects := env.folder(t, "Projects")
	expectStatus(t, env.json(t, "POST", "/emails/move", map[string]any{"email_ids": []int64{}, "folder_type": "trash"}), http.StatusBadRequest)
	expectStatus(t, env.json(t, "POST", "/emails/move", map[string]any{"email_ids": make([]int64, 201), "folder_type": "trash"}), http.StatusBadRequest)
	expectStatus(t, env.json(t, "POST", "/emails/move", map[string]any{"email_ids": []int64{mail.ID}}), http.StatusBadRequest)
	expectStatus(t, env.json(t, "POST", "/emails/move", map[string]any{"email_ids": []int64{mail.ID}, "folder_id": projects.ID}), http.StatusOK)
	if reload(t, mail.ID).FolderID != projects.ID {
		t.Error("the message did not move")
	}
}

func TestRoutesAnswer404ForAnotherUsersAccount(t *testing.T) {
	env := newMailEnv(t)
	body := []byte(`{"name":"Sneaky"}`)
	expectStatus(t, env.call(2, "POST", "/folders", "application/json", body), http.StatusNotFound)
	expectStatus(t, env.call(2, "POST", "/emails/move", "application/json", []byte(`{"email_ids":[1],"folder_type":"trash"}`)), http.StatusNotFound)
	expectStatus(t, env.call(2, "PUT", "/drafts/1", "application/json", []byte(`{}`)), http.StatusNotFound)
	if slices.Contains(env.fake.mailboxes(t), "Sneaky") {
		t.Error("a foreign user created a mailbox")
	}
}

func multipartBody(t *testing.T, fields map[string]string, files map[string]string) ([]byte, string) {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	for k, v := range fields {
		w.WriteField(k, v)
	}
	for name, content := range files {
		part, _ := w.CreateFormFile("attachments", name)
		part.Write([]byte(content))
	}
	w.Close()
	return buf.Bytes(), w.FormDataContentType()
}

func TestDraftRoutesAcceptMultipartAndUpdateInPlace(t *testing.T) {
	env := newMailEnv(t)
	body, ctype := multipartBody(t, map[string]string{"to": "a@x.test, b@x.test", "bcc": "h@x.test", "subject": "m1", "body": "x"}, map[string]string{"one.txt": "111"})
	rec := env.call(1, "POST", "/drafts", ctype, body)
	expectStatus(t, rec, http.StatusOK)
	id := int64(decodeBody(t, rec)["id"].(float64))
	row := reload(t, id)
	if !row.HasAttachments || !strings.Contains(row.BccAddresses, "h@x.test") || !strings.Contains(row.ToAddresses, "b@x.test") {
		t.Fatalf("draft row %+v", row)
	}
	var first schemas.Attachment
	testORM.Where("email_id = ?", id).First(&first)

	body, ctype = multipartBody(t, map[string]string{"to": "a@x.test", "subject": "m2", "body": "y", "remove_attachments": strconv.FormatInt(first.ID, 10)}, map[string]string{"two.txt": "222"})
	rec = env.call(1, "PUT", fmt.Sprintf("/drafts/%d", id), ctype, body)
	expectStatus(t, rec, http.StatusOK)
	if got := int64(decodeBody(t, rec)["id"].(float64)); got != id {
		t.Fatalf("id changed from %d to %d", id, got)
	}
	assertDraftFiles(t, env, id, map[string]string{"two.txt": "222"})
	expectStatus(t, env.json(t, "PUT", "/drafts/999999", map[string]any{"subject": "x"}), http.StatusNotFound)
	expectStatus(t, env.json(t, "PUT", "/drafts/abc", map[string]any{"subject": "x"}), http.StatusBadRequest)
}
