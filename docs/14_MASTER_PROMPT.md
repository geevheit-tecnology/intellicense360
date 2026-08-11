Agora vem a Sprint 014 — CIOT Core

Aqui vou fazer uma distinção importante: não vamos integrar RoadCard/Pamcard/e-Fretes ainda.

Primeiro vamos criar o domínio CIOT corretamente:

CIOT
├── Draft
├── Generated
├── Active
├── Suspended
├── Closed
├── Canceled
└── Error

Incluindo especificamente a realidade que você já encontrou no seu projeto:

TAC Agregado, com CIOT que pode permanecer aberto durante o período operacional.

Também vamos preparar:

CIOT;
contratação;
transportador;
TAC;
veículo;
operação;
valores;
pagamentos;
encerramento;
cancelamento;
histórico;
tentativas de integração;
idempotência;
protocolo externo.

Mas sem conectar ao Fleet, Drivers, Financial ou qualquer outro módulo ainda.

Isso é importante porque depois poderemos adaptar o seu conhecimento e problemas atuais de RoadCard/e-Frete/Pamcard sem contaminar o domínio.

SPRINT 014 — CIOT Core Foundation

Read first:

docs/00_PRODUCT_MANIFESTO.md
docs/01_MASTER_PROMPT.md
docs/ARCHITECTURE.md
docs/BACKEND_ARCHITECTURE.md

IMPORTANT

Create ONLY the CIOT Core domain.

DO NOT integrate with RoadCard.

DO NOT integrate with Pamcard.

DO NOT integrate with e-Frete.

DO NOT integrate with any external CIOT provider.

DO NOT connect CIOT with Fleet.

DO NOT connect CIOT with Drivers.

DO NOT connect CIOT with Financial.

DO NOT connect CIOT with Suppliers.

DO NOT connect CIOT with Maintenance.

DO NOT connect CIOT with Tires.

DO NOT connect CIOT with Checklist.

DO NOT connect CIOT with Intelligence.

DO NOT connect CIOT with Event Bus.

Do not publish events.

Do not create cross-module dependencies.

External integrations will be implemented only in a future dedicated integration sprint.

--------------------------------------------------

OBJECTIVE

Create an enterprise-grade CIOT Core domain.

The domain must represent the complete operational lifecycle of a CIOT independently from external providers.

The implementation must support Brazilian transportation operations and specifically TAC Agregado.

Do not implement provider-specific rules inside the domain.

--------------------------------------------------

MODULE

internal/modules/ciot

Architecture:

domain/
application/
infrastructure/
transport/
ports/

--------------------------------------------------

CORE ENTITIES

CIOT

CIOTContract

CIOTCarrier

CIOTTransporter

CIOTOperation

CIOTVehicleReference

CIOTDriverReference

CIOTPayment

CIOTAmount

CIOTStatusHistory

CIOTProviderAttempt

CIOTExternalReference

CIOTDocument

CIOTStatus

CIOTType

CIOTError

--------------------------------------------------

CIOT TYPES

TAC Agregado

TAC Independente

Other

The domain must be extensible for future CIOT modalities.

--------------------------------------------------

STATUS

Draft

Pending

Generated

Active

Suspended

Closed

Canceled

Error

Do not allow arbitrary status changes.

--------------------------------------------------

CIOT LIFECYCLE

Draft
→ Pending

Pending
→ Generated

Generated
→ Active

Active
→ Suspended

Suspended
→ Active

Active
→ Closed

Pending
→ Canceled

Generated
→ Canceled

Active
→ Canceled

Error
→ Pending

Do not allow invalid transitions.

Document the state machine.

--------------------------------------------------

TAC AGREGADO

The domain must explicitly support long-lived TAC Agregado operations.

A TAC Agregado may remain active during its operational period.

Do not assume that every CIOT is opened and closed within a single trip.

Support:

Start Date

Expected End Date

Actual End Date

Operational Period

Contract Reference

Notes

--------------------------------------------------

CIOT CONTRACT

Represent:

Contract Number

Contract Type

Start Date

End Date

Status

Notes

Do not connect with Suppliers or Financial.

Use references/identifiers only.

--------------------------------------------------

CARRIER

Represent the transport operation parties without importing Driver, Fleet or Supplier domains.

Support:

Document

Legal Name

Trade Name

Registration

Contact Reference

Do not duplicate external domain entities.

Use references only.

--------------------------------------------------

TRANSPORTER / TAC

Represent:

Document

Name

Registration

Contract Reference

Contact Reference

Do not import Driver or Supplier modules.

--------------------------------------------------

OPERATION

Represent:

Operation Number

Origin

Destination

Start Date

Expected End Date

Actual End Date

Cargo Reference

Weight

Distance

Notes

Do not connect to Fleet or Tracking.

--------------------------------------------------

VEHICLE REFERENCE

Store only the external/domain reference:

Vehicle ID

License Plate

Vehicle Type

Do not import Fleet entities.

--------------------------------------------------

DRIVER REFERENCE

Store only:

Driver ID

Name Reference

Document Reference

Do not import Drivers module.

