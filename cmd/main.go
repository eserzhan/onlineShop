package main

import (
	"log"
	"os"

	"github.com/yervsil/onlineShop/pkg/auth"
	"github.com/yervsil/onlineShop/pkg/database/postgres"


	"github.com/yervsil/onlineShop/internal/config"
	"github.com/yervsil/onlineShop/internal/handler"
	"github.com/yervsil/onlineShop/internal/repository"
	"github.com/yervsil/onlineShop/internal/server"
	"github.com/yervsil/onlineShop/internal/service"
	"golang.org/x/exp/slog"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	log := setupLogger(cfg.Env)

	log.Info(
		"starting the app",
		slog.String("env", cfg.Env),
		slog.String("version", "123"),
	)
	log.Debug("debug messages are enabled")

	db, err := postgres.NewPostgresDB(cfg)

	if err != nil {
		log.Error("failed to init storage", err)
		return
	}

	tokenManager, err := auth.NewManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		log.Error("failed to init tokenManager", err)
		return
	}

	repo := repository.NewRepository(db)
	services := service.NewService(service.Deps{Repos: repo, TokenManager: tokenManager, AccessTokenTTL: cfg.Auth.JWT.AccessTokenTTL, RefreshTokenTTL: cfg.Auth.JWT.RefreshTokenTTL})
	handlers := handler.NewHandler(services, tokenManager)

	
	srv := server.NewServer(cfg, handlers.InitRoutes())

	if err := srv.Run(); err != nil {
		log.Error("failed to run server", err)
		return
	}
}


func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default: // If env config is invalid, set prod settings by default due to security
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

