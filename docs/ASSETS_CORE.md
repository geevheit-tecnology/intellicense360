# Fleet Assets & Equipment Foundation

Sprint 006 creates the independent `assets` module:

`backend/api/internal/modules/assets`

It manages assets owned or operated by the company. It does not call Fleet, Maintenance, Checklist, Tires, Fuel, CIOT or Intelligence.

## Architecture

- `domain`: Asset, catalogs, equipment, implement, depreciation, location and attachment models.
- `application`: use cases and validation.
- `ports`: service, repository, audit, import, export and attachment contracts.
- `infrastructure`: in-memory repositories and audit recorder.
- `transport`: HTTP handlers, DTOs and mappers.

## Supported Assets

Truck, Trailer, Semi Trailer, Dolly, Bitrem, Rodotrem, Container Chassis, Tank, Refrigerated Trailer, Platform Trailer, Dump Trailer, Van, Forklift, Crane, Generator, Compressor and Other Equipment.

## Lifecycle

Draft, Available, Assigned, In Operation, Maintenance, Inactive, Sold and Disposed.

## Entity Diagram

```text
Asset
  |-- AssetCategory
  |-- AssetType
  |-- Manufacturer
  |-- Model
  |-- Equipment
  |-- Implement
  |-- Location
  |-- Depreciation
  |-- Attachments
```

## ER Diagram

```text
asset_categories
  |-- asset_types

asset_manufacturers
  |-- asset_models

assets
  |-- asset_equipment
  |-- asset_implements
  |-- asset_attachments
```

## Sequence Diagram

```text
HTTP -> Handler -> Application Service -> Repository
                               |
                               +-> AuditRecorder contract
```

## API

Base path: `/api/v1/assets`

- CRUD `/`
- CRUD `/categories`
- CRUD `/types`
- CRUD `/manufacturers`
- CRUD `/models`
- CRUD `/equipment`

Search supports `search`, `page`, `page_size`, `sort_by`, `sort_order`, `status`, `category_id` and `type_id`.

## Future Contracts

- ImportPort prepares CSV/Excel import.
- ExportPort prepares PDF/Excel export.
- AttachmentPort prepares images, PDFs, warranty files, invoices and manuals.
- AuditRecorder records asset operation events.

## Acceptance Criteria

- Assets module is isolated from other domains.
- UUID-ready persistence migration exists.
- Tenant isolation, soft delete, optimistic version and indexes are covered.
- Internal code, serial number and asset tag uniqueness are validated per tenant.
- Tests cover validation, tenant isolation, soft delete and catalog operations.
