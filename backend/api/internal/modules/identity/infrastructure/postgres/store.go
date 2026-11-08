package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/identity/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/identity/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/identity/ports"
	platform "github.com/geevheit/intelligence360/backend/api/internal/platform/postgres"
	"github.com/lib/pq"
)

type Store struct{ db *sql.DB }

func NewStore(db *sql.DB) *Store { return &Store{db: db} }
func mapError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return application.ErrNotFound
	}
	var pg *pq.Error
	if errors.As(err, &pg) {
		if pg.Code == "23505" {
			return application.ErrConflict
		}
		if pg.Code == "23503" || pg.Code == "23502" || pg.Code == "22P02" {
			return application.ErrValidation
		}
	}
	return err
}
func changed(result sql.Result, err error) error {
	if err != nil {
		return mapError(err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return application.ErrNotFound
	}
	return nil
}
func read[T any](row *sql.Row) (T, error) {
	var item T
	var body []byte
	if err := row.Scan(&body); err != nil {
		return item, mapError(err)
	}
	err := json.Unmarshal(body, &item)
	return item, err
}
func list[T any](ctx context.Context, exec platform.Executor, query string, args ...any) ([]T, error) {
	out := []T{}
	rows, err := exec.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, mapError(err)
	}
	defer rows.Close()
	for rows.Next() {
		var body []byte
		var item T
		if err := rows.Scan(&body); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(body, &item); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}
func (s *Store) Users() ports.UserRepository                 { return UserRepository{s.db} }
func (s *Store) Roles() ports.RoleRepository                 { return RoleRepository{s.db} }
func (s *Store) Permissions() ports.PermissionRepository     { return PermissionRepository{s.db} }
func (s *Store) Tenants() ports.TenantRepository             { return TenantRepository{s.db} }
func (s *Store) Sessions() ports.SessionRepository           { return SessionRepository{s.db} }
func (s *Store) RefreshTokens() ports.RefreshTokenRepository { return RefreshTokenRepository{s.db} }
func (s *Store) Audit() ports.AuditRepository                { return AuditRepository{s.db} }

type UserRepository struct{ db *sql.DB }

var _ ports.UserRepository = UserRepository{}

// Password hashes are intentionally never serialized through the domain JSON tags.
func user(row *sql.Row) (domain.User, error) {
	var u domain.User
	err := row.Scan(&u.ID, &u.TenantID, &u.Name, &u.Email, &u.PasswordHash, &u.Status, &u.CreatedAt, &u.UpdatedAt)
	return u, mapError(err)
}

const userColumns = `id,tenant_id,name,email,password_hash,status,created_at,updated_at`

func (r UserRepository) Create(ctx context.Context, u domain.User) (domain.User, error) {
	return user(platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `INSERT INTO users(id,tenant_id,name,email,password_hash,status,created_at,updated_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8) RETURNING `+userColumns, u.ID, u.TenantID, u.Name, u.Email, u.PasswordHash, u.Status, u.CreatedAt, u.UpdatedAt))
}
func (r UserRepository) FindByID(ctx context.Context, tenant, id string) (domain.User, error) {
	return user(platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `SELECT `+userColumns+` FROM users WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL`, tenant, id))
}
func (r UserRepository) FindByEmail(ctx context.Context, tenant, email string) (domain.User, error) {
	return user(platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `SELECT `+userColumns+` FROM users WHERE tenant_id=$1::uuid AND email=$2 AND deleted_at IS NULL AND EXISTS(SELECT 1 FROM tenants tn WHERE tn.id=users.tenant_id AND tn.active AND tn.deleted_at IS NULL)`, tenant, email))
}
func (r UserRepository) List(ctx context.Context, tenant string) ([]domain.User, error) {
	return list[domain.User](ctx, platform.ExecutorFor(ctx, r.db), `SELECT row_to_json(t) FROM (SELECT id,tenant_id,name,email,status,created_at,updated_at FROM users WHERE tenant_id=$1::uuid AND deleted_at IS NULL ORDER BY created_at,id) t`, tenant)
}
func (r UserRepository) Update(ctx context.Context, u domain.User) (domain.User, error) {
	return user(platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `UPDATE users SET name=$3,email=$4,status=$5,updated_at=$6,version=version+1 WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL RETURNING `+userColumns, u.TenantID, u.ID, u.Name, u.Email, u.Status, u.UpdatedAt))
}
func (r UserRepository) Delete(ctx context.Context, tenant, id string) error {
	return changed(platform.ExecutorFor(ctx, r.db).ExecContext(ctx, `UPDATE users SET deleted_at=now(),updated_at=now(),version=version+1 WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL`, tenant, id))
}

