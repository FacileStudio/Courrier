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

// ContactsResponse is the paginated list of addresses the user has written to.
type ContactsResponse struct {
	Contacts []ContactRow `json:"contacts"`
	Total    int64        `json:"total"`
	Page     int          `json:"page"`
	Limit    int          `json:"limit"`
}

// BackfillContactsResponse reports how many sent envelopes a deep sync read.
type BackfillContactsResponse struct {
	Scanned int  `json:"scanned"`
	Running bool `json:"running,omitempty"`
}

func (s *Service) serveContacts(w http.ResponseWriter, req *http.Request, userID int64, accountID *int64) {
	page, limit := parsePaging(req.URL.Query())
	rows, total, err := s.ListContacts(req.Context(), userID, accountID, page, limit)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, ContactsResponse{Contacts: rows, Total: total, Page: page, Limit: limit})
}

func registerContactRoutes(router chi.Router, service *Service, authService *auth.Service) {
	router.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(authService))
		r.With(middleware.RateLimit(60, time.Minute)).Get("/mail/contacts", service.handleAllContacts)
		r.With(middleware.RateLimit(2, time.Minute)).Post("/accounts/{accountId}/mail/contacts/backfill", service.handleBackfillContacts)
	})
}

func (s *Service) handleAllContacts(w http.ResponseWriter, req *http.Request) {
	uid := middleware.UserID(req)
	accountID, err := optionalAccountID(req)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	s.serveContacts(w, req, uid, accountID)
}

func optionalAccountID(req *http.Request) (*int64, error) {
	raw := req.URL.Query().Get("account")
	if raw == "" {
		return nil, nil
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return nil, errors.Invalid("invalid account id")
	}
	return &id, nil
}

func (s *Service) handleBackfillContacts(w http.ResponseWriter, req *http.Request) {
	uid := middleware.UserID(req)
	accountID, err := strconv.ParseInt(chi.URLParam(req, "accountId"), 10, 64)
	if err != nil {
		httpjson.WriteError(w, errors.Invalid("invalid account id"))
		return
	}
	res, err := s.BackfillContacts(req.Context(), uid, accountID)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, res)
}
