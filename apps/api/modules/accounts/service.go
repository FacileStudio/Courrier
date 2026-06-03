package accounts

import (
	"context"

	"api/internal/errors"
	"api/schemas"

	"gorm.io/gorm"
)

type Service struct {
	orm *gorm.DB
}

func NewService(orm *gorm.DB) *Service {
	return &Service{orm: orm}
}

func (s *Service) Create(ctx context.Context, userID int64, req CreateAccountRequest) (schemas.Account, error) {
	if req.Name == "" {
		return schemas.Account{}, errors.Invalid("name is required")
	}
	if req.Email == "" {
		return schemas.Account{}, errors.Invalid("email is required")
	}
	if req.IMAPHost == "" {
		return schemas.Account{}, errors.Invalid("imap_host is required")
	}
	if req.SMTPHost == "" {
		return schemas.Account{}, errors.Invalid("smtp_host is required")
	}

	if req.IMAPPort == 0 {
		req.IMAPPort = 993
	}
	if req.SMTPPort == 0 {
		req.SMTPPort = 587
	}

	account := schemas.Account{
		UserID:       userID,
		Name:         req.Name,
		Email:        req.Email,
		IMAPHost:     req.IMAPHost,
		IMAPPort:     req.IMAPPort,
		IMAPUser:     req.IMAPUser,
		IMAPPassword: req.IMAPPassword,
		SMTPHost:     req.SMTPHost,
		SMTPPort:     req.SMTPPort,
		SMTPUser:     req.SMTPUser,
		SMTPPassword: req.SMTPPassword,
		Signature:    req.Signature,
		IsDefault:    req.IsDefault,
	}

	if req.IsDefault {
		s.orm.WithContext(ctx).Model(&schemas.Account{}).Where("user_id = ?", userID).Update("is_default", false)
	}

	if err := s.orm.WithContext(ctx).Create(&account).Error; err != nil {
		return schemas.Account{}, errors.Internal("failed to create account", err)
	}
	return account, nil
}

func (s *Service) List(ctx context.Context, userID int64) ([]schemas.Account, error) {
	var accounts []schemas.Account
	if err := s.orm.WithContext(ctx).Where("user_id = ?", userID).Order("is_default DESC, name ASC").Find(&accounts).Error; err != nil {
		return nil, errors.Internal("failed to list accounts", err)
	}
	return accounts, nil
}

func (s *Service) Get(ctx context.Context, userID int64, accountID int64) (schemas.Account, error) {
	var account schemas.Account
	if err := s.orm.WithContext(ctx).Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		return schemas.Account{}, errors.NotFound("account not found")
	}
	return account, nil
}

func (s *Service) Update(ctx context.Context, userID int64, accountID int64, req UpdateAccountRequest) (schemas.Account, error) {
	var account schemas.Account
	if err := s.orm.WithContext(ctx).Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		return schemas.Account{}, errors.NotFound("account not found")
	}

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
	if req.IMAPPassword != nil {
		account.IMAPPassword = *req.IMAPPassword
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
	if req.SMTPPassword != nil {
		account.SMTPPassword = *req.SMTPPassword
	}
	if req.Signature != nil {
		account.Signature = *req.Signature
	}
	if req.IsDefault != nil && *req.IsDefault {
		s.orm.WithContext(ctx).Model(&schemas.Account{}).Where("user_id = ?", userID).Update("is_default", false)
		account.IsDefault = true
	}

	if err := s.orm.WithContext(ctx).Save(&account).Error; err != nil {
		return schemas.Account{}, errors.Internal("failed to update account", err)
	}
	return account, nil
}

func (s *Service) Delete(ctx context.Context, userID int64, accountID int64) error {
	result := s.orm.WithContext(ctx).Where("id = ? AND user_id = ?", accountID, userID).Delete(&schemas.Account{})
	if result.Error != nil {
		return errors.Internal("failed to delete account", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NotFound("account not found")
	}
	s.orm.WithContext(ctx).Where("account_id = ?", accountID).Delete(&schemas.Folder{})
	s.orm.WithContext(ctx).Where("account_id = ?", accountID).Delete(&schemas.Email{})
	return nil
}
