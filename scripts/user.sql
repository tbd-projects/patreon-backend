CREATE TABLE IF NOT EXISTS Users (
    user_id bigserial not null primary key,
    login text not null unique,
    password text not null,
    avatar text not null
);