-- Durable operational profiles and document metadata. No uploads or external calls.
ALTER TABLE fleet_vehicles ADD CONSTRAINT fleet_vehicle_scope_unique UNIQUE(tenant_id,id);

CREATE TABLE drivers (
 id uuid PRIMARY KEY,
 tenant_id uuid NOT NULL REFERENCES tenants(id),
 code text NOT NULL CHECK(length(btrim(code)) BETWEEN 1 AND 64),
 name text NOT NULL CHECK(length(btrim(name)) BETWEEN 1 AND 200),
 tax_id text CHECK(tax_id IS NULL OR tax_id ~ '^[0-9]{11}$'),
 phone text NOT NULL DEFAULT '',
 email text NOT NULL DEFAULT '',
 vehicle_id uuid,
 status text NOT NULL DEFAULT 'active' CHECK(status IN ('active','inactive')),
 notes text NOT NULL DEFAULT '',
 created_at timestamptz NOT NULL,
 updated_at timestamptz NOT NULL,
 deleted_at timestamptz,
 version bigint NOT NULL DEFAULT 1 CHECK(version>0),
 UNIQUE(tenant_id,id),
 FOREIGN KEY(tenant_id,vehicle_id) REFERENCES fleet_vehicles(tenant_id,id)
);
CREATE UNIQUE INDEX drivers_code_active_unique ON drivers(tenant_id,upper(code)) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX drivers_tax_id_active_unique ON drivers(tenant_id,tax_id) WHERE deleted_at IS NULL AND tax_id IS NOT NULL;
CREATE INDEX drivers_status_idx ON drivers(tenant_id,status,created_at,id) WHERE deleted_at IS NULL;
CREATE INDEX drivers_vehicle_idx ON drivers(tenant_id,vehicle_id) WHERE deleted_at IS NULL;

CREATE TABLE documents (
 id uuid PRIMARY KEY,
 tenant_id uuid NOT NULL REFERENCES tenants(id),
 driver_id uuid,
 vehicle_id uuid,
 kind text NOT NULL CHECK(kind IN ('cnh','crlv','insurance','other')),
 title text NOT NULL CHECK(length(btrim(title)) BETWEEN 1 AND 200),
 number text NOT NULL CHECK(length(btrim(number)) BETWEEN 1 AND 100),
 issuer text NOT NULL DEFAULT '',
 issued_on date,
 expires_on date,
 status text NOT NULL DEFAULT 'active' CHECK(status IN ('active','archived')),
 notes text NOT NULL DEFAULT '',
 created_at timestamptz NOT NULL,
 updated_at timestamptz NOT NULL,
 deleted_at timestamptz,
 version bigint NOT NULL DEFAULT 1 CHECK(version>0),
 CHECK((driver_id IS NOT NULL)::int+(vehicle_id IS NOT NULL)::int=1),
 CHECK(kind<>'cnh' OR driver_id IS NOT NULL),
 CHECK(kind<>'crlv' OR vehicle_id IS NOT NULL),
 CHECK(issued_on IS NULL OR issued_on>='1900-01-01'),
 CHECK(expires_on IS NULL OR expires_on>='1900-01-01'),
 CHECK(issued_on IS NULL OR expires_on IS NULL OR expires_on>=issued_on),
 FOREIGN KEY(tenant_id,driver_id) REFERENCES drivers(tenant_id,id),
 FOREIGN KEY(tenant_id,vehicle_id) REFERENCES fleet_vehicles(tenant_id,id)
);
CREATE UNIQUE INDEX documents_driver_number_active_unique ON documents(tenant_id,driver_id,kind,upper(number)) WHERE deleted_at IS NULL AND status='active' AND driver_id IS NOT NULL;
CREATE UNIQUE INDEX documents_vehicle_number_active_unique ON documents(tenant_id,vehicle_id,kind,upper(number)) WHERE deleted_at IS NULL AND status='active' AND vehicle_id IS NOT NULL;
CREATE INDEX documents_expiry_idx ON documents(tenant_id,expires_on,id) WHERE deleted_at IS NULL AND status='active';
CREATE INDEX documents_driver_idx ON documents(tenant_id,driver_id) WHERE deleted_at IS NULL;
CREATE INDEX documents_vehicle_idx ON documents(tenant_id,vehicle_id) WHERE deleted_at IS NULL;

INSERT INTO permissions(key,description) VALUES
 ('drivers.read','Consultar motoristas'),
 ('drivers.manage','Gerenciar motoristas'),
 ('documents.read','Consultar documentos'),
 ('documents.manage','Gerenciar documentos')
ON CONFLICT(key) DO NOTHING;

-- Upgrade only company-admin roles that already manage fleet. No new global powers.
INSERT INTO role_permissions(role_id,permission_id)
SELECT DISTINCT r.id,p.id FROM roles r
JOIN role_permissions rp ON rp.role_id=r.id
JOIN permissions previous ON previous.id=rp.permission_id AND previous.key='fleet.vehicles.manage'
CROSS JOIN permissions p
WHERE r.name='admin' AND r.deleted_at IS NULL
AND p.key IN ('drivers.read','drivers.manage','documents.read','documents.manage')
ON CONFLICT DO NOTHING;
