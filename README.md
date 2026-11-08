# Geevheit Intelligence 360

Enterprise Transportation Intelligence Platform.

This repository follows the structure defined in [docs/estrutura projeto.md](docs/estrutura%20projeto.md) and the engineering direction in [docs/01_MASTER_PROMPT.md](docs/01_MASTER_PROMPT.md).

## Current Bootstrap

- Backend API: Go + Gin in `backend/api`
- Architecture baseline: Modular Monolith under `backend/api/internal`
- Core layer: configuration, dependency injection, HTTP router, middleware, logger and errors under `internal/core`
- Modules layer: bounded contexts under `internal/modules`
- Shared package: reusable technical primitives under `internal/shared`
- Initial endpoint: `GET /health`
- Mission Control endpoint in memory demo mode: `GET /api/v1/mission-control/summary`

## Run Locally

The default runtime uses PostgreSQL for identity, fleet, drivers and document metadata. Run migrations and
provision a company explicitly before starting the API. See
[Persistent identity and fleet](docs/PERSISTENCIA_IDENTIDADE_FROTA.md) for setup,
available routes and local validation.

For the original volatile demonstration with all modules:

```bash
./scripts/dev-api.sh
```

This script selects `STORAGE_DRIVER=memory` unless a driver is explicitly set.
Memory storage is rejected when `APP_ENV=production`. No frontend app has been
implemented yet.

## Delivery status

- [Initial audit](docs/DIAGNOSTICO_ENTREGA.md)
- [Delivery checklist](TASKS.md)
- [AI continuation guide](docs/CONTINUIDADE.md)
- [Drivers and document metadata](docs/MOTORISTAS_DOCUMENTOS.md)
- [Identity and fleet persistence](docs/PERSISTENCIA_IDENTIDADE_FROTA.md)

## Architecture Rule

Flutter is presentation. Go owns business logic, tenant validation, permissions, audit, events and intelligence rules.
