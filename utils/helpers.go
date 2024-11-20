// utils/helpers.go
package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) (uint, bool) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return 0, false
	}
	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return 0, false
	}
	return userID, true
}
