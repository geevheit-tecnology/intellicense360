-- Sprint 006: Fleet Assets & Equipment Foundation.

create extension if not exists "uuid-ossp";

create table if not exists asset_categories (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    name text not null,
    code text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1,
    unique (tenant_id, code)
);

create table if not exists asset_types (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    category_id uuid references asset_categories(id),
    name text not null,
    code text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1,
    unique (tenant_id, code)
);

create table if not exists asset_manufacturers (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    name text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1,
    unique (tenant_id, name)
);

create table if not exists asset_models (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    manufacturer_id uuid references asset_manufacturers(id),
    name text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1,
    unique (tenant_id, manufacturer_id, name)
);

create table if not exists assets (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    audit_id uuid,
    internal_code text not null,
    serial_number text,
    asset_tag text not null,
    name text not null,
    description text not null default '',
    category_id uuid references asset_categories(id),
    type_id uuid references asset_types(id),
    manufacturer_id uuid references asset_manufacturers(id),
    model_id uuid references asset_models(id),
    status text not null default 'draft',
    ownership text not null default 'owned',
    location jsonb not null default '{}'::jsonb,
    depreciation jsonb not null default '{}'::jsonb,
    metadata jsonb not null default '{}'::jsonb,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1
);

create table if not exists asset_equipment (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    asset_id uuid not null references assets(id) on delete cascade,
    category text not null default '',
    type text not null default '',
    capacity text not null default '',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1
);

create table if not exists asset_implements (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    asset_id uuid not null references assets(id) on delete cascade,
    type text not null default '',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1
);

create table if not exists asset_attachments (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    asset_id uuid not null references assets(id) on delete cascade,
    kind text not null,
    file_name text not null,
    mime_type text not null default '',
    uri text not null default '',
    created_at timestamptz not null default now()
);

create unique index if not exists ux_assets_internal_code_active on assets(tenant_id, upper(internal_code)) where deleted_at is null;
create unique index if not exists ux_assets_asset_tag_active on assets(tenant_id, upper(asset_tag)) where deleted_at is null;
create unique index if not exists ux_assets_serial_number_active on assets(tenant_id, upper(serial_number)) where deleted_at is null and serial_number is not null and serial_number <> '';
create index if not exists ix_assets_tenant_status on assets(tenant_id, status) where deleted_at is null;
create index if not exists ix_assets_tenant_category on assets(tenant_id, category_id) where deleted_at is null;
create index if not exists ix_assets_tenant_type on assets(tenant_id, type_id) where deleted_at is null;
create index if not exists ix_asset_equipment_tenant_asset on asset_equipment(tenant_id, asset_id) where deleted_at is null;
create index if not exists ix_asset_attachments_tenant_asset on asset_attachments(tenant_id, asset_id);
