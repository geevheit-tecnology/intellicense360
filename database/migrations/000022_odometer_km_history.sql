-- Sprint 022: Durable Odometer & Vehicle KM History. Append-only, anti-regression, auditable.

CREATE TABLE fleet_odometer_readings (
    id uuid PRIMARY KEY,
    tenant_id uuid NOT NULL REFERENCES tenants(id),
    vehicle_id uuid NOT NULL,
    reading bigint NOT NULL CHECK (reading >= 0),
    read_at timestamptz NOT NULL,
    source text NOT NULL CHECK (source IN ('manual', 'fueling', 'maintenance', 'checklist', 'telemetry', 'correction')),
    driver_id uuid,
    actor_id uuid NOT NULL,
    notes text NOT NULL DEFAULT '',
    is_correction boolean NOT NULL DEFAULT false,
    correction_reason text NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL,
    FOREIGN KEY (tenant_id, vehicle_id) REFERENCES fleet_vehicles(tenant_id, id),
    FOREIGN KEY (tenant_id, driver_id) REFERENCES drivers(tenant_id, id),
    FOREIGN KEY (tenant_id, actor_id) REFERENCES users(tenant_id, id),
    CHECK (is_correction = false OR length(btrim(correction_reason)) >= 10)
);

CREATE INDEX fleet_odometer_vehicle_idx ON fleet_odometer_readings(tenant_id, vehicle_id, read_at DESC, created_at DESC);
CREATE INDEX fleet_odometer_driver_idx ON fleet_odometer_readings(tenant_id, driver_id) WHERE driver_id IS NOT NULL;
CREATE INDEX fleet_odometer_actor_idx ON fleet_odometer_readings(tenant_id, actor_id);

INSERT INTO permissions(key, description) VALUES
    ('fleet.odometer.read', 'Consultar leituras de hodômetro'),
    ('fleet.odometer.record', 'Registrar leitura de hodômetro'),
    ('fleet.odometer.correct', 'Corrigir ou redefinir hodômetro com justificativa')
ON CONFLICT (key) DO NOTHING;

-- Upgrade only company-admin roles that already manage fleet. No new global powers.
INSERT INTO role_permissions(role_id, permission_id)
SELECT DISTINCT r.id, p.id FROM roles r
JOIN role_permissions rp ON rp.role_id = r.id
JOIN permissions prev ON prev.id = rp.permission_id AND prev.key = 'fleet.vehicles.manage'
CROSS JOIN permissions p
WHERE r.name = 'admin' AND r.deleted_at IS NULL
AND p.key IN ('fleet.odometer.read', 'fleet.odometer.record', 'fleet.odometer.correct')
ON CONFLICT DO NOTHING;

