-- Migration 000030: Trips, Cargo Dispatch & Autonomous Operational Control

CREATE TABLE IF NOT EXISTS operational_trips (
    id uuid PRIMARY KEY,
    tenant_id uuid NOT NULL REFERENCES tenants(id),
    trip_number text NOT NULL,
    status text NOT NULL CHECK (status IN ('draft', 'dispatched', 'delivered', 'completed', 'canceled')),
    vehicle_id uuid NOT NULL,
    trailer_id uuid,
    driver_id uuid NOT NULL,
    origin_city text NOT NULL,
    origin_state text NOT NULL,
    destination_city text NOT NULL,
    destination_state text NOT NULL,
    distance_km integer NOT NULL DEFAULT 0,
    departure_at timestamptz,
    expected_arrival_at timestamptz,
    delivered_at timestamptz,
    cargo_description text NOT NULL DEFAULT '',
    cargo_weight_kg numeric(12,2) NOT NULL DEFAULT 0,
    cargo_value numeric(14,2) NOT NULL DEFAULT 0,
    customer_name text NOT NULL DEFAULT '',
    mdfe_key text NOT NULL DEFAULT '',
    cte_key text NOT NULL DEFAULT '',
    ciot_number text NOT NULL DEFAULT '',
    freight_amount numeric(12,2) NOT NULL DEFAULT 0,
    advance_amount numeric(12,2) NOT NULL DEFAULT 0,
    toll_amount numeric(12,2) NOT NULL DEFAULT 0,
    balance_amount numeric(12,2) NOT NULL DEFAULT 0,
    start_odometer_km bigint NOT NULL DEFAULT 0,
    end_odometer_km bigint NOT NULL DEFAULT 0,
    validation_notes text NOT NULL DEFAULT '',
    notes text NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    deleted_at timestamptz,
    version bigint NOT NULL DEFAULT 1,
    FOREIGN KEY (tenant_id, vehicle_id) REFERENCES fleet_vehicles(tenant_id, id),
    FOREIGN KEY (tenant_id, driver_id) REFERENCES drivers(tenant_id, id)
);

CREATE INDEX IF NOT EXISTS operational_trips_tenant_status_idx ON operational_trips(tenant_id, status, departure_at DESC);
CREATE INDEX IF NOT EXISTS operational_trips_vehicle_idx ON operational_trips(tenant_id, vehicle_id);
CREATE INDEX IF NOT EXISTS operational_trips_driver_idx ON operational_trips(tenant_id, driver_id);

INSERT INTO permissions(key, description) VALUES
    ('trips.read', 'Consultar viagens e despachos de cargas'),
    ('trips.manage', 'Criar, despachar e gerenciar viagens de frota')
ON CONFLICT (key) DO NOTHING;

-- Grant to admin roles
INSERT INTO role_permissions(role_id, permission_id)
SELECT DISTINCT r.id, p.id FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'admin' AND r.deleted_at IS NULL
AND p.key IN ('trips.read', 'trips.manage')
ON CONFLICT DO NOTHING;

