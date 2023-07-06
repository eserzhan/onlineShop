package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	//"github.com/onlineShop/internal/domain"
	"github.com/eserzhan/onlineShop/internal/service"
	"github.com/eserzhan/onlineShop/pkg/auth"
)

type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initUsersRoutes(v1)
		h.initAdminRoutes(v1)

		v1.GET("/products", h.getProducts)
		v1.GET("/products/:id", h.getProductById)
	}
}

func (h *Handler) getProducts(c *gin.Context) {
	res, err := h.services.GetProduct()
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return 
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"list": res,
	})
}

func (h *Handler) getProductById(c *gin.Context) {
	id := c.Param("id")

	res, err := h.services.GetProductById(id)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return 
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"item": res,
	})
}