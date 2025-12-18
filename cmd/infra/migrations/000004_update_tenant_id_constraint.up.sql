-- Revertir al constraint original
ALTER TABLE daedalus.users DROP CONSTRAINT check_tenant_id_for_role;

ALTER TABLE daedalus.users
ADD CONSTRAINT check_tenant_id_for_role
CHECK (
    (role = 'root' AND tenant_id IS NULL) OR
    (role != 'root' AND tenant_id IS NOT NULL)
);
