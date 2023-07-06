package v1

import (
	//"net/http"

	"net/http"

	"github.com/eserzhan/onlineShop/internal/domain"
	"github.com/eserzhan/onlineShop/internal/service"
	"github.com/eserzhan/onlineShop/pkg/logger"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initAdminRoutes(api *gin.RouterGroup) {
	admins := api.Group("/admins")
	{
		admins.POST("/sign-in", h.adminSignIn)
		admins.POST("/auth/refresh", h.adminRefresh)

	authenticated := api.Group("/", h.adminIdentity)
	{
	authenticated.POST("/products", h.createProduct)
	authenticated.PUT("/products/:id", h.changeProduct)
	authenticated.DELETE("/products/:id", h.deleteProduct)
	authenticated.GET("/users", h.reponseForAdmin)
	authenticated.GET("/users/:id", h.reponseForAdmin)
	authenticated.PATCH("/users/:id", h.reponseForAdmin)
	}
}
}

func (h *Handler) reponseForAdmin(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}

func (h *Handler) adminSignIn(c *gin.Context) {
	var inp signInInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	res, err := h.services.Admins.SignIn(service.UserSignInInput{
		Email:    inp.Email,
		Password: inp.Password,
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

func (h *Handler) adminRefresh(c *gin.Context) {
	var inp refreshInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	res, err := h.services.Admins.RefreshTokens(inp.Token)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

type Product struct {
	Name string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Price float32 `json:"price" binding:"required"`
	Quantity int `json:"quantity" binding:"required"`
}

func (h *Handler) createProduct (c *gin.Context) {
	var inp Product

	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	res, err := h.services.Admins.CreateProduct(service.Product{Name: inp.Name, Description: inp.Description, Price: inp.Price, Quantity: inp.Quantity})
	if err != nil {
		logger.Error(err)
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"product created": res})
}


func (h *Handler) changeProduct(c *gin.Context) {
	var inp domain.UpdateProduct

	id := c.Param("id")
	if err := c.BindJSON(&inp); err != nil {
		logger.Error(err)
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	err := h.services.Admins.ChangeProduct(inp, id)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return 
	}

	c.JSON(http.StatusOK, map[string]interface{}{"status": "updated"})
}

func (h *Handler) deleteProduct(c *gin.Context) {
	id := c.Param("id")

	err := h.services.Admins.DeleteProduct(id)
	if err != nil {
		logger.Error(err)
		newResponse(c, http.StatusBadRequest, err.Error())
		return 
	}

	c.JSON(http.StatusOK, map[string]interface{}{"status": "ok"})
}