package ports

import (
	"context"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/identity/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	FindByID(ctx context.Context, tenantID string, userID string) (domain.User, error)
	FindByEmail(ctx context.Context, tenantID string, email string) (domain.User, error)
	List(ctx context.Context, tenantID string) ([]domain.User, error)
	Update(ctx context.Context, user domain.User) (domain.User, error)
	Delete(ctx context.Context, tenantID string, userID string) error
}

type RoleRepository interface {
	Create(ctx context.Context, role domain.Role) (domain.Role, error)
	FindByID(ctx context.Context, tenantID string, roleID string) (domain.Role, error)
	List(ctx context.Context, tenantID string) ([]domain.Role, error)
	Update(ctx context.Context, role domain.Role) (domain.Role, error)
	Delete(ctx context.Context, tenantID string, roleID string) error
	AssignPermission(ctx context.Context, tenantID string, roleID string, permissionKey string) error
	PermissionsForUser(ctx context.Context, tenantID string, userID string) ([]domain.Permission, error)
}

type PermissionRepository interface {
	Create(ctx context.Context, permission domain.Permission) (domain.Permission, error)
	FindByKey(ctx context.Context, key string) (domain.Permission, error)
	List(ctx context.Context) ([]domain.Permission, error)
	Update(ctx context.Context, permission domain.Permission) (domain.Permission, error)
	Delete(ctx context.Context, key string) error
}

type TenantRepository interface {
	Create(ctx context.Context, tenant domain.Tenant) (domain.Tenant, error)
	FindByID(ctx context.Context, tenantID string) (domain.Tenant, error)
	List(ctx context.Context) ([]domain.Tenant, error)
	Update(ctx context.Context, tenant domain.Tenant) (domain.Tenant, error)
	Delete(ctx context.Context, tenantID string) error
	AddUser(ctx context.Context, tenantID string, userID string) error
}

type SessionRepository interface {
	Create(ctx context.Context, session domain.Session) (domain.Session, error)
	FindByID(ctx context.Context, sessionID string) (domain.Session, error)
	Revoke(ctx context.Context, sessionID string) error
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, token domain.RefreshToken) (domain.RefreshToken, error)
	FindByHash(ctx context.Context, tokenHash string) (domain.RefreshToken, error)
	Revoke(ctx context.Context, tokenID string, replacedByID string) error
	RevokeSession(ctx context.Context, sessionID string) error
}

type AuditRepository interface {
	Create(ctx context.Context, event domain.AuditLog) (domain.AuditLog, error)
	List(ctx context.Context, tenantID string) ([]domain.AuditLog, error)
}
