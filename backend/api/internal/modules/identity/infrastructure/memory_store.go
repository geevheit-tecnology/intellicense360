package infrastructure

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/identity/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/identity/domain"
)

type MemoryStore struct {
	mu            sync.RWMutex
	users         map[string]domain.User
	roles         map[string]domain.Role
	permissions   map[string]domain.Permission
	tenants       map[string]domain.Tenant
	sessions      map[string]domain.Session
	refreshTokens map[string]domain.RefreshToken
	auditLogs     []domain.AuditLog
	userRoles     map[string][]string
	rolePerms     map[string][]string
}

func NewMemoryStore() *MemoryStore {
	store := &MemoryStore{
		users: map[string]domain.User{}, roles: map[string]domain.Role{}, permissions: map[string]domain.Permission{},
		tenants: map[string]domain.Tenant{}, sessions: map[string]domain.Session{}, refreshTokens: map[string]domain.RefreshToken{},
		userRoles: map[string][]string{}, rolePerms: map[string][]string{},
	}
	now := time.Now().UTC()
	tenant := domain.Tenant{ID: "bootstrap-tenant", Name: "Bootstrap Tenant", Active: true, CreatedAt: now, UpdatedAt: now}
	admin := domain.User{ID: "bootstrap-admin", TenantID: tenant.ID, Name: "Admin", Email: "admin@geevheit.local", PasswordHash: seededPasswordHash("admin1234"), Status: domain.UserStatusActive, CreatedAt: now, UpdatedAt: now}
	role := domain.Role{ID: "bootstrap-role-admin", TenantID: tenant.ID, Name: "admin", Description: "Platform administrator", CreatedAt: now, UpdatedAt: now}
	perms := []domain.Permission{
		{ID: "perm-assets", Key: "assets.assets.manage", Description: "Manage assets", CreatedAt: now},
		{ID: "perm-users", Key: "identity.users.manage", Description: "Manage users", CreatedAt: now},
		{ID: "perm-roles", Key: "identity.roles.manage", Description: "Manage roles", CreatedAt: now},
		{ID: "perm-permissions", Key: "identity.permissions.manage", Description: "Manage permissions", CreatedAt: now},
		{ID: "perm-tenants", Key: "identity.tenants.manage", Description: "Manage tenants", CreatedAt: now},
		{ID: "perm-checklists", Key: "checklist.checklists.manage", Description: "Manage checklists", CreatedAt: now},
		{ID: "perm-checklist-engine", Key: "checklist.checklist.manage", Description: "Manage checklist engine", CreatedAt: now},
		{ID: "perm-fleet-vehicles", Key: "fleet.vehicles.manage", Description: "Manage fleet vehicles", CreatedAt: now},
		{ID: "perm-fleet-assets", Key: "fleet.assets.manage", Description: "Manage fleet assets", CreatedAt: now},
		{ID: "perm-financial", Key: "financial.financial.manage", Description: "Manage financial operations", CreatedAt: now},
		{ID: "perm-ciot", Key: "ciot.ciot.manage", Description: "Manage CIOT core", CreatedAt: now},
		{ID: "perm-fuel", Key: "fuel.fuel.manage", Description: "Manage fuel", CreatedAt: now},
		{ID: "perm-inventory", Key: "inventory.inventory.manage", Description: "Manage inventory", CreatedAt: now},
		{ID: "perm-maintenance", Key: "maintenance.maintenance.manage", Description: "Manage maintenance", CreatedAt: now},
		{ID: "perm-suppliers", Key: "suppliers.suppliers.manage", Description: "Manage suppliers", CreatedAt: now},
		{ID: "perm-tires", Key: "tires.tires.manage", Description: "Manage tires", CreatedAt: now},
	}
	store.tenants[tenant.ID] = tenant
	store.users[userKey(tenant.ID, admin.ID)] = admin
	store.roles[scopedKey(tenant.ID, role.ID)] = role
	store.userRoles[userKey(tenant.ID, admin.ID)] = []string{role.ID}
	for _, permission := range perms {
		store.permissions[permission.Key] = permission
		store.rolePerms[scopedKey(tenant.ID, role.ID)] = append(store.rolePerms[scopedKey(tenant.ID, role.ID)], permission.Key)
	}
	return store
}

