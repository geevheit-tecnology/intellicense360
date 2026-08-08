-- Sprint 003: Fleet Core Foundation.

create extension if not exists "uuid-ossp";

create table if not exists fleet_vehicle_categories (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    audit_id uuid,
    name text not null,
    code text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1,
    unique (tenant_id, code)
);

create table if not exists fleet_vehicle_types (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    audit_id uuid,
    category_id uuid not null references fleet_vehicle_categories(id),
    name text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1,
    unique (tenant_id, category_id, name)
);

create table if not exists fleet_vehicle_brands (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    audit_id uuid,
    name text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1,
    unique (tenant_id, name)
);

create table if not exists fleet_vehicle_models (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    audit_id uuid,
    brand_id uuid not null references fleet_vehicle_brands(id),
    category_id uuid references fleet_vehicle_categories(id),
    type_id uuid references fleet_vehicle_types(id),
    name text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1,
    unique (tenant_id, brand_id, name)
);

create table if not exists fleet_assets (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    audit_id uuid,
    code text not null,
    name text not null,
    description text not null default '',
    status text not null default 'active',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1,
    unique (tenant_id, code)
);

create table if not exists fleet_vehicles (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    audit_id uuid,
    category_id uuid not null references fleet_vehicle_categories(id),
    type_id uuid not null references fleet_vehicle_types(id),
    brand_id uuid not null references fleet_vehicle_brands(id),
    model_id uuid not null references fleet_vehicle_models(id),
    asset_id uuid references fleet_assets(id),
    license_plate text not null,
    renavam text not null,
    chassis text not null,
    engine text,
    color text,
    fuel_type text,
    transmission text,
    axle_configuration text,
    emission_standard text,
    ownership_type text not null,
    status text not null default 'draft',
    year_manufacture integer,
    year_model integer,
    metadata jsonb not null default '{}'::jsonb,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1
);

create unique index if not exists ux_fleet_vehicles_plate_active on fleet_vehicles(tenant_id, upper(replace(license_plate, '-', ''))) where deleted_at is null;
create unique index if not exists ux_fleet_vehicles_chassis_active on fleet_vehicles(tenant_id, upper(chassis)) where deleted_at is null;
create unique index if not exists ux_fleet_vehicles_renavam_active on fleet_vehicles(tenant_id, renavam) where deleted_at is null;
create index if not exists ix_fleet_vehicles_tenant_status on fleet_vehicles(tenant_id, status) where deleted_at is null;
create index if not exists ix_fleet_vehicles_tenant_category on fleet_vehicles(tenant_id, category_id) where deleted_at is null;