type RoleRepository struct{ db *sql.DB }

var _ ports.RoleRepository = RoleRepository{}

func (r RoleRepository) Create(ctx context.Context, v domain.Role) (domain.Role, error) {
	return read[domain.Role](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `INSERT INTO roles AS t(id,tenant_id,name,description,created_at,updated_at) VALUES($1,$2,$3,$4,$5,$6) RETURNING row_to_json(t)`, v.ID, v.TenantID, v.Name, v.Description, v.CreatedAt, v.UpdatedAt))
}
func (r RoleRepository) FindByID(ctx context.Context, tenant, id string) (domain.Role, error) {
	return read[domain.Role](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `SELECT row_to_json(t) FROM roles t WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL`, tenant, id))
}
func (r RoleRepository) List(ctx context.Context, tenant string) ([]domain.Role, error) {
	return list[domain.Role](ctx, platform.ExecutorFor(ctx, r.db), `SELECT row_to_json(t) FROM roles t WHERE tenant_id=$1::uuid AND deleted_at IS NULL ORDER BY created_at,id`, tenant)
}
func (r RoleRepository) Update(ctx context.Context, v domain.Role) (domain.Role, error) {
	return read[domain.Role](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `UPDATE roles t SET name=$3,description=$4,updated_at=$5,version=version+1 WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL RETURNING row_to_json(t)`, v.TenantID, v.ID, v.Name, v.Description, v.UpdatedAt))
}
func (r RoleRepository) Delete(ctx context.Context, tenant, id string) error {
	return changed(platform.ExecutorFor(ctx, r.db).ExecContext(ctx, `UPDATE roles SET deleted_at=now(),updated_at=now(),version=version+1 WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL`, tenant, id))
}
func (r RoleRepository) AssignPermission(ctx context.Context, tenant, id, key string) error {
	return changed(platform.ExecutorFor(ctx, r.db).ExecContext(ctx, `INSERT INTO role_permissions(role_id,permission_id) SELECT r.id,p.id FROM roles r CROSS JOIN permissions p WHERE r.tenant_id=$1::uuid AND r.id=$2::uuid AND r.deleted_at IS NULL AND p.key=$3 ON CONFLICT DO NOTHING`, tenant, id, key))
}
func (r RoleRepository) PermissionsForUser(ctx context.Context, tenant, id string) ([]domain.Permission, error) {
	return list[domain.Permission](ctx, platform.ExecutorFor(ctx, r.db), `SELECT row_to_json(p) FROM permissions p WHERE p.id IN (SELECT rp.permission_id FROM role_permissions rp JOIN roles r ON r.id=rp.role_id JOIN tenant_users tu ON tu.role_id=r.id AND tu.tenant_id=r.tenant_id JOIN users u ON u.id=tu.user_id AND u.tenant_id=tu.tenant_id WHERE tu.tenant_id=$1::uuid AND tu.user_id=$2::uuid AND r.deleted_at IS NULL AND u.deleted_at IS NULL AND u.status='active') ORDER BY p.key`, tenant, id)
}

type PermissionRepository struct{ db *sql.DB }

var _ ports.PermissionRepository = PermissionRepository{}

