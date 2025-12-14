-- Hacer tenant_id nullable para usuarios root
ALTER TABLE daedalus.users 
ALTER COLUMN tenant_id DROP NOT NULL;

-- Agregar constraint para validar que solo usuarios root pueden tener tenant_id NULL
ALTER TABLE daedalus.users 
ADD CONSTRAINT check_tenant_id_for_role 
CHECK (
    (role = 'root' AND tenant_id IS NULL) OR 
    (role != 'root' AND tenant_id IS NOT NULL)
);