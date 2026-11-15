package capsules

import "time"

// SealRequest is the body of POST /capsules: ciphertext the browser already
// encrypted, and how long Capsule keeps it unopened (1h, 24h, 7d or 30d).
type SealRequest struct {
	Ciphertext string `json:"ciphertext"`
	ExpiresIn  string `json:"expires_in"`
}

// SealResponse names the new capsule. URL is its Capsule link without the key,
// which only the sender's browser holds.
type SealResponse struct {
	ID        string     `json:"id"`
	URL       string     `json:"url"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}