func (s *MemoryStore) Users() UserRepository                 { return UserRepository{s: s} }
func (s *MemoryStore) Roles() RoleRepository                 { return RoleRepository{s: s} }
func (s *MemoryStore) Permissions() PermissionRepository     { return PermissionRepository{s: s} }
func (s *MemoryStore) Tenants() TenantRepository             { return TenantRepository{s: s} }
func (s *MemoryStore) Sessions() SessionRepository           { return SessionRepository{s: s} }
func (s *MemoryStore) RefreshTokens() RefreshTokenRepository { return RefreshTokenRepository{s: s} }
func (s *MemoryStore) Audit() AuditRepository                { return AuditRepository{s: s} }

type UserRepository struct{ s *MemoryStore }

func (r UserRepository) Create(_ context.Context, user domain.User) (domain.User, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.users[userKey(user.TenantID, user.ID)] = user
	return user, nil
}
func (r UserRepository) FindByID(_ context.Context, tenantID string, userID string) (domain.User, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	user, ok := r.s.users[userKey(tenantID, userID)]
	if !ok {
		return domain.User{}, application.ErrNotFound
	}
	return user, nil
}
func (r UserRepository) FindByEmail(_ context.Context, tenantID string, email string) (domain.User, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	for _, user := range r.s.users {
		if user.TenantID == tenantID && user.Email == email {
			return user, nil
		}
	}
	return domain.User{}, application.ErrNotFound
}
func (r UserRepository) List(_ context.Context, tenantID string) ([]domain.User, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	users := []domain.User{}
	for _, user := range r.s.users {
		if user.TenantID == tenantID {
			users = append(users, user)
		}
	}
	return users, nil
}
func (r UserRepository) Update(_ context.Context, user domain.User) (domain.User, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.users[userKey(user.TenantID, user.ID)] = user
	return user, nil
}
func (r UserRepository) Delete(_ context.Context, tenantID string, userID string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	delete(r.s.users, userKey(tenantID, userID))
	return nil
}

type RoleRepository struct{ s *MemoryStore }

func (r RoleRepository) Create(_ context.Context, role domain.Role) (domain.Role, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.roles[scopedKey(role.TenantID, role.ID)] = role
	return role, nil
}
func (r RoleRepository) FindByID(_ context.Context, tenantID string, roleID string) (domain.Role, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	role, ok := r.s.roles[scopedKey(tenantID, roleID)]
	if !ok {
		return domain.Role{}, application.ErrNotFound
	}
	return role, nil
}
func (r RoleRepository) List(_ context.Context, tenantID string) ([]domain.Role, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.Role{}
	for _, role := range r.s.roles {
		if role.TenantID == tenantID {
			items = append(items, role)
		}
	}
	return items, nil
}
func (r RoleRepository) Update(_ context.Context, role domain.Role) (domain.Role, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.roles[scopedKey(role.TenantID, role.ID)] = role
	return role, nil
}
func (r RoleRepository) Delete(_ context.Context, tenantID string, roleID string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	delete(r.s.roles, scopedKey(tenantID, roleID))
	return nil
}
func (r RoleRepository) AssignPermission(_ context.Context, tenantID string, roleID string, permissionKey string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.rolePerms[scopedKey(tenantID, roleID)] = append(r.s.rolePerms[scopedKey(tenantID, roleID)], permissionKey)
	return nil
}
func (r RoleRepository) PermissionsForUser(_ context.Context, tenantID string, userID string) ([]domain.Permission, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	out := []domain.Permission{}
	for _, roleID := range r.s.userRoles[userKey(tenantID, userID)] {
		for _, key := range r.s.rolePerms[scopedKey(tenantID, roleID)] {
			if permission, ok := r.s.permissions[key]; ok {
				out = append(out, permission)
			}
		}
	}
	return out, nil
}

type PermissionRepository struct{ s *MemoryStore }

