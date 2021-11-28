ALTER TABLE search_creators
    ADD COLUMN nickname tsvector default to_tsvector('') not null;


CREATE OR REPLACE FUNCTION insert_search_creators()
    RETURNS trigger AS
$$
BEGIN
    INSERT INTO search_creators (ID, description, nickname)
    VALUES (NEW.creator_id, to_tsvector(NEW.description),
            to_tsvector((SELECT nickname FROM users WHERE users_id = NEW.creator_id)));
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

CREATE OR REPLACE FUNCTION upd_user_search_creators()
    RETURNS trigger AS
$$
BEGIN
    IF EXISTS(SELECT FROM creator_profile WHERE creator_id = NEW.users_id) THEN
        UPDATE search_creators
        SET nickname = to_tsvector(NEW.nickname)
        WHERE ID = NEW.users_id;
    END IF;
    RETURN NEW;
END;
$$
    LANGUAGE plpgsql;

CREATE TRIGGER update_user_search
    AFTER UPDATE
    ON users
    FOR EACH ROW
EXECUTE FUNCTION upd_user_search_creators();

DROP INDEX idx_rum_search_creators;

CREATE INDEX idx_rum_search_creators
    ON search_creators
        USING RUM (description rum_tsvector_ops, nickname rum_tsvector_ops);

CREATE OR REPLACE PROCEDURE set_nick_search_creators()
AS
$$
DECLARE
    IdUser  bigint;
    Des text;
BEGIN
    FOR IdUser IN SELECT DISTINCT id FROM search_creators
        LOOP
            UPDATE search_creators SET nickname = to_tsvector((SELECT nickname FROM users WHERE users_id = Id LIMIT 1))
            WHERE search_creators.id = IdUser;
        END LOOP;
END;
$$
    LANGUAGE plpgsql;

CALL set_nick_search_creators();

DROP PROCEDURE set_nick_search_creators();

DROP INDEX idx_like_search_nick_creators;