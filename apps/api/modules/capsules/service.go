// Package capsules lets a Courrier user put a Capsule link in a mail: a secret
// that burns on its first read. The browser encrypts and keeps the key;
// Courrier relays ciphertext to Capsule and never holds a key.
package capsules

import (
	"context"
	"log/slog"

	"github.com/FacileStudio/Courrier/apps/api/internal/capsule"
	"github.com/FacileStudio/tronc/errors"
)

// Service relays capsules to one Capsule deployment.
type Service struct {
	capsule *capsule.Client
	log     *slog.Logger
}

// NewService builds the capsule relay.
func NewService(client *capsule.Client, logger *slog.Logger) *Service {
	return &Service{capsule: client, log: logger}
}

// Seal stores ciphertext as a burn-after-read capsule.
func (s *Service) Seal(ctx context.Context, req SealRequest) (SealResponse, error) {
	if req.Ciphertext == "" {
		return SealResponse{}, errors.Invalid("ciphertext is required")
	}
	if !capsule.ValidExpiry(req.ExpiresIn) {
		return SealResponse{}, errors.Invalid("expires_in must be one of: 1h, 24h, 7d, 30d")
	}
	sealed, err := s.capsule.Seal(ctx, req.Ciphertext, req.ExpiresIn)
	if err != nil {
		s.log.Warn("capsule relay failed", slog.Any("error", err))
		return SealResponse{}, errors.Unavailable("Capsule did not answer")
	}
	return SealResponse{ID: sealed.ID, URL: s.capsule.BaseURL() + "/" + sealed.ID, ExpiresAt: sealed.ExpiresAt}, nil
}
