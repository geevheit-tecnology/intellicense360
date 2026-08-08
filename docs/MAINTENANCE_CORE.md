# Maintenance Core

Maintenance is a business domain, not a simple work-order screen.

It is independent from Fleet, Assets, Tires, Fuel, Checklist, CIOT and Intelligence. It stores opaque references such as `asset_id` and `vehicle_id` without validating them in other modules.

## Subdomains

- Work Orders
- Preventive Plans
- Corrective Records
- Categories
- Types
- Priorities
- Reasons
- Workshops
- Technicians
- Service Types
- Labor
- Downtime
- Availability
- History

## Lifecycle

Work order status:

- `draft`
- `open`
- `waiting`
- `approved`
- `executing`
- `in_progress`
- `paused`
- `completed`
- `canceled`

Maintenance kind:

- `preventive`
- `corrective`
- `predictive`
- `inspection`
- `emergency`
- `warranty`
- `external`
- `internal`

Priority:

- `low`
- `medium`
- `high`
- `critical`

## Architecture

- `domain`: maintenance entities and value enums.
- `application`: lifecycle and validation use cases.
- `ports`: repository and service contracts.
- `infrastructure`: in-memory repositories.
- `transport`: DTOs, mappers and HTTP handlers.

## API

Base path: `/api/v1/maintenance`

- `GET /work-orders`
- `POST /work-orders`
- `GET /work-orders/{id}`
- `PUT /work-orders/{id}`
- `DELETE /work-orders/{id}`
- `POST /work-orders/{id}/start`
- `POST /work-orders/{id}/complete`
- `POST /work-orders/{id}/cancel`
- `GET /work-orders/{id}/labor`
- `POST /work-orders/{id}/labor`
- `GET /work-orders/{id}/downtime`
- `POST /work-orders/{id}/downtime`
- `POST /work-orders/{id}/downtime/{downtime_id}/end`
- `GET /work-orders/{id}/history`
- `GET /preventive-plans`
- `POST /preventive-plans`
- `GET /service-types`
- `POST /service-types`
- `GET /categories`
- `POST /categories`
- `PUT /categories/{id}`
- `DELETE /categories/{id}`
- `GET /types`
- `POST /types`
- `PUT /types/{id}`
- `DELETE /types/{id}`
- `GET /priorities`
- `POST /priorities`
- `PUT /priorities/{id}`
- `DELETE /priorities/{id}`
- `GET /reasons`
- `POST /reasons`
- `PUT /reasons/{id}`
- `DELETE /reasons/{id}`
- `GET /workshops`
- `POST /workshops`
- `PUT /workshops/{id}`
- `DELETE /workshops/{id}`
- `GET /technicians`
- `POST /technicians`
- `PUT /technicians/{id}`
- `DELETE /technicians/{id}`

## Integration Boundaries

- Files are represented by `AttachmentPort` and `MaintenanceAttachment`; storage implementation stays in `internal/platform/storage`.
- Audit is represented by `AuditRecorder`; this sprint does not couple Maintenance to the Identity audit store.
- Reports are represented by `ReportPort`; this sprint exposes report contracts only.
- Cross-domain references such as `asset_id` and `vehicle_id` remain opaque IDs.

## Entity Diagram

```text
WorkOrder
  |-- MaintenanceCategory
  |-- MaintenanceReason
  |-- Workshop
  |-- Technician
  |-- ServiceType
  |-- ServiceItem
  |-- LaborEntry
  |-- Downtime
  |-- CorrectiveRecord
  |-- MaintenanceAttachment
  |-- MaintenanceHistory

PreventivePlan
  |-- ServiceType
  |-- MaintenanceSchedule

AvailabilitySnapshot
```

## ER Diagram

```text
maintenance_service_types
  |-- maintenance_work_orders
  |-- maintenance_preventive_plans

maintenance_categories
maintenance_types
maintenance_priorities
maintenance_reasons
maintenance_workshops
maintenance_technicians

maintenance_work_orders
  |-- maintenance_corrective_records
  |-- maintenance_service_items
  |-- maintenance_labor_entries
  |-- maintenance_downtime
  |-- maintenance_schedules
  |-- maintenance_attachments
  |-- maintenance_history
```

## Acceptance Criteria

- Domain is isolated from other modules.
- Work orders are only one part of Maintenance.
- Preventive, corrective, catalog, workshop, technician, labor, downtime, availability, attachment, audit, report and history contracts exist.
- Tenant isolation, soft delete, pagination, unique order code, enum validation and lifecycle transitions are covered by tests.
