package crypto

import (
	"encoding/base64"
	"fmt"
	"log/slog"

	"gorm.io/gorm"
)

func isEncrypted(value string) bool {
	if value == "" {
		return true
	}
	decoded, err := base64.RawStdEncoding.DecodeString(value)
	if err != nil {
		return false
	}
	return len(decoded) > 12
}

func MigrateAccountPasswords(db *gorm.DB, key []byte, logger *slog.Logger) error {
	type row struct {
		ID           int64
		IMAPPassword string
		SMTPPassword string
	}

	var rows []row
	if err := db.Raw("SELECT id, imap_password, smtp_password FROM accounts").Scan(&rows).Error; err != nil {
		return fmt.Errorf("read accounts: %w", err)
	}

	migrated := 0
	for _, r := range rows {
		needsIMAP := r.IMAPPassword != "" && !isEncrypted(r.IMAPPassword)
		needsSMTP := r.SMTPPassword != "" && !isEncrypted(r.SMTPPassword)
		if !needsIMAP && !needsSMTP {
			continue
		}

		updates := map[string]any{}
		if needsIMAP {
			enc, err := Encrypt(r.IMAPPassword, key)
			if err != nil {
				logger.Error("encrypt IMAP password failed", slog.Int64("account_id", r.ID), slog.Any("error", err))
				continue
			}
			updates["imap_password"] = enc
		}
		if needsSMTP {
			enc, err := Encrypt(r.SMTPPassword, key)
			if err != nil {
				logger.Error("encrypt SMTP password failed", slog.Int64("account_id", r.ID), slog.Any("error", err))
				continue
			}
			updates["smtp_password"] = enc
		}
		if err := db.Table("accounts").Where("id = ?", r.ID).Updates(updates).Error; err != nil {
			logger.Error("update account failed", slog.Int64("account_id", r.ID), slog.Any("error", err))
			continue
		}
		migrated++
	}
	logger.Info("credential migration complete", slog.Int("migrated_accounts", migrated), slog.Int("total_accounts", len(rows)))
	return nil
}

func MigrateOIDCTokens(db *gorm.DB, key []byte, logger *slog.Logger) error {
	type row struct {
		ID               int64
		OIDCAccessToken  string
		OIDCRefreshToken string
	}

	var rows []row
	if err := db.Raw("SELECT id, oidc_access_token, oidc_refresh_token FROM users WHERE oidc_access_token != '' OR oidc_refresh_token != ''").Scan(&rows).Error; err != nil {
		return fmt.Errorf("read users: %w", err)
	}

	migrated := 0
	for _, r := range rows {
		needsAccess := r.OIDCAccessToken != "" && !isEncrypted(r.OIDCAccessToken)
		needsRefresh := r.OIDCRefreshToken != "" && !isEncrypted(r.OIDCRefreshToken)
		if !needsAccess && !needsRefresh {
			continue
		}

		updates := map[string]any{}
		if needsAccess {
			enc, err := Encrypt(r.OIDCAccessToken, key)
			if err != nil {
				logger.Error("encrypt OIDC access token failed", slog.Int64("user_id", r.ID), slog.Any("error", err))
				continue
			}
			updates["oidc_access_token"] = enc
		}
		if needsRefresh {
			enc, err := Encrypt(r.OIDCRefreshToken, key)
			if err != nil {
				logger.Error("encrypt OIDC refresh token failed", slog.Int64("user_id", r.ID), slog.Any("error", err))
				continue
			}
			updates["oidc_refresh_token"] = enc
		}
		if err := db.Table("users").Where("id = ?", r.ID).Updates(updates).Error; err != nil {
			logger.Error("update user failed", slog.Int64("user_id", r.ID), slog.Any("error", err))
			continue
		}
		migrated++
	}
	logger.Info("OIDC token migration complete", slog.Int("migrated_users", migrated), slog.Int("total_users", len(rows)))
	return nil
}
