package v1

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx = "user"
	adminCtx = "admin"
)

func (h *Handler) userIdentity(c *gin.Context) {
	mapClaims, err := h.parseAuthHeader(c)
	if err != nil {
		newResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	if mapClaims["role"] != "user" {
		newResponse(c, http.StatusUnauthorized, "User is not authorized as user")
		return
	}

	c.Set(userCtx, mapClaims["userID"])
}

func (h *Handler) adminIdentity(c *gin.Context) {
	mapClaims, err := h.parseAuthHeader(c)
	if err != nil {
		newResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	if mapClaims["role"] != "admin" {
		newResponse(c, http.StatusUnauthorized, "User is not authorized as admin")
		return
	}

	c.Set(adminCtx, mapClaims["userID"])
}

func (h *Handler) parseAuthHeader(c *gin.Context) (map[string]interface{}, error) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		return map[string]interface{}{}, errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return map[string]interface{}{}, errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return map[string]interface{}{}, errors.New("token is empty")
	}
	return h.tokenManager.Parse(headerParts[1])
}