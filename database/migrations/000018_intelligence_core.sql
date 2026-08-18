create table if not exists intelligence_metrics (
    id uuid primary key,
    tenant_id uuid not null,
    metric_type text not null,
    name text not null,
    value numeric(18,4) not null,
    unit text,
    period_start timestamptz not null,
    period_end timestamptz not null,
    dimension text,
    dimension_value text,
    source text,
    calculated_at timestamptz not null,
    confidence numeric(5,4) not null default 0,
    metadata jsonb not null default '{}'::jsonb,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    version bigint not null default 1
);
create index if not exists idx_intelligence_metrics_tenant_type_period on intelligence_metrics (tenant_id, metric_type, period_start, period_end);
create index if not exists idx_intelligence_metrics_dimension on intelligence_metrics (tenant_id, dimension, dimension_value);

create table if not exists intelligence_anomalies (
    id uuid primary key,
    tenant_id uuid not null,
    type text not null,
    severity text not null,
    metric_id uuid,
    observed_value numeric(18,4) not null,
    expected_value numeric(18,4) not null,
    deviation numeric(18,4) not null,
    deviation_percentage numeric(10,4) not null,
    detected_at timestamptz not null,
    period_start timestamptz not null,
    period_end timestamptz not null,
    evidence jsonb not null default '[]'::jsonb,
    confidence numeric(5,4) not null default 0,
    status text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    version bigint not null default 1
);
create index if not exists idx_intelligence_anomalies_tenant_status on intelligence_anomalies (tenant_id, status, detected_at);
create index if not exists idx_intelligence_anomalies_tenant_severity on intelligence_anomalies (tenant_id, severity, detected_at);

create table if not exists intelligence_risks (
    id uuid primary key,
    tenant_id uuid not null,
    category text not null,
    severity text not null,
    probability numeric(5,4) not null default 0,
    impact jsonb not null default '{}'::jsonb,
    confidence numeric(5,4) not null default 0,
    evidence jsonb not null default '[]'::jsonb,
    detected_at timestamptz not null,
    status text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    version bigint not null default 1
);
create index if not exists idx_intelligence_risks_tenant_category on intelligence_risks (tenant_id, category, detected_at);

create table if not exists intelligence_opportunities (
    id uuid primary key,
    tenant_id uuid not null,
    category text not null,
    estimated_impact jsonb not null default '{}'::jsonb,
    estimated_saving numeric(18,2),
    confidence numeric(5,4) not null default 0,
    evidence jsonb not null default '[]'::jsonb,
    detected_at timestamptz not null,
    valid_until timestamptz,
    status text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    version bigint not null default 1
);
create index if not exists idx_intelligence_opportunities_tenant_status on intelligence_opportunities (tenant_id, status, detected_at);

create table if not exists intelligence_recommendations (
    id uuid primary key,
    tenant_id uuid not null,
    title text not null,
    description text,
    what_happened text,
    why_it_matters text,
    evidence jsonb not null default '[]'::jsonb,
    suggested_action text not null,
    expected_impact jsonb not null default '{}'::jsonb,
    confidence numeric(5,4) not null default 0,
    priority text not null,
    impact_area text,
    status text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    version bigint not null default 1
);
create index if not exists idx_intelligence_recommendations_tenant_priority on intelligence_recommendations (tenant_id, priority, created_at);

create table if not exists intelligence_insights (
    id uuid primary key,
    tenant_id uuid not null,
    type text not null,
    title text not null,
    summary text not null,
    category text not null,
    severity text not null,
    evidence jsonb not null default '[]'::jsonb,
    metric_id uuid,
    anomaly_id uuid,
    risk_id uuid,
    opportunity_id uuid,
    recommendation_id uuid,
    estimated_impact jsonb not null default '{}'::jsonb,
    confidence numeric(5,4) not null default 0,
    priority text not null,
    deduplication_key text not null,
    status text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    version bigint not null default 1
);
create unique index if not exists idx_intelligence_insights_dedup on intelligence_insights (tenant_id, deduplication_key);
create index if not exists idx_intelligence_insights_status on intelligence_insights (tenant_id, status, created_at);

create table if not exists intelligence_evidence (
    id uuid primary key,
    tenant_id uuid not null,
    insight_id uuid,
    source text not null,
    source_type text not null,
    reference_id text,
    metric text,
    observed_value numeric(18,4),
    expected_value numeric(18,4),
    period_start timestamptz,
    period_end timestamptz,
    explanation text not null,
    created_at timestamptz not null default now()
);
create index if not exists idx_intelligence_evidence_tenant_insight on intelligence_evidence (tenant_id, insight_id);

create table if not exists intelligence_rule_definitions (
    id uuid primary key,
    tenant_id uuid,
    name text not null,
    rule_type text not null,
    rule_version integer not null,
    threshold numeric(18,4),
    window text,
    active boolean not null default true,
    explanation text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);
create index if not exists idx_intelligence_rules_tenant_type on intelligence_rule_definitions (tenant_id, rule_type, active);

create table if not exists intelligence_read_models (
    id uuid primary key,
    tenant_id uuid not null,
    source text not null,
    source_type text not null,
    reference_id text not null,
    metric_name text not null,
    value numeric(18,4) not null,
    unit text,
    period_start timestamptz not null,
    period_end timestamptz not null,
    metadata jsonb not null default '{}'::jsonb,
    last_event_id uuid,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    version bigint not null default 1
);
create unique index if not exists idx_intelligence_read_models_identity
    on intelligence_read_models (tenant_id, source_type, reference_id, metric_name);
create index if not exists idx_intelligence_read_models_event on intelligence_read_models (tenant_id, last_event_id);
