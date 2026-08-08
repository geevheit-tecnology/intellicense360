SPRINT 006 — Fleet Assets & Equipment Foundation

Read first:

docs/00_PRODUCT_MANIFESTO.md
docs/01_MASTER_PROMPT.md
docs/ARCHITECTURE.md
docs/BACKEND_ARCHITECTURE.md

IMPORTANT

Do NOT connect this module with Maintenance.

Do NOT connect this module with Checklist.

Do NOT connect this module with Tires.

Do NOT connect this module with Fuel.

Do NOT connect this module with CIOT.

Do NOT connect this module with Intelligence.

Do NOT create cross-module dependencies.

This module must be completely independent.

--------------------------------------------------

OBJECTIVE

Create the Fleet Assets module.

This module is responsible for managing all assets owned or operated by the company.

It must support trucks, trailers, semi-trailers, implements and operational equipment.

No business intelligence.

No recommendations.

No automatic workflows.

--------------------------------------------------

MODULE

internal/modules/assets

Architecture:

domain/
application/
infrastructure/
transport/
ports/

--------------------------------------------------

ENTITIES

Asset

AssetCategory

AssetType

AssetStatus

Equipment

EquipmentCategory

EquipmentType

Implement

ImplementType

Manufacturer

Model

Ownership

Depreciation

Location

--------------------------------------------------

SUPPORTED ASSETS

Truck

Trailer

Semi Trailer

Dolly

Bitrem

Rodotrem

Container Chassis

Tank

Refrigerated Trailer

Platform Trailer

Dump Trailer

Van

Forklift

Crane

Generator

Compressor

Other Equipment

--------------------------------------------------

ASSET LIFECYCLE

Draft

Available

Assigned

In Operation

Maintenance

Inactive

Sold

Disposed

--------------------------------------------------

DATABASE

Create migrations.

UUID primary keys.

TenantId.

AuditId.

Soft Delete.

Optimistic Lock.

Indexes.

CreatedAt.

UpdatedAt.

DeletedAt.

--------------------------------------------------

VALIDATIONS

Unique internal code.

Unique serial number when applicable.

Unique asset tag.

Tenant isolation.

Status validation.

--------------------------------------------------

API

CRUD Assets

CRUD Categories

CRUD Types

CRUD Manufacturers

CRUD Models

CRUD Equipment

Search

Pagination

Sorting

Filtering

--------------------------------------------------

IMPORT

Prepare infrastructure for future CSV and Excel import.

Do not implement import yet.

--------------------------------------------------

EXPORT

Prepare infrastructure for future PDF and Excel export.

Do not implement export yet.

--------------------------------------------------

FILES

Prepare attachment infrastructure.

Support:

Images

PDF

Documents

Warranty files

Invoices

Manuals

Do not connect with Storage Service yet.

--------------------------------------------------

AUDIT

All operations must generate audit events.

Create audit contracts only.

--------------------------------------------------

TESTS

Unit Tests

Repository Tests

Validation Tests

API Tests

--------------------------------------------------

DOCUMENTATION

Generate:

Architecture

Entity Diagram

ER Diagram

Sequence Diagram

REST API Documentation

Acceptance Criteria

README

--------------------------------------------------

NON-FUNCTIONAL REQUIREMENTS

DDD

Clean Architecture

SOLID

Repository Pattern

Dependency Injection

OpenAPI

Production Ready

100% documented

Production quality code

--------------------------------------------------

GOAL

Deliver an enterprise-grade Fleet Assets module fully isolated from other domains and ready for future integrations.