DROP INDEX idx_rum_search_creators;

ALTER TABLE search_creators
    DROP COLUMN nickname;


CREATE OR REPLACE FUNCTION insert_search_creators()
    RETURNS trigger AS
$$
BEGIN
    INSERT INTO search_creators (ID, description)
    VALUES (NEW.creator_id, to_tsvector(NEW.description));
    RETURN NEW;
END;
$$
    LANGUAGE plpgsql;

DROP TRIGGER insert_search ON creator_profile;
CREATE TRIGGER insert_search
    AFTER INSERT
    ON creator_profile
    FOR EACH ROW
EXECUTE FUNCTION insert_search_creators();

DROP TRIGGER update_user_search ON users;
DROP FUNCTION upd_user_search_creators();


CREATE INDEX idx_rum_search_creators
    ON search_creators
        USING RUM (description rum_tsvector_ops);

CREATE INDEX idx_like_search_nick_creators
    on users (nickname) include (users_id);