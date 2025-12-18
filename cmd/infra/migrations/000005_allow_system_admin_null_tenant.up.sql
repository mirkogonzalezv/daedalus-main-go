-- Permitir tenant_id NULL para usuarios root y system_admin
ALTER TABLE daedalus.users DROP CONSTRAINT check_tenant_id_for_role;

ALTER TABLE daedalus.users
ADD CONSTRAINT check_tenant_id_for_role
CHECK (
    (role IN ('root', 'system_admin') AND tenant_id IS NULL) OR
    (role NOT IN ('root', 'system_admin') AND tenant_id IS NOT NULL)
);