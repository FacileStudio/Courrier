package mail

import (
	"net/http"
	"strconv"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/internal/middleware"
	"github.com/FacileStudio/Courrier/apps/api/modules/auth"
	"github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"

	"github.com/go-chi/chi/v5"
)

func registerTriageRoutes(router chi.Router, service *Service, authService *auth.Service) {
	router.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(authService))
		r.Get("/mail/triage", service.handleTriage)
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
	if raw := v.Get("account"); raw != "" {
		id, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return q, errors.Invalid("invalid account id")
		}
		q.AccountID = &id
	}
	if page, _ := strconv.Atoi(v.Get("page")); page > 1 {
		q.Page = page
	}
	if limit, _ := strconv.Atoi(v.Get("limit")); limit >= 1 && limit <= 100 {
		q.Limit = limit
	}
	return q, nil
}

func (s *Service) handleTriage(w http.ResponseWriter, req *http.Request) {
	uid, err := requestUserID(req)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
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
