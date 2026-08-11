create table if not exists ciot_ciots (
    id uuid primary key,
    tenant_id uuid not null references tenants(id),
    audit_id uuid,
    ciot_number text,
    type text not null,
    status text not null,
    contract_id uuid,
    carrier_id uuid,
    transporter_id uuid,
    operation_id uuid,
    vehicle_reference_id uuid,
    driver_reference_id uuid,
    amount_id uuid,
    start_date timestamptz not null,
    expected_end_date timestamptz,
    actual_end_date timestamptz,
    operational_period text,
    contract_reference text,
    external_protocol text,
    idempotency_key text,
    request_hash text,
    notes text,
    error_code text,
    error_message text,
    created_by uuid,
    updated_by uuid,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1,
    constraint ciot_ciots_type_chk check (type in ('tac_agregado', 'tac_independente', 'other')),
    constraint ciot_ciots_status_chk check (status in ('draft', 'pending', 'generated', 'active', 'suspended', 'closed', 'canceled', 'error')),
    constraint ciot_ciots_period_chk check (expected_end_date is null or expected_end_date >= start_date),
    constraint ciot_ciots_actual_chk check (actual_end_date is null or actual_end_date >= start_date)
);

create unique index if not exists idx_ciot_ciots_tenant_idempotency
    on ciot_ciots (tenant_id, idempotency_key)
    where idempotency_key is not null and deleted_at is null;
create index if not exists idx_ciot_ciots_tenant_status on ciot_ciots (tenant_id, status) where deleted_at is null;
create index if not exists idx_ciot_ciots_tenant_type on ciot_ciots (tenant_id, type) where deleted_at is null;

create table if not exists ciot_contracts (
    id uuid primary key,
    tenant_id uuid not null references tenants(id),
    audit_id uuid,
    contract_number text not null,
    contract_type text,
    start_date timestamptz not null,
    end_date timestamptz,
    status text not null,
    notes text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1
);
create index if not exists idx_ciot_contracts_tenant on ciot_contracts (tenant_id) where deleted_at is null;

create table if not exists ciot_carriers (
    id uuid primary key,
    tenant_id uuid not null references tenants(id),
    audit_id uuid,
    document text not null,
    legal_name text not null,
    trade_name text,
    registration text,
    contact_reference text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1
);
create index if not exists idx_ciot_carriers_tenant_document on ciot_carriers (tenant_id, document) where deleted_at is null;

create table if not exists ciot_transporters (
    id uuid primary key,
    tenant_id uuid not null references tenants(id),
    audit_id uuid,
    document text not null,
    name text not null,
    registration text,
    contract_reference text,
    contact_reference text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1
);
create index if not exists idx_ciot_transporters_tenant_document on ciot_transporters (tenant_id, document) where deleted_at is null;

create table if not exists ciot_operations (
    id uuid primary key,
    tenant_id uuid not null references tenants(id),
    audit_id uuid,
    operation_number text not null,
    origin text,
    destination text,
    start_date timestamptz not null,
    expected_end_date timestamptz,
    actual_end_date timestamptz,
    cargo_reference text,
    weight numeric(14,3),
    distance numeric(14,3),
    notes text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1
);
create index if not exists idx_ciot_operations_tenant_number on ciot_operations (tenant_id, operation_number) where deleted_at is null;

create table if not exists ciot_vehicle_references (
    id uuid primary key,
    tenant_id uuid not null references tenants(id),
    audit_id uuid,
    vehicle_id text,
    license_plate text,
    vehicle_type text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1
);
create index if not exists idx_ciot_vehicle_refs_tenant_plate on ciot_vehicle_references (tenant_id, license_plate) where deleted_at is null;

create table if not exists ciot_driver_references (
    id uuid primary key,
    tenant_id uuid not null references tenants(id),
    audit_id uuid,
    driver_id text,
    name_reference text,
    document_reference text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1
);
create index if not exists idx_ciot_driver_refs_tenant_doc on ciot_driver_references (tenant_id, document_reference) where deleted_at is null;

create table if not exists ciot_amounts (
    id uuid primary key,
    tenant_id uuid not null references tenants(id),
    audit_id uuid,
    ciot_id uuid,
    freight_amount numeric(14,2) not null default 0,
    advance_amount numeric(14,2) not null default 0,
    balance_amount numeric(14,2) not null default 0,
    toll_amount numeric(14,2) not null default 0,
    other_amount numeric(14,2) not null default 0,
    total_amount numeric(14,2) not null default 0,
    currency text not null default 'BRL',
    created_at timestamptz not null default now()
);
create index if not exists idx_ciot_amounts_tenant_ciot on ciot_amounts (tenant_id, ciot_id);

create table if not exists ciot_payments (
    id uuid primary key,
    tenant_id uuid not null references tenants(id),
    audit_id uuid,
    ciot_id uuid not null,
    payment_type text not null,
    amount numeric(14,2) not null,
    due_date timestamptz not null,
    payment_date timestamptz,
    status text not null,
    reference text,
    notes text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1
);
create index if not exists idx_ciot_payments_tenant_ciot on ciot_payments (tenant_id, ciot_id) where deleted_at is null;

create table if not exists ciot_status_history (
    id uuid primary key,
    tenant_id uuid not null references tenants(id),
    audit_id uuid,
    ciot_id uuid not null,
    event text not null,
    from_status text,
    to_status text,
    reason text,
    actor_id uuid,
    metadata jsonb not null default '{}'::jsonb,
    created_at timestamptz not null default now()
);
create index if not exists idx_ciot_history_tenant_ciot_created on ciot_status_history (tenant_id, ciot_id, created_at);

create table if not exists ciot_provider_attempts (
    id uuid primary key,
    tenant_id uuid not null references tenants(id),
    audit_id uuid,
    ciot_id uuid not null,
    provider text not null,
    attempt_number integer not null,
    requested_at timestamptz not null,
    completed_at timestamptz,
    status text not null,
    request_reference text,
    response_reference text,
    error_code text,
    error_message text,
    http_status integer,
    latency_ms bigint,
    idempotency_key text,
    request_hash text,
    result_reference text,
    created_at timestamptz not null default now()
);
create unique index if not exists idx_ciot_attempts_tenant_idempotency
    on ciot_provider_attempts (tenant_id, idempotency_key)
    where idempotency_key is not null;
create index if not exists idx_ciot_attempts_tenant_ciot on ciot_provider_attempts (tenant_id, ciot_id);

create table if not exists ciot_external_references (
    id uuid primary key,
    tenant_id uuid not null references tenants(id),
    audit_id uuid,
    ciot_id uuid not null,
    provider text not null,
    external_ciot_number text,
    external_protocol text,
    external_status text,
    generated_at timestamptz,
    last_synchronized_at timestamptz,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    version bigint not null default 1
);
create unique index if not exists idx_ciot_external_refs_tenant_ciot on ciot_external_references (tenant_id, ciot_id);

create table if not exists ciot_documents (
    id uuid primary key,
    tenant_id uuid not null references tenants(id),
    audit_id uuid,
    ciot_id uuid not null,
    type text not null,
    number text,
    reference text,
    notes text,
    created_at timestamptz not null default now()
);
create index if not exists idx_ciot_documents_tenant_ciot on ciot_documents (tenant_id, ciot_id);

create table if not exists ciot_errors (
    id uuid primary key,
    tenant_id uuid not null references tenants(id),
    audit_id uuid,
    ciot_id uuid not null,
    code text not null,
    message text not null,
    details text,
    created_at timestamptz not null default now()
);
create index if not exists idx_ciot_errors_tenant_ciot on ciot_errors (tenant_id, ciot_id);
