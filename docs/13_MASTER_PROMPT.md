Próxima: Sprint 013 — Financial & Cost Core

Aqui eu quero tomar bastante cuidado.

Não vamos criar contabilidade.

O objetivo é criar o domínio de custos operacionais da frota.

Isso será fundamental para o que você quer vender: mostrar para uma transportadora onde ela está gastando dinheiro e onde existe oportunidade de redução de custo.

Mas ainda não vamos ligar Financial a Tires, Fuel, Maintenance etc.

Primeiro construímos o domínio isolado.

SPRINT 013 — Financial & Cost Core

Read first:

docs/00_PRODUCT_MANIFESTO.md
docs/01_MASTER_PROMPT.md
docs/ARCHITECTURE.md
docs/BACKEND_ARCHITECTURE.md

IMPORTANT

This sprint creates ONLY the Financial & Cost Core domain.

DO NOT connect Financial with Fleet.

DO NOT connect Financial with Drivers.

DO NOT connect Financial with Maintenance.

DO NOT connect Financial with Tires.

DO NOT connect Financial with Fuel.

DO NOT connect Financial with Inventory.

DO NOT connect Financial with Suppliers.

DO NOT connect Financial with Checklist.

DO NOT connect Financial with CIOT.

DO NOT connect Financial with Intelligence.

DO NOT connect Financial with Event Bus.

Do not publish events.

Do not create cross-module dependencies.

This is NOT an accounting system.

This module manages operational costs and financial references required for fleet cost analysis.

--------------------------------------------------

OBJECTIVE

Create an enterprise-grade Financial & Cost Core.

The domain must be capable of representing operational expenses and revenues associated with transportation operations without depending on any other domain.

The future Intelligence Engine will consume these records to calculate:

Cost per vehicle

Cost per kilometer

Cost per trip

Cost per operation

Maintenance cost

Tire cost

Fuel cost

Fixed cost

Variable cost

Total cost

Margin

Profitability

Those calculations must NOT be implemented as cross-domain integrations in this sprint.

--------------------------------------------------

MODULE

internal/modules/financial

Architecture:

domain/
application/
infrastructure/
transport/
ports/

--------------------------------------------------

CORE ENTITIES

FinancialTransaction

Expense

Revenue

CostCategory

CostType

CostCenter

Account

PaymentMethod

FinancialStatus

FinancialPeriod

Budget

BudgetItem

FinancialAttachment

FinancialAdjustment

--------------------------------------------------

TRANSACTION TYPES

Expense

Revenue

Adjustment

Transfer

Refund

Other

--------------------------------------------------

COST CLASSIFICATION

Fixed

Variable

Operational

Administrative

Financial

Extraordinary

Other

--------------------------------------------------

EXPENSE

Support:

Description

Amount

Date

Due Date

Payment Date

Category

Cost Type

Cost Center

Supplier Reference

Document Number

Payment Method

Status

Notes

Attachment Reference

Supplier must remain an identifier/reference only.

Do NOT import Suppliers module.

--------------------------------------------------

REVENUE

Support:

Description

Amount

Date

Due Date

Receipt Date

Category

Cost Center

Document Number

Status

Notes

Attachment Reference

--------------------------------------------------

COST CENTER

Support hierarchical cost centers.

Example:

Operations

Fleet

Maintenance

Tires

Fuel

Administration

Technology

Other

Do not create dependencies with these modules.

They are classifications only.

--------------------------------------------------

ACCOUNT

Represent generic financial accounts.

Account Code

Name

Type

Status

Parent Account

Notes

Do not implement accounting double-entry.

--------------------------------------------------

PAYMENT METHODS

Cash

Bank Transfer

PIX

Credit Card

Debit Card

Invoice

Other

Do not hardcode tenant-specific rules.

--------------------------------------------------

STATUS

Draft

Pending

Approved

Paid

Received

Canceled

Overdue

Adjusted

--------------------------------------------------

FINANCIAL PERIOD

Represent:

Year

Month

Start Date

End Date

Status

Closed periods must not allow arbitrary transaction mutation.

--------------------------------------------------

BUDGET

Represent:

Budget

Period

Cost Center

Category

Planned Amount

Actual Amount

Status

Do not automatically calculate actual amounts from other modules.

--------------------------------------------------

ADJUSTMENT

Historical corrections must use explicit adjustment records.

Do not silently mutate finalized financial transactions.

--------------------------------------------------

BUSINESS RULES

Amount must be greater than zero.

Dates must be valid.

Closed financial periods cannot accept arbitrary mutation.

