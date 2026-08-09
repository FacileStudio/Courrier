package auth

import (
	"net/http"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/internal/authcontext"
	"github.com/FacileStudio/Courrier/apps/api/internal/env"
	"github.com/FacileStudio/Courrier/apps/api/internal/middleware"
	"github.com/FacileStudio/Courrier/apps/api/internal/resourcetoken"
	"github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes mounts what porte does not own.
//
// /auth/config, /auth/logout, /auth/sync-profile and the whole OIDC flow come
// from porte's session manager and its OIDC kit, mounted in main.go. What is
// left here is the local password path, which keeps Courrier's own
// {user_id, token} response shape, and the resource token this app signs for
// its own attachment routes.
//
// Under SSO_ONLY the credential routes are not registered rather than
// rejected, so there is no endpoint left to probe for an account. That is the
// behaviour this app already had.
func RegisterRoutes(router chi.Router, service *Service, appEnv env.Config) {
	router.Route("/auth", func(router chi.Router) {
		if !appEnv.SSOOnly {
			router.With(middleware.RateLimit(3, time.Minute)).Post("/register", func(w http.ResponseWriter, request *http.Request) {
				var req RegisterRequest
				if err := httpjson.DecodeJSON(w, request, &req); err != nil {
					httpjson.WriteError(w, err)
					return
				}
				resp, err := service.controller.register(w, request, &req)
				if err != nil {
					httpjson.WriteError(w, err)
					return
				}
				httpjson.WriteJSON(w, http.StatusCreated, resp)
			})

			router.With(middleware.RateLimit(10, time.Minute)).Post("/login", func(w http.ResponseWriter, request *http.Request) {
				var req LoginRequest
				if err := httpjson.DecodeJSON(w, request, &req); err != nil {
					httpjson.WriteError(w, err)
					return
				}
				resp, err := service.controller.login(w, request, &req)
				if err != nil {
					httpjson.WriteError(w, err)
					return
				}
				httpjson.WriteJSON(w, http.StatusOK, resp)
			})
		}

		if len(appEnv.ResourceTokenSecret) > 0 {
			router.With(middleware.RequireAuth(service)).Get("/resource-token", func(w http.ResponseWriter, r *http.Request) {
				identity := authcontext.MustIdentity(r.Context())
				token := resourcetoken.Sign(appEnv.ResourceTokenSecret, identity.UserID, 5*time.Minute)
				httpjson.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
			})
		} else {
			router.With(middleware.RequireAuth(service)).Get("/resource-token", func(w http.ResponseWriter, r *http.Request) {
				httpjson.WriteError(w, errors.Internal("resource tokens not configured (ENCRYPTION_KEY missing)", nil))
			})
		}
	})
}
