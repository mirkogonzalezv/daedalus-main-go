-- ADVERTENCIA: Elimina sesiones globales antes de revertir
DELETE FROM daedalus.sessions
WHERE tenant_id IS NULL;

-- Eliminar índice
DROP INDEX IF EXISTS idx_sessions_global_users;

-- Restaurar NOT NULL constraint
ALTER TABLE daedalus.sessions
ALTER COLUMN tenant_id SET NOT NULL;
