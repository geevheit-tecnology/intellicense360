create table if not exists fuel_types (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    name text not null,
    kind text not null,
    code text,
    description text,
    active boolean not null default true,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    unique (tenant_id, code)
);

create table if not exists fuel_stations (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    name text not null,
    legal_name text,
    cnpj text,
    address text,
    city text,
    state text,
    country text,
    latitude numeric(10,7),
    longitude numeric(10,7),
    active boolean not null default true,
    notes text,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create unique index if not exists fuel_stations_tenant_cnpj_uidx
on fuel_stations (tenant_id, cnpj)
where cnpj is not null and deleted_at is null;

create table if not exists fuel_tanks (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    code text not null,
    name text,
    capacity numeric(14,3) not null default 0,
    current_reading numeric(14,3) not null default 0,
    fuel_type_id uuid references fuel_types(id),
    fuel_kind text not null,
    location_reference text,
    status text not null,
    notes text,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    unique (tenant_id, code)
);

create table if not exists fuel_nozzles (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    code text not null,
    fuel_type_id uuid references fuel_types(id),
    fuel_kind text not null,
    tank_id uuid references fuel_tanks(id),
    status text not null,
    meter_reading numeric(14,3) not null default 0,
    notes text,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    unique (tenant_id, code)
);

create table if not exists fuel_receipts (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    receipt_number text not null,
    receipt_date timestamptz not null,
    amount numeric(14,2) not null default 0,
    attachment_reference text,
    notes text,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    unique (tenant_id, receipt_number)
);

create table if not exists fuel_transactions (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    transaction_date timestamptz not null,
    fuel_type_id uuid references fuel_types(id),
    fuel_kind text not null,
    quantity numeric(14,3) not null,
    unit_price numeric(14,4) not null default 0,
    total_amount numeric(14,2) not null default 0,
    odometer_reading numeric(14,2),
    engine_hour_reading numeric(14,2),
    station_id uuid references fuel_stations(id),
    tank_id uuid references fuel_tanks(id),
    nozzle_id uuid references fuel_nozzles(id),
    receipt_id uuid references fuel_receipts(id),
    receipt_number text,
    driver_reference text,
    vehicle_reference text,
    asset_reference text,
    payment_method text,
    notes text,
    status text not null,
    cancellation_reason text,
    completed_at timestamptz,
    canceled_at timestamptz,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    check (quantity > 0),
    check (unit_price >= 0),
    check (total_amount >= 0),
    check (coalesce(odometer_reading, 0) >= 0),
    check (coalesce(engine_hour_reading, 0) >= 0)
);

create table if not exists fuel_readings (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    reading_type text not null,
    reference_id text,
    value numeric(14,3) not null,
    reading_date timestamptz not null,
    source text,
    notes text,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    check (value >= 0)
);

create table if not exists fuel_prices (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    fuel_type_id uuid references fuel_types(id),
    fuel_kind text not null,
    unit_price numeric(14,4) not null,
    effective_date timestamptz not null,
    station_id uuid references fuel_stations(id),
    source text,
    notes text,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    check (unit_price >= 0)
);

create table if not exists fuel_attachments (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    reference_id text,
    kind text not null,
    file_name text not null,
    mime_type text,
    uri text,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create table if not exists fuel_adjustments (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    transaction_id uuid not null references fuel_transactions(id),
    adjustment_type text not null,
    reason text not null,
    original_reference text,
    adjusted_quantity numeric(14,3),
    adjusted_unit_price numeric(14,4),
    adjusted_total_amount numeric(14,2),
    notes text,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    check (coalesce(adjusted_quantity, 0) >= 0),
    check (coalesce(adjusted_unit_price, 0) >= 0),
    check (coalesce(adjusted_total_amount, 0) >= 0)
);

create index if not exists fuel_transactions_tenant_status_idx on fuel_transactions (tenant_id, status);
create index if not exists fuel_transactions_tenant_date_idx on fuel_transactions (tenant_id, transaction_date);
create index if not exists fuel_transactions_vehicle_ref_idx on fuel_transactions (tenant_id, vehicle_reference);
create index if not exists fuel_readings_tenant_date_idx on fuel_readings (tenant_id, reading_date);
create index if not exists fuel_prices_tenant_effective_idx on fuel_prices (tenant_id, effective_date);
create index if not exists fuel_adjustments_transaction_idx on fuel_adjustments (tenant_id, transaction_id);
