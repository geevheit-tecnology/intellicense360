package middleware

import (
	"context"
	stderrors "errors"
	"log/slog"

	coreerrors "github.com/geevheit/intelligence360/backend/api/internal/core/errors"
	"github.com/geevheit/intelligence360/backend/api/internal/core/logger"
	"github.com/gin-gonic/gin"
)

func ErrorHandler(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		var appErr coreerrors.AppError
		if !stderrors.As(err, &appErr) {
			appErr = coreerrors.Internal("unexpected error")
			log.Error(context.Background(), "request failed", slog.String("error", err.Error()))
		}

		c.AbortWithStatusJSON(appErr.Status, gin.H{"error": appErr})
	}
}
