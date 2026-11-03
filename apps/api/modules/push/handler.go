package push

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/FacileStudio/Courrier/apps/api/internal/authcontext"
	pushsender "github.com/FacileStudio/Courrier/apps/api/internal/push"
	troncerrors "github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"
)

// Handler serves the push routes over one hub.
type Handler struct {
	hub *pushsender.Hub
}

func userID(r *http.Request) (int64, bool) {
	identity, ok := authcontext.IdentityFromContext(r.Context())
	if !ok {
		return 0, false
	}
	id, err := strconv.ParseInt(identity.UserID, 10, 64)
	return id, err == nil
}

func (h *Handler) status(w http.ResponseWriter, _ *http.Request) {
	httpjson.WriteJSON(w, http.StatusOK, StatusResponse{
		Enabled:    h.hub.Enabled(),
		Native:     h.hub.Native(),
		WebPushKey: h.hub.WebPushKey(),
	})
}

func (h *Handler) test(w http.ResponseWriter, r *http.Request) {
	id, ok := userID(r)
	if !ok {
		httpjson.WriteError(w, troncerrors.Unauthorized("not authenticated"))
		return
	}
	h.hub.Notify(r.Context(), id, pushsender.Notification{Title: "Courrier", Body: "Test notification", Path: "mail"})
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) registerDevice(w http.ResponseWriter, r *http.Request) {
	id, ok := userID(r)
	if !ok {
		httpjson.WriteError(w, troncerrors.Unauthorized("not authenticated"))
		return
	}
	var body DeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || strings.TrimSpace(body.Token) == "" {
		httpjson.WriteError(w, troncerrors.Invalid("a device token is required"))
		return
	}
	if body.Kind == pushsender.KindWeb && (body.P256dh == "" || body.Auth == "") {
		httpjson.WriteError(w, troncerrors.Invalid("a browser subscription needs its p256dh and auth keys"))
		return
	}
	err := h.hub.Register(r.Context(), pushsender.Device{
		UserID: id, Kind: body.Kind, Token: strings.TrimSpace(body.Token), P256dh: body.P256dh, Auth: body.Auth,
	})
	if errors.Is(err, pushsender.ErrUnsupported) {
		httpjson.WriteError(w, troncerrors.Conflict("this server cannot send to that kind of device"))
		return
	}
	if err != nil {
		httpjson.WriteError(w, troncerrors.Internal("failed to save the device", err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) forgetDevice(w http.ResponseWriter, r *http.Request) {
	id, ok := userID(r)
	if !ok {
		httpjson.WriteError(w, troncerrors.Unauthorized("not authenticated"))
		return
	}
	var body ForgetDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Token == "" {
		httpjson.WriteError(w, troncerrors.Invalid("a device token is required"))
		return
	}
	if err := h.hub.Unregister(r.Context(), id, body.Token); err != nil {
		httpjson.WriteError(w, troncerrors.Internal("failed to forget the device", err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
