package mail

import (
	"net/http"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/internal/middleware"
	"github.com/FacileStudio/tronc/httpjson"

	"github.com/go-chi/chi/v5"
)

func registerFolderRoutes(r chi.Router, service *Service) {
	folders := r.With(middleware.RateLimit(10, time.Minute))
	folders.Post("/folders", service.handleCreateFolder)
	folders.Patch("/folders/{folderId}", service.handleRenameFolder)
	folders.Delete("/folders/{folderId}", service.handleDeleteFolder)
}

func (s *Service) handleCreateFolder(w http.ResponseWriter, req *http.Request) {
	uid, accountID, err := s.accountScope(req)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	var body FolderNameRequest
	if err := httpjson.DecodeJSON(w, req, &body); err != nil {
		httpjson.WriteError(w, err)
		return
	}
	folder, err := s.CreateFolder(req.Context(), uid, accountID, body.Name)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	httpjson.WriteJSON(w, http.StatusCreated, folderToResponse(folder))
}

func (s *Service) handleRenameFolder(w http.ResponseWriter, req *http.Request) {
	uid, accountID, err := s.accountScope(req)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	folderID, err := urlID(req, "folderId", "folder")
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	var body FolderNameRequest
	if err := httpjson.DecodeJSON(w, req, &body); err != nil {
		httpjson.WriteError(w, err)
		return
	}
	folder, err := s.RenameFolder(req.Context(), uid, accountID, folderID, body.Name)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, folderToResponse(folder))
}

func (s *Service) handleDeleteFolder(w http.ResponseWriter, req *http.Request) {
	uid, accountID, err := s.accountScope(req)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	folderID, err := urlID(req, "folderId", "folder")
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	force := req.URL.Query().Get("force") == "true"
	if err := s.DeleteFolder(req.Context(), uid, accountID, folderID, force); err != nil {
		httpjson.WriteError(w, err)
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, DeleteFolderResponse{Deleted: true})
}
