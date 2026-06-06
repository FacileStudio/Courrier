package mail

import (
	"io"
	"mime"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"api/internal/authcontext"
	"api/internal/errors"
	"api/internal/httpjson"
	"api/internal/middleware"
	"api/internal/resourcetoken"
	"api/modules/auth"

	"github.com/go-chi/chi/v5"
)

func parseSendMultipart(req *http.Request) (SendRequest, error) {
	if err := req.ParseMultipartForm(25 << 20); err != nil {
		return SendRequest{}, errors.TooLarge("request too large or invalid multipart")
	}

	var sr SendRequest
	if v := req.FormValue("to"); v != "" {
		sr.To = strings.Split(v, ",")
		for i := range sr.To {
			sr.To[i] = strings.TrimSpace(sr.To[i])
		}
	}
	if v := req.FormValue("cc"); v != "" {
		sr.Cc = strings.Split(v, ",")
		for i := range sr.Cc {
			sr.Cc[i] = strings.TrimSpace(sr.Cc[i])
		}
	}
	sr.Subject = req.FormValue("subject")
	sr.Body = req.FormValue("body")
	sr.BodyHTML = req.FormValue("body_html")
	sr.InReplyTo = req.FormValue("in_reply_to")
	if v := req.FormValue("references"); v != "" {
		sr.References = strings.Split(v, ",")
		for i := range sr.References {
			sr.References[i] = strings.TrimSpace(sr.References[i])
		}
	}

	if req.MultipartForm != nil && req.MultipartForm.File != nil {
		for _, headers := range req.MultipartForm.File["attachments"] {
			f, err := headers.Open()
			if err != nil {
				continue
			}
			tmpFile, err := os.CreateTemp("", "courrier-attachment-*")
			if err != nil {
				f.Close()
				continue
			}
			_, err = io.Copy(tmpFile, f)
			tmpFile.Close()
			f.Close()
			if err != nil {
				os.Remove(tmpFile.Name())
				continue
			}
			mimeType := headers.Header.Get("Content-Type")
			if mimeType == "" {
				mimeType = mime.TypeByExtension("." + extensionFromFilename(headers.Filename))
				if mimeType == "" {
					mimeType = "application/octet-stream"
				}
			}
			sr.Attachments = append(sr.Attachments, AttachmentUpload{
				Filename: headers.Filename,
				MimeType: mimeType,
				FilePath: tmpFile.Name(),
			})
		}
	}

	return sr, nil
}

func extensionFromFilename(name string) string {
	idx := strings.LastIndex(name, ".")
	if idx < 0 {
		return ""
	}
	return name[idx+1:]
}

