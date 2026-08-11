Sprint 012 — Checklist Engine

A regra será:

primeiro criar o novo domínio isolado; depois migramos/aproveitamos o seu Checklist existente.

Ele precisa suportar:

checklist configurável;
tipos de checklist;
itens;
respostas;
evidências fotográficas;
não conformidades;
severidade;
assinatura/identificação do responsável;
histórico;
conclusão;
bloqueio operacional preparado;
templates;
versões do checklist.

Mas não conectar ainda com Fleet, Tires, Maintenance ou Drivers.
SPRINT 012 — Checklist Engine Core

Read first:

docs/00_PRODUCT_MANIFESTO.md
docs/01_MASTER_PROMPT.md
docs/ARCHITECTURE.md
docs/BACKEND_ARCHITECTURE.md

IMPORTANT

This sprint creates ONLY the Checklist Core domain.

DO NOT connect Checklist with Fleet.

DO NOT connect Checklist with Drivers.

DO NOT connect Checklist with Maintenance.

DO NOT connect Checklist with Tires.

DO NOT connect Checklist with Inventory.

DO NOT connect Checklist with Fuel.

DO NOT connect Checklist with Suppliers.

DO NOT connect Checklist with Financial.

DO NOT connect Checklist with CIOT.

DO NOT connect Checklist with Intelligence.

DO NOT connect Checklist with Event Bus.

Do not publish events.

Do not create cross-module dependencies.

The Checklist module must remain independently testable.

IMPORTANT:

The project already has an existing Checklist application/module.

Do NOT delete it.

Do NOT overwrite existing business logic blindly.

First inspect the existing implementation and preserve compatible functionality.

The new domain must be designed so the existing Checklist implementation can later be migrated into it safely.

--------------------------------------------------

OBJECTIVE

Create an enterprise-grade Checklist Engine.

This is NOT merely a checklist CRUD.

The domain must support configurable checklist templates, versioning, execution, responses, evidence, non-conformities and historical traceability.

--------------------------------------------------

MODULE

internal/modules/checklist

Architecture:

domain/
application/
infrastructure/
transport/
ports/

--------------------------------------------------

CORE ENTITIES

ChecklistTemplate

ChecklistTemplateVersion

ChecklistType

ChecklistSection

ChecklistItem

ChecklistItemOption

ChecklistExecution

ChecklistResponse

ChecklistEvidence

ChecklistNonConformity

ChecklistSignature

ChecklistAssignment

ChecklistStatus

ChecklistSeverity

ChecklistHistory

--------------------------------------------------

CHECKLIST TYPES

Prepare support for:

Pre Trip

Post Trip

Daily

Weekly

Monthly

Inspection

Safety

Operational

Regulatory

Custom

Do not hardcode tenant-specific business rules.

--------------------------------------------------

TEMPLATE

A template must support:

Name

Description

Type

Status

Version

Sections

Items

Required Items

Optional Items

Instructions

Scoring Configuration

Severity Configuration

Evidence Requirement

Signature Requirement

Active/Inactive

--------------------------------------------------

VERSIONING

Checklist templates must be versioned.

A published version must be immutable.

Changes must create a new version.

Historical executions must always reference the exact template version used.

Do NOT allow modification of a published version.

--------------------------------------------------

SECTIONS

A checklist may contain multiple sections.

Examples:

Safety

Documentation

Mechanical

Electrical

Tires

Equipment

Body

Operational

Custom

The section structure must remain generic.

Do NOT import Tire or Maintenance domains.

--------------------------------------------------

ITEMS

Each item must support:

Question

Description

Type

Required

Order

Section

Help Text

Severity

Evidence Required

Response Options

Active

--------------------------------------------------

RESPONSE TYPES

Yes / No

Pass / Fail

OK / Not OK

Numeric

Text

Selection

Multi Selection

Date

Signature

Photo

Do not hardcode all future response types into business logic.

Use an extensible domain model.

--------------------------------------------------

EXECUTION

ChecklistExecution must represent an actual checklist performed by a user.

Support:

Execution ID

Template Version

Start Time

End Time

Status

Performed By

Location Reference

Notes

Score when applicable

Final Result

--------------------------------------------------

EXECUTION STATUS

Draft

In Progress

Completed

Canceled

Invalidated

--------------------------------------------------

RESPONSES

Every execution may contain multiple responses.

Each response must preserve:

Checklist Item

Value

Result

Notes

Timestamp

Responder

Severity

--------------------------------------------------

NON-CONFORMITY

Represent operational findings without connecting them to another module.

Support:

Title

Description

Severity

Status

Response Reference

Evidence

Recommendation

Created At

Resolved At

Do NOT automatically create Maintenance Orders.

Do NOT automatically block Fleet assets.

Those integrations belong to future workflow sprints.

--------------------------------------------------

SEVERITY

Info

Low

Medium

High

Critical

Do not hardcode tenant-specific thresholds.

--------------------------------------------------

EVIDENCE

Support references for:

Photo

Video

Document

Signature

Do not integrate external storage yet.

The domain should store metadata/reference only.

--------------------------------------------------

