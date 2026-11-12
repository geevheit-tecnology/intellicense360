package api

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/core/config"
	"github.com/geevheit/intelligence360/backend/api/internal/core/di"
	platform "github.com/geevheit/intelligence360/backend/api/internal/platform/postgres"
)

// Hosted entrypoint used by platforms that embed the API in another process
// (e.g. the Vercel Go server at the repository root). Unlike cmd/api, a slow or
// sleeping database does not kill the process: initialization is retried on the
// next request and callers receive 503 meanwhile.

const (
	dbConnectTimeout = 25 * time.Second
	scannerInterval  = 60 * time.Second
)

var (
	mu      sync.Mutex
	handler http.Handler
)

func build() (http.Handler, error) {
	cfg := config.Load()
	if cfg.StorageDriver == "postgres" && cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL não configurada para STORAGE_DRIVER=postgres")
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	container := di.NewContainer(cfg)
	if cfg.StorageDriver == "postgres" {
		ctx, cancel := context.WithTimeout(context.Background(), dbConnectTimeout)
		defer cancel()
		db, err := platform.Open(ctx, cfg.DatabaseURL)
		if err != nil {
			return nil, err
		}
		if err := platform.CheckMigrations(ctx, db); err != nil {
			db.Close()
			return nil, err
		}
		container.WithDatabase(db)
	}

	router := container.HTTPRouter()
	if cfg.StorageDriver == "postgres" {
		container.StartAutonomousScanner(context.Background(), scannerInterval)
	}
	return router, nil
}

// Handler returns the API handler, initializing it on first successful use.
func Handler() (http.Handler, error) {
	mu.Lock()
	defer mu.Unlock()
	if handler != nil {
		return handler, nil
	}
	h, err := build()
	if err != nil {
		return nil, err
	}
	handler = h
	return handler, nil
}

// ServeHTTP serves an API request, answering 503 while the backend is unavailable.
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h, err := Handler()
	if err != nil {
		log.Printf("api indisponível: %v", err)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"error":"serviço temporariamente indisponível; tente novamente em instantes"}`))
		return
	}
	h.ServeHTTP(w, r)
}
