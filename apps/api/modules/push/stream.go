package push

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	troncerrors "github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"
)

const streamHeartbeat = 25 * time.Second

// stream holds a server-sent event stream open for the desktop app, which has
// no push service: one "notification" event per notification, and a comment
// every streamHeartbeat so proxies keep the connection. It clears the read
// deadline and moves the write deadline on each write, because the server's
// ReadTimeout and WriteTimeout would otherwise end it.
func (h *Handler) stream(w http.ResponseWriter, r *http.Request) {
	id, ok := userID(r)
	if !ok {
		httpjson.WriteError(w, troncerrors.Unauthorized("not authenticated"))
		return
	}
	notes, unsubscribe := h.hub.Subscribe(id)
	defer unsubscribe()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache, no-transform")
	w.Header().Set("X-Accel-Buffering", "no")
	rc := http.NewResponseController(w)
	_ = rc.SetReadDeadline(time.Time{})
	write := func(frame string) bool {
		_ = rc.SetWriteDeadline(time.Now().Add(time.Minute))
		if _, err := fmt.Fprint(w, frame); err != nil {
			return false
		}
		return rc.Flush() == nil
	}
	w.WriteHeader(http.StatusOK)
	if !write(": ok\n\n") {
		return
	}
	ticker := time.NewTicker(streamHeartbeat)
	defer ticker.Stop()
	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			if !write(": ping\n\n") {
				return
			}
		case n := <-notes:
			data, _ := json.Marshal(StreamNotification{Title: n.Title, Body: n.Body, Path: n.Path})
			if !write("event: notification\ndata: " + string(data) + "\n\n") {
				return
			}
		}
	}
}