SIGNATURE

Prepare signature structure:

Signer

Timestamp

Signature Reference

Signature Type

Do not implement biometric functionality.

--------------------------------------------------

HISTORY

Checklist execution history must be append-only.

Track:

Created

Started

Response Recorded

Evidence Added

Non-Conformity Created

Completed

Canceled

Invalidated

Do not physically delete execution history.

--------------------------------------------------

ASSIGNMENT

Prepare a generic assignment structure.

Do NOT depend on Driver or Fleet entities.

Use references/identifiers only.

--------------------------------------------------

DATABASE

Create a versioned migration.

Requirements:

UUID

TenantId

AuditId

Version

CreatedAt

UpdatedAt

DeletedAt

Indexes

Tenant-aware constraints

Published template versions must be immutable.

Historical execution records must not be physically deleted.

--------------------------------------------------

BUSINESS RULES

A checklist cannot be completed if required items are unanswered.

A published template version cannot be edited.

An execution must reference exactly one template version.

Responses must reference valid checklist items from the execution's template version.

An execution cannot transition from a terminal state arbitrarily.

Evidence requirements must be respected.

Required signatures must be respected before completion.

Invalid state transitions must return domain errors.

--------------------------------------------------

LIFECYCLE

Template:

Draft
→ Published
→ Archived

Execution:

Draft
→ In Progress
→ Completed

In Progress
→ Canceled

Completed
→ Invalidated

Do not allow arbitrary transitions.

--------------------------------------------------

API

Create:

/api/v1/checklists/templates

/api/v1/checklists/templates/{id}

/api/v1/checklists/templates/{id}/versions

/api/v1/checklists/types

/api/v1/checklists/sections

/api/v1/checklists/items

/api/v1/checklists/executions

/api/v1/checklists/executions/{id}

/api/v1/checklists/executions/{id}/responses

/api/v1/checklists/executions/{id}/evidence

/api/v1/checklists/executions/{id}/non-conformities

/api/v1/checklists/executions/{id}/history

Support existing project conventions for:

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

Use explicit application commands:

CreateChecklistTemplate

CreateTemplateVersion

PublishTemplateVersion

ArchiveTemplate

StartChecklistExecution

RecordChecklistResponse

AddChecklistEvidence

CreateNonConformity

CompleteChecklistExecution

CancelChecklistExecution

InvalidateChecklistExecution

--------------------------------------------------

RBAC

Create permission:

checklist.checklist.manage

Follow existing Identity/RBAC architecture.

--------------------------------------------------

AUDIT

Prepare audit contracts for:

Checklist Template Created

Checklist Template Published

Checklist Template Archived

Checklist Execution Started

Checklist Response Recorded

Checklist Evidence Added

Checklist Non-Conformity Created

Checklist Execution Completed

Checklist Execution Canceled

Checklist Execution Invalidated

Do NOT implement Event Bus.

--------------------------------------------------

TESTS

Unit Tests

Domain Tests

Template Version Tests

Published Version Immutability Tests

Execution Lifecycle Tests

Required Item Validation Tests

Required Evidence Tests

Required Signature Tests

Repository Tests

Tenant Isolation Tests

API Tests

RBAC Tests

Historical Immutability Tests

--------------------------------------------------

OPENAPI

Update the existing OpenAPI specification.

Do not create a second OpenAPI specification.

Follow the existing project convention.

--------------------------------------------------

DOCUMENTATION

Create:

docs/CHECKLIST_CORE.md

Document:

Purpose

Architecture

Existing Checklist Compatibility

Domain Model

Template Versioning

Execution Lifecycle

Business Rules

API

Database

Tenant Isolation

Audit

Migration Strategy

Testing

Acceptance Criteria

--------------------------------------------------

MIGRATION STRATEGY

Because the project already contains Checklist functionality:

1. Inspect the existing Checklist implementation.

2. Identify reusable domain concepts.

3. Do not delete existing functionality.

4. Do not break existing endpoints.

5. Document differences between existing and new architecture.

6. Prepare a future migration path.

Do NOT perform the complete migration in this sprint.

--------------------------------------------------

QUALITY

DDD

Clean Architecture

SOLID

Repository Pattern

Dependency Injection

Explicit domain commands

Immutable published templates

Immutable execution history

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

create checklist type

create template

create template version

add section

add required item

publish template version

start execution

record response

attempt completion with unanswered required item

complete execution after requirements are satisfied

create non-conformity

query execution history

RBAC

tenant isolation

invalid lifecycle transition

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

Template versioning rules

Compatibility findings with existing Checklist

Tests executed

Smoke tests executed

Architecture decisions

Migration recommendations

Problems found

Do NOT start Sprint 013.

Wait for approval.

E aqui está a decisão estratégica

O seu Checklist atual não será jogado fora.

Vamos fazer:

Checklist atual
      ↓
análise
      ↓
compatibilidade
      ↓
migração gradual
      ↓
Checklist Engine

Assim preservamos o que você já tem e, ao mesmo tempo, transformamos o checklist em uma das fontes de dados mais importantes do futuro Geevheit Intelligence 360°.