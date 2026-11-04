package accounts

import (
	"cmp"
	"context"

	"github.com/FacileStudio/Courrier/apps/api/internal/spaceid"
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

// prepareCreate validates the request, fills the default ports and resolves
// the space: "" is the personal scope, any other id must be canonical and
// name a space the caller belongs to.
func (s *Service) prepareCreate(ctx context.Context, userID int64, req CreateAccountRequest) (CreateAccountRequest, error) {
	if err := validateCreateAccount(req); err != nil {
		return req, err
	}
	req.IMAPPort = cmp.Or(req.IMAPPort, 993)
	req.SMTPPort = cmp.Or(req.SMTPPort, 587)
	spaceID, err := spaceid.Optional(req.SpaceID)
	if err != nil || spaceID == nil {
		req.SpaceID = nil
		return req, err
	}
	req.SpaceID = spaceID
	return req, s.requireSpaceMember(ctx, *spaceID, userID)
}

func (s *Service) Create(ctx context.Context, userID int64, req CreateAccountRequest) (schemas.Account, error) {
	req, err := s.prepareCreate(ctx, userID, req)
	if err != nil {
		return schemas.Account{}, err
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
