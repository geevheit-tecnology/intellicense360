# Backend Architecture

## Sprint 002 Objective

Refactor the backend into an Enterprise Ready Backend Foundation using a Modular Monolith architecture.

`docs/00_PRODUCT_MANIFESTO.md` was requested for this sprint, but it is not present in the repository at the time of implementation. The available source of product and architecture truth is `docs/01_MASTER_PROMPT.md`.

## Architectural Style

The backend is a Modular Monolith.

The system deploys as one Go process at `backend/api`, while each business capability is isolated under `internal/modules`.

This gives the project:

- clear module boundaries
- one deployment unit during early product development
- simple local operation
- a future path to extract services only when operational pressure justifies it

## Runtime Entry Point

The executable entry point is:

```text
backend/api/cmd/api/main.go
```

The entry point only loads configuration, builds the dependency container and starts the HTTP router.

## Core Layer

`internal/core` contains platform foundations that are shared by modules but do not contain product business rules.

Current packages:

- `config`: environment-based configuration loader
- `di`: dependency injection container
- `errors`: application error model
- `http`: Gin router composition
- `logger`: logging abstraction
- `middleware`: request logging, error handling, tenant, audit and authentication middleware

## Modules Layer

`internal/modules` contains bounded contexts.

Each module should follow this layout:

```text
module/
  domain/
  application/
  ports/
  transport/http/
```

`domain` owns entities, value objects and domain events.

`application` owns commands, queries and use cases.

`ports` owns repository and service interfaces.

`transport/http` exposes HTTP handlers for that module.

## Shared Package

`internal/shared` contains technical primitives that can be reused without coupling modules to each other.

Current shared packages:

- `contextkeys`: tenant and actor context helpers
- `types`: base entity metadata

Shared code must stay generic. Business concepts belong inside modules.

## Repository Interfaces

Repository interfaces live in each module's `ports` package.

They define persistence contracts but do not provide database implementations yet.

Current examples:

- `identity/ports.UserRepository`
- `intelligence/ports.RecommendationRepository`
- `fleet/ports.VehicleRepository`
- `checklist/ports.ChecklistRepository`
- `drivers/ports.DriverRepository`
- `maintenance/ports.MaintenanceRepository`
- `tires/ports.TireRepository`
- `fuel/ports.FuelEventRepository`
- `financial/ports.FinancialEventRepository`
- `documents/ports.DocumentRepository`
- `inventory/ports.InventoryItemRepository`
- `notifications/ports.NotificationRepository`
- `ciot/ports.CIOTRepository`

## Service Interfaces

Service interfaces live in each module's `ports` package.

They define module-facing contracts for application behavior and future cross-module orchestration.

No concrete business rules were implemented in this sprint.

## Dependency Injection

`internal/core/di.Container` is the composition root.

It wires:

- configuration
- logger abstraction
- module application services
- module HTTP handlers
- HTTP router

Concrete databases, queues, object storage, feature flags and external clients should be added to the container only when their adapters exist.

## Middleware

Middleware is centralized under `internal/core/middleware`.

Current middleware:

- `RequestLogger`: logs method, path, status and duration
- `ErrorHandler`: converts application errors to JSON responses
- `Tenant`: attaches tenant information to the request context
- `Audit`: records request audit events through the logger abstraction
- `Authentication`: placeholder JWT boundary that only extracts authenticated context when a bearer token exists

Authentication and tenant validation are foundations only. Real JWT verification, RBAC and persistence-backed tenant validation are future implementation tasks.

## OpenAPI

OpenAPI contracts now live under:

```text
backend/api/openapi/openapi.yaml
```

The legacy root `backend/api/openapi.yaml` remains for compatibility with the initial bootstrap, but new API contract work should happen in the `openapi` folder.

## Current Endpoints

These endpoints remain active:

```text
GET /health
GET /api/v1/mission-control/summary
```

`/api/v1/mission-control/summary` is routed through the Mission Control module handler.

## Business Rules

No business rules were implemented in this sprint.

The current Mission Control response is still a bootstrap response used to keep the endpoint working while the architecture is prepared.
