package http

import (
	"tarantool-app/internal/infrastructure/logging"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Gin logging middleware using zap logger.
func GinLoggerMiddleware(logger *logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		logger.Info("HTTP Request",
			zap.String("method", c.Request.Method),
			zap.Int("status", c.Writer.Status()),
		)
	}
}
