package push_test

import (
	"bufio"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	pushsender "github.com/FacileStudio/Courrier/apps/api/internal/push"
	"github.com/FacileStudio/Courrier/apps/api/modules/push"
	"github.com/go-chi/chi/v5"
)

func TestStreamOutlivesTheServerTimeoutsAndCarriesNotifications(t *testing.T) {
	hub := pushsender.NewHub(pushsender.HubConfig{})
	r := chi.NewRouter()
	push.RegisterRoutes(r, hub, signedIn{})
	srv := httptest.NewUnstartedServer(r)
	srv.Config.ReadTimeout = 100 * time.Millisecond
	srv.Config.WriteTimeout = 100 * time.Millisecond
	srv.Start()
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL+"/push/stream", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("open stream: %v %v", err, resp)
	}
	defer resp.Body.Close()
	lines := bufio.NewScanner(resp.Body)
	if !lines.Scan() || lines.Text() != ": ok" {
		t.Fatalf("first frame %q", lines.Text())
	}
	time.Sleep(300 * time.Millisecond)
	hub.Notify(ctx, 7, pushsender.Notification{Title: "New mail from Ada", Body: "me@example.com", Path: "mail"})
	for lines.Scan() {
		if data, ok := strings.CutPrefix(lines.Text(), "data: "); ok {
			if !strings.Contains(data, `"title":"New mail from Ada"`) {
				t.Fatalf("data %s", data)
			}
			return
		}
	}
	t.Fatalf("stream ended before the notification: %v", lines.Err())
}