func (r PermissionRepository) Create(ctx context.Context, v domain.Permission) (domain.Permission, error) {
	return read[domain.Permission](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `INSERT INTO permissions AS t(id,key,description,created_at) VALUES($1,$2,$3,$4) RETURNING row_to_json(t)`, v.ID, v.Key, v.Description, v.CreatedAt))
}
func (r PermissionRepository) FindByKey(ctx context.Context, key string) (domain.Permission, error) {
	return read[domain.Permission](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `SELECT row_to_json(t) FROM permissions t WHERE key=$1`, key))
}
func (r PermissionRepository) List(ctx context.Context) ([]domain.Permission, error) {
	return list[domain.Permission](ctx, platform.ExecutorFor(ctx, r.db), `SELECT row_to_json(t) FROM permissions t ORDER BY key`)
}
func (r PermissionRepository) Update(ctx context.Context, v domain.Permission) (domain.Permission, error) {
	return read[domain.Permission](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `UPDATE permissions t SET description=$2 WHERE key=$1 RETURNING row_to_json(t)`, v.Key, v.Description))
}
func (r PermissionRepository) Delete(ctx context.Context, key string) error {
	return changed(platform.ExecutorFor(ctx, r.db).ExecContext(ctx, `DELETE FROM permissions WHERE key=$1`, key))
}

type TenantRepository struct{ db *sql.DB }

var _ ports.TenantRepository = TenantRepository{}

func (r TenantRepository) Create(ctx context.Context, v domain.Tenant) (domain.Tenant, error) {
	return read[domain.Tenant](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `INSERT INTO tenants AS t(id,name,active,created_at,updated_at) VALUES($1,$2,$3,$4,$5) RETURNING row_to_json(t)`, v.ID, v.Name, v.Active, v.CreatedAt, v.UpdatedAt))
}
func (r TenantRepository) FindByID(ctx context.Context, id string) (domain.Tenant, error) {
	return read[domain.Tenant](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `SELECT row_to_json(t) FROM tenants t WHERE id=$1::uuid AND deleted_at IS NULL`, id))
}
func (r TenantRepository) List(ctx context.Context) ([]domain.Tenant, error) {
	return list[domain.Tenant](ctx, platform.ExecutorFor(ctx, r.db), `SELECT row_to_json(t) FROM tenants t WHERE deleted_at IS NULL ORDER BY created_at,id`)
}
func (r TenantRepository) Update(ctx context.Context, v domain.Tenant) (domain.Tenant, error) {
	return read[domain.Tenant](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `UPDATE tenants t SET name=$2,active=$3,updated_at=$4,version=version+1 WHERE id=$1::uuid AND deleted_at IS NULL RETURNING row_to_json(t)`, v.ID, v.Name, v.Active, v.UpdatedAt))
}
func (r TenantRepository) Delete(ctx context.Context, id string) error {
	return changed(platform.ExecutorFor(ctx, r.db).ExecContext(ctx, `UPDATE tenants SET deleted_at=now(),active=false,updated_at=now(),version=version+1 WHERE id=$1::uuid AND deleted_at IS NULL`, id))
}
func (r TenantRepository) AddUser(ctx context.Context, tenant, id string) error {
	return changed(platform.ExecutorFor(ctx, r.db).ExecContext(ctx, `INSERT INTO tenant_users(tenant_id,user_id) SELECT tenant_id,id FROM users WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL ON CONFLICT DO NOTHING`, tenant, id))
}

type SessionRepository struct{ db *sql.DB }

var _ ports.SessionRepository = SessionRepository{}

