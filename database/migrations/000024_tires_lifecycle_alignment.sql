-- Sprint 024: Tires Lifecycle & Operational Alignment.
-- Enforces tenant-safe vehicle scoping, unique constraints, and operational permissions.

ALTER TABLE tires ADD CONSTRAINT tires_tenant_scope_unique UNIQUE(tenant_id, id);
ALTER TABLE tires ADD CONSTRAINT tires_vehicle_scope FOREIGN KEY(tenant_id, vehicle_id) REFERENCES fleet_vehicles(tenant_id, id);
ALTER TABLE tire_movements ADD CONSTRAINT tire_movements_tenant_scope FOREIGN KEY(tenant_id, vehicle_id) REFERENCES fleet_vehicles(tenant_id, id);

INSERT INTO permissions(key, description) VALUES
    ('tires.tires.read', 'Consultar pneus, inspeções e movimentações'),
    ('tires.tires.manage', 'Gerenciar pneus, rodízios, recapes e descarte')
ON CONFLICT (key) DO NOTHING;

-- Grant permissions to existing admin roles that manage fleet
INSERT INTO role_permissions(role_id, permission_id)
SELECT DISTINCT r.id, p.id FROM roles r
JOIN role_permissions rp ON rp.role_id = r.id
JOIN permissions prev ON prev.id = rp.permission_id AND prev.key = 'fleet.vehicles.manage'
CROSS JOIN permissions p
WHERE r.name = 'admin' AND r.deleted_at IS NULL
AND p.key IN ('tires.tires.read', 'tires.tires.manage')
ON CONFLICT DO NOTHING;

