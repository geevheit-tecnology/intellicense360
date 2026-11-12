-- Migration 000028: Financial & Cash Flow Operational Reconciliation
-- Connects fuel notes, maintenance work orders and third-party bills into durable financial transactions and cash closing.

alter table if exists financial_transactions
    add column if not exists vehicle_id uuid references fleet_vehicles(id),
    add column if not exists driver_id uuid references drivers(id),
    add column if not exists fuel_transaction_id uuid,
    add column if not exists work_order_id uuid,
    add column if not exists invoice_key text,
    add column if not exists invoice_number text,
    add column if not exists requisition_number text,
    add column if not exists is_third_party boolean not null default false,
    add column if not exists third_party_identifier text,
    add column if not exists payment_method text,
    add column if not exists barcode text,
    add column if not exists bank_account text,
    add column if not exists cash_flow_nature text not null default 'operational';

create index if not exists financial_tx_tenant_vehicle_idx on financial_transactions (tenant_id, vehicle_id);
create index if not exists financial_tx_tenant_invoice_idx on financial_transactions (tenant_id, invoice_number);
create index if not exists financial_tx_tenant_req_idx on financial_transactions (tenant_id, requisition_number);
create index if not exists financial_tx_tenant_third_party_idx on financial_transactions (tenant_id, is_third_party);
create index if not exists financial_tx_tenant_due_date_idx on financial_transactions (tenant_id, due_date);
create index if not exists financial_tx_tenant_settlement_idx on financial_transactions (tenant_id, settlement_date);

-- Table for daily cash closing and cash flow auditing
create table if not exists financial_cash_closings (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    closing_date date not null,
    opening_balance numeric(14,2) not null default 0,
    total_revenues numeric(14,2) not null default 0,
    total_fleet_expenses numeric(14,2) not null default 0,
    total_third_party_expenses numeric(14,2) not null default 0,
    total_other_expenses numeric(14,2) not null default 0,
    closing_balance numeric(14,2) not null default 0,
    status text not null default 'open',
    notes text,
    closed_by uuid references users(id),
    closed_at timestamptz,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    unique (tenant_id, closing_date)
);

create index if not exists financial_closings_tenant_date_idx on financial_cash_closings (tenant_id, closing_date);

-- Permissions
insert into permissions(id, key, description) values
    (gen_random_uuid(), 'financial.financial.read', 'Consultar transações financeiras e fluxo de caixa'),
    (gen_random_uuid(), 'financial.cashflow.close', 'Realizar fechamento diário de caixa da transportadora')
on conflict (key) do nothing;

insert into roles_permissions(tenant_id, role_id, permission_id)
select r.tenant_id, r.id, p.id
from roles r
cross join permissions p
where r.name = 'admin' and p.key in (
    'financial.financial.read',
    'financial.cashflow.close'
)
on conflict do nothing;

