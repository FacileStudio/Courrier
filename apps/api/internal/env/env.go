package env

import (
	"crypto/sha256"
	"fmt"
	"slices"
	"strings"
	"time"

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

	// CLIAudience is the CLI's own client id, and it is what mounts
	// POST /auth/oidc/device/exchange. Empty leaves the route unmounted, so
	// an app that has not been told which CLI to trust has no endpoint to
	// probe rather than one that refuses.
	//
	// porte reads no environment of its own, so a variable set in the panel
	// and never read here reaches nothing.
	CLIAudience string

	// MachineAudience is this app's own Registre client id. It mounts offline
	// verification of the delegated bearer JWTs other suite services mint for
	// this app.
	MachineAudience string

	// ScopePrefix gates those JWTs by method: GET needs courrier.read or
	// courrier.write, anything else needs courrier.write. It requires
	// MachineAudience.
	ScopePrefix string

	// MobileRedirect is where porte's ?flow=mobile login hands the one-time
	// code back to the Android app: its custom scheme, never a value taken
	// from the request.
	MobileRedirect string
}

// Push holds the notification channels. Each one is off while its values
// are empty, and the mail watcher only runs while a channel is on.
type Push struct {
	FCMCredentials string
	VAPIDPublic    string
	VAPIDPrivate   string
	VAPIDSubject   string
	PollInterval   time.Duration
}

// DefaultMobileRedirect is the deep link the Courrier Android app registers.
const DefaultMobileRedirect = "studio.facile.courrier://auth"

// Config is the fully-resolved application configuration: the core tronc
// settings plus the storage, SSO and encryption values this app reads.
type Config struct {
	troncenv.Core
	StorageDir          string
	OIDC                *OIDCConfig
	SSOOnly             bool
	ResourceTokenSecret []byte
	EncryptionKey       []byte
	Push                Push
	BulleURL            string
	MCP                 MCP
	Admins              []string
}

// MCP turns on the remote MCP server. Resource is the canonical URI clients
// connect to, such as https://courrier.facile.studio/mcp, and the RFC 8707
// resource Registre binds its tokens to. It needs SSO. StateKey signs the
// confirmation state handed to clients.
type MCP struct {
	Resource string
	StateKey []byte
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

	cfg.BulleURL = troncenv.String("BULLE_URL", "")
	cfg.Admins = troncenv.List("COURRIER_ADMINS")
	cfg.MCP = MCP{Resource: troncenv.String("MCP_RESOURCE", ""), StateKey: deriveStateKey(encryptionKey)}
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
			CLIAudience:  troncenv.String("OIDC_CLI_AUDIENCE", ""),

			MachineAudience: troncenv.String("OIDC_MACHINE_AUDIENCE", ""),
			ScopePrefix:     troncenv.String("OIDC_SCOPE_PREFIX", ""),
			MobileRedirect:  troncenv.String("OIDC_MOBILE_REDIRECT", DefaultMobileRedirect),
		}
		cfg.CORSAllowedOrigins = withAppOrigins(cfg.CORSAllowedOrigins, cfg.OIDC.MobileRedirect)
	}

	if cfg.IsProduction() && cfg.OIDC == nil {
		return Config{}, fmt.Errorf("OIDC_ISSUER is required in production: sign-in goes through Registre only")
	}

	if cfg.Push, err = loadPush(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// deriveStateKey derives the MCP confirmation key from the app's encryption
// passphrase, namespaced so it collides with no other key.
func deriveStateKey(encryptionKey string) []byte {
	sum := sha256.Sum256([]byte("courrier-mcp-state:" + encryptionKey))
	return sum[:]
}

// withAppOrigins adds the origins the installed apps serve their pages from to
// the CORS list when the app login is on: the Capacitor Android app, then the
// Tauri desktop app on macOS and Linux, then on Windows. A deploy that
// overrides CORS_ALLOWED_ORIGINS therefore cannot lock the apps out. It runs
// after OIDC_SUCCESS_URL has defaulted to the first origin, so an app origin
// never becomes the web success page.
func withAppOrigins(origins []string, mobileRedirect string) []string {
	if mobileRedirect == "" {
		return origins
	}
	for _, origin := range []string{"https://localhost", "tauri://localhost", "http://tauri.localhost"} {
		if !slices.Contains(origins, origin) {
			origins = append(origins, origin)
		}
	}
	return origins
}

// loadPush reads the FCM service account, the VAPID key pair and the mail
// watcher interval. MAIL_POLL_INTERVAL=0 turns the watcher off.
func loadPush() (Push, error) {
	interval, err := troncenv.Duration("MAIL_POLL_INTERVAL", 2*time.Minute)
	if err != nil {
		return Push{}, err
	}
	if interval < 0 {
		return Push{}, fmt.Errorf("MAIL_POLL_INTERVAL must not be negative")
	}
	return Push{
		FCMCredentials: troncenv.String("FCM_SERVICE_ACCOUNT", ""),
		VAPIDPublic:    strings.TrimSpace(troncenv.String("VAPID_PUBLIC_KEY", "")),
		VAPIDPrivate:   strings.TrimSpace(troncenv.String("VAPID_PRIVATE_KEY", "")),
		VAPIDSubject:   troncenv.String("VAPID_SUBJECT", "mailto:contact@facile.studio"),
		PollInterval:   interval,
	}, nil
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
	cfg.CLIAudience = c.OIDC.CLIAudience
	cfg.MachineAudience = c.OIDC.MachineAudience
	cfg.ScopePrefix = c.OIDC.ScopePrefix
	cfg.MobileRedirect = c.OIDC.MobileRedirect
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
