alter table payments
    add column awards_id bigint references awards (awards_id) on delete cascade;

update payments
set awards_id = (select distinct awards_id
                 from awards
                 where payments.amount = awards.price
                   and payments.creator_id = awards.creator_id);

delete
from payments
where awards_id is null;

alter table payments
    alter column awards_id set not null;



