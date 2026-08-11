create table if not exists checklist_types (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    name text not null,
    code text not null,
    description text,
    active boolean not null default true,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    unique (tenant_id, code)
);

create table if not exists checklist_templates (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    name text not null,
    description text,
    type_id uuid references checklist_types(id),
    type text,
    status text not null,
    active boolean not null default true,
    current_version_number integer not null default 0,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create table if not exists checklist_template_versions (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    template_id uuid not null references checklist_templates(id),
    version_number integer not null,
    status text not null,
    instructions text,
    scoring_config jsonb,
    severity_config jsonb,
    evidence_required boolean not null default false,
    signature_required boolean not null default false,
    published_at timestamptz,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    unique (tenant_id, template_id, version_number)
);

create table if not exists checklist_sections (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    template_version_id uuid not null references checklist_template_versions(id),
    name text not null,
    description text,
    order_index integer not null default 0,
    active boolean not null default true,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create table if not exists checklist_engine_items (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    template_version_id uuid not null references checklist_template_versions(id),
    section_id uuid references checklist_sections(id),
    question text not null,
    description text,
    item_type text not null,
    required boolean not null default false,
    order_index integer not null default 0,
    help_text text,
    severity text,
    evidence_required boolean not null default false,
    active boolean not null default true,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create table if not exists checklist_item_options (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    item_id uuid not null references checklist_engine_items(id),
    label text not null,
    value text not null,
    order_index integer not null default 0,
    active boolean not null default true,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create table if not exists checklist_executions (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    template_version_id uuid not null references checklist_template_versions(id),
    status text not null,
    performed_by text,
    location_reference text,
    notes text,
    score numeric(10,2),
    final_result text,
    started_at timestamptz,
    ended_at timestamptz,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create table if not exists checklist_responses (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    execution_id uuid not null references checklist_executions(id),
    item_id uuid not null references checklist_engine_items(id),
    value text not null,
    result text,
    notes text,
    responder text,
    severity text,
    responded_at timestamptz not null,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create table if not exists checklist_evidence (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    execution_id uuid not null references checklist_executions(id),
    response_id text,
    evidence_type text not null,
    reference text not null,
    file_name text,
    mime_type text,
    notes text,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create table if not exists checklist_non_conformities (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    execution_id uuid not null references checklist_executions(id),
    response_id text,
    title text not null,
    description text,
    severity text not null,
    status text not null,
    recommendation text,
    resolved_at timestamptz,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create table if not exists checklist_signatures (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    execution_id uuid not null references checklist_executions(id),
    signer text not null,
    signature_reference text not null,
    signature_type text not null,
    signed_at timestamptz not null,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create table if not exists checklist_assignments (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    template_id uuid not null references checklist_templates(id),
    assigned_to_reference text,
    target_reference text,
    status text not null,
    notes text,
    due_at timestamptz,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create table if not exists checklist_history (
    id uuid primary key default gen_random_uuid(),
    tenant_id uuid not null references tenants(id),
    execution_id uuid not null references checklist_executions(id),
    event text not null,
    actor_id text,
    notes text,
    audit_id uuid,
    version bigint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create index if not exists checklist_templates_tenant_status_idx on checklist_templates (tenant_id, status);
create index if not exists checklist_template_versions_template_idx on checklist_template_versions (tenant_id, template_id);
create index if not exists checklist_engine_items_version_idx on checklist_engine_items (tenant_id, template_version_id);
create index if not exists checklist_executions_tenant_status_idx on checklist_executions (tenant_id, status);
create index if not exists checklist_responses_execution_idx on checklist_responses (tenant_id, execution_id);
create index if not exists checklist_history_execution_idx on checklist_history (tenant_id, execution_id);
