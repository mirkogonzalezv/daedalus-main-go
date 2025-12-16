create extension if not exists "pgcrypto" SCHEMA daedalus;

create table daedalus.tenants(
	id UUID primary key default gen_random_uuid(),
	name text not null,
	slug text unique not null,
	subscription_plan text,
	status text not null default 'active',
	created_at timestamp with time zone default now(),
	updated_at timestamp with time zone default now()
);

create table daedalus.users(
	id UUID primary key default gen_random_uuid(),
	tenant_id uuid not null references daedalus.tenants(id) on delete cascade,
	name text,
	email text not null,
	password_hash text not null,
	role text not null default 'user',
	status text not null default 'active',
	created_at timestamp with time zone default now(),
	updated_at timestamp with time zone default now(),
	unique (tenant_id, email)
);


create table daedalus.conversations(
	id UUID primary key default gen_random_uuid(),
	tenant_id UUID not null references daedalus.tenants(id) on delete cascade,
	user_id UUID references daedalus.users(id) on delete set null,
	title text,
	created_at timestamp with time zone default now(),
	updated_at timestamp with time zone default now()
);

create table daedalus.messages(
	id UUID primary key default gen_random_uuid(),
	conversation_id uuid not null references daedalus.conversations(id) on delete cascade,
	user_id uuid references daedalus.users(id) on delete set null,
	sender text not null,
	content text not null,
	created_at timestamp with time zone default now()
);

create table daedalus.audit_logs(
	id uuid primary key default gen_random_uuid(),
	tenant_id uuid not null references daedalus.tenants(id) on delete cascade,
	user_id uuid references daedalus.users(id) on delete set null,
	action text not null,
	meta JSONB,
	ip text,
	created_at timestamp with time zone default now()
);

create table daedalus.subscriptions(
	id uuid primary key default gen_random_uuid(),
	tenant_id uuid not null references daedalus.tenants(id) on delete cascade,
	stripe_customer_id text,
	stripe_subscription_id text,
	ºperiod_start timestamp with time zone,
	period_end timestamp with time zone,
	status text,
	created_at timestamp with time zone default now(),
	updated_at timestamp with time zone default now()
);

create table daedalus.sessions(
	id uuid primary key default gen_random_uuid(),
	tenant_id uuid not null references daedalus.tenants(id) on delete cascade,
	user_id uuid not null references daedalus.users(id) on delete cascade,
	refresh_hash TEXT not null,
	issued_at timestamp with time zone default now(),
	expires_at timestamp with time zone,
	ip text,
	user_agent text,
	revoked boolean default false
);

-- indexes
create index idx_users_tenant_email on daedalus.users(tenant_id, email);
create index idx_conversations_tenant on daedalus.conversations(tenant_id);
create index idx_conversations_user on daedalus.conversations(user_id);
create index idx_messages_conversation on daedalus.messages(conversation_id);
create index idx_sessions_user on daedalus.sessions(user_id);
create index idx_sessions_refresh_hash on daedalus.sessions(refresh_hash);
create index idx_audit_logs_tenant on daedalus.audit_logs(tenant_id);
