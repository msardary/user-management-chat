package server

import (
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)


func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		fields := log.Fields{
			"status":   status,
			"method":   method,
			"path":     path,
			"ip":       c.ClientIP(),
			"duration": duration.Seconds(),
			"service":  "user-management",
		}

		if status >= 400 {
			log.WithFields(fields).Warn("Request completed with error")
		} else {
			log.WithFields(fields).Info("Request completed successfully")
		}
	}
}