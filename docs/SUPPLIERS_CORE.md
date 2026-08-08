# Suppliers Core

Sprint 009 creates Suppliers Core as an isolated business domain.

Suppliers does not connect to Inventory, Maintenance, Fleet, Tires, Fuel, Checklist, Financial, CIOT or Intelligence. Future modules may reference suppliers through contracts only.

## Purpose

Manage companies and service providers that may participate in transportation operations.

## Architecture

- `domain`: supplier entities and enums.
- `application`: validations and use cases.
- `ports`: services, repositories and audit contracts.
- `infrastructure`: in-memory repositories.
- `transport`: DTOs, mappers and HTTP handlers.

## Entities

- Supplier
- SupplierCategory
- SupplierType
- SupplierStatus
- SupplierContact
- SupplierAddress
- SupplierBankAccount
- SupplierDocument
- SupplierRating
- SupplierContract
- SupplierRepresentative

## Business Rules

- Legal name is required.
- CNPJ must be valid when provided.
- CNPJ is unique per tenant.
- Email must be valid when provided.
- UF must be valid when provided.
- Supplier status and supplier type must be valid.
- Archived suppliers cannot transition back to another status.

## API

Base path: `/api/v1/suppliers`

- `GET /`
- `POST /`
- `GET /{id}`
- `PUT /{id}`
- `DELETE /{id}`
- `GET /categories`
- `POST /categories`
- `GET /types`
- `POST /types`
- `GET /contacts`
- `POST /contacts`
- `GET /addresses`
- `POST /addresses`
- `GET /documents`
- `POST /documents`
- `GET /contracts`
- `POST /contracts`
- `GET /ratings`
- `POST /ratings`

All list endpoints support search, pagination, filters and sort parameters following the backend convention.

## Database

Migration: `database/migrations/000011_suppliers_core.sql`

Includes UUID primary keys, `tenant_id`, `audit_id`, soft delete, optimistic version, timestamps, indexes and tenant-aware uniqueness.

## Security

RBAC permission: `suppliers.suppliers.manage`

All routes require authentication, tenant middleware and permission enforcement.

## Audit

Audit contracts are prepared through `ports.AuditRecorder` for:

- Supplier Created
- Supplier Updated
- Supplier Deleted
- Supplier Status Changed
- Supplier Document Added
- Supplier Contract Added

No Event Bus is implemented and no events are published.

## Acceptance Criteria

- Suppliers is independent from all other business modules.
- CRUD exists for suppliers, categories, types, contacts, addresses, documents, contracts and ratings.
- Tenant isolation, soft delete, CNPJ uniqueness and validation rules are covered by tests.
- OpenAPI and documentation are updated.