--------------------------------------------------

AMOUNTS

Represent:

Freight Amount

Advance Amount

Balance Amount

Toll Amount

Other Amount

Total Amount

Currency

Do not implement Financial integration.

Use value objects where appropriate.

--------------------------------------------------

PAYMENTS

Represent:

Payment ID

Payment Type

Amount

Due Date

Payment Date

Status

Reference

Notes

Do not connect to Financial.

--------------------------------------------------

EXTERNAL PROVIDER

Create provider-agnostic structures.

CIOTProviderAttempt:

Provider

Attempt Number

Requested At

Completed At

Status

Request Reference

Response Reference

Error Code

Error Message

HTTP Status when applicable

Latency

Do not store credentials.

Do not implement external HTTP clients.

--------------------------------------------------

IDEMPOTENCY

Prepare explicit idempotency support.

A CIOT generation request must be safely retryable.

Support:

Idempotency Key

Request Hash

Attempt Number

Result Reference

Do not implement external provider integration.

--------------------------------------------------

EXTERNAL REFERENCE

Represent:

Provider

External CIOT Number

External Protocol

External Status

Generated At

Last Synchronized At

Do not synchronize externally yet.

--------------------------------------------------

ERROR MODEL

Create structured CIOT errors.

Examples:

InvalidTransition

InvalidDocument

InvalidOperationPeriod

DuplicateRequest

ProviderError

ValidationError

ExternalReferenceConflict

AlreadyClosed

AlreadyCanceled

Do not expose internal implementation details through HTTP.

--------------------------------------------------

HISTORY

CIOT lifecycle history must be append-only.

Track:

Created

Submitted

Generated

Activated

Suspended

Reactivated

Closed

Canceled

ProviderAttempted

ProviderSucceeded

ProviderFailed

PaymentRecorded

ErrorOccurred

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

Tenant-aware constraints

Immutable status history.

Provider attempts must be retained.

Do not physically delete operational history.

--------------------------------------------------

API

Create:

/api/v1/ciot

/api/v1/ciot/{id}

/api/v1/ciot/{id}/history

/api/v1/ciot/{id}/activate

/api/v1/ciot/{id}/suspend

/api/v1/ciot/{id}/reactivate

/api/v1/ciot/{id}/close

/api/v1/ciot/{id}/cancel

/api/v1/ciot/{id}/payments

/api/v1/ciot/{id}/provider-attempts

/api/v1/ciot/{id}/external-reference

/api/v1/ciot/types

Support:

GET

POST

PATCH only where appropriate

Pagination

Filtering

Sorting

Search

--------------------------------------------------

IMPORTANT API RULE

Do NOT expose arbitrary status mutation.

Use explicit application commands:

CreateCIOT

SubmitCIOT

ActivateCIOT

SuspendCIOT

ReactivateCIOT

CloseCIOT

CancelCIOT

RecordPayment

RecordProviderAttempt

RegisterExternalReference

--------------------------------------------------

RBAC

Create permission:

ciot.ciot.manage

Follow existing Identity/RBAC architecture.

--------------------------------------------------

AUDIT

Prepare contracts for:

CIOT Created

CIOT Submitted

CIOT Generated

CIOT Activated

CIOT Suspended

CIOT Reactivated

CIOT Closed

CIOT Canceled

CIOT Payment Recorded

CIOT Provider Attempted

CIOT Provider Succeeded

CIOT Provider Failed

Do NOT implement Event Bus.

--------------------------------------------------

TESTS

Unit Tests

Domain Tests

State Machine Tests

TAC Agregado Lifecycle Tests

Idempotency Tests

Provider Attempt Tests

Historical Immutability Tests

Repository Tests

Tenant Isolation Tests

API Tests

RBAC Tests

Validation Tests

--------------------------------------------------

OPENAPI

Update the existing OpenAPI specification.

Do not create a second OpenAPI specification.

--------------------------------------------------

DOCUMENTATION

Create:

docs/CIOT_CORE.md

Document:

Purpose

Architecture

CIOT Types

TAC Agregado

Lifecycle

State Machine

Entities

Provider Abstraction

Idempotency

Error Model

API

Database

Tenant Isolation

Audit

Testing

Future Integration Strategy

--------------------------------------------------

IMPORTANT

Do NOT implement:

RoadCard API

Pamcard API

e-Frete API

External authentication

External credentials

HTTP provider clients

Webhook handling

Automatic synchronization

Financial integration

Fleet integration

Driver integration

These belong to future sprints.

--------------------------------------------------

QUALITY

DDD

Clean Architecture

SOLID

Repository Pattern

Dependency Injection

Explicit domain commands

Immutable history

Provider-agnostic design

Idempotent operations

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

create CIOT

get CIOT

list CIOTs

invalid lifecycle transition

TAC Agregado activation

suspend

reactivate

close

cancel

record payment

provider attempt

external reference

history

RBAC

tenant isolation

idempotency

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

TAC Agregado behavior

Idempotency implementation

Tests executed

Smoke tests executed

Architecture decisions

Problems found

Do NOT start Sprint 015.

Wait for approval.