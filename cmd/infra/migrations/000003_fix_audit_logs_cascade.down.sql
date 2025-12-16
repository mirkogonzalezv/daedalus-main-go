-- Revertir los cambios
ALTER TABLE daedalus.audit_logs DROP CONSTRAINT audit_logs_tenant_id_fkey;
ALTER TABLE daedalus.audit_logs DROP CONSTRAINT audit_logs_user_id_fkey;

-- Hacer tenant_id NOT NULL nuevamente
ALTER TABLE daedalus.audit_logs ALTER COLUMN tenant_id SET NOT NULL;

-- Recrear las foreign keys con CASCADE
ALTER TABLE daedalus.audit_logs 
ADD CONSTRAINT audit_logs_tenant_id_fkey 
FOREIGN KEY (tenant_id) REFERENCES daedalus.tenants(id) ON DELETE CASCADE;

ALTER TABLE daedalus.audit_logs 
ADD CONSTRAINT audit_logs_user_id_fkey 
FOREIGN KEY (user_id) REFERENCES daedalus.users(id) ON DELETE SET NULL;
