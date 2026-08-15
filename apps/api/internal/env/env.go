package env

import (
	"fmt"

	troncenv "github.com/FacileStudio/tronc/env"

	"github.com/FacileStudio/porte"

	"github.com/FacileStudio/Courrier/apps/api/internal/crypto"
	"github.com/FacileStudio/Courrier/apps/api/internal/resourcetoken"
)

// OIDCConfig holds the values needed to federate this app with an OpenID
// Connect identity provider.
type OIDCConfig struct {
	Issuer       string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	SuccessURL   string
}

// Config is the fully-resolved application configuration: the core tronc
// settings plus the storage, SSO and encryption values this app reads.
type Config struct {
	troncenv.Core
	StorageDir          string
	OIDC                *OIDCConfig
	SSOOnly             bool
	ResourceTokenSecret []byte
	EncryptionKey       []byte
}

// Load reads the configuration from the environment, deriving the encryption
// and resource-token keys from ENCRYPTION_KEY and the OIDC block from the
// OIDC_* variables when SSO is configured.
func Load() (Config, error) {
	core, err := troncenv.LoadCore()
	if err != nil {
		return Config{}, err
	}
	if core.Port < 1 || core.Port > 65535 {
		return Config{}, fmt.Errorf("PORT must be a valid TCP port")
	}

	cfg := Config{
		Core:       core,
		StorageDir: troncenv.String("STORAGE_DIR", "./data"),
	}

	encryptionKey, err := troncenv.Required("ENCRYPTION_KEY")
	if err != nil {
		return Config{}, err
	}
	cfg.ResourceTokenSecret = resourcetoken.DeriveSecret(encryptionKey)
	cfg.EncryptionKey = crypto.DeriveKey(encryptionKey)

	if cfg.SSOOnly, err = troncenv.Bool("SSO_ONLY", false); err != nil {
		return Config{}, err
	}

	if issuer := troncenv.String("OIDC_ISSUER", ""); issuer != "" {
		clientID := troncenv.String("OIDC_CLIENT_ID", "")
		clientSecret := troncenv.String("OIDC_CLIENT_SECRET", "")
		redirectURL := troncenv.String("OIDC_REDIRECT_URL", "")
		if clientID == "" || clientSecret == "" || redirectURL == "" {
			return Config{}, fmt.Errorf("OIDC_CLIENT_ID, OIDC_CLIENT_SECRET, and OIDC_REDIRECT_URL are required when OIDC_ISSUER is set")
		}
		successURL := troncenv.String("OIDC_SUCCESS_URL", "")
		if successURL == "" && len(cfg.CORSAllowedOrigins) > 0 {
			successURL = cfg.CORSAllowedOrigins[0]
		}
		cfg.OIDC = &OIDCConfig{
			Issuer:       issuer,
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			SuccessURL:   successURL,
		}
	}

	return cfg, nil
}

// Porte is the one configuration porte's session manager, OIDC kit and local
// login are all built from.
//
// They share it rather than each reading the environment because the fields
// that decide cookie behaviour and the fields that decide the flow are the
// same fields, and porte refuses at boot a kit whose config disagrees with its
// manager's — a mismatch would otherwise change silently whether the session
// cookie is Secure.
//
// The session lifetime is porte's default, which is the thirty days this app
// already used. AcceptLegacyCookie is what keeps this deploy from signing
// everyone out:
// Courrier has always issued its session under the bare `session` name, and
// porte reads a __Host-prefixed one over https. Both are read; only the new
// one is written.
func (c Config) Porte() porte.Config {
	cfg := porte.Config{
		SSOOnly:            c.SSOOnly,
		AcceptLegacyCookie: true,
	}
	if c.OIDC == nil {
		return cfg
	}
	cfg.Issuer = c.OIDC.Issuer
	cfg.ClientID = c.OIDC.ClientID
	cfg.ClientSecret = c.OIDC.ClientSecret
	cfg.RedirectURL = c.OIDC.RedirectURL
	cfg.SuccessURL = c.OIDC.SuccessURL
	return cfg
}

// IssuerForMigration is the issuer the identity backfill keys on, or empty
// when SSO is not configured. It exists so the migration cannot be handed a
// placeholder: an identity row written under the wrong provider matches
// nothing and degrades to the email fallback in silence.
func (c Config) IssuerForMigration() string {
	if c.OIDC == nil {
		return ""
	}
	return c.OIDC.Issuer
}
