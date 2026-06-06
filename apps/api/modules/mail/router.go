package mail

import (
	"net/http"
	"strconv"

	"api/internal/authcontext"
	"api/internal/httpjson"
	"api/internal/middleware"
	"api/modules/auth"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(router chi.Router, service *Service, authService *auth.Service) {
	router.Route("/accounts/{accountId}/mail", func(r chi.Router) {
		r.Use(middleware.RequireAuth(authService))

		r.Post("/sync", func(w http.ResponseWriter, req *http.Request) {
			identity := authcontext.MustIdentity(req.Context())
			uid, _ := strconv.ParseInt(identity.UserID, 10, 64)
			accountID, _ := strconv.ParseInt(chi.URLParam(req, "accountId"), 10, 64)

			if err := service.SyncAccount(req.Context(), uid, accountID); err != nil {
				httpjson.WriteError(w, err)
				return
			}
			httpjson.WriteJSON(w, http.StatusOK, map[string]bool{"synced": true})
		})

		r.Get("/folders", func(w http.ResponseWriter, req *http.Request) {
			identity := authcontext.MustIdentity(req.Context())
			uid, _ := strconv.ParseInt(identity.UserID, 10, 64)
			accountID, _ := strconv.ParseInt(chi.URLParam(req, "accountId"), 10, 64)

			folders, err := service.GetFolders(req.Context(), uid, accountID)
			if err != nil {
				httpjson.WriteError(w, err)
				return
			}
			resp := make([]FolderResponse, len(folders))
			for i, f := range folders {
				resp[i] = folderToResponse(f)
			}
			httpjson.WriteJSON(w, http.StatusOK, map[string]any{"folders": resp})
		})

		r.Post("/folders/{folderId}/sync", func(w http.ResponseWriter, req *http.Request) {
			identity := authcontext.MustIdentity(req.Context())
			uid, _ := strconv.ParseInt(identity.UserID, 10, 64)
			accountID, _ := strconv.ParseInt(chi.URLParam(req, "accountId"), 10, 64)
			folderID, _ := strconv.ParseInt(chi.URLParam(req, "folderId"), 10, 64)

			if err := service.SyncFolderEmails(req.Context(), uid, accountID, folderID); err != nil {
				httpjson.WriteError(w, err)
				return
			}
			httpjson.WriteJSON(w, http.StatusOK, map[string]bool{"synced": true})
		})

		r.Get("/folders/{folderType}/emails", func(w http.ResponseWriter, req *http.Request) {
			identity := authcontext.MustIdentity(req.Context())
			uid, _ := strconv.ParseInt(identity.UserID, 10, 64)
			accountID, _ := strconv.ParseInt(chi.URLParam(req, "accountId"), 10, 64)
			folderType := chi.URLParam(req, "folderType")

			page, _ := strconv.Atoi(req.URL.Query().Get("page"))
			if page < 1 {
				page = 1
			}
			limit, _ := strconv.Atoi(req.URL.Query().Get("limit"))
			if limit < 1 || limit > 100 {
				limit = 50
			}

			emails, total, err := service.GetEmailsByFolderType(req.Context(), uid, accountID, folderType, page, limit)
			if err != nil {
				httpjson.WriteError(w, err)
				return
			}

			resp := make([]EmailResponse, len(emails))
			for i, e := range emails {
				resp[i] = emailToResponse(e)
			}
			httpjson.WriteJSON(w, http.StatusOK, map[string]any{
				"emails": resp,
				"total":  total,
				"page":   page,
				"limit":  limit,
			})
		})

		r.Get("/emails/{emailId}", func(w http.ResponseWriter, req *http.Request) {
			identity := authcontext.MustIdentity(req.Context())
			uid, _ := strconv.ParseInt(identity.UserID, 10, 64)
			accountID, _ := strconv.ParseInt(chi.URLParam(req, "accountId"), 10, 64)
			emailID, _ := strconv.ParseInt(chi.URLParam(req, "emailId"), 10, 64)

			email, attachments, err := service.GetEmailWithAttachments(req.Context(), uid, accountID, emailID)
			if err != nil {
				httpjson.WriteError(w, err)
				return
			}
			httpjson.WriteJSON(w, http.StatusOK, emailToResponse(email, attachments...))
		})

		r.Get("/emails/{emailId}/attachments/{attachmentId}/download", func(w http.ResponseWriter, req *http.Request) {
			identity := authcontext.MustIdentity(req.Context())
			uid, _ := strconv.ParseInt(identity.UserID, 10, 64)
			accountID, _ := strconv.ParseInt(chi.URLParam(req, "accountId"), 10, 64)
			emailID, _ := strconv.ParseInt(chi.URLParam(req, "emailId"), 10, 64)
			attachmentID, _ := strconv.ParseInt(chi.URLParam(req, "attachmentId"), 10, 64)

			service.DownloadAttachment(w, req, uid, accountID, emailID, attachmentID)
		})

		r.Patch("/emails/{emailId}", func(w http.ResponseWriter, req *http.Request) {
			identity := authcontext.MustIdentity(req.Context())
			uid, _ := strconv.ParseInt(identity.UserID, 10, 64)
			accountID, _ := strconv.ParseInt(chi.URLParam(req, "accountId"), 10, 64)
			emailID, _ := strconv.ParseInt(chi.URLParam(req, "emailId"), 10, 64)

			var body UpdateEmailRequest
			if err := httpjson.DecodeJSON(w, req, &body); err != nil {
				httpjson.WriteError(w, err)
				return
			}

			email, err := service.UpdateEmail(req.Context(), uid, accountID, emailID, body)
			if err != nil {
				httpjson.WriteError(w, err)
				return
			}
			httpjson.WriteJSON(w, http.StatusOK, emailToResponse(email))
		})

		r.Post("/send", func(w http.ResponseWriter, req *http.Request) {
			identity := authcontext.MustIdentity(req.Context())
			uid, _ := strconv.ParseInt(identity.UserID, 10, 64)
			accountID, _ := strconv.ParseInt(chi.URLParam(req, "accountId"), 10, 64)

			var body SendRequest
			if err := httpjson.DecodeJSON(w, req, &body); err != nil {
				httpjson.WriteError(w, err)
				return
			}

			if err := service.Send(req.Context(), uid, accountID, body); err != nil {
				httpjson.WriteError(w, err)
				return
			}
			httpjson.WriteJSON(w, http.StatusOK, map[string]bool{"sent": true})
		})
	})

	router.Post("/mail/test-connection", func(w http.ResponseWriter, req *http.Request) {
		var body TestConnectionRequest
		if err := httpjson.DecodeJSON(w, req, &body); err != nil {
			httpjson.WriteError(w, err)
			return
		}
		if err := service.TestConnection(req.Context(), body); err != nil {
			httpjson.WriteError(w, err)
			return
		}
		httpjson.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
	})
}
