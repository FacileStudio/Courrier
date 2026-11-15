// Package capsule relays sealed secrets to Capsule's public paste API. The
// browser encrypts before anything reaches Courrier, so this package only
// carries ciphertext and never sees a key.
package capsule

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// MaxCiphertext is the largest ciphertext Courrier relays.
const MaxCiphertext = 1 << 20

var expiries = map[string]bool{"1h": true, "24h": true, "7d": true, "30d": true}

// ValidExpiry reports whether Capsule accepts expiresIn.
func ValidExpiry(expiresIn string) bool { return expiries[expiresIn] }

// Client talks to one Capsule deployment.
type Client struct {
	baseURL string
	http    *http.Client
}

// Sealed is a freshly created capsule.
type Sealed struct {
	ID        string     `json:"id"`
	ExpiresAt *time.Time `json:"expires_at"`
}

type createRequest struct {
	Content       string `json:"content"`
	BurnAfterRead bool   `json:"burn_after_read"`
	ExpiresIn     string `json:"expires_in"`
}

// New builds a client for the Capsule deployment at baseURL.
func New(baseURL string) *Client {
	return &Client{baseURL: baseURL, http: &http.Client{Timeout: 15 * time.Second}}
}

// BaseURL is the public origin of the Capsule deployment.
func (c *Client) BaseURL() string { return c.baseURL }

// Seal stores ciphertext as a capsule that burns on its first read.
func (c *Client) Seal(ctx context.Context, ciphertext, expiresIn string) (Sealed, error) {
	payload, err := json.Marshal(createRequest{Content: ciphertext, BurnAfterRead: true, ExpiresIn: expiresIn})
	if err != nil {
		return Sealed{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/pastes", bytes.NewReader(payload))
	if err != nil {
		return Sealed{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return Sealed{}, fmt.Errorf("capsule: POST /api/pastes: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 300 {
		return Sealed{}, fmt.Errorf("capsule: POST /api/pastes answered %d", resp.StatusCode)
	}
	var sealed Sealed
	if err := json.NewDecoder(io.LimitReader(resp.Body, 64<<10)).Decode(&sealed); err != nil {
		return Sealed{}, fmt.Errorf("capsule: unreadable answer: %w", err)
	}
	if sealed.ID == "" {
		return Sealed{}, fmt.Errorf("capsule: answer has no id")
	}
	return sealed, nil
}
