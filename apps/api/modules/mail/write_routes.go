package mail

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/internal/middleware"
	"github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"

	"github.com/go-chi/chi/v5"
)

const maxBulkEmails = 200

// registerWriteRoutes mounts the routes that create, move and delete mail
// objects inside the /accounts/{accountId}/mail subtree. They share that
// router with the older routes: a path split between the subtree and a sibling
// group would answer 405 for the methods the other half owns.
func registerWriteRoutes(r chi.Router, service *Service) {
	r.Post("/drafts", service.handleSaveDraft)
	r.Put("/drafts/{emailId}", service.handleUpdateDraft)
	r.With(middleware.RateLimit(30, time.Minute)).Post("/emails/move", service.handleMoveEmails)
	r.Put("/sender-rules", service.handleSenderRule)
	registerFolderRoutes(r, service)
}

func urlID(req *http.Request, name, label string) (int64, error) {
	id, err := strconv.ParseInt(chi.URLParam(req, name), 10, 64)
	if err != nil {
		return 0, errors.Invalid("invalid " + label + " id")
	}
	return id, nil
}

func (s *Service) accountScope(req *http.Request) (userID, accountID int64, err error) {
	userID = middleware.UserID(req)
	accountID, err = urlID(req, "accountId", "account")
	return userID, accountID, err
}

func decodeDraftBody(w http.ResponseWriter, req *http.Request) (DraftRequest, error) {
	if strings.HasPrefix(req.Header.Get("Content-Type"), "multipart/form-data") {
		req.Body = http.MaxBytesReader(w, req.Body, 25<<20)
		sent, err := parseSendMultipart(req)
		draft := DraftRequest{SendRequest: sent}
		if err == nil {
			draft.RemoveAttachments = parseCSVInts(req.FormValue("remove_attachments"))
		}
		return draft, err
	}
	var draft DraftRequest
	err := httpjson.DecodeJSON(w, req, &draft)
	return draft, err
}

func cleanupUploads(uploads []AttachmentUpload) {
	for i := range uploads {
		uploads[i].Cleanup()
	}
}

func (s *Service) handleSaveDraft(w http.ResponseWriter, req *http.Request) {
	uid, accountID, err := s.accountScope(req)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	body, err := decodeDraftBody(w, req)
	defer cleanupUploads(body.Attachments)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	email, err := s.SaveDraft(req.Context(), uid, accountID, body.SendRequest)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, DraftResponse{ID: email.ID})
}

func (s *Service) handleUpdateDraft(w http.ResponseWriter, req *http.Request) {
	uid, accountID, err := s.accountScope(req)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	emailID, err := urlID(req, "emailId", "email")
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	body, err := decodeDraftBody(w, req)
	defer cleanupUploads(body.Attachments)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	email, err := s.UpdateDraft(req.Context(), uid, accountID, emailID, body)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, DraftResponse{ID: email.ID})
}

func (s *Service) handleMoveEmails(w http.ResponseWriter, req *http.Request) {
	uid, accountID, err := s.accountScope(req)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	var body MoveEmailsRequest
	if err := httpjson.DecodeJSON(w, req, &body); err != nil {
		httpjson.WriteError(w, err)
		return
	}
	if len(body.EmailIDs) == 0 {
		httpjson.WriteError(w, errors.Invalid("no email IDs provided"))
		return
	}
	if len(body.EmailIDs) > maxBulkEmails {
		httpjson.WriteError(w, errors.Invalid("too many email IDs (max 200)"))
		return
	}
	if err := s.MoveEmailsTo(req.Context(), uid, accountID, body.EmailIDs, body); err != nil {
		httpjson.WriteError(w, err)
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, BulkActionResponse{OK: true})
}
