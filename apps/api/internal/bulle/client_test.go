package bulle

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/FacileStudio/porte"
)

type memIDs struct {
	mu   sync.Mutex
	rows []porte.StoredIdentity
}

func (m *memIDs) ListByUser(_ context.Context, userID int64) ([]porte.StoredIdentity, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []porte.StoredIdentity
	for _, row := range m.rows {
		if row.UserID == userID {
			out = append(out, row)
		}
	}
	return out, nil
}

func (m *memIDs) Save(_ context.Context, identity porte.StoredIdentity) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, row := range m.rows {
		if row.UserID == identity.UserID && row.Provider == identity.Provider {
			m.rows[i] = identity
		}
	}
	return nil
}

type fakes struct {
	registre, bulle *httptest.Server
	noKey           bool
	exchanges       int
}

func newFakes(t *testing.T) *fakes {
	t.Helper()
	f := &fakes{}
	f.registre = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/.well-known/openid-configuration" {
			_ = json.NewEncoder(w).Encode(map[string]string{"token_endpoint": f.registre.URL + "/token"})
			return
		}
		_ = r.ParseForm()
		switch r.Form.Get("grant_type") {
		case "refresh_token":
			_ = json.NewEncoder(w).Encode(map[string]any{"access_token": "fresh-at", "expires_in": 300})
		case exchangeGrant:
			f.exchanges++
			if r.Form.Get("subject_token") != "fresh-at" || r.Form.Get("audience") != "bulle" {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid_grant"})
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"access_token": "bulle-jwt", "expires_in": 300})
		}
	}))
	f.bulle = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer bulle-jwt" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		switch {
		case r.URL.Path == "/api/llm":
			_ = json.NewEncoder(w).Encode(Offer{Backend: "anthropic", Model: "m", Persona: "You are Pixel."})
		case f.noKey:
			w.WriteHeader(http.StatusConflict)
		case r.URL.Path == "/api/llm/anthropic/v1/messages" && r.Header.Get("anthropic-version") != "":
			var body map[string]any
			_ = json.NewDecoder(r.Body).Decode(&body)
			if !strings.Contains(body["system"].(string), "Pixel") {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"content": []map[string]string{{"type": "text", "text": " Shorter. "}}})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(f.registre.Close)
	t.Cleanup(f.bulle.Close)
	return f
}

func TestCompleteRefreshesExchangesAndCaches(t *testing.T) {
	f := newFakes(t)
	ids := &memIDs{rows: []porte.StoredIdentity{{UserID: 1, Provider: f.registre.URL,
		Tokens: porte.TokenSet{AccessToken: "old", RefreshToken: "rt", Expiry: time.Now().Add(-time.Minute)}}}}
	c := New(Config{BulleURL: f.bulle.URL, Issuer: f.registre.URL, ClientID: "courrier", ClientSecret: "s"}, ids, nil)
	ctx := context.Background()
	offer, err := c.Offer(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	text, err := c.Complete(ctx, 1, offer, Prompt{System: offer.Persona, User: "x", MaxTokens: 10})
	if err != nil || text != "Shorter." {
		t.Fatalf("complete = %q, %v", text, err)
	}
	if f.exchanges != 1 {
		t.Fatalf("the Bulle token must be cached, got %d exchanges", f.exchanges)
	}
	if ids.rows[0].Tokens.AccessToken != "fresh-at" {
		t.Fatal("the refreshed Registre token must be saved")
	}
	f.noKey = true
	if _, err := c.Complete(ctx, 1, offer, Prompt{User: "x"}); !errors.Is(err, ErrNoKey) {
		t.Fatalf("409 must read as no key, got %v", err)
	}
}

func TestUnlinkedUserAndUnconfiguredClient(t *testing.T) {
	f := newFakes(t)
	c := New(Config{BulleURL: f.bulle.URL, Issuer: f.registre.URL, ClientID: "courrier"}, &memIDs{}, nil)
	if _, err := c.Offer(context.Background(), 9); !errors.Is(err, ErrNotLinked) {
		t.Fatalf("a password-only user is not linked, got %v", err)
	}
	if New(Config{Issuer: "x", ClientID: "y"}, &memIDs{}, nil) != nil {
		t.Fatal("no BULLE_URL must leave the helper off")
	}
}
