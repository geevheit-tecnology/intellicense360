-- Sprint 007: Maintenance business domain foundation.

create extension if not exists "uuid-ossp";

create table if not exists maintenance_service_types (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    name text not null,
    code text not null,
    description text not null default '',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1,
    unique (tenant_id, code)
);

create table if not exists maintenance_work_orders (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    asset_id uuid,
    vehicle_id uuid,
    title text not null,
    description text not null default '',
    kind text not null,
    status text not null,
    priority text not null,
    service_type_id uuid references maintenance_service_types(id),
    opened_at timestamptz not null default now(),
    scheduled_at timestamptz,
    started_at timestamptz,
    completed_at timestamptz,
    cancelled_at timestamptz,
    estimated_hours numeric(10,2) not null default 0,
    actual_hours numeric(10,2) not null default 0,
    created_by uuid,
    updated_by uuid,
    metadata jsonb not null default '{}'::jsonb,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1
);

create table if not exists maintenance_preventive_plans (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    name text not null,
    description text not null default '',
    asset_id uuid,
    vehicle_id uuid,
    service_type_id uuid references maintenance_service_types(id),
    frequency text not null,
    interval_value bigint not null,
    next_due_at timestamptz,
    next_due_value bigint not null default 0,
    active boolean not null default true,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1
);

create table if not exists maintenance_corrective_records (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    work_order_id uuid not null references maintenance_work_orders(id) on delete cascade,
    failure_mode text not null,
    root_cause text not null default '',
    resolution text not null default '',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create table if not exists maintenance_labor_entries (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    work_order_id uuid not null references maintenance_work_orders(id) on delete cascade,
    technician text not null,
    hours numeric(10,2) not null,
    cost numeric(14,2) not null default 0,
    started_at timestamptz not null default now(),
    finished_at timestamptz,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create table if not exists maintenance_downtime (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    work_order_id uuid not null references maintenance_work_orders(id) on delete cascade,
    asset_id uuid,
    vehicle_id uuid,
    started_at timestamptz not null default now(),
    ended_at timestamptz,
    reason text not null default '',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create table if not exists maintenance_availability_snapshots (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    asset_id uuid,
    vehicle_id uuid,
    available boolean not null,
    availability_pct numeric(5,2) not null,
    captured_at timestamptz not null default now(),
    created_at timestamptz not null default now()
);

create table if not exists maintenance_history (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    work_order_id uuid not null references maintenance_work_orders(id) on delete cascade,
    event text not null,
    actor_id uuid,
    notes text not null default '',
    created_at timestamptz not null default now()
);

create index if not exists ix_maintenance_work_orders_tenant_status on maintenance_work_orders(tenant_id, status) where deleted_at is null;
create index if not exists ix_maintenance_work_orders_asset on maintenance_work_orders(tenant_id, asset_id) where deleted_at is null;
create index if not exists ix_maintenance_work_orders_vehicle on maintenance_work_orders(tenant_id, vehicle_id) where deleted_at is null;
create index if not exists ix_maintenance_plans_tenant_active on maintenance_preventive_plans(tenant_id, active) where deleted_at is null;
create index if not exists ix_maintenance_labor_work_order on maintenance_labor_entries(tenant_id, work_order_id) where deleted_at is null;
create index if not exists ix_maintenance_downtime_work_order on maintenance_downtime(tenant_id, work_order_id) where deleted_at is null;
create index if not exists ix_maintenance_history_work_order on maintenance_history(tenant_id, work_order_id, created_at);
