package spaces

import (
	"context"
	"strings"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"
)

// WithRegistre hands space management to the Registre at issuer: creating a
// space is refused, and so is every edit of a space linked to Registre. A
// space not linked yet stays editable until the next sync links it. Without
// it, as on a password-only instance, spaces are managed here as before.
func (s *Service) WithRegistre(issuer string) *Service {
	s.registre = strings.TrimSuffix(issuer, "/")
	return s
}

// ManageURL is Registre's spaces page, or empty when Courrier manages its own
// spaces.
func (s *Service) ManageURL() string {
	if s.registre == "" {
		return ""
	}
	return s.registre + "/spaces"
}

// spaceManageURL is Registre's page for one linked space, or empty.
func (s *Service) spaceManageURL(space schemas.Space) string {
	if s.registre == "" || space.RegistreSpaceID == nil {
		return ""
	}
	return s.registre + "/spaces/" + *space.RegistreSpaceID
}

// refuseCreate answers 409 when Registre owns spaces.
func (s *Service) refuseCreate() error {
	if s.registre == "" {
		return nil
	}
	return errors.Conflict("spaces are created in Registre: " + s.ManageURL())
}

// refuseLinked answers 409 for a space Registre owns, naming where to edit it.
func (s *Service) refuseLinked(ctx context.Context, spaceID string) error {
	if s.registre == "" {
		return nil
	}
	var space schemas.Space
	if err := s.orm.WithContext(ctx).Select("id", "registre_space_id").
		Where("id = ?", spaceID).Limit(1).Find(&space).Error; err != nil {
		return errors.Internal("failed to load the space", err)
	}
	if url := s.spaceManageURL(space); url != "" {
		return errors.Conflict("this space is managed in Registre: " + url)
	}
	return nil
}
