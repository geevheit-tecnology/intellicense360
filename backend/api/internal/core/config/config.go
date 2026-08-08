package config

import "os"

type Config struct {
	AppName    string
	Env        string
	HTTPAddr   string
	AuthIssuer string
	AuthSecret string
}

func Load() Config {
	return Config{
		AppName:    getEnv("APP_NAME", "geevheit-intelligence-360"),
		Env:        getEnv("APP_ENV", "development"),
		HTTPAddr:   httpAddr(),
		AuthIssuer: getEnv("AUTH_ISSUER", "geevheit"),
		AuthSecret: getEnv("AUTH_SECRET", "development-identity-secret"),
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
