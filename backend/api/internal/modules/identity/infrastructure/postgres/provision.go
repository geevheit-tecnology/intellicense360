package postgres

import (
	"context"
	"database/sql"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/identity/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/identity/domain"
	platform "github.com/geevheit/intelligence360/backend/api/internal/platform/postgres"
	"github.com/google/uuid"
	"strings"
	"time"
)

// Provision creates a new company and its administrator in a single transaction.
// It never changes an existing company, user or password.
func Provision(ctx context.Context, db *sql.DB, name, email, password string) (domain.Tenant, error) {
	if strings.TrimSpace(name) == "" || !strings.Contains(email, "@") {
		return domain.Tenant{}, application.ErrValidation
	}
	// Validate/hash before opening a transaction.
	hash, err := application.HashPassword(password)
	if err != nil {
		return domain.Tenant{}, err
	}
	var tenant domain.Tenant
	err = (platform.Transactions{DB: db}).Run(ctx, func(ctx context.Context) error {
		store := NewStore(db)
		var err error
		tenant, err = application.NewTenantService(store.Tenants()).Create(ctx, domain.Tenant{Name: name})
		if err != nil {
			return err
		}
		now := time.Now().UTC()
		user, err := store.Users().Create(ctx, domain.User{ID: uuid.NewString(), TenantID: tenant.ID, Name: "Administrador", Email: strings.TrimSpace(email), PasswordHash: hash, Status: domain.UserStatusActive, CreatedAt: now, UpdatedAt: now})
		if err != nil {
			return err
		}
		role, err := application.NewRoleService(store.Roles()).Create(ctx, domain.Role{TenantID: tenant.ID, Name: "admin", Description: "Administrador da empresa"})
		if err != nil {
			return err
		}
		// Global platform-management permissions are not granted to a company admin.
		for _, key := range []string{
			"assets.assets.manage",
			"identity.users.manage",
			"identity.roles.manage",
			"checklist.checklists.manage",
			"checklist.checklist.manage",
			"fleet.vehicles.manage",
			"fleet.odometer.read",
			"fleet.odometer.record",
			"fleet.odometer.correct",
			"drivers.read",
			"drivers.manage",
			"documents.read",
			"documents.manage",
			"fleet.assets.manage",
			"financial.financial.manage",
			"ciot.ciot.manage",
			"fuel.fuel.read",
			"fuel.fuel.manage",
			"intelligence.intelligence.manage",
			"intelligence.intelligence.read",
			"inventory.inventory.manage",
			"maintenance.maintenance.read",
			"maintenance.maintenance.manage",
			"mission_control.read",
			"mission_control.create",
			"mission_control.update",
			"mission_control.acknowledge",
			"mission_control.resolve",
			"mission_control.dismiss",
			"mission_control.snapshot",
			"mission_control.admin",
			"suppliers.suppliers.manage",
			"tires.tires.read",
			"tires.tires.manage",
			"trips.read",
			"trips.manage",
		} {
			_, err := platform.ExecutorFor(ctx, db).ExecContext(ctx, `INSERT INTO permissions(id,key,description) VALUES($1,$2,$2) ON CONFLICT(key) DO NOTHING`, uuid.NewString(), key)
			if err != nil {
				return err
			}
			if err := store.Roles().AssignPermission(ctx, tenant.ID, role.ID, key); err != nil {
				return err
			}
		}
		_, err = platform.ExecutorFor(ctx, db).ExecContext(ctx, `INSERT INTO tenant_users(tenant_id,user_id,role_id) VALUES($1,$2,$3)`, tenant.ID, user.ID, role.ID)
		if err != nil {
			return err
		}
		return application.NewAuditService(store.Audit()).Record(ctx, domain.AuditLog{TenantID: tenant.ID, ActorID: user.ID, Action: "tenant.provisioned", ResourceType: "tenant", ResourceID: tenant.ID})
	})
	return tenant, err
}
