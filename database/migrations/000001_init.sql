-- Geevheit Intelligence 360 initial database baseline.
-- Every operational table must include tenant_id, audit metadata and optimistic locking.

create extension if not exists "uuid-ossp";

create table if not exists tenants (
    id uuid primary key default uuid_generate_v4(),
    name text not null,
    active boolean not null default true,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1
);

create table if not exists audit_events (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    actor_id uuid,
    action text not null,
    resource_type text not null,
    resource_id uuid,
    metadata jsonb not null default '{}'::jsonb,
    created_at timestamptz not null default now()
);
