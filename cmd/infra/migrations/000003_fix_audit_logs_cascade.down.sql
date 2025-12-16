-- Revertir los cambios: eliminar las foreign keys actuales
ALTER TABLE daedalus.audit_logs DROP CONSTRAINT audit_logs_tenant_id_fkey;
ALTER TABLE daedalus.audit_logs DROP CONSTRAINT audit_logs_user_id_fkey;

-- Hacer tenant_id NOT NULL nuevamente (solo si no hay registros con NULL)
ALTER TABLE daedalus.audit_logs ALTER COLUMN tenant_id SET NOT NULL;

-- Recrear las foreign keys con CASCADE como estaban originalmente
ALTER TABLE daedalus.audit_logs 
ADD CONSTRAINT audit_logs_tenant_id_fkey 
FOREIGN KEY (tenant_id) REFERENCES daedalus.tenants(id) ON DELETE CASCADE;

ALTER TABLE daedalus.audit_logs 
ADD CONSTRAINT audit_logs_user_id_fkey 
FOREIGN KEY (user_id) REFERENCES daedalus.users(id) ON DELETE SET NULL;