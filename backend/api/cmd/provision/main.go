package main

import (
	"context"
	"flag"
	"fmt"
	identity "github.com/geevheit/intelligence360/backend/api/internal/modules/identity/infrastructure/postgres"
	platform "github.com/geevheit/intelligence360/backend/api/internal/platform/postgres"
	"log"
	"os"
	"time"
)

func main() {
	name := flag.String("name", "", "company name")
	email := flag.String("email", "", "administrator email")
	flag.Parse()
	url := os.Getenv("DATABASE_URL")
	password := os.Getenv("PROVISION_ADMIN_PASSWORD")
	if url == "" || password == "" {
		log.Fatal("DATABASE_URL and PROVISION_ADMIN_PASSWORD are required; do not pass passwords as command arguments")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	db, err := platform.Open(ctx, url)
	if err != nil {
		log.Fatal("database connection failed")
	}
	defer db.Close()
	tenant, err := identity.Provision(ctx, db, *name, *email, password)
	if err != nil {
		log.Fatal("provisioning failed; verify input and migrations")
	}
	fmt.Printf("Company created. tenant_id=%s\n", tenant.ID)
}
