SPRINT 009 — Suppliers Core

Read first:

docs/00_PRODUCT_MANIFESTO.md
docs/01_MASTER_PROMPT.md
docs/ARCHITECTURE.md
docs/BACKEND_ARCHITECTURE.md

IMPORTANT

This sprint creates ONLY the Suppliers Core domain.

DO NOT connect Suppliers with Inventory.

DO NOT connect Suppliers with Maintenance.

DO NOT connect Suppliers with Fleet.

DO NOT connect Suppliers with Tires.

DO NOT connect Suppliers with Fuel.

DO NOT connect Suppliers with Checklist.

DO NOT connect Suppliers with Financial.

DO NOT connect Suppliers with CIOT.

DO NOT connect Suppliers with Intelligence.

DO NOT publish domain events.

Do not create cross-module dependencies.

--------------------------------------------------

OBJECTIVE

Create an enterprise-grade Suppliers Core module.

The module must manage companies and service providers that may participate in the transportation operation.

The domain must remain independent and reusable by future modules.

--------------------------------------------------

MODULE

internal/modules/suppliers

Architecture:

domain/
application/
infrastructure/
transport/
ports/

--------------------------------------------------

DOMAIN ENTITIES

Supplier

SupplierCategory

SupplierType

SupplierStatus

SupplierContact

SupplierAddress

SupplierBankAccount

SupplierDocument

SupplierRating

SupplierContract

SupplierRepresentative

--------------------------------------------------

SUPPLIER TYPES

Parts Supplier

Tire Supplier

Fuel Supplier

Workshop

Service Provider

Technology Provider

Insurance Provider

Transport Provider

Other

--------------------------------------------------

SUPPLIER STATUS

Draft

Active

Inactive

Blocked

Archived

--------------------------------------------------

SUPPLIER DATA

Legal Name

Trade Name

CNPJ

State Registration

Municipal Registration

Email

Phone

Website

Notes

Status

Category

Type

--------------------------------------------------

CONTACTS

Name

Role

Email

Phone

Mobile

Primary Contact

--------------------------------------------------

ADDRESS

Street

Number

Complement

Neighborhood

City

State

Postal Code

Country

Address Type

--------------------------------------------------

BANK ACCOUNT

Bank

Branch

Account

Account Type

PIX Key

Holder Name

Holder Document

Do not expose sensitive banking data unnecessarily through API responses.

--------------------------------------------------

DOCUMENTS

Document Type

Document Number

Issue Date

Expiration Date

Status

Attachment Reference

Do not integrate external storage yet.

--------------------------------------------------

RATING

Create the domain structure for supplier evaluation.

Criteria:

Quality

Price

Delivery

Service

Reliability

Overall Score

Do not implement automatic scoring.

--------------------------------------------------

CONTRACT

Create the structure for supplier contracts.

Contract Number

Start Date

End Date

Status

Notes

Attachment Reference

Do not implement contract workflows.

--------------------------------------------------

DATABASE

Create a versioned migration.

Requirements:

UUID

TenantId

AuditId

Soft Delete

Optimistic Lock / Version

CreatedAt

UpdatedAt

DeletedAt

Indexes

Tenant-aware unique constraints.

CNPJ must be unique per tenant when provided.

--------------------------------------------------

VALIDATIONS

Tenant isolation.

CNPJ format validation.

Unique CNPJ per tenant.

Required legal name.

Valid status.

Valid supplier type.

Valid category.

Valid email when provided.

Valid state/UF when provided.

Prevent invalid terminal status transitions if applicable.

--------------------------------------------------

API

Create CRUD endpoints:

/api/v1/suppliers

/api/v1/suppliers/categories

/api/v1/suppliers/types

/api/v1/suppliers/contacts

/api/v1/suppliers/addresses

/api/v1/suppliers/documents

/api/v1/suppliers/contracts

/api/v1/suppliers/ratings

Support:

GET

POST

PUT/PATCH according to existing project convention

DELETE using soft delete

Pagination

Filtering

Sorting

Search

--------------------------------------------------

SECURITY

Add appropriate RBAC permission:

suppliers.suppliers.manage

Respect existing authentication.

Respect tenant middleware.

Respect existing authorization conventions.

--------------------------------------------------

AUDIT

Prepare audit contracts for:

Supplier Created

Supplier Updated

Supplier Deleted

Supplier Status Changed

Supplier Document Added

Supplier Contract Added

Do not implement Event Bus.

Do not publish events.

--------------------------------------------------

TESTS

Unit Tests

Domain Validation Tests

Repository Tests

Tenant Isolation Tests

HTTP/API Tests

RBAC Tests

Soft Delete Tests

--------------------------------------------------

OPENAPI

Update the existing OpenAPI specification.

Do not create a second unrelated OpenAPI specification.

Follow the existing project convention.

--------------------------------------------------

DOCUMENTATION

Create:

docs/SUPPLIERS_CORE.md

Document:

Purpose

Architecture

Domain Model

Entities

Business Rules

API

Database

Security

Tenant Isolation

Audit

Testing

Acceptance Criteria

--------------------------------------------------

QUALITY REQUIREMENTS

DDD

Clean Architecture

SOLID

Repository Pattern

Dependency Injection

Production-ready code

No duplicated architecture

Follow existing project conventions.

Do not refactor unrelated modules.

Do not break existing endpoints.

--------------------------------------------------

VALIDATION

Run:

go test ./...

Run the API.

Perform real HTTP smoke tests for:

login

create supplier

get supplier

list suppliers

create category

create contact

tenant isolation

RBAC

soft delete

Verify configurable PORT is respected.

--------------------------------------------------

FINAL REPORT

At the end report:

Files created

Files modified

Migration created

Endpoints created

Permissions created

Tests executed

Smoke tests executed

Architecture decisions

Any problems found

Do NOT start Sprint 010 automatically.

Wait for approval.