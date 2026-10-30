package mail

import (
	"net/http"
	"strconv"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/internal/authcontext"
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
	Scanned int `json:"scanned"`
}

func (s *Service) serveContacts(w http.ResponseWriter, req *http.Request, userID int64, accountID *int64) {
	page, _ := strconv.Atoi(req.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(req.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 50
	}
	rows, total, err := s.ListContacts(req.Context(), userID, accountID, page, limit)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, ContactsResponse{Contacts: rows, Total: total, Page: page, Limit: limit})
}

func requestUserID(req *http.Request) (int64, error) {
	identity := authcontext.MustIdentity(req.Context())
	uid, err := strconv.ParseInt(identity.UserID, 10, 64)
	if err != nil {
		return 0, errors.Invalid("invalid user id")
	}
	return uid, nil
}

func registerContactRoutes(router chi.Router, service *Service, authService *auth.Service) {
	router.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(authService))
		r.Get("/mail/contacts", service.handleAllContacts)
		r.With(middleware.RateLimit(2, time.Minute)).Post("/accounts/{accountId}/mail/contacts/backfill", service.handleBackfillContacts)
	})
}

func (s *Service) handleAllContacts(w http.ResponseWriter, req *http.Request) {
	uid, err := requestUserID(req)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	s.serveContacts(w, req, uid, nil)
}

func (s *Service) handleBackfillContacts(w http.ResponseWriter, req *http.Request) {
	uid, err := requestUserID(req)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	accountID, err := strconv.ParseInt(chi.URLParam(req, "accountId"), 10, 64)
	if err != nil {
		httpjson.WriteError(w, errors.Invalid("invalid account id"))
		return
	}
	scanned, err := s.deepSyncSent(req.Context(), uid, accountID, maxDeepSyncSent)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, BackfillContactsResponse{Scanned: scanned})
}
