package application

import (
	"context"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/identity/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/identity/ports"
)

type Services struct {
	Auth          ports.AuthenticationService
	Users         ports.UserService
	Roles         ports.RoleService
	Permissions   ports.PermissionService
	Tenants       ports.TenantService
	Authorization ports.AuthorizationService
	Audit         ports.AuditService
}

type AuthService struct {
	users         ports.UserRepository
	roles         ports.RoleRepository
	sessions      ports.SessionRepository
	refreshTokens ports.RefreshTokenRepository
	tokens        TokenService
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewAuthService(users ports.UserRepository, roles ports.RoleRepository, sessions ports.SessionRepository, refreshTokens ports.RefreshTokenRepository, tokens TokenService) AuthService {
	return AuthService{users: users, roles: roles, sessions: sessions, refreshTokens: refreshTokens, tokens: tokens, accessTTL: 15 * time.Minute, refreshTTL: 24 * time.Hour * 30}
}

func (s AuthService) Login(ctx context.Context, tenantID string, email string, password string) (ports.AuthResult, error) {
	user, err := s.users.FindByEmail(ctx, tenantID, email)
	if err != nil || user.Status != domain.UserStatusActive || !verifyPassword(user.PasswordHash, password) {
		return ports.AuthResult{}, ErrInvalidCredentials
	}
	return s.createSession(ctx, user)
}

func (s AuthService) Refresh(ctx context.Context, refreshToken string) (ports.AuthResult, error) {
	stored, err := s.refreshTokens.FindByHash(ctx, hashValue(refreshToken))
	if err != nil || stored.RevokedAt != nil || stored.ExpiresAt.Before(time.Now()) {
		return ports.AuthResult{}, ErrInvalidToken
	}
	user, err := s.users.FindByID(ctx, stored.TenantID, stored.UserID)
	if err != nil {
		return ports.AuthResult{}, err
	}
	result, err := s.createSession(ctx, user)
	if err != nil {
		return ports.AuthResult{}, err
	}
	_ = s.refreshTokens.Revoke(ctx, stored.ID, hashValue(result.RefreshToken))
	return result, nil
}

func (s AuthService) Logout(ctx context.Context, refreshToken string) error {
	stored, err := s.refreshTokens.FindByHash(ctx, hashValue(refreshToken))
	if err != nil {
		return ErrInvalidToken
	}
	if err := s.refreshTokens.RevokeSession(ctx, stored.SessionID); err != nil {
		return err
	}
	return s.sessions.Revoke(ctx, stored.SessionID)
}

func (s AuthService) ValidateAccessToken(ctx context.Context, accessToken string) (ports.AuthClaims, error) {
	return s.tokens.Validate(ctx, accessToken)
}

func (s AuthService) Me(ctx context.Context, tenantID string, userID string) (domain.User, error) {
	return s.users.FindByID(ctx, tenantID, userID)
}

func (s AuthService) createSession(ctx context.Context, user domain.User) (ports.AuthResult, error) {
	now := time.Now().UTC()
	session, err := s.sessions.Create(ctx, domain.Session{ID: newID("ses"), TenantID: user.TenantID, UserID: user.ID, ExpiresAt: now.Add(s.refreshTTL), CreatedAt: now})
	if err != nil {
		return ports.AuthResult{}, err
	}
	perms, _ := s.roles.PermissionsForUser(ctx, user.TenantID, user.ID)
	permissionKeys := make([]string, 0, len(perms))
	for _, permission := range perms {
		permissionKeys = append(permissionKeys, permission.Key)
	}
	expiresAt := now.Add(s.accessTTL)
	accessToken, err := s.tokens.Issue(ctx, ports.AuthClaims{TenantID: user.TenantID, UserID: user.ID, SessionID: session.ID, Permissions: permissionKeys, ExpiresAt: expiresAt.Unix()})
	if err != nil {
		return ports.AuthResult{}, err
	}
	refreshToken := randomToken()
	_, err = s.refreshTokens.Create(ctx, domain.RefreshToken{ID: newID("rt"), TenantID: user.TenantID, UserID: user.ID, SessionID: session.ID, TokenHash: hashValue(refreshToken), ExpiresAt: now.Add(s.refreshTTL), CreatedAt: now})
	if err != nil {
		return ports.AuthResult{}, err
	}
	return ports.AuthResult{AccessToken: accessToken, RefreshToken: refreshToken, TokenType: "Bearer", ExpiresIn: int64(s.accessTTL.Seconds()), User: user}, nil
}

type UserService struct {
	repo   ports.UserRepository
	policy domain.PasswordPolicy
}

func NewUserService(repo ports.UserRepository) UserService {
	return UserService{repo: repo, policy: domain.PasswordPolicy{MinLength: 8, RequireNumber: true, RequireSpecial: false}}
}

func (s UserService) Create(ctx context.Context, user domain.User, password string) (domain.User, error) {
	if err := validatePassword(s.policy, password); err != nil {
		return domain.User{}, err
	}
	now := time.Now().UTC()
	user.ID = newID("usr")
	user.PasswordHash = hashPassword(password)
	user.Status = domain.UserStatusActive
	user.CreatedAt = now
	user.UpdatedAt = now
	return s.repo.Create(ctx, user)
}

func (s UserService) List(ctx context.Context, tenantID string) ([]domain.User, error) {
	return s.repo.List(ctx, tenantID)
}
func (s UserService) FindByID(ctx context.Context, tenantID string, userID string) (domain.User, error) {
	return s.repo.FindByID(ctx, tenantID, userID)
}
func (s UserService) Update(ctx context.Context, user domain.User) (domain.User, error) {
	user.UpdatedAt = time.Now().UTC()
	return s.repo.Update(ctx, user)
}
func (s UserService) Delete(ctx context.Context, tenantID string, userID string) error {
	return s.repo.Delete(ctx, tenantID, userID)
}

type AuthorizationService struct{ roles ports.RoleRepository }

func NewAuthorizationService(roles ports.RoleRepository) AuthorizationService {
	return AuthorizationService{roles: roles}
}

func (s AuthorizationService) Can(ctx context.Context, tenantID string, userID string, permission string) (bool, error) {
	permissions, err := s.roles.PermissionsForUser(ctx, tenantID, userID)
	if err != nil {
		return false, err
	}
	for _, item := range permissions {
		if item.Key == permission {
			return true, nil
		}
	}
	return false, nil
}

func (s AuthorizationService) Require(ctx context.Context, tenantID string, userID string, permission string) error {
	ok, err := s.Can(ctx, tenantID, userID, permission)
	if err != nil {
		return err
	}
	if !ok {
		return ErrForbidden
	}
	return nil
}

type RoleService struct{ repo ports.RoleRepository }

func NewRoleService(repo ports.RoleRepository) RoleService { return RoleService{repo: repo} }
func (s RoleService) Create(ctx context.Context, role domain.Role) (domain.Role, error) {
	role.ID = newID("rol")
	role.CreatedAt = time.Now().UTC()
	role.UpdatedAt = role.CreatedAt
	return s.repo.Create(ctx, role)
}
func (s RoleService) List(ctx context.Context, tenantID string) ([]domain.Role, error) {
	return s.repo.List(ctx, tenantID)
}
func (s RoleService) Update(ctx context.Context, role domain.Role) (domain.Role, error) {
	role.UpdatedAt = time.Now().UTC()
	return s.repo.Update(ctx, role)
}
func (s RoleService) Delete(ctx context.Context, tenantID string, roleID string) error {
	return s.repo.Delete(ctx, tenantID, roleID)
}

type PermissionService struct{ repo ports.PermissionRepository }

func NewPermissionService(repo ports.PermissionRepository) PermissionService {
	return PermissionService{repo: repo}
}
func (s PermissionService) Create(ctx context.Context, permission domain.Permission) (domain.Permission, error) {
	permission.ID = newID("per")
	permission.CreatedAt = time.Now().UTC()
	return s.repo.Create(ctx, permission)
}
func (s PermissionService) List(ctx context.Context) ([]domain.Permission, error) {
	return s.repo.List(ctx)
}
func (s PermissionService) Update(ctx context.Context, permission domain.Permission) (domain.Permission, error) {
	return s.repo.Update(ctx, permission)
}
func (s PermissionService) Delete(ctx context.Context, key string) error {
	return s.repo.Delete(ctx, key)
}

type TenantService struct{ repo ports.TenantRepository }

func NewTenantService(repo ports.TenantRepository) TenantService { return TenantService{repo: repo} }
func (s TenantService) Resolve(ctx context.Context, tenantID string) (domain.Tenant, error) {
	return s.repo.FindByID(ctx, tenantID)
}
func (s TenantService) Create(ctx context.Context, tenant domain.Tenant) (domain.Tenant, error) {
	tenant.ID = newID("ten")
	tenant.Active = true
	tenant.CreatedAt = time.Now().UTC()
	tenant.UpdatedAt = tenant.CreatedAt
	return s.repo.Create(ctx, tenant)
}
func (s TenantService) List(ctx context.Context) ([]domain.Tenant, error) { return s.repo.List(ctx) }
func (s TenantService) Update(ctx context.Context, tenant domain.Tenant) (domain.Tenant, error) {
	tenant.UpdatedAt = time.Now().UTC()
	return s.repo.Update(ctx, tenant)
}
func (s TenantService) Delete(ctx context.Context, tenantID string) error {
	return s.repo.Delete(ctx, tenantID)
}

type AuditService struct{ repo ports.AuditRepository }

func NewAuditService(repo ports.AuditRepository) AuditService { return AuditService{repo: repo} }
func (s AuditService) Record(ctx context.Context, event domain.AuditLog) error {
	event.ID = newID("aud")
	event.CreatedAt = time.Now().UTC()
	_, err := s.repo.Create(ctx, event)
	return err
}
func (s AuditService) List(ctx context.Context, tenantID string) ([]domain.AuditLog, error) {
	return s.repo.List(ctx, tenantID)
}
