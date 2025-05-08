package auth

import (
	"net/http"
	"strings"
	"user-management/pkg/response"
	"user-management/pkg/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			response.Error(c, http.StatusUnauthorized, "Authorization header is missing!")
			c.Abort()
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		userID, isAdmin, err := utils.ValidateAccessToken(tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Token is not invalid!")
			c.Abort()
			return
		}

		id := int(userID)
		valid, err := service.IsRefreshTokenValid(c, int32(id))
		if err != nil || !valid {
			response.Error(c, http.StatusUnauthorized, "Refresh token is revoked or expired!")
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
			response.Error(c, http.StatusForbidden, "You are not authorized to access this resource!")
			c.Abort()
			return
		}
		c.Next()
	}
}
