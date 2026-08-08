-- Sprint 004: Checklist Core.

create extension if not exists "uuid-ossp";

create table if not exists checklists (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    vehicle_id uuid not null,
    title text not null,
    description text not null default '',
    type text not null default '',
    status text not null default 'draft',
    started_at timestamptz,
    finished_at timestamptz,
    driver_name text not null default '',
    driver_document text not null default '',
    created_by uuid,
    updated_by uuid,
    deleted_at timestamptz,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table if not exists checklist_items (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    checklist_id uuid not null references checklists(id) on delete cascade,
    title text not null,
    description text not null default '',
    category text not null default '',
    required boolean not null default false,
    order_index integer not null default 0,
    answer_type text not null,
    expected_value text not null default '',
    deleted_at timestamptz,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table if not exists checklist_answers (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    checklist_id uuid not null references checklists(id) on delete cascade,
    checklist_item_id uuid not null references checklist_items(id) on delete cascade,
    answer text not null default '',
    notes text not null default '',
    photo_url text not null default '',
    answered_by uuid,
    answered_at timestamptz not null default now(),
    deleted_at timestamptz,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create index if not exists ix_checklists_tenant on checklists(tenant_id) where deleted_at is null;
create index if not exists ix_checklists_vehicle on checklists(tenant_id, vehicle_id) where deleted_at is null;
create index if not exists ix_checklists_status on checklists(tenant_id, status) where deleted_at is null;
create index if not exists ix_checklists_created_at on checklists(tenant_id, created_at desc) where deleted_at is null;
create index if not exists ix_checklist_items_tenant_checklist on checklist_items(tenant_id, checklist_id) where deleted_at is null;
create index if not exists ix_checklist_answers_tenant_checklist on checklist_answers(tenant_id, checklist_id) where deleted_at is null;
