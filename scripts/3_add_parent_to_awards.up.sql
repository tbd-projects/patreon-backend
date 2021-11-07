CREATE TABLE IF NOT EXISTS parents_awards
(
    awards_id bigint not null references awards (awards_id) on delete cascade,
    parent_id bigint not null references awards (awards_id) on delete cascade
);