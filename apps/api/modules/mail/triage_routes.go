package mail

import (
	"net/http"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/internal/middleware"
	"github.com/FacileStudio/Courrier/apps/api/modules/auth"
	"github.com/FacileStudio/tronc/httpjson"

	"github.com/go-chi/chi/v5"
)

func registerTriageRoutes(router chi.Router, service *Service, authService *auth.Service) {
	router.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(authService))
		r.With(middleware.RateLimit(60, time.Minute)).Get("/mail/triage", service.handleTriage)
		r.With(middleware.RateLimit(2, time.Minute)).Post("/accounts/{accountId}/mail/backfill", service.handleBackfill)
	})
}

func parseTriageQuery(req *http.Request) (TriageQuery, error) {
	v := req.URL.Query()
	q := TriageQuery{UnreadOnly: v.Get("unread") == "true", Page: 1, Limit: 50}
	var err error
	if q.Classes, err = parseClasses(v.Get("class")); err != nil {
		return q, err
	}
	if q.Since, err = parseSince(v.Get("since"), time.Now()); err != nil {
		return q, err
	}
	if q.AccountID, err = optionalAccountID(req); err != nil {
		return q, err
	}
	q.Page, q.Limit = parsePaging(v)
	return q, nil
}

func (s *Service) handleTriage(w http.ResponseWriter, req *http.Request) {
	uid := middleware.UserID(req)
	q, err := parseTriageQuery(req)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	resp, err := s.Triage(req.Context(), uid, q)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, resp)
}
