package config

import "testing"

func TestLoadHTTPAddrFromPORT(t *testing.T) {
	t.Setenv("HTTP_ADDR", "")
	t.Setenv("PORT", "18087")

	cfg := Load()
	if cfg.HTTPAddr != ":18087" {
		t.Fatalf("expected :18087, got %q", cfg.HTTPAddr)
	}
}

func TestLoadHTTPAddrPrecedence(t *testing.T) {
	t.Setenv("HTTP_ADDR", "127.0.0.1:19000")
	t.Setenv("PORT", "18087")

	cfg := Load()
	if cfg.HTTPAddr != "127.0.0.1:19000" {
		t.Fatalf("expected HTTP_ADDR precedence, got %q", cfg.HTTPAddr)
	}
}

func TestLoadHTTPAddrFallback(t *testing.T) {
	t.Setenv("HTTP_ADDR", "")
	t.Setenv("PORT", "")

	cfg := Load()
	if cfg.HTTPAddr != ":8080" {
		t.Fatalf("expected :8080 fallback, got %q", cfg.HTTPAddr)
	}
}

func TestPersistentConfigurationFailsClosed(t *testing.T) {
	tests := []struct {
		name  string
		cfg   Config
		valid bool
	}{
		{"missing database", Config{StorageDriver: "postgres", AuthSecret: "12345678901234567890123456789012"}, false},
		{"missing secret", Config{StorageDriver: "postgres", DatabaseURL: "postgres://localhost/db"}, false},
		{"demo secret", Config{StorageDriver: "postgres", DatabaseURL: "postgres://localhost/db", AuthSecret: "development-identity-secret"}, false},
		{"memory production", Config{Env: "production", StorageDriver: "memory"}, false},
		{"memory development", Config{StorageDriver: "memory"}, true},
		{"postgres", Config{StorageDriver: "postgres", DatabaseURL: "postgres://localhost/db", AuthSecret: "12345678901234567890123456789012"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cfg.Validate() == nil; got != tt.valid {
				t.Fatalf("valid=%v, expected %v", got, tt.valid)
			}
		})
	}
}
