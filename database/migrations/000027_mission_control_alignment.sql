-- Sprint 027: Mission Control operational alignment & permissions.

ALTER TABLE mission_control_command_items ADD CONSTRAINT mission_control_command_items_tenant_scope_unique UNIQUE(tenant_id, id);

INSERT INTO permissions(key, description) VALUES
    ('mission_control.read', 'Consultar painel mission control e alertas'),
    ('mission_control.create', 'Criar itens de comando e alertas'),
    ('mission_control.update', 'Atualizar itens do mission control'),
    ('mission_control.acknowledge', 'Reconhecer alertas e sinais críticos'),
    ('mission_control.resolve', 'Resolver alertas e itens de comando'),
    ('mission_control.dismiss', 'Dispensar alertas não críticos'),
    ('mission_control.snapshot', 'Gerar snapshot do status operacional'),
    ('mission_control.admin', 'Administração geral do mission control')
ON CONFLICT (key) DO NOTHING;

INSERT INTO role_permissions(role_id, permission_id)
SELECT DISTINCT r.id, p.id FROM roles r
JOIN role_permissions rp ON rp.role_id = r.id
JOIN permissions prev ON prev.id = rp.permission_id AND prev.key = 'fleet.vehicles.manage'
CROSS JOIN permissions p
WHERE r.name = 'admin' AND r.deleted_at IS NULL
AND p.key LIKE 'mission_control.%'
ON CONFLICT DO NOTHING;

