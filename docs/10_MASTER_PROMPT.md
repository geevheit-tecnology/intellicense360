SPRINT 010 — Tire Lifecycle Core

Read first:

docs/00_PRODUCT_MANIFESTO.md
docs/01_MASTER_PROMPT.md
docs/ARCHITECTURE.md
docs/BACKEND_ARCHITECTURE.md

IMPORTANT

This sprint creates ONLY the Tire Lifecycle Core.

DO NOT connect Tires with Fleet.

DO NOT connect Tires with Maintenance.

DO NOT connect Tires with Inventory.

DO NOT connect Tires with Checklist.

DO NOT connect Tires with Fuel.

DO NOT connect Tires with Suppliers.

DO NOT connect Tires with Financial.

DO NOT connect Tires with Intelligence.

DO NOT connect Tires with Event Bus.

Do not publish events.

Do not create cross-module dependencies.

The Tire module must remain independently testable.

--------------------------------------------------

OBJECTIVE

Create an enterprise-grade Tire Lifecycle domain.

This is NOT a simple tire CRUD.

The system must represent the complete lifecycle of each physical tire from acquisition to final disposal.

Every tire must have traceable history.

No history may be physically deleted.

--------------------------------------------------

MODULE

internal/modules/tires

Architecture:

domain/
application/
infrastructure/
transport/
ports/

--------------------------------------------------

CORE ENTITIES

Tire

TireBrand

TireModel

TireSpecification

TirePosition

TireInstallation

TireRemoval

TireMovement

TireInspection

TireMeasurement

TireRetread

TireRetreadEvent

TireStatus

TireCondition

TireHistory

TireCost

TireDisposal

TireAttachment

--------------------------------------------------

TIRE IDENTIFICATION

Support:

Serial Number

DOT

Brand

Model

Dimension

Construction

Load Index

Speed Rating

Original Tread Depth

Current Tread Depth

Purchase Date

Purchase Cost

Supplier Reference

Warranty

Manufacturing Date

--------------------------------------------------

TIRE STATUS

New

In Stock

Installed

Removed

Under Inspection

Under Retread

Retreaded

Reserved

Damaged

End Of Life

Disposed

Lost

--------------------------------------------------

TIRE CONDITION

Excellent

Good

Attention

Heavy Wear

At Limit

Critical

Damaged

End Of Life

--------------------------------------------------

PHYSICAL POSITION

Create a generic position model.

Do NOT depend on Fleet Vehicle entities.

A tire position must support:

Axle

Side

Position

Inner/Outer

Position Code

Position Label

This allows future integration with trucks, trailers and other assets without coupling this module to Fleet.

--------------------------------------------------

INSTALLATION

Represent:

Installation Date

Installation Position

Initial KM

Initial Tread Depth

Installation Reason

Notes

--------------------------------------------------

REMOVAL

Represent:

Removal Date

Removal Position

Removal KM

Remaining Tread Depth

Removal Reason

Condition

Notes

--------------------------------------------------

MOVEMENTS

Represent complete physical movement history.

Types:

Purchase

Receipt

Installation

Removal

Transfer

Retread

Return

Disposal

Adjustment

Loss

Every movement must be immutable from the domain perspective.

--------------------------------------------------

INSPECTION

Represent tire inspections.

Data:

Inspection Date

Tread Depth

Pressure

Condition

Visual Condition

Sidewall Condition

Shoulder Condition

Crown Condition

Irregular Wear

Damage

Cracks

Cuts

Object Penetration

Recommendation

Notes

--------------------------------------------------

MEASUREMENTS

Support multiple measurements.

Tread depth in millimeters.

Pressure.

Measurement position.

Measurement date.

Do not hardcode a single threshold.

Thresholds must be configurable in future.

--------------------------------------------------

RETREAD

Represent complete retread lifecycle.

Retread Number

Retread Date

Provider Reference

Tread Brand

Tread Model

Cost

Before Retread Tread Depth

After Retread Tread Depth

Warranty

Result

Status

Notes

A tire may have zero, one or multiple retread events.

--------------------------------------------------

COST

Prepare tire cost structure.

Track:

Purchase Cost

Retread Cost

Repair Cost

Other Cost

Total Cost

Cost Per KM

