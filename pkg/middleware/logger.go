package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"goweb/pkg/logger"
)

func AccessLogger(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		log.Info("access",
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.FullPath()),
			zap.String("client_ip", c.ClientIP()),
			zap.Duration("latency", latency),
			zap.String("request_id", c.GetString("request_id")),
		)
	}
}
