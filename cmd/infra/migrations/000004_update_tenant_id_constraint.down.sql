-- Eliminar constraint existente
ALTER TABLE daedalus.users DROP CONSTRAINT check_tenant_id_for_role;

-- Crear nuevo constraint que incluya system_admin
ALTER TABLE daedalus.users
ADD CONSTRAINT check_tenant_id_for_role
CHECK (
    (role IN ('root', 'system_admin') AND tenant_id IS NULL) OR
    (role NOT IN ('root', 'system_admin') AND tenant_id IS NOT NULL)
);
