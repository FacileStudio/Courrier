// Package mcp is Courrier's remote MCP server: a curated set of mail tools on
// the go-sdk streamable HTTP transport in stateless mode, behind Registre
// access tokens bound to this server's resource URI.
//
// The transport, auth, Origin, rate-limit, pagination and confirmation parts
// do not know about mail, so they can move to tronc/mcp when a second app
// adopts them. tools*.go, pending*.go and execute.go are Courrier's own.
package mcp

import (
	"errors"
	"log/slog"
	"net/url"
	"time"

	"github.com/FacileStudio/porte/mcpauth"
	"github.com/modelcontextprotocol/go-sdk/auth"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"gorm.io/gorm"
)

// scopePrefix names this app's MCP scopes, courrier.read and courrier.write.
const scopePrefix = "courrier"

// Config describes the endpoint.
type Config struct {
	// Resource is the canonical URI people paste into their MCP client, such
	// as https://courrier.facile.studio/mcp. Registre binds tokens to it.
	Resource string

	// Issuer is Registre, the only authorization server.
	Issuer string

	// Version is the server version reported to clients.
	Version string

	// Origins are the browser origins allowed besides the resource's own.
	Origins []string

	// StateKey signs the confirmation state handed to clients.
	StateKey []byte

	// Verifier checks access tokens, normally mcpauth.Verifier.
	Verifier auth.TokenVerifier

	// PerMinute caps tool calls per person. Zero means 120.
	PerMinute int

	// Logger receives the write audit, which Journal collects.
	Logger *slog.Logger
}

// Service is the MCP endpoint and the pending-action store behind its
// confirmation links.
type Service struct {
	cfg         Config
	resourceURL *url.URL
	appURL      string
	metadataURL string
	scopes      mcpauth.Scopes
	tools       *toolset
	server      *sdk.Server
}

// NewService builds the MCP server over mailbox, keeping pending actions in
// db.
func NewService(db *gorm.DB, mailbox Mailbox, cfg Config) (*Service, error) {
	if cfg.Verifier == nil || len(cfg.StateKey) == 0 {
		return nil, errors.New("mcp: a token verifier and a state key are required")
	}
	resource, err := url.Parse(cfg.Resource)
	if err != nil || resource.Host == "" {
		return nil, errors.New("mcp: the resource must be an absolute URI")
	}
	metadataURL, err := mcpauth.MetadataURL(cfg.Resource)
	if err != nil {
		return nil, err
	}
	if cfg.PerMinute == 0 {
		cfg.PerMinute = 120
	}
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}
	appURL := resource.Scheme + "://" + resource.Host
	tools := &toolset{
		mail:    mailbox,
		pending: pendingStore{orm: db, ttl: 7 * 24 * time.Hour},
		confirm: confirmer{key: cfg.StateKey, ttl: 15 * time.Minute},
		appURL:  appURL,
		logger:  cfg.Logger,
	}
	s := &Service{cfg: cfg, resourceURL: resource, appURL: appURL, metadataURL: metadataURL, scopes: mcpauth.ScopesFor(scopePrefix), tools: tools}
	s.server = newServer(cfg.Version, tools)
	return s, nil
}
