create table memory (
	id            bigserial primary key,
	title         text not null,
	details       text not null,
	latitude      numeric(14,11) not null,
	longitude     numeric(14,11) not null,
	author        text not null default 'Anonymous',
	is_approved   boolean default false,
	approval_uuid uuid,
	inserted_at   timestamp with time zone not null default now(),
	updated_at    timestamp with time zone not null
);

commit;