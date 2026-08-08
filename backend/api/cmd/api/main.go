package main

import (
	"log"

	"github.com/geevheit/intelligence360/backend/api/internal/core/config"
	"github.com/geevheit/intelligence360/backend/api/internal/core/di"
)

func main() {
	cfg := config.Load()

	container := di.NewContainer(cfg)
	router := container.HTTPRouter()

	if err := router.Run(cfg.HTTPAddr); err != nil {
		log.Fatalf("api stopped: %v", err)
	}
}