func (r SessionRepository) Create(ctx context.Context, v domain.Session) (domain.Session, error) {
	return read[domain.Session](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `INSERT INTO sessions AS t(id,tenant_id,user_id,expires_at,created_at) VALUES($1,$2,$3,$4,$5) RETURNING row_to_json(t)`, v.ID, v.TenantID, v.UserID, v.ExpiresAt, v.CreatedAt))
}
func (r SessionRepository) FindByID(ctx context.Context, id string) (domain.Session, error) {
	return read[domain.Session](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `SELECT row_to_json(t) FROM sessions t JOIN tenants tn ON tn.id=t.tenant_id WHERE t.id=$1::uuid AND tn.active AND tn.deleted_at IS NULL`, id))
}
func (r SessionRepository) Revoke(ctx context.Context, id string) error {
	return changed(platform.ExecutorFor(ctx, r.db).ExecContext(ctx, `UPDATE sessions SET revoked_at=coalesce(revoked_at,now()) WHERE id=$1::uuid`, id))
}

type RefreshTokenRepository struct{ db *sql.DB }

var _ ports.RefreshTokenRepository = RefreshTokenRepository{}

func (r RefreshTokenRepository) Create(ctx context.Context, v domain.RefreshToken) (domain.RefreshToken, error) {
	// token_hash is excluded from domain JSON, so reads use an explicit scanner.
	_, err := platform.ExecutorFor(ctx, r.db).ExecContext(ctx, `INSERT INTO refresh_tokens(id,tenant_id,user_id,session_id,token_hash,expires_at,created_at) VALUES($1,$2,$3,$4,$5,$6,$7)`, v.ID, v.TenantID, v.UserID, v.SessionID, v.TokenHash, v.ExpiresAt, v.CreatedAt)
	return v, mapError(err)
}
func (r RefreshTokenRepository) FindByHash(ctx context.Context, hash string) (domain.RefreshToken, error) {
	query := `SELECT id,tenant_id,user_id,session_id,token_hash,expires_at,revoked_at,coalesce(replaced_by_id::text,''),created_at FROM refresh_tokens WHERE token_hash=$1`
	if platform.InTransaction(ctx) {
		query += " FOR UPDATE"
	}
	var v domain.RefreshToken
	err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query, hash).Scan(&v.ID, &v.TenantID, &v.UserID, &v.SessionID, &v.TokenHash, &v.ExpiresAt, &v.RevokedAt, &v.ReplacedByID, &v.CreatedAt)
	return v, mapError(err)
}
func (r RefreshTokenRepository) Revoke(ctx context.Context, id, replaced string) error {
	return changed(platform.ExecutorFor(ctx, r.db).ExecContext(ctx, `UPDATE refresh_tokens SET revoked_at=now(),replaced_by_id=NULLIF($2,'')::uuid WHERE id=$1::uuid AND revoked_at IS NULL`, id, replaced))
}
func (r RefreshTokenRepository) RevokeSession(ctx context.Context, id string) error {
	_, err := platform.ExecutorFor(ctx, r.db).ExecContext(ctx, `UPDATE refresh_tokens SET revoked_at=coalesce(revoked_at,now()) WHERE session_id=$1::uuid`, id)
	return mapError(err)
}

type AuditRepository struct{ db *sql.DB }

var _ ports.AuditRepository = AuditRepository{}

func (r AuditRepository) Create(ctx context.Context, v domain.AuditLog) (domain.AuditLog, error) {
	body, err := json.Marshal(v.Metadata)
	if err != nil {
		return v, err
	}
	if v.Metadata == nil {
		body = []byte("{}")
	}
	return read[domain.AuditLog](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `INSERT INTO audit_logs AS t(id,tenant_id,actor_id,action,resource_type,resource_id,metadata,created_at) VALUES($1,$2,NULLIF($3,'')::uuid,$4,$5,NULLIF($6,'')::uuid,$7::jsonb,$8) RETURNING row_to_json(t)`, v.ID, v.TenantID, v.ActorID, v.Action, v.ResourceType, v.ResourceID, string(body), v.CreatedAt))
}
func (r AuditRepository) List(ctx context.Context, tenant string) ([]domain.AuditLog, error) {
	return list[domain.AuditLog](ctx, platform.ExecutorFor(ctx, r.db), `SELECT row_to_json(t) FROM audit_logs t WHERE tenant_id=$1::uuid ORDER BY created_at,id`, tenant)
}
