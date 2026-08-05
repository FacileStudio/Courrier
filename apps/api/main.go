package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/internal/crypto"
	"github.com/FacileStudio/Courrier/apps/api/internal/database"
	"github.com/FacileStudio/Courrier/apps/api/internal/env"
	"github.com/FacileStudio/Courrier/apps/api/internal/middleware"
	"github.com/FacileStudio/Courrier/apps/api/modules/accounts"
	"github.com/FacileStudio/Courrier/apps/api/modules/auth"
	"github.com/FacileStudio/Courrier/apps/api/modules/mail"
	"github.com/FacileStudio/Courrier/apps/api/modules/settings"
	"github.com/FacileStudio/Courrier/apps/api/modules/spaces"
	"github.com/FacileStudio/Courrier/apps/api/modules/users"
	"github.com/FacileStudio/Courrier/apps/api/schemas"

	"github.com/FacileStudio/Journal/sdk/journal"
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

	appEnv, err := env.Load()
	appLogger := logger.New(logger.Config{})
	if err != nil {
		appLogger.Error("failed to load config", slog.Any("error", err))
		return
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
		return
	}

	if err := schemas.Migrate(db); err != nil {
		appLogger.Error("failed to run migrations", slog.Any("error", err))
		return
	}
	go func() {
		if err := mail.BackfillThreads(db); err != nil {
			appLogger.Warn("thread backfill failed", slog.Any("error", err))
		}
	}()
	if len(appEnv.EncryptionKey) > 0 {
		if err := crypto.MigrateAccountPasswords(db, appEnv.EncryptionKey, appLogger); err != nil {
			appLogger.Warn("credential migration failed", slog.Any("error", err))
		}
		if err := crypto.MigrateOIDCTokens(db, appEnv.EncryptionKey, appLogger); err != nil {
			appLogger.Warn("OIDC token migration failed", slog.Any("error", err))
		}
	}

	if err := os.MkdirAll(filepath.Join(appEnv.StorageDir, "avatars"), 0o755); err != nil {
		appLogger.Error("failed to prepare storage", slog.Any("error", err))
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		appLogger.Error("failed to access database handle", slog.Any("error", err))
		return
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			appLogger.Error("failed to close database", slog.Any("error", err))
		}
	}()

	authService := auth.NewService(db, appEnv.StorageDir, appLogger, appEnv.EncryptionKey)
	accountService := accounts.NewService(db, appEnv.EncryptionKey)
	mailService := mail.NewService(db, appEnv.EncryptionKey)
	userService := users.NewService(db, appEnv.StorageDir)
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

	health.Mount(router, health.DB(sqlDB))
	router.Handle("/api/files/*", http.StripPrefix("/api/files/", http.FileServer(http.Dir(appEnv.StorageDir))))

	auth.RegisterRoutes(router, authService, appEnv)
	accounts.RegisterRoutes(router, accountService, authService)
	mail.RegisterRoutes(router, mailService, authService, appEnv.ResourceTokenSecret)
	users.RegisterRoutes(router, userService, authService)
	settings.RegisterRoutes(router, settingsService, authService)
	spaces.RegisterRoutes(router, spaceService, authService)

	clientDir := spa.DirFromEnv()
	if spa.Available(clientDir) {
		router.Handle("/*", middleware.Gzip(spa.Handler(spa.Config{Dir: clientDir})))
		appLogger.Info("serving client", slog.String("dir", clientDir))
	}

	addr := ":" + appEnv.Port
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
		}
	case <-shutdownSignal.Done():
		appLogger.Info("server shutting down")
		shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownContext); err != nil {
			appLogger.Error("server shutdown failed", slog.Any("error", err))
			return
		}
		appLogger.Info("server stopped")
	}
}
