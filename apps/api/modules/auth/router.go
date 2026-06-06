package auth

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"api/internal/authcontext"
	"api/internal/env"
	"api/internal/errors"
	"api/internal/httpjson"
	"api/internal/middleware"
	"api/internal/resourcetoken"

	"github.com/go-chi/chi/v5"
)

const sessionCookieName = "session"

func isSecure(r *http.Request) bool {
	if r.TLS != nil {
		return true
	}
	return strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
}

func setSessionCookie(w http.ResponseWriter, r *http.Request, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(SessionTTL.Seconds()),
		HttpOnly: true,
		Secure:   isSecure(r),
		SameSite: http.SameSiteLaxMode,
	})
}

func clearSessionCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   isSecure(r),
		SameSite: http.SameSiteLaxMode,
	})
}

func RegisterRoutes(router chi.Router, service *Service, appEnv env.Config) {
	oidcEnabled := appEnv.OIDC != nil

	router.Route("/auth", func(router chi.Router) {
		router.Get("/config", func(w http.ResponseWriter, r *http.Request) {
			httpjson.WriteJSON(w, http.StatusOK, map[string]bool{
				"sso_only":     appEnv.SSOOnly,
				"oidc_enabled": oidcEnabled,
			})
		})

		if !appEnv.SSOOnly {
			router.With(middleware.RateLimit(3, time.Minute)).Post("/register", func(w http.ResponseWriter, request *http.Request) {
				var req RegisterRequest
				if err := httpjson.DecodeJSON(w, request, &req); err != nil {
					httpjson.WriteError(w, err)
					return
				}
				resp, err := service.controller.register(request.Context(), &req)
				if err != nil {
					httpjson.WriteError(w, err)
					return
				}
				setSessionCookie(w, request, resp.Token)
				httpjson.WriteJSON(w, http.StatusCreated, resp)
			})

			router.With(middleware.RateLimit(10, time.Minute)).Post("/login", func(w http.ResponseWriter, request *http.Request) {
				var req LoginRequest
				if err := httpjson.DecodeJSON(w, request, &req); err != nil {
					httpjson.WriteError(w, err)
					return
				}
				resp, err := service.controller.login(request.Context(), &req)
				if err != nil {
					httpjson.WriteError(w, err)
					return
				}
				setSessionCookie(w, request, resp.Token)
				httpjson.WriteJSON(w, http.StatusOK, resp)
			})
		}

		if len(appEnv.ResourceTokenSecret) > 0 {
			router.With(middleware.RequireAuth(service)).Get("/resource-token", func(w http.ResponseWriter, r *http.Request) {
				identity := authcontext.MustIdentity(r.Context())
				token := resourcetoken.Sign(appEnv.ResourceTokenSecret, identity.UserID, 5*time.Minute)
				httpjson.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
			})
		} else {
			router.With(middleware.RequireAuth(service)).Get("/resource-token", func(w http.ResponseWriter, r *http.Request) {
				httpjson.WriteError(w, errors.Internal("resource tokens not configured (ENCRYPTION_KEY missing)", nil))
			})
		}

		router.Post("/logout", func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(sessionCookieName)
			if err == nil && cookie.Value != "" {
				_ = service.deleteSession(r.Context(), cookie.Value)
			}
			clearSessionCookie(w, r)
			httpjson.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
		})

		if oidcEnabled {
			oidc, err := newOIDCHandler(context.Background(), appEnv.OIDC, service)
			if err != nil {
				slog.Error("failed to initialize OIDC provider", slog.Any("error", err))
			} else {
				router.Get("/oidc", oidc.login)
				router.Get("/oidc/callback", oidc.callback)
				router.With(middleware.RequireAuth(service)).Post("/sync-profile", oidc.syncProfile)
			}
		}
	})
}
