# Architecture

Geevheit Intelligence 360 is organized as a monorepo with apps, backend services, packages, database assets and infrastructure.

## Backend

The first executable boundary is `backend/api`, a Go/Gin API gateway prepared for vertical slices.

Each business module should be organized around:

- `domain`: entities, value objects, domain events and domain services
- `application`: commands, queries, use cases and ports
- infrastructure adapters only when concrete persistence, messaging or external integrations are added

## Tenancy

Every persisted business record must carry `tenant_id`.

Repository implementations must receive tenant context and filter every query by tenant. Presentation-layer filtering is not sufficient.

## Intelligence

Modules must emit operational events that feed Mission Control, Cost Radar, Fleet Pulse and recommendation use cases.
