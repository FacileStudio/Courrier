package accounts

import (
	"context"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"

	"gorm.io/gorm"
)

func validateCreateAccount(req CreateAccountRequest) error {
	if req.Name == "" {
		return errors.Invalid("name is required")
	}
	if req.Email == "" {
		return errors.Invalid("email is required")
	}
	if req.IMAPHost == "" {
		return errors.Invalid("imap_host is required")
	}
	if req.SMTPHost == "" {
		return errors.Invalid("smtp_host is required")
	}
	return nil
}

func (s *Service) buildAccount(userID int64, req CreateAccountRequest) (schemas.Account, error) {
	encIMAPPw, err := s.encryptPassword(req.IMAPPassword)
	if err != nil {
		return schemas.Account{}, errors.Internal("failed to encrypt IMAP password", err)
	}
	encSMTPPw, err := s.encryptPassword(req.SMTPPassword)
	if err != nil {
		return schemas.Account{}, errors.Internal("failed to encrypt SMTP password", err)
	}
	return schemas.Account{
		UserID:       userID,
		SpaceID:      req.SpaceID,
		Name:         req.Name,
		Email:        req.Email,
		IMAPHost:     req.IMAPHost,
		IMAPPort:     req.IMAPPort,
		IMAPUser:     req.IMAPUser,
		IMAPPassword: encIMAPPw,
		SMTPHost:     req.SMTPHost,
		SMTPPort:     req.SMTPPort,
		SMTPUser:     req.SMTPUser,
		SMTPPassword: encSMTPPw,
		Signature:    req.Signature,
		IsDefault:    req.IsDefault,
	}, nil
}

func (s *Service) Create(ctx context.Context, userID int64, req CreateAccountRequest) (schemas.Account, error) {
	if err := validateCreateAccount(req); err != nil {
		return schemas.Account{}, err
	}
	if req.IMAPPort == 0 {
		req.IMAPPort = 993
	}
	if req.SMTPPort == 0 {
		req.SMTPPort = 587
	}
	if req.SpaceID != nil && *req.SpaceID != "" {
		if err := s.requireSpaceMember(ctx, *req.SpaceID, userID); err != nil {
			return schemas.Account{}, err
		}
	}
	account, err := s.buildAccount(userID, req)
	if err != nil {
		return schemas.Account{}, err
	}
	if err := s.orm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return s.insert(tx, &account)
	}); err != nil {
		return schemas.Account{}, err
	}
	return account, nil
}

// insert creates the account once its scope's plan has room for it. It runs
// inside the creating transaction, which the plan check locks the scope in.
func (s *Service) insert(tx *gorm.DB, account *schemas.Account) error {
	if err := s.plans.CheckMailbox(tx, account.UserID, account.SpaceID); err != nil {
		return err
	}
	if account.IsDefault {
		unsetDefault(tx, account.UserID, account.SpaceID)
	}
	if err := tx.Create(account).Error; err != nil {
		return errors.Internal("failed to create account", err)
	}
	return nil
}
