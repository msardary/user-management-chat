package server

import (
	"context"
	"net/http"
	"user-management/internal/auth"
	"user-management/internal/chat"
	"user-management/internal/db"
	"user-management/internal/health"
	"user-management/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

func SetupRouter(authService *auth.Service, userService *user.Service, chatService *chat.Service, rdb *redis.Client) *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery(), LoggingMiddleware())

	r.GET("/", func(c *gin.Context) {
		httpRequests.Inc()
		c.String(http.StatusOK, "Welcome to User Management API")
	})
	
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	healthHandler := health.NewHandler(rdb, func(ctx context.Context) error {
		return db.Ping(ctx)
	})

	r.GET("/health", healthHandler.HealthCheck)


	api := r.Group("/api")
	{
		api.POST("/register", auth.RegisterHandler(authService))
		api.POST("/login", auth.LoginHandler(authService))
		api.POST("/refresh", auth.RefreshTokenHandler(authService))
		api.GET("/logout", auth.LogoutHandler(authService))

		api.GET("/chat/ws", chat.ChatHandler(chatService))
	}

	authGroup := r.Group("/api", auth.AuthMiddleware())
	{
		authGroup.GET("/users/me", user.GetProfile(userService))
		authGroup.PUT("/users/:id", user.UpdateProfile(userService))
	}

	adminGroup := r.Group("/admin/api", auth.AuthMiddleware(), auth.AdminMiddleware())
	{
		adminGroup.GET("/users", user.GetAllUsers(userService))
		adminGroup.PUT("/users/:id", user.UserUpdateHandler(userService))
		adminGroup.DELETE("/users/:id", user.UserDeleteHandler(userService))
	}

	return r

}