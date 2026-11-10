package assist

import (
	"net/http"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/internal/middleware"
	"github.com/FacileStudio/tronc/httpjson"
	"github.com/go-chi/chi/v5"
)

// RegisterRoutes mounts GET and POST /assist behind the session.
func RegisterRoutes(router chi.Router, service *Service, authenticator middleware.Authenticator) {
	router.Route("/assist", func(r chi.Router) {
		r.Use(middleware.RequireAuth(authenticator))
		r.Get("/", service.handleStatus)
		r.With(middleware.RateLimit(30, time.Minute)).Post("/", service.handleRun)
	})
}

func (s *Service) handleStatus(w http.ResponseWriter, req *http.Request) {
	uid := middleware.UserID(req)
	httpjson.WriteJSON(w, http.StatusOK, s.Status(req.Context(), uid))
}

func (s *Service) handleRun(w http.ResponseWriter, req *http.Request) {
	uid := middleware.UserID(req)
	req.Body = http.MaxBytesReader(w, req.Body, 256<<10)
	var body Request
	if err := httpjson.DecodeJSON(w, req, &body); err != nil {
		httpjson.WriteError(w, err)
		return
	}
	out, err := s.Run(req.Context(), uid, body)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, out)
}
