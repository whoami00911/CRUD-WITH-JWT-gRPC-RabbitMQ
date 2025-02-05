package handlers

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func (har *HandlerAssetsResponse) userIdentity(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		c.JSON(401, gin.H{"error": "Invalid authorization token (is empty)"})
		c.Abort()
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		c.JSON(401, gin.H{"error": "Invalid authorization token format"})
		c.Abort()
		return
	}

	if len(headerParts[1]) == 0 {
		c.JSON(401, gin.H{"error": "Empty token"})
		c.Abort()
		return
	}

	userId, err := har.service.Session.ParseToken(headerParts[1])
	if err != nil {
		har.Logger.Error(fmt.Sprintf("failed to parse token: %s", err))
		c.JSON(401, gin.H{"error": "Invalid or expired token"})
		c.Abort() // Прерываем выполнение запроса
		return
	}

	c.Set("UserId", userId)
}
