package app

import (
	//"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/eserzhan/onlineShop/pkg/auth"
	"github.com/eserzhan/onlineShop/pkg/database/postgres"

	"github.com/eserzhan/onlineShop/pkg/logger"

	"github.com/eserzhan/onlineShop/internal/config"
	"github.com/eserzhan/onlineShop/internal/handler"
	"github.com/eserzhan/onlineShop/internal/repository"
	"github.com/eserzhan/onlineShop/internal/server"
	"github.com/eserzhan/onlineShop/internal/service"
)

func Run() {
	cfg, err := config.InitConfig()
	if err != nil {
		logger.Error(err)
		return
	}
	db, err := postgres.NewPostgresDB(cfg)

	if err != nil {
		logger.Error(err)
		return
	}

	tokenManager, err := auth.NewManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		logger.Error(err)
		return
	}

	repo := repository.NewRepository(db)
	services := service.NewService(service.Deps{Repos: repo, TokenManager: tokenManager, AccessTokenTTL: cfg.Auth.JWT.AccessTokenTTL, RefreshTokenTTL: cfg.Auth.JWT.RefreshTokenTTL})
	handlers := handler.NewHandler(services, tokenManager)

	
	srv := server.NewServer(cfg, handlers.InitRoutes())

	if err := srv.Run(); err != nil {
		logger.Error(err)
	}

}
