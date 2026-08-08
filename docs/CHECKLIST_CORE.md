# Checklist Core

Sprint 004 implements Checklist Core as an independent module:

`backend/api/internal/modules/checklist`

The module stores `vehicle_id` as an opaque UUID/string. It does not call Fleet and does not validate vehicle existence.

## Architecture

- `domain`: Checklist, ChecklistItem and ChecklistAnswer.
- `application`: use cases and validation.
- `ports`: service and repository contracts.
- `infrastructure`: in-memory repositories.
- `transport/dto`: HTTP request DTOs.
- `transport/mapper`: DTO to domain mapping.
- `transport/http`: Gin handlers.

Handlers do not contain business rules.

## Status

- `draft`
- `in_progress`
- `completed`
- `cancelled`

## Answer Types

- `boolean`
- `text`
- `number`
- `photo`
- `signature`
- `select`

## Use Cases

- CreateChecklist
- UpdateChecklist
- DeleteChecklist
- GetChecklist
- ListChecklists
- StartChecklist
- FinishChecklist
- CancelChecklist
- AddItem
- UpdateItem
- DeleteItem
- ListItems
- AnswerItem
- ListAnswers

## Routes

Base path: `/api/v1/checklists`

- `GET /`
- `POST /`
- `GET /{id}`
- `PUT /{id}`
- `DELETE /{id}`
- `POST /{id}/start`
- `POST /{id}/finish`
- `POST /{id}/cancel`
- `GET /{id}/items`
- `POST /{id}/items`
- `PUT /{id}/items/{item_id}`
- `DELETE /{id}/items/{item_id}`
- `POST /{id}/answers`
- `GET /{id}/answers`

## Data Rules

- Tenant isolation.
- Soft delete for checklists and items.
- Search, pagination, filters and sorting.
- `vehicle_id` is stored only, with no Fleet dependency.

## ER Diagram

```text
checklists
  |-- checklist_items
  |-- checklist_answers

checklist_items
  |-- checklist_answers
```