func (r PermissionRepository) Create(_ context.Context, permission domain.Permission) (domain.Permission, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.permissions[permission.Key] = permission
	return permission, nil
}
func (r PermissionRepository) FindByKey(_ context.Context, key string) (domain.Permission, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	permission, ok := r.s.permissions[key]
	if !ok {
		return domain.Permission{}, application.ErrNotFound
	}
	return permission, nil
}
func (r PermissionRepository) List(_ context.Context) ([]domain.Permission, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.Permission{}
	for _, permission := range r.s.permissions {
		items = append(items, permission)
	}
	return items, nil
}
func (r PermissionRepository) Update(_ context.Context, permission domain.Permission) (domain.Permission, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.permissions[permission.Key] = permission
	return permission, nil
}
func (r PermissionRepository) Delete(_ context.Context, key string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	delete(r.s.permissions, key)
	return nil
}

type TenantRepository struct{ s *MemoryStore }

func (r TenantRepository) Create(_ context.Context, tenant domain.Tenant) (domain.Tenant, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.tenants[tenant.ID] = tenant
	return tenant, nil
}
func (r TenantRepository) FindByID(_ context.Context, tenantID string) (domain.Tenant, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	tenant, ok := r.s.tenants[tenantID]
	if !ok {
		return domain.Tenant{}, application.ErrNotFound
	}
	return tenant, nil
}
func (r TenantRepository) List(_ context.Context) ([]domain.Tenant, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.Tenant{}
	for _, tenant := range r.s.tenants {
		items = append(items, tenant)
	}
	return items, nil
}
func (r TenantRepository) Update(_ context.Context, tenant domain.Tenant) (domain.Tenant, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.tenants[tenant.ID] = tenant
	return tenant, nil
}
func (r TenantRepository) Delete(_ context.Context, tenantID string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	delete(r.s.tenants, tenantID)
	return nil
}
func (r TenantRepository) AddUser(_ context.Context, _ string, _ string) error { return nil }

type SessionRepository struct{ s *MemoryStore }

func (r SessionRepository) Create(_ context.Context, session domain.Session) (domain.Session, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.sessions[session.ID] = session
	return session, nil
}
func (r SessionRepository) FindByID(_ context.Context, sessionID string) (domain.Session, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	session, ok := r.s.sessions[sessionID]
	if !ok {
		return domain.Session{}, application.ErrNotFound
	}
	return session, nil
}
func (r SessionRepository) Revoke(_ context.Context, sessionID string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	session := r.s.sessions[sessionID]
	now := time.Now().UTC()
	session.RevokedAt = &now
	r.s.sessions[sessionID] = session
	return nil
}

type RefreshTokenRepository struct{ s *MemoryStore }

func (r RefreshTokenRepository) Create(_ context.Context, token domain.RefreshToken) (domain.RefreshToken, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.refreshTokens[token.TokenHash] = token
	return token, nil
}
func (r RefreshTokenRepository) FindByHash(_ context.Context, tokenHash string) (domain.RefreshToken, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	token, ok := r.s.refreshTokens[tokenHash]
	if !ok {
		return domain.RefreshToken{}, application.ErrNotFound
	}
	return token, nil
}
func (r RefreshTokenRepository) Revoke(_ context.Context, tokenID string, replacedByID string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	now := time.Now().UTC()
	for key, token := range r.s.refreshTokens {
		if token.ID == tokenID {
			token.RevokedAt = &now
			token.ReplacedByID = replacedByID
			r.s.refreshTokens[key] = token
			return nil
		}
	}
	return application.ErrNotFound
}
func (r RefreshTokenRepository) RevokeSession(_ context.Context, sessionID string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	now := time.Now().UTC()
	for key, token := range r.s.refreshTokens {
		if token.SessionID == sessionID {
			token.RevokedAt = &now
			r.s.refreshTokens[key] = token
		}
	}
	return nil
}

type AuditRepository struct{ s *MemoryStore }

func (r AuditRepository) Create(_ context.Context, event domain.AuditLog) (domain.AuditLog, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.auditLogs = append(r.s.auditLogs, event)
	return event, nil
}
func (r AuditRepository) List(_ context.Context, tenantID string) ([]domain.AuditLog, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.AuditLog{}
	for _, event := range r.s.auditLogs {
		if event.TenantID == tenantID {
			items = append(items, event)
		}
	}
	return items, nil
}

func scopedKey(tenantID string, id string) string   { return tenantID + ":" + id }
func userKey(tenantID string, userID string) string { return scopedKey(tenantID, userID) }
func seededPasswordHash(password string) string {
	sum := sha256.Sum256([]byte("seed:" + password))
	return "seed:" + hex.EncodeToString(sum[:])
}
