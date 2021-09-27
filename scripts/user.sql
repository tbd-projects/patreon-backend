\connect restapi_dev
CREATE TABLE IF NOT EXISTS users
(
    user_id            bigserial not null primary key,
    login              text      not null unique,
    nickname           text      not null unique,
    encrypted_password text      not null,
    avatar             text      not null
);
CREATE TABLE IF NOT EXISTS creator_profile
(
    creator_id  integer not null primary key,
    avatar      text    not null,
    cover       text    not null,
    description text    not null,
    category    text    not null,
    foreign key (creator_id) references users (user_id) on delete cascade
);
CREATE TABLE IF NOT EXISTS subscribers
(
    id         bigserial not null primary key,
    user_id    integer   not null,
    creator_id integer   not null,
    foreign key (creator_id) references creator_profile (creator_id),
    foreign key (user_id) references users (user_id)
);
\disconnect