-- Development seed. Do not use these credentials in production.

insert into tenants (id, slug, name, active)
values ('00000000-0000-0000-0000-000000000001', 'bootstrap-tenant', 'Bootstrap Tenant', true)
on conflict (id) do nothing;

insert into permissions (key, description)
values
    ('assets.assets.manage', 'Manage assets'),
    ('identity.users.manage', 'Manage users'),
    ('identity.roles.manage', 'Manage roles'),
    ('identity.permissions.manage', 'Manage permissions'),
    ('identity.tenants.manage', 'Manage tenants'),
    ('checklist.checklists.manage', 'Manage checklists'),
    ('checklist.checklist.manage', 'Manage checklist engine'),
    ('fleet.vehicles.manage', 'Manage fleet vehicles'),
    ('fleet.assets.manage', 'Manage fleet assets'),
    ('financial.financial.manage', 'Manage financial operations'),
    ('ciot.ciot.manage', 'Manage CIOT core'),
    ('fuel.fuel.manage', 'Manage fuel'),
    ('inventory.inventory.manage', 'Manage inventory'),
    ('maintenance.maintenance.manage', 'Manage maintenance'),
    ('suppliers.suppliers.manage', 'Manage suppliers'),
    ('tires.tires.manage', 'Manage tires')
on conflict (key) do nothing;

insert into roles (id, tenant_id, name, description)
values ('00000000-0000-0000-0000-000000000101', '00000000-0000-0000-0000-000000000001', 'admin', 'Development administrator')
on conflict (tenant_id, name) do nothing;

-- Password here is illustrative for dev seed; production adapters must use the application hasher.
insert into users (id, tenant_id, name, email, password_hash, status)
values ('00000000-0000-0000-0000-000000000201', '00000000-0000-0000-0000-000000000001', 'Admin', 'admin@geevheit.local', 'seed:35bd6520ed10376664a480ae1bdcf05d915b7fb853c67c77ae43efb18c2c5a98', 'active')
on conflict (tenant_id, email) do nothing;

insert into tenant_users (tenant_id, user_id, role_id)
values ('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000201', '00000000-0000-0000-0000-000000000101')
on conflict (tenant_id, user_id) do nothing;

insert into role_permissions (role_id, permission_id)
select '00000000-0000-0000-0000-000000000101', p.id
from permissions p
where p.key in (
    'assets.assets.manage',
    'identity.users.manage',
    'identity.roles.manage',
    'identity.permissions.manage',
    'identity.tenants.manage',
    'checklist.checklists.manage',
    'checklist.checklist.manage',
    'fleet.vehicles.manage',
    'fleet.assets.manage',
    'financial.financial.manage',
    'ciot.ciot.manage',
    'fuel.fuel.manage',
    'inventory.inventory.manage',
    'maintenance.maintenance.manage',
    'suppliers.suppliers.manage',
    'tires.tires.manage'
)
on conflict (role_id, permission_id) do nothing;
