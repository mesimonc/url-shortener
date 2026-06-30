package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// NewLogger creates a new zap logger.
func NewLogger() *zap.Logger {
	logger, _ := zap.NewProduction()
	return logger
}

// LoggerMiddleware logs each request in structured JSON format.
func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		logger.Info("request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
			zap.String("ip", c.ClientIP()),
		)
	}
}
