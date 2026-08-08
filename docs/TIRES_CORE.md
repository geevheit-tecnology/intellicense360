# Tire Lifecycle Core

Sprint 010 evolves Tires into an independent lifecycle domain.

The module stores `vehicle_id`, `position` and provider/supplier references as opaque values. It does not call Fleet, Checklist, Maintenance, Inventory, Fuel, Suppliers, Financial, CIOT or Intelligence.

## Lifecycle

```text
Purchase
  -> Receipt
  -> Stock
  -> Installation
  -> Position
  -> KM accumulation
  -> Inspections
  -> Rotations
  -> Retreads
  -> Reinstallations
  -> Wear tracking
  -> Cost per KM contracts
  -> Removal
  -> Disposal
```

## Entities

- Tire
- TireBrand
- TireModel
- TireSpecification
- TirePosition
- TireInstallation
- TireRemoval
- TireMovement
- TireInspection
- TireMeasurement
- TireRetread
- TireRetreadEvent
- TireStatus
- TireCondition
- TireHistory
- TireCost
- TireDisposal
- TireAttachment

## Status

- `new`
- `in_stock`
- `installed`
- `removed`
- `under_inspection`
- `under_retread`
- `retreaded`
- `reserved`
- `damaged`
- `end_of_life`
- `disposed`
- `lost`

## State Machine

Allowed transitions:

- `new -> in_stock`
- `in_stock -> installed`
- `installed -> removed`
- `removed -> in_stock`
- `removed -> under_retread`
- `under_retread -> retreaded`
- `retreaded -> in_stock`
- `installed -> damaged`
- `damaged -> end_of_life`
- `end_of_life -> disposed`

The application does not expose arbitrary lifecycle status mutation. Transitions go through explicit commands such as receive, install, remove, rotate, retread, return and dispose.

## API

Base path: `/api/v1/tires`

- `GET /`
- `POST /`
- `GET /{id}`
- `PUT /{id}`
- `DELETE /{id}`
- `POST /{id}/receive`
- `POST /{id}/install`
- `POST /{id}/remove`
- `POST /{id}/rotate`
- `POST /{id}/recap`
- `POST /{id}/return`
- `POST /{id}/dispose`
- `GET /{id}/inspections`
- `POST /{id}/inspections`
- `GET /{id}/movements`
- `POST /{id}/movements`
- `GET /{id}/history`
- `GET /brands`
- `GET /models`
- `GET /specifications`
- `GET /positions`
- `GET /inspections`
- `POST /inspections`
- `GET /measurements`
- `GET /installations`
- `GET /removals`
- `GET /movements`
- `GET /retreads`
- `GET /history`
- `GET /attachments`

## Rules

- Tenant isolation.
- Serial number is unique per tenant when provided.
- Fire number remains unique per tenant.
- Status and condition must be valid.
- Tread depth and pressure cannot be negative.
- Dimension must be valid when provided.
- History and movement records are append-only from the business perspective.
- Master records may use soft delete where appropriate.

## Contracts

- `ports.CostCalculator` prepares future cost-per-KM calculation.
- `ports.AttachmentPort` prepares storage integration for photos, invoices, inspection photos, retread documents, warranty documents and disposal documents.
- No Event Bus is integrated and no events are published.

## ER Diagram

```text
tires
  |-- tire_inspections
  |-- tire_measurements
  |-- tire_movements
  |-- tire_installations
  |-- tire_removals
  |-- tire_retreads
  |-- tire_retread_events
  |-- tire_history
  |-- tire_costs
  |-- tire_disposals
  |-- tire_attachments

tire_brands
  |-- tire_models

tire_specifications
tire_positions
```
