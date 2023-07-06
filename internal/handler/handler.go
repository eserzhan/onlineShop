package handler

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/eserzhan/onlineShop/internal/handler/v1"
	"github.com/eserzhan/onlineShop/internal/service"
	"github.com/eserzhan/onlineShop/pkg/auth"
)

type Handler struct {
	services *service.Services
	tokenManager auth.TokenManager
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler{
	return &Handler{
		services: services,
		tokenManager: tokenManager,
	}
}


func (h *Handler) InitRoutes() *gin.Engine {
	// Init gin handler
	router := gin.Default()



	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.services, h.tokenManager)
	log.Println(h.tokenManager)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}