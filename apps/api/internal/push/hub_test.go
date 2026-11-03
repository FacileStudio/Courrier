package push

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type memDevices struct{ devices []Device }

func (m *memDevices) SaveDevice(_ context.Context, d Device) error {
	m.devices = append(m.devices, d)
	return nil
}

func (m *memDevices) DeleteDevice(_ context.Context, userID int64, token string) error {
	kept := m.devices[:0]
	for _, d := range m.devices {
		if d.UserID != userID || d.Token != token {
			kept = append(kept, d)
		}
	}
	m.devices = kept
	return nil
}

func (m *memDevices) Devices(_ context.Context, userID int64) ([]Device, error) {
	var out []Device
	for _, d := range m.devices {
		if d.UserID == userID {
			out = append(out, d)
		}
	}
	return out, nil
}

func TestHubDeliversToFCMAndForgetsGoneTokens(t *testing.T) {
	var got []string
	fcm := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		got = append(got, string(body))
		if strings.Contains(string(body), `"token":"dead"`) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":{"status":"NOT_FOUND","details":[{"errorCode":"UNREGISTERED"}]}}`))
		}
	}))
	defer fcm.Close()

	devices := &memDevices{}
	hub := NewHub(HubConfig{FCM: &FCM{endpoint: fcm.URL, client: fcm.Client()}, Devices: devices})
	ctx := context.Background()
	for _, token := range []string{"live", "dead"} {
		if err := hub.Register(ctx, Device{UserID: 1, Kind: KindFCM, Token: token}); err != nil {
			t.Fatal(err)
		}
	}
	if err := hub.Register(ctx, Device{UserID: 1, Kind: KindWeb, Token: "x"}); err != ErrUnsupported {
		t.Fatalf("web push without VAPID keys: %v", err)
	}

	hub.Notify(ctx, 1, Notification{Title: "New mail from Ada", Body: "hello@facile.studio", Path: "mail"})

	if len(got) != 2 || !strings.Contains(got[0], `"path":"mail"`) {
		t.Fatalf("fcm requests: %q", got)
	}
	left, _ := devices.Devices(ctx, 1)
	if len(left) != 1 || left[0].Token != "live" {
		t.Fatalf("devices after delivery: %+v", left)
	}
}

func TestHubWithoutChannelsKeepsOnlyTheLiveStream(t *testing.T) {
	hub := NewHub(HubConfig{Devices: &memDevices{}})
	if !hub.Enabled() || hub.Native() || hub.WebPushKey() != "" {
		t.Fatal("a hub with no FCM and no VAPID keys must report only the live stream on")
	}
	if err := hub.Register(context.Background(), Device{UserID: 1, Kind: KindFCM, Token: "t"}); err != ErrUnsupported {
		t.Fatalf("fcm device without FCM: %v", err)
	}
	var none *Hub
	none.Notify(context.Background(), 1, Notification{Title: "x"})
	if none.Enabled() {
		t.Fatal("a nil hub is off")
	}
}
