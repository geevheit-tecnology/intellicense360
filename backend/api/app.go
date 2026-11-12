package api

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/core/config"
	"github.com/geevheit/intelligence360/backend/api/internal/core/di"
	platform "github.com/geevheit/intelligence360/backend/api/internal/platform/postgres"
	"github.com/gin-gonic/gin"
)

var (
	router     *gin.Engine
	routerOnce sync.Once
)

// InitRouter initializes and returns the singleton Gin engine for serverless environments.
func InitRouter() *gin.Engine {
	routerOnce.Do(func() {
		cfg := config.Load()
		if cfg.AuthSecret == "" {
			cfg.AuthSecret = "geevheit-serverless-jwt-secret-min-32-chars-long"
		}
		if cfg.DatabaseURL == "" {
			cfg.StorageDriver = "memory"
		}

		container := di.NewContainer(cfg)

		if cfg.StorageDriver == "postgres" && cfg.DatabaseURL != "" {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			db, err := platform.Open(ctx, cfg.DatabaseURL)
			if err == nil {
				container.WithDatabase(db)
			}
		}

		router = container.HTTPRouter()
	})
	return router
}

// ServeHTTP delegates HTTP requests to the initialized Gin engine.
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rtr := InitRouter()
	if rtr != nil {
		rtr.ServeHTTP(w, r)
		return
	}
	http.Error(w, "internal server error: router could not be initialized", http.StatusInternalServerError)
}