Paid/Received transactions cannot be directly modified.

Corrections must use FinancialAdjustment.

Tenant isolation is mandatory.

Status transitions must be explicit.

--------------------------------------------------

LIFECYCLE

Expense/Revenue:

Draft
→ Pending
→ Approved
→ Paid/Received

Pending
→ Canceled

Approved
→ Canceled

Paid/Received
→ Adjusted

Do not allow arbitrary status transitions.

--------------------------------------------------

DATABASE

Create versioned migration.

Requirements:

UUID

TenantId

AuditId

Version

CreatedAt

UpdatedAt

DeletedAt

Indexes

Tenant-aware constraints

Financial transaction history must remain auditable.

--------------------------------------------------

AUDIT

Prepare contracts for:

Financial Transaction Created

Expense Created

Expense Approved

Expense Paid

Revenue Created

Revenue Received

Financial Adjustment Created

Financial Period Closed

Budget Created

Do NOT implement Event Bus.

--------------------------------------------------

API

Create:

/api/v1/financial/transactions

/api/v1/financial/expenses

/api/v1/financial/revenues

/api/v1/financial/categories

/api/v1/financial/types

/api/v1/financial/cost-centers

/api/v1/financial/accounts

/api/v1/financial/payment-methods

/api/v1/financial/periods

/api/v1/financial/budgets

/api/v1/financial/adjustments

Support existing conventions:

GET

POST

PATCH/PUT where appropriate

DELETE only where business rules permit

Pagination

Filtering

Sorting

Search

--------------------------------------------------

COMMANDS

Use explicit application commands:

CreateExpense

ApproveExpense

PayExpense

CancelExpense

CreateRevenue

ApproveRevenue

ReceiveRevenue

CancelRevenue

AdjustFinancialTransaction

CloseFinancialPeriod

CreateBudget

--------------------------------------------------

RBAC

Create permission:

financial.financial.manage

Follow existing Identity/RBAC architecture.

--------------------------------------------------

TESTS

Unit Tests

Domain Tests

Lifecycle Tests

Closed Period Tests

Historical Immutability Tests

Adjustment Tests

Repository Tests

Tenant Isolation Tests

API Tests

RBAC Tests

Validation Tests

--------------------------------------------------

OPENAPI

Update the existing OpenAPI specification.

Do not create another OpenAPI specification.

--------------------------------------------------

DOCUMENTATION

Create:

docs/FINANCIAL_CORE.md

Document:

Purpose

Architecture

Domain Model

Expense Lifecycle

Revenue Lifecycle

Cost Classification

Cost Centers

Financial Periods

Budgets

Adjustments

Business Rules

API

Database

Tenant Isolation

Audit

Testing

Acceptance Criteria

--------------------------------------------------

IMPORTANT ARCHITECTURAL RULE

Do NOT implement:

Accounting

Double-entry bookkeeping

Tax calculation

Payroll

Bank reconciliation

Invoice integration

Payment gateway integration

ERP integration

External financial APIs

Cross-module cost calculation

Intelligence

Analytics

These belong to future phases.

--------------------------------------------------

QUALITY

DDD

Clean Architecture

SOLID

Repository Pattern

Dependency Injection

Explicit domain commands

Historical traceability

Production-ready code

Do not refactor unrelated modules.

Do not break existing endpoints.

--------------------------------------------------

VALIDATION

Run:

go test ./...

Start API using configurable PORT.

Perform real HTTP smoke tests for:

login

create category

create cost center

create expense

approve expense

pay expense

attempt invalid mutation of paid expense

create adjustment

create revenue

approve revenue

receive revenue

create financial period

close period

attempt invalid transaction mutation in closed period

RBAC

tenant isolation

Verify configurable PORT is respected.

--------------------------------------------------

FINAL REPORT

Report:

Files created

Files modified

Migration created

Endpoints created

Permissions created

Domain entities

Lifecycle rules

Tests executed

Smoke tests executed

Architecture decisions

Problems found

Do NOT start Sprint 014.

Wait for approval.

E uma observação estratégica

Depois da Sprint 013, teremos:

Operação
├── Fleet
├── Assets
├── Drivers
├── Maintenance
├── Tires
├── Fuel
├── Inventory
├── Checklist
└── Suppliers

Financeiro
└── Financial Core

Aí estaremos muito próximos do ponto em que o Geevheit poderá responder uma pergunta extremamente valiosa para uma transportadora:

"Quanto realmente custa manter e operar cada ativo?"

Mas ainda não vamos conectar nada.