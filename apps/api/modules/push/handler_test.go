package push_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	pushsender "github.com/FacileStudio/Courrier/apps/api/internal/push"
	"github.com/FacileStudio/Courrier/apps/api/modules/push"
	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/porte/kit"
	"github.com/go-chi/chi/v5"
)

type signedIn struct{}

func (signedIn) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(porte.WithIdentity(r.Context(), porte.Identity{UserID: 7})))
	})
}

func (signedIn) Identity(_ context.Context, userID int64) (kit.Identity, error) {
	return kit.Identity{UserID: userID, Email: "ada@example.com"}, nil
}

func serve(t *testing.T, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	r := chi.NewRouter()
	push.RegisterRoutes(r, pushsender.NewHub(pushsender.HubConfig{}), signedIn{})
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(method, path, strings.NewReader(body)))
	return rec
}

func TestOnlyTheLiveStreamIsOnWithoutChannels(t *testing.T) {
	rec := serve(t, http.MethodGet, "/push", "")
	var status push.StatusResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &status); err != nil || rec.Code != http.StatusOK {
		t.Fatalf("status %d %s", rec.Code, rec.Body.String())
	}
	if !status.Enabled || status.Native || status.WebPushKey != "" {
		t.Fatalf("only the live stream should be on: %+v", status)
	}
	if rec := serve(t, http.MethodPost, "/push/test", ""); rec.Code != http.StatusNoContent {
		t.Fatalf("test send to the live stream: want 204, got %d", rec.Code)
	}
	if rec := serve(t, http.MethodPost, "/push/devices", `{"kind":"fcm","token":"t"}`); rec.Code != http.StatusConflict {
		t.Fatalf("fcm device without FCM: want 409, got %d", rec.Code)
	}
	if rec := serve(t, http.MethodPost, "/push/devices", `{"kind":"web","token":"t"}`); rec.Code != http.StatusBadRequest {
		t.Fatalf("web device without keys: want 400, got %d", rec.Code)
	}
}
