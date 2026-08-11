create table if not exists financial_categories (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    name text not null,
    code text,
    description text,
    classification text,
    active boolean not null default true,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    unique (tenant_id, code)
);

create table if not exists financial_cost_types (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    name text not null,
    code text,
    description text,
    classification text,
    active boolean not null default true,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    unique (tenant_id, code)
);

create table if not exists financial_cost_centers (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    parent_id uuid references financial_cost_centers(id),
    name text not null,
    code text,
    notes text,
    active boolean not null default true,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    unique (tenant_id, code)
);

create table if not exists financial_accounts (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    parent_id uuid references financial_accounts(id),
    account_code text,
    name text not null,
    type text,
    status text,
    notes text,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    unique (tenant_id, account_code)
);

create table if not exists financial_payment_methods (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    name text not null,
    code text,
    method_type text,
    notes text,
    active boolean not null default true,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    unique (tenant_id, code)
);

create table if not exists financial_periods (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    year integer not null,
    month integer not null,
    start_date date not null,
    end_date date not null,
    status text not null,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    unique (tenant_id, year, month)
);

create table if not exists financial_transactions (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    kind text not null,
    description text not null,
    amount numeric(14,2) not null,
    date date not null,
    due_date date,
    settlement_date date,
    category_id uuid references financial_categories(id),
    cost_type_id uuid references financial_cost_types(id),
    cost_center_id uuid references financial_cost_centers(id),
    account_id uuid references financial_accounts(id),
    supplier_reference text,
    operation_reference text,
    document_number text,
    payment_method_id uuid references financial_payment_methods(id),
    status text not null,
    notes text,
    attachment_reference text,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    check (amount > 0)
);

create table if not exists financial_budgets (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    period_id uuid references financial_periods(id),
    cost_center_id uuid references financial_cost_centers(id),
    category_id uuid references financial_categories(id),
    name text not null,
    status text not null,
    planned_amount numeric(14,2) not null default 0,
    actual_amount numeric(14,2) not null default 0,
    notes text,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create table if not exists financial_adjustments (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    transaction_id uuid not null references financial_transactions(id),
    adjustment_type text not null,
    reason text not null,
    adjusted_amount numeric(14,2) not null,
    notes text,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    check (adjusted_amount > 0)
);

create index if not exists financial_transactions_tenant_kind_idx on financial_transactions (tenant_id, kind);
create index if not exists financial_transactions_tenant_status_idx on financial_transactions (tenant_id, status);
create index if not exists financial_transactions_tenant_date_idx on financial_transactions (tenant_id, date);
create index if not exists financial_periods_tenant_status_idx on financial_periods (tenant_id, status);
create index if not exists financial_budgets_tenant_period_idx on financial_budgets (tenant_id, period_id);
