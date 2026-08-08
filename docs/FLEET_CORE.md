# Fleet Core Foundation

Sprint 003 creates Fleet Core as an independent module:

`backend/api/internal/modules/fleet`

Fleet Core represents vehicles and assets only. It does not connect to Checklist, Tires, Maintenance, Fuel, CIOT or Intelligence.

## Domain

Core entities:

- Vehicle
- VehicleBrand
- VehicleModel
- VehicleCategoryEntity
- VehicleTypeEntity
- Asset
- LicensePlate
- Renavam
- Chassis
- Engine
- Color
- FuelType
- Transmission
- AxleConfiguration
- EmissionStandard
- OwnershipType
- VehicleStatus

Supported vehicle categories:

- Truck
- Trailer
- Semi Trailer
- Van
- Car
- Motorcycle
- Machinery
- Equipment

Lifecycle:

- Draft
- Active
- Inactive
- Maintenance
- Sold
- Disposed

## Entity Diagram

```text
Vehicle
  |-- VehicleCategory
  |-- VehicleType
  |-- VehicleBrand
  |-- VehicleModel
  |-- Asset
  |-- LicensePlate
  |-- Renavam
  |-- Chassis
```

## ER Diagram

```text
fleet_vehicle_categories
  |-- fleet_vehicle_types

fleet_vehicle_brands
  |-- fleet_vehicle_models

fleet_assets
  |-- fleet_vehicles

fleet_vehicle_categories
  |-- fleet_vehicles

fleet_vehicle_types
  |-- fleet_vehicles

fleet_vehicle_brands
  |-- fleet_vehicles

fleet_vehicle_models
  |-- fleet_vehicles
```

## Validation

- License plate is unique per tenant among non-deleted vehicles.
- Chassis is unique per tenant among non-deleted vehicles.
- Renavam is unique per tenant among non-deleted vehicles.
- Repository methods receive `tenantID` explicitly.
- Delete is soft delete.
- Version is incremented on updates/deletes.

## API

Base path: `/api/v1/fleet`

- `POST /vehicles`
- `GET /vehicles`
- `GET /vehicles/{id}`
- `PUT /vehicles/{id}`
- `DELETE /vehicles/{id}`
- `POST /brands`
- `GET /brands`
- `PUT /brands/{id}`
- `DELETE /brands/{id}`
- `POST /models`
- `GET /models`
- `PUT /models/{id}`
- `DELETE /models/{id}`
- `POST /categories`
- `GET /categories`
- `PUT /categories/{id}`
- `DELETE /categories/{id}`
- `POST /types`
- `GET /types`
- `PUT /types/{id}`
- `DELETE /types/{id}`
- `POST /assets`
- `GET /assets`
- `PUT /assets/{id}`
- `DELETE /assets/{id}`

Search supports:

- `search`
- `page`
- `page_size`
- `sort_by`
- `sort_order`
- `status`
- `category_id`

## Acceptance Criteria

- Fleet module remains independent.
- No business relationships with Checklist, Tires, Maintenance, Fuel or CIOT.
- CRUD contracts exist for vehicles, brands, models, categories, types and assets.
- Database migration uses UUID, tenant_id, audit_id, soft delete, optimistic locking and timestamps.
- Unit tests validate tenant isolation, uniqueness, pagination and soft delete.
