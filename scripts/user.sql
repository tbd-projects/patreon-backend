\connect restapi_dev
CREATE TABLE IF NOT EXISTS users (
    user_id bigserial not null primary key,
    login text not null unique,
    encrypted_password text not null,
    avatar text not null
);
\disconnect