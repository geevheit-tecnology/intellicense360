-- Migration 000031: Driver Journey Control & Work Shifts (Lei 13.103/2015 Compliance)

CREATE TABLE IF NOT EXISTS driver_journey_records (
    id uuid PRIMARY KEY,
    tenant_id uuid NOT NULL REFERENCES tenants(id),
    driver_id uuid NOT NULL,
    vehicle_id uuid,
    trip_id uuid,
    event_type text NOT NULL CHECK (event_type IN ('driving', 'break_rest', 'meal', 'daily_rest', 'waiting', 'standby')),
    start_at timestamptz NOT NULL,
    end_at timestamptz,
    duration_minutes integer NOT NULL DEFAULT 0,
    location text NOT NULL DEFAULT '',
    odometer_km bigint NOT NULL DEFAULT 0,
    is_violation boolean NOT NULL DEFAULT false,
    violation_reason text NOT NULL DEFAULT '',
    daily_allowance_amount numeric(10,2) NOT NULL DEFAULT 0,
    notes text NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    version bigint NOT NULL DEFAULT 1,
    FOREIGN KEY (tenant_id, driver_id) REFERENCES drivers(tenant_id, id)
);

CREATE INDEX IF NOT EXISTS driver_journey_driver_idx ON driver_journey_records(tenant_id, driver_id, start_at DESC);
CREATE INDEX IF NOT EXISTS driver_journey_violations_idx ON driver_journey_records(tenant_id, is_violation) WHERE is_violation = true;