func RegisterRoutes(router chi.Router, service *Service, authService *auth.Service, rtSecret []byte) {
	router.Get("/accounts/{accountId}/mail/emails/{emailId}/cid/{cid}", func(w http.ResponseWriter, req *http.Request) {
		var userID string

		if token := req.URL.Query().Get("token"); token != "" && len(rtSecret) > 0 {
			uid, err := resourcetoken.Verify(rtSecret, token)
			if err == nil {
				userID = uid
			}
		}

		if userID == "" {
			authHeader := req.Header.Get("Authorization")
			if authHeader == "" {
				token := req.URL.Query().Get("token")
				if token == "" {
					httpjson.WriteError(w, errors.Unauthorized("missing token"))
					return
				}
				authHeader = "Bearer " + token
			}
			uid, _, err := authService.Authenticate(req.Context(), authHeader)
			if err != nil {
				httpjson.WriteError(w, err)
				return
			}
			userID = uid
		}

		uid, _ := strconv.ParseInt(userID, 10, 64)
		accountID, _ := strconv.ParseInt(chi.URLParam(req, "accountId"), 10, 64)
		emailID, _ := strconv.ParseInt(chi.URLParam(req, "emailId"), 10, 64)
		cid := chi.URLParam(req, "cid")

		service.ServeCIDImage(w, req, uid, accountID, emailID, cid)
	})

	router.Route("/accounts/{accountId}/mail", func(r chi.Router) {
		r.Use(middleware.RequireAuth(authService))

		r.With(middleware.RateLimit(5, time.Minute)).Post("/sync", func(w http.ResponseWriter, req *http.Request) {
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

		r.With(middleware.RateLimit(10, time.Minute)).Post("/folders/{folderId}/sync", func(w http.ResponseWriter, req *http.Request) {
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

		r.Get("/contacts", func(w http.ResponseWriter, req *http.Request) {
			identity := authcontext.MustIdentity(req.Context())
			uid, _ := strconv.ParseInt(identity.UserID, 10, 64)
			accountID, _ := strconv.ParseInt(chi.URLParam(req, "accountId"), 10, 64)
			q := req.URL.Query().Get("q")
			if q == "" {
				httpjson.WriteJSON(w, http.StatusOK, map[string]any{"contacts": []ContactResult{}})
				return
			}
			contacts, err := service.SearchContacts(req.Context(), uid, accountID, q)
			if err != nil {
				httpjson.WriteError(w, err)
				return
			}
			httpjson.WriteJSON(w, http.StatusOK, map[string]any{"contacts": contacts})
		})

		r.With(middleware.RateLimit(10, time.Minute)).Post("/send", func(w http.ResponseWriter, req *http.Request) {
			identity := authcontext.MustIdentity(req.Context())
			uid, _ := strconv.ParseInt(identity.UserID, 10, 64)
			accountID, _ := strconv.ParseInt(chi.URLParam(req, "accountId"), 10, 64)

			contentType := req.Header.Get("Content-Type")
			var body SendRequest

			if strings.HasPrefix(contentType, "multipart/form-data") {
				req.Body = http.MaxBytesReader(w, req.Body, 25<<20)
				parsed, err := parseSendMultipart(req)
				if err != nil {
					for i := range parsed.Attachments {
						parsed.Attachments[i].Cleanup()
					}
					httpjson.WriteError(w, err)
					return
				}
				body = parsed
			} else {
				if err := httpjson.DecodeJSON(w, req, &body); err != nil {
					httpjson.WriteError(w, err)
					return
				}
			}

			defer func() {
				for i := range body.Attachments {
					body.Attachments[i].Cleanup()
				}
			}()

			if err := service.Send(req.Context(), uid, accountID, body); err != nil {
				httpjson.WriteError(w, err)
				return
			}
			httpjson.WriteJSON(w, http.StatusOK, map[string]bool{"sent": true})
		})

		r.Post("/drafts", func(w http.ResponseWriter, req *http.Request) {
			identity := authcontext.MustIdentity(req.Context())
			uid, _ := strconv.ParseInt(identity.UserID, 10, 64)
			accountID, _ := strconv.ParseInt(chi.URLParam(req, "accountId"), 10, 64)

			var body SendRequest
			if err := httpjson.DecodeJSON(w, req, &body); err != nil {
				httpjson.WriteError(w, err)
				return
			}

			email, err := service.SaveDraft(req.Context(), uid, accountID, body)
			if err != nil {
				httpjson.WriteError(w, err)
				return
			}
			httpjson.WriteJSON(w, http.StatusOK, map[string]any{"id": email.ID})
		})

		r.Delete("/drafts/{emailId}", func(w http.ResponseWriter, req *http.Request) {
			identity := authcontext.MustIdentity(req.Context())
			uid, _ := strconv.ParseInt(identity.UserID, 10, 64)
			accountID, _ := strconv.ParseInt(chi.URLParam(req, "accountId"), 10, 64)
			emailID, _ := strconv.ParseInt(chi.URLParam(req, "emailId"), 10, 64)

			if err := service.DeleteDraft(req.Context(), uid, accountID, emailID); err != nil {
				httpjson.WriteError(w, err)
				return
			}
			httpjson.WriteJSON(w, http.StatusOK, map[string]bool{"deleted": true})
		})
	})

	router.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(authService))
		r.With(middleware.RateLimit(5, time.Minute)).Post("/mail/test-connection", func(w http.ResponseWriter, req *http.Request) {
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
	})
}
