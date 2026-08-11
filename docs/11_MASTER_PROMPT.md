O objetivo não é apenas registrar abastecimentos. Precisamos guardar uma base capaz de futuramente responder:

Quanto cada veículo consome?
Quanto custa cada km?
Onde existem desvios?
Qual posto é mais eficiente?
Qual combustível está gerando maior custo?

Mas essas análises ficam para o Intelligence Engine.

Envie ao Codex:

SPRINT 011 — Fuel Core

Read first:

docs/00_PRODUCT_MANIFESTO.md
docs/01_MASTER_PROMPT.md
docs/ARCHITECTURE.md
docs/BACKEND_ARCHITECTURE.md

IMPORTANT

This sprint creates ONLY the Fuel Core domain.

DO NOT connect Fuel with Fleet.

DO NOT connect Fuel with Drivers.

DO NOT connect Fuel with Maintenance.

DO NOT connect Fuel with Inventory.

DO NOT connect Fuel with Tires.

DO NOT connect Fuel with Checklist.

DO NOT connect Fuel with Suppliers.

DO NOT connect Fuel with Financial.

DO NOT connect Fuel with Intelligence.

DO NOT connect Fuel with Event Bus.

Do not publish events.

Do not create cross-module dependencies.

Fuel must remain independently testable.

--------------------------------------------------

OBJECTIVE

Create an enterprise-grade Fuel Core domain.

The module must record fuel operations and preserve the complete history of fuel transactions.

This is NOT an intelligence module.

This is NOT a cost-analysis module.

This is NOT a fleet integration.

The domain must provide reliable operational data for future intelligence.

--------------------------------------------------

MODULE

internal/modules/fuel

Architecture:

domain/
application/
infrastructure/
transport/
ports/

--------------------------------------------------

ENTITIES

FuelTransaction

FuelType

FuelStation

FuelStationLocation

FuelTank

FuelNozzle

FuelReading

FuelPrice

FuelReceipt

FuelAttachment

FuelAdjustment

FuelTransactionStatus

--------------------------------------------------

FUEL TYPES

Diesel S10

Diesel S500

Gasoline

Ethanol

GNV

ARLA 32

Other

Do not hardcode business rules that may vary by tenant.

--------------------------------------------------

TRANSACTION

Represent:

Transaction ID

Date/Time

Fuel Type

Quantity

Unit Price

Total Amount

Odometer Reading

Engine Hour Reading when applicable

Station

Pump/Nozzle when applicable

Receipt Number

Driver Reference

Vehicle/Asset Reference

Payment Method

Notes

Status

The references to Driver and Vehicle/Asset must remain identifiers/contracts only.

Do NOT import or depend on Fleet or Driver modules.

--------------------------------------------------

STATUS

Draft

Completed

Canceled

Adjusted

Rejected

--------------------------------------------------

FUEL STATION

Represent:

Name

Legal Name

CNPJ

Address

City

State

Country

Latitude

Longitude

Active

Notes

CNPJ must be unique per tenant when provided.

--------------------------------------------------

FUEL TANK

Represent internal company fuel tanks.

Capacity

Current Reading

Fuel Type

Location Reference

Status

Notes

Do not implement automatic stock synchronization with Inventory.

--------------------------------------------------

FUEL NOZZLE

Represent:

Code

Fuel Type

Tank Reference

Status

Meter Reading

--------------------------------------------------

READINGS

Create immutable fuel readings.

Support:

Odometer

Engine Hours

Fuel Tank Reading

Reading Date

Source

Notes

Do not integrate telemetry yet.

--------------------------------------------------

FUEL PRICE

Represent:

Fuel Type

Unit Price

Effective Date

Station Reference

Source

Notes

Do not implement automatic price synchronization.

--------------------------------------------------

RECEIPT

Represent:

Receipt Number

Date

Amount

Attachment Reference

Notes

Do not integrate external storage.

--------------------------------------------------

ADJUSTMENT

Represent corrections to operational records.

Adjustment Type

Reason

Original Reference

Adjusted Value

Notes

Do not allow silent mutation of historical completed transactions.

--------------------------------------------------

BUSINESS RULES

Completed fuel transactions must preserve historical values.

Do not allow arbitrary modification of completed transactions.

Corrections must use an explicit adjustment mechanism.

Quantity must be greater than zero.

Unit price must be zero or greater.

Total amount must be valid.

Odometer must not be negative.

Engine hours must not be negative.

Fuel type must be valid.

Tenant isolation is mandatory.

--------------------------------------------------

DATABASE

Create versioned migration.

Requirements:

UUID

TenantId

AuditId

Version

CreatedAt

UpdatedAt

DeletedAt

Indexes

Tenant-aware unique constraints

Historical transaction records must not be physically deleted.

--------------------------------------------------

AUDIT

Prepare audit contracts for:

Fuel Transaction Created

Fuel Transaction Completed

Fuel Transaction Canceled

Fuel Transaction Adjusted

Fuel Station Created

Fuel Reading Recorded

Fuel Price Recorded

Do NOT implement Event Bus.

Do NOT publish events.

--------------------------------------------------

API

Create:

/api/v1/fuel/transactions

/api/v1/fuel/types

/api/v1/fuel/stations

/api/v1/fuel/tanks

/api/v1/fuel/nozzles

/api/v1/fuel/readings

/api/v1/fuel/prices

/api/v1/fuel/receipts

/api/v1/fuel/adjustments

Support existing project conventions:

GET

POST

PATCH/PUT where appropriate

DELETE only where business rules permit

Pagination

Filtering

Sorting

Search

--------------------------------------------------

COMMANDS

Do not expose arbitrary status mutation.

Use explicit application commands where lifecycle matters:

CreateFuelTransaction

CompleteFuelTransaction

CancelFuelTransaction

AdjustFuelTransaction

RecordFuelReading

RecordFuelPrice

--------------------------------------------------

RBAC

Create permission:

fuel.fuel.manage

Follow the existing Identity/RBAC implementation.

--------------------------------------------------

TESTS

Unit Tests

Domain Tests

Transaction Lifecycle Tests

Repository Tests

Tenant Isolation Tests

API Tests

RBAC Tests

Validation Tests

Historical Immutability Tests

Adjustment Tests

--------------------------------------------------

OPENAPI

Update the existing OpenAPI specification.

Do not create another OpenAPI specification.

Follow the current project convention.

--------------------------------------------------

DOCUMENTATION

Create:

docs/FUEL_CORE.md

Document:

Purpose

Architecture

Domain Model

Entities

Transaction Lifecycle

Business Rules

API

Database

Tenant Isolation

Audit

Testing

Acceptance Criteria

--------------------------------------------------

QUALITY

DDD

Clean Architecture

SOLID

Repository Pattern

Dependency Injection

Explicit domain commands

Historical immutability

Production-ready code

Do not refactor unrelated modules.

Do not break existing endpoints.

--------------------------------------------------

VALIDATION

Run:

go test ./...

Start API using configurable PORT.

Perform real HTTP smoke tests for:

login

create fuel type

create station

create fuel transaction

complete transaction

record reading

record price

query transaction

query history where implemented

attempt invalid completed-transaction mutation

adjust transaction

RBAC

tenant isolation

Verify configurable PORT is respected.

--------------------------------------------------

FINAL REPORT

Report:

Files created

Files modified

Migration created

Endpoints created

Permissions created

Domain entities

Lifecycle rules

Tests executed

Smoke tests executed

Architecture decisions

Problems found

Do NOT start Sprint 012.

Wait for approval.