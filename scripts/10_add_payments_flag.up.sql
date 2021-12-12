alter table payments
    add column status   bool not null default false,
    add column pay_token text;
