-- Sprint 007 alignment: Maintenance catalog, order code, workshops, technicians, scheduling, attachments.

alter table maintenance_work_orders
    add column if not exists code text,
    add column if not exists category_id uuid,
    add column if not exists workshop_id uuid,
    add column if not exists technician_id uuid,
    add column if not exists reason_id uuid;

update maintenance_work_orders
set code = 'MO-' || upper(left(replace(id::text, '-', ''), 8))
where code is null or code = '';

alter table maintenance_work_orders
    alter column code set not null;

create unique index if not exists ux_maintenance_work_orders_tenant_code
    on maintenance_work_orders(tenant_id, code)
    where deleted_at is null;

create table if not exists maintenance_categories (
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

create table if not exists maintenance_types (
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

create table if not exists maintenance_priorities (
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

create table if not exists maintenance_reasons (
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

create table if not exists maintenance_workshops (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    name text not null,
    document text not null default '',
    phone text not null default '',
    email text not null default '',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1
);

create table if not exists maintenance_technicians (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    name text not null,
    document text not null default '',
    phone text not null default '',
    email text not null default '',
    active boolean not null default true,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1
);

create table if not exists maintenance_service_items (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    work_order_id uuid not null references maintenance_work_orders(id) on delete cascade,
    description text not null,
    quantity numeric(10,2) not null default 1,
    unit_cost numeric(14,2) not null default 0,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create table if not exists maintenance_schedules (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    work_order_id uuid references maintenance_work_orders(id) on delete cascade,
    preventive_plan_id uuid references maintenance_preventive_plans(id) on delete cascade,
    scheduled_at timestamptz not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    check (work_order_id is not null or preventive_plan_id is not null)
);

create table if not exists maintenance_attachments (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    work_order_id uuid not null references maintenance_work_orders(id) on delete cascade,
    kind text not null,
    file_name text not null,
    mime_type text not null default '',
    uri text not null default '',
    created_at timestamptz not null default now()
);

create index if not exists ix_maintenance_categories_tenant on maintenance_categories(tenant_id) where deleted_at is null;
create index if not exists ix_maintenance_types_tenant on maintenance_types(tenant_id) where deleted_at is null;
create index if not exists ix_maintenance_priorities_tenant on maintenance_priorities(tenant_id) where deleted_at is null;
create index if not exists ix_maintenance_reasons_tenant on maintenance_reasons(tenant_id) where deleted_at is null;
create index if not exists ix_maintenance_workshops_tenant on maintenance_workshops(tenant_id) where deleted_at is null;
create index if not exists ix_maintenance_technicians_tenant on maintenance_technicians(tenant_id, active) where deleted_at is null;
create index if not exists ix_maintenance_service_items_work_order on maintenance_service_items(tenant_id, work_order_id) where deleted_at is null;
create index if not exists ix_maintenance_schedules_tenant_scheduled on maintenance_schedules(tenant_id, scheduled_at) where deleted_at is null;
create index if not exists ix_maintenance_attachments_work_order on maintenance_attachments(tenant_id, work_order_id);
