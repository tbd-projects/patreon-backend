CREATE TABLE push_history (
    id bigserial not null primary key,
    users_id bigint not null  references users(users_id),
    push_type text not null,
    push jsonb not null,
    date timestamptz not null default now()::timestamptz,
    is_viewed bool not null default false
)



