-- Sprint 010: Tire Lifecycle Core.

alter table tires
    add column if not exists brand_id uuid,
    add column if not exists model_id uuid,
    add column if not exists dimension text not null default '',
    add column if not exists load_index text not null default '',
    add column if not exists speed_rating text not null default '',
    add column if not exists supplier_reference text not null default '',
    add column if not exists warranty text not null default '',
    add column if not exists condition text not null default '',
    add column if not exists version bigint not null default 1;

create table if not exists tire_brands (
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

create table if not exists tire_models (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    audit_id uuid,
    brand_id uuid,
    name text not null,
    code text not null,
    description text not null default '',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1,
    unique (tenant_id, code)
);

create table if not exists tire_specifications (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    dimension text not null,
    construction text not null default '',
    load_index text not null default '',
    speed_rating text not null default '',
    original_tread_depth_mm numeric(6,2) not null default 0,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1
);

create table if not exists tire_positions (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    axle text not null default '',
    side text not null default '',
    position text not null default '',
    inner_outer text not null default '',
    position_code text not null,
    position_label text not null default '',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1,
    unique (tenant_id, position_code)
);

create table if not exists tire_installations (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    tire_id uuid not null references tires(id),
    installation_date timestamptz not null,
    position_code text not null default '',
    initial_km bigint not null default 0,
    initial_tread_depth numeric(6,2) not null default 0,
    installation_reason text not null default '',
    notes text not null default '',
    created_at timestamptz not null default now()
);

create table if not exists tire_removals (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    tire_id uuid not null references tires(id),
    removal_date timestamptz not null,
    position_code text not null default '',
    removal_km bigint not null default 0,
    remaining_tread_depth numeric(6,2) not null default 0,
    removal_reason text not null default '',
    condition text not null default '',
    notes text not null default '',
    created_at timestamptz not null default now()
);

create table if not exists tire_measurements (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    tire_id uuid not null references tires(id),
    tread_depth_mm numeric(6,2) not null,
    pressure numeric(8,2) not null default 0,
    measurement_position text not null default '',
    measurement_date timestamptz not null,
    created_at timestamptz not null default now()
);

create table if not exists tire_retreads (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    tire_id uuid not null references tires(id),
    retread_number integer not null,
    retread_date timestamptz not null,
    provider_reference text not null default '',
    tread_brand text not null default '',
    tread_model text not null default '',
    cost numeric(14,2) not null default 0,
    before_retread_tread_depth numeric(6,2) not null default 0,
    after_retread_tread_depth numeric(6,2) not null default 0,
    warranty text not null default '',
    result text not null default '',
    status text not null default '',
    notes text not null default '',
    created_at timestamptz not null default now()
);

create table if not exists tire_retread_events (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    tire_id uuid not null references tires(id),
    retread_id uuid not null references tire_retreads(id),
    event text not null,
    notes text not null default '',
    created_at timestamptz not null default now()
);

create table if not exists tire_history (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    tire_id uuid not null references tires(id),
    event text not null,
    from_status text not null default '',
    to_status text not null default '',
    actor_id uuid,
    notes text not null default '',
    created_at timestamptz not null default now()
);

create table if not exists tire_costs (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    tire_id uuid not null references tires(id),
    purchase_cost numeric(14,2) not null default 0,
    retread_cost numeric(14,2) not null default 0,
    repair_cost numeric(14,2) not null default 0,
    other_cost numeric(14,2) not null default 0,
    total_cost numeric(14,2) not null default 0,
    cost_per_km numeric(14,6) not null default 0,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1
);

create table if not exists tire_disposals (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    tire_id uuid not null references tires(id),
    disposal_date timestamptz not null,
    reason text not null default '',
    attachment_reference text not null default '',
    created_at timestamptz not null default now()
);

create table if not exists tire_attachments (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    tire_id uuid not null references tires(id),
    kind text not null,
    file_name text not null,
    mime_type text not null default '',
    uri text not null default '',
    created_at timestamptz not null default now()
);

create index if not exists ix_tire_history_tire on tire_history(tenant_id, tire_id, created_at);
create index if not exists ix_tire_movements_immutable on tire_movements(tenant_id, tire_id, movement_date);
create index if not exists ix_tire_measurements_tire on tire_measurements(tenant_id, tire_id, measurement_date);
create index if not exists ix_tire_retreads_tire on tire_retreads(tenant_id, tire_id, retread_date);
