package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string
	Error   string
}

func errorResponse(c *gin.Context, err error, msg string) {
	response := ErrorResponse{
		Message: msg,
		Error:   err.Error(),
	}

	// Отправляем JSON-ответ с кодом статуса 400 (Bad Request)
	c.JSON(http.StatusBadRequest, response)
}