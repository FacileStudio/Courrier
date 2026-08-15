package spaces

import (
	"net/http"

	"github.com/FacileStudio/Courrier/apps/api/internal/authcontext"
	"github.com/FacileStudio/Courrier/apps/api/internal/middleware"
	"github.com/FacileStudio/Courrier/apps/api/modules/auth"
	"github.com/FacileStudio/tronc/httpjson"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes mounts the authenticated /spaces CRUD and membership
// endpoints on the router.
func RegisterRoutes(router chi.Router, service *Service, authService *auth.Service) {
	router.Route("/spaces", func(r chi.Router) {
		r.Use(middleware.RequireAuth(authService))

		r.Post("/", func(w http.ResponseWriter, req *http.Request) {
			uid, err := currentUserID(req)
			if err != nil {
				httpjson.WriteError(w, err)
				return
			}
			var body CreateSpaceRequest
			if err := httpjson.DecodeJSON(w, req, &body); err != nil {
				httpjson.WriteError(w, err)
				return
			}
			space, role, err := service.Create(req.Context(), uid, body)
			if err != nil {
				httpjson.WriteError(w, err)
				return
			}
			httpjson.WriteJSON(w, http.StatusCreated, toSpaceResponse(space, role))
		})

		r.Get("/", func(w http.ResponseWriter, req *http.Request) {
			uid, err := currentUserID(req)
			if err != nil {
				httpjson.WriteError(w, err)
				return
			}
			spaceList, roleMap, err := service.List(req.Context(), uid)
			if err != nil {
				httpjson.WriteError(w, err)
				return
			}
			resp := make([]SpaceResponse, len(spaceList))
			for i, s := range spaceList {
				resp[i] = toSpaceResponse(s, roleMap[s.ID])
			}
			httpjson.WriteJSON(w, http.StatusOK, map[string]any{"spaces": resp})
		})

		r.Route("/{spaceId}", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, req *http.Request) {
				uid, err := currentUserID(req)
				if err != nil {
					httpjson.WriteError(w, err)
					return
				}
				spaceID := chi.URLParam(req, "spaceId")
				space, role, err := service.Get(req.Context(), uid, spaceID)
				if err != nil {
					httpjson.WriteError(w, err)
					return
				}
				members, err := service.ListMembers(req.Context(), uid, spaceID)
				if err != nil {
					httpjson.WriteError(w, err)
					return
				}
				httpjson.WriteJSON(w, http.StatusOK, toSpaceDetailResponse(space, role, members))
			})

			r.Put("/", func(w http.ResponseWriter, req *http.Request) {
				uid, err := currentUserID(req)
				if err != nil {
					httpjson.WriteError(w, err)
					return
				}
				spaceID := chi.URLParam(req, "spaceId")
				var body UpdateSpaceRequest
				if err := httpjson.DecodeJSON(w, req, &body); err != nil {
					httpjson.WriteError(w, err)
					return
				}
				space, err := service.Update(req.Context(), uid, spaceID, body)
				if err != nil {
					httpjson.WriteError(w, err)
					return
				}
				role, _ := service.memberRole(req.Context(), spaceID, uid)
				httpjson.WriteJSON(w, http.StatusOK, toSpaceResponse(space, role))
			})

			r.Delete("/", func(w http.ResponseWriter, req *http.Request) {
				uid, err := currentUserID(req)
				if err != nil {
					httpjson.WriteError(w, err)
					return
				}
				spaceID := chi.URLParam(req, "spaceId")
				if err := service.Delete(req.Context(), uid, spaceID); err != nil {
					httpjson.WriteError(w, err)
					return
				}
				httpjson.WriteJSON(w, http.StatusOK, map[string]bool{"deleted": true})
			})

			r.Post("/leave", func(w http.ResponseWriter, req *http.Request) {
				uid, err := currentUserID(req)
				if err != nil {
					httpjson.WriteError(w, err)
					return
				}
				spaceID := chi.URLParam(req, "spaceId")
				if err := service.Leave(req.Context(), uid, spaceID); err != nil {
					httpjson.WriteError(w, err)
					return
				}
				httpjson.WriteJSON(w, http.StatusOK, map[string]bool{"left": true})
			})

			r.Route("/members", func(r chi.Router) {
				r.Get("/", func(w http.ResponseWriter, req *http.Request) {
					uid, err := currentUserID(req)
					if err != nil {
						httpjson.WriteError(w, err)
						return
					}
					spaceID := chi.URLParam(req, "spaceId")
					members, err := service.ListMembers(req.Context(), uid, spaceID)
					if err != nil {
						httpjson.WriteError(w, err)
						return
					}
					httpjson.WriteJSON(w, http.StatusOK, map[string]any{"members": members})
				})

				r.Post("/", func(w http.ResponseWriter, req *http.Request) {
					uid, err := currentUserID(req)
					if err != nil {
						httpjson.WriteError(w, err)
						return
					}
					spaceID := chi.URLParam(req, "spaceId")
					var body AddMemberRequest
					if err := httpjson.DecodeJSON(w, req, &body); err != nil {
						httpjson.WriteError(w, err)
						return
					}
					member, err := service.AddMember(req.Context(), uid, spaceID, body)
					if err != nil {
						httpjson.WriteError(w, err)
						return
					}
					httpjson.WriteJSON(w, http.StatusCreated, map[string]any{
						"id":       member.ID,
						"space_id": member.SpaceID,
						"user_id":  member.UserID,
						"role":     member.Role,
					})
				})

				r.Put("/{memberId}", func(w http.ResponseWriter, req *http.Request) {
					uid, err := currentUserID(req)
					if err != nil {
						httpjson.WriteError(w, err)
						return
					}
					spaceID := chi.URLParam(req, "spaceId")
					memberID := chi.URLParam(req, "memberId")
					var body UpdateMemberRequest
					if err := httpjson.DecodeJSON(w, req, &body); err != nil {
						httpjson.WriteError(w, err)
						return
					}
					member, err := service.UpdateMember(req.Context(), uid, spaceID, memberID, body)
					if err != nil {
						httpjson.WriteError(w, err)
						return
					}
					httpjson.WriteJSON(w, http.StatusOK, map[string]any{
						"id":   member.ID,
						"role": member.Role,
					})
				})

				r.Delete("/{memberId}", func(w http.ResponseWriter, req *http.Request) {
					uid, err := currentUserID(req)
					if err != nil {
						httpjson.WriteError(w, err)
						return
					}
					spaceID := chi.URLParam(req, "spaceId")
					memberID := chi.URLParam(req, "memberId")
					if err := service.RemoveMember(req.Context(), uid, spaceID, memberID); err != nil {
						httpjson.WriteError(w, err)
						return
					}
					httpjson.WriteJSON(w, http.StatusOK, map[string]bool{"removed": true})
				})
			})
		})
	})
}

func currentUserID(req *http.Request) (int64, error) {
	identity := authcontext.MustIdentity(req.Context())
	return parseUserID(identity.UserID)
}
