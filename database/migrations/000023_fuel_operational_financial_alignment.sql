-- Sprint 023: Fuel operational & financial alignment.
-- Adds support for fueling requisitions, third-party vehicles, invoice reconciliation, and tenant-safe scoping.

ALTER TABLE fuel_stations ADD CONSTRAINT fuel_station_tenant_scope_unique UNIQUE(tenant_id, id);
ALTER TABLE fuel_types ADD CONSTRAINT fuel_type_tenant_scope_unique UNIQUE(tenant_id, id);
ALTER TABLE fuel_transactions ADD CONSTRAINT fuel_transaction_tenant_scope_unique UNIQUE(tenant_id, id);

ALTER TABLE fuel_transactions
    ADD COLUMN requisition_number text NOT NULL DEFAULT '',
    ADD COLUMN is_third_party boolean NOT NULL DEFAULT false,
    ADD COLUMN third_party_identifier text NOT NULL DEFAULT '',
    ADD COLUMN invoice_number text NOT NULL DEFAULT '',
    ADD COLUMN invoice_key text NOT NULL DEFAULT '',
    ADD COLUMN payment_status text NOT NULL DEFAULT 'pending' CHECK (payment_status IN ('pending', 'approved', 'paid', 'canceled'));

CREATE INDEX fuel_transactions_requisition_idx ON fuel_transactions(tenant_id, requisition_number) WHERE requisition_number <> '';
CREATE INDEX fuel_transactions_invoice_idx ON fuel_transactions(tenant_id, invoice_number) WHERE invoice_number <> '';
CREATE INDEX fuel_transactions_payment_status_idx ON fuel_transactions(tenant_id, payment_status);
CREATE INDEX fuel_transactions_third_party_idx ON fuel_transactions(tenant_id, is_third_party);

INSERT INTO permissions(key, description) VALUES
    ('fuel.fuel.read', 'Consultar abastecimentos e postos'),
    ('fuel.fuel.manage', 'Gerenciar abastecimentos, requisições e postos')
ON CONFLICT (key) DO NOTHING;

-- Grant permissions to existing admin roles that manage fleet
INSERT INTO role_permissions(role_id, permission_id)
SELECT DISTINCT r.id, p.id FROM roles r
JOIN role_permissions rp ON rp.role_id = r.id
JOIN permissions prev ON prev.id = rp.permission_id AND prev.key = 'fleet.vehicles.manage'
CROSS JOIN permissions p
WHERE r.name = 'admin' AND r.deleted_at IS NULL
AND p.key IN ('fuel.fuel.read', 'fuel.fuel.manage')
ON CONFLICT DO NOTHING;

