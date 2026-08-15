Essa sprint é crítica: não é para criar telas nem conectar tudo de qualquer maneira. Ela cria a infraestrutura de eventos que permitirá que Tires, Fuel, Checklist, Maintenance, Financial, CIOT etc. conversem sem perder o desacoplamento.

SPRINT 015 — EVENT BUS & DOMAIN EVENTS

Read first:

docs/00_PRODUCT_MANIFESTO.md
docs/01_MASTER_PROMPT.md
docs/ARCHITECTURE.md
docs/BACKEND_ARCHITECTURE.md

Also inspect:

backend/api/internal/contracts/events/
backend/api/internal/core/
backend/api/internal/modules/

IMPORTANT

This sprint creates ONLY the Domain Events and Event Bus infrastructure.

DO NOT build Intelligence yet.

DO NOT build dashboards.

DO NOT build Flutter UI.

DO NOT integrate external providers.

DO NOT integrate RoadCard.

DO NOT integrate Pamcard.

DO NOT integrate e-Frete.

DO NOT create business intelligence.

DO NOT create analytics.

DO NOT create cross-module direct dependencies.

The existing modules must remain independently testable.

--------------------------------------------------
OBJECTIVE
--------------------------------------------------

Create an internal, provider-agnostic Domain Event infrastructure.

The objective is to allow modules to publish business events without importing or directly depending on other business modules.

Architecture must support:

Module
   ↓
Domain Event
   ↓
Event Bus
   ↓
Subscribers / Handlers
   ↓
Future Intelligence / Workflows

The Event Bus must NOT contain business rules.

The Event Bus is infrastructure.

--------------------------------------------------
ARCHITECTURE
--------------------------------------------------

Create:

backend/api/internal/events/

Suggested structure:

events/
  domain_event.go
  event_bus.go
  event_handler.go
  event_registry.go
  event_metadata.go
  event_errors.go

events/inmemory/
  bus.go

events/outbox/
  repository.go
  service.go

events/handlers/
  registry.go

Use interfaces so the transport implementation can be replaced later.

Do not hardcode a specific external message broker.

--------------------------------------------------
DOMAIN EVENT CONTRACT
--------------------------------------------------

Every event must contain:

EventID

EventType

EventVersion

OccurredAt

TenantID

AggregateID

AggregateType

CorrelationID

CausationID

ActorID when available

ActorType when available

Metadata

Payload

The event contract must be serializable.

Use JSON-compatible payloads.

--------------------------------------------------
EVENT TYPE
--------------------------------------------------

Use stable event type names.

Examples:

tire.removed.v1

tire.installed.v1

tire.inspected.v1

tire.retreaded.v1

fuel.transaction.completed.v1

fuel.transaction.adjusted.v1

checklist.execution.completed.v1

checklist.non_conformity.created.v1

maintenance.order.created.v1

maintenance.order.completed.v1

financial.expense.created.v1

financial.expense.paid.v1

financial.adjustment.created.v1

ciot.created.v1

ciot.activated.v1

ciot.closed.v1

Do not use Go package names as event identifiers.

--------------------------------------------------
EVENT VERSIONING
--------------------------------------------------

Events must be versioned.

Example:

tire.removed.v1

Future incompatible schema:

tire.removed.v2

Do not silently change the meaning of an existing event version.

Document compatibility rules.

--------------------------------------------------
EVENT BUS
--------------------------------------------------

Create:

type EventBus interface

Support:

Publish

PublishBatch

Subscribe

Unsubscribe if appropriate

The interface must not depend on any specific broker.

--------------------------------------------------
IN-MEMORY BUS
--------------------------------------------------

Implement an in-memory Event Bus for:

local development

unit tests

integration tests

The in-memory implementation must:

- dispatch events to registered handlers
- preserve event metadata
- isolate handler failures
- return meaningful errors
- support multiple handlers
- not panic the application because one handler failed

Do not introduce uncontrolled goroutines.

Prefer deterministic behavior for tests.

--------------------------------------------------
HANDLERS
--------------------------------------------------

Create:

EventHandler interface

Each handler must declare the event types it handles.

Example:

TireRemovedHandler

FuelTransactionCompletedHandler

ChecklistNonConformityHandler

Do NOT implement cross-module business behavior yet.

For this sprint, handlers may be:

registered

invoked

logged/tested

but must not alter other business domains.

--------------------------------------------------
OUTBOX PATTERN
--------------------------------------------------

Implement an Outbox abstraction.

The purpose is future reliable event publication.

Create an interface similar to:

OutboxRepository

Support:

Save

GetPending

MarkPublished

MarkFailed

Retry

The Outbox record must contain:

ID

TenantID

EventID

EventType

AggregateID

AggregateType

Payload

OccurredAt

Status

Attempts

AvailableAt

