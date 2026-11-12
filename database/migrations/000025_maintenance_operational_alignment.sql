-- Sprint 025: Maintenance operational & financial alignment.
-- Adds invoice/cost reconciliation fields, third-party vehicle support, tenant-safe vehicle scoping, and permissions.

ALTER TABLE maintenance_service_types ADD CONSTRAINT maintenance_service_types_tenant_scope_unique UNIQUE(tenant_id, id);
ALTER TABLE maintenance_work_orders ADD CONSTRAINT maintenance_work_orders_tenant_scope_unique UNIQUE(tenant_id, id);
ALTER TABLE maintenance_preventive_plans ADD CONSTRAINT maintenance_preventive_plans_tenant_scope_unique UNIQUE(tenant_id, id);
ALTER TABLE maintenance_workshops ADD CONSTRAINT maintenance_workshops_tenant_scope_unique UNIQUE(tenant_id, id);
ALTER TABLE maintenance_technicians ADD CONSTRAINT maintenance_technicians_tenant_scope_unique UNIQUE(tenant_id, id);

ALTER TABLE maintenance_work_orders
    ADD COLUMN invoice_number text NOT NULL DEFAULT '',
    ADD COLUMN invoice_key text NOT NULL DEFAULT '',
    ADD COLUMN total_cost numeric(14,2) NOT NULL DEFAULT 0,
    ADD COLUMN payment_status text NOT NULL DEFAULT 'pending' CHECK (payment_status IN ('pending', 'approved', 'paid', 'canceled')),
    ADD COLUMN is_third_party boolean NOT NULL DEFAULT false,
    ADD COLUMN third_party_identifier text NOT NULL DEFAULT '';

CREATE INDEX ix_maintenance_work_orders_invoice ON maintenance_work_orders(tenant_id, invoice_number) WHERE invoice_number <> '';
CREATE INDEX ix_maintenance_work_orders_payment ON maintenance_work_orders(tenant_id, payment_status);
CREATE INDEX ix_maintenance_work_orders_third_party ON maintenance_work_orders(tenant_id, is_third_party);

INSERT INTO permissions(key, description) VALUES
    ('maintenance.maintenance.read', 'Consultar ordens de serviço e planos preventivos'),
    ('maintenance.maintenance.manage', 'Gerenciar ordens de serviço, oficinas e planos')
ON CONFLICT (key) DO NOTHING;

-- Grant permissions to existing admin roles that manage fleet
INSERT INTO role_permissions(role_id, permission_id)
SELECT DISTINCT r.id, p.id FROM roles r
JOIN role_permissions rp ON rp.role_id = r.id
JOIN permissions prev ON prev.id = rp.permission_id AND prev.key = 'fleet.vehicles.manage'
CROSS JOIN permissions p
WHERE r.name = 'admin' AND r.deleted_at IS NULL
AND p.key IN ('maintenance.maintenance.read', 'maintenance.maintenance.manage')
ON CONFLICT DO NOTHING;

