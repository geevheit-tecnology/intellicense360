# Mission Control

Sprint 017 introduces the backend core for the operational command center.

## Boundary

Mission Control is a read-model and orchestration module. It does not import or call Fleet, Tires, Maintenance, Inventory, Fuel, Suppliers, Financial, CIOT, Checklist, or Intelligence. Future integrations must publish normalized signals into Mission Control contracts instead of coupling this module to operational domains.

## Module Layout

`backend/api/internal/modules/missioncontrol/`

- `domain`: command center entities, controlled vocabularies, SLA/risk/state rules.
- `application`: deterministic services for lifecycle, ranking, risk aggregation, snapshots, and recommendation evaluation.
- `infrastructure`: in-memory repositories for local development and tests.
- `transport`: Gin HTTP endpoints.
- `ports`: repository and service contracts.

## Lifecycle

Command items follow this state machine:

- `open -> acknowledged | dismissed | resolved`
- `acknowledged -> in_progress | dismissed | resolved`
- `in_progress -> resolved | dismissed`

Invalid transitions return `invalid_status_transition`.

## Endpoints

- `GET /api/v1/mission-control/summary`
- `GET /api/v1/mission-control/items`
- `GET /api/v1/mission-control/items/:id`
- `POST /api/v1/mission-control/items`
- `PATCH /api/v1/mission-control/items/:id`
- `POST /api/v1/mission-control/items/:id/acknowledge`
- `POST /api/v1/mission-control/items/:id/start`
- `POST /api/v1/mission-control/items/:id/resolve`
- `POST /api/v1/mission-control/items/:id/dismiss`
- `GET /api/v1/mission-control/items/:id/actions`
- `POST /api/v1/mission-control/items/:id/actions`
- `GET /api/v1/mission-control/items/:id/history`
- `GET /api/v1/mission-control/snapshot`
- `POST /api/v1/mission-control/snapshot/rebuild`
- `GET /api/v1/mission-control/recommendations`
- `POST /api/v1/mission-control/recommendations/evaluate`

## RBAC

Permissions:

- `mission_control.read`
- `mission_control.create`
- `mission_control.update`
- `mission_control.acknowledge`
- `mission_control.resolve`
- `mission_control.dismiss`
- `mission_control.snapshot`
- `mission_control.admin`

## Persistence

The versioned PostgreSQL contract is `database/migrations/000019_mission_control_core.sql`.

Tables:

- `mission_control_command_items`
- `mission_control_command_actions`
- `mission_control_command_events`
- `mission_control_snapshots`
- `mission_control_idempotency`

The current runtime adapter remains in-memory, consistent with the current sprint boundary.