PublishedAt

LastError

CreatedAt

UpdatedAt

--------------------------------------------------
OUTBOX STATUS

Pending

Processing

Published

Failed

DeadLetter

--------------------------------------------------
IDEMPOTENCY
--------------------------------------------------

Event processing must support idempotency.

The same EventID must not be processed twice by the same logical consumer.

Create an abstraction:

EventConsumerStore

or equivalent.

Support:

HasProcessed

MarkProcessed

Do not require an external database implementation if the project currently uses memory repositories.

Create the contracts and an in-memory implementation.

--------------------------------------------------
RETRY
--------------------------------------------------

Prepare retry handling.

Support:

Attempt count

Last error

Next attempt time

Maximum retry configuration

Do not create an infinite retry loop.

Do not implement distributed scheduling yet.

--------------------------------------------------
DEAD LETTER
--------------------------------------------------

Events that exceed retry limits must be representable as:

DeadLetter

Store:

EventID

EventType

Payload

Failure reason

Attempts

Last error

Timestamp

Do not silently discard events.

--------------------------------------------------
TRANSACTIONAL BOUNDARY
--------------------------------------------------

Document the difference between:

Domain transaction

Outbox persistence

Event publication

Do not pretend the in-memory Event Bus provides distributed transactional guarantees.

--------------------------------------------------
EXISTING EVENTS
--------------------------------------------------

Inspect the existing files under:

backend/api/internal/contracts/events/

Do NOT duplicate event definitions.

Unify the current contracts with the new Domain Event model.

Preserve compatible event names where possible.

If an existing event contract conflicts with the new model:

1. Document the conflict.
2. Prefer backward-compatible adaptation.
3. Do not silently delete existing contracts.

--------------------------------------------------
MODULE EVENT CONTRACTS
--------------------------------------------------

Create or adapt events for:

ASSETS

asset.created.v1
asset.updated.v1

MAINTENANCE

maintenance.order.created.v1
maintenance.order.completed.v1
maintenance.order.canceled.v1

INVENTORY

inventory.part.created.v1
inventory.part.updated.v1

SUPPLIERS

supplier.created.v1
supplier.updated.v1

TIRES

tire.created.v1
tire.installed.v1
tire.removed.v1
tire.inspected.v1
tire.retreaded.v1
tire.damaged.v1
tire.disposed.v1

FUEL

fuel.transaction.created.v1
fuel.transaction.completed.v1
fuel.transaction.canceled.v1
fuel.transaction.adjusted.v1

CHECKLIST

checklist.template.published.v1
checklist.execution.started.v1
checklist.execution.completed.v1
checklist.non_conformity.created.v1

FINANCIAL

financial.expense.created.v1
financial.expense.approved.v1
financial.expense.paid.v1
financial.revenue.created.v1
financial.revenue.received.v1
financial.adjustment.created.v1

CIOT

ciot.created.v1
ciot.activated.v1
ciot.suspended.v1
ciot.reactivated.v1
ciot.closed.v1
ciot.canceled.v1

Do not require every module to publish these events automatically in this sprint.

The objective is to establish standardized contracts.

--------------------------------------------------
EVENT ENVELOPE
--------------------------------------------------

Example conceptual structure:

{
  "event_id": "...",
  "event_type": "tire.removed.v1",
  "event_version": 1,
  "occurred_at": "...",
  "tenant_id": "...",
  "aggregate_id": "...",
  "aggregate_type": "tire",
  "correlation_id": "...",
  "causation_id": "...",
  "actor_id": "...",
  "metadata": {},
  "payload": {}
}

Do not expose internal database implementation details in event payloads.

--------------------------------------------------
SECURITY
--------------------------------------------------

Never include:

passwords

JWT tokens

secrets

API keys

payment credentials

provider credentials

sensitive authentication information

in event payloads.

TenantID must always be explicit.

--------------------------------------------------
OBSERVABILITY
--------------------------------------------------

Prepare structured logging for:

event published

event handled

event failed

event retried

event moved to dead letter

Include:

EventID

EventType

TenantID

CorrelationID

Handler

Duration

Result

Do not add an external observability platform.

--------------------------------------------------
CORRELATION
--------------------------------------------------

Support:

CorrelationID

CausationID

This will later allow tracing flows such as:

Checklist
→ NonConformity
→ Maintenance
→ Financial
→ Intelligence

without directly coupling the domains.

--------------------------------------------------
DATABASE
--------------------------------------------------

Create a versioned migration for Outbox / Event Processing if the current project architecture uses PostgreSQL migrations.

Suggested:

database/migrations/000016_event_bus_outbox.sql

Tables should support:

event_outbox

event_processing

event_dead_letters

Use:

UUID

TenantID

indexes

timestamps

status

retry metadata

Do not create foreign keys to business modules unless absolutely required.

