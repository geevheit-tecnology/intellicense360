create table if not exists mission_control_command_items (
    id uuid primary key,
    tenant_id uuid not null references tenants(id),
    type text not null,
    category text not null,
    severity text not null,
    priority text not null,
    status text not null,
    title text not null,
    description text,
    source text,
    source_type text,
    source_id text,
    confidence numeric(6,5) not null default 0,
    impact_score numeric(6,5) not null default 0,
    risk_score numeric(6,5) not null default 0,
    urgency_score numeric(6,5) not null default 0,
    due_at timestamptz,
    sla_status text not null default 'not_applicable',
    sla_hours integer not null default 0,
    detected_at timestamptz not null default now(),
    acknowledged_at timestamptz,
    resolved_at timestamptz,
    assigned_to uuid,
    fingerprint text not null,
    metadata jsonb not null default '{}'::jsonb,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1,
    constraint mission_control_command_items_type_chk check (type in ('alert', 'risk', 'incident', 'opportunity', 'recommendation', 'task', 'anomaly', 'warning', 'insight')),
    constraint mission_control_command_items_category_chk check (category in ('operational', 'maintenance', 'fleet', 'tire', 'fuel', 'inventory', 'financial', 'compliance', 'driver', 'document', 'ciot', 'safety', 'performance', 'cost')),
    constraint mission_control_command_items_severity_chk check (severity in ('critical', 'high', 'medium', 'low', 'info')),
    constraint mission_control_command_items_priority_chk check (priority in ('low', 'normal', 'high', 'urgent', 'critical')),
    constraint mission_control_command_items_status_chk check (status in ('open', 'acknowledged', 'in_progress', 'resolved', 'dismissed', 'expired')),
    constraint mission_control_command_items_sla_chk check (sla_status in ('within_sla', 'at_risk', 'breached', 'not_applicable')),
    constraint mission_control_command_items_scores_chk check (
        confidence between 0 and 1 and impact_score between 0 and 1 and risk_score between 0 and 1 and urgency_score between 0 and 1
    )
);

create unique index if not exists mission_control_command_items_active_fingerprint_idx
    on mission_control_command_items (tenant_id, fingerprint)
    where status not in ('resolved', 'dismissed') and deleted_at is null;

create index if not exists mission_control_command_items_tenant_status_idx
    on mission_control_command_items (tenant_id, status, severity, priority);

create table if not exists mission_control_command_actions (
    id uuid primary key,
    tenant_id uuid not null references tenants(id),
    command_item_id uuid not null references mission_control_command_items(id),
    type text not null,
    label text not null,
    status text not null default 'pending',
    priority text not null default 'normal',
    payload jsonb not null default '{}'::jsonb,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    constraint mission_control_command_actions_type_chk check (type in ('acknowledge', 'assign', 'start', 'resolve', 'dismiss', 'escalate', 'review')),
    constraint mission_control_command_actions_status_chk check (status in ('pending', 'completed', 'canceled')),
    constraint mission_control_command_actions_priority_chk check (priority in ('low', 'normal', 'high', 'urgent', 'critical'))
);

create table if not exists mission_control_command_events (
    id uuid primary key,
    tenant_id uuid not null references tenants(id),
    command_item_id uuid not null references mission_control_command_items(id),
    event_type text not null,
    previous_status text,
    new_status text,
    actor_id uuid,
    payload jsonb not null default '{}'::jsonb,
    occurred_at timestamptz not null default now()
);

create index if not exists mission_control_command_events_item_idx
    on mission_control_command_events (tenant_id, command_item_id, occurred_at desc);

create table if not exists mission_control_snapshots (
    id uuid primary key,
    tenant_id uuid not null references tenants(id),
    snapshot_at timestamptz not null default now(),
    open_items integer not null default 0,
    critical_items integer not null default 0,
    high_priority_items integer not null default 0,
    active_risks integer not null default 0,
    active_alerts integer not null default 0,
    open_incidents integer not null default 0,
    opportunities integer not null default 0,
    breached_slas integer not null default 0,
    average_resolution_time numeric(12,2) not null default 0,
    operational_score numeric(6,2) not null default 100,
    risk_score numeric(6,5) not null default 0,
    health_score numeric(6,2) not null default 100
);

create index if not exists mission_control_snapshots_latest_idx
    on mission_control_snapshots (tenant_id, snapshot_at desc);

create table if not exists mission_control_idempotency (
    tenant_id uuid not null references tenants(id),
    idempotency_key text not null,
    resource_id uuid not null,
    created_at timestamptz not null default now(),
    primary key (tenant_id, idempotency_key)
);
