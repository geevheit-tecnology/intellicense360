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