Prefer aggregate identifiers.

--------------------------------------------------
RBAC
--------------------------------------------------

Do not create a normal business permission for Event Bus.

Event infrastructure is internal.

Administrative inspection endpoints, if created, must be protected separately and must not expose event payloads to normal users.

--------------------------------------------------
API
--------------------------------------------------

DO NOT create public event publishing endpoints.

Business modules must publish through application/domain infrastructure.

If an internal diagnostic endpoint is useful, keep it disabled by default or protected for development only.

--------------------------------------------------
TESTING
--------------------------------------------------

Create comprehensive tests for:

Event creation

Event metadata

Event serialization

Event versioning

Event Bus publish

Multiple handlers

Handler failure isolation

Event ordering within deterministic in-memory dispatch

Idempotency

Duplicate EventID

Retry

Maximum retry

Dead Letter

Tenant isolation

CorrelationID

CausationID

Outbox persistence

MarkPublished

MarkFailed

Processing status

--------------------------------------------------
INTEGRATION TEST

Create an integration test representing:

TireRemoved

      ↓

Event Bus

      ↓

TireRemovedHandler

The handler must execute successfully.

Do NOT create Maintenance or Intelligence behavior.

Another integration test:

FuelTransactionCompleted

      ↓

Event Bus

      ↓

FuelTransactionHandler

Again, handler may only record/observe the event.

--------------------------------------------------
ARCHITECTURAL TEST

Add a test or validation mechanism that ensures:

Tires does NOT import Maintenance.

Fuel does NOT import Financial.

Checklist does NOT import Maintenance.

Financial does NOT import Fuel.

CIOT does NOT import Fleet.

No business module directly imports another business module.

The Event Bus must be the abstraction used for future communication.

--------------------------------------------------
DOCUMENTATION
--------------------------------------------------

Create:

docs/EVENT_BUS.md

Document:

Purpose

Architecture

Domain Event

Event Envelope

Event Naming

Versioning

Event Bus

Handlers

Outbox Pattern

Idempotency

Retry

Dead Letter

Correlation

Causation

Security

Observability

Testing

Future Broker Strategy

--------------------------------------------------
FUTURE BROKER

Do NOT install Kafka, RabbitMQ, NATS, Redis Streams or another external broker in this sprint unless the existing project already requires one.

The architecture must make future replacement possible.

Possible future adapters:

KafkaEventBus

RabbitMQEventBus

NATSEventBus

CloudEventBus

Do not implement them now.

--------------------------------------------------
QUALITY REQUIREMENTS
--------------------------------------------------

DDD

Clean Architecture

SOLID

Dependency Inversion

Hexagonal Architecture

Provider-agnostic infrastructure

Idempotent processing

Immutable event identity

Versioned contracts

Tenant isolation

Production-oriented code

Do not refactor unrelated business logic.

Do not break existing endpoints.

--------------------------------------------------
VALIDATION
--------------------------------------------------

Run:

go test ./...

Verify all existing module tests continue to pass.

Start API with configurable PORT.

Run smoke/integration tests for:

login

event creation

event publish

multiple handlers

handler failure

retry

dead letter

duplicate EventID

idempotent processing

outbox persistence

tenant isolation

correlation ID

causation ID

Verify:

/health

existing authentication

existing RBAC

existing business endpoints

continue working.

--------------------------------------------------
IMPORTANT

DO NOT START SPRINT 016.

DO NOT BUILD INTELLIGENCE.

DO NOT BUILD DASHBOARDS.

DO NOT BUILD FLUTTER.

DO NOT DEPLOY.

DO NOT CONNECT BUSINESS MODULES TO EACH OTHER YET.

Only build and validate the Event Bus infrastructure and standardized Domain Event contracts.

--------------------------------------------------
FINAL REPORT

Report:

1. Files created
2. Files modified
3. Migration created
4. Event contracts created
5. Event Bus architecture
6. Outbox architecture
7. Idempotency strategy
8. Retry strategy
9. Dead Letter strategy
10. Handler strategy
11. Existing event contracts adapted
12. Tests executed
13. Integration tests executed
14. Architecture validation
15. Problems found
16. Recommended next step

STOP AFTER SPRINT 015.

WAIT FOR APPROVAL.

Depois dessa sprint

Não deixe o Codex sair executando a 016 automaticamente.

O resultado que quero ver no relatório é algo parecido com:

Módulos
   ↓
Domain Events
   ↓
Event Bus
   ↓
Outbox
   ↓
Handlers
   ↓
[ futuro ]
Intelligence

Só depois de validarmos isso partimos para a Sprint 016 — Geevheit Intelligence Engine, que será a primeira sprint em que começaremos a transformar os dados dos módulos em indicadores, anomalias, oportunidades de economia e recomendações acionáveis.