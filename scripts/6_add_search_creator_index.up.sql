CREATE EXTENSION rum;

CREATE TABLE search_creators
(
    ID          bigint   not null references creator_profile (creator_id) on delete cascade,
    description tsvector not null,
    nickname    tsvector not null
);


CREATE FUNCTION upd_creators_search_creators()
    RETURNS trigger AS
$$
BEGIN
    UPDATE search_creators
    SET description = to_tsvector(NEW.description)
    WHERE ID = NEW.creator_id;
    RETURN NEW;
END;
$$
    LANGUAGE plpgsql;

CREATE TRIGGER update_creator_search
    AFTER UPDATE
    ON creator_profile
    FOR EACH ROW
EXECUTE FUNCTION upd_creators_search_creators();


CREATE FUNCTION insert_search_creators()
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


CREATE FUNCTION delete_search_creators()
    RETURNS trigger AS
$$
BEGIN
    DELETE FROM search_creators WHERE ID = NEW.creator_id;
    RETURN NEW;
END;
$$
    LANGUAGE plpgsql;

CREATE TRIGGER delete_search
    AFTER DELETE
    ON creator_profile
    FOR EACH ROW
EXECUTE FUNCTION delete_search_creators();


CREATE OR REPLACE PROCEDURE update_search_creators()
AS
$$
DECLARE
    Id  bigint;
    Des text;
BEGIN
    FOR Id, Des IN SELECT DISTINCT creator_id, description FROM creator_profile
        LOOP
            INSERT INTO search_creators(ID, nickname, description)
            VALUES (Id, to_tsvector((SELECT nickname FROM users WHERE users_id = Id LIMIT 1)), to_tsvector(Des));
        END LOOP;
END;
$$
    LANGUAGE plpgsql;

CALL update_search_creators();

CREATE INDEX idx_rum_search_creators
    ON search_creators
        USING RUM (description rum_tsvector_ops, nickname rum_tsvector_ops);