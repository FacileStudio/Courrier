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
	"github.com/FacileStudio/Courrier/apps/api/internal/middleware"
	"github.com/FacileStudio/Courrier/apps/api/internal/push"
	"github.com/FacileStudio/Courrier/apps/api/modules/accounts"
	"github.com/FacileStudio/Courrier/apps/api/modules/auth"
	"github.com/FacileStudio/Courrier/apps/api/modules/mail"
	pushmodule "github.com/FacileStudio/Courrier/apps/api/modules/push"
	"github.com/FacileStudio/Courrier/apps/api/modules/settings"
	"github.com/FacileStudio/Courrier/apps/api/modules/spaces"
	"github.com/FacileStudio/Courrier/apps/api/modules/users"
	"github.com/FacileStudio/Courrier/apps/api/schemas"

	"github.com/FacileStudio/porte/local"
	"github.com/FacileStudio/porte/oidc"
	portepg "github.com/FacileStudio/porte/pg"
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
			},
		},
	}
}

// buildAuth constructs porte: one session manager, shared by the OIDC kit and
// the local login, over the identity tables.
//
// One manager and not two: they would each keep their own idea of the clock
// and of whether the cookie is Secure, and porte refuses a kit whose config
// disagrees with its manager's for exactly that reason. Discovery runs here,
// so an unreachable or half-configured issuer fails at boot rather than on
// somebody's first login — a change from what this app did, where a discovery
// failure at route-registration time logged an error and left SSO 404ing until
// the next restart.
//
// Courrier's password floor has always been eight characters while porte
// defaults to twelve; raising it here would reject a password this app
// accepted yesterday — a product decision, not a migration.
func buildAuth(ctx context.Context, db *gorm.DB, appEnv env.Config, appLogger *slog.Logger) (services, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return services{}, err
	}
	store := portepg.New(sqlDB)
	users := auth.NewUserStore(db)
	cfg := appEnv.Porte()

	sessions, err := session.New(cfg, session.Deps{Sessions: store.Sessions(), Logger: appLogger})
	if err != nil {
		return services{}, err
	}
	kit, err := oidc.New(ctx, cfg, oidc.Deps{
		Users:      users,
		Identities: store.Identities(),
		Sessions:   sessions,
		Codes:      store.LoginCodes(),
		Logger:     appLogger,
	})
	if err != nil {
		return services{}, err
	}
	passwords, err := local.New(local.Config{AllowRegistration: !appEnv.SSOOnly, MinPasswordLength: 8}, local.Deps{
		Users:      users,
		Identities: store.Identities(),
		Sessions:   sessions,
		Logger:     appLogger,
		Count:      users.CountUsers,
	})
	if err != nil {
		return services{}, err
	}
	return services{sessions: sessions, passwords: passwords, kit: kit}, nil
}

// services are the long-lived pieces run builds once: porte's session manager,
// password login and OIDC kit, the mail service the watcher shares with the
// HTTP routes, the push hub, and the Registre space sync when SSO is on.
type services struct {
	sessions  *session.Manager
	passwords *local.Kit
	kit       *oidc.Kit
	mail      *mail.Service
	push      *push.Hub
	directory *spaces.Directory
}

// buildRouter wires every module service and the shared middleware onto the
// chi router, then mounts the API routes and the static client files.
//
// Behind Traefik and Cloudflare, RemoteAddr is only the visitor if both are
// trusted: Traefik replaces the forwarded chain rather than extending it, so
// the visitor survives in Cf-Connecting-Ip alone. TRUSTED_PROXIES=private,cloudflare
// fills all three.
func buildRouter(db *gorm.DB, dbCheck health.Check, appEnv env.Config, appLogger *slog.Logger, svc services) chi.Router {
	authService := auth.NewService(db, svc.sessions, svc.passwords, appLogger)
	accountService := accounts.NewService(db, appEnv.EncryptionKey)
	userService := users.NewService(db, appEnv.StorageDir, authService)
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
	})
	router.Use(middleware.SecurityHeaders)

	health.Mount(router, dbCheck)
	apiref.Mount(router, referenceConfig())

	router.Route("/api", func(r chi.Router) {
		r.Handle("/files/*", http.StripPrefix("/api/files/", http.FileServer(http.Dir(appEnv.StorageDir))))

		svc.sessions.Mount(r)
		svc.kit.Mount(r)
		auth.RegisterRoutes(r, authService, appEnv)
		accounts.RegisterRoutes(r, accountService, authService)
		mail.RegisterRoutes(r, svc.mail, authService, appEnv.ResourceTokenSecret)
		pushmodule.RegisterRoutes(r, svc.push, authService)
		users.RegisterRoutes(r, userService, authService)
		settings.RegisterRoutes(r, settingsService, authService)
		registerSpaces(r, db, appEnv, svc.directory, authService)
	})

	clientDir := spa.DirFromEnv()
	if spa.Available(clientDir) {
		router.Handle("/*", middleware.Gzip(spa.Handler(spa.Config{Dir: clientDir})))
		appLogger.Info("serving client", slog.String("dir", clientDir))
	}

	return router
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

	svc, err := buildAuth(context.Background(), db, appEnv, appLogger)
	if err != nil {
		appLogger.Error("failed to build authentication", slog.Any("error", err))
		return 1
	}
	svc.push = newPushHub(db, appEnv.Push, appLogger)
	svc.mail = mail.NewService(db, appEnv.EncryptionKey).WithNotifier(svc.push)
	svc.directory = newDirectory(db, appEnv, appLogger)

	router := buildRouter(db, health.DB(sqlDB), appEnv, appLogger, svc)

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
