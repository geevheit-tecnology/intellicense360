SPRINT 002 — Identity & Platform Foundation

Read before starting:

docs/00_PRODUCT_MANIFESTO.md
docs/01_MASTER_PROMPT.md
docs/ARCHITECTURE.md
docs/BACKEND_ARCHITECTURE.md

Do NOT implement Fleet business rules yet.

Objective:

Build the complete Identity Platform and prepare the enterprise backend foundation.

Tasks

1. Refactor all modules to follow the same architecture:

domain/
application/
infrastructure/
transport/
ports/

2. Create internal/platform

database
cache
messaging
storage
mailer
scheduler
telemetry
search

3. Implement Identity module

User

Role

Permission

Group

Tenant

Session

Refresh Token

Password Policy

Authentication

Authorization

4. JWT Authentication

Access Token

Refresh Token

Token Rotation

Logout

Revoke Token

Expiration

5. RBAC

Permissions

Roles

Permission Matrix

Middleware

Authorization Service

6. Multi Tenant

Tenant Resolver

Tenant Middleware

Tenant Context

Tenant Repository

Tenant Service

7. Audit

Audit Events

Audit Repository

Audit Service

Audit Middleware

Audit Trail

8. Database

Create migrations for

users

roles

permissions

role_permissions

tenants

tenant_users

sessions

audit_logs

refresh_tokens

9. API

Authentication

POST /login

POST /refresh

POST /logout

GET /me

Users

CRUD

Roles

CRUD

Permissions

CRUD

Tenants

CRUD

10. OpenAPI

Generate complete Swagger documentation.

11. Testing

Unit Tests

Integration Tests

Authentication Tests

RBAC Tests

Tenant Tests

12. Documentation

Generate:

Identity Architecture

Authentication Flow

Authorization Flow

Tenant Flow

Audit Flow

Database Diagram

API Documentation

Acceptance Criteria

Do NOT implement Fleet.

Do NOT implement Checklist.

Do NOT implement Tires.

Do NOT implement CIOT.

Goal:

Enterprise Identity Platform ready for production.

IMPORTANT

Do not connect modules.

Do not create cross-module dependencies.

Every module must be independently testable.

Business integrations will be implemented only in future sprints.

Focus only on creating stable contracts, interfaces, repositories, entities, migrations and APIs.

No orchestration between domains yet.