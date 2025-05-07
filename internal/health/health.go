package health

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)


type Handler struct {
	RedisClient *redis.Client
	DBChecker func(ctx context.Context) error
}

func NewHandler(redisClient *redis.Client, dbChecker func(ctx context.Context) error) *Handler {
	return &Handler{
		RedisClient: redisClient,
		DBChecker: dbChecker,
	}
}

func (h *Handler) HealthCheck(c *gin.Context) {
		ctx := context.Background()

		if _, err := h.RedisClient.Ping(ctx).Result(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "failed",
				"message": "Redis unavailable",
			})
			return
		}

		if err:= h.DBChecker(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "failed",
				"message": "Database unavailable",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Service is healthy",
		})
	}