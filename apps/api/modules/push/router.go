// Package push mounts the notification routes: the channel status the
// clients read before asking for permission, a test send, and device
// registration for FCM tokens and Web Push subscriptions.
package push

import (
	"github.com/FacileStudio/Courrier/apps/api/internal/middleware"
	pushsender "github.com/FacileStudio/Courrier/apps/api/internal/push"
	"github.com/go-chi/chi/v5"
)

// RegisterRoutes mounts /push, /push/stream and /push/devices behind the session.
func RegisterRoutes(router chi.Router, hub *pushsender.Hub, authenticator middleware.Authenticator) {
	handler := &Handler{hub: hub}
	router.Route("/push", func(r chi.Router) {
		r.Use(middleware.RequireAuth(authenticator))
		r.Get("/", handler.status)
		r.Post("/test", handler.test)
		r.Get("/stream", handler.stream)
		r.Post("/devices", handler.registerDevice)
		r.Delete("/devices", handler.forgetDevice)
	})
}
