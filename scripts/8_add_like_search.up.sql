CREATE EXTENSION hunspell_ru_ru;

DROP INDEX idx_rum_search_creators;

ALTER TABLE search_creators
    DROP COLUMN nickname,
    ADD COLUMN description_ru tsvector default to_tsvector('') not null;

ALTER TABLE search_creators RENAME COLUMN description TO description_en;

CREATE OR REPLACE FUNCTION upd_creators_search_creators()
    RETURNS trigger AS
$$
BEGIN
    UPDATE search_creators
    SET description_en = to_tsvector('english', NEW.description),
        description_ru = to_tsvector('russian_hunspell', NEW.description)
    WHERE ID = NEW.creator_id;
    RETURN NEW;
END;
$$
    LANGUAGE plpgsql;

DROP TRIGGER update_creator_search ON creator_profile;

CREATE TRIGGER update_creator_search
    AFTER UPDATE
    ON creator_profile
    FOR EACH ROW
EXECUTE FUNCTION upd_creators_search_creators();

CREATE OR REPLACE FUNCTION insert_search_creators()
    RETURNS trigger AS
$$
BEGIN
    INSERT INTO search_creators (ID, description_en, description_ru)
    VALUES (NEW.creator_id, to_tsvector('english', NEW.description), to_tsvector('russian_hunspell', NEW.description));
    RETURN NEW;
END
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

CREATE OR REPLACE PROCEDURE add_rus_search_creators()
AS
$$
DECLARE
    IdCreator  bigint;
    Des text;
BEGIN
    FOR IdCreator, Des IN SELECT DISTINCT creator_id, description FROM creator_profile
        LOOP
            UPDATE search_creators SET description_en = to_tsvector('english', Des),
                                       description_ru = to_tsvector('russian_hunspell', Des)
            WHERE id = IdCreator;
        END LOOP;
END
$$
    LANGUAGE plpgsql;

CALL add_rus_search_creators();

CREATE INDEX idx_rum_search_creators
    ON search_creators
        USING RUM (description_en rum_tsvector_ops, description_ru rum_tsvector_ops);

CREATE INDEX idx_like_search_nick_creators
    on users (nickname) include (users_id);
