package capsules

import (
	"net/http"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/internal/capsule"
	"github.com/FacileStudio/Courrier/apps/api/internal/middleware"
	"github.com/FacileStudio/tronc/httpjson"
	"github.com/go-chi/chi/v5"
)

// RegisterRoutes mounts POST /capsules behind the session.
func RegisterRoutes(router chi.Router, service *Service, authenticator middleware.Authenticator) {
	router.Route("/capsules", func(r chi.Router) {
		r.Use(middleware.RequireAuth(authenticator))
		r.With(middleware.RateLimit(30, time.Minute)).Post("/", service.handleSeal)
	})
}

func (s *Service) handleSeal(w http.ResponseWriter, req *http.Request) {
	req.Body = http.MaxBytesReader(w, req.Body, capsule.MaxCiphertext+4096)
	var body SealRequest
	if err := httpjson.DecodeJSON(w, req, &body); err != nil {
		httpjson.WriteError(w, err)
		return
	}
	out, err := s.Seal(req.Context(), body)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, out)
}
