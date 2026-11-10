package auth

import (
	"net/http"
	"strconv"
	"time"

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
// from porte's session manager and its OIDC kit, mounted in main.go. Courrier
// has no register or login route: sign-in goes through Registre only. What is
// left here is the resource token this app signs for its own attachment routes.
func RegisterRoutes(router chi.Router, service *Service, appEnv env.Config) {
	router.Route("/auth", func(router chi.Router) {
		if len(appEnv.ResourceTokenSecret) > 0 {
			router.With(middleware.RequireAuth(service)).Get("/resource-token", func(w http.ResponseWriter, r *http.Request) {
				userID := strconv.FormatInt(middleware.UserID(r), 10)
				token := resourcetoken.Sign(appEnv.ResourceTokenSecret, userID, 5*time.Minute)
				httpjson.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
			})
		} else {
			router.With(middleware.RequireAuth(service)).Get("/resource-token", func(w http.ResponseWriter, r *http.Request) {
				httpjson.WriteError(w, errors.Internal("resource tokens not configured (ENCRYPTION_KEY missing)", nil))
			})
		}
	})
}
