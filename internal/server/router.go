package server

import (
	"context"
	"net/http"
	"user-management/internal/auth"
	"user-management/internal/chat"
	"user-management/internal/config"
	"user-management/internal/db"
	"user-management/internal/health"
	"user-management/internal/user"

	_ "user-management/docs"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	if !config.PRODUCTION {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	api := r.Group("/api/v1")
	{
		api.POST("/register", auth.RegisterHandler(authService))
		api.POST("/login", auth.LoginHandler(authService))
		api.POST("/refresh", auth.RefreshTokenHandler(authService))
		
		api.GET("/chat/ws", chat.ChatHandler(chatService))

	}

	authGroup := r.Group("/api/v1", auth.AuthMiddleware(authService))
	{

		usersGroup := authGroup.Group("/user")
		{
			usersGroup.GET("/", user.GetProfileHandler(userService))
			usersGroup.PUT("/", user.UpdateMyProfileHandler(userService))
			usersGroup.GET("/logout", auth.LogoutHandler(authService))
		}
	}

	adminGroup := r.Group("/admin/api/v1", auth.AuthMiddleware(authService), auth.AdminMiddleware())
	{
		usersAdminGroup := adminGroup.Group("/users")
        {
            usersAdminGroup.GET("", user.GetAllUsersHandler(userService))
            usersAdminGroup.PUT("/:id", user.UserUpdateHandler(userService))
            usersAdminGroup.DELETE("/:id", user.UserDeleteHandler(userService))
        }
	}

	return r

}