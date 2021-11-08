ALTER TABLE parents_awards
    ADD COLUMN id bigserial not null primary key;

INSERT INTO parents_awards (awards_id, parent_id)
SELECT child.awards_id, aw.awards_id
FROM awards as aw
         JOIN awards as child on aw.creator_id = child.creator_id
WHERE aw.price > child.price
  AND aw.creator_id = child.creator_id;