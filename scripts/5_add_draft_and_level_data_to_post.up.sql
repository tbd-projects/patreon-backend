ALTER TABLE posts
    ADD COLUMN is_draft bool default false not null;

ALTER TABLE posts_data
    ADD COLUMN level bigint default -1 not null;

CREATE PROCEDURE update_level_post_data()
AS $$
DECLARE
    Id bigint;
BEGIN
    CREATE SEQUENCE data_level;
    FOR Id IN SELECT DISTINCT post_id FROM posts_data LOOP
            PERFORM setval('data_level', 1, false);
            UPDATE posts_data SET level = nextval('data_level') WHERE post_id = Id;
        END LOOP;
    DROP SEQUENCE data_level;
END;
$$
LANGUAGE plpgsql;
CALL update_level_post_data()