Do not implement automatic financial calculations that depend on Fleet or telemetry yet.

Create domain contracts/interfaces where future calculations can be added.

--------------------------------------------------

HISTORY

Tire history must be append-only from the business perspective.

Every relevant lifecycle event must be traceable.

Examples:

Created

Purchased

Received

Installed

Inspected

Moved

Removed

Retreaded

Reinstalled

Damaged

Disposed

Do not physically delete history.

--------------------------------------------------

ATTACHMENTS

Prepare contracts for:

Photos

Invoices

Inspection Photos

Retread Documents

Warranty Documents

Disposal Documents

Do not integrate external storage yet.

--------------------------------------------------

DATABASE

Create a versioned migration.

Requirements:

UUID

TenantId

AuditId

Soft Delete only for master entities where appropriate

Version

CreatedAt

UpdatedAt

DeletedAt

Indexes

Tenant-aware unique constraints

Important:

Do NOT soft-delete immutable lifecycle history records.

--------------------------------------------------

VALIDATIONS

Tenant isolation.

Unique tire serial number per tenant when provided.

Valid tire status.

Valid condition.

Valid tread depth.

Valid pressure.

Valid dimension.

Valid lifecycle transition.

Prevent invalid transitions.

Examples:

Disposed -> Installed = invalid

Disposed -> Retread = invalid

Lost -> Installed = invalid

Installed -> Installed = invalid without a removal/reinstallation lifecycle.

--------------------------------------------------

LIFECYCLE STATE MACHINE

Implement explicit domain transition validation.

Example:

New
→ In Stock

In Stock
→ Installed

Installed
→ Removed

Removed
→ In Stock

Removed
→ Under Retread

Under Retread
→ Retreaded

Retreaded
→ In Stock

Installed
→ Damaged

Damaged
→ End Of Life

End Of Life
→ Disposed

Do not allow arbitrary status changes.

--------------------------------------------------

API

Create REST APIs:

/api/v1/tires

/api/v1/tires/brands

/api/v1/tires/models

/api/v1/tires/specifications

/api/v1/tires/positions

/api/v1/tires/inspections

/api/v1/tires/measurements

/api/v1/tires/installations

/api/v1/tires/removals

/api/v1/tires/movements

/api/v1/tires/retreads

/api/v1/tires/history

/api/v1/tires/attachments

Support existing project conventions for:

GET

POST

PATCH/PUT

DELETE where appropriate

Pagination

Filtering

Sorting

Search

--------------------------------------------------

IMPORTANT API RULE

Do not expose mutation endpoints that allow arbitrary lifecycle status changes.

Lifecycle transitions must go through explicit application/domain commands.

Examples:

InstallTire

RemoveTire

MoveTire

InspectTire

SendToRetread

CompleteRetread

DisposeTire

Do not implement integrations with other modules.

--------------------------------------------------

RBAC

Create permission:

tires.tires.manage

Follow the existing Identity/RBAC architecture.

--------------------------------------------------

AUDIT

Prepare audit contracts for:

Tire Created

Tire Updated

Tire Installed

Tire Removed

Tire Inspected

Tire Moved

Tire Sent To Retread

Tire Retreaded

Tire Damaged

Tire Disposed

Do not implement Event Bus.

--------------------------------------------------

TESTS

Unit Tests

Domain Tests

Lifecycle State Machine Tests

Repository Tests

Tenant Isolation Tests

API Tests

RBAC Tests

Soft Delete Tests where applicable

Immutable History Tests

Invalid Transition Tests

--------------------------------------------------

OPENAPI

Update the existing OpenAPI specification.

Do not create another OpenAPI specification.

Follow the current project convention.

--------------------------------------------------

DOCUMENTATION

Create:

docs/TIRES_CORE.md

Document:

Purpose

Architecture

Domain Model

Lifecycle

State Machine

Entities

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

Immutable lifecycle history

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

create tire brand

create tire model

create tire

get tire

list tires

inspect tire

install tire

remove tire

create retread event

query tire history

RBAC

tenant isolation

invalid lifecycle transition

soft delete where applicable

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

Lifecycle transitions

Tests executed

Smoke tests executed

Architecture decisions

Problems found

Do NOT start Sprint 011.

Wait for approval.