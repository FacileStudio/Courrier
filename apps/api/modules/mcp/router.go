package mcp

import (
	"github.com/FacileStudio/Courrier/apps/api/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// wellKnown is the RFC 9728 metadata path at the root of the host.
const wellKnown = "/.well-known/oauth-protected-resource"

// RegisterRoutes mounts the MCP endpoint at the resource's path and its RFC
// 9728 metadata at the well-known path for that resource, plus the bare
// well-known path for clients that look only there. router is the root
// router, outside /api.
func RegisterRoutes(router chi.Router, service *Service) {
	path := service.resourcePath()
	router.Handle(path, service.endpoint())
	router.Handle(wellKnown+path, service.metadata())
	router.Handle(wellKnown, service.metadata())
}

// RegisterPendingRoutes mounts the confirm page's API under /mcp/pending on
// the /api router, for the signed-in person.
func RegisterPendingRoutes(router chi.Router, service *Service, authService middleware.Authenticator) {
	router.Route("/mcp/pending", func(r chi.Router) {
		r.Use(middleware.RequireAuth(authService))
		r.Get("/", service.handleList)
		r.Post("/{id}/confirm", service.handleConfirm)
		r.Post("/{id}/dismiss", service.handleDismiss)
	})
}

// resourcePath is the path part of the resource URI, such as /mcp.
func (s *Service) resourcePath() string {
	return s.resourceURL.EscapedPath()
}
