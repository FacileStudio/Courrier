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
	"github.com/FacileStudio/Courrier/apps/api/modules/accounts"
	"github.com/FacileStudio/Courrier/apps/api/modules/auth"
	"github.com/FacileStudio/Courrier/apps/api/modules/mail"
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
func buildAuth(ctx context.Context, db *gorm.DB, appEnv env.Config, appLogger *slog.Logger) (*session.Manager, *local.Kit, *oidc.Kit, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, nil, err
	}
	store := portepg.New(sqlDB)
	users := auth.NewUserStore(db)
	cfg := appEnv.Porte()

	sessions, err := session.New(cfg, session.Deps{Sessions: store.Sessions(), Logger: appLogger})
	if err != nil {
		return nil, nil, nil, err
	}
	kit, err := oidc.New(ctx, cfg, oidc.Deps{
		Users:      users,
		Identities: store.Identities(),
		Sessions:   sessions,
		Codes:      store.LoginCodes(),
		Logger:     appLogger,
	})
	if err != nil {
		return nil, nil, nil, err
	}
	// Courrier's floor has always been eight characters. porte defaults to
	// twelve, and raising it here would reject a password this app accepted
	// yesterday — a product decision, not a migration.
	passwords, err := local.New(local.Config{AllowRegistration: !appEnv.SSOOnly, MinPasswordLength: 8}, local.Deps{
		Users:      users,
		Identities: store.Identities(),
		Sessions:   sessions,
		Logger:     appLogger,
		Count:      users.CountUsers,
	})
	if err != nil {
		return nil, nil, nil, err
	}
	return sessions, passwords, kit, nil
}

func buildRouter(db *gorm.DB, dbCheck health.Check, appEnv env.Config, appLogger *slog.Logger, sessions *session.Manager, passwords *local.Kit, kit *oidc.Kit) chi.Router {
	authService := auth.NewService(db, sessions, passwords, appLogger)
	accountService := accounts.NewService(db, appEnv.EncryptionKey)
	mailService := mail.NewService(db, appEnv.EncryptionKey)
	userService := users.NewService(db, appEnv.StorageDir, authService)
	settingsService := settings.NewService(db)
	spaceService := spaces.NewService(db)

	router := httpx.NewRouter(httpx.Config{
		Logger: appLogger,
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

		sessions.Mount(r)
		kit.Mount(r)
		auth.RegisterRoutes(r, authService, appEnv)
		accounts.RegisterRoutes(r, accountService, authService)
		mail.RegisterRoutes(r, mailService, authService, appEnv.ResourceTokenSecret)
		users.RegisterRoutes(r, userService, authService)
		settings.RegisterRoutes(r, settingsService, authService)
		spaces.RegisterRoutes(r, spaceService, authService)
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

	sessions, passwords, kit, err := buildAuth(context.Background(), db, appEnv, appLogger)
	if err != nil {
		appLogger.Error("failed to build authentication", slog.Any("error", err))
		return 1
	}

	router := buildRouter(db, health.DB(sqlDB), appEnv, appLogger, sessions, passwords, kit)

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
