# Financial & Cost Core

Sprint 013 creates an isolated operational cost domain. This is not accounting: it does not implement double-entry bookkeeping, taxes, payroll, bank reconciliation, invoices, payment gateways, ERP integration, analytics or cross-module cost calculation.

## Architecture

Module path: `internal/modules/financial`

- `domain`: financial transaction, expense, revenue, categories, cost types, cost centers, accounts, periods, budgets and adjustments.
- `application`: explicit lifecycle commands and validations.
- `ports`: service and repository contracts.
- `infrastructure`: in-memory repositories.
- `transport`: HTTP routes.

## Domain Model

- FinancialTransaction
- Expense
- Revenue
- CostCategory
- CostType
- CostCenter
- Account
- PaymentMethod
- FinancialPeriod
- Budget
- BudgetItem
- FinancialAttachment
- FinancialAdjustment

Supplier, operation and asset references remain opaque text identifiers.

## Lifecycle

Expense:

- `draft -> approved`
- `pending -> approved`
- `approved -> paid`
- `pending/approved -> canceled`
- `paid -> adjusted` through `FinancialAdjustment`

Revenue:

- `draft -> approved`
- `pending -> approved`
- `approved -> received`
- `pending/approved -> canceled`
- `received -> adjusted` through `FinancialAdjustment`

Paid and received transactions cannot be mutated directly.

## Financial Periods

Periods represent year, month, start date, end date and status. Closed periods block arbitrary transaction creation/update for dates inside the period.

## API

Base path: `/api/v1/financial`

- `GET /transactions`
- `GET /transactions/{id}`
- `PUT /transactions/{id}`
- `POST /transactions/{id}/approve`
- `POST /transactions/{id}/pay`
- `POST /transactions/{id}/receive`
- `POST /transactions/{id}/cancel`
- `POST /transactions/{id}/adjust`
- `GET/POST /expenses`
- `GET/POST /revenues`
- `GET/POST /categories`
- `GET/POST /types`
- `GET/POST /cost-centers`
- `GET/POST /accounts`
- `GET/POST /payment-methods`
- `GET/POST /periods`
- `POST /periods/{id}/close`
- `GET/POST /budgets`
- `GET/POST /adjustments`

All routes require `financial.financial.manage`.

## Database

Migration: `database/migrations/000015_financial_cost_core.sql`

The schema uses UUID primary keys, `tenant_id`, `audit_id`, versioning, timestamps, soft delete, tenant-aware unique constraints and indexes.

## Audit

Audit contracts are prepared in `internal/contracts/events`. No Event Bus is implemented and no events are published.

## Testing

Tests cover lifecycle, finalized immutability, adjustments, closed period validation and tenant isolation.

## Acceptance Criteria

- Financial domain is independently testable.
- No cross-module imports were added.
- Lifecycle transitions are explicit.
- Closed periods block mutation.
- Migration, docs, OpenAPI and RBAC are in place.
