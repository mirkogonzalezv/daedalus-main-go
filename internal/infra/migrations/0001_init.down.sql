-- Eliminar índices
DROP INDEX IF EXISTS daedalus.idx_audit_logs_tenant;
DROP INDEX IF EXISTS daedalus.idx_sessions_refresh_hash;
DROP INDEX IF EXISTS daedalus.idx_sessions_user;
DROP INDEX IF EXISTS daedalus.idx_messages_conversation;
DROP INDEX IF EXISTS daedalus.idx_conversations_user;
DROP INDEX IF EXISTS daedalus.idx_conversations_tenant;
DROP INDEX IF EXISTS daedalus.idx_users_tenant_email;

-- Eliminar tablas (en orden inverso por dependencias)
DROP TABLE IF EXISTS daedalus.sessions CASCADE;
DROP TABLE IF EXISTS daedalus.subscriptions CASCADE;
DROP TABLE IF EXISTS daedalus.audit_logs CASCADE;
DROP TABLE IF EXISTS daedalus.messages CASCADE;
DROP TABLE IF EXISTS daedalus.conversations CASCADE;
DROP TABLE IF EXISTS daedalus.users CASCADE;
DROP TABLE IF EXISTS daedalus.tenants CASCADE;

-- Eliminar extensión
DROP EXTENSION IF EXISTS "pgcrypto" CASCADE;

-- Eliminar schema
DROP SCHEMA IF EXISTS daedalus CASCADE;