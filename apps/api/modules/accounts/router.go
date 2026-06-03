package accounts

import (
	"net/http"
	"strconv"

	"api/internal/authcontext"
	"api/internal/httpjson"
	"api/internal/middleware"
	"api/modules/auth"
	"api/schemas"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(router chi.Router, service *Service, authService *auth.Service) {
	router.Route("/accounts", func(r chi.Router) {
		r.Use(middleware.RequireAuth(authService))

		r.Post("/", func(w http.ResponseWriter, req *http.Request) {
			identity := authcontext.MustIdentity(req.Context())
			uid, _ := strconv.ParseInt(identity.UserID, 10, 64)
			var body CreateAccountRequest
			if err := httpjson.DecodeJSON(w, req, &body); err != nil {
				httpjson.WriteError(w, err)
				return
			}
			account, err := service.Create(req.Context(), uid, body)
			if err != nil {
				httpjson.WriteError(w, err)
				return
			}
			httpjson.WriteJSON(w, http.StatusCreated, toResponse(account))
		})

		r.Get("/", func(w http.ResponseWriter, req *http.Request) {
			identity := authcontext.MustIdentity(req.Context())
			uid, _ := strconv.ParseInt(identity.UserID, 10, 64)
			accounts, err := service.List(req.Context(), uid)
			if err != nil {
				httpjson.WriteError(w, err)
				return
			}
			resp := make([]AccountResponse, len(accounts))
			for i, a := range accounts {
				resp[i] = toResponse(a)
			}
			httpjson.WriteJSON(w, http.StatusOK, map[string]any{"accounts": resp})
		})

		r.Get("/{id}", func(w http.ResponseWriter, req *http.Request) {
			identity := authcontext.MustIdentity(req.Context())
			uid, _ := strconv.ParseInt(identity.UserID, 10, 64)
			id, _ := strconv.ParseInt(chi.URLParam(req, "id"), 10, 64)
			account, err := service.Get(req.Context(), uid, id)
			if err != nil {
				httpjson.WriteError(w, err)
				return
			}
			httpjson.WriteJSON(w, http.StatusOK, toResponse(account))
		})

		r.Put("/{id}", func(w http.ResponseWriter, req *http.Request) {
			identity := authcontext.MustIdentity(req.Context())
			uid, _ := strconv.ParseInt(identity.UserID, 10, 64)
			id, _ := strconv.ParseInt(chi.URLParam(req, "id"), 10, 64)
			var body UpdateAccountRequest
			if err := httpjson.DecodeJSON(w, req, &body); err != nil {
				httpjson.WriteError(w, err)
				return
			}
			account, err := service.Update(req.Context(), uid, id, body)
			if err != nil {
				httpjson.WriteError(w, err)
				return
			}
			httpjson.WriteJSON(w, http.StatusOK, toResponse(account))
		})

		r.Delete("/{id}", func(w http.ResponseWriter, req *http.Request) {
			identity := authcontext.MustIdentity(req.Context())
			uid, _ := strconv.ParseInt(identity.UserID, 10, 64)
			id, _ := strconv.ParseInt(chi.URLParam(req, "id"), 10, 64)
			if err := service.Delete(req.Context(), uid, id); err != nil {
				httpjson.WriteError(w, err)
				return
			}
			httpjson.WriteJSON(w, http.StatusOK, map[string]bool{"deleted": true})
		})
	})
}

func toResponse(a schemas.Account) AccountResponse {
	return AccountResponse{
		ID:        a.ID,
		Name:      a.Name,
		Email:     a.Email,
		IMAPHost:  a.IMAPHost,
		IMAPPort:  a.IMAPPort,
		IMAPUser:  a.IMAPUser,
		SMTPHost:  a.SMTPHost,
		SMTPPort:  a.SMTPPort,
		SMTPUser:  a.SMTPUser,
		Signature: a.Signature,
		IsDefault: a.IsDefault,
		CreatedAt: a.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: a.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
