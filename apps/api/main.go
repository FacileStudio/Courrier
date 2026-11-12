package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/internal/crypto"
	"github.com/FacileStudio/Courrier/apps/api/internal/database"
	documentation "github.com/FacileStudio/Courrier/apps/api/internal/documentation"
	"github.com/FacileStudio/Courrier/apps/api/internal/env"
	"github.com/FacileStudio/Courrier/apps/api/internal/ideas"
	"github.com/FacileStudio/Courrier/apps/api/internal/middleware"
	"github.com/FacileStudio/Courrier/apps/api/modules/accounts"
	"github.com/FacileStudio/Courrier/apps/api/modules/assist"
	"github.com/FacileStudio/Courrier/apps/api/modules/auth"
	"github.com/FacileStudio/Courrier/apps/api/modules/mail"
	"github.com/FacileStudio/Courrier/apps/api/modules/mcp"
	pushmodule "github.com/FacileStudio/Courrier/apps/api/modules/push"
	"github.com/FacileStudio/Courrier/apps/api/modules/settings"
	"github.com/FacileStudio/Courrier/apps/api/modules/spaces"
	"github.com/FacileStudio/Courrier/apps/api/modules/users"
	"github.com/FacileStudio/Courrier/apps/api/schemas"

	"github.com/FacileStudio/porte/kit"
	"github.com/FacileStudio/porte/oidc"
	"github.com/FacileStudio/porte/session"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"github.com/FacileStudio/Journal/sdk/journal"
	"github.com/FacileStudio/tronc/apiref"
	"github.com/FacileStudio/tronc/health"
	"github.com/FacileStudio/tronc/healthcheck"
	"github.com/FacileStudio/tronc/httpx"
	"github.com/FacileStudio/tronc/logger"
	troncmiddleware "github.com/FacileStudio/tronc/middleware"
	"github.com/FacileStudio/tronc/spa"
)

func main() {
	if healthcheck.Handle(os.Args) {
		return
	}

	os.Exit(run())
}

func referenceConfig() apiref.Config {
	return apiref.Config{
		Title:       "Courrier API",
		Description: "Self-hosted email client for creative studios.",
		Servers:     []string{"/api"},
		Registry: documentation.Response{
			Modules: []documentation.Module{
				auth.Documentation,
				accounts.Documentation,
				mail.Documentation,
				users.Documentation,
				settings.Documentation,
				spaces.Documentation,
				pushmodule.Documentation,
				assist.Documentation,
				mcp.Documentation,
				ideas.Documentation,
			},
		},
	}
}

// buildAuth constructs porte: one session manager, porte/kit's store over
// Courrier's users table, the Registre deletion receiver, and the OIDC kit
// over the identity tables, which refuses identities the receiver erased. Sign-in
// goes through Registre only; there is no local password kit. Discovery runs
// here, so an unreachable or half-configured issuer fails at boot rather than
// on somebody's first login.
func buildAuth(ctx context.Context, db *gorm.DB, appEnv env.Config, appLogger *slog.Logger) (services, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return services{}, err
	}
	store, err := newPorteStore(sqlDB)
	if err != nil {
		return services{}, err
	}
	resealInBackground(ctx, store, appLogger)
	users, err := auth.NewUserStore(ctx, sqlDB, appLogger)
	if err != nil {
		return services{}, err
	}
	cfg := appEnv.Porte()

	sessions, err := session.New(cfg, session.Deps{Sessions: store.Sessions(), Logger: appLogger})
	if err != nil {
		return services{}, err
	}
	deletions := newDeletions(db, users, store, appEnv, appLogger)
	deps := oidc.Deps{
		Users: users, Identities: store.Identities(), Sessions: sessions, Codes: store.LoginCodes(), Logger: appLogger,
	}
	if deletions != nil {
		deps.Erased = deletions
	}
	oidcKit, err := oidc.New(ctx, cfg, deps)
	if err != nil {
		return services{}, err
	}
	return services{
		sessions: sessions, kit: oidcKit, users: users, deletions: deletions,
		assist: newAssist(appEnv, store.Identities()), identity: store.Identities(),
	}, nil
}

// buildRouter wires every module service and the shared middleware onto the
// chi router, then mounts the API routes and the static client files.
//
// Behind Traefik and Cloudflare, RemoteAddr is only the visitor if both are
// trusted: Traefik replaces the forwarded chain rather than extending it, so
// the visitor survives in Cf-Connecting-Ip alone. TRUSTED_PROXIES=private,cloudflare
// fills all three.
func buildRouter(db *gorm.DB, dbCheck health.Check, appEnv env.Config, appLogger *slog.Logger, svc services) (chi.Router, error) {
	clientDir := spa.DirFromEnv()
	authService := auth.NewService(svc.sessions, svc.users)
	accountService := accounts.NewService(db, appEnv.EncryptionKey, svc.plans)
	userService := users.NewService(db, appEnv.StorageDir, authService, appEnv.Admins)
	settingsService := settings.NewService(len(appEnv.EncryptionKey) > 0)

	router := httpx.NewRouter(httpx.Config{
		TrustedProxies: appEnv.TrustedProxies,
		CDNProxies:     appEnv.CDNProxies,
		CDNHeader:      appEnv.CDNHeader,
		Logger:         appLogger,
		CORS: troncmiddleware.CORSConfig{
			AllowedOrigins:   appEnv.CORSAllowedOrigins,
			AllowCredentials: true,
		},
		Security: securityConfig(clientDir, appLogger),
	})

	mountRoot(router, dbCheck, svc)

	var mountErr error
	router.Route("/api", func(r chi.Router) {
		r.Handle("/files/*", http.StripPrefix("/api/files/", http.FileServer(http.Dir(appEnv.StorageDir))))

		mountErr = kit.Mount(r, kit.Config{Sessions: svc.sessions, OIDC: svc.kit})
		auth.RegisterRoutes(r, authService, appEnv)
		accounts.RegisterRoutes(r, accountService, authService)
		mail.RegisterRoutes(r, svc.mail, authService, appEnv.ResourceTokenSecret)
		pushmodule.RegisterRoutes(r, svc.push, authService)
		assist.RegisterRoutes(r, svc.assist, authService)
		users.RegisterRoutes(r, userService, authService)
		settings.RegisterRoutes(r, settingsService, authService)
		registerSpaces(r, db, appEnv, svc.directory, svc.deletions, authService)
		mountPending(r, svc, authService)
		mountErr = errors.Join(mountErr, mountIdeas(r, db, svc, authService),
			mountInvitations(r, appEnv, svc, authService, appLogger))
	})
	if mountErr != nil {
		return nil, mountErr
	}

	if spa.Available(clientDir) {
		router.Handle("/*", middleware.Gzip(spa.Handler(spa.Config{Dir: clientDir})))
		appLogger.Info("serving client", slog.String("dir", clientDir))
	}

	return router, nil
}

