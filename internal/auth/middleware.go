package auth

import (
	"net/http"
	"strings"
	"user-management/pkg/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing!"})
			c.Abort()
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		userID, isAdmin, err := utils.ValidateAccessToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is not invalid!"})
			c.Abort()
			return
		}

		id := int(userID)
		valid, err := service.IsRefreshTokenValid(c, int32(id))
		if err != nil || !valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token is revoked or expired!"})
			c.Abort()
			return
		}
		
		c.Set("userID", userID)
		c.Set("isAdmin", isAdmin)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("isAdmin")
		if !exists || isAdmin == false {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to access this resource!"})
			c.Abort()
			return
		}
		c.Next()
	}
}
