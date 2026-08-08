# Identity API

Base path: `/api/v1`

## Public

- `POST /login`
- `POST /refresh`

## Authenticated

- `POST /logout`
- `GET /me`
- `GET /audit-logs`

## Users

Requires `identity.users.manage`.

- `POST /users`
- `GET /users`
- `GET /users/{id}`
- `PUT /users/{id}`
- `DELETE /users/{id}`

## Roles

Requires `identity.roles.manage`.

- `POST /roles`
- `GET /roles`
- `PUT /roles/{id}`
- `DELETE /roles/{id}`

## Permissions

Requires `identity.permissions.manage`.

- `POST /permissions`
- `GET /permissions`
- `PUT /permissions/{key}`
- `DELETE /permissions/{key}`

## Tenants

Requires `identity.tenants.manage`.

- `POST /tenants`
- `GET /tenants`
- `PUT /tenants/{id}`
- `DELETE /tenants/{id}`

## Login Example

```json
{
  "tenant_id": "bootstrap-tenant",
  "email": "admin@geevheit.local",
  "password": "admin1234"
}
```
