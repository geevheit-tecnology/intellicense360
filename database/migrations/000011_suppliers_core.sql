-- Sprint 009: Suppliers Core.

create extension if not exists "uuid-ossp";

create table if not exists supplier_categories (
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

create table if not exists supplier_types (
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

create table if not exists suppliers (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    audit_id uuid,
    legal_name text not null,
    trade_name text not null default '',
    cnpj text,
    state_registration text not null default '',
    municipal_registration text not null default '',
    email text not null default '',
    phone text not null default '',
    website text not null default '',
    notes text not null default '',
    status text not null default 'draft',
    category_id uuid,
    type text not null default 'other',
    metadata jsonb not null default '{}'::jsonb,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1
);

create unique index if not exists ux_suppliers_tenant_cnpj
    on suppliers(tenant_id, cnpj)
    where deleted_at is null and cnpj is not null and cnpj <> '';

create table if not exists supplier_contacts (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    supplier_id uuid not null references suppliers(id) on delete cascade,
    name text not null,
    role text not null default '',
    email text not null default '',
    phone text not null default '',
    mobile text not null default '',
    primary_contact boolean not null default false,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1
);

create table if not exists supplier_addresses (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    supplier_id uuid not null references suppliers(id) on delete cascade,
    street text not null,
    number text not null default '',
    complement text not null default '',
    neighborhood text not null default '',
    city text not null,
    state text not null default '',
    postal_code text not null default '',
    country text not null default 'BR',
    address_type text not null default '',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1
);

create table if not exists supplier_bank_accounts (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    supplier_id uuid not null references suppliers(id) on delete cascade,
    bank text not null,
    branch text not null default '',
    account text not null default '',
    account_type text not null default '',
    pix_key text not null default '',
    holder_name text not null default '',
    holder_document text not null default '',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1
);

create table if not exists supplier_documents (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    supplier_id uuid not null references suppliers(id) on delete cascade,
    document_type text not null,
    document_number text not null default '',
    issue_date date,
    expiration_date date,
    status text not null default '',
    attachment_reference text not null default '',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1
);

create table if not exists supplier_ratings (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    supplier_id uuid not null references suppliers(id) on delete cascade,
    quality numeric(5,2) not null default 0,
    price numeric(5,2) not null default 0,
    delivery numeric(5,2) not null default 0,
    service numeric(5,2) not null default 0,
    reliability numeric(5,2) not null default 0,
    overall_score numeric(5,2) not null default 0,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1
);

create table if not exists supplier_contracts (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    supplier_id uuid not null references suppliers(id) on delete cascade,
    contract_number text not null,
    start_date date,
    end_date date,
    status text not null default '',
    notes text not null default '',
    attachment_reference text not null default '',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1
);

create table if not exists supplier_representatives (
    id uuid primary key default uuid_generate_v4(),
    tenant_id uuid not null references tenants(id),
    supplier_id uuid not null references suppliers(id) on delete cascade,
    name text not null,
    document text not null default '',
    email text not null default '',
    phone text not null default '',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    version bigint not null default 1
);

create index if not exists ix_suppliers_tenant_status on suppliers(tenant_id, status) where deleted_at is null;
create index if not exists ix_suppliers_tenant_type on suppliers(tenant_id, type) where deleted_at is null;
create index if not exists ix_suppliers_category on suppliers(tenant_id, category_id) where deleted_at is null;
create index if not exists ix_supplier_contacts_supplier on supplier_contacts(tenant_id, supplier_id) where deleted_at is null;
create index if not exists ix_supplier_addresses_supplier on supplier_addresses(tenant_id, supplier_id) where deleted_at is null;
create index if not exists ix_supplier_documents_supplier on supplier_documents(tenant_id, supplier_id) where deleted_at is null;
create index if not exists ix_supplier_contracts_supplier on supplier_contracts(tenant_id, supplier_id) where deleted_at is null;
create index if not exists ix_supplier_ratings_supplier on supplier_ratings(tenant_id, supplier_id) where deleted_at is null;
