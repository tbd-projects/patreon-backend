\connect restapi_dev
-- create type posts_type as enum ('podcasts','article', 'music', 'designer', 'illustrator', 'journalist','mass media');
CREATE TABLE IF NOT EXISTS users
(
    users_id           bigserial not null primary key,
    login              text      not null unique,
    nickname           text      not null unique,
    encrypted_password text      not null,
    avatar             text
);

CREATE TABLE IF NOT EXISTS creator_category
(
    category_id bigserial not null primary key,
    name        text      not null
);

CREATE TABLE IF NOT EXISTS creator_profile
(
    creator_id  bigint not null primary key,
    category    bigint not null references creator_category (category_id) on delete cascade,
    description text   not null,
    avatar      text,
    cover       text,
    foreign key (creator_id) references users (users_id) on delete cascade
);

CREATE TABLE IF NOT EXISTS subscribers
(
    id         bigserial not null primary key,
    users_id   bigint    not null,
    creator_id bigint    not null,
    foreign key (creator_id) references creator_profile (creator_id),
    foreign key (users_id) references users (users_id)
);


CREATE TABLE IF NOT EXISTS posts_type
(
    posts_type_id bigserial not null primary key,
    type          text      not null
);


CREATE TABLE IF NOT EXISTS awards
(
    awards_id   bigserial        not null primary key,
    name        text             not null,
    description text             not null,
    price       integer          not null,
    color       bigint default 0 not null,
    creator_id  bigint references creator_profile (creator_id) on delete cascade,
    UNIQUE (name, creator_id)
);

CREATE TABLE IF NOT EXISTS posts
(
    posts_id    bigserial               not null primary key,
    title       text                    not null,
    description text                    not null,
    likes       bigint    default 0     not null,
    date        timestamp default now() not null,
    cover       text,
    type_awards bigint references awards (awards_id),
    creator_id  bigint references creator_profile (creator_id) on delete cascade
);

CREATE TABLE IF NOT EXISTS posts_data
(
    data_id   bigserial not null primary key,
    post_id   bigint references posts (posts_id) on delete cascade,
    type      bigint    not null references posts_type (posts_type_id) on delete cascade,
    data_path text      not null
);

CREATE TABLE IF NOT EXISTS likes
(
    likes_id bigserial not null primary key,
    value    bool      not null,
    post_id  bigint    not null references posts (posts_id) on delete cascade,
    users_id bigint    not null references users (users_id) on delete cascade
);

CREATE TABLE IF NOT EXISTS comments
(
    comments_id bigserial not null primary key,
    body        text      not null,
    posts_id    bigint references posts (posts_id) on delete cascade,
    users_id    bigint references users (users_id) on delete cascade
);

CREATE TABLE IF NOT EXISTS payments
(
    payments_id bigserial not null primary key,
    amount      integer   not null,
    date        date      not null,
    creator_id  bigint references creator_profile (creator_id) on delete cascade,
    users_id    bigint references users (users_id) on delete cascade
)

\disconnect