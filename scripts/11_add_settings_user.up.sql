CREATE TABLE user_settings (
    user_id bigint not null references users (users_id) on delete cascade,
    get_post bool not null,
    get_comment bool not null,
    get_sub bool not null
);
