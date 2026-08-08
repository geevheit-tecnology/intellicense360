-- Sprint 008: Inventory & Parts Core.

create extension if not exists "uuid-ossp";

create table if not exists inventory_categories (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    audit_id uuid,
    name text not null,
    code text not null,
    description text not null default '',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1,
    unique (tenant_id, code)
);

create table if not exists inventory_units (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    audit_id uuid,
    name text not null,
    code text not null,
    description text not null default '',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1,
    unique (tenant_id, code)
);

create table if not exists inventory_brands (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    audit_id uuid,
    name text not null,
    code text not null,
    description text not null default '',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1,
    unique (tenant_id, code)
);

create table if not exists inventory_models (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    audit_id uuid,
    name text not null,
    code text not null,
    description text not null default '',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1,
    unique (tenant_id, code)
);

create table if not exists inventory_parts (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    audit_id uuid,
    sku text not null,
    internal_code text not null,
    name text not null,
    description text not null default '',
    category_id uuid,
    brand_id uuid,
    model_id uuid,
    unit_id text not null,
    status text not null default 'active',
    metadata jsonb not null default '{}'::jsonb,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1
);

create unique index if not exists ux_inventory_parts_tenant_sku
    on inventory_parts(tenant_id, lower(sku))
    where deleted_at is null;

create unique index if not exists ux_inventory_parts_tenant_internal_code
    on inventory_parts(tenant_id, lower(internal_code))
    where deleted_at is null;

create table if not exists inventory_supplier_references (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    part_id uuid not null references inventory_parts(id) on delete cascade,
    supplier_id uuid not null,
    supplier_code text not null default '',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1
);

create table if not exists inventory_warehouses (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    audit_id uuid,
    name text not null,
    code text not null,
    address text not null default '',
    active boolean not null default true,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1,
    unique (tenant_id, code)
);

create table if not exists inventory_locations (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    warehouse_id uuid not null references inventory_warehouses(id) on delete cascade,
    name text not null,
    code text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1,
    unique (tenant_id, warehouse_id, code)
);

create table if not exists inventory_stock_items (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    part_id uuid not null references inventory_parts(id),
    warehouse_id uuid not null references inventory_warehouses(id),
    location_id uuid,
    quantity numeric(14,3) not null default 0,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1
);

create table if not exists inventory_stock_batches (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    part_id uuid not null references inventory_parts(id),
    batch_code text not null,
    expires_at timestamptz,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1
);

create table if not exists inventory_stock_levels (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    part_id uuid not null references inventory_parts(id),
    warehouse_id uuid not null references inventory_warehouses(id),
    quantity numeric(14,3) not null default 0,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1,
    unique (tenant_id, part_id, warehouse_id)
);

create table if not exists inventory_stock_limits (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    part_id uuid not null references inventory_parts(id),
    minimum_quantity numeric(14,3) not null default 0,
    maximum_quantity numeric(14,3) not null default 0,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1,
    check (minimum_quantity <= maximum_quantity)
);

create table if not exists inventory_attachments (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    part_id uuid not null references inventory_parts(id) on delete cascade,
    file_name text not null,
    mime_type text not null default '',
    uri text not null default '',
    created_at timestamptz not null default now()
);

create index if not exists ix_inventory_parts_tenant_status on inventory_parts(tenant_id, status) where deleted_at is null;
create index if not exists ix_inventory_parts_category on inventory_parts(tenant_id, category_id) where deleted_at is null;
create index if not exists ix_inventory_warehouses_tenant on inventory_warehouses(tenant_id, active) where deleted_at is null;
create index if not exists ix_inventory_locations_warehouse on inventory_locations(tenant_id, warehouse_id) where deleted_at is null;
create index if not exists ix_inventory_stock_items_part on inventory_stock_items(tenant_id, part_id) where deleted_at is null;
create index if not exists ix_inventory_stock_levels_part on inventory_stock_levels(tenant_id, part_id) where deleted_at is null;
