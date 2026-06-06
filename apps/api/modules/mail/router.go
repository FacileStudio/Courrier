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
			token := ""
			if c, err := req.Cookie("session"); err == nil && c.Value != "" {
				token = c.Value
			} else if h := req.Header.Get("Authorization"); h != "" {
				token = h
			} else if q := req.URL.Query().Get("token"); q != "" {
				token = q
			}

			if token == "" {
				httpjson.WriteError(w, errors.Unauthorized("missing token"))
				return
			}
			uid, _, err := authService.Authenticate(req.Context(), token)
			if err != nil {
				httpjson.WriteError(w, err)
				return
			}
			userID = uid
		}

		uid, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			httpjson.WriteError(w, errors.Invalid("invalid user id"))
			return
		}
		accountID, err := strconv.ParseInt(chi.URLParam(req, "accountId"), 10, 64)
		if err != nil {
			httpjson.WriteError(w, errors.Invalid("invalid account id"))
			return
		}
		emailID, err := strconv.ParseInt(chi.URLParam(req, "emailId"), 10, 64)
		if err != nil {
			httpjson.WriteError(w, errors.Invalid("invalid email id"))
			return
		}
		cid := chi.URLParam(req, "cid")

		service.ServeCIDImage(w, req, uid, accountID, emailID, cid)
	})

	router.Route("/accounts/{accountId}/mail", func(r chi.Router) {
		r.Use(middleware.RequireAuth(authService))

		r.With(middleware.RateLimit(5, time.Minute)).Post("/sync", func(w http.ResponseWriter, req *http.Request) {
			identity := authcontext.MustIdentity(req.Context())
			uid, err := strconv.ParseInt(identity.UserID, 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid user id"))
				return
			}
			accountID, err := strconv.ParseInt(chi.URLParam(req, "accountId"), 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid account id"))
				return
			}

			if err := service.SyncAccount(req.Context(), uid, accountID); err != nil {
				httpjson.WriteError(w, err)
				return
			}
			httpjson.WriteJSON(w, http.StatusOK, map[string]bool{"synced": true})
		})

		r.Get("/folders", func(w http.ResponseWriter, req *http.Request) {
			identity := authcontext.MustIdentity(req.Context())
			uid, err := strconv.ParseInt(identity.UserID, 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid user id"))
				return
			}
			accountID, err := strconv.ParseInt(chi.URLParam(req, "accountId"), 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid account id"))
				return
			}

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
			uid, err := strconv.ParseInt(identity.UserID, 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid user id"))
				return
			}
			accountID, err := strconv.ParseInt(chi.URLParam(req, "accountId"), 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid account id"))
				return
			}
			folderID, err := strconv.ParseInt(chi.URLParam(req, "folderId"), 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid folder id"))
				return
			}

			if err := service.SyncFolderEmails(req.Context(), uid, accountID, folderID); err != nil {
				httpjson.WriteError(w, err)
				return
			}
			httpjson.WriteJSON(w, http.StatusOK, map[string]bool{"synced": true})
		})

		r.Get("/folders/{folderType}/emails", func(w http.ResponseWriter, req *http.Request) {
			identity := authcontext.MustIdentity(req.Context())
			uid, err := strconv.ParseInt(identity.UserID, 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid user id"))
				return
			}
			accountID, err := strconv.ParseInt(chi.URLParam(req, "accountId"), 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid account id"))
				return
			}
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
			uid, err := strconv.ParseInt(identity.UserID, 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid user id"))
				return
			}
			accountID, err := strconv.ParseInt(chi.URLParam(req, "accountId"), 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid account id"))
				return
			}
			emailID, err := strconv.ParseInt(chi.URLParam(req, "emailId"), 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid email id"))
				return
			}

			email, attachments, err := service.GetEmailWithAttachments(req.Context(), uid, accountID, emailID)
			if err != nil {
				httpjson.WriteError(w, err)
				return
			}
			httpjson.WriteJSON(w, http.StatusOK, emailToResponse(email, attachments...))
		})

		r.Get("/emails/{emailId}/attachments/{attachmentId}/download", func(w http.ResponseWriter, req *http.Request) {
			identity := authcontext.MustIdentity(req.Context())
			uid, err := strconv.ParseInt(identity.UserID, 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid user id"))
				return
			}
			accountID, err := strconv.ParseInt(chi.URLParam(req, "accountId"), 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid account id"))
				return
			}
			emailID, err := strconv.ParseInt(chi.URLParam(req, "emailId"), 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid email id"))
				return
			}
			attachmentID, err := strconv.ParseInt(chi.URLParam(req, "attachmentId"), 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid attachment id"))
				return
			}

			service.DownloadAttachment(w, req, uid, accountID, emailID, attachmentID)
		})

		r.Patch("/emails/{emailId}", func(w http.ResponseWriter, req *http.Request) {
			identity := authcontext.MustIdentity(req.Context())
			uid, err := strconv.ParseInt(identity.UserID, 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid user id"))
				return
			}
			accountID, err := strconv.ParseInt(chi.URLParam(req, "accountId"), 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid account id"))
				return
			}
			emailID, err := strconv.ParseInt(chi.URLParam(req, "emailId"), 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid email id"))
				return
			}

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
			uid, err := strconv.ParseInt(identity.UserID, 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid user id"))
				return
			}
			accountID, err := strconv.ParseInt(chi.URLParam(req, "accountId"), 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid account id"))
				return
			}
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
			uid, err := strconv.ParseInt(identity.UserID, 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid user id"))
				return
			}
			accountID, err := strconv.ParseInt(chi.URLParam(req, "accountId"), 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid account id"))
				return
			}

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
			uid, err := strconv.ParseInt(identity.UserID, 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid user id"))
				return
			}
			accountID, err := strconv.ParseInt(chi.URLParam(req, "accountId"), 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid account id"))
				return
			}

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
			uid, err := strconv.ParseInt(identity.UserID, 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid user id"))
				return
			}
			accountID, err := strconv.ParseInt(chi.URLParam(req, "accountId"), 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid account id"))
				return
			}
			emailID, err := strconv.ParseInt(chi.URLParam(req, "emailId"), 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid email id"))
				return
			}

			if err := service.DeleteDraft(req.Context(), uid, accountID, emailID); err != nil {
				httpjson.WriteError(w, err)
				return
			}
			httpjson.WriteJSON(w, http.StatusOK, map[string]bool{"deleted": true})
		})

		r.With(middleware.RateLimit(30, time.Minute)).Post("/emails/bulk-action", func(w http.ResponseWriter, req *http.Request) {
			identity := authcontext.MustIdentity(req.Context())
			uid, err := strconv.ParseInt(identity.UserID, 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid user id"))
				return
			}
			accountID, err := strconv.ParseInt(chi.URLParam(req, "accountId"), 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid account id"))
				return
			}

			var body BulkActionRequest
			if err := httpjson.DecodeJSON(w, req, &body); err != nil {
				httpjson.WriteError(w, err)
				return
			}

			if len(body.EmailIDs) == 0 {
				httpjson.WriteError(w, errors.Invalid("no email IDs provided"))
				return
			}
			if len(body.EmailIDs) > 200 {
				httpjson.WriteError(w, errors.Invalid("too many email IDs (max 200)"))
				return
			}

			switch body.Action {
			case "delete":
				if err := service.DeleteEmails(req.Context(), uid, accountID, body.EmailIDs); err != nil {
					httpjson.WriteError(w, err)
					return
				}
			case "archive":
				if err := service.ArchiveEmails(req.Context(), uid, accountID, body.EmailIDs); err != nil {
					httpjson.WriteError(w, err)
					return
				}
			case "mark_read":
				isRead := true
				for _, emailID := range body.EmailIDs {
					if _, err := service.UpdateEmail(req.Context(), uid, accountID, emailID, UpdateEmailRequest{IsRead: &isRead}); err != nil {
						httpjson.WriteError(w, err)
						return
					}
				}
			case "mark_unread":
				isRead := false
				for _, emailID := range body.EmailIDs {
					if _, err := service.UpdateEmail(req.Context(), uid, accountID, emailID, UpdateEmailRequest{IsRead: &isRead}); err != nil {
						httpjson.WriteError(w, err)
						return
					}
				}
			default:
				httpjson.WriteError(w, errors.Invalid("unknown action"))
				return
			}

			httpjson.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
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

	router.Route("/templates", func(r chi.Router) {
		r.Use(middleware.RequireAuth(authService))
		r.Use(middleware.RateLimit(30, time.Minute))

		r.Get("/", func(w http.ResponseWriter, req *http.Request) {
			identity := authcontext.MustIdentity(req.Context())
			uid, err := strconv.ParseInt(identity.UserID, 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid user id"))
				return
			}

			templates, err := service.ListTemplates(req.Context(), uid)
			if err != nil {
				httpjson.WriteError(w, err)
				return
			}
			resp := make([]EmailTemplateResponse, len(templates))
			for i, t := range templates {
				resp[i] = EmailTemplateResponse{
					ID:        t.ID,
					Name:      t.Name,
					Subject:   t.Subject,
					BodyHTML:  t.BodyHTML,
					BodyText:  t.BodyText,
					CreatedAt: t.CreatedAt.Format("2006-01-02T15:04:05Z"),
					UpdatedAt: t.UpdatedAt.Format("2006-01-02T15:04:05Z"),
				}
			}
			httpjson.WriteJSON(w, http.StatusOK, map[string]any{"templates": resp})
		})

		r.Post("/", func(w http.ResponseWriter, req *http.Request) {
			identity := authcontext.MustIdentity(req.Context())
			uid, err := strconv.ParseInt(identity.UserID, 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid user id"))
				return
			}

			var body EmailTemplateRequest
			if err := httpjson.DecodeJSON(w, req, &body); err != nil {
				httpjson.WriteError(w, err)
				return
			}

			tmpl, err := service.CreateTemplate(req.Context(), uid, body)
			if err != nil {
				httpjson.WriteError(w, err)
				return
			}
			httpjson.WriteJSON(w, http.StatusCreated, EmailTemplateResponse{
				ID:        tmpl.ID,
				Name:      tmpl.Name,
				Subject:   tmpl.Subject,
				BodyHTML:  tmpl.BodyHTML,
				BodyText:  tmpl.BodyText,
				CreatedAt: tmpl.CreatedAt.Format("2006-01-02T15:04:05Z"),
				UpdatedAt: tmpl.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			})
		})

		r.Put("/{templateId}", func(w http.ResponseWriter, req *http.Request) {
			identity := authcontext.MustIdentity(req.Context())
			uid, err := strconv.ParseInt(identity.UserID, 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid user id"))
				return
			}
			templateID, err := strconv.ParseInt(chi.URLParam(req, "templateId"), 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid template id"))
				return
			}

			var body EmailTemplateRequest
			if err := httpjson.DecodeJSON(w, req, &body); err != nil {
				httpjson.WriteError(w, err)
				return
			}

			tmpl, err := service.UpdateTemplate(req.Context(), uid, templateID, body)
			if err != nil {
				httpjson.WriteError(w, err)
				return
			}
			httpjson.WriteJSON(w, http.StatusOK, EmailTemplateResponse{
				ID:        tmpl.ID,
				Name:      tmpl.Name,
				Subject:   tmpl.Subject,
				BodyHTML:  tmpl.BodyHTML,
				BodyText:  tmpl.BodyText,
				CreatedAt: tmpl.CreatedAt.Format("2006-01-02T15:04:05Z"),
				UpdatedAt: tmpl.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			})
		})

		r.Delete("/{templateId}", func(w http.ResponseWriter, req *http.Request) {
			identity := authcontext.MustIdentity(req.Context())
			uid, err := strconv.ParseInt(identity.UserID, 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid user id"))
				return
			}
			templateID, err := strconv.ParseInt(chi.URLParam(req, "templateId"), 10, 64)
			if err != nil {
				httpjson.WriteError(w, errors.Invalid("invalid template id"))
				return
			}

			if err := service.DeleteTemplate(req.Context(), uid, templateID); err != nil {
				httpjson.WriteError(w, err)
				return
			}
			httpjson.WriteJSON(w, http.StatusOK, map[string]bool{"deleted": true})
		})
	})
}
