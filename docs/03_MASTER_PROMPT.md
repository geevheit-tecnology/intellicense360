SPRINT 003 — Fleet Core Foundation

Read first:

docs/00_PRODUCT_MANIFESTO.md
docs/01_MASTER_PROMPT.md
docs/ARCHITECTURE.md
docs/BACKEND_ARCHITECTURE.md

IMPORTANT

Do NOT connect this module with Checklist.

Do NOT connect this module with Tires.

Do NOT connect this module with Maintenance.

Do NOT connect this module with Fuel.

Do NOT connect this module with CIOT.

Do NOT create cross-module dependencies.

The Fleet module must be completely independent.

--------------------------------------------------

Objective

Create the Fleet Core domain following DDD and Clean Architecture.

The Fleet module represents vehicles and assets only.

No intelligence.

No recommendations.

No automations.

--------------------------------------------------

Create module:

internal/modules/fleet

with:

domain/
application/
infrastructure/
transport/
ports/

--------------------------------------------------

Entities

Vehicle

VehicleCategory

VehicleType

VehicleBrand

VehicleModel

VehicleStatus

OwnershipType

Asset

LicensePlate

Renavam

Chassis

Engine

Color

FuelType

Transmission

AxleConfiguration

EmissionStandard

--------------------------------------------------

Vehicle must support

Truck

Trailer

Semi Trailer

Van

Car

Motorcycle

Machinery

Equipment

--------------------------------------------------

Vehicle lifecycle

Draft

Active

Inactive

Maintenance

Sold

Disposed

--------------------------------------------------

Database

Create migrations.

Use UUID.

Soft Delete.

Optimistic Lock.

TenantId.

AuditId.

CreatedAt.

UpdatedAt.

DeletedAt.

--------------------------------------------------

API

CRUD Vehicle

CRUD Vehicle Brand

CRUD Vehicle Model

CRUD Vehicle Category

CRUD Vehicle Type

CRUD Asset

Search

Pagination

Filtering

Sorting

--------------------------------------------------

Validation

Unique License Plate

Unique Chassis

Unique Renavam

Tenant isolation

--------------------------------------------------

Documentation

Generate:

Domain documentation

Entity diagram

ER diagram

API documentation

Acceptance criteria

--------------------------------------------------

Testing

Unit Tests

Repository Tests

API Tests

Validation Tests

--------------------------------------------------

Do NOT create relationships with any other module.

Fleet Core must be fully independent and production-ready.

Wait for approval before Sprint 004.