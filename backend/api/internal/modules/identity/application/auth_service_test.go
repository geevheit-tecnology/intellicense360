package application_test

import (
	"context"
	"testing"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/identity/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/identity/infrastructure"
)

func buildTestAuth() application.AuthService {
	store := infrastructure.NewMemoryStore()
	return application.NewAuthService(store.Users(), store.Roles(), store.Sessions(), store.RefreshTokens(), application.NewTokenService("test", "secret"))
}

func TestLoginIssuesValidAccessToken(t *testing.T) {
	auth := buildTestAuth()
	result, err := auth.Login(context.Background(), "bootstrap-tenant", "admin@geevheit.local", "admin1234")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if result.AccessToken == "" || result.RefreshToken == "" {
		t.Fatal("expected access and refresh tokens")
	}
	claims, err := auth.ValidateAccessToken(context.Background(), result.AccessToken)
	if err != nil {
		t.Fatalf("token validation failed: %v", err)
	}
	if claims.UserID != "bootstrap-admin" || claims.TenantID != "bootstrap-tenant" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestRefreshRotatesToken(t *testing.T) {
	auth := buildTestAuth()
	result, err := auth.Login(context.Background(), "bootstrap-tenant", "admin@geevheit.local", "admin1234")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	rotated, err := auth.Refresh(context.Background(), result.RefreshToken)
	if err != nil {
		t.Fatalf("refresh failed: %v", err)
	}
	if rotated.RefreshToken == result.RefreshToken {
		t.Fatal("expected rotated refresh token")
	}
	if _, err := auth.Refresh(context.Background(), result.RefreshToken); err == nil {
		t.Fatal("expected old refresh token to be revoked")
	}
}

func TestInvalidPasswordIsRejected(t *testing.T) {
	auth := buildTestAuth()
	if _, err := auth.Login(context.Background(), "bootstrap-tenant", "admin@geevheit.local", "wrong"); err == nil {
		t.Fatal("expected invalid credentials")
	}
}
