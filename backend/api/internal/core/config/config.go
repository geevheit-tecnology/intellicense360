package config

import (
	"errors"
	"os"
	"strings"
)

type Config struct {
	AppName       string
	Env           string
	HTTPAddr      string
	AuthIssuer    string
	AuthSecret    string
	DatabaseURL   string
	StorageDriver string
}

func Load() Config {
	return Config{
		AppName:       getEnv("APP_NAME", "geevheit-intelligence-360"),
		Env:           getEnv("APP_ENV", "development"),
		HTTPAddr:      httpAddr(),
		AuthIssuer:    getEnv("AUTH_ISSUER", "geevheit"),
		AuthSecret:    getEnv("AUTH_SECRET", "development-identity-secret"),
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		StorageDriver: getEnv("STORAGE_DRIVER", "postgres"),
	}
}

func httpAddr() string {
	if value := os.Getenv("HTTP_ADDR"); value != "" {
		return value
	}
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}
	return ":8080"
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

// Validate prevents accidental fallback to volatile storage or demo secrets.
func (c Config) Validate() error {
	if c.StorageDriver != "postgres" && c.StorageDriver != "memory" {
		return errors.New("STORAGE_DRIVER must be postgres or memory")
	}
	if c.StorageDriver == "postgres" && c.DatabaseURL == "" {
		return errors.New("DATABASE_URL is required for postgres storage")
	}
	if c.Env == "production" && c.StorageDriver != "postgres" {
		return errors.New("production requires postgres storage")
	}
	if c.StorageDriver == "postgres" && (len(strings.TrimSpace(c.AuthSecret)) < 32 || c.AuthSecret == "development-identity-secret") {
		return errors.New("postgres storage requires AUTH_SECRET with at least 32 characters")
	}
	return nil
}
