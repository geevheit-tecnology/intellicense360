-- Sprint 026: Checklist operational & fleet alignment.
-- Adds tenant-safe vehicle scoping, vehicle binding to executions, non-conformity tracking, and permissions.

ALTER TABLE checklists ADD CONSTRAINT checklists_tenant_scope_unique UNIQUE(tenant_id, id);
ALTER TABLE checklist_templates ADD CONSTRAINT checklist_templates_tenant_scope_unique UNIQUE(tenant_id, id);
ALTER TABLE checklist_template_versions ADD CONSTRAINT checklist_template_versions_tenant_scope_unique UNIQUE(tenant_id, id);
ALTER TABLE checklist_executions ADD CONSTRAINT checklist_executions_tenant_scope_unique UNIQUE(tenant_id, id);

ALTER TABLE checklist_executions ADD COLUMN IF NOT EXISTS vehicle_id uuid;

ALTER TABLE checklists
    ADD CONSTRAINT fk_checklists_vehicle_tenant
    FOREIGN KEY (tenant_id, vehicle_id)
    REFERENCES fleet_vehicles(tenant_id, id);

ALTER TABLE checklist_executions
    ADD CONSTRAINT fk_checklist_executions_vehicle_tenant
    FOREIGN KEY (tenant_id, vehicle_id)
    REFERENCES fleet_vehicles(tenant_id, id);

CREATE INDEX IF NOT EXISTS ix_checklist_executions_vehicle ON checklist_executions(tenant_id, vehicle_id);
CREATE INDEX IF NOT EXISTS ix_checklist_non_conformities_severity ON checklist_non_conformities(tenant_id, severity, status);

INSERT INTO permissions(key, description) VALUES
    ('checklist.checklists.read', 'Consultar checklists e templates'),
    ('checklist.checklists.manage', 'Gerenciar templates e itens de checklist'),
    ('checklist.checklist.manage', 'Executar e gerenciar checklists operacionais'),
    ('checklist.executions.read', 'Consultar execuções de checklist'),
    ('checklist.executions.execute', 'Realizar preenchimento de checklist de veículo')
ON CONFLICT (key) DO NOTHING;

-- Grant to admin roles
INSERT INTO role_permissions(role_id, permission_id)
SELECT DISTINCT r.id, p.id FROM roles r
JOIN role_permissions rp ON rp.role_id = r.id
JOIN permissions prev ON prev.id = rp.permission_id AND prev.key = 'fleet.vehicles.manage'
CROSS JOIN permissions p
WHERE r.name = 'admin' AND r.deleted_at IS NULL
AND p.key IN (
    'checklist.checklists.read',
    'checklist.checklists.manage',
    'checklist.checklist.manage',
    'checklist.executions.read',
    'checklist.executions.execute'
)
ON CONFLICT DO NOTHING;

