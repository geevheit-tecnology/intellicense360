-- Migration 000029: Inventory & Stock Operational Alignment
-- Enables durable stock tracking for parts, tires, oils, onboard equipment and maintenance supplies.

alter table if exists inventory_parts
    add column if not exists unit_cost numeric(14,2) not null default 0,
    add column if not exists min_stock numeric(14,3) not null default 0,
    add column if not exists current_stock numeric(14,3) not null default 0,
    add column if not exists location text not null default '';

create index if not exists ix_inventory_parts_tenant_sku on inventory_parts(tenant_id, sku);
create index if not exists ix_inventory_parts_tenant_name on inventory_parts(tenant_id, name);

-- Permissions
insert into permissions(id, key, description) values
    (gen_random_uuid(), 'inventory.inventory.read', 'Consultar estoque, peças e equipamentos'),
    (gen_random_uuid(), 'inventory.inventory.manage', 'Gerenciar cadastro de peças, entradas e saídas de estoque')
on conflict (key) do nothing;

insert into roles_permissions(tenant_id, role_id, permission_id)
select r.tenant_id, r.id, p.id
from roles r
cross join permissions p
where r.name = 'admin' and p.key in (
    'inventory.inventory.read',
    'inventory.inventory.manage'
)
on conflict do nothing;

