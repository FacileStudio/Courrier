package settings

import (
	"net/http"

	"github.com/FacileStudio/Courrier/apps/api/internal/middleware"
	"github.com/FacileStudio/Courrier/apps/api/modules/auth"
	"github.com/FacileStudio/tronc/httpjson"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes mounts the authenticated /settings endpoint on the router.
func RegisterRoutes(router chi.Router, service *Service, authService *auth.Service) {
	router.Route("/settings", func(r chi.Router) {
		r.With(middleware.RequireAuth(authService)).Get("/", func(w http.ResponseWriter, req *http.Request) {
			httpjson.WriteJSON(w, http.StatusOK, service.controller.getSettings())
		})
	})
}
