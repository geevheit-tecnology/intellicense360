SPRINT 007 — Maintenance Core Foundation

Read first:

docs/00_PRODUCT_MANIFESTO.md
docs/01_MASTER_PROMPT.md
docs/ARCHITECTURE.md
docs/BACKEND_ARCHITECTURE.md

IMPORTANT

Do NOT connect with Fleet.

Do NOT connect with Tires.

Do NOT connect with Checklist.

Do NOT connect with Fuel.

Do NOT connect with Inventory.

Do NOT connect with Financial.

Do NOT connect with Intelligence.

No cross-domain dependencies.

Maintenance must be completely independent.

--------------------------------------------------

OBJECTIVE

Create an enterprise Maintenance Core module.

This module manages maintenance processes only.

No business intelligence.

No automatic scheduling.

No integrations.

--------------------------------------------------

MODULE

internal/modules/maintenance

Architecture

domain/

application/

infrastructure/

transport/

ports/

--------------------------------------------------

ENTITIES

MaintenanceOrder

MaintenanceType

MaintenanceStatus

MaintenancePriority

MaintenanceCategory

Workshop

Technician

Labor

ServiceItem

Downtime

Availability

MaintenanceSchedule

MaintenanceReason

MaintenanceHistory

--------------------------------------------------

MAINTENANCE TYPES

Preventive

Corrective

Predictive

Inspection

Emergency

Warranty

External

Internal

--------------------------------------------------

STATUS

Draft

Open

Waiting

Approved

Executing

Paused

Completed

Canceled

--------------------------------------------------

DATABASE

Create migrations.

UUID.

TenantId.

AuditId.

Soft Delete.

Optimistic Lock.

Indexes.

--------------------------------------------------

VALIDATIONS

Status transitions.

Priority validation.

Unique maintenance code.

Tenant isolation.

--------------------------------------------------

API

CRUD Maintenance Orders

CRUD Workshops

CRUD Technicians

CRUD Categories

CRUD Types

CRUD Priorities

CRUD Reasons

Search

Pagination

Filtering

Sorting

--------------------------------------------------

FILES

Prepare attachment contracts.

Photos

Invoices

Reports

Warranty

Manuals

No storage integration.

--------------------------------------------------

AUDIT

Generate audit contracts.

--------------------------------------------------

REPORTS

Prepare report interfaces.

Do not implement reports.

--------------------------------------------------

TESTS

Unit

Repository

API

Validation

--------------------------------------------------

DOCUMENTATION

Architecture

ER Diagram

Entity Diagram

REST Documentation

Acceptance Criteria

README

--------------------------------------------------

NON FUNCTIONAL

DDD

Clean Architecture

SOLID

Repository Pattern

Dependency Injection

OpenAPI

Production Ready

100% documented

--------------------------------------------------

GOAL

Enterprise Maintenance Core ready for future integrations.