package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/vo1dFl0w/auth-service/internal/app/adapters/storage/postgres"
	httpadapter "github.com/vo1dFl0w/auth-service/internal/app/transport/http"
	"github.com/vo1dFl0w/auth-service/internal/app/usecase"
	"github.com/vo1dFl0w/auth-service/internal/config"
	"github.com/vo1dFl0w/auth-service/internal/gen"
	"github.com/vo1dFl0w/auth-service/internal/pkg/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := run(ctx); err != nil {
		log.Println(ctx, "startup", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger := logger.LoadLogger(cfg.Env)

	databaseDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.Username, cfg.DB.Password, cfg.DB.DBname, cfg.DB.Sslmode,
	)

	db, err := sql.Open("postgres", databaseDSN)
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}

	storage := postgres.New(db)
	tokenService := usecase.NewTokenService([]byte(cfg.JWTsecret), storage.Token())
	authService := usecase.NewAuthService(storage.Auth(), storage.Token(), tokenService)

	handler := httpadapter.NewHandler(logger, authService, cfg.CookieSecure)
	secHandler := httpadapter.NewSecuredHandler(tokenService)

	server, err := gen.NewServer(handler, secHandler)
	if err != nil {
		return fmt.Errorf("failed to generate new server: %w", err)
	}

	middlewares := handler.LoggerMiddleware(handler.CorsMiddleware(server))

	srv := http.Server{
		Addr:    cfg.HttpAddr,
		Handler: middlewares,
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	serverErr := make(chan error, 1)
	go func() {
		logger.Info("server started", "host", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		} else {
			serverErr <- nil
		}
	}()

	select {
	case err := <-serverErr:
		return fmt.Errorf("server error: %w", err)
	case s := <-sig:
		logger.Info("initialization gracefull shutdown", "signal", s)
		shutdownCtx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown error: %w", err)
		}
		logger.Info("server gracefully stopped")
		return nil
	}
}
