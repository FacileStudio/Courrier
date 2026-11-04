package mail

import (
	"context"

	"github.com/FacileStudio/Courrier/apps/api/internal/spaceid"
)

// templateSpace resolves the space a template request names: nil for the
// personal scope, else the canonical id of a space the caller belongs to.
func (s *Service) templateSpace(ctx context.Context, userID int64, raw *string) (*string, error) {
	spaceID, err := spaceid.Optional(raw)
	if err != nil || spaceID == nil {
		return nil, err
	}
	return spaceID, s.requireSpaceMember(ctx, *spaceID, userID)
}
