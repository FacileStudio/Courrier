package main

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FacileStudio/Courrier/apps/api/internal/env"
	"github.com/FacileStudio/Courrier/apps/api/modules/mail"
	"github.com/FacileStudio/Courrier/apps/api/modules/spaces"
	"github.com/FacileStudio/porte/directory"
	"github.com/FacileStudio/porte/oidc"
	portepg "github.com/FacileStudio/porte/pg"
	"github.com/FacileStudio/porte/session"
	"github.com/FacileStudio/tronc/apiref"
	"github.com/go-chi/chi/v5"
)

// referenceRouter is buildRouter as run() wires it with SSO on, so the
// Registre nudge route is mounted and has to be documented too. The OIDC
// block is set after the kit is built, which keeps discovery off the network.
func referenceRouter(t *testing.T) chi.Router {
	t.Helper()
	appEnv := env.Config{StorageDir: t.TempDir()}
	appLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	store := portepg.New(nil)
	sessions, err := session.New(appEnv.Porte(), session.Deps{Sessions: store.Sessions(), Logger: appLogger})
	if err != nil {
		t.Fatalf("session.New: %v", err)
	}
	kit, err := oidc.New(context.Background(), appEnv.Porte(), oidc.Deps{Sessions: sessions, Logger: appLogger})
	if err != nil {
		t.Fatalf("oidc.New: %v", err)
	}
	appEnv.OIDC = &env.OIDCConfig{Issuer: "https://sso.test"}
	svc := services{
		sessions: sessions, kit: kit, mail: mail.NewService(nil, nil),
		directory: spaces.NewDirectory(nil, directory.Client{Issuer: appEnv.OIDC.Issuer}, appLogger),
		deletions: spaces.NewDeletions(nil, directory.Client{Issuer: appEnv.OIDC.Issuer}, nil, appEnv.StorageDir, appLogger),
	}
	return buildRouter(nil, func(context.Context) error { return nil }, appEnv, appLogger, svc)
}

func TestEveryRouteIsDocumented(t *testing.T) {
	if missing := apiref.Undocumented(referenceRouter(t), referenceConfig(), "/version"); len(missing) > 0 {
		t.Errorf("routes missing from the API registry: %v", missing)
	}
}

func TestRegistryIsComplete(t *testing.T) {
	if issues := apiref.Incomplete(
		referenceConfig(),
		"/auth/logout",
		"/auth/oidc",
		"/auth/oidc/callback",
		"/auth/sync-profile",
		"/auth/backchannel-logout",
		"/accounts/{accountId}/mail/check",
		"/accounts/{accountId}/mail/sync",
		"/accounts/{accountId}/mail/backfill",
		"/accounts/{accountId}/mail/contacts/backfill",
		"/accounts/{accountId}/mail/folders/{folderId}/sync",
		"/accounts/{accountId}/mail/emails/{emailId}/attachments/{attachmentId}/download",
		"/accounts/{accountId}/mail/emails/{emailId}/cid/{cid}",
		"/users/me/avatar",
		"/spaces/{spaceId}/leave",
		"/registre/sync",
		"/registre/deletions",
	); len(issues) > 0 {
		t.Errorf("incomplete documentation routes: %v", issues)
	}
}

func TestReferenceIsServedAtDocs(t *testing.T) {
	router := referenceRouter(t)

	page := httptest.NewRecorder()
	router.ServeHTTP(page, httptest.NewRequest(http.MethodGet, "/docs", nil))
	if page.Code != http.StatusOK {
		t.Fatalf("GET /docs = %d, want 200", page.Code)
	}

	spec := httptest.NewRecorder()
	router.ServeHTTP(spec, httptest.NewRequest(http.MethodGet, "/docs/openapi.json", nil))
	if spec.Code != http.StatusOK {
		t.Fatalf("GET /docs/openapi.json = %d, want 200", spec.Code)
	}
	var document struct {
		OpenAPI string         `json:"openapi"`
		Paths   map[string]any `json:"paths"`
	}
	if err := json.Unmarshal(spec.Body.Bytes(), &document); err != nil {
		t.Fatalf("spec is not JSON: %v", err)
	}
	if document.OpenAPI == "" || len(document.Paths) == 0 {
		t.Fatalf("spec is empty: %+v", document)
	}
}
