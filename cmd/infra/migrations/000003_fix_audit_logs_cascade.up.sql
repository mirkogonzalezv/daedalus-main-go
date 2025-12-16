-- Eliminar las foreign keys existentes con cascade
ALTER TABLE daedalus.audit_logs DROP CONSTRAINT audit_logs_tenant_id_fkey;
ALTER TABLE daedalus.audit_logs DROP CONSTRAINT audit_logs_user_id_fkey;

-- Recrear las foreign keys sin cascade (SET NULL para preservar los logs)
ALTER TABLE daedalus.audit_logs 
ADD CONSTRAINT audit_logs_tenant_id_fkey 
FOREIGN KEY (tenant_id) REFERENCES daedalus.tenants(id) ON DELETE SET NULL;

ALTER TABLE daedalus.audit_logs 
ADD CONSTRAINT audit_logs_user_id_fkey 
FOREIGN KEY (user_id) REFERENCES daedalus.users(id) ON DELETE SET NULL;

-- Cambiar tenant_id para que pueda ser NULL
ALTER TABLE daedalus.audit_logs ALTER COLUMN tenant_id DROP NOT NULL;
