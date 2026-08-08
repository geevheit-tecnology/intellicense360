-- Sprint 002: Identity & Platform Foundation.

create extension if not exists "uuid-ossp";
create extension if not exists "citext";

alter table tenants
    add column if not exists slug text,
    add column if not exists metadata jsonb not null default '{}'::jsonb;

create unique index if not exists ux_tenants_slug on tenants(slug) where slug is not null;

create table if not exists users (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    name text not null,
    email citext not null,
    password_hash text not null,
    status text not null default 'active',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1,
    unique (tenant_id, email)
);

create table if not exists roles (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    name text not null,
    description text not null default '',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1,
    unique (tenant_id, name)
);

create table if not exists permissions (
    id uuid primary key default uuid_generate_v4(),
    key text not null unique,
    description text not null default '',
    created_at timestamptz not null default now()
);

create table if not exists role_permissions (
    role_id uuid not null references roles(id) on delete cascade,
    permission_id uuid not null references permissions(id) on delete cascade,
    created_at timestamptz not null default now(),
    primary key (role_id, permission_id)
);

create table if not exists tenant_users (
    tenant_id uuid not null references tenants(id) on delete cascade,
    user_id uuid not null references users(id) on delete cascade,
    role_id uuid references roles(id) on delete set null,
    created_at timestamptz not null default now(),
    primary key (tenant_id, user_id)
);

create table if not exists groups (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    name text not null,
    description text not null default '',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    unique (tenant_id, name)
);

create table if not exists sessions (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    user_id uuid not null references users(id),
    expires_at timestamptz not null,
    revoked_at timestamptz,
    created_at timestamptz not null default now()
);

create table if not exists refresh_tokens (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    user_id uuid not null references users(id),
    session_id uuid not null references sessions(id) on delete cascade,
    token_hash text not null unique,
    expires_at timestamptz not null,
    revoked_at timestamptz,
    replaced_by_id uuid references refresh_tokens(id),
    created_at timestamptz not null default now()
);

create table if not exists audit_logs (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    actor_id uuid references users(id),
    action text not null,
    resource_type text not null,
    resource_id uuid,
    metadata jsonb not null default '{}'::jsonb,
    created_at timestamptz not null default now()
);

create index if not exists ix_users_tenant on users(tenant_id);
create index if not exists ix_roles_tenant on roles(tenant_id);
create index if not exists ix_sessions_user on sessions(user_id);
create index if not exists ix_refresh_tokens_session on refresh_tokens(session_id);
create index if not exists ix_audit_logs_tenant_created on audit_logs(tenant_id, created_at desc);
