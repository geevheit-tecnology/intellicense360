# Inventory Core

Sprint 008 creates Inventory & Parts Core as an isolated business domain.

It does not connect to Maintenance, Fleet, Tires, Fuel, Checklist, Financial or Intelligence. Other domains will reference inventory later through contracts only.

## Subdomains

- Parts
- Categories
- Brands
- Models
- Supplier References
- Warehouses
- Warehouse Locations
- Stock Items
- Stock Batches
- Stock Levels
- Units of Measure
- Minimum Stock
- Maximum Stock
- Attachments

## Architecture

- `domain`: entities and enums.
- `application`: validation and use cases.
- `ports`: service, repository and attachment contracts.
- `infrastructure`: in-memory repositories.
- `transport`: DTOs, mappers and HTTP handlers.

## API

Base path: `/api/v1/inventory`

- `GET /parts`
- `POST /parts`
- `GET /parts/{id}`
- `PUT /parts/{id}`
- `DELETE /parts/{id}`
- `GET /categories`
- `POST /categories`
- `PUT /categories/{id}`
- `DELETE /categories/{id}`
- `GET /units`
- `POST /units`
- `PUT /units/{id}`
- `DELETE /units/{id}`
- `GET /warehouses`
- `POST /warehouses`
- `PUT /warehouses/{id}`
- `DELETE /warehouses/{id}`
- `GET /locations`
- `POST /locations`
- `PUT /locations/{id}`
- `DELETE /locations/{id}`

## Validations

- SKU is unique per tenant.
- Internal code is unique per tenant.
- Unit of measure must exist in the tenant unit catalog.
- Minimum stock cannot be greater than maximum stock.
- Tenant isolation, soft delete, pagination, filtering and sorting are enforced in repositories/use cases.

## Event Contracts

Inventory defines event contracts only:

- `part.created`
- `part.updated`

No event bus or publisher is implemented in this sprint.

## ER Diagram

```text
inventory_categories
inventory_units
inventory_brands
inventory_models

inventory_parts
  |-- inventory_supplier_references
  |-- inventory_stock_items
  |-- inventory_stock_batches
  |-- inventory_stock_levels
  |-- inventory_stock_limits
  |-- inventory_attachments

inventory_warehouses
  |-- inventory_locations
  |-- inventory_stock_items
  |-- inventory_stock_levels
```

## Acceptance Criteria

- Inventory remains independent from all other business domains.
- CRUD exists for Parts, Categories, Warehouses, Units and Locations.
- Core entities and storage schema exist for the full inventory model.
- Attachment contracts exist without storage integration.
- Unit, repository and validation tests cover the critical rules.
