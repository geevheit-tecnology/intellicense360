# Checklist Engine Core

Sprint 012 evolves the existing Checklist module into an isolated Checklist Engine.

The previous checklist CRUD remains available for compatibility. The new engine adds templates, versions, executions, responses, evidence, non-conformities, signatures, assignments and append-only history without connecting to Fleet, Drivers, Maintenance, Tires, Inventory, Fuel, Suppliers, Financial, CIOT, Intelligence or Event Bus.

## Architecture

Module path: `internal/modules/checklist`

- `domain`: legacy checklist entities plus engine entities.
- `application`: legacy services plus explicit engine commands.
- `ports`: repository and service contracts.
- `infrastructure`: in-memory repositories.
- `transport`: Gin HTTP routes.

## Existing Checklist Compatibility

Existing endpoints are preserved:

- `GET /api/v1/checklists`
- `POST /api/v1/checklists`
- `GET /api/v1/checklists/{id}`
- `PUT /api/v1/checklists/{id}`
- `DELETE /api/v1/checklists/{id}`
- `POST /api/v1/checklists/{id}/start`
- `POST /api/v1/checklists/{id}/finish`
- `POST /api/v1/checklists/{id}/cancel`
- `GET/POST /api/v1/checklists/{id}/items`
- `GET/POST /api/v1/checklists/{id}/answers`

## Domain Model

- ChecklistTemplate
- ChecklistTemplateVersion
- ChecklistType
- ChecklistSection
- ChecklistEngineItem
- ChecklistItemOption
- ChecklistExecution
- ChecklistResponse
- ChecklistEvidence
- ChecklistNonConformity
- ChecklistSignature
- ChecklistAssignment
- ChecklistHistory

## Template Versioning

Templates start as `draft`. Versions are created from templates and can be published. A published version is immutable from the application layer: changes require creating a new version.

Historical executions always reference a single exact template version.

## Execution Lifecycle

- `draft`
- `in_progress`
- `completed`
- `canceled`
- `invalidated`

Supported commands:

- CreateChecklistTemplate
- CreateTemplateVersion
- PublishTemplateVersion
- ArchiveTemplate
- StartChecklistExecution
- RecordChecklistResponse
- AddChecklistEvidence
- CreateNonConformity
- CompleteChecklistExecution
- CancelChecklistExecution
- InvalidateChecklistExecution

## Business Rules

- Tenant isolation is mandatory.
- Published template versions cannot be edited.
- Execution must reference exactly one published template version.
- Required items must be answered before completion.
- Evidence-required items need evidence before completion.
- Signature-required template versions need a signature before completion.
- Terminal transitions are explicit.
- Execution history is append-only.
- No integration creates Maintenance orders or blocks Fleet assets in this sprint.

## API

New engine endpoints:

- `GET /api/v1/checklists/templates`
- `POST /api/v1/checklists/templates`
- `GET /api/v1/checklists/templates/{id}`
- `PUT /api/v1/checklists/templates/{id}`
- `POST /api/v1/checklists/templates/{id}/archive`
- `GET /api/v1/checklists/templates/{id}/versions`
- `POST /api/v1/checklists/templates/{id}/versions`
- `POST /api/v1/checklists/templates/versions/{version_id}/publish`
- `GET /api/v1/checklists/types`
- `POST /api/v1/checklists/types`
- `GET /api/v1/checklists/sections`
- `POST /api/v1/checklists/sections`
- `GET /api/v1/checklists/items`
- `POST /api/v1/checklists/items`
- `GET /api/v1/checklists/executions`
- `POST /api/v1/checklists/executions`
- `GET /api/v1/checklists/executions/{id}`
- `POST /api/v1/checklists/executions/{id}/complete`
- `POST /api/v1/checklists/executions/{id}/cancel`
- `POST /api/v1/checklists/executions/{id}/invalidate`
- `GET/POST /api/v1/checklists/executions/{id}/responses`
- `GET/POST /api/v1/checklists/executions/{id}/evidence`
- `GET/POST /api/v1/checklists/executions/{id}/non-conformities`
- `GET /api/v1/checklists/executions/{id}/history`

All Checklist routes require the existing `checklist.checklists.manage` permission and the new `checklist.checklist.manage` engine permission.

## Database

Migration: `database/migrations/000014_checklist_engine_core.sql`

The migration adds versioned, tenant-scoped tables with UUID, `audit_id`, `version`, timestamps, soft delete, indexes and tenant-aware constraints.

## Audit

Audit/event contracts are prepared in `internal/contracts/events`. No Event Bus is implemented and no events are published.

## Migration Strategy

The legacy Checklist entities are left in place. Future migration can map legacy `Checklist`, `ChecklistItem` and `ChecklistAnswer` into `ChecklistTemplate`, `ChecklistTemplateVersion`, `ChecklistEngineItem`, `ChecklistExecution` and `ChecklistResponse`.

## Testing

Focused tests cover:

- Published version immutability.
- Required item validation.
- Required evidence validation.
- Execution completion.
- Existing checklist service compatibility through the existing tests.

## Acceptance Criteria

- Existing Checklist routes still compile and remain registered.
- Checklist Engine has explicit lifecycle commands.
- No cross-module imports were added.
- Migration is versioned.
- OpenAPI and docs are updated.
