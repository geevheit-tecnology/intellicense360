# Identity Architecture

Sprint 002 establishes Identity as an independent module under:

`backend/api/internal/modules/identity`

The module follows the standard structure:

- `domain`: User, Role, Permission, Group, Tenant, Session, RefreshToken, PasswordPolicy and AuditLog.
- `application`: authentication, authorization, tenant, user, role, permission and audit services.
- `ports`: repository and service contracts.
- `infrastructure`: current in-memory repositories, ready to be replaced by PostgreSQL adapters.
- `transport`: Gin HTTP handlers.

Identity does not depend on Fleet, Checklist, Tires, CIOT or any other business module.

## Authentication Flow

1. `POST /api/v1/login` receives tenant, email and password.
2. The authentication service validates the user and password policy hash.
3. A session is created.
4. An access token is issued as HS256 JWT.
5. A refresh token is generated, hashed and stored.
6. `POST /api/v1/refresh` rotates the refresh token and revokes the old token.
7. `POST /api/v1/logout` revokes all refresh tokens for the session.

## Authorization Flow

1. Protected routes require a Bearer access token.
2. Authentication middleware validates JWT issuer, signature and expiration.
3. Claims are added to Gin and request context.
4. Permission middleware checks required permission keys.
5. Application services remain testable through `AuthorizationService`.

## Tenant Flow

1. Public login accepts `tenant_id` in JSON or `X-Tenant-ID`.
2. Protected requests resolve tenant from token claims.
3. Tenant context is propagated through `contextkeys`.
4. Repositories receive `tenantID` explicitly to keep data access tenant-scoped.

## Audit Flow

1. Global middleware records request metadata through logs.
2. Identity contains `AuditLog`, `AuditRepository` and `AuditService` contracts.
3. PostgreSQL migration includes `audit_logs` for durable audit trail.
4. Business-level audit events should be recorded by application services, not by transport-only code.

## Database Diagram

```text
tenants
  |-- users
  |-- roles
  |     |-- role_permissions -- permissions
  |-- tenant_users
  |-- groups
  |-- sessions
  |     |-- refresh_tokens
  |-- audit_logs
```

## Bootstrap

The in-memory adapter seeds:

- tenant: `bootstrap-tenant`
- user: `admin@geevheit.local`
- password: `admin1234`
- permissions: `identity.users.manage`, `identity.roles.manage`, `identity.permissions.manage`, `identity.tenants.manage`
