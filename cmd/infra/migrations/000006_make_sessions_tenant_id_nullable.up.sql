-- Hacer tenant_id nullable en sessions para usuarios globales
ALTER TABLE daedalus.sessions
ALTER COLUMN tenant_id DROP NOT NULL;

-- Índice para optimizar queries de sesiones globales
CREATE INDEX idx_sessions_global_users
ON daedalus.sessions(user_id, expires_at)
WHERE tenant_id IS NULL;
