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
- Initial intelligence endpoint: `GET /api/v1/mission-control/summary`

## Run Locally

```bash
cd backend/api
go mod tidy
go run ./cmd/api
```

Then check:

```bash
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/mission-control/summary
```

## Architecture Rule

Flutter is presentation. Go owns business logic, tenant validation, permissions, audit, events and intelligence rules.
