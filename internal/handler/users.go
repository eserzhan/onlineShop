package handler

import (
	"net/http"

	"github.com/yervsil/onlineShop/internal/service"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.POST("/sign-up", h.userSignUp)
		users.POST("/sign-in", h.userSignIn)
		users.POST("/auth/refresh", h.userRefresh)
	}

	authenticated := api.Group("/", h.userIdentity)
		{
		authenticated.POST("/products/:id/cart", h.addToCart)
		authenticated.GET("/myCart", h.getCart)

		authenticated.POST("/createOrder", h.createOrder)
		authenticated.GET("/myOrders", h.orders)
		authenticated.GET("/myOrders/:id", h.getOrderById)
		}
}

type userSignUpInput struct {
	Name     string `json:"name" binding:"required,min=2,max=64"`
	Email    string `json:"email" binding:"required,email,max=64"`
	Phone    string `json:"phone" binding:"required,max=13"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

func (h *Handler) userSignUp(c *gin.Context) {
	var inp userSignUpInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	id, err := h.services.Users.SignUp(service.UserSignUpInput{
		Name:     inp.Name,
		Email:    inp.Email,
		Phone:    inp.Phone,
		Password: inp.Password,
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(200, gin.H{"id": id})
}

type signInInput struct {
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type tokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (h *Handler) userSignIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	res, err := h.services.Users.SignIn(service.UserSignInInput{
		Email:    input.Email,
		Password: input.Password,
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

type refreshInput struct {
	Token string `json:"token" binding:"required"`
}

func (h *Handler) userRefresh(c *gin.Context) {
	var inp refreshInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	res, err := h.services.Users.RefreshTokens(inp.Token)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

type quantityInput struct {
	Quantity int `json:"quantity" binding:"required"`
}

func (h *Handler) addToCart(c *gin.Context) {
	var inp quantityInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	id := c.Param("id")

	userId, ok := c.Get("user")
	if !ok {
		newResponse(c, http.StatusInternalServerError, "userId not found")
		return 
	}
	
	usrId, ok := userId.(string)
	if !ok {
		newResponse(c, http.StatusInternalServerError, "invalid type of userId")
		return 
	}

	err := h.services.Users.AddToCart(inp.Quantity, usrId, id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, "item added")
}

type orderInput struct {
	Delivery_method string `json:"delivery_method" binding:"required"`
	Payment_method string  `json:"payment_method" binding:"required"`
}

func (h *Handler) createOrder(c *gin.Context) {
	var inp orderInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}


	userId, ok := c.Get("user")
	if !ok {
		newResponse(c, http.StatusInternalServerError, "userId not found")
		return 
	}
	
	usrId, ok := userId.(string)
	if !ok {
		newResponse(c, http.StatusInternalServerError, "invalid type of userId")
		return 
	}

	orderID, err := h.services.Users.CreateOrder(usrId, inp.Delivery_method, inp.Payment_method)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"order created": orderID})
}

func (h *Handler) orders(c *gin.Context) {
	userID, ok := c.Get("user")
	if !ok {
		newResponse(c, http.StatusInternalServerError, "userId not found")
		return 
	}
	
	usrId, ok := userID.(string)
	if !ok {
		newResponse(c, http.StatusInternalServerError, "invalid type of userId")
		return 
	}

	res, err := h.services.Users.GetOrders(usrId)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}


	c.JSON(200, map[string]interface{}{
		"listOfOrders": res,
	})
}

func (h *Handler) getOrderById(c *gin.Context){
	

	userId, ok := c.Get("user")

	if !ok {
		newResponse(c, http.StatusInternalServerError, "userId not found")
		return 
	}
	
	usrid, ok := userId.(string)

	if !ok {
		newResponse(c, http.StatusInternalServerError, "invalid type of userId")
		return 
	}

	lstId := c.Param("id")

	res, err := h.services.Users.GetOrderById(usrid, lstId)

	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return 
	}

	c.JSON(200, map[string]interface{}{
		"order": res,
	})
}

func (h *Handler) getCart(c *gin.Context) {
	userID, ok := c.Get("user")
	if !ok {
		newResponse(c, http.StatusInternalServerError, "userId not found")
		return 
	}
	
	usrId, ok := userID.(string)
	if !ok {
		newResponse(c, http.StatusInternalServerError, "invalid type of userId")
		return 
	}

	res, err := h.services.Users.GetCart(usrId)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}


	c.JSON(200, map[string]interface{}{
		"cart": res,
	})
}