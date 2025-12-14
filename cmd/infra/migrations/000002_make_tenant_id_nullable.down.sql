-- Eliminar constraint
ALTER TABLE daedalus.users 
DROP CONSTRAINT IF EXISTS check_tenant_id_for_role;

-- Eliminar usuarios root antes de hacer NOT NULL (si los hay)
DELETE FROM daedalus.users WHERE role = 'root';

-- Volver a hacer tenant_id NOT NULL
ALTER TABLE daedalus.users 
ALTER COLUMN tenant_id SET NOT NULL;