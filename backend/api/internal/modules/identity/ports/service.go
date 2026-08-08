package ports

import (
	"context"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/identity/domain"
)

type AuthResult struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	TokenType    string      `json:"token_type"`
	ExpiresIn    int64       `json:"expires_in"`
	User         domain.User `json:"user"`
}

type AuthClaims struct {
	TenantID    string   `json:"tenant_id"`
	UserID      string   `json:"user_id"`
	SessionID   string   `json:"session_id"`
	Permissions []string `json:"permissions"`
	ExpiresAt   int64    `json:"exp"`
	Issuer      string   `json:"iss"`
}

type AuthenticationService interface {
	Login(ctx context.Context, tenantID string, email string, password string) (AuthResult, error)
	Refresh(ctx context.Context, refreshToken string) (AuthResult, error)
	Logout(ctx context.Context, refreshToken string) error
	ValidateAccessToken(ctx context.Context, accessToken string) (AuthClaims, error)
	Me(ctx context.Context, tenantID string, userID string) (domain.User, error)
}

type AuthorizationService interface {
	Can(ctx context.Context, tenantID string, userID string, permission string) (bool, error)
	Require(ctx context.Context, tenantID string, userID string, permission string) error
}

type TenantService interface {
	Resolve(ctx context.Context, tenantID string) (domain.Tenant, error)
	Create(ctx context.Context, tenant domain.Tenant) (domain.Tenant, error)
	List(ctx context.Context) ([]domain.Tenant, error)
	Update(ctx context.Context, tenant domain.Tenant) (domain.Tenant, error)
	Delete(ctx context.Context, tenantID string) error
}

type UserService interface {
	Create(ctx context.Context, user domain.User, password string) (domain.User, error)
	List(ctx context.Context, tenantID string) ([]domain.User, error)
	FindByID(ctx context.Context, tenantID string, userID string) (domain.User, error)
	Update(ctx context.Context, user domain.User) (domain.User, error)
	Delete(ctx context.Context, tenantID string, userID string) error
}

type RoleService interface {
	Create(ctx context.Context, role domain.Role) (domain.Role, error)
	List(ctx context.Context, tenantID string) ([]domain.Role, error)
	Update(ctx context.Context, role domain.Role) (domain.Role, error)
	Delete(ctx context.Context, tenantID string, roleID string) error
}

type PermissionService interface {
	Create(ctx context.Context, permission domain.Permission) (domain.Permission, error)
	List(ctx context.Context) ([]domain.Permission, error)
	Update(ctx context.Context, permission domain.Permission) (domain.Permission, error)
	Delete(ctx context.Context, key string) error
}

type AuditService interface {
	Record(ctx context.Context, event domain.AuditLog) error
	List(ctx context.Context, tenantID string) ([]domain.AuditLog, error)
}
