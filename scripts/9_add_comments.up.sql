DROP INDEX idx_like_search_nick_creators;

CREATE INDEX idx_like_search_nick_creators
    on users (LOWER(nickname) text_pattern_ops);

ALTER TABLE comments
    ADD COLUMN as_creator bool        default false              not null,
    ADD COLUMN date       timestamptz default now()::timestamptz not null;

ALTER TABLE posts
    ADD COLUMN number_comments bigint default 0 not null;