package main

import (
	"context"
	"flag"
	dbplatform "github.com/geevheit/intelligence360/backend/api/internal/platform/postgres"
	"log"
	"os"
	"time"
)

func main() {
	directory := flag.String("dir", "../../database/migrations", "versioned SQL directory")
	flag.Parse()
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		log.Fatal("DATABASE_URL is required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	db, err := dbplatform.Open(ctx, url)
	if err != nil {
		log.Fatal("database connection failed; check DATABASE_URL and PostgreSQL availability")
	}
	defer db.Close()
	files, err := dbplatform.Migrate(ctx, db, *directory)
	if err != nil {
		log.Fatal(err)
	}
	for _, name := range files {
		log.Printf("applied %s", name)
	}
	log.Printf("identity/fleet/drivers/documents migrations complete: %d applied; other modules and development seeds excluded", len(files))
}
