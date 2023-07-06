package v1

import (

	"github.com/gin-gonic/gin"
	//"github.com/eserzhan/onlineShop/pkg/logger"
)

func newResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, gin.H{"error": message})
}
