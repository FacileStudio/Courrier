package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/FacileStudio/tronc/health"

	"github.com/FacileStudio/Courrier/apps/api/internal/env"
	"github.com/FacileStudio/Courrier/apps/api/modules/accounts"
	"github.com/FacileStudio/Courrier/apps/api/modules/auth"
	"github.com/FacileStudio/Courrier/apps/api/modules/mail"
	"github.com/FacileStudio/Courrier/apps/api/modules/settings"
	"github.com/FacileStudio/Courrier/apps/api/modules/spaces"
	"github.com/FacileStudio/Courrier/apps/api/modules/users"
)

const (
	spaFallbackBody  = "<html>spa</html>"
	filesHandlerBody = "static-file"
)

// testRouter wires the modules exactly as run() does. The services are nil
// because route registration never dereferences them; only handlers do, and
// these cases stop at routing.
func testRouter() *chi.Mux {
	router := chi.NewRouter()

	health.Mount(router)

	router.Route("/api", func(r chi.Router) {
		r.Handle("/files/*", http.StripPrefix("/api/files/", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(filesHandlerBody))
		})))

		auth.RegisterRoutes(r, nil, env.Config{})
		accounts.RegisterRoutes(r, nil, nil)
		mail.RegisterRoutes(r, nil, nil, nil)
		users.RegisterRoutes(r, nil, nil)
		settings.RegisterRoutes(r, nil, nil)
		spaces.RegisterRoutes(r, nil, nil)
	})

	router.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(spaFallbackBody))
	}))

	return router
}

func TestUnknownAPIRoutesDoNotFallThroughToTheSPA(t *testing.T) {
	router := testRouter()

	for _, path := range []string{"/api/nimporte-quoi", "/api", "/api/auth/nope", "/api/users/me/nope"} {
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))

		if recorder.Code != http.StatusNotFound {
			t.Errorf("GET %s = %d, want 404", path, recorder.Code)
		}
		if recorder.Body.String() == spaFallbackBody {
			t.Errorf("GET %s was answered by the SPA catch-all", path)
		}
	}
}

func TestPublicAPIURLsAreUnchanged(t *testing.T) {
	router := testRouter()

	registered := map[string]bool{}
	if err := chi.Walk(router, func(method, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		registered[method+" "+route] = true
		return nil
	}); err != nil {
		t.Fatalf("walk: %v", err)
	}

	want := []string{
		"GET /api/auth/config",
		"POST /api/auth/login",
		"POST /api/auth/register",
		"POST /api/auth/logout",
		"GET /api/auth/resource-token",
		"GET /api/accounts/",
		"POST /api/accounts/",
		"GET /api/accounts/{id}",
		"POST /api/accounts/{accountId}/mail/sync",
		"GET /api/accounts/{accountId}/mail/folders",
		"GET /api/accounts/{accountId}/mail/emails/{emailId}",
		"GET /api/accounts/{accountId}/mail/emails/{emailId}/cid/{cid}",
		"POST /api/mail/test-connection",
		"GET /api/templates/",
		"GET /api/users/me",
		"GET /api/settings/",
		"GET /api/spaces/",
		"GET /api/spaces/{spaceId}/members/",
	}
	for _, route := range want {
		if !registered[route] {
			t.Errorf("%s is not registered", route)
		}
	}
}

func TestHealthRoutesSurviveTheAPISubtree(t *testing.T) {
	router := testRouter()

	for _, path := range []string{"/health", "/ready", "/api/health", "/api/ready"} {
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))

		if recorder.Code != http.StatusOK {
			t.Errorf("GET %s = %d, want 200", path, recorder.Code)
		}
	}
}

func TestFilesRouteStillServesFromStorage(t *testing.T) {
	recorder := httptest.NewRecorder()
	testRouter().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/files/avatars/1.png", nil))

	if recorder.Body.String() != filesHandlerBody {
		t.Fatalf("GET /api/files/avatars/1.png body = %q, want the file handler", recorder.Body.String())
	}
}

func TestSPAFallbackStillServesNonAPIPaths(t *testing.T) {
	recorder := httptest.NewRecorder()
	testRouter().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/mail/inbox", nil))

	if recorder.Body.String() != spaFallbackBody {
		t.Fatalf("GET /mail/inbox body = %q, want the SPA fallback", recorder.Body.String())
	}
}
