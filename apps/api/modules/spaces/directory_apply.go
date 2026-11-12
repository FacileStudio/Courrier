package spaces

import (
	"context"
	"log/slog"
	"strings"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/porte/directory"
	"github.com/FacileStudio/porte/kit"

	"gorm.io/gorm"
)

// Apply makes Courrier's spaces match a Registre snapshot, in one
// transaction, and changes nothing when they already do.
//
// Only linked spaces are touched, and inside them only members with a
// Registre identity are ever removed: a local password account added on a
// password-only day keeps its seat. A linked space missing from the snapshot
// is left alone and logged, because Registre v1 never deletes a space.
func (d *Directory) Apply(ctx context.Context, snapshot directory.Snapshot) error {
	return d.orm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ids := []string{""}
		for _, space := range snapshot.Spaces {
			if err := d.applySpace(tx, space); err != nil {
				return err
			}
			ids = append(ids, space.ID)
		}
		return d.reportMissing(tx, ids)
	})
}

// applySpace links one Registre space and sets its membership.
func (d *Directory) applySpace(tx *gorm.DB, space directory.Space) error {
	local, err := linkSpace(tx, space)
	if err != nil {
		return err
	}
	keep := []int64{0}
	for _, member := range space.Members {
		userID, err := d.resolveMember(tx, member)
		if err != nil {
			return err
		}
		if userID == 0 {
			continue
		}
		if err := setMembership(tx, local.ID, userID, localRole(member.Role)); err != nil {
			return err
		}
		keep = append(keep, userID)
	}
	return tx.Exec(`DELETE FROM space_members m USING porte_identities i
		WHERE m.space_id = ? AND i.user_id = m.user_id AND i.provider = ? AND m.user_id NOT IN ?`,
		local.ID, d.issuer, keep).Error
}

// linkSpace finds the local space for a Registre one, by link and then by an
// unlinked space of the same name, ignoring case, or creates it. The name
// always follows Registre.
func linkSpace(tx *gorm.DB, space directory.Space) (schemas.Space, error) {
	var local schemas.Space
	found := tx.Where("registre_space_id = ?", space.ID).Limit(1).Find(&local)
	if found.Error == nil && found.RowsAffected == 0 {
		found = tx.Where("registre_space_id IS NULL AND lower(name) = lower(?)", space.Name).
			Order("created_at").Limit(1).Find(&local)
	}
	if found.Error != nil {
		return local, found.Error
	}
	if found.RowsAffected == 0 {
		local = schemas.Space{Name: space.Name, RegistreSpaceID: &space.ID}
		return local, tx.Create(&local).Error
	}
	if local.Name == space.Name && local.RegistreSpaceID != nil {
		return local, nil
	}
	return local, tx.Model(&local).Updates(map[string]any{"name": space.Name, "registre_space_id": space.ID}).Error
}

// resolveMember finds the Courrier user for a Registre member: by porte
// identity, else by email, else by pre-provisioning one so the member's first
// sign-in lands on it. Zero means the member cannot be placed and is skipped.
func (d *Directory) resolveMember(tx *gorm.DB, member directory.Member) (int64, error) {
	var userID int64
	err := tx.Raw(`SELECT user_id FROM porte_identities WHERE provider = ? AND subject = ?`,
		d.issuer, member.Subject).Scan(&userID).Error
	if err != nil || userID != 0 || member.Subject == "" || strings.TrimSpace(member.Email) == "" {
		return userID, err
	}
	if erased, err := d.isErased(tx.Statement.Context, member.Subject); err != nil || erased {
		return 0, err
	}
	userID, err = d.userByEmail(tx, member)
	if err != nil || userID == 0 {
		return 0, err
	}
	return userID, tx.Exec(`INSERT INTO porte_identities (user_id, provider, subject) VALUES (?, ?, ?)
		ON CONFLICT (provider, subject) DO NOTHING`, userID, d.issuer, member.Subject).Error
}

// isErased reports whether Courrier already erased the member's identity, and
// logs the skip. Without WithErased nothing is erased.
func (d *Directory) isErased(ctx context.Context, subject string) (bool, error) {
	if d.erased == nil {
		return false, nil
	}
	erased, err := d.erased.Erased(ctx, d.issuer, subject)
	if erased {
		d.logger.Warn("registre: member was erased here, skipped", slog.String("subject", subject))
	}
	return erased, err
}

// userByEmail returns the user with the member's address, or a new one. A user
// already holding a different Registre identity is not taken over: zero.
func (d *Directory) userByEmail(tx *gorm.DB, member directory.Member) (int64, error) {
	email := strings.ToLower(strings.TrimSpace(member.Email))
	var user schemas.User
	found := tx.Where("lower(email) = ?", email).Limit(1).Find(&user)
	if found.Error != nil {
		return 0, found.Error
	}
	if found.RowsAffected == 0 {
		user = schemas.User{Email: email, Name: member.Name, OIDCPictureURL: kit.PictureURL(member.Picture)}
		err := tx.Create(&user).Error
		return user.ID, err
	}
	var linked int64
	if err := tx.Raw(`SELECT count(*) FROM porte_identities WHERE user_id = ? AND provider = ?`,
		user.ID, d.issuer).Scan(&linked).Error; err != nil || linked == 0 {
		return user.ID, err
	}
	d.logger.Warn("registre: member's email belongs to another identity, skipped", slog.String("email", email))
	return 0, nil
}

// setMembership adds the user to the space at role, or corrects their role.
func setMembership(tx *gorm.DB, spaceID string, userID int64, role string) error {
	if err := tx.Exec(`UPDATE space_members SET role = ? WHERE space_id = ? AND user_id = ? AND role <> ?`,
		role, spaceID, userID, role).Error; err != nil {
		return err
	}
	return tx.Exec(`INSERT INTO space_members (space_id, user_id, role, joined_at)
		SELECT ?, ?, ?, now()
		 WHERE NOT EXISTS (SELECT 1 FROM space_members WHERE space_id = ? AND user_id = ?)`,
		spaceID, userID, role, spaceID, userID).Error
}

// reportMissing logs linked spaces the snapshot no longer carries.
func (d *Directory) reportMissing(tx *gorm.DB, ids []string) error {
	var missing []string
	if err := tx.Model(&schemas.Space{}).
		Where("registre_space_id IS NOT NULL AND registre_space_id NOT IN ?", ids).
		Pluck("name", &missing).Error; err != nil {
		return err
	}
	if len(missing) > 0 {
		d.logger.Info("registre: linked spaces absent from the snapshot, left alone", slog.Any("spaces", missing))
	}
	return nil
}

// localRole maps a Registre role onto Courrier's, which uses the same three
// words; anything unknown is an ordinary member.
func localRole(role string) string {
	switch role {
	case directory.RoleOwner, directory.RoleAdmin:
		return role
	default:
		return schemas.SpaceRoleMember
	}
}
