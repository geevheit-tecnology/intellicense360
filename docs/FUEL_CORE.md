# Fuel Core

Sprint 011 creates an isolated Fuel Core domain.

Fuel records operational data for future cost, deviation and consumption intelligence, but it does not calculate intelligence and does not connect to Fleet, Drivers, Maintenance, Inventory, Tires, Checklist, Suppliers, Financial or Event Bus.

## Architecture

Module path: `internal/modules/fuel`

- `domain`: entities and enums.
- `application`: commands, validation and lifecycle rules.
- `ports`: service and repository contracts.
- `infrastructure`: in-memory repository used by the current API runtime.
- `transport`: HTTP DTOs, mappers and Gin routes.

## Domain Model

- FuelTransaction
- FuelType
- FuelStation
- FuelStationLocation fields on FuelStation
- FuelTank
- FuelNozzle
- FuelReading
- FuelPrice
- FuelReceipt
- FuelAttachment
- FuelAdjustment
- FuelTransactionStatus

External references such as driver, vehicle and asset are stored as opaque identifiers only.

## Transaction Lifecycle

Statuses:

- `draft`
- `completed`
- `canceled`
- `adjusted`
- `rejected`

Completed transactions preserve historical values. Direct update is blocked after a transaction reaches a final historical status. Corrections must be represented by `FuelAdjustment`.

## Business Rules

- Tenant isolation is mandatory.
- Quantity must be greater than zero.
- Unit price must be zero or greater.
- Total amount must match quantity times unit price within cent precision.
- Odometer and engine hour readings cannot be negative.
- Fuel kind must be valid.
- Fuel station CNPJ is unique per tenant when provided.
- Completed records cannot be silently mutated.
- Final historical transactions are not physically deleted.

## API

Base path: `/api/v1/fuel`

- `GET /transactions`
- `POST /transactions`
- `GET /transactions/{id}`
- `PUT /transactions/{id}`
- `DELETE /transactions/{id}`
- `POST /transactions/{id}/complete`
- `POST /transactions/{id}/cancel`
- `POST /transactions/{id}/adjust`
- `GET /types`
- `POST /types`
- `GET /types/{id}`
- `PUT /types/{id}`
- `DELETE /types/{id}`
- `GET /stations`
- `POST /stations`
- `GET /stations/{id}`
- `PUT /stations/{id}`
- `DELETE /stations/{id}`
- `GET /tanks`
- `POST /tanks`
- `PUT /tanks/{id}`
- `DELETE /tanks/{id}`
- `GET /nozzles`
- `POST /nozzles`
- `PUT /nozzles/{id}`
- `DELETE /nozzles/{id}`
- `GET /readings`
- `POST /readings`
- `GET /prices`
- `POST /prices`
- `PUT /prices/{id}`
- `DELETE /prices/{id}`
- `GET /receipts`
- `POST /receipts`
- `PUT /receipts/{id}`
- `DELETE /receipts/{id}`
- `GET /adjustments`
- `POST /adjustments`

All routes require `fuel.fuel.manage`.

## Database

Migration: `database/migrations/000013_fuel_core.sql`

Tables include UUID primary keys, `tenant_id`, `audit_id`, `version`, timestamps, soft delete and tenant-aware indexes/unique constraints.

## Audit

Contracts are defined in `internal/contracts/events` only. No Event Bus is implemented and no events are published.

## Testing

Focused coverage validates transaction lifecycle, historical immutability, explicit adjustment, tenant isolation, numeric validation and CNPJ uniqueness per tenant.

## Acceptance Criteria

- Fuel domain exists with the standard module folders.
- Fuel is independently testable.
- No cross-module imports were added.
- API routes are protected by RBAC.
- Migration is versioned.
- OpenAPI was updated in the existing specification.
