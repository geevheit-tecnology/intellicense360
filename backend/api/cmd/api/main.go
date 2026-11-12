package main

import (
	"context"
	platform "github.com/geevheit/intelligence360/backend/api/internal/platform/postgres"
	"log"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/core/config"
	"github.com/geevheit/intelligence360/backend/api/internal/core/di"
)

func main() {
	cfg := config.Load()

	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}
	container := di.NewContainer(cfg)
	if cfg.StorageDriver == "postgres" {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		db, err := platform.Open(ctx, cfg.DatabaseURL)
		if err != nil {
			log.Fatal("database connection failed; check configuration and PostgreSQL availability")
		}
		defer db.Close()
		if err := platform.CheckMigrations(ctx, db); err != nil {
			cancel()
			log.Fatal(err)
		}
		cancel()
		container.WithDatabase(db)
	}
	router := container.HTTPRouter()

	if cfg.StorageDriver == "postgres" {
		// Inicia motor autônomo 98% IA para escaneamento proativo de sinais operacionais a cada 60s
		container.StartAutonomousScanner(context.Background(), 60*time.Second)
	}

	if err := router.Run(cfg.HTTPAddr); err != nil {
		log.Fatalf("api stopped: %v", err)
	}
}
