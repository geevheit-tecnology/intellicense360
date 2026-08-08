package middleware

import (
	"context"
	"log/slog"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/core/logger"
	"github.com/gin-gonic/gin"
)

func RequestLogger(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startedAt := time.Now()
		c.Next()

		log.Info(context.Background(), "http.request",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("duration", time.Since(startedAt)),
		)
	}
}