func run() int {
	appEnv, err := env.Load()
	appLogger := logger.New(logger.Config{})
	if err != nil {
		appLogger.Error("failed to load config", slog.Any("error", err))
		return 1
	}
	var journalClient *journal.Client
	appLogger = logger.New(logger.Config{
		Level: appEnv.LogLevel,
		Wrap: func(handler slog.Handler) slog.Handler {
			if appEnv.JournalURL == "" || appEnv.JournalToken == "" {
				return handler
			}
			journalClient = journal.New(journal.Config{URL: appEnv.JournalURL, Token: appEnv.JournalToken})
			return journal.NewHandler(journalClient, handler)
		},
	})
	if journalClient != nil {
		defer journalClient.Close()
	}

	db, err := database.Open(appEnv.DatabaseURL)
	if err != nil {
		appLogger.Error("failed to open database", slog.Any("error", err))
		return 1
	}

	if err := schemas.MigrateWithIssuer(db, appEnv.IssuerForMigration()); err != nil {
		appLogger.Error("failed to run migrations", slog.Any("error", err))
		return 1
	}
	go func() {
		if err := mail.BackfillThreads(db); err != nil {
			appLogger.Warn("thread backfill failed", slog.Any("error", err))
		}
	}()
	if err := crypto.MigrateAccountPasswords(db, appEnv.EncryptionKey, appLogger); err != nil {
		appLogger.Warn("credential migration failed", slog.Any("error", err))
	}
	if err := crypto.MigrateOIDCTokens(db, appEnv.EncryptionKey, appLogger); err != nil {
		appLogger.Warn("OIDC token migration failed", slog.Any("error", err))
	}

	if err := os.MkdirAll(filepath.Join(appEnv.StorageDir, "avatars"), 0o755); err != nil {
		appLogger.Error("failed to prepare storage", slog.Any("error", err))
		return 1
	}
	sqlDB, err := db.DB()
	if err != nil {
		appLogger.Error("failed to access database handle", slog.Any("error", err))
		return 1
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			appLogger.Error("failed to close database", slog.Any("error", err))
		}
	}()

	svc, err := buildServices(context.Background(), db, appEnv, appLogger)
	if err != nil {
		appLogger.Error("failed to build services", slog.Any("error", err))
		return 1
	}

	router, err := buildRouter(db, health.DB(sqlDB), appEnv, appLogger, svc)
	if err != nil {
		appLogger.Error("failed to build router", slog.Any("error", err))
		return 1
	}

	addr := ":" + strconv.Itoa(appEnv.Port)
	server := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	serverErrCh := make(chan error, 1)
	go func() {
		serverErrCh <- server.ListenAndServe()
	}()

	shutdownSignal, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if svc.push.Enabled() && appEnv.Push.PollInterval > 0 {
		go svc.mail.WatchInboxes(shutdownSignal, appEnv.Push.PollInterval, appLogger)
	}
	if svc.directory != nil {
		go svc.directory.Run(shutdownSignal)
	}
	if svc.deletions != nil {
		go svc.deletions.Run(shutdownSignal)
	}

	appLogger.Info("server starting", slog.String("addr", addr))
	select {
	case err := <-serverErrCh:
		if !errors.Is(err, http.ErrServerClosed) {
			appLogger.Error("server stopped", slog.Any("error", err))
			return 1
		}
	case <-shutdownSignal.Done():
		appLogger.Info("server shutting down")
		shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownContext); err != nil {
			appLogger.Error("server shutdown failed", slog.Any("error", err))
			return 1
		}
		appLogger.Info("server stopped")
	}

	return 0
}

// externalConnect lists the origins the built client calls besides itself:
// Ardoise's public plan catalog and the icon API behind iconify-icon.
var externalConnect = []string{"https://api-ardoise.facile.studio", "https://api.iconify.design"}

// securityConfig extends tronc's default security headers with what the built
// client needs: the inline scripts adapter-static writes into its pages, which
// script-src 'self' would otherwise block, and the external origins above. A
// client directory that cannot be hashed is logged and leaves the default
// script policy, which only reports.
func securityConfig(clientDir string, appLogger *slog.Logger) troncmiddleware.SecurityConfig {
	sources := map[string][]string{"connect-src": externalConnect}
	if spa.Available(clientDir) {
		hashes, err := spa.InlineScriptHashes(clientDir)
		if err != nil {
			appLogger.Warn("failed to hash the client's inline scripts", slog.Any("error", err))
		} else {
			sources["script-src"] = hashes
		}
	}
	return troncmiddleware.SecurityConfig{CSPSources: sources}
}
