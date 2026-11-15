package capsule

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSealSendsOnlyKnownFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/pastes" {
			http.NotFound(w, r)
			return
		}
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		var req createRequest
		if err := dec.Decode(&req); err != nil || !req.BurnAfterRead || req.Content != "cipher" || req.ExpiresIn != "24h" {
			http.Error(w, "bad", http.StatusBadRequest)
			return
		}
		if err := json.NewEncoder(w).Encode(map[string]any{"id": "cap_0123456789abcdef", "delete_token": "secret"}); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	sealed, err := New(server.URL).Seal(context.Background(), "cipher", "24h")
	if err != nil || sealed.ID != "cap_0123456789abcdef" {
		t.Fatalf("seal: %v %+v", err, sealed)
	}
}

func TestSealReportsRefusal(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "slow down", http.StatusTooManyRequests)
	}))
	t.Cleanup(server.Close)

	if _, err := New(server.URL).Seal(context.Background(), "cipher", "1h"); err == nil {
		t.Fatal("a refused seal must fail")
	}
}

func TestValidExpiry(t *testing.T) {
	for _, ok := range []string{"1h", "24h", "7d", "30d"} {
		if !ValidExpiry(ok) {
			t.Fatalf("%s must be valid", ok)
		}
	}
	if ValidExpiry("2d") || ValidExpiry("") {
		t.Fatal("unknown expiries must be refused")
	}
}
