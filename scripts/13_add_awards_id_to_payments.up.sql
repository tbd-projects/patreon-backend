alter table payments
    add column awards_id bigint not null references awards (awards_id) on delete cascade;