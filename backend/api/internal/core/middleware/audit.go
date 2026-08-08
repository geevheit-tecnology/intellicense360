package middleware

import (
	"context"
	"log/slog"

	"github.com/geevheit/intelligence360/backend/api/internal/core/logger"
	"github.com/gin-gonic/gin"
)

func Audit(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		log.Info(context.Background(), "audit.request",
			slog.String("method", c.Request.Method),
			slog.String("path", c.FullPath()),
			slog.Int("status", c.Writer.Status()),
		)
	}
}
