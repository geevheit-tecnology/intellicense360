-- Sprint 005: Tires Core.

create extension if not exists "uuid-ossp";

create table if not exists tires (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    serial_number text not null,
    fire_number text not null,
    brand text not null default '',
    model text not null default '',
    size text not null default '',
    construction text not null default '',
    tire_type text not null default '',
    position_type text not null default '',
    manufacturing_date date,
    purchase_date date,
    purchase_value numeric(14,2),
    supplier text not null default '',
    dot text not null,
    current_tread_mm numeric(8,2) not null default 0,
    original_tread_mm numeric(8,2) not null default 0,
    minimum_tread_mm numeric(8,2) not null default 0,
    status text not null default 'new',
    vehicle_id uuid,
    position text not null default '',
    current_km bigint not null default 0,
    total_km bigint not null default 0,
    recap_count integer not null default 0,
    notes text not null default '',
    created_by uuid,
    updated_by uuid,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    constraint ck_tires_tread_non_negative check (current_tread_mm >= 0 and original_tread_mm >= 0 and minimum_tread_mm >= 0),
    constraint ck_tires_min_tread_lte_original check (minimum_tread_mm <= original_tread_mm)
);

create table if not exists tire_inspections (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    tire_id uuid not null references tires(id) on delete cascade,
    inspection_date timestamptz not null default now(),
    tread_mm numeric(8,2) not null,
    pressure numeric(8,2),
    temperature numeric(8,2),
    condition text not null default '',
    observations text not null default '',
    inspector text not null default '',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    constraint ck_tire_inspections_tread_non_negative check (tread_mm >= 0)
);

create table if not exists tire_movements (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    tire_id uuid not null references tires(id) on delete cascade,
    movement_type text not null,
    vehicle_id uuid,
    position text not null default '',
    km bigint not null default 0,
    reason text not null default '',
    performed_by uuid,
    movement_date timestamptz not null default now(),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create unique index if not exists ux_tires_serial_number_active on tires(tenant_id, upper(serial_number)) where deleted_at is null;
create unique index if not exists ux_tires_fire_number_active on tires(tenant_id, upper(fire_number)) where deleted_at is null;
create index if not exists ix_tires_tenant on tires(tenant_id) where deleted_at is null;
create index if not exists ix_tires_vehicle on tires(tenant_id, vehicle_id) where deleted_at is null;
create index if not exists ix_tires_status on tires(tenant_id, status) where deleted_at is null;
create index if not exists ix_tire_inspections_tenant_tire on tire_inspections(tenant_id, tire_id) where deleted_at is null;
create index if not exists ix_tire_movements_tenant_tire on tire_movements(tenant_id, tire_id) where deleted_at is null;
