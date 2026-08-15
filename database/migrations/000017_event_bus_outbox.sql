create table if not exists event_outbox (
    id uuid primary key,
    tenant_id uuid not null,
    event_id uuid not null,
    event_type text not null,
    aggregate_id text not null,
    aggregate_type text not null,
    payload jsonb not null,
    occurred_at timestamptz not null,
    status text not null,
    attempts integer not null default 0,
    available_at timestamptz not null default now(),
    published_at timestamptz,
    last_error text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    constraint event_outbox_status_chk check (status in ('pending', 'processing', 'published', 'failed', 'dead_letter'))
);

create unique index if not exists idx_event_outbox_tenant_event
    on event_outbox (tenant_id, event_id);
create index if not exists idx_event_outbox_pending
    on event_outbox (status, available_at, created_at)
    where status in ('pending', 'failed');
create index if not exists idx_event_outbox_tenant_type
    on event_outbox (tenant_id, event_type, occurred_at);

create table if not exists event_processing (
    id uuid primary key,
    tenant_id uuid not null,
    event_id uuid not null,
    event_type text not null,
    consumer_name text not null,
    processed_at timestamptz not null default now(),
    created_at timestamptz not null default now()
);

create unique index if not exists idx_event_processing_consumer_event
    on event_processing (tenant_id, consumer_name, event_id);
create index if not exists idx_event_processing_tenant_type
    on event_processing (tenant_id, event_type, processed_at);

create table if not exists event_dead_letters (
    id uuid primary key,
    tenant_id uuid not null,
    event_id uuid not null,
    event_type text not null,
    payload jsonb not null,
    failure_reason text not null,
    attempts integer not null,
    last_error text,
    created_at timestamptz not null default now()
);

create index if not exists idx_event_dead_letters_tenant_type
    on event_dead_letters (tenant_id, event_type, created_at);
