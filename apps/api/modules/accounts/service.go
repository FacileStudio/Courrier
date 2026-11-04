package accounts

import (
	"context"

	"github.com/FacileStudio/Courrier/apps/api/internal/crypto"
	"github.com/FacileStudio/Courrier/apps/api/internal/plans"
	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"

	"gorm.io/gorm"
)

// Service is the accounts module's data access: it owns the accounts table,
// encrypts IMAP and SMTP passwords at rest and refuses an account the scope's
// plan has no room for.
type Service struct {
	orm           *gorm.DB
	encryptionKey []byte
	plans         *plans.Gate
}

// NewService builds an accounts Service over the database, the app's
// encryption key and the plan gate, which may be nil when no plan applies.
func NewService(orm *gorm.DB, encryptionKey []byte, gate *plans.Gate) *Service {
	return &Service{orm: orm, encryptionKey: encryptionKey, plans: gate}
}

func (s *Service) encryptPassword(plain string) (string, error) {
	if len(s.encryptionKey) == 0 || plain == "" {
		return plain, nil
	}
	return crypto.Encrypt(plain, s.encryptionKey)
}

func (s *Service) DecryptPassword(cipher string) (string, error) {
	if len(s.encryptionKey) == 0 || cipher == "" {
		return cipher, nil
	}
	return crypto.Decrypt(cipher, s.encryptionKey)
}

func (s *Service) requireSpaceMember(ctx context.Context, spaceID string, userID int64) error {
	var count int64
	if err := s.orm.WithContext(ctx).Model(&schemas.SpaceMember{}).
		Where("space_id = ? AND user_id = ?", spaceID, userID).
		Count(&count).Error; err != nil {
		return errors.Internal("failed to check space membership", err)
	}
	if count == 0 {
		return errors.Forbidden("not a member of this space")
	}
	return nil
}

func unsetDefault(db *gorm.DB, userID int64, spaceID *string) {
	if spaceID != nil && *spaceID != "" {
		db.Model(&schemas.Account{}).Where("space_id = ?", *spaceID).Update("is_default", false)
		return
	}
	db.Model(&schemas.Account{}).Where("user_id = ? AND space_id IS NULL", userID).Update("is_default", false)
}

func (s *Service) List(ctx context.Context, userID int64, spaceID *string) ([]schemas.Account, error) {
	var accounts []schemas.Account
	query := s.orm.WithContext(ctx)
	if spaceID != nil && *spaceID != "" {
		if err := s.requireSpaceMember(ctx, *spaceID, userID); err != nil {
			return nil, err
		}
		query = query.Where("space_id = ?", *spaceID)
	} else {
		query = query.Where("user_id = ? AND space_id IS NULL", userID)
	}
	if err := query.Order("is_default DESC, name ASC").Find(&accounts).Error; err != nil {
		return nil, errors.Internal("failed to list accounts", err)
	}
	return accounts, nil
}

func (s *Service) Get(ctx context.Context, userID int64, accountID int64) (schemas.Account, error) {
	var account schemas.Account
	if err := s.orm.WithContext(ctx).Where("id = ?", accountID).First(&account).Error; err != nil {
		return schemas.Account{}, errors.NotFound("account not found")
	}
	if account.SpaceID != nil && *account.SpaceID != "" {
		if err := s.requireSpaceMember(ctx, *account.SpaceID, userID); err != nil {
			return schemas.Account{}, err
		}
	} else if account.UserID != userID {
		return schemas.Account{}, errors.NotFound("account not found")
	}
	return account, nil
}

func applyBasicUpdates(account *schemas.Account, req UpdateAccountRequest) {
	if req.Name != nil {
		account.Name = *req.Name
	}
	if req.Email != nil {
		account.Email = *req.Email
	}
	if req.IMAPHost != nil {
		account.IMAPHost = *req.IMAPHost
	}
	if req.IMAPPort != nil {
		account.IMAPPort = *req.IMAPPort
	}
	if req.IMAPUser != nil {
		account.IMAPUser = *req.IMAPUser
	}
	if req.SMTPHost != nil {
		account.SMTPHost = *req.SMTPHost
	}
	if req.SMTPPort != nil {
		account.SMTPPort = *req.SMTPPort
	}
	if req.SMTPUser != nil {
		account.SMTPUser = *req.SMTPUser
	}
	if req.Signature != nil {
		account.Signature = *req.Signature
	}
}

func (s *Service) applyAccountUpdates(account *schemas.Account, req UpdateAccountRequest) error {
	applyBasicUpdates(account, req)
	if req.IMAPPassword != nil {
		enc, err := s.encryptPassword(*req.IMAPPassword)
		if err != nil {
			return errors.Internal("failed to encrypt IMAP password", err)
		}
		account.IMAPPassword = enc
	}
	if req.SMTPPassword != nil {
		enc, err := s.encryptPassword(*req.SMTPPassword)
		if err != nil {
			return errors.Internal("failed to encrypt SMTP password", err)
		}
		account.SMTPPassword = enc
	}
	return nil
}

func (s *Service) Update(ctx context.Context, userID int64, accountID int64, req UpdateAccountRequest) (schemas.Account, error) {
	account, err := s.Get(ctx, userID, accountID)
	if err != nil {
		return schemas.Account{}, err
	}
	if err := s.applyAccountUpdates(&account, req); err != nil {
		return schemas.Account{}, err
	}
	if req.IsDefault != nil && *req.IsDefault {
		unsetDefault(s.orm.WithContext(ctx), account.UserID, account.SpaceID)
		account.IsDefault = true
	}
	if err := s.orm.WithContext(ctx).Save(&account).Error; err != nil {
		return schemas.Account{}, errors.Internal("failed to update account", err)
	}
	return account, nil
}

func (s *Service) Delete(ctx context.Context, userID int64, accountID int64) error {
	account, err := s.Get(ctx, userID, accountID)
	if err != nil {
		return err
	}
	result := s.orm.WithContext(ctx).Where("id = ?", account.ID).Delete(&schemas.Account{})
	if result.Error != nil {
		return errors.Internal("failed to delete account", result.Error)
	}
	s.orm.WithContext(ctx).Where("account_id = ?", accountID).Delete(&schemas.Folder{})
	s.orm.WithContext(ctx).Where("account_id = ?", accountID).Delete(&schemas.Email{})
	s.orm.WithContext(ctx).Where("account_id = ?", accountID).Delete(&schemas.Contact{})
	return nil
}
