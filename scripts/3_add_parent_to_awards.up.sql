CREATE TABLE IF NOT EXISTS parents_awards
(
    id        bigserial not null primary key,
    awards_id bigint    not null references awards (awards_id) on delete cascade,
    parent_id bigint    not null references awards (awards_id) on delete cascade
);

INSERT INTO parents_awards (awards_id, parent_id)
SELECT child.awards_id, aw.awards_id
FROM awards as aw
         JOIN awards as child on aw.creator_id = child.creator_id
WHERE aw.price > child.price
  AND aw.creator_id = child.creator_id;