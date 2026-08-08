package application_test

import (
	"context"
	"testing"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/identity/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/identity/infrastructure"
)

func TestAuthorizationAllowsSeededAdminPermission(t *testing.T) {
	store := infrastructure.NewMemoryStore()
	service := application.NewAuthorizationService(store.Roles())
	ok, err := service.Can(context.Background(), "bootstrap-tenant", "bootstrap-admin", "identity.users.manage")
	if err != nil {
		t.Fatalf("authorization failed: %v", err)
	}
	if !ok {
		t.Fatal("expected seeded admin to have identity.users.manage")
	}
}

func TestAuthorizationRejectsMissingPermission(t *testing.T) {
	store := infrastructure.NewMemoryStore()
	service := application.NewAuthorizationService(store.Roles())
	if err := service.Require(context.Background(), "bootstrap-tenant", "bootstrap-admin", "billing.invoices.manage"); err == nil {
		t.Fatal("expected forbidden error")
	}
}
